@echo off
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
        echo [OK] Python ランチャー (py) が検出されました。
    ) else (
        echo [!] 警告: Python が PATH に見つかりません。サイドカー収集機能を利用する場合は Python をインストールしてください。
    )
)

if not exist "dozou_katanuki.exe" (
    echo [ERROR] dozou_katanuki.exe が見つかりません。
    pause
    exit /b 1
)

echo [*] dozou_katanuki を起動しています...
start "" "dozou_katanuki.exe"

exit /b 0
