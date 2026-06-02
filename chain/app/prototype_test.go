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

func TestPhase2RunIncludesAuditAndFindings(t *testing.T) {
	prototype := NewPhase2Prototype(42, 6)
	var out bytes.Buffer

	if err := prototype.Run(30, &out); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if len(prototype.validators) != 6 {
		t.Fatalf("expected 6 validators, got %d", len(prototype.validators))
	}

	output := out.String()
	if !strings.Contains(output, "PHASE 2 GOVERNANCE AUDIT") {
		t.Fatal("expected governance audit section in phase 2 output")
	}
	if !strings.Contains(output, "PHASE 2 STABILITY FINDINGS") {
		t.Fatal("expected stability findings section in phase 2 output")
	}
	if !strings.Contains(output, "scenario=") {
		t.Fatal("expected scenario labels in phase 2 KPI output")
	}
}

func TestPhase2ValidatorBounds(t *testing.T) {
	low := NewPhase2Prototype(42, 2)
	if len(low.validators) != 4 {
		t.Fatalf("expected low bound to clamp at 4 validators, got %d", len(low.validators))
	}

	high := NewPhase2Prototype(42, 50)
	if len(high.validators) != 10 {
		t.Fatalf("expected high bound to clamp at 10 validators, got %d", len(high.validators))
	}
}
