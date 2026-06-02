param(
    [string]$ForkDir = "e:\side-gigs\Lala\cosmos-sdk-fork",
    [string]$Version = "v0.50.9"
)

$ErrorActionPreference = "Stop"

function Invoke-Git {
    param(
        [Parameter(Mandatory = $true)]
        [string[]]$Arguments,
        [Parameter(Mandatory = $true)]
        [string]$FailureMessage
    )

    $escapedArgs = ($Arguments | ForEach-Object {
        if ($_ -match '[\s"]') {
            '"' + ($_ -replace '"', '\\"') + '"'
        }
        else {
            $_
        }
    }) -join " "

    cmd /c "git $escapedArgs 2>&1"
    if ($LASTEXITCODE -ne 0) {
        throw "$FailureMessage (exit code: $LASTEXITCODE)"
    }
}

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    throw "git is required but was not found in PATH"
}

$forkExists = Test-Path (Join-Path $ForkDir ".git")
if (-not $forkExists) {
    $parent = Split-Path -Parent $ForkDir
    if (-not (Test-Path $parent)) {
        New-Item -ItemType Directory -Path $parent -Force | Out-Null
    }

    Write-Output "Cloning Cosmos SDK into $ForkDir (version $Version)..."
    Invoke-Git -Arguments @("clone", "--depth", "1", "--branch", $Version, "https://github.com/cosmos/cosmos-sdk.git", $ForkDir) -FailureMessage "git clone failed"
}
else {
    Write-Output "Existing fork found at $ForkDir; refreshing checkout to $Version..."
    Push-Location $ForkDir
    Invoke-Git -Arguments @("fetch", "origin", "tag", $Version, "--force") -FailureMessage "git fetch failed"
    Invoke-Git -Arguments @("checkout", $Version) -FailureMessage "git checkout failed"
    Pop-Location
}

Push-Location $ForkDir

$branchName = "lalachain-phase1"
$branchExists = (& git branch --list $branchName)
if ([string]::IsNullOrWhiteSpace($branchExists)) {
    Invoke-Git -Arguments @("checkout", "-b", $branchName) -FailureMessage "failed to create branch $branchName"
}
else {
    Invoke-Git -Arguments @("checkout", $branchName) -FailureMessage "failed to checkout branch $branchName"
}

$notes = @"
# LalaChain Fork Notes

This local fork workspace is prepared for roadmap alignment with
IMPLEMENTATION_FEASIBILITY.md:

- Base version: $Version
- Working branch: $branchName

Next integration tasks:
1. Port x/telemetry, x/aiadvisor, x/gov logic from prototype/lalachain.
2. Wire lalachain app module manager and parameter store handlers.
3. Run local multi-validator testnet and compare against prototype behavior.
"@

Set-Content -Path (Join-Path $ForkDir "LALACHAIN_FORK_NOTES.md") -Value $notes -Encoding UTF8

Pop-Location

Write-Output "Cosmos SDK fork workspace is ready: $ForkDir"
