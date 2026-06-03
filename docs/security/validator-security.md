---
title: "Validator Security"
description: "Operational security for LalaChain validator operators."
---

# Validator Security

**Security practices specific to LalaChain validator operators.**

---

## Threat Model

| Threat | Impact | Likelihood |
|--------|--------|-----------|
| Server compromise | Key theft, double-signing | Medium |
| DDoS attack | Downtime, missed blocks | High |
| Insider attack | Key exposure | Low |
| Supply chain attack | Malicious binary | Low |
| Social engineering | Key disclosure | Medium |

---

## Sentry Node Architecture

```
[Internet/Peers] → [Sentry 1] → [Validator (hidden)]
                 → [Sentry 2] ↗
```

- Validator has NO public IP
- Sentries handle peer connections
- DDoS targets sentries, not validator
- Validator connects only to trusted sentries

### Configuration

Validator `config.toml`:
```toml
[p2p]
pex = false                    # Disable peer exchange
persistent_peers = "sentry1-id@sentry1-ip:26656,sentry2-id@sentry2-ip:26656"
addr_book_strict = false
```

Sentry `config.toml`:
```toml
[p2p]
pex = true
private_peer_ids = "validator-node-id"  # Hide validator from peer list
persistent_peers = "validator-id@validator-ip:26656"
```

---

## Key Security

### Use KMS (Key Management Service)

Never store validator keys on the same machine that faces the internet:

1. Install `tmkms` on a separate secure machine
2. Connect validator to KMS via private network
3. Keys never leave the KMS machine

### Hardware Security Module (HSM)

For maximum security:
- YubiHSM2 holds signing keys
- Key operations happen inside the HSM
- Key can never be extracted, even if server is compromised

---

## Access Control

| Practice | Implementation |
|----------|---------------|
| SSH key-only access | Disable password auth |
| Minimal users | Only operator accounts needed |
| Firewall | Allow only required ports |
| Fail2ban | Block brute-force attempts |
| 2FA | For all administrative access |
| Audit logs | Track all commands executed |

---

## Monitoring for Security Events

Watch for:
- Unexpected SSH logins
- Process changes
- Unusual network traffic
- Validator signing anomalies
- Disk/CPU spikes (crypto mining)

---

## Incident Response Plan

1. **Detection** — Monitoring alert fires
2. **Containment** — Stop validator if key may be compromised
3. **Assessment** — Determine what happened
4. **Recovery** — Rotate keys, restore from backup
5. **Communication** — Inform delegators and community
6. **Post-mortem** — Document and prevent recurrence
