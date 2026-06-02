# Evaluating AI Proposals as a Validator

**LalaChain validators have the unique responsibility of evaluating and voting on AI-generated parameter optimization proposals.**

---

## Your Role

Unlike other chains where validators mostly sign blocks, on LalaChain you are an active **governance participant**. Every epoch, you may need to evaluate proposals like:

> "AI Advisor proposes: Increase block_gas_limit by 5% because avg_block_utilization has been below 0.40 for 3 consecutive epochs."

Your job: decide if this change is good for the network.

---

## Proposal Anatomy

Every AI proposal contains:

| Field | What to Check |
|-------|---------------|
| `parameter` | Which parameter is being changed? |
| `current_value` | What is it now? |
| `proposed_value` | What will it become? |
| `magnitude` | How big is the change? (5-10%) |
| `rule` | Which AI rule triggered? (R1-R4) |
| `rationale` | Why does the AI think this is needed? |
| `evidence` | KPI data that triggered the rule |

---

## Decision Framework

### Vote YES when:

1. **Data supports the claim** — Verify the KPIs match your own observations
2. **Change is modest** — 5-10% adjustments are safe by design
3. **Pattern is sustained** — Streaks show a real trend, not a blip
4. **No external explanation** — The low/high utilization isn't caused by a known temporary event

### Vote NO when:

1. **External factor explains the anomaly** — Known maintenance, event, or market condition
2. **Recent change still settling** — A previous proposal just took effect; give it time
3. **Contradicts network goals** — The change moves against a known community preference
4. **Evidence looks wrong** — KPI data doesn't match your own node's observations

### Vote ABSTAIN when:

1. **Insufficient information** — You can't verify the data
2. **Genuinely uncertain** — Reasonable arguments both ways
3. **Conflict of interest** — You benefit disproportionately from the outcome

---

## Quick Evaluation Checklist

```
□ Is the parameter within safe bounds? (Check min/max limits)
□ Does the streak length justify action? (3+ for low, 2+ for high)
□ Is there a known external cause for the anomaly?
□ Has a recent proposal already addressed this area?
□ Will this change benefit users or hurt them?
□ Am I seeing the same pattern on my own node?
```

---

## Monitoring AI State

Keep an eye on the AI advisor's current state:

```bash
# Check current streaks and last proposal
curl http://localhost:1317/lala/aiadvisor/v1/state
```

If `low_streak` is approaching 3 or `high_streak` approaching 2, a proposal is likely imminent. Prepare to evaluate.

---

## Rule Reference

| Rule | Trigger | Proposal | Your Question |
|------|---------|----------|---------------|
| R1 | Low util 3 epochs + low fee | +5% gas limit | "Is the chain genuinely underused?" |
| R2 | High util 2 epochs | -5% gas limit | "Is the chain genuinely overloaded?" |
| R3 | Fee above 5B | -10% fee | "Are fees hurting users?" |
| R4 | Fee below 800M | +5% fee | "Are fees too low to sustain validators?" |

---

## Best Practices

1. **Don't auto-approve everything** — The AI is a tool, not an oracle
2. **Don't auto-reject everything** — The data is usually correct
3. **Communicate your reasoning** — Tell your delegators why you voted a certain way
4. **Track proposal outcomes** — Did past approvals improve things?
5. **Report anomalies** — If AI seems to be proposing incorrectly, alert the community

---

**Next:** Continue to [Governance: Overview](../governance/overview.md)
