package lalagov

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	ltypes "github.com/lala-protocol/lalachain/chain/types"
)

// testSecret is used by tests to sign proposals with HMAC.
var testSecret = []byte("test-secret-key")

func signTestProposal(p ltypes.AIProposal) string {
	payload := fmt.Sprintf(
		"%s|%d|%s|%d|%d|%d",
		p.AdvisoryKey, p.EpochSubmitted, p.Parameter,
		p.CurrentValue, p.ProposedValue, p.ActivationEpoch,
	)
	mac := hmac.New(sha256.New, testSecret)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestGovernancePassAndActivation(t *testing.T) {
	keeper := NewKeeper(nil, DefaultConfig(), map[string][]byte{"ai-key": testSecret})
	proposal := ltypes.AIProposal{
		ProposalID:      1,
		EpochSubmitted:  5,
		Parameter:       "block_gas_limit",
		CurrentValue:    15_000_000,
		ProposedValue:   15_750_000,
		Rationale:       "test",
		ActivationEpoch: 7,
		AdvisoryKey:     "ai-key",
	}
	proposal.AdvisorySignature = signTestProposal(proposal)

	if err := keeper.RegisterAIProposal(proposal); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := keeper.CastVote(1, ltypes.Validator{Address: "val1", Stake: 0.7}, true); err != nil {
		t.Fatalf("vote failed: %v", err)
	}
	if err := keeper.CastVote(1, ltypes.Validator{Address: "val2", Stake: 0.3}, false); err != nil {
		t.Fatalf("vote failed: %v", err)
	}

	resolved := keeper.EndEpoch(6, 1.0)
	if len(resolved) != 1 {
		t.Fatalf("expected 1 resolved proposal, got %d", len(resolved))
	}
	if resolved[0].Outcome != "passed" {
		t.Fatalf("expected passed outcome, got %s", resolved[0].Outcome)
	}

	params := ltypes.DefaultNetworkParams()
	applied := keeper.ApplyScheduledActivations(7, &params)
	if len(applied) != 1 {
		t.Fatalf("expected 1 applied activation, got %d", len(applied))
	}
	if params.BlockGasLimit != 15_750_000 {
		t.Fatalf("expected block gas limit activation, got %d", params.BlockGasLimit)
	}
}

func TestRejectsUnwhitelistedProposalKey(t *testing.T) {
	keeper := NewKeeper(nil, DefaultConfig(), map[string][]byte{"trusted-key": testSecret})
	proposal := ltypes.AIProposal{
		ProposalID:      1,
		EpochSubmitted:  1,
		Parameter:       "base_fee_per_gas",
		CurrentValue:    1_000_000_000,
		ProposedValue:   1_100_000_000,
		Rationale:       "test",
		ActivationEpoch: 3,
		AdvisoryKey:     "untrusted-key",
	}
	proposal.AdvisorySignature = signTestProposal(proposal)

	if err := keeper.RegisterAIProposal(proposal); err == nil {
		t.Fatal("expected unwhitelisted key to be rejected")
	}
}
