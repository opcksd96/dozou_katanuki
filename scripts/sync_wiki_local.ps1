<#
.SYNOPSIS
    docs/wiki の内容を隣の dozou_katanuki.wiki ディレクトリに同期するスクリプト (100行以下)
#>
$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$wikiSource = Join-Path $projectRoot "docs\wiki"
$wikiDest = Join-Path (Split-Path -Parent $projectRoot) "dozou_katanuki.wiki"

if (-not (Test-Path $wikiDest)) {
    Write-Host "[!] Target wiki repository not found at: $wikiDest" -ForegroundColor Yellow
    exit 1
}

Write-Host "[*] Syncing markdown files from docs/wiki to dozou_katanuki.wiki..." -ForegroundColor Cyan
$copiedCount = 0

Get-ChildItem -Path $wikiSource -Filter "*.md" | ForEach-Object {
    Copy-Item -Path $_.FullName -Destination $wikiDest -Force
    $copiedCount++
}

Write-Host "[+] Successfully copied $copiedCount markdown files to $wikiDest" -ForegroundColor Green

# Optional: Run combine_wiki if python is available
$combineScript = Join-Path $scriptDir "docs\combine_wiki.py"
if (Test-Path $combineScript) {
    Write-Host "[*] Regenerating COMBINED_WIKI.md..." -ForegroundColor Cyan
    python $combineScript
}
