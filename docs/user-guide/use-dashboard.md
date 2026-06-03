---
title: "Use the Dashboard"
description: "Navigate the LalaChain monitoring dashboard."
---

# Use the Dashboard

**The LalaChain dashboard is a web interface for monitoring chain performance, viewing proposals, and tracking governance activity.**

---

## Accessing the Dashboard

```
URL: http://localhost:3000
```

The dashboard is a Next.js application that connects to your LalaChain node's REST API.

---

## Dashboard Sections

### 1. KPIs (Key Performance Indicators)

Displays real-time epoch metrics:
- **Block Utilization** — How full blocks are (gas used / gas limit)
- **Base Fee** — Current fee per gas unit (in ulala)
- **Block Time** — Seconds between blocks
- **Transaction Count** — Transactions per epoch

Data source: `GET /lala/telemetry/v1/kpis`

### 2. Governance / Proposals

Shows proposal history:
- Proposal ID and type
- Parameter changed
- Before/after values
- Vote outcome (passed/rejected)
- Epoch created and resolved

Data source: `GET /lala/lalagov/v1/history`

### 3. AI Advisor State

Displays current AI status:
- Streak counters (low/high utilization streaks)
- Rule configuration
- Last proposal epoch
- Current thresholds

Data source: `GET /lala/aiadvisor/v1/state`

### 4. Network Status

Shows general chain info:
- Current block height
- Active validators
- Total staked tokens
- Latest epoch number

---

## Reading the Charts

### Utilization Chart
- **Green zone (40-80%)** — Healthy utilization
- **Yellow zone (<40%)** — Underutilized (AI may propose gas limit increase)
- **Red zone (>80%)** — Overloaded (AI may propose gas limit decrease)

### Fee Chart
- **Green zone (800M-5B ulala)** — Healthy fee range
- **Below 800M** — Too low (AI may propose fee increase)
- **Above 5B** — Too high (AI may propose fee decrease)

---

## Dashboard Configuration

The dashboard connects to these endpoints:

```env
NEXT_PUBLIC_API_URL=http://localhost:1317
NEXT_PUBLIC_RPC_URL=http://localhost:26657
```

To point at a different node, update these environment variables and restart the dashboard.

