package aiadvisor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	ltypes "github.com/lala-protocol/whitepaper/prototype/lalachain/types"
)

const (
	ParamBlockGasLimit     = "block_gas_limit"
	ParamBaseFeePerGas     = "base_fee_per_gas"
	ParamTargetBlockTimeMS = "target_block_time_ms"
)

type Config struct {
	MinFeeTarget      int64
	MaxFeeTarget      int64
	LowUtilThreshold  float64
	HighUtilThreshold float64
	MinBlockGasLimit  int64
	MaxBlockGasLimit  int64
	MinBaseFee        int64
	MaxBaseFee        int64
}

func DefaultConfig() Config {
	return Config{
		MinFeeTarget:      800_000_000,   // 0.8 Gwei
		MaxFeeTarget:      5_000_000_000, // 5.0 Gwei
		LowUtilThreshold:  0.40,
		HighUtilThreshold: 0.80,
		MinBlockGasLimit:  10_000_000,
		MaxBlockGasLimit:  30_000_000,
		MinBaseFee:        100_000_000,
		MaxBaseFee:        10_000_000_000,
	}
}

func (c Config) withDefaults() Config {
	def := DefaultConfig()
	if c.MinFeeTarget == 0 {
		c.MinFeeTarget = def.MinFeeTarget
	}
	if c.MaxFeeTarget == 0 {
		c.MaxFeeTarget = def.MaxFeeTarget
	}
	if c.LowUtilThreshold == 0 {
		c.LowUtilThreshold = def.LowUtilThreshold
	}
	if c.HighUtilThreshold == 0 {
		c.HighUtilThreshold = def.HighUtilThreshold
	}
	if c.MinBlockGasLimit == 0 {
		c.MinBlockGasLimit = def.MinBlockGasLimit
	}
	if c.MaxBlockGasLimit == 0 {
		c.MaxBlockGasLimit = def.MaxBlockGasLimit
	}
	if c.MinBaseFee == 0 {
		c.MinBaseFee = def.MinBaseFee
	}
	if c.MaxBaseFee == 0 {
		c.MaxBaseFee = def.MaxBaseFee
	}
	return c
}

// Keeper holds AI advisory state for streak-based rules.
type Keeper struct {
	advisoryKey    string
	nextProposalID int64
	lowUtilStreak  int64
	highUtilStreak int64
	config         Config
}

func NewKeeper(advisoryKey string, cfg Config) *Keeper {
	return &Keeper{
		advisoryKey:    advisoryKey,
		nextProposalID: 1,
		config:         cfg.withDefaults(),
	}
}

// AnalyzeAndPropose executes rule-based advisory logic.
func (k *Keeper) AnalyzeAndPropose(
	kpis ltypes.EpochKPIs,
	params ltypes.NetworkParams,
) *ltypes.AIProposal {
	if kpis.AvgUtilization < k.config.LowUtilThreshold {
		k.lowUtilStreak++
		k.highUtilStreak = 0
	} else if kpis.AvgUtilization > k.config.HighUtilThreshold {
		k.highUtilStreak++
		k.lowUtilStreak = 0
	} else {
		k.lowUtilStreak = 0
		k.highUtilStreak = 0
	}

	// R1: low utilization for 3 epochs and low fees -> increase block gas limit by 5%.
	if k.lowUtilStreak >= 3 && kpis.AvgBaseFee < k.config.MinFeeTarget {
		newVal := scaled(params.BlockGasLimit, 1.05, k.config.MinBlockGasLimit, k.config.MaxBlockGasLimit)
		if newVal != params.BlockGasLimit {
			return k.makeProposal(
				kpis,
				ParamBlockGasLimit,
				params.BlockGasLimit,
				newVal,
				"low utilization (3 epochs) and low fees: increase block gas limit by 5%",
			)
		}
	}

	// R2: high utilization for 2 epochs -> decrease block gas limit by 5%.
	if k.highUtilStreak >= 2 {
		newVal := scaled(params.BlockGasLimit, 0.95, k.config.MinBlockGasLimit, k.config.MaxBlockGasLimit)
		if newVal != params.BlockGasLimit {
			return k.makeProposal(
				kpis,
				ParamBlockGasLimit,
				params.BlockGasLimit,
				newVal,
				"high utilization (2 epochs): decrease block gas limit by 5%",
			)
		}
	}

	// R3: average base fee above max target -> decrease base fee by 10%.
	if kpis.AvgBaseFee > k.config.MaxFeeTarget {
		newVal := scaled(params.BaseFeePerGas, 0.90, k.config.MinBaseFee, k.config.MaxBaseFee)
		if newVal != params.BaseFeePerGas {
			return k.makeProposal(
				kpis,
				ParamBaseFeePerGas,
				params.BaseFeePerGas,
				newVal,
				"base fee above target: decrease base_fee_per_gas by 10%",
			)
		}
	}

	// R4: average base fee below min target -> increase base fee by 5%.
	if kpis.AvgBaseFee < k.config.MinFeeTarget {
		newVal := scaled(params.BaseFeePerGas, 1.05, k.config.MinBaseFee, k.config.MaxBaseFee)
		if newVal != params.BaseFeePerGas {
			return k.makeProposal(
				kpis,
				ParamBaseFeePerGas,
				params.BaseFeePerGas,
				newVal,
				"base fee below target: increase base_fee_per_gas by 5%",
			)
		}
	}

	return nil
}

func (k *Keeper) makeProposal(
	kpis ltypes.EpochKPIs,
	parameter string,
	currentValue int64,
	proposedValue int64,
	rationale string,
) *ltypes.AIProposal {
	proposal := &ltypes.AIProposal{
		ProposalID:      k.nextProposalID,
		EpochSubmitted:  kpis.Epoch,
		Parameter:       parameter,
		CurrentValue:    currentValue,
		ProposedValue:   proposedValue,
		Rationale:       rationale,
		ActivationEpoch: kpis.Epoch + 2,
		AdvisoryKey:     k.advisoryKey,
	}
	proposal.AdvisorySignature = signProposal(*proposal)
	k.nextProposalID++
	return proposal
}

func signProposal(proposal ltypes.AIProposal) string {
	payload := fmt.Sprintf(
		"%s|%d|%s|%d|%d|%d",
		proposal.AdvisoryKey,
		proposal.EpochSubmitted,
		proposal.Parameter,
		proposal.CurrentValue,
		proposal.ProposedValue,
		proposal.ActivationEpoch,
	)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func scaled(current int64, multiplier float64, minValue int64, maxValue int64) int64 {
	next := int64(float64(current) * multiplier)
	if next < minValue {
		next = minValue
	}
	if next > maxValue {
		next = maxValue
	}
	return next
}
