# Smart Contract Setup

**Install the tools needed to write, compile, and deploy CosmWasm smart contracts on LalaChain.**

---

## Prerequisites

### 1. Rust Toolchain

```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Add Wasm target
rustup target add wasm32-unknown-unknown

# Verify
rustc --version
cargo --version
```

### 2. cargo-generate (Project Templates)

```bash
cargo install cargo-generate
```

### 3. wasm-opt (Optimization)

```bash
# Install binaryen (includes wasm-opt)
# macOS:
brew install binaryen

# Ubuntu:
apt install binaryen

# Windows:
# Download from https://github.com/WebAssembly/binaryen/releases
```

### 4. LalaChain Node

Ensure you have a running LalaChain node with the wasm module enabled.

---

## Project Setup

### Generate from Template

```bash
cargo generate --git https://github.com/CosmWasm/cw-template.git --name my-contract
cd my-contract
```

### Project Structure

```
my-contract/
├── Cargo.toml          # Dependencies
├── src/
│   ├── contract.rs     # Main contract logic
│   ├── error.rs        # Custom errors
│   ├── msg.rs          # Message definitions
│   ├── state.rs        # State storage
│   └── lib.rs          # Module exports
└── tests/
    └── integration.rs  # Integration tests
```

---

## Compile Contract

```bash
# Debug build
cargo build --target wasm32-unknown-unknown

# Optimized production build
cargo build --release --target wasm32-unknown-unknown

# Optimize with wasm-opt (reduces size significantly)
wasm-opt -Os target/wasm32-unknown-unknown/release/my_contract.wasm \
  -o artifacts/my_contract.wasm
```

### Using Docker Optimizer (Recommended)

```bash
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.15.0
```

This produces a reproducible, optimized Wasm binary.

---

## Deploy to LalaChain

### Upload Code

```bash
lalachaind tx wasm store artifacts/my_contract.wasm \
  --from my-wallet \
  --gas auto \
  --fees 50000ulala \
  --chain-id lalachain-local -y
```

Note the **code_id** from the response.

### Instantiate Contract

```bash
lalachaind tx wasm instantiate <code_id> '{"count": 0}' \
  --from my-wallet \
  --label "my-first-contract" \
  --admin $(lalachaind keys show my-wallet -a) \
  --gas auto \
  --fees 50000ulala \
  --chain-id lalachain-local -y
```

Note the **contract address** from the response.

---

## Verify Deployment

```bash
# Query contract info
lalachaind query wasm contract <contract-address>

# Query contract state
lalachaind query wasm contract-state smart <contract-address> '{"get_count": {}}'
```

---

**Next:** [First Contract](first-contract.md)
