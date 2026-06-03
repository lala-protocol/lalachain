---
title: "Network Security"
description: "How LalaChain defends against network-level attacks."
---

# Network Security

**How LalaChain protects against network-level attacks.**

---

## Attack Vectors & Mitigations

### 1. 2/3 Stake Attack

**Attack:** Acquire 2/3+ of staked tokens to control consensus.

**Mitigation:**
- Economic cost is prohibitive (requires billions in token value)
- Token distribution designed to prevent concentration
- Slashing discourages large-scale malicious staking
- Community monitoring for unusual stake accumulation

---

### 2. Long-Range Attack

**Attack:** Fork from a historical block using keys that have since unstaked.

**Mitigation:**
- Unbonding period (21 days) means old keys can still be slashed
- Weak subjectivity checkpoints
- Social consensus on canonical chain

---

### 3. DDoS

**Attack:** Overwhelm nodes with traffic to prevent block production.

**Mitigation:**
- Sentry node architecture shields validators
- Rate limiting on public endpoints
- Transaction fees make spam expensive
- Mempool size limits
- Geographic distribution of validators

---

### 4. Eclipse Attack

**Attack:** Surround a node with attacker-controlled peers, isolating it from the real network.

**Mitigation:**
- Persistent peer connections to known validators
- Seed nodes for discovery
- Node ID verification
- Multiple connection paths

---

### 5. Transaction Spam

**Attack:** Flood the network with low-value transactions to degrade performance.

**Mitigation:**
- Gas fees make spam expensive
- Mempool limits reject excess transactions
- Dynamic fee adjustment (EIP-1559) raises prices during spam
- AI Advisor detects and responds to utilization spikes

---

## Network Hardening

### P2P Layer
- Authenticated connections (node IDs)
- Encrypted communication
- Gossip protocol limits propagation of invalid data
- Connection limits per peer

### Consensus Layer
- BFT tolerates up to 1/3 malicious validators
- Evidence system detects and punishes misbehavior
- Instant finality prevents rollback attacks

### Application Layer
- Transaction validation before mempool acceptance
- State transition verification by all nodes
- Deterministic execution prevents consensus splits

---

## Security Assumptions

LalaChain's security holds when:
1. Less than 1/3 of stake is controlled by adversaries
2. Network can deliver messages within reasonable time
3. Cryptographic primitives (SHA-256, secp256k1, Ed25519) remain secure
4. At least one honest node is connected to you
