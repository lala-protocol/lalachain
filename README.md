# LalaChain

[![License: CC BY 4.0](https://img.shields.io/badge/License-CC%20BY%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)

**AI-Advised Adaptive Blockchain Protocol**

LalaChain is a Layer 1 proof-of-stake blockchain built on the Cosmos SDK that integrates a deterministic AI Advisor into its governance loop. The AI monitors real-time network metrics and proposes parameter adjustments—subject to validator vote—to maintain optimal throughput, fair fees, and network health without constant human intervention.

---

## Repository Structure

```
lalachain/
├── chain/              → Go blockchain node (Cosmos SDK + CometBFT)
│   ├── app/            → Application logic & epoch system
│   ├── cmd/lalachaind/ → Chain binary entrypoint
│   ├── x/telemetry/    → On-chain KPI tracking module
│   ├── x/aiadvisor/    → Deterministic AI rule engine
│   ├── x/lalagov/      → Governance with AI proposal integration
│   └── api/            → Generated protobuf/gRPC code
├── dashboard/          → Next.js monitoring dashboard
├── docs/               → Full documentation site (GitBook-compatible)
├── paper/              → Whitepaper (LaTeX source + PDF)
├── simulation/         → Python network simulation
└── scripts/            → Testnet initialization scripts
```

---

## Quick Start

### Build the chain binary

```bash
make build
# Output: ./build/lalachaind
```

### Run tests

```bash
make test
```

### Start a local testnet (4 validators)

```bash
make testnet
# or with Docker:
make docker-testnet
```

### Run the simulation

```bash
make simulation
# or directly:
python simulation/simulation.py --epochs 30 --seed 42
```

---

## Key Features

| Feature | Description |
|---------|-------------|
| **AI Advisor** | Deterministic rule engine proposes gas limit and fee adjustments based on sustained network patterns |
| **Bounded Safety** | All AI proposals are capped at ±5-10% per epoch with hard parameter bounds |
| **Validator Gating** | Every proposal requires 66% quorum + 51% approval from validators |
| **Instant Finality** | CometBFT consensus provides single-block finality (~5s) |
| **Epoch System** | 10-block epochs for KPI evaluation and governance cycles |
| **Full Cosmos Stack** | auth, bank, staking, distribution, slashing, mint + custom modules |

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│                  CometBFT Consensus              │
├─────────────────────────────────────────────────┤
│  x/telemetry  →  x/aiadvisor  →  x/lalagov     │
│  (KPI tracking)  (rule engine)   (vote + apply) │
├─────────────────────────────────────────────────┤
│  Standard Cosmos Modules (auth, bank, staking…) │
└─────────────────────────────────────────────────┘
```

---

## Documentation

Full documentation is available in [`docs/`](docs/):

- [Getting Started](docs/getting-started/) — Blockchain fundamentals
- [How It Works](docs/how-it-works/) — Consensus, AI, fees, epochs
- [Developer Docs](docs/developers/) — API reference, SDK, modules
- [Validator Guide](docs/validators/) — Setup, monitoring, rewards
- [Governance](docs/governance/) — Proposals, voting, AI transparency
- [Whitepaper Summary](docs/whitepaper-summary.md) — Executive overview

---

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /lala/telemetry/v1/kpis` | Current epoch KPIs |
| `GET /lala/aiadvisor/v1/state` | AI Advisor state & streaks |
| `GET /lala/lalagov/v1/config` | Governance configuration |
| `GET /lala/lalagov/v1/history` | Proposal history |

---

## Development

### Prerequisites

- Go 1.22+
- Node.js 18+ (for dashboard)
- Python 3.8+ (for simulation)
- Docker & Docker Compose (for multi-validator testnet)

### Useful commands

```bash
make build          # Build lalachaind binary
make test           # Run Go tests
make lint           # Run golangci-lint
make proto          # Regenerate protobuf code
make tidy           # go mod tidy
make docker         # Build Docker image
make docker-testnet # Launch 4-validator Docker testnet
```

---

## License

This work is licensed under [CC BY 4.0](LICENSE.md).

Please see the `LICENSE.md` file for the full license text. 


---

We are excited about the potential of adaptive blockchain technology and look forward to community engagement as Lala Protocol evolves!
