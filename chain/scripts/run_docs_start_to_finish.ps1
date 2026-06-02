param(
    [int]$Seed = 42,
    [int]$Phase0Epochs = 30,
    [int]$Phase1Epochs = 30,
    [int]$Phase2Epochs = 20,
    [int]$Validators = 8,
    [int]$Nodes = 8,
    [switch]$SkipFork
)

$ErrorActionPreference = "Stop"

$goExe = "C:\Program Files\Go\bin\go.exe"
if (-not (Test-Path $goExe)) {
    throw "Go executable not found at $goExe"
}

$lalachainRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$repoRoot = (Resolve-Path (Join-Path $lalachainRoot "..\..")).Path
$reportsDir = (Resolve-Path (Join-Path $lalachainRoot "reports")).Path

$stamp = Get-Date -Format "yyyyMMdd-HHmmss"
$phase0Log = Join-Path $reportsDir "phase0-simulation-$stamp.log"
$phase1Log = Join-Path $reportsDir "phase1-cosmos-single-$stamp.log"
$phase2Log = Join-Path $reportsDir "phase2-cosmos-multi-$stamp.log"
$phase3Log = Join-Path $reportsDir "phase3-shadow-$stamp.log"
$phase3Json = Join-Path $reportsDir "phase3-shadow-$stamp.json"
$forkLog = Join-Path $reportsDir "phase1-fork-bootstrap-$stamp.log"
$summaryMd = Join-Path $reportsDir "docs-start-to-finish-$stamp.md"

if (-not $SkipFork) {
    Write-Output "[Phase 1.a] Bootstrapping local Cosmos SDK fork workspace..."
    & "$PSScriptRoot\bootstrap_cosmos_fork.ps1" 2>&1 | Tee-Object -FilePath $forkLog
    if ($LASTEXITCODE -ne 0) {
        throw "Fork bootstrap failed. See $forkLog"
    }
}

Write-Output "[Phase 0] Running adaptation loop simulation..."
Set-Location $repoRoot
python -X utf8 simulation/simulation.py --epochs $Phase0Epochs --seed $Seed 2>&1 | Tee-Object -FilePath $phase0Log
if ($LASTEXITCODE -ne 0) {
    throw "Phase 0 simulation failed. See $phase0Log"
}

Write-Output "[Phase 1.b] Running Cosmos single-node phase 1 flow..."
Set-Location $lalachainRoot
& $goExe run ./cmd/lalachain --runtime cosmos --network single --phase 1 --epochs $Phase1Epochs --seed $Seed 2>&1 | Tee-Object -FilePath $phase1Log
if ($LASTEXITCODE -ne 0) {
    throw "Phase 1 single-node run failed. See $phase1Log"
}

Write-Output "[Phase 2] Running local multi-validator Cosmos testnet flow..."
& $goExe run ./cmd/lalachain --runtime cosmos --network multi --phase 2 --validators $Validators --nodes $Nodes --epochs $Phase2Epochs --seed $Seed 2>&1 | Tee-Object -FilePath $phase2Log
if ($LASTEXITCODE -ne 0) {
    throw "Phase 2 multi-validator run failed. See $phase2Log"
}

Write-Output "[Phase 3] Running research-track shadow-mode harness..."
Set-Location $lalachainRoot
python -X utf8 research/phase3_shadow_mode.py --input-log $phase2Log --output $phase3Json --seed $Seed --epochs $Phase2Epochs 2>&1 | Tee-Object -FilePath $phase3Log
if ($LASTEXITCODE -ne 0) {
    throw "Phase 3 shadow-mode run failed. See $phase3Log"
}

$forkStatus = if ($SkipFork) { "skipped" } else { "completed" }
$summary = @"
# LalaChain Docs Start-to-Finish Run

Generated: $(Get-Date -Format "u")

## Configuration
- Seed: $Seed
- Phase 0 epochs: $Phase0Epochs
- Phase 1 epochs: $Phase1Epochs
- Phase 2 epochs: $Phase2Epochs
- Validators: $Validators
- Nodes: $Nodes
- Fork bootstrap: $forkStatus

## Outputs
- Phase 0 log: $(Split-Path -Leaf $phase0Log)
- Phase 1 log: $(Split-Path -Leaf $phase1Log)
- Phase 2 log: $(Split-Path -Leaf $phase2Log)
- Phase 3 log: $(Split-Path -Leaf $phase3Log)
- Phase 3 report: $(Split-Path -Leaf $phase3Json)

## Notes
- Phase 3 remains research-dependent; this run executes a shadow-mode model comparison harness and governance process draft.
- For strict roadmap alignment, continue porting prototype modules into the local Cosmos SDK fork workspace created by bootstrap_cosmos_fork.ps1.
"@

Set-Content -Path $summaryMd -Value $summary -Encoding UTF8
Write-Output "Start-to-finish run complete. Summary: $summaryMd"
