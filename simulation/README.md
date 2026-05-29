# Lala Protocol — Adaptation Loop Simulation

A self-contained Python simulation that demonstrates the full four-layer
architecture described in the [Lala Protocol whitepaper](../lala-protocol-whitepaper.pdf):

1. **Base Consensus Layer (PoS)** — epoch-driven block production with
   configurable gas limits and fee parameters.
2. **Network Telemetry Module** — deterministically computes KPIs (block
   time, utilization, fees, slashing events) from synthetic block data.
3. **AI Advisory Module** — rule-based advisor that proposes parameter
   adjustments when KPIs fall outside target ranges.
4. **Governance Layer** — stake-weighted validator voting with configurable
   quorum (66 %) and approval (51 %) thresholds.

## Requirements

Python 3.8 or later. No external dependencies — the simulation uses only the
standard library (`dataclasses`, `statistics`, `random`, `argparse`).

## Running the simulation

```bash
# Default: 30 epochs
python simulation.py

# Custom epoch count
python simulation.py --epochs 60

# Reproducible run
python simulation.py --epochs 30 --seed 42
```

## What to expect

The simulation runs a synthetic demand curve:

| Epochs | Demand level | Expected behaviour |
|---|---|---|
| 1–10 | Low (≈30 %) | AI proposes block gas limit increase after 3 consecutive low-utilisation epochs |
| 11–20 | High spike (≈90 %) | AI proposes block gas limit decrease to ease congestion |
| 21–N | Moderate (≈55 %) | System stabilises; fewer or no proposals generated |

Each epoch prints live KPIs, proposal registrations, and governance outcomes.
A summary at the end lists all ratified parameter changes and aggregate KPI
statistics.

## Sample output (seed 42, 30 epochs)

```
======================================================================
 Lala Protocol — Adaptation Loop Simulation
======================================================================
Initial params: NetworkParams(block_gas_limit=15000000, ...)

── Epoch 001 ──────────────────────────────────────────────────────────
  KPIs  util=27.3%  fee=0.821 Gwei  blk_time=6003 ms  demand≈0.28
...
── Epoch 005 ──────────────────────────────────────────────────────────
  KPIs  util=31.1%  fee=0.798 Gwei  blk_time=5991 ms  demand≈0.31
    [Gov] Proposal #1 registered: block_gas_limit 15,000,000 → 15,750,000
...
    [Gov] Proposal #1 ✅ PASSED  (approve=78.3%, reject=12.1%)
    [Gov] Scheduled activation of block_gas_limit=15750000 at epoch 7
...
======================================================================
 SIMULATION SUMMARY
======================================================================
Final network parameters:
  block_gas_limit    : 14,250,000
  base_fee_per_gas   : 1.050 Gwei
  ...
```

## Relationship to the whitepaper

The simulation implements the core adaptation loop from **Figure 2** of the
whitepaper (`Lala Protocol Parameter Adaptation Loop`) and validates the
mechanisms described in **Section 3** (Core Mechanisms).

See [`../IMPLEMENTATION_FEASIBILITY.md`](../IMPLEMENTATION_FEASIBILITY.md) for
a full feasibility analysis and recommended production implementation roadmap.
