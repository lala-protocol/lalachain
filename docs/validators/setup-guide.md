---
title: "Validator Setup Guide"
description: "Complete guide to setting up a LalaChain validator node."
---

# Validator Setup Guide

**Step-by-step guide to setting up a LalaChain validator node.**

---

## Step 1: Install Prerequisites

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y build-essential git curl jq

# Install Go 1.21+
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

---

## Step 2: Build LalaChain

```bash
# Clone repository
git clone https://github.com/lalachain/lalachain.git
cd lalachain/chain

# Build
go build -o build/lalachaind ./cmd/lalachain

# Install globally
sudo cp build/lalachaind /usr/local/bin/
```

---

## Step 3: Initialize Node

```bash
# Initialize (choose a descriptive moniker)
lalachaind init "my-validator-name" --chain-id lalachain-1

# Download genesis file (from network coordinator)
curl -o ~/.lalachaind/config/genesis.json https://raw.githubusercontent.com/lalachain/mainnet/main/genesis.json

# Set seeds/peers in config.toml
sed -i 's/seeds = ""/seeds = "seed-node-id@seed-ip:26656"/' ~/.lalachaind/config/config.toml
```

---

## Step 4: Create Validator Key

```bash
# Create key (SAVE THE MNEMONIC!)
lalachaind keys add validator

# Fund the account (get tokens from faucet or transfer)
# Verify balance
lalachaind query bank balances $(lalachaind keys show validator -a)
```

---

## Step 5: Start the Node and Sync

```bash
# Start node
lalachaind start

# Wait for full sync (check status)
curl localhost:26657/status | jq '.result.sync_info.catching_up'
# Should return "false" when synced
```

---

## Step 6: Create Validator Transaction

Once synced:

```bash
lalachaind tx staking create-validator \
  --amount 1000000000ulala \
  --pubkey $(lalachaind tendermint show-validator) \
  --moniker "my-validator-name" \
  --chain-id lalachain-1 \
  --commission-rate 0.10 \
  --commission-max-rate 0.20 \
  --commission-max-change-rate 0.01 \
  --min-self-delegation 1000000 \
  --from validator \
  --fees 50000ulala -y
```

---

## Step 7: Verify Validator Status

```bash
# Check your validator
lalachaind query staking validator $(lalachaind keys show validator --bech val -a)

# Check if in active set
lalachaind query staking validators --output json | jq '.validators[] | select(.description.moniker=="my-validator-name")'
```

---

## Step 8: Set Up as System Service

```bash
sudo tee /etc/systemd/system/lalachaind.service > /dev/null <<EOF
[Unit]
Description=LalaChain Node
After=network.target

[Service]
Type=simple
User=$USER
ExecStart=/usr/local/bin/lalachaind start
Restart=on-failure
RestartSec=10
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable lalachaind
sudo systemctl start lalachaind
```

---

## Step 9: Configure Monitoring

See [Monitoring](monitoring.md) for Prometheus + Grafana setup.
