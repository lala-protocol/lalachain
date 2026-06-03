---
title: "LalaChain Validator Guide"
description: "Comprehensive guide for LalaChain validator operators."
---

# LalaChain Validator Guide

## Overview

LalaChain is a Cosmos SDK v0.50 blockchain with AI-governed parameter optimization. Validators participate in both consensus and governance by voting on AI-originated proposals.

**Token:** $LALA (denom: `ulala`, 6 decimals)  
**Chain ID:** `lalachain-testnet-1`  
**Block Time:** ~5 seconds  
**Epoch:** 100 blocks (~8.3 minutes)

## Hardware Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 4 cores | 8 cores |
| RAM | 8 GB | 16 GB |
| Storage | 100 GB SSD | 500 GB NVMe |
| Network | 100 Mbps | 1 Gbps |

## Installation

### Option 1: Build from source

```bash
git clone https://github.com/lala-protocol/whitepaper.git
cd whitepaper
make build
```

The binary will be at `./build/lalachaind`.

### Option 2: Docker

```bash
docker build -t lalachain:latest .
```

## Setup

### 1. Initialize the node

```bash
lalachaind init <moniker> --chain-id lalachain-testnet-1
```

### 2. Get the genesis file

Download the official genesis file:

```bash
curl -o ~/.lalachaind/config/genesis.json https://raw.githubusercontent.com/lala-protocol/whitepaper/main/testnet/genesis.json
```

### 3. Configure persistent peers

Edit `~/.lalachaind/config/config.toml`:

```toml
persistent_peers = "<node_id>@<ip>:26656,..."
```

### 4. Set minimum gas prices

Edit `~/.lalachaind/config/app.toml`:

```toml
minimum-gas-prices = "0ulala"
```

### 5. Start the node

```bash
lalachaind start
```

## Creating a Validator

### 1. Get tokens

Request tokens from the faucet:

```bash
curl -X POST http://faucet.lalachain.network/faucet \
  -H "Content-Type: application/json" \
  -d '{"address": "lala1..."}'
```

### 2. Create validator transaction

```bash
lalachaind tx staking create-validator \
  --amount=1000000000ulala \
  --pubkey=$(lalachaind comet show-validator) \
  --moniker="<your_moniker>" \
  --chain-id=lalachain-testnet-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --from=<your_key>
```

## Governance Participation

LalaChain uses AI-originated proposals. As a validator, you vote on parameter changes proposed by the AI advisor.

### View pending proposals

```bash
lalachaind query lalagov proposals
```

### Vote on a proposal

```bash
lalachaind tx lalagov vote <proposal_id> --approve=true --from=<your_key> --chain-id=lalachain-testnet-1
```

## Monitoring

### Dashboard

Access the web dashboard at: `http://dashboard.lalachain.network`

### Prometheus metrics

Metrics are exposed at `:26660/metrics` by default. Enable in `config.toml`:

```toml
[instrumentation]
prometheus = true
```

## Key Addresses

| Prefix | Usage |
|--------|-------|
| `lala` | Account addresses |
| `lalapub` | Account public keys |
| `lalavaloper` | Validator operator addresses |
| `lalavaloperpub` | Validator public keys |
| `lalavalcons` | Consensus node addresses |

## Tokenomics

- **Total Initial Supply:** 100,000,000 LALA
- **Inflation:** 7–20% annually (targets 67% bonded)
- **Initial Inflation:** 13%
- **Community Tax:** 2%
- **Unbonding Period:** 21 days
- **Max Validators:** 100
- **Min Commission:** 5%

## Slashing

| Infraction | Slash Fraction | Jail Duration |
|-----------|---------------|---------------|
| Double Sign | 5% | Permanent (tombstoned) |
| Downtime | 0.01% | 10 minutes |
| Signed Blocks Window | 100 blocks | — |
| Min Signed Per Window | 50% | — |

## Troubleshooting

### Node won't sync
- Check persistent_peers in config.toml
- Ensure ports 26656 (P2P) and 26657 (RPC) are open
- Verify genesis.json hash matches the network

### Validator jailed
- If jailed for downtime: `lalachaind tx slashing unjail --from=<key> --chain-id=lalachain-testnet-1`
- If double-sign: cannot be unjailed (tombstoned)

### Out of memory
- Increase `cache_size` in app.toml
- Enable state-sync to reduce local storage
