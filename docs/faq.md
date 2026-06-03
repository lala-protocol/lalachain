---
title: "Frequently Asked Questions"
description: "Answers to frequently asked questions about LalaChain."
---

# Frequently Asked Questions

---

## Beginner Questions

### What is LalaChain?
LalaChain is a Layer 1 blockchain built on the Cosmos SDK that uses an AI Advisor to automatically tune network parameters like gas limits and fees. It combines decentralized governance with intelligent automation.

### Is LalaChain its own blockchain or built on another chain?
LalaChain is its own independent Layer 1 blockchain. It uses the Cosmos SDK framework and CometBFT consensus engine, but operates as a sovereign chain.

### What makes LalaChain different from other blockchains?
LalaChain's AI Advisor monitors network health metrics and proposes parameter adjustments automatically. Traditional blockchains require manual governance proposals for parameter changes, which are slow and often never happen.

### What is the LALA token used for?
LALA is used for gas fees (paying for transactions), staking (securing the network), and governance (voting on proposals including AI-suggested changes).

### How do I get LALA tokens?
On testnet, use the faucet. On mainnet (future), you'll be able to acquire tokens through exchanges, staking rewards, or direct purchase.

### Is LalaChain proof-of-work or proof-of-stake?
Proof-of-Stake. Validators stake LALA tokens as collateral and are selected to produce blocks based on their stake weight.

### How fast are transactions?
Transactions achieve finality in a single block (~5 seconds target block time). Once confirmed, they cannot be reversed.

### What are gas fees on LalaChain?
Gas fees are small payments in LALA that compensate validators for processing your transaction. LalaChain's AI Advisor helps keep fees reasonable by adjusting parameters based on network usage.

### Can the AI steal my tokens?
No. The AI Advisor has zero access to user funds. It can only propose changes to network parameters (like gas limits), and even those require validator votes to take effect.

### What is an epoch?
An epoch is a fixed period of 10 blocks (~50 seconds). The AI Advisor evaluates network metrics and potentially proposes changes at each epoch boundary.

### What is staking?
Staking means locking your LALA tokens to help secure the network. In return, you earn staking rewards (newly minted tokens + transaction fees).

### Can I lose my staked tokens?
Yes, through "slashing" — if the validator you delegate to misbehaves (double-signing or extended downtime), a portion of staked tokens can be destroyed. Choose validators carefully.

### What is delegation?
Delegation lets you stake your tokens through an existing validator without running your own node. You earn rewards proportional to your delegation.

### How do I vote on governance proposals?
Use the CLI or dashboard to submit your vote (Yes/No/Abstain/NoWithVeto) during the voting period. Your voting power equals your staked amount.

### What happens if I don't vote?
Your tokens still count toward quorum but don't influence the outcome. Active participation strengthens the network.

### What's the minimum amount to stake?
There is no protocol-enforced minimum for delegation. However, very small amounts may not be worthwhile given transaction fees for claiming rewards.

### How long does unstaking take?
The unbonding period is 21 days. During this time, tokens earn no rewards and cannot be transferred.

### Is LalaChain compatible with Ethereum?
Not directly. LalaChain uses Cosmos SDK architecture. However, future IBC (Inter-Blockchain Communication) bridges could enable cross-chain token transfers.

### What wallets support LalaChain?
Any Cosmos-compatible wallet (Keplr, Leap) can be configured for LalaChain. The native prefix is "lala".

### Is LalaChain decentralized?
Yes. Validators are independent operators, governance is on-chain, and the AI Advisor cannot act without validator approval.

---

## AI Advisor Questions

### Can the AI Advisor be manipulated?
The AI uses streak-based detection (sustained patterns over multiple epochs) which makes manipulation expensive. An attacker would need to sustain artificial activity for 2-3 full epochs while paying real gas fees.

### What happens if the AI proposes a bad change?
Validators vote on every proposal. If they vote NO, nothing changes. Even if a bad change passes, it's bounded to ±5-10% and can be reversed in the next epoch.

### Is the AI learning from the network?
No. The AI Advisor uses fixed deterministic rules, not machine learning. It evaluates KPIs against thresholds — the same inputs always produce the same outputs.

### Can the AI be upgraded?
Yes, through governance. A proposal to modify AI configuration parameters or upgrade the module can be submitted and voted on by validators.

### What metrics does the AI monitor?
Block utilization (% of gas limit used), fee levels, and streak counts (consecutive epochs of high or low usage).

### Are AI decisions transparent?
Completely. Every proposal includes the triggering rule, input data, and reasoning. Anyone can verify by checking the same KPIs.

### What parameters can the AI change?
- `block_gas_limit` (range: 10M-30M)
- `base_fee_per_gas` (range: 100M-10B)

### How often does the AI make proposals?
At most once per epoch (every 10 blocks). Many epochs produce no proposal if metrics are within normal ranges.

---

## Technical Questions

### What consensus algorithm does LalaChain use?
CometBFT (formerly Tendermint), which provides instant finality and Byzantine Fault Tolerance with up to 1/3 malicious validators.

### What programming language is LalaChain built in?
Go (Golang), using the Cosmos SDK v0.50.9 framework.

### Can I build smart contracts on LalaChain?
Yes. LalaChain supports CosmWasm smart contracts written in Rust, compiled to WebAssembly.

### What's the block size limit?
Configurable via the `block_gas_limit` parameter. Default range is 10M-30M gas units per block.

### How does fee calculation work?
LalaChain uses an EIP-1559-style mechanism with a base fee that adjusts based on block utilization. The formula includes a decay factor: `baseFee * 7 / (7 + decayFactor)`.

### What's the RPC endpoint?
- REST API: `localhost:1317`
- RPC: `localhost:26657`
- gRPC: `localhost:9090`

### How do I query AI Advisor state?
```
GET /lala/aiadvisor/v1/state
```

### How do I query current KPIs?
```
GET /lala/telemetry/v1/kpis
```

### What's the account prefix?
`lala` — addresses look like `lala1abc123...`

### What's the bond denomination?
`ulala` (micro-LALA). 1 LALA = 1,000,000 ulala.

### How do I run a full node without being a validator?
Run `lalachaind start` without creating a validator transaction. You'll sync blocks and serve queries but won't produce blocks.

### Can I run a node on Windows?
The chain binary compiles on Windows (lalachaind.exe). For production validators, Linux is recommended.

---

## Staking & Economics Questions

### What's the current APY for staking?
Depends on total staked amount, inflation parameters, and fee revenue. Check the dashboard or query the mint module.

### Do I earn rewards while unbonding?
No. Once you initiate unbonding, rewards stop immediately even though tokens are locked for 21 days.

### Can I redelegate without waiting?
Yes. You can redelegate from one validator to another instantly, but the destination validator inherits the original unbonding timer.

### What's the commission rate?
Each validator sets their own commission (0-100%). They take this percentage from delegator rewards. Lower commission = more rewards for delegators.

### Is there inflation?
Yes. New LALA tokens are minted each block as staking rewards, incentivizing network security. The inflation rate is governed by protocol parameters.

---

## Security Questions

### Has LalaChain been audited?
The codebase uses battle-tested components (Cosmos SDK, CometBFT) with extensive audit histories. Custom modules (telemetry, AI advisor, governance) require separate auditing before mainnet.

### What happens during a network attack?
CometBFT halts rather than producing incorrect blocks. Validators can coordinate response through social consensus and emergency governance.

### Can validators collude?
If 2/3+ of stake colludes, they could theoretically control the network. Economic incentives (slashing, reputation loss, token devaluation) make this irrational.

### Is my data private?
All transactions are public on the blockchain. LalaChain does not provide privacy features by default. Don't put sensitive data in transaction memos.
