$bytes = [System.IO.File]::ReadAllBytes("e:\testlala\config\genesis.json")
$first3 = $bytes[0..2] | ForEach-Object { '{0:X2}' -f $_ }
Write-Host "First 3 bytes: $($first3 -join ' ')"
if ($bytes[0] -eq 0xFF -or $bytes[0] -eq 0xFE -or ($bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB)) {
    Write-Host "WARNING: BOM detected!"
} else {
    Write-Host "No BOM - file looks ok"
}
