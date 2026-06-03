---
title: "Voting Guide"
description: "How to vote on governance proposals as a validator or delegator."
---

# Voting Guide

**How to evaluate and vote on LalaChain governance proposals.**

---

## Voting Options

| Vote | Meaning | Effect |
|------|---------|--------|
| **YES** | Approve the proposal | Counts toward passage |
| **NO** | Reject the proposal | Counts against passage |
| **ABSTAIN** | No position | Counts toward quorum only |

---

## When to Vote YES

- The KPI evidence clearly shows a sustained problem
- The proposed change is within safe bounds (5-10%)
- No external factors explain the anomaly
- Previous similar proposals had good outcomes
- The change aligns with network health goals

## When to Vote NO

- External events explain the anomaly (planned maintenance, market event)
- A recent proposal already addressed this; give it time to take effect
- The change contradicts known community preferences
- You believe the AI is reacting to noise
- The KPI data seems inconsistent with your own observations

## When to Vote ABSTAIN

- You don't have enough information to decide
- You're genuinely uncertain about the impact
- The proposal is borderline and you defer to other validators

---

## Voting Best Practices

1. **Review every proposal** — Don't auto-vote on everything
2. **Check the evidence** — Verify KPIs against your own node data
3. **Consider context** — Is anything unusual happening in the broader ecosystem?
4. **Vote promptly** — Voting period is only 1 epoch (~50 seconds)
5. **Communicate** — Share your reasoning with delegators
6. **Track outcomes** — Review whether past votes improved the network

---

## Voting Power

Your vote weight equals your total stake:

```
Voting Power = Self-Bond + All Delegations
```

A validator with 10% of total stake has 10% of voting power.

---

## Quorum and Threshold

For a vote to be valid:
- **66% of total voting power must participate** (vote YES, NO, or ABSTAIN)
- Of those who voted YES or NO, **>51% must be YES** for passage

Example with 100M total staked:
- Quorum requires 66M in voting power to participate
- If 70M participates: 51% of (YES + NO) must be YES
