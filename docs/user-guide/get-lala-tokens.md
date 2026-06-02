# Get LALA Tokens

**You need LALA tokens to pay transaction fees, stake, and participate in governance. Here's how to get them.**

---

## Methods to Acquire LALA

### 1. Testnet Faucet (Development)

For testnet usage, request tokens from the faucet:

```bash
# Request testnet tokens (if faucet is available)
curl -X POST https://faucet.lalachain.io/request \
  -H "Content-Type: application/json" \
  -d '{"address": "lala1your-address-here"}'
```

Or use the CLI genesis command for local development:

```bash
# Add tokens to an address in genesis (local testnet only)
lalachaind genesis add-genesis-account my-wallet 1000000000000ulala
```

### 2. Receive from Another User

Ask someone to send you tokens:

```bash
# The sender runs:
lalachaind tx bank send sender-wallet lala1your-address 1000000ulala \
  --fees 5000ulala
```

### 3. Decentralized Exchange (Mainnet)

Once LalaChain is live on mainnet with IBC connections:
- Trade on Osmosis DEX
- Bridge from other Cosmos chains
- Swap via integrated DEX interfaces

### 4. Staking Rewards

If you already have some LALA staked, you earn more through block rewards:

```bash
# Claim pending staking rewards
lalachaind tx distribution withdraw-all-rewards --from my-wallet --fees 5000ulala
```

---

## Check Your Balance

```bash
# CLI
lalachaind query bank balances lala1your-address

# REST API
curl http://localhost:1317/cosmos/bank/v1beta1/balances/lala1your-address
```

Response:
```json
{
  "balances": [
    {
      "denom": "ulala",
      "amount": "1000000000"
    }
  ]
}
```

Remember: `1000000000 ulala = 1000 LALA`

---

## Understanding Denominations

| You See | It Means |
|---------|----------|
| 1,000,000 ulala | 1 LALA |
| 5,000 ulala | 0.005 LALA |
| 1,000,000,000 ulala | 1,000 LALA |

The CLI and API always work in **ulala**. User interfaces typically display in **LALA**.

---

**Next:** [Send Tokens](send-tokens.md)
