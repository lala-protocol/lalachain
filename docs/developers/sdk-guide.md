---
title: "SDK Guide"
description: "Using the LalaChain SDK to build applications."
---

# SDK Guide

**Libraries and tools for building applications on LalaChain.**

---

## JavaScript / TypeScript (CosmJS)

The recommended library for web and Node.js applications.

### Installation

```bash
npm install @cosmjs/stargate @cosmjs/proto-signing
```

### Connect to Node

```typescript
import { StargateClient } from "@cosmjs/stargate";

const client = await StargateClient.connect("http://localhost:26657");
const height = await client.getHeight();
console.log(`Current height: ${height}`);
```

### Query Balance

```typescript
const balance = await client.getBalance("lala1abc...", "ulala");
console.log(`Balance: ${balance.amount} ulala`);
```

### Send Tokens

```typescript
import { SigningStargateClient } from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

// Create wallet from mnemonic
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
  prefix: "lala"
});

// Connect signing client
const client = await SigningStargateClient.connectWithSigner(
  "http://localhost:26657",
  wallet
);

// Send tokens
const result = await client.sendTokens(
  senderAddress,
  recipientAddress,
  [{ denom: "ulala", amount: "1000000" }],
  { amount: [{ denom: "ulala", amount: "5000" }], gas: "200000" }
);

console.log(`TX hash: ${result.transactionHash}`);
```

### Query LalaChain Custom Endpoints

```typescript
// Custom endpoints use the REST API directly
const response = await fetch("http://localhost:1317/lala/telemetry/v1/kpis");
const kpis = await response.json();
```

---

## Go (Cosmos SDK Client)

For Go applications and module development.

### Query Client

```go
package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func main() {
    conn, _ := grpc.Dial("localhost:9090", grpc.WithInsecure())
    defer conn.Close()
    
    client := banktypes.NewQueryClient(conn)
    res, _ := client.Balance(context.Background(), &banktypes.QueryBalanceRequest{
        Address: "lala1abc...",
        Denom:   "ulala",
    })
    
    fmt.Printf("Balance: %s ulala\n", res.Balance.Amount)
}
```

---

## Python (cosmpy)

```python
from cosmpy.aerial.client import LedgerClient, NetworkConfig

# Connect
cfg = NetworkConfig(
    chain_id="lalachain-local",
    url="rest+http://localhost:1317",
    fee_minimum_gas_price=0.025,
    fee_denomination="ulala",
    staking_denomination="ulala",
)
client = LedgerClient(cfg)

# Query balance
balance = client.query_bank_balance("lala1abc...", "ulala")
print(f"Balance: {balance}")
```

---

## CLI (lalachaind)

The binary itself serves as an SDK for scripting:

```bash
# Query in JSON format
lalachaind query bank balances lala1abc... --output json

# Pipe to jq for processing
lalachaind query staking validators --output json | jq '.validators[].tokens'
```

---

## Choosing an SDK

| Language | Library | Best For |
|----------|---------|----------|
| TypeScript | @cosmjs/stargate | Web apps, React/Next.js |
| Go | cosmos-sdk/client | Backend services, modules |
| Python | cosmpy | Scripts, data analysis |
| Rust | cosmrs | Smart contracts, performance-critical |
| CLI | lalachaind | DevOps, automation |
