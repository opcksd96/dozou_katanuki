@echo off
echo ===================================================
echo   dozou_katanuki Hot Reload Dev Launcher
echo ===================================================
echo.
echo [*] Checking external processes...

:: Stash check
tasklist /FI "IMAGENAME eq stash-win.exe" 2>NUL | find /I /N "stash-win.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [OK] Stash is running.
) else (
    echo [!] Stash is NOT running. Starting it in the background...
    if exist "bin\stash\stash-win.exe" (
        start "Stash" /D "bin\stash" /MIN "stash-win.exe"
    ) else if exist "bin\stash-win.exe" (
        start "Stash" /D "bin" /MIN "stash-win.exe"
    ) else (
        echo [ERROR] Stash binary not found.
    )
)

:: Thunder check
tasklist /FI "IMAGENAME eq Thunder.exe" 2>NUL | find /I /N "Thunder.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [OK] Thunder is running.
) else (
    echo [!] Thunder is NOT running. (Please kick it manually or via app if needed)
)

:: Motrix check
tasklist /FI "IMAGENAME eq Motrix.exe" 2>NUL | find /I /N "Motrix.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [OK] Motrix Next is running.
) else (
    echo [!] Motrix Next is NOT running. (Please kick it manually if needed)
)

echo.
echo [*] Starting Backend (Air) and Frontend (Vite) in new windows...
start "Katanuki Backend (Air)" cmd /k "air"
start "Katanuki Frontend (Vite)" cmd /k "cd frontend && npm run dev"

echo.
echo [OK] Development environments launched!
echo Please access http://localhost:5173 in your browser after a few seconds.
pause
