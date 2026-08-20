@echo off
title dozou_katanuki (土蔵・型抜き) Emergency Stopper

echo ===================================================
echo   dozou_katanuki (土蔵・型抜き) 一括終了スクリプト
echo ===================================================

cd /d "%~dp0"

echo [*] dozou_katanuki 関連プロセスを終了しています...
taskkill /F /IM dozou_katanuki.exe >nul 2>&1
taskkill /F /IM stash-win.exe >nul 2>&1

echo [OK] すべての関連プロセスを終了・クレンジングしました。
timeout /t 2 >nul
exit /b 0
