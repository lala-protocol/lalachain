---
title: "Validator Requirements"
description: "Hardware, software, and network requirements for running a validator."
---

# Validator Requirements

**Hardware, software, and operational requirements for running a LalaChain validator.**

---

## Hardware Requirements

### Minimum

| Component | Specification |
|-----------|--------------|
| CPU | 4 cores (x86_64) |
| RAM | 16 GB |
| Storage | 500 GB NVMe SSD |
| Network | 100 Mbps, stable connection |

### Recommended

| Component | Specification |
|-----------|--------------|
| CPU | 8+ cores (x86_64) |
| RAM | 32 GB |
| Storage | 1 TB NVMe SSD |
| Network | 1 Gbps, redundant connection |

---

## Software Requirements

| Software | Version |
|----------|---------|
| Operating System | Ubuntu 22.04+ / RHEL 8+ / Debian 12+ |
| Go | 1.21+ (for building from source) |
| lalachaind | Latest stable release |

---

## Network Requirements

- **Static IP** or reliable DNS
- **Low latency** to other validators (<200ms recommended)
- **Uptime** target: 99.9% (max ~8h downtime/year)
- **Ports:** 26656 (P2P), 26657 (RPC, optional public)
- **DDoS protection** recommended (sentry node architecture)

---

## Staking Requirements

| Parameter | Value |
|-----------|-------|
| Minimum self-bond | Configurable (governance-set) |
| Commission range | 0-100% (recommend 5-20%) |
| Top N validators | Active set limited (e.g., top 100 by stake) |

---

## Operational Requirements

### Must Have
- 24/7 monitoring with alerting
- Automated restart on crash
- Key backup and recovery plan
- Software update process
- Understanding of AI governance proposals

### Should Have
- Sentry node architecture (DDoS protection)
- Hardware Security Module (HSM) for validator keys
- Geographic redundancy
- On-call rotation for incidents

---

## LalaChain-Specific Requirements

Unlike other chains, LalaChain validators must also:

1. **Understand AI proposals** — Evaluate parameter change proposals every epoch
2. **Vote actively** — Governance participation expected
3. **Monitor proposal quality** — Identify if AI generates unexpected proposals
4. **Communicate with community** — Explain voting decisions

---

## Estimated Costs

| Deployment | Monthly Cost | Suitable For |
|-----------|-------------|-------------|
| Cloud (basic) | $200-400 | New validators |
| Cloud (production) | $500-1,000 | Established validators |
| Bare metal + colo | $300-600 | Large operations |
