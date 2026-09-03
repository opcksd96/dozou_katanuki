# scripts/run.ps1
[CmdletBinding()]
param ()

$ErrorActionPreference = "Stop"

Write-Host "==================================================="
Write-Host "  dozou_katanuki (土蔵・型抜き) Production Conductor"
Write-Host "==================================================="

# プロジェクトルートに移動
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = (Resolve-Path "$ScriptDir\..").Path
Set-Location $ProjectRoot

$global:managedProcesses = @()

function Cleanup {
    Write-Host "`n[*] Stopping managed processes..." -ForegroundColor Yellow
    foreach ($proc in $global:managedProcesses) {
        if ($proc -and !$proc.HasExited) {
            Write-Host "    Killing $($proc.ProcessName) (PID: $($proc.Id))"
            Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
        }
    }
    
    taskkill /F /IM stash-win.exe 2>$null
    taskkill /F /IM dozou_katanuki.exe 2>$null
    
    Write-Host "[OK] Cleanup completed." -ForegroundColor Green
}

[console]::TreatControlCAsInput = $false
Register-EngineEvent -SourceIdentifier Console.CancelKeyPress -Action {
    Cleanup
    exit 0
}

try {
    taskkill /F /IM stash-win.exe 2>$null
    taskkill /F /IM dozou_katanuki.exe 2>$null

    # Stash のパス解決
    $stashPath = "stash-win.exe"
    if (Test-Path "bin\stash\stash-win.exe") {
        $stashPath = "bin\stash\stash-win.exe"
    } elseif (Test-Path "bin\stash-win.exe") {
        $stashPath = "bin\stash-win.exe"
    }

    Write-Host "[*] Starting Stash-win ($stashPath)..." -ForegroundColor Cyan
    $stashProc = Start-Process -FilePath $stashPath -WindowStyle Hidden -PassThru
    $global:managedProcesses += $stashProc

    if (!(Test-Path "dozou_katanuki.exe")) {
        Write-Host "[ERROR] dozou_katanuki.exe not found. Please build first." -ForegroundColor Red
        exit 1
    }

    Write-Host "[*] Starting dozou_katanuki.exe..." -ForegroundColor Cyan
    $appProc = Start-Process -FilePath "dozou_katanuki.exe" -PassThru
    $global:managedProcesses += $appProc

    Write-Host "[*] System is running. Press Ctrl+C to stop all processes." -ForegroundColor White
    
    # メインプロセスが終了するまで待機
    $appProc.WaitForExit()

} finally {
    Cleanup
}
