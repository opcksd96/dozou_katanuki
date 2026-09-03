@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion
title dozou_katanuki (土蔵・型抜き) Launcher

echo ===================================================
echo   dozou_katanuki (土蔵・型抜き) ポータブル起動ランチャー
echo ===================================================

cd /d "%~dp0"

echo [*] 作業ディレクトリを初期化中...
if not exist "blobs" mkdir blobs
if not exist "backups\dumps" mkdir backups\dumps
if not exist "assets\avatars" mkdir assets\avatars
if not exist "cache" mkdir cache

echo [*] 既存の残存プロセスをクレンジング中...
taskkill /F /IM stash-win.exe >nul 2>&1

echo [*] Python 実行環境を確認中...
where python >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo [OK] システム Python が検出されました。
) else (
    where py >nul 2>&1
    if !ERRORLEVEL! equ 0 (
        echo [OK] Python ランチャー 'py' が検出されました。
    ) else (
        echo [!] 警告: Python が PATH に見つかりません。サイドカー収集機能を利用する場合は Python をインストールしてください。
    )
)

echo [*] Production Conductor (run.ps1) へ起動を委譲します...
if exist "scripts\run.ps1" (
    where pwsh >nul 2>&1
    if !ERRORLEVEL! equ 0 (
        pwsh -NoProfile -ExecutionPolicy Bypass -File "scripts\run.ps1"
    ) else (
        powershell -NoProfile -ExecutionPolicy Bypass -File "scripts\run.ps1"
    )
) else (
    echo [ERROR] scripts\run.ps1 が見つかりません。
    pause
    exit /b 1
)

exit /b %ERRORLEVEL%
