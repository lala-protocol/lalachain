package cosmosapp

import (
	"bytes"
	"testing"

	protoapp "github.com/lala-protocol/lalachain/chain/app"
)

func TestModuleManagerIntegration(t *testing.T) {
	app, err := New(42, protoapp.Phase1, 10)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	var buf bytes.Buffer
	if err := app.Run(30, &buf); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// Verify snapshot was produced.
	snapshot, ok := app.Snapshot()
	if !ok {
		t.Fatal("expected snapshot after run")
	}
	if snapshot.Epochs != 30 {
		t.Fatalf("expected 30 epochs, got %d", snapshot.Epochs)
	}
	if snapshot.Phase != "phase1" {
		t.Fatalf("expected phase1, got %s", snapshot.Phase)
	}

	// Verify state hash is deterministic.
	hash := app.StateHash()
	if hash == "" {
		t.Fatal("expected non-empty state hash")
	}

	// Verify module manager is wired.
	if app.ModuleManager() == nil {
		t.Fatal("expected non-nil module manager")
	}
}

func TestModuleManagerPhase2(t *testing.T) {
	app, err := New(99, protoapp.Phase2, 6)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	var buf bytes.Buffer
	if err := app.Run(30, &buf); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	snapshot, ok := app.Snapshot()
	if !ok {
		t.Fatal("expected snapshot after run")
	}
	if snapshot.Phase != "phase2" {
		t.Fatalf("expected phase2, got %s", snapshot.Phase)
	}
	if snapshot.Validators != 6 {
		t.Fatalf("expected 6 validators, got %d", snapshot.Validators)
	}
}

func TestStoreBackedKPIPersistence(t *testing.T) {
	app, err := New(42, protoapp.Phase1, 10)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	var buf bytes.Buffer
	if err := app.Run(5, &buf); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// Verify KPIs were persisted to the telemetry store.
	ctx := app.newContext(5)
	kpis, found := app.telemetryKeeper.GetKPIs(ctx, 5)
	if !found {
		t.Fatal("expected KPIs for epoch 5 to be persisted")
	}
	if kpis.Epoch != 5 {
		t.Fatalf("expected epoch 5, got %d", kpis.Epoch)
	}
	if kpis.AvgUtilization <= 0 {
		t.Fatal("expected positive avg utilization")
	}
}

func TestStoreBackedAdvisorState(t *testing.T) {
	app, err := New(42, protoapp.Phase1, 10)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	var buf bytes.Buffer
	if err := app.Run(10, &buf); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// Verify aiadvisor state was persisted via EndBlock hook.
	// The aiKeeper should have been SaveState'd by the module EndBlocker.
	ctx := app.newContext(10)
	store := ctx.KVStore(app.keys[storeKeyAIAdvisor])
	stateBz := store.Get([]byte("state"))
	if stateBz == nil {
		t.Fatal("expected aiadvisor state to be persisted in store")
	}
}

func TestDeterministicStateHash(t *testing.T) {
	// Two apps with same seed must produce identical state hashes.
	app1, _ := New(123, protoapp.Phase1, 10)
	app2, _ := New(123, protoapp.Phase1, 10)

	var buf1, buf2 bytes.Buffer
	_ = app1.Run(20, &buf1)
	_ = app2.Run(20, &buf2)

	hash1 := app1.StateHash()
	hash2 := app2.StateHash()
	if hash1 != hash2 {
		t.Fatalf("state hashes diverged: %s != %s", hash1, hash2)
	}
}
