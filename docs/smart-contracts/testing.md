---
title: "Testing Smart Contracts"
description: "Testing strategies for CosmWasm smart contracts."
---

# Testing Smart Contracts

**How to test CosmWasm contracts locally before deploying to LalaChain.**

---

## Unit Tests

Test individual functions in isolation:

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg { count: 17 };
        let info = mock_info("creator", &[]);
        
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());
        
        // Query the state
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCount {}).unwrap();
        let value: CountResponse = from_json(res).unwrap();
        assert_eq!(17, value.count);
    }

    #[test]
    fn increment() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg { count: 0 };
        let info = mock_info("creator", &[]);
        instantiate(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();
        
        // Increment
        let msg = ExecuteMsg::Increment {};
        execute(deps.as_mut(), mock_env(), info, msg).unwrap();
        
        // Should be 1
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCount {}).unwrap();
        let value: CountResponse = from_json(res).unwrap();
        assert_eq!(1, value.count);
    }
}
```

Run tests:
```bash
cargo test
```

---

## Integration Tests

Test contract interactions using `cw-multi-test`:

```rust
use cosmwasm_std::Addr;
use cw_multi_test::{App, ContractWrapper, Executor};

#[test]
fn multi_test_example() {
    let mut app = App::default();
    
    // Store contract code
    let code = ContractWrapper::new(execute, instantiate, query);
    let code_id = app.store_code(Box::new(code));
    
    // Instantiate
    let contract_addr = app
        .instantiate_contract(
            code_id,
            Addr::unchecked("creator"),
            &InstantiateMsg { count: 0 },
            &[],
            "Counter",
            None,
        )
        .unwrap();
    
    // Execute
    app.execute_contract(
        Addr::unchecked("user"),
        contract_addr.clone(),
        &ExecuteMsg::Increment {},
        &[],
    ).unwrap();
    
    // Query
    let resp: CountResponse = app
        .wrap()
        .query_wasm_smart(contract_addr, &QueryMsg::GetCount {})
        .unwrap();
    
    assert_eq!(resp.count, 1);
}
```

---

## Testing Checklist

- [ ] All entry points tested (instantiate, execute, query)
- [ ] Error cases tested (unauthorized, invalid input)
- [ ] Edge cases (zero values, max values, empty state)
- [ ] Multi-contract interactions (if applicable)
- [ ] Gas consumption reasonable

---

## Running Tests

```bash
# All tests
cargo test

# Specific test
cargo test proper_initialization

# With output
cargo test -- --nocapture
```
