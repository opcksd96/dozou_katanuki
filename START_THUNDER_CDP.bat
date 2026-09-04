@echo off
setlocal
echo ===================================================
echo   Thunder CDP (Port 9222) Launcher
echo ===================================================
echo.

echo [*] Checking if Thunder is already running with CDP...

set "PS_CHECK=$proc = Get-CimInstance Win32_Process -Filter \"Name='Thunder.exe'\" -ErrorAction SilentlyContinue | Select-Object -First 1; if ($proc) { if ($proc.CommandLine -match '9222') { exit 0 } else { exit 2 } } else { exit 1 }"
powershell -NoProfile -Command "%PS_CHECK%"
set PS_EXIT=%errorlevel%

if %PS_EXIT% equ 0 (
    echo [OK] Thunder is already running in CDP mode (Port 9222).
    exit /b 0
) else if %PS_EXIT% equ 2 (
    echo [*] Thunder is running, but WITHOUT CDP enabled. Terminating to restart with CDP...
    taskkill /F /IM Thunder.exe /T 2>NUL
    taskkill /F /IM ThunderApp.exe /T 2>NUL
    timeout /t 2 /nobreak >NUL
)

echo [*] Resolving Thunder installation path...
:: PowerShell script to find Thunder path from registry or running process
set "PS_CMD=$path = ''; $proc = Get-CimInstance Win32_Process -Filter \"Name='Thunder.exe'\" -ErrorAction SilentlyContinue | Select-Object -First 1; if ($proc) { $path = $proc.ExecutablePath } else { $regPaths = @('HKCU:\Software\Thunder Network\ThunderOem\thunder_backwnd', 'HKLM:\SOFTWARE\WOW6432Node\Thunder Network\ThunderOem\thunder_backwnd', 'HKCU:\Software\Thunder Network\Thunder\Path'); foreach ($rp in $regPaths) { $val = Get-ItemPropertyValue -Path $rp -Name 'Path' -ErrorAction SilentlyContinue; if ($val -and (Test-Path $val)) { $path = $val; break } }; if (-not $path) { $default = 'C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe'; if (Test-Path $default) { $path = $default } } }; Write-Output $path"

for /f "usebackq tokens=*" %%A in (`powershell -NoProfile -Command "%PS_CMD%"`) do set "THUNDER_EXE=%%A"

if "%THUNDER_EXE%"=="" (
    echo [ERROR] Thunder executable not found. Please install Thunder or start it manually with --remote-debugging-port=9222.
    exit /b 1
)

echo [*] Thunder path resolved: %THUNDER_EXE%

echo [*] Terminating existing Thunder processes...
taskkill /F /IM Thunder.exe /T 2>NUL
taskkill /F /IM ThunderApp.exe /T 2>NUL

:: Wait a moment for processes to exit fully
timeout /t 2 /nobreak >NUL

echo [*] Starting Thunder with CDP enabled...
start "Thunder" "%THUNDER_EXE%" --remote-debugging-port=9222

echo [OK] Thunder is now running with CDP enabled.
exit /b 0
