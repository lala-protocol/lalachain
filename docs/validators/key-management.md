# Key Management

**Securing your validator keys is the most critical operational task. A compromised key means lost funds and potential slashing.**

---

## Key Types

| Key | Purpose | Risk if Compromised |
|-----|---------|-------------------|
| **Validator consensus key** | Signs blocks and consensus votes | Double-signing → 5% slash + jail |
| **Operator key** | Manages validator (unjail, edit) | Validator takeover |
| **Delegator key** | Holds and delegates tokens | Fund theft |

---

## Key Storage Options

### Option 1: File-Based (Default)

Keys stored in `~/.lalachaind/config/priv_validator_key.json`

- **Pros:** Simple, no extra infrastructure
- **Cons:** If server is compromised, key is exposed
- **Suitable for:** Testnet, small validators

### Option 2: Key Management Service (KMS)

Use Tendermint KMS (tmkms) to sign remotely:

```toml
# tmkms.toml
[[chain]]
id = "lalachain-1"
key_format = { type = "bech32", account_key_prefix = "lalapub" }

[[validator]]
chain_id = "lalachain-1"
addr = "tcp://validator-ip:26658"
secret_key = "path/to/secret_connection.key"
```

- **Pros:** Key never on validator server, HSM support
- **Cons:** Additional infrastructure, complexity
- **Suitable for:** Production validators

### Option 3: Hardware Security Module (HSM)

YubiHSM2 or similar hardware:

- **Pros:** Key never leaves hardware device
- **Cons:** Expensive, single point of failure without redundancy
- **Suitable for:** High-value validators

---

## Backup Procedures

### Validator Key Backup

```bash
# Backup consensus key (KEEP SECURE!)
cp ~/.lalachaind/config/priv_validator_key.json /secure/backup/location/

# Backup node key
cp ~/.lalachaind/config/node_key.json /secure/backup/location/
```

### Operator Key Backup

```bash
# Export encrypted key
lalachaind keys export validator > /secure/backup/validator-key.armor

# To restore on new machine:
lalachaind keys import validator /path/to/validator-key.armor
```

---

## Security Practices

1. **Never share validator key** — Even with team members
2. **Rotate access credentials** — SSH keys, server passwords regularly
3. **Use firewall** — Only allow necessary ports
4. **Encrypt backups** — Use GPG or similar for key backups
5. **Test recovery** — Periodically verify you can restore from backup
6. **Sentry node architecture** — Hide validator IP behind sentry nodes

---

## Double-Sign Prevention

The most severe slashing offense. Prevent by:

- **Never run two validators with the same key** simultaneously
- Use `priv_validator_state.json` to track signing state
- If migrating: stop old node, verify it's dead, then start new node
- Consider tmkms for hardware-enforced single-signing

---

**Next:** [Monitoring](monitoring.md)
