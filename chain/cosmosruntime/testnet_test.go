package cosmosruntime

import (
	"bytes"
	"testing"
)

func TestRunLocalTestnetConsensus(t *testing.T) {
	var out bytes.Buffer
	report, err := RunLocalTestnet(LocalTestnetConfig{
		Phase:      "phase2",
		Epochs:     12,
		Seed:       42,
		Validators: 6,
		Nodes:      4,
	}, &out)
	if err != nil {
		t.Fatalf("run local testnet failed: %v", err)
	}
	if len(report.Nodes) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(report.Nodes))
	}
	if !report.Consensus {
		t.Fatalf("expected consensus, divergences=%v", report.Divergences)
	}
	if report.ReferenceHash == "" {
		t.Fatal("expected non-empty reference hash")
	}
	if out.Len() == 0 {
		t.Fatal("expected testnet report output")
	}
}

func TestRunLocalTestnetRejectsValidatorOutOfRange(t *testing.T) {
	_, err := RunLocalTestnet(LocalTestnetConfig{
		Phase:      "phase2",
		Epochs:     5,
		Seed:       42,
		Validators: 3,
	}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected validator range error")
	}
}
