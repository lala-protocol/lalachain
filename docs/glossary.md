---
title: "glossary"
description: "Definitions of key terms used in LalaChain documentation."
---

# Glossary

**Definitions of key terms used throughout the LalaChain documentation.**

---

## A

### Account
A unique identity on LalaChain, derived from a public/private key pair. Identified by an address starting with `lala1`.

### Activation Delay
A 2-epoch waiting period between a governance proposal passing and the parameter change taking effect. Gives validators and users time to prepare.

### AI Advisor
LalaChain's deterministic rule engine that monitors network KPIs and proposes parameter adjustments. Cannot act without validator approval.

### APY (Annual Percentage Yield)
The annualized return rate for staking, including compound effects from reinvesting rewards.

---

## B

### Base Fee
The minimum fee per unit of gas required for transaction inclusion. Adjusts dynamically based on block utilization (EIP-1559 style).

### Block
A batch of transactions bundled together, proposed by a validator, and agreed upon by consensus. Each block references the previous block's hash.

### Block Gas Limit
The maximum total gas that can be consumed by all transactions in a single block. Range: 10M-30M on LalaChain.

### Block Height
The sequential number of a block in the chain. Genesis is block 0 (or 1).

### Block Time
The target interval between blocks. LalaChain targets ~5 seconds per block.

### Bond
To lock tokens as collateral for validating or delegating. Bonded tokens earn rewards but cannot be transferred.

### Byzantine Fault Tolerance (BFT)
The ability of a distributed system to function correctly even when up to 1/3 of participants are malicious or faulty.

---

## C

### CometBFT
The consensus engine used by LalaChain (formerly called Tendermint). Provides instant finality and BFT guarantees.

### Commission
The percentage of delegator rewards that a validator keeps as payment for operating their node.

### Consensus
The process by which validators agree on the next block to add to the chain. Requires 2/3+ of voting power.

### Cosmos SDK
The open-source framework used to build LalaChain. Provides modular architecture for blockchain development.

---

## D

### Decay Factor
A parameter in the fee formula that determines how quickly the base fee decreases when blocks are underutilized. Formula: `baseFee * 7 / (7 + decayFactor)`.

### Delegation
The act of staking tokens through an existing validator without running your own node.

### Deterministic
Producing the same output given the same input, every time. Critical for consensus — all validators must compute identical results.

### Double-Signing
A validator signing two different blocks at the same height. A severe offense resulting in slashing and potential tombstoning.

---

## E

### EIP-1559
An Ethereum fee mechanism that LalaChain adapts. Uses a dynamic base fee that adjusts based on block utilization, plus optional tips.

### Epoch
A fixed period of 10 blocks on LalaChain. The AI Advisor evaluates network metrics and may propose changes at each epoch boundary.

### Evidence
Proof of validator misbehavior (e.g., double-signing) submitted to the chain for slashing.

---

## F

### Fee
Payment in ulala attached to a transaction to compensate validators for processing it. Must meet or exceed the base fee.

### Finality
The guarantee that a confirmed transaction cannot be reversed. LalaChain provides instant finality — once a block is committed, it's permanent.

### Full Node
A node that stores the full blockchain state and validates all transactions but doesn't participate in block production.

---

## G

### Gas
A unit measuring computational work required to process a transaction. More complex operations require more gas.

### Genesis
The first block (or configuration) of a blockchain. Contains initial parameters, accounts, and validator set.

### Governance
The on-chain system for proposing and voting on changes to the network. Includes both human-initiated and AI-initiated proposals.

---

## H

### Hard Bounds
Absolute minimum and maximum values for parameters that cannot be exceeded, regardless of governance votes or AI proposals.

### Hash
A fixed-length fingerprint of data, computed by a cryptographic function. Used to link blocks and verify data integrity.

### HSM (Hardware Security Module)
A physical device that stores cryptographic keys and performs signing operations in a tamper-resistant environment.

---

## I

### IBC (Inter-Blockchain Communication)
A Cosmos protocol for transferring tokens and data between independent blockchains.

### Inflation
The creation of new tokens over time, typically as staking rewards. Incentivizes network security.

---

## J

### Jailing
Temporary removal of a misbehaving validator from the active set. Can be reversed by the validator after fixing the issue and sending an unjail transaction.

---

## K

### Keeper
In Cosmos SDK architecture, a module's core logic component that manages state reads/writes and business rules.

### Key Management Service (KMS)
A separate system (like `tmkms`) that holds validator signing keys on a dedicated secure machine.

### KPI (Key Performance Indicator)
Network metrics tracked by the telemetry module: block utilization, fees, streaks, and gas usage.

---

## L

### LALA
The native token of LalaChain. Used for fees, staking, and governance. 1 LALA = 1,000,000 ulala.

### Layer 1 (L1)
A base blockchain that processes and finalizes transactions on its own chain (as opposed to Layer 2 solutions built on top).

---

## M

### Mempool
The waiting area where valid transactions sit before being included in a block by a validator.

### Module
A self-contained unit of blockchain logic in Cosmos SDK. LalaChain has custom modules: telemetry, aiadvisor, and gov.

---

## N

### Node
A computer running the LalaChain software, maintaining a copy of the blockchain state.

### Nonce
See "Sequence Number." A counter that prevents transaction replay.

---

## P

### Persistent Peers
Nodes that your node maintains constant connections to, configured in `config.toml`.

### Proposal
A governance action submitted for validator vote. Can be human-initiated (text, parameter change, upgrade) or AI-initiated (parameter adjustment).

### Pruning
Discarding historical state to save disk space. Pruned nodes can't serve historical queries but require less storage.

---

## Q

### Quorum
The minimum voting power (66% on LalaChain) that must participate for a vote to be valid.

---

## R

### RPC (Remote Procedure Call)
An interface for querying blockchain state and submitting transactions. LalaChain exposes REST (port 1317), Tendermint RPC (port 26657), and gRPC (port 9090).

### Rule Engine
The deterministic logic in the AI Advisor module that evaluates KPIs against thresholds and generates proposals. Four rules: R1-R4.

---

## S

### Seed Phrase (Mnemonic)
A 24-word phrase that can regenerate your private key. The master backup for your wallet.

### Sentry Node
A public-facing node that shields a validator from direct internet exposure. Part of the validator security architecture.

### Sequence Number
A per-account counter that increments with each transaction, preventing replay attacks.

### Slashing
The destruction of staked tokens as punishment for validator misbehavior (double-signing or extended downtime).

### Staking
Locking tokens as collateral to participate in consensus (as validator or delegator) and earn rewards.

### State Sync
A method for new nodes to quickly sync by downloading a recent state snapshot rather than replaying all historical blocks.

### Streak
The count of consecutive epochs where a metric exceeds a threshold. The AI Advisor requires sustained streaks before proposing changes.

---

## T

### Tombstoning
Permanent removal of a validator for severe misbehavior (double-signing). Cannot be reversed — requires creating a new validator.

### Transaction (Tx)
A signed message submitted to the blockchain requesting a state change (transfer, stake, vote, etc.).

---

## U

### ulala
The smallest denomination of LALA. 1 LALA = 1,000,000 ulala. All on-chain amounts are expressed in ulala.

### Unbonding Period
The 21-day waiting period after initiating unstaking, during which tokens earn no rewards and cannot be transferred.

### Utilization
The percentage of block gas limit actually used by transactions. Key metric for the AI Advisor's decisions.

---

## V

### Validator
A node operator who stakes tokens, produces blocks, and votes on governance proposals. Earns rewards for securing the network.

### Voting Period
The time window during which validators can cast votes on a proposal. One epoch on LalaChain.

### Voting Power
A validator's influence in consensus and governance, proportional to their total stake (self-bonded + delegated).

---

## W

### Wallet
Software that manages your private keys and allows you to sign transactions on LalaChain.

### WebAssembly (Wasm)
A binary instruction format used by CosmWasm smart contracts. Rust code compiles to Wasm for execution on-chain.
