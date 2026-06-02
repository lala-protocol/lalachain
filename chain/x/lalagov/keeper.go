package lalagov

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ltypes "github.com/lala-protocol/lalachain/chain/types"
)

type ParameterBounds struct {
	Min int64
	Max int64
}

type Config struct {
	Quorum             float64
	Approval           float64
	VotingPeriodEpochs int64
	AllowedParameters  map[string]ParameterBounds
}

func DefaultConfig() Config {
	return Config{
		Quorum:             0.66,
		Approval:           0.51,
		VotingPeriodEpochs: 1,
		AllowedParameters: map[string]ParameterBounds{
			"block_gas_limit":      {Min: 10_000_000, Max: 30_000_000},
			"base_fee_per_gas":     {Min: 100_000_000, Max: 10_000_000_000},
			"target_block_time_ms": {Min: 1_000, Max: 20_000},
		},
	}
}

func (c Config) withDefaults() Config {
	def := DefaultConfig()
	if c.Quorum == 0 {
		c.Quorum = def.Quorum
	}
	if c.Approval == 0 {
		c.Approval = def.Approval
	}
	if c.VotingPeriodEpochs == 0 {
		c.VotingPeriodEpochs = def.VotingPeriodEpochs
	}
	if len(c.AllowedParameters) == 0 {
		c.AllowedParameters = def.AllowedParameters
	}
	return c
}

// MaxPendingProposals is the upper bound on proposals awaiting a vote.
const MaxPendingProposals = 100

// MaxStakeValue is the safety ceiling for individual validator stake.
const MaxStakeValue = 1e18

type proposalRecord struct {
	proposal ltypes.AIProposal
	votes    map[string]ltypes.Vote
}

// Keeper implements minimal governance flow for AI-originated proposals.
// When storeKey is non-nil, state is persisted to the KV store.
type Keeper struct {
	mu            sync.RWMutex
	storeKey      storetypes.StoreKey
	config        Config
	whitelisted   map[string]struct{}
	advisoryKeys  map[string][]byte // advisory key -> HMAC secret
	pending       map[int64]*proposalRecord
	history       []ltypes.ResolvedProposal
	scheduled     map[int64][]ltypes.ParameterChange
	registeredIDs map[int64]struct{}
}

// NewKeeper creates a governance keeper.
// advisoryKeys maps each advisory key name to its HMAC-SHA256 secret for signature verification.
func NewKeeper(storeKey storetypes.StoreKey, cfg Config, advisoryKeys map[string][]byte) *Keeper {
	whitelisted := make(map[string]struct{}, len(advisoryKeys))
	for key := range advisoryKeys {
		whitelisted[key] = struct{}{}
	}
	return &Keeper{
		storeKey:      storeKey,
		config:        cfg.withDefaults(),
		whitelisted:   whitelisted,
		advisoryKeys:  advisoryKeys,
		pending:       map[int64]*proposalRecord{},
		scheduled:     map[int64][]ltypes.ParameterChange{},
		registeredIDs: map[int64]struct{}{},
	}
}

// govState is the JSON-serializable state for persistence.
type govState struct {
	Pending       map[int64]*proposalRecordJSON `json:"pending"`
	History       []ltypes.ResolvedProposal     `json:"history"`
	Scheduled     map[int64][]ltypes.ParameterChange `json:"scheduled"`
	RegisteredIDs []int64                       `json:"registered_ids"`
}

type proposalRecordJSON struct {
	Proposal ltypes.AIProposal      `json:"proposal"`
	Votes    map[string]ltypes.Vote `json:"votes"`
}

// LoadState loads persisted state from the KV store into memory.
func (k *Keeper) LoadState(ctx sdk.Context) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.storeKey == nil {
		return
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("state"))
	if bz == nil {
		return
	}
	var state govState
	if err := json.Unmarshal(bz, &state); err != nil {
		ctx.Logger().Error("lalagov: failed to unmarshal state", "error", err)
		return
	}
	k.pending = make(map[int64]*proposalRecord, len(state.Pending))
	for id, rec := range state.Pending {
		k.pending[id] = &proposalRecord{
			proposal: rec.Proposal,
			votes:    rec.Votes,
		}
	}
	k.history = state.History
	k.scheduled = state.Scheduled
	if k.scheduled == nil {
		k.scheduled = map[int64][]ltypes.ParameterChange{}
	}
	k.registeredIDs = make(map[int64]struct{}, len(state.RegisteredIDs))
	for _, id := range state.RegisteredIDs {
		k.registeredIDs[id] = struct{}{}
	}
}

// SaveState persists the current in-memory state to the KV store.
func (k *Keeper) SaveState(ctx sdk.Context) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	if k.storeKey == nil {
		return
	}
	store := ctx.KVStore(k.storeKey)
	pending := make(map[int64]*proposalRecordJSON, len(k.pending))
	for id, rec := range k.pending {
		pending[id] = &proposalRecordJSON{
			Proposal: rec.proposal,
			Votes:    rec.votes,
		}
	}
	ids := make([]int64, 0, len(k.registeredIDs))
	for id := range k.registeredIDs {
		ids = append(ids, id)
	}
	state := govState{
		Pending:       pending,
		History:       k.history,
		Scheduled:     k.scheduled,
		RegisteredIDs: ids,
	}
	bz, err := json.Marshal(state)
	if err != nil {
		return
	}
	store.Set([]byte("state"), bz)
}

func (k *Keeper) RegisterAIProposal(proposal ltypes.AIProposal) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if len(k.pending) >= MaxPendingProposals {
		return fmt.Errorf("max pending proposals (%d) reached", MaxPendingProposals)
	}
	if _, exists := k.registeredIDs[proposal.ProposalID]; exists {
		return fmt.Errorf("proposal id %d already exists", proposal.ProposalID)
	}
	if _, ok := k.whitelisted[proposal.AdvisoryKey]; !ok {
		return fmt.Errorf("proposal key %q is not whitelisted", proposal.AdvisoryKey)
	}
	if err := k.verifySignature(proposal); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	if err := k.validateProposal(proposal); err != nil {
		return err
	}

	k.pending[proposal.ProposalID] = &proposalRecord{
		proposal: proposal,
		votes:    map[string]ltypes.Vote{},
	}
	k.registeredIDs[proposal.ProposalID] = struct{}{}
	return nil
}

func (k *Keeper) validateProposal(proposal ltypes.AIProposal) error {
	if proposal.ProposedValue == proposal.CurrentValue {
		return errors.New("proposal has no value change")
	}
	bounds, ok := k.config.AllowedParameters[proposal.Parameter]
	if !ok {
		return fmt.Errorf("unsupported parameter %q", proposal.Parameter)
	}
	if proposal.ProposedValue < bounds.Min || proposal.ProposedValue > bounds.Max {
		return fmt.Errorf(
			"proposed value %d out of range for %s [%d,%d]",
			proposal.ProposedValue,
			proposal.Parameter,
			bounds.Min,
			bounds.Max,
		)
	}
	if proposal.ActivationEpoch <= proposal.EpochSubmitted {
		return errors.New("activation epoch must be after submission epoch")
	}
	// Prevent scheduling too far in the future (max 1000 epochs ahead).
	if proposal.ActivationEpoch > proposal.EpochSubmitted+1000 {
		return errors.New("activation epoch too far in the future")
	}
	if proposal.AdvisorySignature == "" {
		return errors.New("proposal must include advisory signature")
	}
	return nil
}

func (k *Keeper) CastVote(proposalID int64, validator ltypes.Validator, approve bool) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	record, exists := k.pending[proposalID]
	if !exists {
		return fmt.Errorf("unknown proposal id %d", proposalID)
	}
	if validator.Stake <= 0 || validator.Stake > MaxStakeValue ||
		math.IsNaN(validator.Stake) || math.IsInf(validator.Stake, 0) {
		return fmt.Errorf("validator %s has invalid stake %.4f", validator.Address, validator.Stake)
	}
	if _, voted := record.votes[validator.Address]; voted {
		return fmt.Errorf("validator %s already voted on proposal %d", validator.Address, proposalID)
	}

	record.votes[validator.Address] = ltypes.Vote{
		Validator: validator.Address,
		Stake:     validator.Stake,
		Approve:   approve,
	}
	return nil
}

func (k *Keeper) PendingProposals() []ltypes.AIProposal {
	k.mu.RLock()
	defer k.mu.RUnlock()

	ids := make([]int64, 0, len(k.pending))
	for id := range k.pending {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	proposals := make([]ltypes.AIProposal, 0, len(ids))
	for _, id := range ids {
		proposals = append(proposals, k.pending[id].proposal)
	}
	return proposals
}

func (k *Keeper) EndEpoch(currentEpoch int64, totalStake float64) []ltypes.ResolvedProposal {
	k.mu.Lock()
	defer k.mu.Unlock()

	if totalStake <= 0 || math.IsNaN(totalStake) || math.IsInf(totalStake, 0) {
		return nil
	}

	ids := make([]int64, 0, len(k.pending))
	for id := range k.pending {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	resolved := make([]ltypes.ResolvedProposal, 0)
	for _, id := range ids {
		record := k.pending[id]
		if currentEpoch < record.proposal.EpochSubmitted+k.config.VotingPeriodEpochs {
			continue
		}

		votesApprove := 0.0
		votesReject := 0.0
		for _, vote := range record.votes {
			if vote.Approve {
				votesApprove += vote.Stake
			} else {
				votesReject += vote.Stake
			}
		}
		totalVoted := votesApprove + votesReject

		quorumMet := (totalVoted / totalStake) >= k.config.Quorum
		approvalMet := totalVoted > 0 && ((votesApprove / totalVoted) >= k.config.Approval)

		outcome := "failed"
		if quorumMet && approvalMet {
			outcome = "passed"
			k.scheduled[record.proposal.ActivationEpoch] = append(
				k.scheduled[record.proposal.ActivationEpoch],
				ltypes.ParameterChange{
					Parameter: record.proposal.Parameter,
					Value:     record.proposal.ProposedValue,
				},
			)
		}

		r := ltypes.ResolvedProposal{
			Proposal:     record.proposal,
			VotesApprove: votesApprove,
			VotesReject:  votesReject,
			Outcome:      outcome,
		}
		k.history = append(k.history, r)
		resolved = append(resolved, r)
		delete(k.pending, id)
	}
	return resolved
}

func (k *Keeper) ApplyScheduledActivations(epoch int64, params *ltypes.NetworkParams) []ltypes.ParameterChange {
	k.mu.Lock()
	defer k.mu.Unlock()

	changes, ok := k.scheduled[epoch]
	if !ok || len(changes) == 0 {
		return nil
	}
	delete(k.scheduled, epoch)

	applied := make([]ltypes.ParameterChange, 0, len(changes))
	for _, change := range changes {
		switch change.Parameter {
		case "block_gas_limit":
			params.BlockGasLimit = change.Value
		case "base_fee_per_gas":
			params.BaseFeePerGas = change.Value
		case "target_block_time_ms":
			params.TargetBlockTimeMS = change.Value
		default:
			continue
		}
		applied = append(applied, change)
	}
	return applied
}

func (k *Keeper) History() []ltypes.ResolvedProposal {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return append([]ltypes.ResolvedProposal(nil), k.history...)
}

// GetConfig returns a copy of the governance configuration.
func (k *Keeper) GetConfig() Config {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.config
}

// verifySignature checks the HMAC-SHA256 signature on an AI proposal.
func (k *Keeper) verifySignature(proposal ltypes.AIProposal) error {
	secret, ok := k.advisoryKeys[proposal.AdvisoryKey]
	if !ok {
		return errors.New("unknown advisory key")
	}
	payload := fmt.Sprintf(
		"%s|%d|%s|%d|%d|%d",
		proposal.AdvisoryKey,
		proposal.EpochSubmitted,
		proposal.Parameter,
		proposal.CurrentValue,
		proposal.ProposedValue,
		proposal.ActivationEpoch,
	)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(proposal.AdvisorySignature), []byte(expected)) {
		return errors.New("invalid advisory signature")
	}
	return nil
}
