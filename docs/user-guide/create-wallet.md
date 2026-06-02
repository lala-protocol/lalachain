# Create a Wallet

**A wallet is your gateway to LalaChain. This guide walks you through creating one.**

---

## Option 1: CLI Wallet (Recommended for Developers)

### Prerequisites
- `lalachaind` binary installed

### Steps

```bash
# Create a new wallet
lalachaind keys add my-wallet

# Output:
# - name: my-wallet
#   type: local
#   address: lala1qnk2n4nlkpw9xfqntladh74w6ux37lk3pzd7yh
#   pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"..."}'
#   mnemonic: ""
#
# **Important** write this mnemonic phrase in a safe place:
# word1 word2 word3 ... word24
```

### Critical: Save Your Seed Phrase

The 24-word mnemonic phrase is your **master backup**. Write it down on paper and store it securely.

- Never store it digitally (no screenshots, no text files, no email)
- Never share it with anyone
- If you lose it, you lose access to your tokens permanently

### Verify Your Wallet

```bash
# List all wallets
lalachaind keys list

# Show specific wallet address
lalachaind keys show my-wallet -a
```

---

## Option 2: Browser Extension (Keplr)

### Steps

1. Install [Keplr Wallet](https://www.keplr.app/) browser extension
2. Create a new account or import existing seed phrase
3. Add LalaChain network:
   - Chain ID: `lalachain-1`
   - RPC: Your node's RPC endpoint
   - REST: Your node's API endpoint
   - Currency: LALA (ulala, 6 decimals)
   - Address prefix: `lala`

4. Your wallet address will appear starting with `lala1...`

---

## Option 3: Recover Existing Wallet

If you have a seed phrase from a previous wallet:

```bash
lalachaind keys add my-wallet --recover
# Enter your 24-word mnemonic when prompted
```

---

## Wallet Security Checklist

- [ ] Seed phrase written on paper (not digital)
- [ ] Stored in a secure location (safe, safety deposit box)
- [ ] Never shared with anyone
- [ ] Tested recovery process (optional second device)
- [ ] Understand that transactions are irreversible

---

## Next Steps

Once you have a wallet, you'll need tokens to use the network:

**Next:** [Get LALA Tokens](get-lala-tokens.md)
