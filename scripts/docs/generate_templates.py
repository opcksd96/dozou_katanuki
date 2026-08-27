#!/usr/bin/env python3
import os
import sys

# 各アトミック仕様書のひな型定義 (100行以下の美しい構造)
TEMPLATES = {
    "part2_01_00_index.md": """# SPEC-ADMIN-000: Part 2-1 管理・設定・災害復旧 総合インデックス

## 1. 概要と目的
本ドキュメントは、システムにおける「設定一元化」「プロセスライフサイクル制御」「閉塞通信」「バックアップ・DR」を担う、Admin & Governance階層（第5層）の総合索引である。各詳細仕様はアトミックに分割され、相互参照される。

## 2. アトミック仕様リンク一覧（親・子・関連性）
*   [SPEC-CONFIG-001: 統合設定ファイル (config.json) 仕様](./part2_01_01_config.md)
*   [SPEC-LIFECYCLE-001: Wails-Stash 起動・終了プロセス制御仕様](./part2_01_02_lifecycle.md)
*   [SPEC-PROXY-001: Wailsインメモリプロキシ（閉塞通信）仕様](./part2_01_03_proxy.md)
*   [SPEC-ADMINBOARD-001: Settings UI (7大制御ビュー) 仕様](./part2_01_04_admin_board.md)
*   [SPEC-BACKUP-001: 二重化バックアップ（Dual-Source DR）仕様](./part2_01_05_backup.md)
*   [SPEC-RECOVERY-001: 災害復旧（オフライン自動リストア）手順](./part2_01_06_recovery.md)
*   [SPEC-AUDIT-001: SQLite3 整合性監査＆パージプロトコル](./part2_01_07_audit.md)
*   [SPEC-SCHEDULER-001: 常駐スケジューラー＆キャスト配信仕様](./part2_01_08_scheduler.md)
""",

    "part2_01_01_config.md": """# SPEC-CONFIG-001: 統合設定ポータル (config.json) 仕様

## 1. 設定構造ハイエラルキー (SSOT)
システム全体の設定は、ルート直下の `config.json` にて一元管理される。

root (config.json)
┣ 1. システム設定 (system) ➔ [詳細仕様](./part2_01_01_config.md#2-システム設定)
┣ 2. 表示・外観設定 (appearance)
┣ 3. フォルダ・ストレージ設定 (storage) ➔ [「Stash使わんし！」モード挙動](./part2_01_01_config.md#3-stash使わんしモード挙動)
┗ 4. 運行・タスクスケジュール設定 (scheduler)

## 2. システム設定 (system)
*   `stash_enabled`: boolean (デフォルト: true。false時はポータブルモード移行)
*   `stash_path`: string (デフォルト: "./bin/stash-win.exe")
*   `listen_port`: number (デフォルト: 9999)

## 3. 「Stash使わんし！」モード挙動
*   **物理保存**: `{local_media_dir}/{platform}/{username}/{media_id}` (URL BaseName命名原則準拠)
*   **相対パス動的解決**: ミドルウェアが `/media-local/{platform}/{username}/{media_id}` へURLを置換。
""",

    "part2_01_02_lifecycle.md": """# SPEC-LIFECYCLE-001: Wails-Stash プロセスライフサイクル制御仕様

## 1. 状態定義 (AppStatus)
*   `BOOTING`: Wails `OnStartup` 起動、前回の `stash-win.exe` ゾンビプロセスの事前駆逐。
*   `INITIALIZING`: Stashプロセスの非表示起動（`CREATE_NO_WINDOW`）、GraphQLヘルスチェック。
*   `READY`: Wailsイベント `stash:statusCONNECTED` 送信完了、UDFイベントループ稼働。
*   `TERMINATING`: Wails `OnShutdown` / `Process.Kill()` によるStash強制道連れ終了。

## 2. ゾンビ回避（Lifeline Sync）
*   Wails終了時に確実に子プロセスハンドルを回収。
*   Windows環境では `JobObject`（Windows API）をバインドし、親プロセスが急死した場合もカーネルレベルでStashを追従消滅させる。
""",

    "part2_01_03_proxy.md": """# SPEC-PROXY-001: Wailsインメモリプロキシ（閉塞通信）仕様

## 1. 外部プロキシポートの完全廃止
*   ローカルネットワークへの露出ポートを全廃し、セキュリティ衝突を防ぐ。
*   Wailsの `AssetHandler` を利用し、メモリ内でリクエストをインターセプトする。

## 2. インメモリ・リバースプロキシ挙動
*   `/stash-proxy/*` ➔ `http://127.0.0.1:9999/*` への透過中継。
*   ブラウザの Same-Origin Policy を満たすため、CORS制限ヘッダーを Go のメモリ内で自動付与。
""",

    "part2_01_04_admin_board.md": """# SPEC-ADMINBOARD-001: Settings UI (7大制御ビュー) 仕様

## 1. Dumb UI原則に基づく責務
フロントエンド（Vue 3）は一切の状態ロジックを持たず、Propsのレンダリングと、Go APIを呼ぶActionイベントの発行にのみ専念する。

## 2. 7大制御ビューの一覧
1.  **Scraper View**: 標準出力 `PROGRESS: {current}/{total} | {msg}` パース＆疑似ターミナル中継。
2.  **Whitelist管理**: 対象アカウント `whitelist` テーブル of CRUD操作。
3.  **Article Editor**: キャッシュ翻訳テキスト `full_text_ja/en/zh` の個別手動上書き。
4.  **Stashスマート別窓導線**: `http://127.0.0.1:9999/scenes/{stash_scene_id}` への `_blank` 誘導。
5.  **デフォルトCSSエディタ**: `plugins/{platform}/skin/design.css` のブラウザ経由物理上書き。
6.  **フォント微調整パネル**: Vue 3のCSSカスタム変数へ優先フォントをシグナルバインド。
7.  **「Stash使わんし！」モードトグル**: `storage.stash_enabled` フラグの即時トグラー。
""",

    "part2_01_05_backup.md": """# SPEC-BACKUP-001: 二重化バックアップ（Dual-Source DR）仕様

## 1. 二重系統データ保全（Dual-Source DR）
*   **Layer 1 (Fast Path)**: SQLite3のオンラインバックアップ（`VACUUM INTO`）を用い、無停止で整合コピー `backups/database/archive_YYYYMMDD_HHMMSS.db` を生成。
*   **Layer 2 (Deep Path)**: ポータビリティ不変原本。`metadata.json` と生パケット原本 `snapshot.warc.gz` を `backups/dumps/{platform}/{username}/{post_id}/` に永続化。

## 2. アバター完全隔離ポリシー (Avatar Isolation)
*   Stashライブラリ（Scene / Image）へのアバター画像混入を厳格に禁止（汚染パージ）。
*   アバター実体は `backups/dumps/{platform}/{username}/avatars/` に隔離。
""",

    "part2_01_06_recovery.md": """# SPEC-RECOVERY-001: 災害復旧（オフライン自動リストア）手順

## 1. ゼロからの復旧設計 (Disaster Recovery)
実稼働データベース `archive.db` が完全破壊された場合、Layer 2 (生データ原本) のみから完全オフラインで自動再構築を行う。

## 2. 対称リストアシーケンス
1.  **原本スキャン**: `backups/dumps/` 配下を走査し、各投稿の `metadata.json` を全ロード。
2.  **DB再構築 (Go/GORM)**: ロードしたデータから SQLite3 の schema / tables を自動修復・再バインド。
3.  **Stashインジェクション**: 実メディアファイルをStashへ自動再登録、GraphQL UUIDを逆引きバインド。
""",

    "part2_01_07_audit.md": """# SPEC-AUDIT-001: SQLite3 整合性監査＆パージプロトコル

## 1. データベース監査規約
*   `PRAGMA integrity_check;`: SQLite3のデータページ、B-Tree破損検知。
*   `PRAGMA foreign_key_check;`: テーブル間の外部キーリレーションの破綻を検知。

## 2. ゾンビメディア・キャッシュパージ
*   **SQLite3 孤立メディア削除**: `media` テーブルに存在するがStashapp側にない `stash_scene_id` レコードを自動パージ。
*   **Stash 孤立ファイルパージ**: `stash/scenes/` 物理ディレクトリ内を自動スキャンし、DBに登録されていない不要ファイルをOSのゴミ箱へ自動退避。
""",

    "part2_01_08_scheduler.md": """# SPEC-SCHEDULER-001: 常駐スケジューラー＆キャスト配信仕様

## 1. 常駐型ワーカースケジューラー
Wails起動時にGoバックエンド（Goroutine）として常駐。
*   **完了フォルダ巡回**: `scheduler.polling_interval_minutes` ごとにMotrix等の完了フォルダを走査し、自律的にStashにインジェクションを実行。
*   **自動オンラインバックアップ**: `scheduler.backup_interval_hours` 周期で SQLite3 の Layer 1 バックアップを実行、世代数管理。

## 2. メディア Broadcast（キャスト配信）
*   **バインド**: `network.public_bind_address`（デフォルト 0.0.0.0）で起動。
*   **IP制限**: Goの `net.IP.Contains` を用い、リクエスト元IPが `broadcast.allowed_networks`（CIDR形式）内にあるかを厳格に検証（不一致時は403）。
"""
}

def main():
    out_dir = "/workspace/scratch/docs"
    if len(sys.argv) > 1:
        out_dir = sys.argv[1]

    os.makedirs(out_dir, exist_ok=True)
    print(f"[*] Output directory configured: {out_dir}")

    for filename, content in TEMPLATES.items():
        filepath = os.path.join(out_dir, filename)
        lines = content.strip().split("\n")
        line_count = len(lines)
        
        print(f"[+] Generating {filename} ({line_count} lines)...")
        
        if line_count > 100:
            print(f"[!] Warning: {filename} exceeds 100 lines limit ({line_count} lines)!")
            
        with open(filepath, "w", encoding="utf-8") as f:
            f.write(content.strip() + "\n")
            
    print("[*] All templates generated successfully!")

if __name__ == "__main__":
    main()
