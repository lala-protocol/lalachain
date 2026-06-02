# LalaChain Prototype (Phase 1 + Phase 2)

This directory contains a Cosmos-style prototype scaffold focused on the
implementation track defined in `IMPLEMENTATION_FEASIBILITY.md`:

- `x/telemetry` - deterministic KPI aggregation from finalized epoch data.
- `x/aiadvisor` - rule-based AI advisory engine that emits signed proposals.
- `x/gov` - governance flow for AI proposal validation, voting, and scheduled activation.
- `app` - wiring layer that executes the adaptation loop end-to-end.

This is a single-process prototype to validate module boundaries and behavior
before full Cosmos SDK integration.

Phase support:

- Phase 1 mode (`--phase 1`): baseline adaptation loop validation.
- Phase 2 mode (`--phase 2`): multi-validator (4-10) stress scenarios with governance/activation audit and stability findings.

Runtime support:

- `--runtime prototype` (default): in-memory standalone prototype execution.
- `--runtime cosmos`: Cosmos SDK scaffold runtime (accepts `sdk.Context`) for integration-oriented module boundaries.
- `--network single|multi`: Cosmos runtime can run one node or an orchestrated local multi-validator testnet.

Research support:

- `research/phase3_shadow_mode.py`: Phase 3 shadow-mode harness for model-vs-rule comparison and governance process draft output.

## Run

```bash
# Phase 1 (default)
go run ./cmd/lalachain --runtime prototype --phase 1 --epochs 30 --seed 42

# Phase 2
go run ./cmd/lalachain --runtime prototype --phase 2 --validators 10 --epochs 30 --seed 42

# Cosmos scaffold runtime (phase 2)
go run ./cmd/lalachain --runtime cosmos --phase 2 --validators 10 --epochs 30 --seed 42

# Cosmos local multi-validator testnet (phase 2)
go run ./cmd/lalachain --runtime cosmos --network multi --phase 2 --validators 8 --nodes 8 --epochs 20 --seed 42

# Docs-aligned start-to-finish pipeline (Phase 0 -> Phase 3 research harness)
powershell -ExecutionPolicy Bypass -File ./scripts/run_docs_start_to_finish.ps1
```

The start-to-finish script writes logs and reports under `reports/`, including:

- Phase 0 simulation log
- Phase 1 single-node cosmos run log
- Phase 2 multi-validator cosmos run log
- Phase 3 shadow-mode report (`.json`)
- consolidated markdown summary

## Test

```bash
go test ./...
```

Can you prove that AI-guided parameter governance produces a blockchain that users and validators prefer over existing alternatives?