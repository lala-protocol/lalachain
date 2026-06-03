---
title: "Vote on Proposals"
description: "Participate in governance by voting on proposals."
---

# Vote on Proposals

**Validators vote on AI-generated proposals to approve or reject parameter changes. This guide explains the process.**

---

## Who Can Vote?

Currently, only **active validators** can vote on LalaChain governance proposals. If you're a delegator, your stake's voting power is exercised by your validator.

---

## Viewing Active Proposals

### Dashboard
Visit `http://localhost:3000` to see active proposals with:
- Parameter being changed
- Current vs. proposed value
- AI rationale
- Voting status
- Time remaining

### CLI
```bash
# View governance history
curl http://localhost:1317/lala/lalagov/v1/history
```

### API
```bash
# Get proposal history with vote details
curl http://localhost:1317/lala/lalagov/v1/history
```

---

## Understanding a Proposal

Each AI proposal contains:

```json
{
  "id": "prop-epoch-15-R1",
  "parameter": "block_gas_limit",
  "current_value": 10000000,
  "proposed_value": 10500000,
  "direction": "increase",
  "magnitude": "5%",
  "rule": "R1",
  "rationale": "avg_block_utilization below 0.40 for 3 consecutive epochs",
  "evidence": {
    "utilization_streak": [0.32, 0.28, 0.35],
    "avg_base_fee": 750000000
  }
}
```

### Questions to Ask Before Voting

1. **Is the data accurate?** — Does the KPI evidence match what you observe?
2. **Is the change appropriate?** — Does the magnitude (5-10%) seem reasonable?
3. **Are there external factors?** — Is low utilization due to a known temporary event?
4. **What are the downstream effects?** — Will this change affect fees or user experience?

---

## Casting a Vote (Validators Only)

Voting happens automatically through the validator's consensus participation. In the current implementation, validators signal their vote via the governance module:

```bash
# Vote YES on a proposal (validator only)
lalachaind tx gov vote <proposal-id> yes \
  --from validator-key \
  --fees 5000ulala

# Vote NO
lalachaind tx gov vote <proposal-id> no \
  --from validator-key \
  --fees 5000ulala

# Abstain (counts toward quorum but not approval)
lalachaind tx gov vote <proposal-id> abstain \
  --from validator-key \
  --fees 5000ulala
```

---

## Vote Outcomes

| Outcome | Condition | Result |
|---------|-----------|--------|
| **Passed** | Quorum met (≥66%) AND approval (>51% Yes) | Change activates after 2-epoch delay |
| **Rejected** | Quorum met but <51% Yes | No change; AI may re-propose later |
| **Expired** | Quorum NOT met within voting period | No change; treated as rejected |

---

## After a Vote Passes

1. **Activation delay** — 2 epochs (~100 seconds) safety buffer
2. **Change applied** — Parameter updated at epoch boundary
3. **Visible on dashboard** — Historical record maintained
4. **AI adjusts** — Streak counters reset; AI monitors new baseline

---

## Governance Philosophy

- **When to vote YES:** The data supports the change, and the magnitude is safe
- **When to vote NO:** External factors explain the anomaly, or the change is premature
- **When to ABSTAIN:** You don't have enough information to make an informed decision

The AI is a tool — it provides data-driven suggestions, but validators are the final decision-makers.
