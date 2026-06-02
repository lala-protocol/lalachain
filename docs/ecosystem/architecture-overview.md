# Architecture Overview

**LalaChain is a modular Layer 1 blockchain built on Cosmos SDK with three custom modules that enable AI-driven governance.**

---

## High-Level Architecture

```mermaid
flowchart TB
    subgraph "Application Layer"
        Dashboard[Web Dashboard]
        CLI[CLI Client]
        API[REST API]
    end
    
    subgraph "LalaChain Node"
        subgraph "Custom Modules"
            TEL[x/telemetry]
            AI[x/aiadvisor]
            GOV[x/lalagov]
        end
        
        subgraph "Cosmos SDK Core"
            BANK[x/bank]
            STAKING[x/staking]
            AUTH[x/auth]
            DIST[x/distribution]
        end
        
        subgraph "Consensus Engine"
            CMT[CometBFT]
        end
    end
    
    subgraph "Network"
        V1[Validator 1]
        V2[Validator 2]
        V3[Validator N]
    end
    
    Dashboard --> API
    CLI --> API
    API --> TEL
    API --> AI
    API --> GOV
    TEL --> AI
    AI --> GOV
    CMT --> V1
    CMT --> V2
    CMT --> V3
```

---

## Layer Breakdown

### 1. Consensus Layer (CometBFT)

The foundation. CometBFT handles:
- Peer-to-peer networking between validators
- Block proposal and voting (Byzantine Fault Tolerant)
- Transaction mempool management
- Finality guarantees (no chain reorganizations)

**Protocol:** BFT consensus requiring 2/3+ validator agreement  
**Block time:** ~5 seconds  
**Finality:** Instant (single-slot finality)

### 2. Application Layer (Cosmos SDK)

The SDK provides standard blockchain functionality through modules:

| Module | Purpose |
|--------|---------|
| `x/auth` | Account management, transaction authentication |
| `x/bank` | Token transfers, balance tracking |
| `x/staking` | Validator set management, delegation |
| `x/distribution` | Reward distribution to validators/delegators |
| `x/gov` | (Base) governance infrastructure |

### 3. LalaChain Custom Layer

Three modules that provide the unique AI governance functionality:

| Module | Purpose | Runs Every |
|--------|---------|-----------|
| `x/telemetry` | Collects block metrics, computes KPIs | Every epoch (10 blocks) |
| `x/aiadvisor` | Evaluates KPIs against rules, generates proposals | Every epoch |
| `x/lalagov` | Manages proposal lifecycle, voting, activation | Every epoch |

---

## Data Flow

```mermaid
sequenceDiagram
    participant Block as New Block
    participant Tel as x/telemetry
    participant AI as x/aiadvisor
    participant Gov as x/lalagov
    participant Val as Validators
    
    loop Every block
        Block->>Tel: Record gas used, fees, block time
    end
    
    Note over Block,Val: End of Epoch (every 10 blocks)
    
    Tel->>Tel: Compute epoch KPIs (avg utilization, avg fee, etc.)
    Tel->>AI: Pass KPI data
    AI->>AI: Evaluate rule engine (streaks, thresholds)
    AI-->>Gov: Generate proposal (if rule triggered)
    Gov->>Val: Present proposal for voting
    Val->>Gov: Cast votes (Yes/No/Abstain)
    Gov->>Gov: Tally votes after voting period
    Gov-->>Block: Apply approved changes (after delay)
```

---

## Module Interactions

### x/telemetry → x/aiadvisor
- Passes computed KPIs at epoch end
- KPIs include: avg_block_utilization, avg_base_fee, avg_block_time, tx_count

### x/aiadvisor → x/lalagov
- Generates signed proposals when rules trigger
- Each proposal specifies: parameter to change, direction, magnitude, rationale

### x/lalagov → Chain State
- Applies approved parameter changes to the chain configuration
- Tracks proposal history for transparency

---

## State Storage

Each module maintains its own state in the application's key-value store:

| Module | State Contents |
|--------|---------------|
| `x/telemetry` | KPI history (per epoch), raw block metrics |
| `x/aiadvisor` | Rule configuration, streak counters, proposal log |
| `x/lalagov` | Active proposals, vote tallies, resolved proposals, config |

State is persisted via the Cosmos SDK `KVStore` and included in the app hash (Merkle root) committed with each block.

---

## API Layer

LalaChain exposes a REST API for external access:

| Endpoint | Returns |
|----------|---------|
| `GET /lala/telemetry/v1/kpis` | Historical KPI data |
| `GET /lala/aiadvisor/v1/state` | AI rule configuration and streak state |
| `GET /lala/lalagov/v1/history` | Resolved proposal history |
| `GET /lala/lalagov/v1/config` | Governance parameters |

Plus standard Cosmos SDK endpoints for accounts, balances, staking, etc.

---

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| Rule-based AI (not ML) | Deterministic, auditable, reproducible across nodes |
| Epoch-based analysis | Smooths out noise, prevents overreaction |
| Hard parameter bounds | Safety rails — AI can't propose dangerous values |
| Validator voting | Human oversight maintained; AI is advisory only |
| Cosmos SDK framework | Battle-tested, modular, IBC-compatible |

---

**Next:** [Components](components.md)
