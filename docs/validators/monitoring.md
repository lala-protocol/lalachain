---
title: "Validator Monitoring"
description: "Monitor your validator's performance, uptime, and health."
---

# Validator Monitoring

**Set up monitoring to ensure high uptime and catch issues before they cause slashing.**

---

## Prometheus Metrics

LalaChain exposes Prometheus metrics on port 26660:

```toml
# config/config.toml
[instrumentation]
prometheus = true
prometheus_listen_addr = ":26660"
```

### Key Metrics to Monitor

| Metric | Alert Threshold | Meaning |
|--------|----------------|---------|
| `cometbft_consensus_height` | Stalling | Node stopped producing |
| `cometbft_consensus_validators` | Changes | Validator set modified |
| `cometbft_consensus_missing_validators` | > 0 (your node) | You're missing votes |
| `cometbft_p2p_peers` | < 3 | Losing connectivity |
| `cometbft_mempool_size` | > 1000 | Mempool congestion |
| `process_resident_memory_bytes` | > 80% of RAM | Memory pressure |

---

## Grafana Dashboard

Recommended panels:
1. **Block height** (real-time counter)
2. **Consensus rounds** (should be mostly 0-1)
3. **Peer count** (stable, 10-50)
4. **Mempool size** (should drain regularly)
5. **Block time** (target ~5s)
6. **Disk usage** (trending)

---

## Alert Rules

```yaml
# alertmanager rules
groups:
  - name: validator
    rules:
      - alert: NodeDown
        expr: up{job="lalachaind"} == 0
        for: 1m
        annotations:
          summary: "Validator node is down"
          
      - alert: MissedBlocks
        expr: increase(cometbft_consensus_missing_validators[5m]) > 0
        annotations:
          summary: "Validator missing consensus votes"
          
      - alert: LowPeers
        expr: cometbft_p2p_peers < 3
        for: 5m
        annotations:
          summary: "Low peer count - connectivity issue"
          
      - alert: DiskSpace
        expr: node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes < 0.1
        annotations:
          summary: "Disk space below 10%"
```

---

## Health Check Script

```bash
#!/bin/bash
# health-check.sh

# Check if node is synced
CATCHING_UP=$(curl -s localhost:26657/status | jq -r '.result.sync_info.catching_up')
if [ "$CATCHING_UP" = "true" ]; then
    echo "WARNING: Node is catching up"
    exit 1
fi

# Check latest block time
LATEST_TIME=$(curl -s localhost:26657/status | jq -r '.result.sync_info.latest_block_time')
BLOCK_AGE=$(($(date +%s) - $(date -d "$LATEST_TIME" +%s)))
if [ $BLOCK_AGE -gt 30 ]; then
    echo "WARNING: Last block was ${BLOCK_AGE}s ago"
    exit 1
fi

echo "OK: Node healthy"
exit 0
```

---

## Log Monitoring

Watch for these log patterns:

| Pattern | Meaning | Action |
|---------|---------|--------|
| `committed state` | Normal block commit | None |
| `validator not in set` | You're not active | Check stake/jail status |
| `connection refused` | Peer connectivity issue | Check firewall |
| `out of memory` | OOM kill imminent | Increase RAM |
| `WAL replay` | Node crashed, recovering | Monitor recovery |
