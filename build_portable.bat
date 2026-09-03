@echo off
chcp 65001 >nul
title dozou_katanuki Portable Builder

echo ===================================================
echo   dozou_katanuki (土蔵・型抜き) ポータブルビルド
echo ===================================================

cd /d "%~dp0"

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\package_portable.ps1"

if %ERRORLEVEL% equ 0 (
    echo.
    echo [SUCCESS] ビルドとパッケージングが完了しました。
    echo release\ ディレクトリをご確認ください。
) else (
    echo.
    echo [ERROR] ビルド中にエラーが発生しました。
)

pause
