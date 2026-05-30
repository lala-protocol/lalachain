# LalaChain Phase 1 Prototype

This directory contains a Cosmos-style Phase 1 prototype scaffold focused on the
implementation track defined in `IMPLEMENTATION_FEASIBILITY.md`:

- `x/telemetry` - deterministic KPI aggregation from finalized epoch data.
- `x/aiadvisor` - rule-based AI advisory engine that emits signed proposals.
- `x/gov` - governance flow for AI proposal validation, voting, and scheduled activation.
- `app` - wiring layer that executes the adaptation loop end-to-end.

This is a single-process prototype to validate module boundaries and behavior
before full Cosmos SDK integration.

## Run

```bash
go run ./cmd/lalachain --epochs 30 --seed 42
```

## Test

```bash
go test ./...
```
