package aiadvisor

import (
	"testing"

	ltypes "github.com/lala-protocol/whitepaper/prototype/lalachain/types"
)

func TestAnalyzeAndProposeLowUtilizationRule(t *testing.T) {
	keeper := NewKeeper("ai-key", DefaultConfig())
	params := ltypes.DefaultNetworkParams()

	// First two epochs may trigger base fee proposals.
	_ = keeper.AnalyzeAndPropose(ltypes.EpochKPIs{
		Epoch:          1,
		AvgUtilization: 0.30,
		AvgBaseFee:     700_000_000,
	}, params)
	_ = keeper.AnalyzeAndPropose(ltypes.EpochKPIs{
		Epoch:          2,
		AvgUtilization: 0.35,
		AvgBaseFee:     700_000_000,
	}, params)

	proposal := keeper.AnalyzeAndPropose(ltypes.EpochKPIs{
		Epoch:          3,
		AvgUtilization: 0.32,
		AvgBaseFee:     700_000_000,
	}, params)
	if proposal == nil {
		t.Fatal("expected proposal, got nil")
	}
	if proposal.Parameter != ParamBlockGasLimit {
		t.Fatalf("expected %s proposal, got %s", ParamBlockGasLimit, proposal.Parameter)
	}
	if proposal.ProposedValue != int64(float64(params.BlockGasLimit)*1.05) {
		t.Fatalf("unexpected proposal value: %d", proposal.ProposedValue)
	}
	if proposal.AdvisoryKey != "ai-key" {
		t.Fatalf("unexpected advisory key: %s", proposal.AdvisoryKey)
	}
	if proposal.AdvisorySignature == "" {
		t.Fatal("expected non-empty advisory signature")
	}
}

func TestAnalyzeAndProposeHighUtilizationRule(t *testing.T) {
	keeper := NewKeeper("ai-key", DefaultConfig())
	params := ltypes.DefaultNetworkParams()

	first := keeper.AnalyzeAndPropose(ltypes.EpochKPIs{
		Epoch:          1,
		AvgUtilization: 0.90,
		AvgBaseFee:     1_500_000_000,
	}, params)
	if first != nil {
		t.Fatalf("expected no proposal on first high util epoch, got %v", first.Parameter)
	}

	second := keeper.AnalyzeAndPropose(ltypes.EpochKPIs{
		Epoch:          2,
		AvgUtilization: 0.88,
		AvgBaseFee:     1_700_000_000,
	}, params)
	if second == nil {
		t.Fatal("expected proposal, got nil")
	}
	if second.Parameter != ParamBlockGasLimit {
		t.Fatalf("expected %s proposal, got %s", ParamBlockGasLimit, second.Parameter)
	}
	if second.ProposedValue != int64(float64(params.BlockGasLimit)*0.95) {
		t.Fatalf("unexpected proposal value: %d", second.ProposedValue)
	}
}
