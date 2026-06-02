$ErrorActionPreference = "Continue"
$BIN = "e:\side-gigs\Lala\whitepaper\prototype\lalachain\build\lalachaind.exe"
$HOME_DIR = "e:\testlala"
$CHAIN_ID = "lalachain-local-1"
$DENOM = "ulala"
$KEY_NAME = "validator"

Write-Host "=== Step 1: Init chain ==="
& $BIN init "lala-node" --chain-id $CHAIN_ID --home $HOME_DIR 2>&1 | Out-Null
Write-Host "Init done"

Write-Host "=== Step 2: Fix genesis denom ==="
$genesis = "$HOME_DIR\config\genesis.json"
$content = Get-Content $genesis -Raw
$content = $content -replace '"stake"', "`"$DENOM`""
Set-Content -Path $genesis -Value $content
Write-Host "Denom fixed to $DENOM"

Write-Host "=== Step 3: Create key ==="
& $BIN keys add $KEY_NAME --keyring-backend test --home $HOME_DIR 2>&1 | Out-Null
Write-Host "Key created"

Write-Host "=== Step 4: Get address ==="
$ADDR = & $BIN keys show $KEY_NAME -a --keyring-backend test --home $HOME_DIR
Write-Host "Address: $ADDR"

Write-Host "=== Step 5: Add genesis account ==="
& $BIN add-genesis-account $ADDR "100000000${DENOM}" --home $HOME_DIR --keyring-backend test
Write-Host "Genesis account added"

Write-Host "=== Step 6: Create gentx ==="
& $BIN gentx $KEY_NAME "50000000${DENOM}" --chain-id $CHAIN_ID --home $HOME_DIR --keyring-backend test 2>&1 | Out-Null
Write-Host "Gentx created"

Write-Host "=== Step 7: Collect gentxs ==="
& $BIN collect-gentxs --home $HOME_DIR 2>&1 | Out-Null
Write-Host "Gentxs collected"

Write-Host "=== Step 8: Validate genesis ==="
& $BIN validate-genesis --home $HOME_DIR
Write-Host "Genesis valid"

Write-Host ""
Write-Host "=== SETUP COMPLETE ==="
Write-Host "Validator address: $ADDR"
Write-Host "To start: $BIN start --home $HOME_DIR --minimum-gas-prices=0.025ulala --api.enable=true --api.address=tcp://0.0.0.0:1317"
