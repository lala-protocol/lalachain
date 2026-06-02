# LalaChain Incentivized Testnet Program

## Overview

The LalaChain Incentivized Testnet rewards early validators and community members who contribute to network stability, governance participation, and bug discovery.

## Rewards Structure

| Category | Allocation | Criteria |
|----------|-----------|----------|
| Validator Uptime | 40% | >95% uptime across testnet period |
| Governance Participation | 25% | Vote on all AI proposals |
| Bug Reports | 20% | Valid bugs reported via GitHub |
| Community Contributions | 15% | Documentation, tooling, guides |

## Validator Rewards

### Uptime Scoring

- **100% uptime:** Full tier allocation
- **95–99%:** 75% of tier
- **90–94%:** 50% of tier
- **<90%:** No reward

### Requirements

1. Run a validator node for the entire testnet period
2. Maintain >95% uptime (missed <5 blocks per 100)
3. Do not get jailed for double-signing
4. Participate in at least 80% of governance votes

## Governance Rewards

Validators earn additional rewards for active governance participation:

- Vote on every AI-originated proposal
- Provide rationale for controversial votes (bonus)
- Participate in parameter discussions

## Bug Bounty

| Severity | Reward |
|----------|--------|
| Critical (consensus break, fund loss) | 5,000 LALA |
| High (state inconsistency, DoS vector) | 2,000 LALA |
| Medium (performance issue, edge case) | 1,000 LALA |
| Low (UI bug, documentation error) | 250 LALA |

### Submission Process

1. Open a GitHub issue with the `bug-bounty` label
2. Include: reproduction steps, expected vs actual behavior, severity assessment
3. Do NOT publicly disclose critical vulnerabilities before fix
4. Team responds within 48 hours

## Timeline

| Phase | Duration | Focus |
|-------|----------|-------|
| Phase 1: Genesis | Week 1–2 | Validator onboarding, network stability |
| Phase 2: Governance | Week 3–4 | AI proposals, voting mechanics |
| Phase 3: Stress Test | Week 5–6 | High tx load, parameter extremes |
| Phase 4: Graduation | Week 7–8 | Final scoring, mainnet preparation |

## How to Participate

1. **Join Discord:** [discord.gg/lalachain](https://discord.gg/lalachain)
2. **Set up validator:** Follow the [Validator Guide](./VALIDATOR_GUIDE.md)
3. **Register:** Submit your validator address in the `#testnet-validators` channel
4. **Monitor:** Use the dashboard to track your performance
5. **Govern:** Vote on all proposals via CLI or dashboard

## Disqualification

Validators may be disqualified for:

- Double-signing (intentional or due to misconfiguration)
- Running multiple validators from the same entity (Sybil)
- Attempting to exploit the network maliciously
- Failing to maintain minimum uptime (<90%)
- Not voting on any proposals during a full phase

## FAQ

**Q: When does the testnet start?**
A: Genesis will be announced in Discord with 48-hour notice.

**Q: Can I run multiple validators?**
A: One validator per team/individual. Multiple nodes for redundancy are fine.

**Q: What happens if I get jailed?**
A: Unjail within 1 hour for downtime. Double-sign jailing is permanent for that testnet.

**Q: Will testnet tokens convert to mainnet?**
A: No. Rewards are allocated separately at mainnet genesis based on testnet performance.
