# Building Modules

**How to create custom Cosmos SDK modules for LalaChain.**

---

## Module Structure

A LalaChain module follows this layout:

```
x/mymodule/
├── keeper.go         # State management (reads/writes)
├── keeper_test.go    # Unit tests
├── types.go          # Message types, state types
├── handler.go        # Transaction message handlers
└── genesis.go        # Genesis state import/export
```

---

## Creating a Basic Module

### Step 1: Define Types

```go
// x/mymodule/types.go
package mymodule

type ModuleState struct {
    Counter int64  `json:"counter"`
    LastUpdated string `json:"last_updated"`
}

type MsgIncrement struct {
    Sender string `json:"sender"`
    Amount int64  `json:"amount"`
}
```

### Step 2: Implement Keeper

```go
// x/mymodule/keeper.go
package mymodule

import (
    "encoding/json"
    storetypes "cosmossdk.io/store/types"
)

type Keeper struct {
    storeKey storetypes.StoreKey
    state    ModuleState
}

func NewKeeper(storeKey storetypes.StoreKey) *Keeper {
    return &Keeper{
        storeKey: storeKey,
        state:    ModuleState{Counter: 0},
    }
}

func (k *Keeper) Increment(amount int64) {
    k.state.Counter += amount
}

func (k *Keeper) GetState() ModuleState {
    return k.state
}

func (k *Keeper) SaveState(store storetypes.KVStore) {
    data, _ := json.Marshal(k.state)
    store.Set([]byte("mymodule_state"), data)
}

func (k *Keeper) LoadState(store storetypes.KVStore) {
    data := store.Get([]byte("mymodule_state"))
    if data != nil {
        json.Unmarshal(data, &k.state)
    }
}
```

### Step 3: Register with App

```go
// In app/prototype.go, add to module registration:
myModuleKeeper := mymodule.NewKeeper(storeKey)

// Register REST endpoint
router.HandleFunc("/lala/mymodule/v1/state", func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(myModuleKeeper.GetState())
}).Methods("GET")
```

---

## Module Lifecycle Hooks

Modules can hook into the block lifecycle:

```go
// BeginBlocker - runs at start of each block
func (k *Keeper) BeginBlocker(ctx context.Context) {
    // Called before transaction processing
}

// EndBlocker - runs at end of each block
func (k *Keeper) EndBlocker(ctx context.Context) {
    // Called after transaction processing
    // LalaChain's epoch logic runs here
}
```

---

## Best Practices

1. **Deterministic state** — No random values, no external API calls, no floating-point
2. **Explicit persistence** — Call SaveState/LoadState at known points
3. **Bounded operations** — Don't iterate over unbounded collections in EndBlocker
4. **Test thoroughly** — Unit test all keeper methods
5. **Version state** — Include version in serialized state for upgradability

---

## Example: Complete Module Test

```go
// x/mymodule/keeper_test.go
package mymodule

import "testing"

func TestIncrement(t *testing.T) {
    k := NewKeeper(nil)
    
    k.Increment(5)
    if k.GetState().Counter != 5 {
        t.Errorf("expected 5, got %d", k.GetState().Counter)
    }
    
    k.Increment(3)
    if k.GetState().Counter != 8 {
        t.Errorf("expected 8, got %d", k.GetState().Counter)
    }
}
```

---

**Next:** [Network Endpoints](network-endpoints.md)
