@echo off
chcp 65001 >nul
title dozou_katanuki (土蔵・型抜き) Emergency Stopper

echo ===================================================
echo   dozou_katanuki (土蔵・型抜き) 緊急停止ツール
echo ===================================================
echo.
echo バックグラウンドプロセス (Stash 等) を強制終了します...

taskkill /F /IM stash-win.exe >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo [OK] stash-win.exe を停止しました。
) else (
    echo [-] stash-win.exe は実行されていません。
)

taskkill /F /IM dozou_katanuki.exe >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo [OK] dozou_katanuki.exe を停止しました。
) else (
    echo [-] dozou_katanuki.exe は実行されていません。
)

echo.
echo 停止プロセスが完了しました。
pause
