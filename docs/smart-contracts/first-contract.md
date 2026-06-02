# First Contract

**Build a simple counter contract to learn CosmWasm fundamentals.**

---

## The Counter Contract

A basic contract that:
- Stores a counter value
- Allows incrementing the counter
- Allows querying the current count

---

## Step 1: Define Messages (msg.rs)

```rust
use cosmwasm_schema::{cw_serde, QueryResponses};

// Message sent when instantiating the contract
#[cw_serde]
pub struct InstantiateMsg {
    pub count: i32,
}

// Messages that modify state
#[cw_serde]
pub enum ExecuteMsg {
    Increment {},
    Reset { count: i32 },
}

// Messages that read state (no gas cost)
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(CountResponse)]
    GetCount {},
}

#[cw_serde]
pub struct CountResponse {
    pub count: i32,
}
```

---

## Step 2: Define State (state.rs)

```rust
use cw_storage_plus::Item;

pub const COUNT: Item<i32> = Item::new("count");
```

---

## Step 3: Implement Contract Logic (contract.rs)

```rust
use cosmwasm_std::{entry_point, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult, to_json_binary};
use crate::msg::{CountResponse, ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::COUNT;

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    COUNT.save(deps.storage, &msg.count)?;
    Ok(Response::new().add_attribute("method", "instantiate"))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> StdResult<Response> {
    match msg {
        ExecuteMsg::Increment {} => {
            COUNT.update(deps.storage, |count| -> StdResult<_> {
                Ok(count + 1)
            })?;
            Ok(Response::new().add_attribute("method", "increment"))
        }
        ExecuteMsg::Reset { count } => {
            COUNT.save(deps.storage, &count)?;
            Ok(Response::new().add_attribute("method", "reset"))
        }
    }
}

#[entry_point]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetCount {} => {
            let count = COUNT.load(deps.storage)?;
            to_json_binary(&CountResponse { count })
        }
    }
}
```

---

## Step 4: Build and Deploy

```bash
# Build
cargo build --release --target wasm32-unknown-unknown

# Upload
lalachaind tx wasm store target/wasm32-unknown-unknown/release/my_contract.wasm \
  --from my-wallet --gas auto --fees 50000ulala -y

# Instantiate with initial count of 0
lalachaind tx wasm instantiate 1 '{"count": 0}' \
  --from my-wallet --label "counter" --no-admin \
  --gas auto --fees 50000ulala -y
```

---

## Step 5: Interact

```bash
# Increment the counter
lalachaind tx wasm execute <contract-addr> '{"increment": {}}' \
  --from my-wallet --gas auto --fees 50000ulala -y

# Query the count
lalachaind query wasm contract-state smart <contract-addr> '{"get_count": {}}'
# Response: {"data":{"count":1}}

# Reset to 42
lalachaind tx wasm execute <contract-addr> '{"reset": {"count": 42}}' \
  --from my-wallet --gas auto --fees 50000ulala -y
```

---

## What You Learned

1. **InstantiateMsg** — Initial configuration when deploying
2. **ExecuteMsg** — State-changing operations (costs gas)
3. **QueryMsg** — Read-only operations (free)
4. **State** — Persistent storage using `cw_storage_plus`
5. **Entry points** — `instantiate`, `execute`, `query`

---

**Next:** [Testing](testing.md)
