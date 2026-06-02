# LalaChain API Reference

## REST API (LCD)

Default: `http://localhost:1317`

### Standard Cosmos Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /cosmos/bank/v1beta1/balances/{address}` | Account balances |
| `GET /cosmos/staking/v1beta1/validators` | All validators |
| `GET /cosmos/staking/v1beta1/pool` | Staking pool |
| `GET /cosmos/mint/v1beta1/inflation` | Current inflation |
| `GET /cosmos/mint/v1beta1/params` | Mint parameters |
| `GET /cosmos/distribution/v1beta1/community_pool` | Community pool |
| `GET /cosmos/slashing/v1beta1/params` | Slashing parameters |
| `GET /cosmos/base/tendermint/v1beta1/blocks/latest` | Latest block |
| `GET /cosmos/base/tendermint/v1beta1/blocks/{height}` | Block at height |

### Custom: Telemetry Module

| Endpoint | Description |
|----------|-------------|
| `GET /lala/telemetry/v1/kpis?epoch={n}` | KPIs for specific epoch |
| `GET /lala/telemetry/v1/kpis` | All KPI history |

#### Response: EpochKPIs

```json
{
  "kpis": {
    "epoch": "1",
    "avg_block_time_ms": "5012",
    "block_time_variance": "234.5",
    "avg_utilization": "0.65",
    "avg_base_fee": "1000000000",
    "validator_count": 4,
    "total_staked_ratio": "0.67",
    "mempool_size_trend": "0.1",
    "slashing_events": "0"
  }
}
```

### Custom: AI Advisor Module

| Endpoint | Description |
|----------|-------------|
| `GET /lala/aiadvisor/v1/state` | Current advisor state |

#### Response: AdvisorState

```json
{
  "state": {
    "next_proposal_id": "5",
    "low_util_streak": "2",
    "high_util_streak": "0"
  }
}
```

### Custom: LalaGov Module

| Endpoint | Description |
|----------|-------------|
| `GET /lala/lalagov/v1/config` | Governance config |
| `GET /lala/lalagov/v1/history` | Resolved proposal history |

#### Response: Config

```json
{
  "config": {
    "quorum": "0.66",
    "approval": "0.51",
    "voting_period_epochs": "1"
  }
}
```

#### Response: History

```json
{
  "proposals": [
    {
      "proposal": {
        "proposal_id": "1",
        "parameter": "block_gas_limit",
        "current_value": "15000000",
        "proposed_value": "15750000",
        "rationale": "low utilization (3 epochs): increase by 5%"
      },
      "votes_approve": "0.85",
      "votes_reject": "0.15",
      "outcome": "passed"
    }
  ]
}
```

## gRPC

Default: `localhost:9090`

### Custom Services

```protobuf
service lala.telemetry.v1.Query {
  rpc KPIs(QueryKPIsRequest) returns (QueryKPIsResponse);
  rpc AllKPIs(QueryAllKPIsRequest) returns (QueryAllKPIsResponse);
}

service lala.aiadvisor.v1.Query {
  rpc State(QueryStateRequest) returns (QueryStateResponse);
}

service lala.lalagov.v1.Query {
  rpc Config(QueryConfigRequest) returns (QueryConfigResponse);
  rpc History(QueryHistoryRequest) returns (QueryHistoryResponse);
}

service lala.lalagov.v1.Msg {
  rpc Vote(MsgVote) returns (MsgVoteResponse);
}
```

## Transaction Types

### MsgVote (LalaGov)

Vote on an AI-originated governance proposal.

```json
{
  "@type": "/lala.lalagov.v1.MsgVote",
  "proposal_id": "1",
  "voter": "lalavaloper1...",
  "approve": true
}
```

## WebSocket (Tendermint)

Default: `ws://localhost:26657/websocket`

### Subscribe to new blocks

```json
{
  "jsonrpc": "2.0",
  "method": "subscribe",
  "id": "1",
  "params": {
    "query": "tm.event='NewBlock'"
  }
}
```

### Subscribe to transactions

```json
{
  "jsonrpc": "2.0",
  "method": "subscribe",
  "id": "2",
  "params": {
    "query": "tm.event='Tx'"
  }
}
```

## SDK Integration (JavaScript/TypeScript)

```typescript
import { StargateClient } from "@cosmjs/stargate";

const client = await StargateClient.connect("http://localhost:26657");
const balance = await client.getBalance("lala1...", "ulala");
console.log(`Balance: ${parseInt(balance.amount) / 1e6} LALA`);
```

## Keplr Wallet Integration

```typescript
await window.keplr.experimentalSuggestChain({
  chainId: "lalachain-testnet-1",
  chainName: "LalaChain Testnet",
  rpc: "http://localhost:26657",
  rest: "http://localhost:1317",
  bip44: { coinType: 118 },
  bech32Config: {
    bech32PrefixAccAddr: "lala",
    bech32PrefixAccPub: "lalapub",
    bech32PrefixValAddr: "lalavaloper",
    bech32PrefixValPub: "lalavaloperpub",
    bech32PrefixConsAddr: "lalavalcons",
    bech32PrefixConsPub: "lalavalconspub",
  },
  currencies: [{ coinDenom: "LALA", coinMinimalDenom: "ulala", coinDecimals: 6 }],
  feeCurrencies: [{ coinDenom: "LALA", coinMinimalDenom: "ulala", coinDecimals: 6 }],
  stakeCurrency: { coinDenom: "LALA", coinMinimalDenom: "ulala", coinDecimals: 6 },
});
```
