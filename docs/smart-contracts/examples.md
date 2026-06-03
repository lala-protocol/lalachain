---
title: "Smart Contract Examples"
description: "Example smart contracts with full source code."
---

# Smart Contract Examples

**Common patterns and templates for building on LalaChain.**

---

## Example 1: Token Vesting

A contract that releases tokens over time:

```rust
#[cw_serde]
pub struct VestingSchedule {
    pub beneficiary: Addr,
    pub total_amount: Uint128,
    pub start_time: u64,
    pub end_time: u64,
    pub claimed: Uint128,
}

pub fn execute_claim(deps: DepsMut, env: Env, info: MessageInfo) -> StdResult<Response> {
    let mut schedule = VESTING.load(deps.storage)?;
    
    if info.sender != schedule.beneficiary {
        return Err(StdError::generic_err("unauthorized"));
    }
    
    let now = env.block.time.seconds();
    let elapsed = now.saturating_sub(schedule.start_time);
    let duration = schedule.end_time - schedule.start_time;
    
    let vested = if now >= schedule.end_time {
        schedule.total_amount
    } else {
        schedule.total_amount.multiply_ratio(elapsed, duration)
    };
    
    let claimable = vested - schedule.claimed;
    schedule.claimed = vested;
    VESTING.save(deps.storage, &schedule)?;
    
    let msg = BankMsg::Send {
        to_address: info.sender.to_string(),
        amount: vec![Coin::new(claimable.u128(), "ulala")],
    };
    
    Ok(Response::new().add_message(msg))
}
```

---

## Example 2: Simple Escrow

Hold funds until conditions are met:

```rust
pub fn execute_create_escrow(
    deps: DepsMut,
    info: MessageInfo,
    recipient: String,
    arbiter: String,
) -> StdResult<Response> {
    let escrow = Escrow {
        sender: info.sender.clone(),
        recipient: deps.api.addr_validate(&recipient)?,
        arbiter: deps.api.addr_validate(&arbiter)?,
        amount: info.funds.clone(),
        is_released: false,
    };
    ESCROW.save(deps.storage, &escrow)?;
    Ok(Response::new())
}

pub fn execute_release(deps: DepsMut, info: MessageInfo) -> StdResult<Response> {
    let escrow = ESCROW.load(deps.storage)?;
    
    if info.sender != escrow.arbiter {
        return Err(StdError::generic_err("only arbiter can release"));
    }
    
    let msg = BankMsg::Send {
        to_address: escrow.recipient.to_string(),
        amount: escrow.amount,
    };
    
    Ok(Response::new().add_message(msg))
}
```

---

## Example 3: On-Chain Voting

Leverage LalaChain's governance patterns in a contract:

```rust
#[cw_serde]
pub struct Poll {
    pub question: String,
    pub options: Vec<String>,
    pub votes: Vec<Uint128>,  // votes per option
    pub voters: Vec<Addr>,    // who already voted
    pub end_time: u64,
}

pub fn execute_vote(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    option: u32,
) -> StdResult<Response> {
    let mut poll = POLL.load(deps.storage)?;
    
    if env.block.time.seconds() > poll.end_time {
        return Err(StdError::generic_err("voting ended"));
    }
    
    if poll.voters.contains(&info.sender) {
        return Err(StdError::generic_err("already voted"));
    }
    
    if option as usize >= poll.options.len() {
        return Err(StdError::generic_err("invalid option"));
    }
    
    poll.votes[option as usize] += Uint128::one();
    poll.voters.push(info.sender);
    POLL.save(deps.storage, &poll)?;
    
    Ok(Response::new().add_attribute("action", "vote"))
}
```

---

## Patterns Reference

| Pattern | Use Case | Key Concept |
|---------|----------|-------------|
| **Escrow** | Conditional payments | Hold funds, release on condition |
| **Vesting** | Token distribution | Time-locked release |
| **Voting** | Governance | Track voters, prevent doubles |
| **NFT (CW721)** | Digital ownership | Unique tokens, metadata |
| **Fungible (CW20)** | Custom tokens | ERC20-equivalent |
| **Staking** | Yield farming | Deposit, track shares, withdraw |

