---
title: "Error Codes"
description: "Error codes, meanings, and resolution steps."
---

# Error Codes

**Common error codes returned by LalaChain and how to resolve them.**

---

## Cosmos SDK Error Codes

| Code | Name | Description | Resolution |
|------|------|-------------|-----------|
| 2 | `ErrTxDecode` | Transaction encoding error | Ensure correct protobuf encoding |
| 3 | `ErrInvalidSequence` | Account sequence mismatch | Query current sequence, resend |
| 4 | `ErrUnauthorized` | Signature verification failed | Check private key matches sender |
| 5 | `ErrInsufficientFunds` | Not enough balance | Check balance, reduce amount |
| 6 | `ErrUnknownRequest` | Unknown message type | Check message type registration |
| 7 | `ErrInvalidAddress` | Malformed address | Verify bech32 format with `lala` prefix |
| 8 | `ErrInvalidPubKey` | Invalid public key | Regenerate key pair |
| 9 | `ErrUnknownAddress` | Account not found | Fund the account first |
| 10 | `ErrInvalidCoins` | Invalid coin denomination | Use `ulala` as denom |
| 11 | `ErrOutOfGas` | Exceeded gas limit | Increase gas limit |
| 12 | `ErrMemoTooLarge` | Memo exceeds max length | Shorten memo (max 256 bytes) |
| 13 | `ErrInsufficientFee` | Fee below minimum | Increase fee amount |
| 14 | `ErrTooManySignatures` | Too many signers | Reduce multisig participants |
| 19 | `ErrTxInMempoolCouldNotCheckTx` | Mempool rejection | Retry with higher fee |
| 20 | `ErrMempoolIsFull` | Mempool at capacity | Wait and retry |

---

## Transaction-Level Errors

### Sequence Mismatch (Code 3)

```json
{
  "code": 3,
  "message": "account sequence mismatch, expected 5, got 4"
}
```

**Cause:** You're using a stale sequence number (nonce).

**Fix:**
```bash
# Query current sequence
lalachaind query auth account lala1your-address --output json | jq '.sequence'
# Use the returned value in your next transaction
```

### Insufficient Funds (Code 5)

```json
{
  "code": 5,
  "message": "insufficient funds: 1000ulala is smaller than 5000ulala"
}
```

**Cause:** Account balance is less than amount + fees.

**Fix:** Reduce amount or acquire more tokens.

### Out of Gas (Code 11)

```json
{
  "code": 11,
  "message": "out of gas in location: ...; gasWanted: 100000, gasUsed: 100001"
}
```

**Cause:** Transaction needed more gas than the limit you set.

**Fix:** Increase gas limit:
```bash
lalachaind tx bank send ... --gas 300000
```

---

## Network-Level Errors

| Error | Meaning | Resolution |
|-------|---------|-----------|
| `connection refused` | Node not running | Start the node |
| `post failed: status 503` | Node syncing | Wait for sync to complete |
| `tx already exists in cache` | Duplicate transaction | Use different sequence number |
| `validator not found` | Invalid validator address | Check validator operator address |

---

## LalaChain-Specific Errors

| Error | Context | Meaning |
|-------|---------|---------|
| `proposal already active` | Governance | Can't submit when one is pending |
| `voting period ended` | Governance | Vote submitted too late |
| `parameter out of bounds` | AI Advisor | Proposed value exceeds safety limits |

---

## Debugging Tips

1. **Always check the `code` field** — 0 means success, anything else is an error
2. **Use `--output json`** for machine-readable errors
3. **Check `raw_log`** for detailed error messages
4. **Query account state** before building transactions
5. **Simulate first** with `--dry-run` flag to catch errors without spending gas

```bash
# Simulate transaction (no broadcast)
lalachaind tx bank send ... --dry-run
```

