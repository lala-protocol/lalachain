# AI Advisor Transparency

**Every AI decision on LalaChain is fully explainable, verifiable, and auditable.**

---

## Transparency Principles

1. **No black boxes** — The AI is a deterministic rule engine with published rules
2. **Full evidence** — Every proposal includes the data that triggered it
3. **Reproducible** — Any validator can independently verify why a proposal was generated
4. **Auditable** — Complete history of all proposals, votes, and outcomes is on-chain
5. **Open source** — Rule engine code is publicly available

---

## What You Can Verify

### 1. Rule Logic

The complete rule set is defined in source code:

```go
// R1: Low utilization + low fees → increase gas limit
if lowStreak >= 3 && avgBaseFee < MinFeeTarget {
    propose(IncreaseGasLimit, 5%)
}

// R2: High utilization → decrease gas limit
if highStreak >= 2 {
    propose(DecreaseGasLimit, 5%)
}

// R3: Fee too high → decrease fee
if avgBaseFee > MaxFeeTarget {
    propose(DecreaseFee, 10%)
}

// R4: Fee too low → increase fee
if avgBaseFee < MinFeeTarget {
    propose(IncreaseFee, 5%)
}
```

### 2. Input Data

Every proposal includes the KPI evidence that triggered it:

```json
{
  "evidence": {
    "avg_block_utilization_history": [0.32, 0.28, 0.35],
    "avg_base_fee": 750000000,
    "low_streak": 3,
    "high_streak": 0
  }
}
```

### 3. Deterministic Output

Given the same inputs, the same proposal is always generated. You can verify:
- Take the KPI history from `/lala/telemetry/v1/kpis`
- Run it through the published rules
- Confirm the proposal matches what the AI generated

### 4. Vote Records

Every vote is publicly recorded:

```json
{
  "proposal_id": "prop-epoch-15-R1",
  "votes": [
    {"validator": "lalavaloper1abc...", "vote": "yes", "power": 5000000},
    {"validator": "lalavaloper1xyz...", "vote": "no", "power": 3000000}
  ]
}
```

---

## Verification Process

Anyone can verify a proposal:

```bash
# 1. Get the proposal details
curl http://localhost:1317/lala/lalagov/v1/history

# 2. Get the KPI data for the relevant epochs
curl http://localhost:1317/lala/telemetry/v1/kpis

# 3. Get the AI configuration
curl http://localhost:1317/lala/aiadvisor/v1/state

# 4. Manually check: do the KPIs + config → produce this proposal?
```

---

## What the AI Cannot Do

| Action | Possible? |
|--------|-----------|
| Execute changes without vote | NO |
| Propose values outside bounds | NO |
| Access external data | NO |
| Learn or adapt its rules | NO (without code upgrade) |
| Hide its reasoning | NO (evidence on-chain) |
| Vote on its own proposals | NO |
| Ignore validator rejection | NO |

---

## Comparison to Other AI Governance

| System | Transparency | Deterministic | Bounded |
|--------|-------------|---------------|---------|
| ChatGPT-style LLM | Low (black box) | No | No |
| ML-based optimization | Medium (feature importance) | No | Configurable |
| **LalaChain rule engine** | **High (fully explainable)** | **Yes** | **Yes** |

---

**Next:** [Governance Philosophy](governance-philosophy.md)
