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
    Get-Process -Name "stash-win", "stash-win.exe" -ErrorAction SilentlyContinue | Stop-Process -Force
    Get-Process -Name "dozou_katanuki", "dozou_katanuki.exe" -ErrorAction SilentlyContinue | Stop-Process -Force
    
    Write-Host "[OK] Cleanup completed." -ForegroundColor Green
}

[console]::TreatControlCAsInput = $false
Register-EngineEvent -SourceIdentifier Console.CancelKeyPress -Action {
    Cleanup
    exit 0
}

function Test-Port {
    param([int]$port)
    $tcp = New-Object System.Net.Sockets.TcpClient
    try {
        $task = $tcp.ConnectAsync("127.0.0.1", $port)
        if ($task.Wait(100)) {
            if ($tcp.Connected) { return $true }
        }
    } catch {
        # ignore
    } finally { 
        $tcp.Dispose() 
    }
    return $false
}

function Send-Beacon {
    param([string]$service, [string]$status)
    $payload = @{
        service = $service
        status = $status
    } | ConvertTo-Json -Compress
    
    try {
        Invoke-RestMethod -Uri "http://127.0.0.1:5175/api/internal/beacon" -Method Post -Body $payload -ContentType "application/json" -ErrorAction SilentlyContinue | Out-Null
    } catch {
        # ignore
    }
}

try {
    Get-Process -Name "stash-win", "stash-win.exe" -ErrorAction SilentlyContinue | Stop-Process -Force
    Get-Process -Name "dozou_katanuki", "dozou_katanuki.exe" -ErrorAction SilentlyContinue | Stop-Process -Force

    $stashPath = "stash-win.exe"
    $stashDir = $ProjectRoot
    if (Test-Path "bin\stash\stash-win.exe") {
        $stashPath = "bin\stash\stash-win.exe"
        $stashDir = "$ProjectRoot\bin\stash"
    } elseif (Test-Path "bin\stash-win.exe") {
        $stashPath = "bin\stash-win.exe"
        $stashDir = "$ProjectRoot\bin"
    }

    Write-Host "[*] Starting Stash-win ($stashPath) in ($stashDir)..." -ForegroundColor Cyan
    $stashProc = Start-Process -FilePath $stashPath -WorkingDirectory $stashDir -WindowStyle Hidden -PassThru
    $global:managedProcesses += $stashProc

    if (!(Test-Path "dozou_katanuki.exe")) {
        Write-Host "[ERROR] dozou_katanuki.exe not found. Please build first." -ForegroundColor Red
        exit 1
    }

    Write-Host "[*] Starting dozou_katanuki.exe..." -ForegroundColor Cyan
    $appProc = Start-Process -FilePath "dozou_katanuki.exe" -PassThru
    $global:managedProcesses += $appProc

    Write-Host "[*] System is running. Beacon loop activated. Press Ctrl+C to stop all processes." -ForegroundColor White
    
    # メインプロセスが終了するまでビーコンループ
    while (!$appProc.HasExited) {
        # Check Stash (:9999)
        if (Test-Port 9999) { Send-Beacon "stash" "ready" } else { Send-Beacon "stash" "busy" }
        # Check Thunder CDP (:9222)
        if (Test-Port 9222) { Send-Beacon "thunder" "ready" } else { Send-Beacon "thunder" "stopped" }
        # Check Motrix (:16800)
        if (Test-Port 16800) { Send-Beacon "motrix" "ready" } else { Send-Beacon "motrix" "stopped" }

        Start-Sleep -Seconds 5
    }
} catch {
    Write-Host "[ERROR] $($_.Exception.Message)" -ForegroundColor Red
} finally {
    Cleanup
}
