# Network Map

**A visual guide to LalaChain's network topology, data flows, and connectivity.**

---

## Network Topology

```mermaid
flowchart TB
    subgraph "Validator Set"
        V1[Validator 1<br/>Stake: 30%]
        V2[Validator 2<br/>Stake: 25%]
        V3[Validator 3<br/>Stake: 20%]
        V4[Validator 4<br/>Stake: 15%]
        V5[Validator 5<br/>Stake: 10%]
    end
    
    subgraph "Full Nodes"
        FN1[Full Node - Americas]
        FN2[Full Node - Europe]
        FN3[Full Node - Asia]
    end
    
    subgraph "Client Layer"
        W[Wallets]
        D[Dashboard]
        A[Applications]
    end
    
    V1 <--> V2
    V2 <--> V3
    V3 <--> V4
    V4 <--> V5
    V5 <--> V1
    V1 <--> V3
    
    V1 --> FN1
    V2 --> FN2
    V3 --> FN3
    
    FN1 --> W
    FN2 --> D
    FN3 --> A
```

---

## Node Types

| Node Type | Purpose | Count | Requirement |
|-----------|---------|-------|-------------|
| **Validator** | Produce blocks, vote on proposals | Active set (up to 100) | High uptime, stake |
| **Full Node** | Serve API queries, relay transactions | Unlimited | Moderate hardware |
| **Archive Node** | Store complete historical state | Few | Large storage |
| **Seed Node** | Help new nodes discover peers | 2-3 | Reliable connectivity |
| **Sentry Node** | DDoS protection for validators | Per validator | Behind validator |

---

## Communication Protocols

### Validator ↔ Validator (P2P)

- **Protocol:** CometBFT gossip
- **Port:** 26656 (default)
- **Data:** Blocks, votes, transactions, consensus messages
- **Security:** Authenticated encryption, node IDs

### Node ↔ Client (API)

- **Protocol:** HTTP REST
- **Port:** 1317 (API), 26657 (RPC)
- **Data:** Queries, transaction broadcasts
- **Security:** TLS recommended for production

---

## Data Flow: Transaction Lifecycle

```mermaid
sequenceDiagram
    participant User as User Wallet
    participant FN as Full Node
    participant MP as Mempool
    participant Prop as Block Proposer
    participant VS as Validator Set
    
    User->>FN: Submit transaction
    FN->>FN: Validate format & signature
    FN->>MP: Add to mempool
    MP->>Prop: Include in block proposal
    Prop->>VS: Broadcast proposed block
    VS->>VS: Pre-vote → Pre-commit
    VS->>FN: Commit block
    FN->>User: Transaction confirmed
```

---

## Data Flow: AI Governance Cycle

```mermaid
sequenceDiagram
    participant Blocks as Block Production
    participant Tel as Telemetry Module
    participant AI as AI Advisor
    participant Gov as Governance Module
    participant Vals as Validator Votes
    
    loop Every Block
        Blocks->>Tel: Record metrics
    end
    
    Note over Blocks,Vals: Epoch Boundary
    
    Tel->>AI: Epoch KPIs
    AI->>AI: Evaluate rules
    
    alt Rule triggered
        AI->>Gov: Create proposal
        Gov->>Vals: Open voting
        Vals->>Gov: Cast votes
        Gov->>Gov: Tally at epoch end
        
        alt Proposal passes
            Gov->>Blocks: Apply change (after delay)
        else Proposal fails
            Gov->>Gov: Archive as rejected
        end
    end
```

---

## Network Ports

| Port | Service | Access |
|------|---------|--------|
| 26656 | P2P (CometBFT) | Validators/nodes only |
| 26657 | RPC (CometBFT) | Public or restricted |
| 1317 | REST API (Cosmos) | Public or restricted |
| 26660 | Prometheus metrics | Internal monitoring |
| 9090 | gRPC | Applications |

---

## Geographic Distribution

For a healthy network, validators should be geographically distributed:

| Region | Recommended Validators | Latency to Others |
|--------|----------------------|-------------------|
| North America | 30-40% | <100ms to EU |
| Europe | 30-40% | <150ms to Asia |
| Asia-Pacific | 20-30% | <200ms to NA |

Block time of ~5 seconds comfortably accommodates global validator distribution.

---

**Next:** Continue to [How It Works: Consensus](../how-it-works/consensus.md)
