package app

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"

	ltypes "github.com/lala-protocol/whitepaper/prototype/lalachain/types"
	"github.com/lala-protocol/whitepaper/prototype/lalachain/x/aiadvisor"
	"github.com/lala-protocol/whitepaper/prototype/lalachain/x/gov"
	"github.com/lala-protocol/whitepaper/prototype/lalachain/x/telemetry"
)

const advisoryKey = "lala-ai-advisor-phase1"

// Prototype wires deterministic telemetry, rule-based AI advisory, and governance.
type Prototype struct {
	params     ltypes.NetworkParams
	validators []ltypes.Validator

	telemetry  telemetry.Keeper
	ai         *aiadvisor.Keeper
	governance *gov.Keeper

	rng        *rand.Rand
	totalStake float64
	kpiHistory []ltypes.EpochKPIs
}

func NewPrototype(seed int64) *Prototype {
	validators := []ltypes.Validator{
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

	totalStake := 0.0
	for _, validator := range validators {
		totalStake += validator.Stake
	}

	return &Prototype{
		params:     ltypes.DefaultNetworkParams(),
		validators: validators,
		telemetry:  telemetry.NewKeeper(),
		ai:         aiadvisor.NewKeeper(advisoryKey, aiadvisor.DefaultConfig()),
		governance: gov.NewKeeper(gov.DefaultConfig(), []string{advisoryKey}),
		rng:        rand.New(rand.NewSource(seed)),
		totalStake: totalStake,
	}
}

func (p *Prototype) Run(epochs int, out io.Writer) error {
	if epochs <= 0 {
		return errors.New("epochs must be > 0")
	}

	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintln(out, " LalaChain Phase 1 Prototype")
	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintf(out, "Initial params: %+v\n\n", p.params)

	for epoch := int64(1); epoch <= int64(epochs); epoch++ {
		fmt.Fprintf(out, "-- Epoch %03d --------------------------------------------------\n", epoch)

		applied := p.governance.ApplyScheduledActivations(epoch, &p.params)
		for _, change := range applied {
			fmt.Fprintf(
				out,
				"  [Activation] %s = %d\n",
				change.Parameter,
				change.Value,
			)
		}

		demand := p.demandForEpoch(epoch)
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
			"  [KPI] util=%.1f%% fee=%.3f gwei block_time=%.0fms demand=%.2f\n",
			kpis.AvgUtilization*100,
			float64(kpis.AvgBaseFee)/1e9,
			kpis.AvgBlockTimeMS,
			demand,
		)

		if proposal := p.ai.AnalyzeAndPropose(kpis, p.params); proposal != nil {
			if err := p.governance.RegisterAIProposal(*proposal); err != nil {
				return err
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
	}

	p.printSummary(out, epochs)
	return nil
}

func (p *Prototype) Params() ltypes.NetworkParams {
	return p.params
}

func (p *Prototype) History() []ltypes.ResolvedProposal {
	return p.governance.History()
}

func (p *Prototype) KPIHistory() []ltypes.EpochKPIs {
	return append([]ltypes.EpochKPIs(nil), p.kpiHistory...)
}

func (p *Prototype) demandForEpoch(epoch int64) float64 {
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
