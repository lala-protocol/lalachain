---
title: "Send Tokens"
description: "Send LALA tokens to another address using CLI or dashboard."
---

# Send Tokens

**Transfer LALA tokens to any address on LalaChain.**

---

## CLI Method

```bash
lalachaind tx bank send <from-key> <to-address> <amount> --fees <fee>
```

### Example

```bash
# Send 100 LALA (= 100,000,000 ulala) to another address
lalachaind tx bank send my-wallet lala1recipient-address-here 100000000ulala \
  --fees 5000ulala \
  --chain-id lalachain-1
```

### Confirmation

You'll see a transaction hash:
```
txhash: A1B2C3D4E5F6...
```

### Verify

```bash
# Check transaction status
lalachaind query tx A1B2C3D4E5F6...

# Check recipient balance
lalachaind query bank balances lala1recipient-address
```

---

## REST API Method

```bash
# 1. Create and sign the transaction
# 2. Broadcast via API
curl -X POST http://localhost:1317/cosmos/tx/v1beta1/txs \
  -H "Content-Type: application/json" \
  -d '{
    "tx_bytes": "<base64-encoded-signed-tx>",
    "mode": "BROADCAST_MODE_SYNC"
  }'
```

---

## Important Notes

- **Double-check the recipient address** — transactions are irreversible
- **Ensure sufficient balance** — you need the amount + fee
- **Fees are non-refundable** — even if the transaction fails
- **Amount is in ulala** — multiply LALA by 1,000,000

---

## Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `insufficient funds` | Balance too low | Check balance, reduce amount |
| `account sequence mismatch` | Nonce out of sync | Wait for pending tx or query correct sequence |
| `invalid address` | Wrong format | Ensure address starts with `lala1` |
| `out of gas` | Gas limit too low | Increase `--gas` flag |
