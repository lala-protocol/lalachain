package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrototypeRunEndToEnd(t *testing.T) {
	prototype := NewPrototype(42)
	var out bytes.Buffer

	if err := prototype.Run(25, &out); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if len(prototype.KPIHistory()) != 25 {
		t.Fatalf("expected 25 KPI snapshots, got %d", len(prototype.KPIHistory()))
	}
	if len(prototype.History()) == 0 {
		t.Fatal("expected at least one governance outcome")
	}
	if !strings.Contains(out.String(), "PROTOTYPE SUMMARY") {
		t.Fatal("expected summary output")
	}
}
