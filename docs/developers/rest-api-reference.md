# REST API Reference

**Complete documentation for LalaChain's REST API endpoints.**

---

## Base URLs

| Service | URL | Description |
|---------|-----|-------------|
| Cosmos REST | `http://localhost:1317` | Standard + custom endpoints |
| CometBFT RPC | `http://localhost:26657` | Block/tx queries |

---

## LalaChain Custom Endpoints

### GET /lala/telemetry/v1/kpis

Returns epoch KPI history.

**Response:**
```json
[
  {
    "epoch": 1,
    "avg_block_utilization": 0.35,
    "avg_base_fee": 950000000,
    "avg_block_time_ms": 5100,
    "tx_count": 12
  },
  {
    "epoch": 2,
    "avg_block_utilization": 0.42,
    "avg_base_fee": 880000000,
    "avg_block_time_ms": 5050,
    "tx_count": 18
  }
]
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `epoch` | int | Epoch number |
| `avg_block_utilization` | float | Gas used / gas limit (0-1) |
| `avg_base_fee` | int | Average base fee in ulala/gas |
| `avg_block_time_ms` | int | Average milliseconds between blocks |
| `tx_count` | int | Total transactions in epoch |

---

### GET /lala/aiadvisor/v1/state

Returns AI Advisor current state and configuration.

**Response:**
```json
{
  "low_streak": 2,
  "high_streak": 0,
  "last_proposal_epoch": 5,
  "config": {
    "min_fee_target": 800000000,
    "max_fee_target": 5000000000,
    "low_util_threshold": 0.40,
    "high_util_threshold": 0.80
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `low_streak` | int | Consecutive low-utilization epochs |
| `high_streak` | int | Consecutive high-utilization epochs |
| `last_proposal_epoch` | int | Last epoch a proposal was generated |
| `config.min_fee_target` | int | Below this fee → R4 may trigger |
| `config.max_fee_target` | int | Above this fee → R3 may trigger |
| `config.low_util_threshold` | float | Below this utilization → low streak increments |
| `config.high_util_threshold` | float | Above this utilization → high streak increments |

---

### GET /lala/lalagov/v1/history

Returns resolved governance proposals.

**Response:**
```json
[
  {
    "id": "prop-epoch-5-R4",
    "parameter": "base_fee_per_gas",
    "current_value": 750000000,
    "proposed_value": 787500000,
    "rule": "R4",
    "status": "passed",
    "epoch_created": 5,
    "epoch_resolved": 6,
    "votes_for": 1,
    "votes_against": 0,
    "votes_abstain": 0
  }
]
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique proposal identifier |
| `parameter` | string | Parameter being changed |
| `current_value` | int | Value before proposed change |
| `proposed_value` | int | Value after change (if approved) |
| `rule` | string | AI rule that triggered this (R1-R4) |
| `status` | string | "passed" or "rejected" |
| `epoch_created` | int | Epoch when proposed |
| `epoch_resolved` | int | Epoch when vote tallied |
| `votes_for` | int | YES votes (by voting power) |
| `votes_against` | int | NO votes |
| `votes_abstain` | int | ABSTAIN votes |

---

### GET /lala/lalagov/v1/config

Returns governance configuration.

**Response:**
```json
{
  "quorum": 0.66,
  "threshold": 0.51,
  "voting_period_epochs": 1,
  "activation_delay_epochs": 2
}
```

---

## Standard Cosmos SDK Endpoints

### Accounts & Balances

| Endpoint | Description |
|----------|-------------|
| `GET /cosmos/bank/v1beta1/balances/{address}` | Account balance |
| `GET /cosmos/bank/v1beta1/supply` | Total token supply |
| `GET /cosmos/auth/v1beta1/accounts/{address}` | Account info |

### Staking

| Endpoint | Description |
|----------|-------------|
| `GET /cosmos/staking/v1beta1/validators` | All validators |
| `GET /cosmos/staking/v1beta1/delegations/{address}` | Delegations for address |
| `GET /cosmos/staking/v1beta1/validators/{addr}/delegations` | Delegators for validator |

### Transactions

| Endpoint | Description |
|----------|-------------|
| `POST /cosmos/tx/v1beta1/txs` | Broadcast transaction |
| `GET /cosmos/tx/v1beta1/txs/{hash}` | Get transaction by hash |
| `GET /cosmos/tx/v1beta1/txs?events=...` | Search transactions |

### Governance

| Endpoint | Description |
|----------|-------------|
| `GET /cosmos/gov/v1beta1/proposals` | List proposals |
| `GET /cosmos/gov/v1beta1/proposals/{id}` | Get specific proposal |

---

## Error Responses

All endpoints return standard HTTP status codes:

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad request (invalid parameters) |
| 404 | Resource not found |
| 500 | Internal server error |

Error body format:
```json
{
  "code": 5,
  "message": "account not found",
  "details": []
}
```

---

## Rate Limiting

Default node configuration has no rate limiting. For production:
- Use a reverse proxy (nginx, Caddy) with rate limiting
- Limit by IP: 100 requests/second recommended
- Limit by endpoint: Heavy queries (tx search) should be stricter

---

**Next:** [SDK Guide](sdk-guide.md)
