---
title: "troubleshooting"
description: "Solutions to common issues when using LalaChain."
---

# Troubleshooting

**Solutions to common issues when using LalaChain.**

---

## Wallet Issues

### Can't create a wallet
**Symptoms:** `lalachaind keys add` fails or returns an error.

**Solutions:**
1. Ensure `lalachaind` is in your PATH or use the full path to the binary
2. Check that the home directory exists: `lalachaind init <moniker> --home ~/.lalachain`
3. If keyring backend error: specify `--keyring-backend test` (for testing) or `--keyring-backend os`

---

### Lost seed phrase
**Reality:** If you've lost your seed phrase and have no backup, your funds are permanently inaccessible. There is no recovery mechanism.

**Prevention:** Always store seed phrase backups in multiple secure physical locations.

---

### Wrong address format
**Symptoms:** Address doesn't start with `lala1`.

**Solutions:**
- Ensure you're using the correct chain binary (`lalachaind`, not `gaiad` or others)
- Verify account prefix in genesis: should be `lala`

---

## Transaction Issues

### Transaction rejected: insufficient fees
**Symptoms:** `insufficient fees; got: X, required: Y`

**Solutions:**
1. Query current minimum fee: `GET /lala/telemetry/v1/kpis` → check `BaseFeePerGas`
2. Increase fee: `--fees <higher_amount>ulala`
3. Use `--gas auto --gas-adjustment 1.3` for automatic estimation

---

### Transaction rejected: out of gas
**Symptoms:** `out of gas in location: ...`

**Solutions:**
1. Increase gas limit: `--gas 300000` (or higher)
2. Use `--gas auto --gas-adjustment 1.5` for complex transactions
3. Simplify transaction if hitting block gas limit

---

### Transaction stuck / not included
**Symptoms:** Transaction submitted but never appears in a block.

**Solutions:**
1. Check fee is above minimum (base fee may have increased since submission)
2. Verify node is synced: `lalachaind status` → `catching_up` should be `false`
3. Check mempool: transaction may be waiting behind others
4. Resubmit with higher fee and same sequence number

---

### Account sequence mismatch
**Symptoms:** `account sequence mismatch, expected X, got Y`

**Solutions:**
1. Query your current sequence: `lalachaind query account <address>`
2. Use `--sequence <correct_number>` flag
3. Wait for pending transactions to be included before submitting new ones

---

## Node Issues

### Node won't start
**Symptoms:** `lalachaind start` immediately exits with error.

**Common causes & fixes:**

| Error | Fix |
|-------|-----|
| `genesis file not found` | Run `lalachaind init <moniker>` first |
| `address already in use` | Another instance is running, or change ports in `config.toml` |
| `wrong Block.Header.AppHash` | State corruption — reset: `lalachaind tendermint unsafe-reset-all` |
| `permission denied` | Check file permissions on home directory |

---

### Node not syncing
**Symptoms:** `catching_up: true` indefinitely, block height not increasing.

**Solutions:**
1. Check peers: `curl localhost:26657/net_info` → should have connected peers
2. Add persistent peers to `config.toml`:
   ```toml
   [p2p]
   persistent_peers = "node-id@ip:26656"
   ```
3. Check firewall allows port 26656 inbound/outbound
4. Verify genesis file matches the network

---

### Node crashes with out of memory
**Symptoms:** Process killed, OOM errors in system logs.

**Solutions:**
1. Increase available RAM (4GB minimum recommended)
2. Enable state sync to avoid full history
3. Prune aggressively:
   ```toml
   [pruning]
   pruning = "everything"  # Keep minimal state
   ```

---

## Validator Issues

### Validator jailed
**Symptoms:** Validator stops earning rewards, shows `jailed: true`.

**Causes:**
- Missed too many blocks (downtime)
- Double-signed a block

**Solutions:**
1. Fix the underlying issue (connectivity, disk space, etc.)
2. Unjail: `lalachaind tx slashing unjail --from <key>`
3. Wait for jail duration to expire (if applicable)

---

### Validator not producing blocks
**Symptoms:** Validator is active but never selected as proposer.

**Solutions:**
1. Verify validator is in the active set: `lalachaind query staking validators`
2. Check your voting power (low power = rare proposals, but should still happen)
3. Ensure clock is synchronized (NTP)
4. Check CometBFT logs for consensus errors

---

### Double-sign slashing
**Symptoms:** Large token slash, validator tombstoned.

**Cause:** Same validator key signed two different blocks at the same height (usually from running two instances).

**Prevention:**
- NEVER run the same validator key on two machines
- Use KMS to prevent duplicate signing
- Careful during migrations — shut down old node completely first

**Recovery:** Tombstoned validators cannot be unjailed. You must create a new validator with new keys.

---

## API / Query Issues

### API returns 404
**Symptoms:** REST endpoint returns "Not Implemented" or 404.

**Solutions:**
1. Ensure API is enabled in `app.toml`:
   ```toml
   [api]
   enable = true
   address = "tcp://0.0.0.0:1317"
   ```
2. Restart the node after config changes
3. Verify the endpoint path matches documentation

---

### gRPC connection refused
**Symptoms:** Can't connect to gRPC on port 9090.

**Solutions:**
1. Enable in `app.toml`:
   ```toml
   [grpc]
   enable = true
   address = "0.0.0.0:9090"
   ```
2. Check firewall allows the port
3. Restart node after changes

---

### Stale data / old block height
**Symptoms:** Queries return outdated information.

**Solutions:**
1. Check node sync status: `lalachaind status`
2. If `catching_up: true`, wait for sync to complete
3. For indexing issues, rebuild index: restart with `--x-crisis-skip-assert-invariants`

---

## Governance Issues

### Vote not counted
**Symptoms:** Submitted vote but not reflected in proposal tally.

**Solutions:**
1. Verify you voted during the voting period (not before or after)
2. Check you had staked tokens at the time of the vote
3. Confirm transaction was included in a block (check tx hash)

---

### Proposal passed but nothing changed
**Symptoms:** Proposal shows "PASSED" but parameter unchanged.

**Explanation:** LalaChain has an activation delay of 2 epochs. The change will take effect after the delay period.

---

## Performance Issues

### High memory usage
- Enable pruning in `app.toml`
- Reduce `cache_size` in config
- Use state sync for new nodes

### Slow block processing
- Check disk I/O (SSD required for validators)
- Ensure adequate CPU
- Monitor for network congestion

### Large database size
- Enable pruning
- Use `lalachaind tendermint unsafe-reset-all` and state sync to rebuild
