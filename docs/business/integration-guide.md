# Integration Guide

**How to integrate your application with LalaChain using the REST API.**

---

## API Overview

LalaChain exposes standard Cosmos SDK endpoints plus custom governance endpoints:

| Base URL | Default |
|----------|---------|
| REST API | `http://localhost:1317` |
| RPC | `http://localhost:26657` |
| gRPC | `localhost:9090` |

---

## Authentication

LalaChain's API is **publicly readable** (no API key needed for queries). Write operations (transactions) require a signed transaction from the sender's private key.

---

## Core Integration Patterns

### Pattern 1: Read Chain State

```javascript
// Check a user's balance
const response = await fetch(
  `${API_URL}/cosmos/bank/v1beta1/balances/${address}`
);
const { balances } = await response.json();
const lalaBalance = balances.find(b => b.denom === 'ulala');
```

### Pattern 2: Submit Transactions

```javascript
// 1. Build the transaction
// 2. Sign with user's private key
// 3. Broadcast
const response = await fetch(`${API_URL}/cosmos/tx/v1beta1/txs`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    tx_bytes: signedTxBase64,
    mode: 'BROADCAST_MODE_SYNC'
  })
});
```

### Pattern 3: Monitor Network Health

```javascript
// Get LalaChain-specific KPIs
const kpis = await fetch(`${API_URL}/lala/telemetry/v1/kpis`).then(r => r.json());
const latest = kpis[kpis.length - 1];

console.log(`Utilization: ${(latest.avg_block_utilization * 100).toFixed(1)}%`);
console.log(`Base Fee: ${(latest.avg_base_fee / 1e6).toFixed(2)} ulala`);
```

---

## LalaChain-Specific Endpoints

| Endpoint | Method | Returns |
|----------|--------|---------|
| `/lala/telemetry/v1/kpis` | GET | Epoch KPI history |
| `/lala/aiadvisor/v1/state` | GET | AI advisor state and config |
| `/lala/lalagov/v1/history` | GET | Governance proposal history |
| `/lala/lalagov/v1/config` | GET | Governance parameters |

---

## SDK Options

| Language | Library | Status |
|----------|---------|--------|
| JavaScript/TypeScript | `@cosmjs/stargate` | Available |
| Go | Cosmos SDK client | Available |
| Python | `cosmpy` | Available |
| Rust | `cosmrs` | Available |

---

## Integration Checklist

- [ ] Node accessible (API + RPC ports open)
- [ ] Can query balances
- [ ] Can build and sign transactions
- [ ] Can broadcast transactions
- [ ] Error handling for common failure modes
- [ ] Fee estimation using current base fee
- [ ] Transaction confirmation polling

---

**Next:** [Cost Analysis](cost-analysis.md)
