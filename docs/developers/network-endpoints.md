---
title: "Network Endpoints"
description: "Available RPC, REST, and gRPC endpoints."
---

# Network Endpoints

**All available network endpoints for LalaChain nodes.**

---

## Default Ports

| Port | Protocol | Service | Access Level |
|------|----------|---------|-------------|
| 26656 | TCP | P2P (CometBFT) | Validators/nodes only |
| 26657 | HTTP/WS | RPC (CometBFT) | Public or restricted |
| 1317 | HTTP | REST API (Cosmos) | Public or restricted |
| 9090 | gRPC | gRPC interface | Applications |
| 26660 | HTTP | Prometheus metrics | Internal |

---

## CometBFT RPC (port 26657)

### Common Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/status` | GET | Node status, chain info |
| `/block` | GET | Latest block |
| `/block?height=N` | GET | Block at height N |
| `/blockchain?minHeight=A&maxHeight=B` | GET | Block range |
| `/tx?hash=0xABC` | GET | Transaction by hash |
| `/tx_search?query="..."` | GET | Search transactions |
| `/validators` | GET | Current validator set |
| `/net_info` | GET | Network peer info |
| `/health` | GET | Node health check |

### WebSocket

```
ws://localhost:26657/websocket
```

Subscribe to events:
```json
{
  "jsonrpc": "2.0",
  "method": "subscribe",
  "params": ["tm.event='NewBlock'"],
  "id": 1
}
```

---

## REST API (port 1317)

### LalaChain Custom

| Endpoint | Returns |
|----------|---------|
| `GET /lala/telemetry/v1/kpis` | Epoch KPI history |
| `GET /lala/aiadvisor/v1/state` | AI advisor state |
| `GET /lala/lalagov/v1/history` | Governance proposal history |
| `GET /lala/lalagov/v1/config` | Governance configuration |

### Cosmos Standard

| Category | Endpoint |
|----------|----------|
| **Bank** | `/cosmos/bank/v1beta1/balances/{address}` |
| **Bank** | `/cosmos/bank/v1beta1/supply` |
| **Auth** | `/cosmos/auth/v1beta1/accounts/{address}` |
| **Staking** | `/cosmos/staking/v1beta1/validators` |
| **Staking** | `/cosmos/staking/v1beta1/delegations/{address}` |
| **Distribution** | `/cosmos/distribution/v1beta1/delegators/{address}/rewards` |
| **Tx** | `POST /cosmos/tx/v1beta1/txs` |
| **Tx** | `GET /cosmos/tx/v1beta1/txs/{hash}` |

---

## gRPC (port 9090)

For high-performance applications. Supports all Cosmos SDK query services:

```bash
# List available services
grpcurl -plaintext localhost:9090 list

# Query balance
grpcurl -plaintext -d '{"address":"lala1abc...","denom":"ulala"}' \
  localhost:9090 cosmos.bank.v1beta1.Query/Balance
```

---

## Configuration

In `config/app.toml`:
```toml
[api]
enable = true
address = "tcp://0.0.0.0:1317"

[grpc]
enable = true
address = "0.0.0.0:9090"
```

In `config/config.toml`:
```toml
[rpc]
laddr = "tcp://0.0.0.0:26657"

[p2p]
laddr = "tcp://0.0.0.0:26656"
```

---

## Security Recommendations

- **P2P (26656):** Only expose to known peers/sentries
- **RPC (26657):** Restrict to trusted clients or use reverse proxy
- **REST (1317):** Rate-limit for public access
- **gRPC (9090):** Internal applications only
- **Metrics (26660):** Internal monitoring only
