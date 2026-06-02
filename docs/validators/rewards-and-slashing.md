# Rewards & Slashing

**How validators earn rewards and what happens when they misbehave.**

---

## Reward Sources

| Source | Distribution |
|--------|-------------|
| Block inflation | Proportional to voting power |
| Transaction fees | Proposer bonus + proportional split |
| Commission | Percentage of delegator rewards |

### Reward Calculation

```
Validator Annual Reward = (Network Inflation × Validator Stake / Total Stake) 
                        + (Transaction Fees × Validator Share)
                        + (Commission × Delegator Rewards)
```

---

## Commission

Validators earn commission on all delegator rewards:

```bash
# Set during validator creation
--commission-rate 0.10          # 10% commission
--commission-max-rate 0.20      # Can never exceed 20%
--commission-max-change-rate 0.01  # Can change max 1% per day
```

| Commission | Validator Keeps | Delegators Get |
|-----------|----------------|----------------|
| 5% | 5% of delegation rewards | 95% |
| 10% | 10% of delegation rewards | 90% |
| 20% | 20% of delegation rewards | 80% |

---

## Slashing Conditions

### Downtime (Liveness Fault)

| Parameter | Value |
|-----------|-------|
| Window | 100 blocks |
| Missed blocks threshold | 95% (miss 95 of 100) |
| Slash amount | 0.01% of stake |
| Jail duration | 10 minutes |
| Recovery | Unjail transaction after jail period |

### Double-Signing (Safety Fault)

| Parameter | Value |
|-----------|-------|
| Detection | Automatic (evidence in blocks) |
| Slash amount | 5% of stake |
| Jail duration | Permanent (tombstoned) |
| Recovery | None — validator permanently removed |

---

## Slashing Impact on Delegators

**Delegators are slashed proportionally:**

If a validator with 1M LALA stake (including delegations) is slashed 5%:
- All delegators lose 5% of their delegation
- Validator loses 5% of self-bond
- Total: 50,000 LALA destroyed

This is why choosing a reliable validator matters.

---

## Maximizing Rewards

1. **High uptime** — Don't miss blocks, don't get jailed
2. **Attract delegations** — More stake = more rewards
3. **Reasonable commission** — Too high loses delegators, too low is unsustainable
4. **Active governance** — Builds reputation, attracts delegates
5. **Quick upgrades** — Don't miss blocks during network upgrades

---

## Reward Claiming

```bash
# Validator claims their own rewards
lalachaind tx distribution withdraw-rewards $(lalachaind keys show validator --bech val -a) \
  --from validator --commission --fees 10000ulala -y
```

The `--commission` flag also withdraws accumulated commission.

---

**Next:** [Evaluating AI Proposals](ai-proposals-as-validator.md)
