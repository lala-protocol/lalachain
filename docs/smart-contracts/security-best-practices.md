---
title: "Security Best Practices"
description: "Security patterns and common pitfalls in smart contract development."
---

# Security Best Practices

**Common vulnerabilities in smart contracts and how to prevent them on LalaChain.**

---

## CosmWasm Security Advantages

CosmWasm avoids many Solidity pitfalls by design:

| Vulnerability | Solidity Risk | CosmWasm |
|---------------|--------------|----------|
| Reentrancy | High | **Impossible** (actor model) |
| Integer overflow | Medium | **Caught** (Rust panics) |
| Unchecked return values | Medium | **Enforced** (Result types) |
| Access control | High | **Explicit** (info.sender) |

---

## Best Practices

### 1. Validate All Inputs

```rust
pub fn execute_transfer(
    deps: DepsMut,
    info: MessageInfo,
    recipient: String,
    amount: Uint128,
) -> Result<Response, ContractError> {
    // Validate address
    let recipient_addr = deps.api.addr_validate(&recipient)?;
    
    // Validate amount
    if amount.is_zero() {
        return Err(ContractError::InvalidAmount {});
    }
    
    // Validate sender has permission
    let config = CONFIG.load(deps.storage)?;
    if info.sender != config.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // ... proceed
}
```

### 2. Use Access Control

```rust
// Store the owner at instantiation
pub fn instantiate(deps: DepsMut, info: MessageInfo, ...) {
    let config = Config { owner: info.sender.clone() };
    CONFIG.save(deps.storage, &config)?;
}

// Check ownership before privileged operations
pub fn execute_admin_action(deps: DepsMut, info: MessageInfo, ...) {
    let config = CONFIG.load(deps.storage)?;
    if info.sender != config.owner {
        return Err(ContractError::Unauthorized {});
    }
}
```

### 3. Handle Funds Carefully

```rust
// Verify sent funds match expectations
pub fn execute_deposit(info: MessageInfo, ...) {
    let payment = info
        .funds
        .iter()
        .find(|c| c.denom == "ulala")
        .ok_or(ContractError::NoFundsProvided {})?;
    
    if payment.amount < minimum_deposit {
        return Err(ContractError::InsufficientDeposit {});
    }
}
```

### 4. Avoid Unbounded Iterations

```rust
// BAD: Iterating over all users
fn bad_distribute(deps: DepsMut) {
    let all_users = USERS.range(deps.storage, None, None, Order::Ascending);
    for user in all_users { ... } // Could run out of gas
}

// GOOD: Paginate or use batch limits
fn good_distribute(deps: DepsMut, limit: u32) {
    let users = USERS.range(deps.storage, None, None, Order::Ascending)
        .take(limit as usize);
    for user in users { ... }
}
```

### 5. Use Checked Math

```rust
use cosmwasm_std::Uint128;

// Uint128 automatically panics on overflow
let total = amount_a.checked_add(amount_b)?;
let share = total.checked_div(participants)?;
```

---

## Security Audit Checklist

- [ ] All entry points have proper access control
- [ ] No unbounded iterations
- [ ] All external inputs validated
- [ ] Fund handling is correct (check sent vs expected)
- [ ] State transitions are atomic
- [ ] No sensitive data in error messages
- [ ] Migration path secured (admin-only or no-admin)
- [ ] Contract cannot be bricked (recovery mechanism exists)

---

## Common Attack Vectors

| Attack | Prevention |
|--------|-----------|
| Unauthorized access | Check `info.sender` against stored permissions |
| Fund draining | Validate withdrawals against balances |
| DoS via gas | Limit iterations, paginate queries |
| Price manipulation | Use TWAP oracles, not spot prices |
| Flash loan attacks | Add time delays for sensitive operations |
