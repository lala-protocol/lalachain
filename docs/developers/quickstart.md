# Developer Quickstart

**Get a LalaChain node running locally in 5 minutes and make your first API call.**

---

## Prerequisites

- Go 1.21+ installed
- Git
- ~500MB disk space

---

## Step 1: Build the Binary

```bash
cd chain
go build -o build/lalachaind ./cmd/lalachain
```

---

## Step 2: Initialize the Node

```bash
# Initialize with a chain ID
./build/lalachaind init my-node --chain-id lalachain-local

# Create a validator key
./build/lalachaind keys add validator

# Add genesis account with tokens
./build/lalachaind genesis add-genesis-account validator 1000000000000ulala

# Create genesis transaction
./build/lalachaind genesis gentx validator 500000000000ulala --chain-id lalachain-local

# Collect genesis transactions
./build/lalachaind genesis collect-gentxs
```

---

## Step 3: Start the Node

```bash
./build/lalachaind start
```

You should see blocks being produced every ~5 seconds:
```
INF committed state height=1 ...
INF committed state height=2 ...
```

---

## Step 4: Make Your First API Call

In a new terminal:

```bash
# Check node status
curl http://localhost:26657/status

# Get KPIs (after at least one epoch = 10 blocks)
curl http://localhost:1317/lala/telemetry/v1/kpis

# Get AI Advisor state
curl http://localhost:1317/lala/aiadvisor/v1/state

# Get governance history
curl http://localhost:1317/lala/lalagov/v1/history

# Check balance
curl http://localhost:1317/cosmos/bank/v1beta1/balances/$(./build/lalachaind keys show validator -a)
```

---

## Step 5: Send a Transaction

```bash
# Create a second wallet
./build/lalachaind keys add user1

# Send tokens
./build/lalachaind tx bank send validator $(./build/lalachaind keys show user1 -a) 1000000ulala \
  --fees 5000ulala \
  --chain-id lalachain-local -y
```

---

## What's Running

After starting, you have:

| Service | URL | Purpose |
|---------|-----|---------|
| CometBFT RPC | `localhost:26657` | Block/tx queries, WebSocket |
| Cosmos REST API | `localhost:1317` | Full REST API |
| P2P | `localhost:26656` | Peer connections |

---

## Project Structure

```
chain/
├── cmd/lalachain/main.go    # Entry point
├── app/
│   ├── prototype.go         # App definition, module wiring
│   └── prototype_test.go    # Integration tests
├── types/types.go           # Shared type definitions
└── x/                       # Custom modules
    ├── telemetry/           # KPI collection and computation
    ├── aiadvisor/           # AI rule engine
    └── gov/                 # Governance (proposals, voting)
```

---

## Next Steps

- [Architecture Deep Dive](architecture-deep-dive.md) — Understand the module system
- [REST API Reference](rest-api-reference.md) — Complete endpoint documentation
- [Building Modules](building-modules.md) — Create your own modules

---

**Next:** [Architecture Deep Dive](architecture-deep-dive.md)
