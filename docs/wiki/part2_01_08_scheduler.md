[← 07 DB健全性監査](part2_01_07_audit) | [📚 目次 (Home)](Home.md) | [第2編第2章：プラグインアーキテクチャとサイドカー →](part2_02_plugin_architecture)

# SPEC-SCHEDULER-001: 常駐スケジューラー＆キャスト配信仕様

## 1. cron-like 常駐型ワーカースケジューラー
Wails起動時にGo常駐Goroutineとして立ち上がる軽量スケジューラーです。
* **完了フォルダ自動巡回**: `scheduler.polling_interval_minutes` 周期でダウンロード完了フォルダをスキャンし、対象ファイルをStashへ自動インジェクション。
* **Layer 1 自動オンラインバックアップ**: `scheduler.backup_interval_hours` ごとに `VACUUM INTO` を実行し、世代数上限（`max_backup_files`）超過分を安全退避・パージ。

## 2. メディア Broadcast（家庭内LANキャスト）
* **ネットワークバインド**: `network.public_bind_address`（`0.0.0.0`）にバインドして家庭内デバイスへメディアを中継。
* **IP / CIDR サブネット制限**: 送信元IPが `broadcast.allowed_networks`（例: `192.168.1.0/24`）に合致するかを厳格に検証し、不正アクセスを即座に `403 Forbidden` で遮断。

---

[← 07 DB健全性監査](part2_01_07_audit) | [📚 目次 (Home)](Home) | [第2編第2章：プラグインアーキテクチャとサイドカー →](part2_02_plugin_architecture)
