param (
    [string]$Version = "v4.1.0",
    [switch]$SkipBuild = $false,
    [switch]$CreateZip = $true
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir
Set-Location $RootDir

Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host "  dozou_katanuki ($Version) Portable Packaging" -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan

if (-not $SkipBuild) {
    Write-Host "`n[1/3] Building Frontend (Vue 3 / Vite)..." -ForegroundColor Yellow
    Set-Location (Join-Path $RootDir "frontend")
    npm run build
    Set-Location $RootDir

    Write-Host "`n[2/3] Building Application Binary..." -ForegroundColor Yellow
    if (Get-Command "wails" -ErrorAction SilentlyContinue) {
        wails build -clean -o dozou_katanuki.exe
    } else {
        go build -tags desktop,production -ldflags "-w -s -H windowsgui" -o dozou_katanuki.exe .
    }
}

Write-Host "`n[3/3] Assembling Portable Package..." -ForegroundColor Yellow
$OutDir = Join-Path $RootDir "release\dozou_katanuki_portable"
if (Test-Path $OutDir) {
    Remove-Item -Path $OutDir -Recurse -Force
}
New-Item -ItemType Directory -Path $OutDir -Force | Out-Null

$FilesToCopy = @(
    "dozou_katanuki.exe",
    "START_DOZOU.bat",
    "STOP_DOZOU.bat",
    "README_PORTABLE.txt",
    "config.json",
    "archive_schema.sql"
)

foreach ($f in $FilesToCopy) {
    $src = Join-Path $RootDir $f
    if (Test-Path $src) {
        Copy-Item -Path $src -Destination $OutDir -Force
        Write-Host "  + $f" -ForegroundColor Green
    } else {
        Write-Host "  ! Warning: $f not found" -ForegroundColor Magenta
    }
}

$DirsToCopy = @("bin", "plugins", "assets")
foreach ($d in $DirsToCopy) {
    $srcDir = Join-Path $RootDir $d
    if (Test-Path $srcDir) {
        Copy-Item -Path $srcDir -Destination (Join-Path $OutDir $d) -Recurse -Force
        Write-Host "  + $d/ (Directory)" -ForegroundColor Green
    }
}

$EmptyDirs = @(
    "blobs",
    "backups\dumps",
    "stash\scenes",
    "stash\images",
    "cache"
)
foreach ($ed in $EmptyDirs) {
    $targetPath = Join-Path $OutDir $ed
    if (-not (Test-Path $targetPath)) {
        New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
    }
}

Write-Host "`n[OK] Assembled Portable Package: $OutDir" -ForegroundColor Cyan

if ($CreateZip) {
    $ZipPath = Join-Path $RootDir "release\dozou_katanuki_${Version}_portable.zip"
    if (Test-Path $ZipPath) {
        Remove-Item -Path $ZipPath -Force
    }
    Write-Host "`n[*] Compressing ZIP Archive: $ZipPath" -ForegroundColor Yellow
    Compress-Archive -Path "$OutDir\*" -DestinationPath $ZipPath -CompressionLevel Optimal
    $ZipSizeMB = [math]::Round((Get-Item $ZipPath).Length / 1MB, 2)
    Write-Host "[OK] ZIP Archive Created ($ZipSizeMB MB)" -ForegroundColor Green
}

Write-Host "`n=========================================================" -ForegroundColor Cyan
Write-Host "  Packaging Completed Successfully!" -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan
