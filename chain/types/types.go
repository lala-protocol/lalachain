package types

import "math"

// NetworkParams are the mutable protocol values governed through on-chain voting.
type NetworkParams struct {
	BlockGasLimit     int64
	BaseFeePerGas     int64
	TargetBlockTimeMS int64
}

func DefaultNetworkParams() NetworkParams {
	return NetworkParams{
		BlockGasLimit:     15_000_000,
		BaseFeePerGas:     1_000_000_000,
		TargetBlockTimeMS: 6_000,
	}
}

// FinalizedBlock is the deterministic input used by telemetry.
type FinalizedBlock struct {
	BlockTimeMS int64
	GasUsed     int64
	BaseFee     int64
}

// EpochKPIs are computed once per epoch from finalized data.
type EpochKPIs struct {
	Epoch               int64
	AvgBlockTimeMS      float64
	BlockTimeVarianceMS float64
	AvgUtilization      float64
	AvgBaseFee          int64
	ValidatorCount      int
	TotalStakedRatio    float64
	MempoolSizeTrend    float64
	SlashingEvents      int64
}

// AIProposal mirrors MsgSubmitAIProposal payload intent from the feasibility plan.
type AIProposal struct {
	ProposalID        int64
	EpochSubmitted    int64
	Parameter         string
	CurrentValue      int64
	ProposedValue     int64
	Rationale         string
	ActivationEpoch   int64
	AdvisoryKey       string
	AdvisorySignature string
}

func (p AIProposal) ChangePercent() float64 {
	if p.CurrentValue == 0 {
		return 0
	}
	delta := float64(p.ProposedValue-p.CurrentValue) / float64(p.CurrentValue)
	return math.Abs(delta)
}

type Validator struct {
	Address      string
	Stake        float64
	Conservative bool
}

type Vote struct {
	Validator string
	Stake     float64
	Approve   bool
}

type ParameterChange struct {
	Parameter string
	Value     int64
}

type ResolvedProposal struct {
	Proposal     AIProposal
	VotesApprove float64
	VotesReject  float64
	Outcome      string // "passed" | "failed"
}
