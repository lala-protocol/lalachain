---
title: "Validator Maintenance"
description: "Upgrades, backups, and routine validator maintenance."
---

# Validator Maintenance

**Ongoing operational tasks to keep your validator healthy and performant.**

---

## Routine Tasks

### Daily
- Check monitoring dashboard for anomalies
- Verify block signing (no missed blocks)
- Review disk space trends

### Weekly
- Review and apply OS security updates
- Check peer connectivity quality
- Review governance proposals (vote if needed)
- Verify backup integrity

### Monthly
- Review validator performance metrics
- Audit access logs
- Test disaster recovery procedure
- Update documentation

---

## Software Upgrades

### Standard Upgrade (Governance-Coordinated)

```bash
# 1. Watch for upgrade proposal on-chain
# 2. Download new binary BEFORE upgrade height
wget https://github.com/lalachain/lalachain/releases/download/v2.0.0/lalachaind-v2.0.0-linux-amd64

# 3. Prepare binary
chmod +x lalachaind-v2.0.0-linux-amd64
sudo cp lalachaind-v2.0.0-linux-amd64 /usr/local/bin/lalachaind-new

# 4. At upgrade height, node will halt automatically
# 5. Swap binary
sudo systemctl stop lalachaind
sudo cp /usr/local/bin/lalachaind-new /usr/local/bin/lalachaind
sudo systemctl start lalachaind
```

### Using Cosmovisor (Automatic Upgrades)

```bash
# Install cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

# Setup directory structure
mkdir -p ~/.lalachaind/cosmovisor/genesis/bin
mkdir -p ~/.lalachaind/cosmovisor/upgrades/v2/bin

cp /usr/local/bin/lalachaind ~/.lalachaind/cosmovisor/genesis/bin/
cp lalachaind-v2.0.0 ~/.lalachaind/cosmovisor/upgrades/v2/bin/lalachaind
```

---

## State Pruning

Reduce disk usage by pruning old state:

```toml
# app.toml
[pruning]
pruning = "custom"
pruning-keep-recent = "100"
pruning-interval = "10"
```

For archive nodes, use `pruning = "nothing"`.

---

## Unjailing

If your validator gets jailed for downtime:

```bash
# Wait until jail time expires (10 minutes for downtime)
# Then unjail:
lalachaind tx slashing unjail --from validator --fees 10000ulala -y

# Verify status
lalachaind query staking validator $(lalachaind keys show validator --bech val -a) | grep jailed
```

---

## Emergency Procedures

### Node Won't Start
1. Check logs: `journalctl -u lalachaind -n 100`
2. Verify disk space: `df -h`
3. Try resetting state (DANGER — only for non-validator or with backup):
   ```bash
   lalachaind tendermint unsafe-reset-all
   ```

### Validator Double-Signing (Critical)
1. **IMMEDIATELY STOP THE NODE**
2. Investigate cause (two instances? failed migration?)
3. If key is compromised: rotate immediately
4. Contact community for guidance
5. You will be tombstoned (permanent jail)
