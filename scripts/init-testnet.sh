#!/bin/bash
# init-testnet.sh — Initialize a 4-validator LalaChain testnet
set -euo pipefail

CHAIN_ID="lalachain-testnet-1"
DENOM="ulala"
VALIDATORS=4
INITIAL_SUPPLY="100000000000000"  # 100M LALA in ulala
VALIDATOR_STAKE="25000000000000"  # 25M LALA per validator
HOME_BASE="./testnet"

echo "=== LalaChain Testnet Initialization ==="
echo "Chain ID: $CHAIN_ID"
echo "Validators: $VALIDATORS"
echo "Initial Supply: $INITIAL_SUPPLY $DENOM (100M LALA)"
echo ""

# Clean previous state
rm -rf "$HOME_BASE"
mkdir -p "$HOME_BASE"

# Initialize each validator
for i in $(seq 1 $VALIDATORS); do
    MONIKER="validator$i"
    NODE_HOME="$HOME_BASE/$MONIKER"

    echo "--- Initializing $MONIKER ---"
    lalachaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "$NODE_HOME" > /dev/null 2>&1

    # Create validator key
    lalachaind keys add "$MONIKER" --keyring-backend test --home "$NODE_HOME" > /dev/null 2>&1
done

# Use validator1's genesis as the canonical genesis
GENESIS="$HOME_BASE/validator1/config/genesis.json"

# Add genesis accounts with initial balances
for i in $(seq 1 $VALIDATORS); do
    MONIKER="validator$i"
    NODE_HOME="$HOME_BASE/$MONIKER"
    ADDR=$(lalachaind keys show "$MONIKER" -a --keyring-backend test --home "$NODE_HOME")
    lalachaind genesis add-genesis-account "$ADDR" "${VALIDATOR_STAKE}${DENOM}" --home "$HOME_BASE/validator1" --keyring-backend test
done

# Add faucet account with remaining supply
FAUCET_HOME="$HOME_BASE/faucet"
mkdir -p "$FAUCET_HOME"
cp -r "$HOME_BASE/validator1/config" "$FAUCET_HOME/"
lalachaind keys add faucet --keyring-backend test --home "$FAUCET_HOME" > /dev/null 2>&1
FAUCET_ADDR=$(lalachaind keys show faucet -a --keyring-backend test --home "$FAUCET_HOME")
FAUCET_AMOUNT=$((INITIAL_SUPPLY - VALIDATORS * VALIDATOR_STAKE))
lalachaind genesis add-genesis-account "$FAUCET_ADDR" "${FAUCET_AMOUNT}${DENOM}" --home "$HOME_BASE/validator1" --keyring-backend test

# Create genesis transactions (each validator self-delegates)
for i in $(seq 1 $VALIDATORS); do
    MONIKER="validator$i"
    NODE_HOME="$HOME_BASE/$MONIKER"

    # Copy updated genesis to each node
    cp "$GENESIS" "$NODE_HOME/config/genesis.json"

    # Generate gentx (25M LALA self-delegation)
    lalachaind genesis gentx "$MONIKER" "${VALIDATOR_STAKE}${DENOM}" \
        --chain-id "$CHAIN_ID" \
        --moniker "$MONIKER" \
        --commission-rate "0.10" \
        --commission-max-rate "0.20" \
        --commission-max-change-rate "0.01" \
        --min-self-delegation "1" \
        --keyring-backend test \
        --home "$NODE_HOME" > /dev/null 2>&1
done

# Collect gentxs into validator1 genesis
for i in $(seq 1 $VALIDATORS); do
    MONIKER="validator$i"
    NODE_HOME="$HOME_BASE/$MONIKER"
    cp "$NODE_HOME/config/gentx/"*.json "$HOME_BASE/validator1/config/gentx/" 2>/dev/null || true
done
lalachaind genesis collect-gentxs --home "$HOME_BASE/validator1" > /dev/null 2>&1

# Validate genesis
lalachaind genesis validate-genesis --home "$HOME_BASE/validator1"
echo ""

# Copy final genesis to all validators
FINAL_GENESIS="$HOME_BASE/validator1/config/genesis.json"
for i in $(seq 2 $VALIDATORS); do
    MONIKER="validator$i"
    NODE_HOME="$HOME_BASE/$MONIKER"
    cp "$FINAL_GENESIS" "$NODE_HOME/config/genesis.json"
done

# Configure persistent peers
echo "--- Configuring Peers ---"
PEERS=""
for i in $(seq 1 $VALIDATORS); do
    NODE_HOME="$HOME_BASE/validator$i"
    NODE_ID=$(lalachaind comet show-node-id --home "$NODE_HOME")
    if [ -n "$PEERS" ]; then
        PEERS="$PEERS,"
    fi
    PORT=$((26656 + (i-1) * 100))
    PEERS="${PEERS}${NODE_ID}@validator${i}:26656"
done

for i in $(seq 1 $VALIDATORS); do
    NODE_HOME="$HOME_BASE/validator$i"
    sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PEERS\"/" "$NODE_HOME/config/config.toml"
    # Enable API
    sed -i 's/enable = false/enable = true/' "$NODE_HOME/config/app.toml"
    # Set minimum gas prices
    sed -i "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0${DENOM}\"/" "$NODE_HOME/config/app.toml"
done

echo ""
echo "=== Testnet Initialized Successfully ==="
echo "Genesis: $FINAL_GENESIS"
echo "Validators: $VALIDATORS"
echo "Faucet address: $FAUCET_ADDR"
echo ""
echo "Start with: docker-compose up -d"
echo "Or manually: lalachaind start --home $HOME_BASE/validator1"
