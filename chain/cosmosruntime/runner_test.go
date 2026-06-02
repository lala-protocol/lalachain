package cosmosruntime

import (
	"bytes"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewRunnerRejectsInvalidPhase(t *testing.T) {
	if _, err := NewRunner(42, "phase9", 10); err == nil {
		t.Fatal("expected invalid phase error")
	}
}

func TestRunPhase2ThroughCosmosRuntime(t *testing.T) {
	runner, err := NewRunner(42, "phase2", 8)
	if err != nil {
		t.Fatalf("new runner failed: %v", err)
	}

	var out bytes.Buffer
	if err := runner.Run(sdk.Context{}, 20, &out); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	log := out.String()
	if !strings.Contains(log, "Cosmos SDK Scaffold Runtime") {
		t.Fatal("expected cosmos scaffold header")
	}
	if !strings.Contains(log, "PHASE 2 GOVERNANCE AUDIT") {
		t.Fatal("expected phase 2 governance audit output")
	}

	snapshot, ok := runner.Snapshot()
	if !ok {
		t.Fatal("expected snapshot after run")
	}
	if snapshot.Phase != "phase2" {
		t.Fatalf("expected phase2 snapshot, got %s", snapshot.Phase)
	}
	if runner.StateHash() == "" {
		t.Fatal("expected non-empty state hash")
	}
}
