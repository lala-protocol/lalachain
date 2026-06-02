package app

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"

	ltypes "github.com/lala-protocol/lalachain/chain/types"
	"github.com/lala-protocol/lalachain/chain/x/aiadvisor"
	"github.com/lala-protocol/lalachain/chain/x/gov"
	"github.com/lala-protocol/lalachain/chain/x/telemetry"
)

const (
	advisoryKeyPhase1 = "lala-ai-advisor-phase1"
	advisoryKeyPhase2 = "lala-ai-advisor-phase2"
)

type Phase string

const (
	Phase1 Phase = "phase1"
	Phase2 Phase = "phase2"
)

type phase2Scenario struct {
	Name       string
	StartEpoch int64
	EndEpoch   int64
	BaseDemand float64
	Volatility float64
}

type activationAuditKey struct {
	Epoch     int64
	Parameter string
	Value     int64
}

type phase2Audit struct {
	registeredCount     int
	resolvedCount       int
	passedCount         int
	failedCount         int
	appliedCount        int
	expectedActivations map[activationAuditKey]int
	anomalies           []string
}

func newPhase2Audit() *phase2Audit {
	return &phase2Audit{
		expectedActivations: map[activationAuditKey]int{},
	}
}

func (a *phase2Audit) recordRegistered() {
	a.registeredCount++
}

func (a *phase2Audit) recordResolved(resolved ltypes.ResolvedProposal) {
	a.resolvedCount++
	if resolved.Outcome == "passed" {
		a.passedCount++
		key := activationAuditKey{
			Epoch:     resolved.Proposal.ActivationEpoch,
			Parameter: resolved.Proposal.Parameter,
			Value:     resolved.Proposal.ProposedValue,
		}
		a.expectedActivations[key]++
		return
	}
	a.failedCount++
}

func (a *phase2Audit) recordApplied(epoch int64, change ltypes.ParameterChange) {
	a.appliedCount++
	key := activationAuditKey{
		Epoch:     epoch,
		Parameter: change.Parameter,
		Value:     change.Value,
	}
	remaining := a.expectedActivations[key]
	if remaining == 0 {
		a.anomalies = append(
			a.anomalies,
			fmt.Sprintf("unexpected activation at epoch %d for %s=%d", epoch, change.Parameter, change.Value),
		)
		return
	}
	if remaining == 1 {
		delete(a.expectedActivations, key)
		return
	}
	a.expectedActivations[key] = remaining - 1
}

func (a *phase2Audit) finalize(finalEpoch int64) {
	for key, count := range a.expectedActivations {
		if key.Epoch > finalEpoch {
			continue
		}
		a.anomalies = append(
			a.anomalies,
			fmt.Sprintf("missing activation by epoch %d for %s=%d (%d pending)", key.Epoch, key.Parameter, key.Value, count),
		)
	}
}

func (a *phase2Audit) isClean() bool {
	return len(a.anomalies) == 0
}

// EpochHook is called by the prototype at each epoch boundary when set.
// This allows the cosmos app to interleave module lifecycle calls.
type EpochHook interface {
	BeforeEpoch(epoch int64)
	AfterEpoch(epoch int64)
}

// Prototype wires deterministic telemetry, rule-based AI advisory, and governance.
type Prototype struct {
	phase      Phase
	params     ltypes.NetworkParams
	validators []ltypes.Validator

	telemetry  telemetry.Keeper
	ai         *aiadvisor.Keeper
	governance *gov.Keeper

	rng        *rand.Rand
	totalStake float64
	kpiHistory []ltypes.EpochKPIs

	phase2Scenarios []phase2Scenario
	phase2Audit     *phase2Audit

	epochHook EpochHook
}

// SetEpochHook sets an optional hook that is called before/after each epoch.
func (p *Prototype) SetEpochHook(hook EpochHook) {
	p.epochHook = hook
}

func NewPrototype(seed int64) *Prototype {
	return newPrototype(seed, Phase1, 10)
}

func NewPhase2Prototype(seed int64, validatorCount int) *Prototype {
	if validatorCount < 4 {
		validatorCount = 4
	}
	if validatorCount > 10 {
		validatorCount = 10
	}
	return newPrototype(seed, Phase2, validatorCount)
}

func newPrototype(seed int64, phase Phase, validatorCount int) *Prototype {
	validators := validatorSet(validatorCount)
	advisoryKey := advisoryKeyPhase1
	if phase == Phase2 {
		advisoryKey = advisoryKeyPhase2
	}

	totalStake := 0.0
	for _, validator := range validators {
		totalStake += validator.Stake
	}

	p := &Prototype{
		phase:      phase,
		params:     ltypes.DefaultNetworkParams(),
		validators: validators,
		telemetry:  telemetry.NewKeeper(nil),
		ai:         aiadvisor.NewKeeper(nil, advisoryKey, []byte("prototype-secret"), aiadvisor.DefaultConfig()),
		governance: gov.NewKeeper(nil, gov.DefaultConfig(), []string{advisoryKey}),
		rng:        rand.New(rand.NewSource(seed)),
		totalStake: totalStake,
	}

	if phase == Phase2 {
		p.phase2Scenarios = defaultPhase2Scenarios()
		p.phase2Audit = newPhase2Audit()
	}

	return p
}

func validatorSet(count int) []ltypes.Validator {
	base := []ltypes.Validator{
		{Address: "val01", Stake: 0.18},
		{Address: "val02", Stake: 0.15},
		{Address: "val03", Stake: 0.13},
		{Address: "val04", Stake: 0.11},
		{Address: "val05", Stake: 0.10},
		{Address: "val06", Stake: 0.09},
		{Address: "val07", Stake: 0.08},
		{Address: "val08", Stake: 0.07, Conservative: true},
		{Address: "val09", Stake: 0.05, Conservative: true},
		{Address: "val10", Stake: 0.04, Conservative: true},
	}
	if count <= 0 || count >= len(base) {
		count = len(base)
	}

	validators := append([]ltypes.Validator(nil), base[:count]...)
	totalStake := 0.0
	for _, validator := range validators {
		totalStake += validator.Stake
	}
	if totalStake == 0 {
		return validators
	}

	for i := range validators {
		validators[i].Stake = validators[i].Stake / totalStake
	}

	return validators
}

func defaultPhase2Scenarios() []phase2Scenario {
	return []phase2Scenario{
		{Name: "low-load warmup", StartEpoch: 1, EndEpoch: 8, BaseDemand: 0.32, Volatility: 0.06},
		{Name: "surge stress", StartEpoch: 9, EndEpoch: 16, BaseDemand: 0.94, Volatility: 0.08},
		{Name: "oscillation stress", StartEpoch: 17, EndEpoch: 24, BaseDemand: 0.66, Volatility: 0.20},
		{Name: "recovery stabilization", StartEpoch: 25, EndEpoch: 1 << 62, BaseDemand: 0.56, Volatility: 0.05},
	}
}

func (p *Prototype) phase2Scenario(epoch int64) phase2Scenario {
	for _, scenario := range p.phase2Scenarios {
		if epoch >= scenario.StartEpoch && epoch <= scenario.EndEpoch {
			return scenario
		}
	}
	return phase2Scenario{Name: "fallback", StartEpoch: 1, EndEpoch: 1 << 62, BaseDemand: 0.55, Volatility: 0.05}
}

func phaseLabel(phase Phase) string {
	if phase == Phase2 {
		return "Phase 2"
	}
	return "Phase 1"
}

func (p *Prototype) Run(epochs int, out io.Writer) error {
	if epochs <= 0 {
		return errors.New("epochs must be > 0")
	}

	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintf(out, " LalaChain %s Prototype\n", phaseLabel(p.phase))
	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintf(out, "Mode: %s  Validators: %d\n", p.phase, len(p.validators))
	fmt.Fprintf(out, "Initial params: %+v\n\n", p.params)

	for epoch := int64(1); epoch <= int64(epochs); epoch++ {
		if p.epochHook != nil {
			p.epochHook.BeforeEpoch(epoch)
		}

		fmt.Fprintf(out, "-- Epoch %03d --------------------------------------------------\n", epoch)

		applied := p.governance.ApplyScheduledActivations(epoch, &p.params)
		for _, change := range applied {
			if p.phase2Audit != nil {
				p.phase2Audit.recordApplied(epoch, change)
			}
			fmt.Fprintf(
				out,
				"  [Activation] %s = %d\n",
				change.Parameter,
				change.Value,
			)
		}

		demand, scenario := p.demandForEpoch(epoch)
		blocks := p.produceEpochBlocks(demand)

		kpis, err := p.telemetry.CalculateEpochKPIs(
			epoch,
			blocks,
			p.params,
			p.validators,
			(demand-0.5)*10,
			p.syntheticSlashingEvents(),
		)
		if err != nil {
			return err
		}
		p.kpiHistory = append(p.kpiHistory, kpis)

		fmt.Fprintf(
			out,
			"  [KPI] util=%.1f%% fee=%.3f gwei block_time=%.0fms demand=%.2f scenario=%s\n",
			kpis.AvgUtilization*100,
			float64(kpis.AvgBaseFee)/1e9,
			kpis.AvgBlockTimeMS,
			demand,
			scenario,
		)

		if proposal := p.ai.AnalyzeAndPropose(kpis, p.params); proposal != nil {
			if err := p.governance.RegisterAIProposal(*proposal); err != nil {
				return err
			}
			if p.phase2Audit != nil {
				p.phase2Audit.recordRegistered()
			}
			fmt.Fprintf(
				out,
				"  [Gov] proposal #%d registered: %s %d -> %d\n",
				proposal.ProposalID,
				proposal.Parameter,
				proposal.CurrentValue,
				proposal.ProposedValue,
			)
		}

		p.castVotes(epoch)
		resolved := p.governance.EndEpoch(epoch, p.totalStake)
		for _, item := range resolved {
			if p.phase2Audit != nil {
				p.phase2Audit.recordResolved(item)
			}
			fmt.Fprintf(
				out,
				"  [Gov] proposal #%d %s (approve=%.1f%% reject=%.1f%%)\n",
				item.Proposal.ProposalID,
				strStatus(item.Outcome),
				item.VotesApprove*100,
				item.VotesReject*100,
			)
			if item.Outcome == "passed" {
				fmt.Fprintf(
					out,
					"  [Gov] activation scheduled at epoch %d\n",
					item.Proposal.ActivationEpoch,
				)
			}
		}
		fmt.Fprintln(out)

		if p.epochHook != nil {
			p.epochHook.AfterEpoch(epoch)
		}
	}

	if p.phase2Audit != nil {
		p.phase2Audit.finalize(int64(epochs))
	}

	p.printSummary(out, epochs)
	return nil
}

func (p *Prototype) Params() ltypes.NetworkParams {
	return p.params
}

func (p *Prototype) SetParams(params ltypes.NetworkParams) {
	p.params = params
}

func (p *Prototype) TelemetryKeeper() telemetry.Keeper {
	return p.telemetry
}

func (p *Prototype) AIAdvisorKeeper() *aiadvisor.Keeper {
	return p.ai
}

func (p *Prototype) GovKeeper() *gov.Keeper {
	return p.governance
}

func (p *Prototype) Validators() []ltypes.Validator {
	return p.validators
}

func (p *Prototype) TotalStake() float64 {
	return p.totalStake
}

func (p *Prototype) History() []ltypes.ResolvedProposal {
	return p.governance.History()
}

func (p *Prototype) KPIHistory() []ltypes.EpochKPIs {
	return append([]ltypes.EpochKPIs(nil), p.kpiHistory...)
}

func (p *Prototype) demandForEpoch(epoch int64) (float64, string) {
	if p.phase == Phase2 {
		return p.phase2DemandForEpoch(epoch)
	}
	return p.phase1DemandForEpoch(epoch), "baseline"
}

func (p *Prototype) phase1DemandForEpoch(epoch int64) float64 {
	var base float64
	switch {
	case epoch <= 10:
		base = 0.30
	case epoch <= 20:
		base = 0.90
	default:
		base = 0.55
	}
	return clamp(base+p.rng.NormFloat64()*0.05, 0, 1.2)
}

func (p *Prototype) phase2DemandForEpoch(epoch int64) (float64, string) {
	scenario := p.phase2Scenario(epoch)
	demand := scenario.BaseDemand + p.rng.NormFloat64()*scenario.Volatility

	if scenario.Name == "oscillation stress" {
		if epoch%2 == 0 {
			demand += 0.15
		} else {
			demand -= 0.15
		}
	}

	if epoch%15 == 0 {
		demand += 0.18
	}
	if epoch%22 == 0 {
		demand -= 0.15
	}

	return clamp(demand, 0, 1.25), scenario.Name
}

func (p *Prototype) produceEpochBlocks(demand float64) []ltypes.FinalizedBlock {
	blocks := make([]ltypes.FinalizedBlock, 0, 50)
	for i := 0; i < 50; i++ {
		blockTime := int64(math.Max(
			1_000,
			float64(p.params.TargetBlockTimeMS)+(p.rng.NormFloat64()*float64(p.params.TargetBlockTimeMS)*0.05),
		))

		utilization := clamp(demand*(1+p.rng.NormFloat64()*0.15), 0, 1)
		gasUsed := int64(float64(p.params.BlockGasLimit) * utilization)

		feeAdjustment := (utilization - 0.5) * 0.125
		baseFee := int64(math.Max(
			100_000_000,
			float64(p.params.BaseFeePerGas)*(1+feeAdjustment),
		))

		blocks = append(blocks, ltypes.FinalizedBlock{
			BlockTimeMS: blockTime,
			GasUsed:     gasUsed,
			BaseFee:     baseFee,
		})
	}
	return blocks
}

func (p *Prototype) syntheticSlashingEvents() int64 {
	var count int64
	for range p.validators {
		if p.rng.Float64() < 0.005 {
			count++
		}
	}
	return count
}

func (p *Prototype) castVotes(currentEpoch int64) {
	for _, proposal := range p.governance.PendingProposals() {
		if currentEpoch <= proposal.EpochSubmitted {
			continue
		}
		for _, validator := range p.validators {
			if p.rng.Float64() < 0.10 {
				continue
			}
			approve := p.voteForProposal(validator, proposal)
			_ = p.governance.CastVote(proposal.ProposalID, validator, approve)
		}
	}
}

func (p *Prototype) voteForProposal(validator ltypes.Validator, proposal ltypes.AIProposal) bool {
	changePct := proposal.ChangePercent()
	if validator.Conservative && changePct > 0.03 {
		return false
	}

	approvalProb := 0.90 - changePct*2
	if approvalProb < 0.30 {
		approvalProb = 0.30
	}
	return p.rng.Float64() < approvalProb
}

func (p *Prototype) printSummary(out io.Writer, epochs int) {
	history := p.governance.History()
	passed := 0
	failed := 0
	for _, item := range history {
		if item.Outcome == "passed" {
			passed++
		} else {
			failed++
		}
	}

	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintln(out, " PROTOTYPE SUMMARY")
	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintf(out, "Phase: %s\n", p.phase)
	fmt.Fprintf(out, "Validators: %d\n", len(p.validators))
	fmt.Fprintf(out, "Final params: %+v\n", p.params)
	fmt.Fprintf(out, "Epochs: %d\n", epochs)
	fmt.Fprintf(out, "Governance outcomes: passed=%d failed=%d\n", passed, failed)

	if len(p.kpiHistory) > 0 {
		utilSum := 0.0
		feeSum := 0.0
		for _, kpi := range p.kpiHistory {
			utilSum += kpi.AvgUtilization
			feeSum += float64(kpi.AvgBaseFee) / 1e9
		}
		fmt.Fprintf(
			out,
			"KPI averages: utilization=%.1f%% base_fee=%.3f gwei\n",
			(utilSum/float64(len(p.kpiHistory)))*100,
			feeSum/float64(len(p.kpiHistory)),
		)
	}

	if p.phase == Phase2 {
		p.printPhase2GovernanceAudit(out)
		p.printPhase2StabilityFindings(out, history)
	}
}

func (p *Prototype) printPhase2GovernanceAudit(out io.Writer) {
	if p.phase2Audit == nil {
		return
	}

	fmt.Fprintln(out, "----------------------------------------------------------------------")
	fmt.Fprintln(out, " PHASE 2 GOVERNANCE AUDIT")
	fmt.Fprintln(out, "----------------------------------------------------------------------")
	fmt.Fprintf(out, "Registered proposals: %d\n", p.phase2Audit.registeredCount)
	fmt.Fprintf(out, "Resolved proposals: %d (passed=%d failed=%d)\n", p.phase2Audit.resolvedCount, p.phase2Audit.passedCount, p.phase2Audit.failedCount)
	fmt.Fprintf(out, "Applied activations: %d\n", p.phase2Audit.appliedCount)

	if p.phase2Audit.isClean() {
		fmt.Fprintln(out, "Audit status: OK (no activation path anomalies)")
		return
	}

	fmt.Fprintf(out, "Audit status: FAIL (%d anomalies)\n", len(p.phase2Audit.anomalies))
	for _, anomaly := range p.phase2Audit.anomalies {
		fmt.Fprintf(out, "  - %s\n", anomaly)
	}
}

func (p *Prototype) printPhase2StabilityFindings(out io.Writer, history []ltypes.ResolvedProposal) {
	if len(p.kpiHistory) == 0 {
		return
	}

	utilizations := make([]float64, 0, len(p.kpiHistory))
	blockTimes := make([]float64, 0, len(p.kpiHistory))
	baseFees := make([]float64, 0, len(p.kpiHistory))
	outsideTarget := 0

	for _, kpi := range p.kpiHistory {
		utilizations = append(utilizations, kpi.AvgUtilization)
		blockTimes = append(blockTimes, kpi.AvgBlockTimeMS)
		baseFees = append(baseFees, float64(kpi.AvgBaseFee)/1e9)
		if kpi.AvgUtilization < 0.40 || kpi.AvgUtilization > 0.80 {
			outsideTarget++
		}
	}

	totalOutcomes := len(history)
	passRate := 0.0
	if totalOutcomes > 0 {
		passes := 0
		for _, outcome := range history {
			if outcome.Outcome == "passed" {
				passes++
			}
		}
		passRate = float64(passes) / float64(totalOutcomes)
	}

	oscillationEvents := countOscillationEvents(history)
	verdict := "stable under synthetic stress"
	if oscillationEvents > 3 || outsideTarget > len(p.kpiHistory)/2 {
		verdict = "needs tighter control limits"
	}
	if p.phase2Audit != nil && !p.phase2Audit.isClean() {
		verdict = "governance activation path requires fixes"
	}

	fmt.Fprintln(out, "----------------------------------------------------------------------")
	fmt.Fprintln(out, " PHASE 2 STABILITY FINDINGS")
	fmt.Fprintln(out, "----------------------------------------------------------------------")
	fmt.Fprintf(
		out,
		"Utilization: avg=%.1f%% std=%.1f%% min=%.1f%% max=%.1f%% outside_target=%d/%d\n",
		mean(utilizations)*100,
		stddev(utilizations)*100,
		minFloat(utilizations)*100,
		maxFloat(utilizations)*100,
		outsideTarget,
		len(p.kpiHistory),
	)
	fmt.Fprintf(
		out,
		"Block time: avg=%.0fms std=%.0fms\n",
		mean(blockTimes),
		stddev(blockTimes),
	)
	fmt.Fprintf(
		out,
		"Base fee: avg=%.3f gwei min=%.3f max=%.3f\n",
		mean(baseFees),
		minFloat(baseFees),
		maxFloat(baseFees),
	)
	fmt.Fprintf(
		out,
		"Governance: pass_rate=%.1f%% outcomes=%d oscillation_events=%d\n",
		passRate*100,
		totalOutcomes,
		oscillationEvents,
	)
	fmt.Fprintf(out, "Stability verdict: %s\n", verdict)
}

func countOscillationEvents(history []ltypes.ResolvedProposal) int {
	lastDirectionByParam := map[string]int{}
	count := 0

	for _, item := range history {
		if item.Outcome != "passed" {
			continue
		}
		delta := item.Proposal.ProposedValue - item.Proposal.CurrentValue
		direction := 0
		if delta > 0 {
			direction = 1
		} else if delta < 0 {
			direction = -1
		}
		if direction == 0 {
			continue
		}

		if lastDirection, ok := lastDirectionByParam[item.Proposal.Parameter]; ok && lastDirection != direction {
			count++
		}
		lastDirectionByParam[item.Proposal.Parameter] = direction
	}

	return count
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func stddev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	m := mean(values)
	acc := 0.0
	for _, value := range values {
		delta := value - m
		acc += delta * delta
	}
	return math.Sqrt(acc / float64(len(values)))
}

func minFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
	}
	return min
}

func maxFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}
	return max
}

func strStatus(outcome string) string {
	if outcome == "passed" {
		return "PASSED"
	}
	return "FAILED"
}

func clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
