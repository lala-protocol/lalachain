$ErrorActionPreference = "SilentlyContinue"
$HOME_DIR = "e:\testlala"
$BINARY = "e:\side-gigs\Lala\whitepaper\prototype\lalachain\build\lalachaind.exe"
$CHAIN_ID = "lalachain-local-1"

# Fix genesis denom
$genesisPath = "$HOME_DIR\config\genesis.json"
$content = Get-Content $genesisPath -Raw
$content = $content -replace '"bond_denom": "stake"', '"bond_denom": "ulala"'
$content = $content -replace '"mint_denom": "stake"', '"mint_denom": "ulala"'
Set-Content -Path $genesisPath -Value $content -NoNewline
Write-Host "[OK] Genesis denom set to ulala"

# Create key
$null = & $BINARY keys add myvalidator --keyring-backend test --home $HOME_DIR 2>&1
$addr = (& $BINARY keys show myvalidator -a --keyring-backend test --home $HOME_DIR 2>&1) | Select-Object -First 1
Write-Host "[OK] Validator address: $addr"

# Add genesis account
& $BINARY add-genesis-account $addr 100000000000000ulala --keyring-backend test --home $HOME_DIR
Write-Host "[OK] Genesis account funded with 100M LALA"

# Create gentx
$null = & $BINARY gentx myvalidator 50000000000000ulala --chain-id $CHAIN_ID --keyring-backend test --moniker testnode --home $HOME_DIR 2>&1
Write-Host "[OK] Gentx created (50M LALA self-delegation)"

# Collect gentxs
$null = & $BINARY collect-gentxs --home $HOME_DIR 2>&1
Write-Host "[OK] Gentxs collected"

# Validate
$null = & $BINARY validate --home $HOME_DIR 2>&1
Write-Host "[OK] Genesis validated"
Write-Host ""
Write-Host "=== Ready to start! ==="
Write-Host "Run: $BINARY start --minimum-gas-prices=0.025ulala --api.enable=true --home $HOME_DIR"
