# View Transactions

**Query and track your transactions on LalaChain.**

---

## Query a Specific Transaction

```bash
# By transaction hash
lalachaind query tx <txhash>

# Via REST API
curl http://localhost:1317/cosmos/tx/v1beta1/txs/<txhash>
```

---

## Query Transactions by Address

```bash
# All transactions sent by an address
lalachaind query txs --events 'message.sender=lala1your-address'

# All transactions received by an address
lalachaind query txs --events 'transfer.recipient=lala1your-address'
```

---

## REST API Queries

```bash
# Get transactions by sender
curl "http://localhost:1317/cosmos/tx/v1beta1/txs?events=message.sender%3D%27lala1abc%27"

# Get transactions by height
curl "http://localhost:1317/cosmos/tx/v1beta1/txs?events=tx.height%3D100"
```

---

## Transaction Details

A transaction response includes:

| Field | Description |
|-------|-------------|
| `txhash` | Unique transaction identifier |
| `height` | Block number it was included in |
| `code` | 0 = success, non-zero = failure |
| `gas_wanted` | Gas limit set |
| `gas_used` | Actual gas consumed |
| `timestamp` | When the block was produced |
| `tx.body.messages` | What the transaction did |
| `logs` | Detailed execution events |

---

## Using the Dashboard

The LalaChain dashboard at `http://localhost:3000` provides:
- Real-time transaction feed
- Block explorer with transaction details
- Account balance history
- Network activity charts

---

**Next:** [Stake Tokens](stake-tokens.md)
