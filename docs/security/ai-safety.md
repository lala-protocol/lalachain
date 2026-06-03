---
title: "AI Safety"
description: "Safety guarantees ensuring the AI Advisor cannot harm the network."
---

# AI Safety

**How LalaChain ensures the AI Advisor cannot harm the network, even in adversarial scenarios.**

---

## AI Safety Principles

1. **Bounded authority** — The AI can only propose, never execute
2. **Bounded parameters** — Hard limits on all adjustable values
3. **Bounded magnitude** — Maximum 5-10% change per proposal
4. **Human gating** — Every proposal requires validator vote
5. **Deterministic behavior** — No learning, no adaptation, no external inputs
6. **Full transparency** — Every decision is explainable and reproducible

---

## Threat Scenarios

### Scenario 1: AI Proposes Dangerous Value

**Risk:** Rule logic has a bug that proposes an extreme value.

**Mitigations:**
- Hard parameter bounds (min/max) catch extreme proposals
- Maximum ±5-10% per proposal limits damage
- Validator review catches anomalous proposals
- Activation delay provides time to react

**Worst case:** One parameter moves 5-10% in wrong direction. Reversed next epoch.

---

### Scenario 2: Sustained Bad Proposals

**Risk:** AI consistently proposes harmful changes that validators keep approving.

**Mitigations:**
- Hard bounds eventually stop the drift (e.g., gas limit can't go below 10M)
- Community monitoring detects degradation
- Validator reputation suffers for rubber-stamping bad proposals
- Governance can modify AI configuration parameters

---

### Scenario 3: AI Manipulation via Data

**Risk:** Attacker manipulates chain activity to trick AI into bad proposals.

**Example:** Artificially inflate block utilization to trigger gas limit reduction.

**Mitigations:**
- Streak requirements (2-3 epochs sustained) make manipulation expensive
- Each epoch costs real gas fees to fill blocks
- Validators can recognize artificial patterns
- The resulting change is only ±5% (reversible)

---

### Scenario 4: Code Vulnerability in Rule Engine

**Risk:** Bug in the AI module causes non-deterministic behavior or consensus failure.

**Mitigations:**
- Code audit before deployment
- Extensive test coverage
- Deterministic logic only (no floating point, no randomness)
- Same code runs on all validators (consensus catches discrepancies)

---

## Safety Guarantees

| Guarantee | How It's Enforced |
|-----------|------------------|
| AI cannot move parameters beyond hard bounds | Clamping in code |
| AI cannot skip validator voting | Governance module enforces vote requirement |
| AI cannot execute changes instantly | Activation delay is protocol-enforced |
| AI behavior is identical on all nodes | Deterministic code + consensus verification |
| AI decisions are auditable | Evidence included in every proposal |

---

## Comparing AI Risk to Status Quo

| Risk Source | Traditional Chains | LalaChain |
|-------------|-------------------|-----------|
| Bad parameters set by humans | Months to fix | AI proposes fix in epochs |
| Governance capture | Political (hard to detect) | Transparent (all votes on-chain) |
| Parameter never updated | Common (inertia) | Impossible (AI monitors continuously) |
| Catastrophic parameter change | Possible (no bounds) | Bounded (max ±10% per change) |

LalaChain's AI governance arguably **reduces** total governance risk compared to human-only systems.

---

## Emergency Shutdown

If the AI is determined to be malfunctioning:

1. **Validators vote NO** on all proposals (immediate)
2. **Governance proposal** to modify AI configuration (adjusts thresholds)
3. **Software upgrade** to fix or disable the rule engine (requires coordinated upgrade)
4. **Chain halt** as absolute last resort (1/3+ validators stop)

