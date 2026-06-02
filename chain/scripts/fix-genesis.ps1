$genesisPath = "$env:USERPROFILE\.lalachaind\config\genesis.json"
$content = Get-Content $genesisPath -Raw
$content = $content -replace '"bond_denom": "stake"', '"bond_denom": "ulala"'
$content = $content -replace '"mint_denom": "stake"', '"mint_denom": "ulala"'
Set-Content -Path $genesisPath -Value $content -NoNewline
Write-Host "Genesis updated: bond_denom and mint_denom set to ulala"
