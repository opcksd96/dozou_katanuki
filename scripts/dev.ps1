# scripts/dev.ps1
[CmdletBinding()]
param ()

$ErrorActionPreference = "Stop"

Write-Host "==================================================="
Write-Host "  dozou_katanuki (土蔵・型抜き) Development Conductor"
Write-Host "==================================================="

# プロジェクトルートに移動
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = (Resolve-Path "$ScriptDir\..").Path
Set-Location $ProjectRoot

# クリーンアップ処理の定義
$global:managedProcesses = @()

function Cleanup {
    Write-Host "`n[*] Stopping managed processes..." -ForegroundColor Yellow
    foreach ($proc in $global:managedProcesses) {
        if ($proc -and !$proc.HasExited) {
            Write-Host "    Killing $($proc.ProcessName) (PID: $($proc.Id))"
            Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
        }
    }
    
    # フォールバック: taskkill
    taskkill /F /IM stash-win.exe 2>$null
    
    Write-Host "[OK] Cleanup completed." -ForegroundColor Green
}

# Ctrl+C ハンドラの登録
[console]::TreatControlCAsInput = $false
$Host.UI.RawUI.FlushInputBuffer()
Register-EngineEvent -SourceIdentifier Console.CancelKeyPress -Action {
    Cleanup
    exit 0
}

try {
    # 既存の Stash プロセスをクリーンアップ
    taskkill /F /IM stash-win.exe 2>$null

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

    Write-Host "[*] Starting Wails Dev Server..." -ForegroundColor Cyan
    # wails dev はコンソールを占有するため直接実行。
    # 終了後、finallyブロックにて一括クリーンアップが走る。
    wails dev

} finally {
    Cleanup
}
