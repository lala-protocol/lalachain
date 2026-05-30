package telemetry

import (
	"math"
	"testing"

	ltypes "github.com/lala-protocol/whitepaper/prototype/lalachain/types"
)

func TestCalculateEpochKPIsDeterministic(t *testing.T) {
	keeper := NewKeeper()
	params := ltypes.NetworkParams{
		BlockGasLimit:     100,
		BaseFeePerGas:     1_000,
		TargetBlockTimeMS: 6_000,
	}
	blocks := []ltypes.FinalizedBlock{
		{BlockTimeMS: 6_000, GasUsed: 50, BaseFee: 900},
		{BlockTimeMS: 6_200, GasUsed: 60, BaseFee: 1_100},
		{BlockTimeMS: 5_800, GasUsed: 70, BaseFee: 1_000},
	}
	validators := []ltypes.Validator{
		{Address: "val1", Stake: 0.6},
		{Address: "val2", Stake: 0.4},
	}

	kpis, err := keeper.CalculateEpochKPIs(7, blocks, params, validators, 2.5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if kpis.Epoch != 7 {
		t.Fatalf("unexpected epoch: %d", kpis.Epoch)
	}
	if math.Abs(kpis.AvgBlockTimeMS-6000) > 0.001 {
		t.Fatalf("unexpected avg block time: %.4f", kpis.AvgBlockTimeMS)
	}
	if math.Abs(kpis.BlockTimeVarianceMS-26666.6667) > 0.01 {
		t.Fatalf("unexpected variance: %.4f", kpis.BlockTimeVarianceMS)
	}
	if math.Abs(kpis.AvgUtilization-0.6) > 1e-9 {
		t.Fatalf("unexpected avg utilization: %.10f", kpis.AvgUtilization)
	}
	if kpis.AvgBaseFee != 1000 {
		t.Fatalf("unexpected avg base fee: %d", kpis.AvgBaseFee)
	}
	if kpis.ValidatorCount != 2 {
		t.Fatalf("unexpected validator count: %d", kpis.ValidatorCount)
	}
	if math.Abs(kpis.TotalStakedRatio-1.0) > 1e-9 {
		t.Fatalf("unexpected total staked ratio: %.10f", kpis.TotalStakedRatio)
	}
	if math.Abs(kpis.MempoolSizeTrend-2.5) > 1e-9 {
		t.Fatalf("unexpected mempool trend: %.10f", kpis.MempoolSizeTrend)
	}
	if kpis.SlashingEvents != 1 {
		t.Fatalf("unexpected slashing events: %d", kpis.SlashingEvents)
	}
}
