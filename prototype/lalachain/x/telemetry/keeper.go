package telemetry

import (
	"errors"
	"fmt"

	ltypes "github.com/lala-protocol/whitepaper/prototype/lalachain/types"
)

// Keeper computes deterministic telemetry from finalized block data.
type Keeper struct{}

func NewKeeper() Keeper {
	return Keeper{}
}

func (k Keeper) CalculateEpochKPIs(
	epoch int64,
	blocks []ltypes.FinalizedBlock,
	params ltypes.NetworkParams,
	validators []ltypes.Validator,
	mempoolTrend float64,
	slashingEvents int64,
) (ltypes.EpochKPIs, error) {
	if params.BlockGasLimit <= 0 {
		return ltypes.EpochKPIs{}, errors.New("block gas limit must be > 0")
	}
	if len(blocks) == 0 {
		return ltypes.EpochKPIs{}, errors.New("cannot compute telemetry from empty block set")
	}

	blockTimes := make([]float64, 0, len(blocks))
	utilizations := make([]float64, 0, len(blocks))
	baseFees := make([]float64, 0, len(blocks))

	for _, block := range blocks {
		if block.GasUsed < 0 {
			return ltypes.EpochKPIs{}, fmt.Errorf("invalid gas used: %d", block.GasUsed)
		}
		blockTimes = append(blockTimes, float64(block.BlockTimeMS))
		utilizations = append(utilizations, float64(block.GasUsed)/float64(params.BlockGasLimit))
		baseFees = append(baseFees, float64(block.BaseFee))
	}

	totalStake := 0.0
	for _, validator := range validators {
		totalStake += validator.Stake
	}

	return ltypes.EpochKPIs{
		Epoch:               epoch,
		AvgBlockTimeMS:      mean(blockTimes),
		BlockTimeVarianceMS: variance(blockTimes),
		AvgUtilization:      mean(utilizations),
		AvgBaseFee:          int64(mean(baseFees)),
		ValidatorCount:      len(validators),
		TotalStakedRatio:    totalStake,
		MempoolSizeTrend:    mempoolTrend,
		SlashingEvents:      slashingEvents,
	}, nil
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

func variance(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	m := mean(values)
	sum := 0.0
	for _, value := range values {
		delta := value - m
		sum += delta * delta
	}
	return sum / float64(len(values))
}
