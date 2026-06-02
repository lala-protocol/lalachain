$ErrorActionPreference = "Continue"
$BIN = "e:\side-gigs\Lala\whitepaper\prototype\lalachain\build\lalachaind.exe"
$HOME_DIR = "e:\testlala"
$CHAIN_ID = "lalachain-local-1"

# Init
& $BIN init "lala-node" --chain-id $CHAIN_ID --home $HOME_DIR 2>$null

# Fix denom
$genesis = "$HOME_DIR\config\genesis.json"
(Get-Content $genesis -Raw) -replace '"stake"', '"ulala"' | Set-Content $genesis

# Create key + genesis
& $BIN keys add validator --keyring-backend test --home $HOME_DIR 2>$null
$ADDR = & $BIN keys show validator -a --keyring-backend test --home $HOME_DIR
& $BIN add-genesis-account $ADDR "100000000ulala" --home $HOME_DIR --keyring-backend test 2>$null
& $BIN gentx validator "50000000ulala" --chain-id $CHAIN_ID --home $HOME_DIR --keyring-backend test 2>$null
& $BIN collect-gentxs --home $HOME_DIR 2>$null
& $BIN validate-genesis --home $HOME_DIR 2>$null

Write-Host "SETUP DONE - Address: $ADDR"
