<div align="center">

# LalaChain

### The Self-Optimizing Blockchain

[![Build Status](https://img.shields.io/github/actions/workflow/status/lala-protocol/lalachain/ci.yml?branch=main&style=flat-square&logo=github)](https://github.com/lala-protocol/lalachain/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![Cosmos SDK](https://img.shields.io/badge/Cosmos%20SDK-v0.50.9-6F42C1?style=flat-square)](https://github.com/cosmos/cosmos-sdk)
[![License: CC BY 4.0](https://img.shields.io/badge/License-CC%20BY%204.0-lightgrey.svg?style=flat-square)](https://creativecommons.org/licenses/by/4.0/)
[![Testnet](https://img.shields.io/badge/Testnet-Live-brightgreen?style=flat-square)](docs/validators/setup-guide.md)

**A Layer 1 blockchain that tunes itself.**

LalaChain embeds a deterministic AI Advisor directly into its consensus layer. It watches. It learns patterns. It proposes. Validators decide. The chain evolves — every epoch, autonomously, transparently, safely.

[Whitepaper](paper/lala-protocol-whitepaper.pdf) · [Documentation](docs/) · [Run a Validator](docs/validators/setup-guide.md) · [API Reference](docs/developers/rest-api-reference.md)

</div>

---

## Why LalaChain?

Every blockchain launches with parameters — gas limits, fee floors, block sizes — tuned for day one. But networks grow. Usage shifts. Conditions change.

**Most chains don't adapt.** Governance proposals take weeks. Voter turnout is abysmal. Parameters rot.

LalaChain solves this with a closed-loop control system built into the protocol itself:

```
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│   Observe          Analyze           Propose        Decide   │
│                                                              │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────┐ │
│   │Telemetry │───▶│AI Advisor│───▶│ Proposal │───▶│ Vote │ │
│   │  Module  │    │  Engine  │    │  Module  │    │      │ │
│   └──────────┘    └──────────┘    └──────────┘    └──────┘ │
│        ▲                                              │      │
│        └──────────────────────────────────────────────┘      │
│                    Continuous Feedback Loop                   │
└──────────────────────────────────────────────────────────────┘
```

The result: a chain that maintains peak performance without human babysitting.

---

## Core Principles

| | Principle | How |
|---|---|---|
| 🔒 | **Safety First** | AI proposals bounded to ±5-10% per epoch. Hard min/max on all parameters. Cannot exceed bounds even if every validator votes yes. |
| 🗳️ | **Human Sovereignty** | The AI *proposes*. Validators *decide*. Every change requires 66% quorum + 51% approval. Nothing happens without consensus. |
| 🔍 | **Radical Transparency** | Every AI decision includes the triggering rule, input KPIs, and reasoning. Fully reproducible. Fully auditable. On-chain forever. |
| ⚡ | **Instant Finality** | CometBFT consensus — once a block is committed, it's permanent. No rollbacks. No reorgs. ~5 second blocks. |
| 🧮 | **Deterministic AI** | No neural networks. No randomness. No external data. Same inputs → same outputs on every validator, every time. |

---

## Quick Start

```bash
# Clone
git clone https://github.com/lala-protocol/lalachain.git
cd lalachain

# Build
make build

# Run tests
make test

# Launch a 4-validator local testnet
make docker-testnet

# Or run the network simulation (no dependencies)
python simulation/simulation.py --epochs 50 --seed 42
```

<details>
<summary><strong>System Requirements</strong></summary>

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| Go | 1.22+ | Latest |
| Docker | 20.10+ | Latest |
| Python | 3.8+ | 3.11+ |
| Node.js | 18+ | 20+ (dashboard only) |
| RAM | 4 GB | 8 GB |
| Disk | 20 GB SSD | 100 GB NVMe |

</details>

---

## Architecture

```
lalachain/
├── chain/                  Blockchain node (Go)
│   ├── cmd/lalachaind/     Binary entrypoint
│   ├── app/                Application wiring + epoch logic
│   ├── x/telemetry/        On-chain metrics collection
│   ├── x/aiadvisor/        Deterministic rule engine (4 rules)
│   ├── x/lalagov/          Governance: proposals, votes, activation
│   └── api/                Generated protobuf/gRPC interfaces
├── dashboard/              Real-time monitoring UI (Next.js)
├── docs/                   Full documentation (70+ pages)
├── paper/                  Whitepaper (LaTeX + PDF)
├── simulation/             Python network simulation
├── scripts/                Deployment & testnet tooling
├── Dockerfile              Multi-stage production build
├── docker-compose.yml      4-validator local cluster
└── Makefile                Build, test, lint, proto, deploy
```

---

## The AI Advisor

Four deterministic rules. No magic. No black boxes.

| Rule | Condition | Action | Rationale |
|------|-----------|--------|-----------|
| **R1** | Low utilization ≥3 epochs AND fee below target | Increase gas limit +5% | Network has headroom — expand capacity |
| **R2** | High utilization ≥2 epochs | Decrease gas limit -5% | Sustained congestion — protect block times |
| **R3** | Fee exceeds maximum target | Decrease base fee -10% | Users overpaying — reduce cost pressure |
| **R4** | Fee below minimum target | Increase base fee +5% | Spam risk — raise the floor |

**Safety constraints:**
- `block_gas_limit`: 10M – 30M (hard bounds)
- `base_fee_per_gas`: 100M – 10B ulala (hard bounds)
- Maximum one proposal per epoch
- 2-epoch activation delay after approval

---

## Network Parameters

| Parameter | Value | Governed By |
|-----------|-------|-------------|
| Block Time | ~5 seconds | CometBFT |
| Epoch Length | 10 blocks | Protocol |
| Consensus | CometBFT BFT | 2/3+ stake |
| Finality | Instant (1 block) | Protocol |
| Token | LALA (ulala = 10⁻⁶) | — |
| Unbonding Period | 21 days | Governance |
| Quorum | 66% | Protocol |
| Approval Threshold | 51% | Protocol |
| Activation Delay | 2 epochs | Protocol |

---

## API

```bash
# Current network KPIs
curl localhost:1317/lala/telemetry/v1/kpis

# AI Advisor state (streaks, last proposal)
curl localhost:1317/lala/aiadvisor/v1/state

# Governance config
curl localhost:1317/lala/lalagov/v1/config

# Proposal history
curl localhost:1317/lala/lalagov/v1/history
```

Full API reference: [`docs/developers/rest-api-reference.md`](docs/developers/rest-api-reference.md)

---

## Documentation

| Section | Description |
|---------|-------------|
| [Getting Started](docs/getting-started/) | Blockchain fundamentals for newcomers |
| [How It Works](docs/how-it-works/) | Consensus, epochs, AI, fee model |
| [Tokenomics](docs/tokenomics/) | LALA supply, staking, fee economics |
| [Validator Guide](docs/validators/) | Setup, monitoring, rewards, slashing |
| [Developer Docs](docs/developers/) | API, SDK, building modules |
| [Smart Contracts](docs/smart-contracts/) | CosmWasm development |
| [Governance](docs/governance/) | Proposals, voting, AI transparency |
| [Security](docs/security/) | Wallet safety, network security, AI safety |
| [FAQ](docs/faq.md) | 50+ answered questions |

---

## For Validators

```bash
# Initialize node
lalachaind init my-validator --chain-id lalachain-testnet-1

# Create keys
lalachaind keys add validator --keyring-backend os

# Start node
lalachaind start --minimum-gas-prices=0.025ulala
```

Hardware requirements, monitoring setup, and operational guides: [`docs/validators/`](docs/validators/)

---

## Development

```bash
make build          # Build lalachaind binary → ./build/lalachaind
make test           # Run all Go tests
make lint           # golangci-lint
make proto          # Regenerate protobuf code
make tidy           # go mod tidy
make docker         # Build Docker image
make docker-testnet # Launch 4-validator Docker cluster
make simulation     # Run Python network simulation
make clean          # Remove build artifacts
```

---

## Contributing

We welcome contributions from the community. Whether you're fixing a typo, improving documentation, or proposing a new module — we'd love to hear from you.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Commit your changes (`git commit -m 'feat: add new feature'`)
4. Push to the branch (`git push origin feature/my-change`)
5. Open a Pull Request

Please read the [documentation](docs/) to understand the architecture before making protocol-level changes.

---

## Roadmap

- [x] Whitepaper & formal specification
- [x] Network simulation (Python)
- [x] Chain prototype with AI governance loop
- [x] Full Cosmos SDK integration (auth, bank, staking, mint, slashing)
- [x] CometBFT consensus with instant finality
- [x] REST/gRPC API endpoints
- [x] Multi-validator Docker testnet
- [x] Real-time monitoring dashboard
- [x] Comprehensive documentation (70+ pages)
- [ ] Security audit
- [ ] Public incentivized testnet
- [ ] IBC integration
- [ ] Mainnet launch

---

## Community

- [GitHub Issues](https://github.com/lala-protocol/lalachain/issues) — Bug reports & feature requests
- [Documentation](docs/) — Everything you need to know

---

## License

This project is licensed under [Creative Commons Attribution 4.0 International](LICENSE.md).

---

<div align="center">
<sub>Built with Cosmos SDK · Secured by CometBFT · Governed by Validators · Advised by AI</sub>
</div>

Please see the `LICENSE.md` file for the full license text. 


---

We are excited about the potential of adaptive blockchain technology and look forward to community engagement as Lala Protocol evolves!
