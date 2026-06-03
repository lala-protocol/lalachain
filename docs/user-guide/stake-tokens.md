---
title: "Stake Tokens"
description: "Delegate your LALA to validators and earn staking rewards."
---

# Stake Tokens

**Staking your LALA tokens helps secure the network and earns you rewards.**

---

## How to Stake (Delegate)

### Step 1: Choose a Validator

```bash
# List all validators
lalachaind query staking validators

# Look for:
# - High uptime
# - Reasonable commission (5-15%)
# - Active governance participation
# - Not jailed
```

### Step 2: Delegate Tokens

```bash
# Delegate 500 LALA to a validator
lalachaind tx staking delegate <validator-operator-address> 500000000ulala \
  --from my-wallet \
  --fees 10000ulala \
  --chain-id lalachain-1
```

### Step 3: Verify Delegation

```bash
# Check your delegations
lalachaind query staking delegations lala1your-address
```

---

## Managing Staked Tokens

### Claim Rewards

```bash
# Claim from all validators
lalachaind tx distribution withdraw-all-rewards \
  --from my-wallet \
  --fees 10000ulala

# Claim from specific validator
lalachaind tx distribution withdraw-rewards <validator-address> \
  --from my-wallet \
  --fees 10000ulala
```

### Redelegate (Move to Different Validator)

```bash
# Move stake from validator A to validator B (instant, no unbonding)
lalachaind tx staking redelegate <from-validator> <to-validator> 100000000ulala \
  --from my-wallet \
  --fees 10000ulala
```

Note: You can only redelegate from the same source validator once every 21 days.

### Unstake (Undelegate)

```bash
# Unstake 200 LALA (21-day unbonding period)
lalachaind tx staking unbond <validator-address> 200000000ulala \
  --from my-wallet \
  --fees 10000ulala
```

**Warning:** During the 21-day unbonding period:
- You earn NO rewards
- Your tokens CAN still be slashed
- You CANNOT cancel the unbonding

---

## Staking Tips

1. **Don't stake 100%** — keep some liquid for fees
2. **Diversify** — delegate to multiple validators to spread slashing risk
3. **Compound rewards** — regularly claim and re-stake rewards
4. **Monitor your validator** — redelegate away if uptime drops or they get jailed
5. **Check commission changes** — validators can increase commission rates

---

## Expected Returns

| Staking Amount | Est. Annual Reward (21% APR) | Daily |
|---------------|------------------------------|-------|
| 100 LALA | 21 LALA | 0.058 LALA |
| 1,000 LALA | 210 LALA | 0.575 LALA |
| 10,000 LALA | 2,100 LALA | 5.75 LALA |
| 100,000 LALA | 21,000 LALA | 57.5 LALA |

*Actual returns vary based on network inflation, staking ratio, and validator commission.*
