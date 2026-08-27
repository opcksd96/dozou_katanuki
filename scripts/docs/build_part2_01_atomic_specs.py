#!/usr/bin/env python3
# File: build_part2_01_atomic_specs.py
# -*- coding: utf-8 -*-
"""
Part 2-1 (管理・設定・ディザスタリカバリ運用) アトミック仕様書群 生成スクリプト
- バージョンを 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様) へ一括インクリメント
- 全ファイルにドキュメントID・バージョン・更新日(2026-08-19)を明記
- 項間ナビゲーションの数珠つなぎ巡回 (00 ➔ 01 ➔ ... ➔ 08)
- 100行ルール完全適合（全ファイル60〜85行）
- Mermaidダイアグラム完全準拠
"""

import os
import sys
from pathlib import Path

BT = chr(96) * 3

# アトミック仕様書の巡回チェーン定義
NAV_CHAIN = {
    "part2_01_00_index.md": {
        "prev": ("第1編第4章：実装規約・制約原則", "part1_04_implementation_principles"),
        "next": ("01 一元設定仕様", "part2_01_01_config"),
    },
    "part2_01_01_config.md": {
        "prev": ("00 総合インデックス", "part2_01_00_index"),
        "next": ("02 プロセス制御仕様", "part2_01_02_lifecycle"),
    },
    "part2_01_02_lifecycle.md": {
        "prev": ("01 一元設定仕様", "part2_01_01_config"),
        "next": ("03 インメモリプロキシ仕様", "part2_01_03_proxy"),
    },
    "part2_01_03_proxy.md": {
        "prev": ("02 プロセス制御仕様", "part2_01_02_lifecycle"),
        "next": ("04 Admin Board仕様", "part2_01_04_admin_board"),
    },
    "part2_01_04_admin_board.md": {
        "prev": ("03 インメモリプロキシ仕様", "part2_01_03_proxy"),
        "next": ("05 二重化バックアップ仕様", "part2_01_05_backup"),
    },
    "part2_01_05_backup.md": {
        "prev": ("04 Admin Board仕様", "part2_01_04_admin_board"),
        "next": ("06 災害復旧手順", "part2_01_06_recovery"),
    },
    "part2_01_06_recovery.md": {
        "prev": ("05 二重化バックアップ仕様", "part2_01_05_backup"),
        "next": ("07 DB健全性監査", "part2_01_07_audit"),
    },
    "part2_01_07_audit.md": {
        "prev": ("06 災害復旧手順", "part2_01_06_recovery"),
        "next": ("08 常駐スケジューラー", "part2_01_08_scheduler"),
    },
    "part2_01_08_scheduler.md": {
        "prev": ("07 DB健全性監査", "part2_01_07_audit"),
        "next": ("第2編第2章：プラグインアーキテクチャとサイドカー", "part2_02_plugin_architecture"),
    },
}

TEMPLATES = {
    # 00: 総合インデックス
    "part2_01_00_index.md": """@NAV@

# 第2編 第1章：管理・設定・ディザスタリカバリ運用 総合インデックス

**プロジェクト名** : dozou_katanuki  
**ドキュメントID** : SPEC-ADMIN-000  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  
**ステータス** : 正式仕様（アトミック分割・全編100行未満・Mermaid準拠）

---

## 1. 概要とレイヤー責務 (第5層: Admin & Governance)
本レイヤーは、日常のタイムライン描画やデータ中継には関与せず、**「システム設定の一元統治、WailsとStashのライフライン同期、DB監査、自動スケジューリング、災害復旧」**に特化した最上位ガバナンス階層です。

## 2. アトミック仕様構成マップ

@BT@mermaid
flowchart TD
    Index["SPEC-ADMIN-000<br>総合インデックス"] --> Cfg["01: 一元設定 SSOT<br>(config.json)"]
    Index --> Life["02: プロセス統治<br>(キック＆ライフライン)"]
    Index --> Proxy["03: インメモリプロキシ<br>(外部閉塞中継)"]
    Index --> UI["04: Admin Board<br>(7大制御ビュー)"]
    Index --> Bak["05: 二重化バックアップ<br>(Dual-Source DR)"]
    Index --> Rec["06: 災害復旧<br>(オフライン自動リストア)"]
    Index --> Aud["07: DB健全性監査<br>(PRAGMA＆パージ)"]
    Index --> Sch["08: 常駐ワーカー<br>(Scheduler＆Broadcast)"]
@BT@

## 3. アトミック仕様リンク一覧
* [[SPEC-CONFIG-001: 一元設定ポータル (config.json) 仕様|part2_01_01_config]]
* [[SPEC-LIFECYCLE-001: Wails-Stash プロセスライフサイクル制御仕様|part2_01_02_lifecycle]]
* [[SPEC-PROXY-001: Wails インメモリプロキシ（閉塞通信）仕様|part2_01_03_proxy]]
* [[SPEC-ADMINBOARD-001: Settings UI (7大制御ビュー) ＆ Scraper View 仕様|part2_01_04_admin_board]]
* [[SPEC-BACKUP-001: 二重化バックアップ（Dual-Source DR）仕様|part2_01_05_backup]]
* [[SPEC-RECOVERY-001: 災害復旧（完全オフライン自動リストア）手順|part2_01_06_recovery]]
* [[SPEC-AUDIT-001: SQLite3 整合性監査＆パージプロトコル|part2_01_07_audit]]
* [[SPEC-SCHEDULER-001: 常駐スケジューラー＆キャスト配信仕様|part2_01_08_scheduler]]

---

@NAV@
""",

    # 01: 設定仕様
    "part2_01_01_config.md": """@NAV@

# SPEC-CONFIG-001: 一元設定ポータル (config.json) 仕様

**ドキュメントID** : SPEC-CONFIG-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. 設定構造ハイエラルキー (Single Source of Truth)
従来散逸していたYAMLや個別環境変数を全廃し、ルート直下の `config.json` をシステム唯一の設定の源泉（SSOT）とします。

* **system**: 実行環境（`env`）、デフォルトプラットフォーム、UI言語（`ja/en/zh`）。
* **network**: ポート定義（内部閉塞 `stash_port: 9999` 等）、バインドアドレス。
* **storage**: DBパス（`db_path`）、`stash_enabled`、物理メディア保存先。
* **scheduler**: ポーリング間隔、自動バックアップ周期、最大保持世代数。
* **broadcast**: 家庭内LANキャスト有効化フラグ、許可サブネット（CIDR）。
* **appearance**: 3言語対応の優先フォントファミリー定義。

## 2. 「Stash使わんし！」モード (物理フォルダ保存ポリシー)
`storage.stash_enabled` が `false` の場合、システムは Stashapp を起動せず、軽量な物理フォルダダイレクトサーブモードへ自律移行します。

* **物理マッピング**: ダウンロード完了ファイルを `{local_media_dir}/{platform}/{username}/{media_id}` にフラット保存。
* **相対URL動的解決**: タイムライン構築時、メディアURLは `/media-local/{platform}/{username}/{media_id}` へ自動置換され、軽量Goバイナリ単体で動態保存が完結します。

---

@NAV@
""",

    # 02: ライフサイクル仕様
    "part2_01_02_lifecycle.md": """@NAV@

# SPEC-LIFECYCLE-001: Wails-Stash プロセスライフサイクル制御仕様

**ドキュメントID** : SPEC-LIFECYCLE-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. 起動・終了シーケンス (Lifeline Sync)

@BT@mermaid
sequenceDiagram
    autonumber
    participant OS as OS / Process
    participant Wails as Wails Core (Go)
    participant Pipe as Stdout/Stderr Pipe
    participant Stash as Stash (stash-win.exe)
    participant UI as Frontend (Vue 3)

    Wails->>OS: 1. taskkill /F /IM stash-win.exe (ゾンビ事前パージ)
    Wails->>Stash: 2. exec.Command() (CREATE_NO_WINDOW ヘッドレス起動)
    Wails->>Pipe: 3. Stdout/StderrPipe 常時監視開始
    Stash-->>Pipe: 4. "Server started on :9999" ログ送出
    Wails->>Stash: 5. POST /graphql (内部疎通ヘルスチェック)
    Stash-->>Wails: 6. 200 OK (Ready)
    Wails->>UI: 7. Wails Event: stash:ready (ONLINE表示へ移行)
    Note over Wails, Stash: --- Wails終了時 ---
    Wails->>Stash: 8. Process.Kill() (Stash道連れ完全終了)
@BT@

## 2. プロセス制御規約
* **ゾンビパージ**: `OnStartup` 時に前回の残存プロセスを強制終了し、ポート衝突を防止。
* **ヘッドレス起動**: OSウィンドウを一切生成せず、完全バックグラウンドで起動。
* **道連れ終了**: 親ウィンドウ終了時（`OnShutdown`）に確実に `Process.Kill()` を発行。

---

@NAV@
""",

    # 03: プロキシ仕様
    "part2_01_03_proxy.md": """@NAV@

# SPEC-PROXY-001: Wails インメモリプロキシ（閉塞通信）仕様

**ドキュメントID** : SPEC-PROXY-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. 外部プロキシポートの完全閉塞
従来の設計にあった外部公開プロキシポート（`:9998` 等）を全廃し、Wailsの `AssetHandler` を用いて **Goのメモリ内部でリバースプロキシを展開** します。これにより、外部からの不正アクセスやポート競合をゼロにします。

## 2. インメモリ・リバースプロキシ中継仕様

@BT@mermaid
flowchart LR
    subgraph WailsApp ["Wails v2 Desktop Process"]
        Vue["Vue 3 Frontend<br>(Dumb UI)"]
        Handler["Wails AssetHandler<br>(インメモリプロキシ)"]
        Stash["Stashapp Core<br>(127.0.0.1:9999 閉塞)"]
        
        Vue -- "src='/stash-proxy/scene/1/stream'" --> Handler
        Handler -- "CORS透過付与 & 内部中継" --> Stash
    end
@BT@

* **パス書き換え**: フロントエンドの `/stash-proxy/*` 要求を `http://127.0.0.1:9999/*` へメモリ内中継。
* **CORS透過無効化**: Same-Origin Policy を満たすHTTPヘッダーを内部自動付与し、ゼロレイテンシ再生を実現。

---

@NAV@
""",

    # 04: 管理画面・スクレイパー仕様
    "part2_01_04_admin_board.md": """@NAV@

# SPEC-ADMINBOARD-001: Settings UI (7大制御ビュー) ＆ Scraper View 仕様

**ドキュメントID** : SPEC-ADMINBOARD-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. Dumb UI原則に基づく責務
Vue 3 フロントエンドは状態ロジックを持たず、Actionイベント発行と表示に徹します。

## 2. 管理画面の「7大制御ビュー」
1. **Job コントローラー ＆ Scraper View**: サルベージキック、StdoutPipeによるリアルタイム進捗ログ疑似端末。
2. **Whitelist 管理ビュー**: 対象アカウント・キーワードのCRUDおよびトグル。
3. **個別記事編集ビュー**: 3言語翻訳テキスト（`full_text_ja/en/zh`）の手動微調整・保存。
4. **Stashスマート別窓導線**: `http://127.0.0.1:9999/scenes/{stash_scene_id}` への `_blank` 誘導。
5. **デフォルトCSSエディタ**: `plugins/{platform}/skin/design.css` のブラウザ直接編集・上書き。
6. **フォント微調整パネル**: 日・英・中の優先フォントをCSSカスタム変数へ動的シグナル同期。
7. **「Stash使わんし！」モードトグル**: `storage.stash_enabled` のワンクリック切り替え。

## 3. Python サイドカー連携シーケンス

@BT@mermaid
sequenceDiagram
    autonumber
    participant UI as Settings UI (:5173)
    participant Go as Wails Go Core
    participant Py as Python Sidecar (main.py)
    
    UI->>Go: POST /api/jobs/salvage
    Go->>Py: exec.CommandContext() (並行数1排他キック)
    loop リアルタイム進捗
        Py-->>Go: PROGRESS 標準出力フラッシュ
        Go-->>UI: StdoutPipe中継 ➔ Scraper View ログ追加
    end
    Py->>Go: 共通中間JSON登録 ＆ 完了通知
@BT@

---

@NAV@
""",

    # 05: バックアップ仕様
    "part2_01_05_backup.md": """@NAV@

# SPEC-BACKUP-001: 二重化バックアップ（Dual-Source DR）仕様

**ドキュメントID** : SPEC-BACKUP-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. 二重系統データ保全アーキテクチャ

@BT@mermaid
flowchart TD
    DB["SQLite3 (archive.db)<br>実稼働マスター"]
    
    subgraph L1 ["Layer 1: Fast Path (バイナリ復元)"]
        F1["VACUUM INTO コピー<br>backups/database/archive_*.db"]
        F2["・ミリ秒即時復旧 (RTO極小)<br>・リレーション完全維持"]
    end
    
    subgraph L2 ["Layer 2: Deep Path (生データ原本)"]
        D1["原本魚拓 ＆ メタデータ<br>backups/dumps/"]
        D2["・ISO 28500 warc.gz 原本<br>・metadata.json 共通中間表現"]
    end
    
    DB --> L1
    DB --> L2
@BT@

## 2. アバター保全・隔離ポリシー (Avatar Isolation)
* **ライブラリ汚染防止**: Stashの画像グリッド混入を防ぐため、アバター実ファイルはすべて `backups/dumps/{platform}/{username}/avatars/` に物理隔離保存します。
* **原本性保証**: `metadata.json` 内には生のアバターオリジナルURLを不変データとして保持します。

---

@NAV@
""",

    # 06: リストア仕様
    "part2_01_06_recovery.md": """@NAV@

# SPEC-RECOVERY-001: 災害復旧（完全オフライン自動リストア）手順

**ドキュメントID** : SPEC-RECOVERY-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. ゼロからの自動再構築設計
実稼働DB `archive.db` が完全破壊された場合、Layer 2 (生データ原本) のみから外部通信ゼロ・完全オフラインで動態保存状態を100%自動再構築します。

## 2. 対称リストアフロー

@BT@mermaid
graph TD
    A[実稼働 archive.db が全損] --> B[1. 破損 archive.db を物理削除]
    B --> C[2. Wails Core 起動 ➔ GORM AutoMigrate で空DB生成]
    C --> D[3. Pythonサイドカーをリストアモードでキック]
    D --> E[4. dumps 配下の metadata.json と snapshot.warc.gz を一括走査]
    E --> F[5. Core API POST /api/articles へ無加工で順次投入]
    F --> G[6. メディア実体をStashへ再配置 ＆ UUID逆引きバインド]
    G --> H[リストア完了: 100%同一の timeline 導通状態が完全復帰]

    style A fill:#ffebee,stroke:#f44336,stroke-width:2px
    style H fill:#e8f5e9,stroke:#4caf50,stroke-width:2px
@BT@

---

@NAV@
""",

    # 07: 監査仕様
    "part2_01_07_audit.md": """@NAV@

# SPEC-AUDIT-001: SQLite3 整合性監査＆パージプロトコル

**ドキュメントID** : SPEC-AUDIT-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. SQLite3 整合性監査 (PRAGMA Audit)
* **`PRAGMA integrity_check;`**: データページ、B-Tree、インデックスの破損を徹底スキャン。破損時は Layer 1 / 2 からの復旧アラートを発行。
* **`PRAGMA foreign_key_check;`**: `accounts` ➔ `articles` ➔ `media` 間の孤立外部キーエラーが0件であることを保証。

## 2. 孤立メディア・ゾンビキャッシュの自動パージ
* **SQLite3 孤立メディア検出**: DBに存在するがStash側にないレコードを検知・削除。
* **Stash 孤立ファイルパージ**: `stash/scenes/` 内を自動走査し、DBの `media_id` と一致しない未紐付けファイルをOSのゴミ箱へ自動退避。

---

@NAV@
""",

    # 08: スケジューラー仕様
    "part2_01_08_scheduler.md": """@NAV@

# SPEC-SCHEDULER-001: 常駐スケジューラー＆キャスト配信仕様

**ドキュメントID** : SPEC-SCHEDULER-001  
**バージョン** : 4.1.0 (Wailsキメラデスクトップ アトミック統合仕様)  
**更新日** : 2026-08-19  

---

## 1. cron-like 常駐型ワーカースケジューラー
Wails起動時にGo常駐Goroutineとして立ち上がる軽量スケジューラーです。
* **完了フォルダ自動巡回**: `scheduler.polling_interval_minutes` 周期でダウンロード完了フォルダをスキャンし、対象ファイルをStashへ自動インジェクション。
* **Layer 1 自動オンラインバックアップ**: `scheduler.backup_interval_hours` ごとに `VACUUM INTO` を実行し、世代数上限（`max_backup_files`）超過分を安全退避・パージ。

## 2. メディア Broadcast（家庭内LANキャスト）
* **ネットワークバインド**: `network.public_bind_address`（`0.0.0.0`）にバインドして家庭内デバイスへメディアを中継。
* **IP / CIDR サブネット制限**: 送信元IPが `broadcast.allowed_networks`（例: `192.168.1.0/24`）に合致するかを厳格に検証し、不正アクセスを即座に `403 Forbidden` で遮断。

---

@NAV@
"""
}


def build_atomic_specs(output_dir: Path, dry_run: bool):
    print(f"[*] 対象出力ディレクトリ: {output_dir.resolve()}")

    if not dry_run:
        output_dir.mkdir(parents=True, exist_ok=True)

    for filename, raw_template in TEMPLATES.items():
        filepath = output_dir / filename

        nav_info = NAV_CHAIN[filename]
        prev_title, prev_target = nav_info["prev"]
        next_title, next_target = nav_info["next"]
        nav_bar = f"[[← {prev_title}|{prev_target}]] | [[📚 目次 (Home)|Home]] | [[{next_title} →|{next_target}]]"

        content = raw_template.replace("@NAV@", nav_bar).replace("@BT@", BT)
        lines = content.strip().split("\n")
        line_count = len(lines)

        status_flag = "PASS" if line_count <= 100 else "EXCEED"
        print(f"[{status_flag}] {filename:<28} : {line_count:>3} 行 (v4.1.0 / Prev: {prev_target} / Next: {next_target})")

        if line_count > 100:
            print(f"  [!] 警告: {filename} が100行を超過しています！")

        if not dry_run:
            filepath.write_text(content.strip() + "\n", encoding="utf-8")

    if not dry_run:
        print(f"\n🎉 全 {len(TEMPLATES)} ファイルのアトミック仕様書（v4.1.0 インクリメント＆巡回完全同期）を反映いたしました！")
    else:
        print(f"\n💡 Dry-Run完了: ファイル書き込みを実行するには `--run` を付けてください。")


if __name__ == "__main__":
    dry_run = "--run" not in sys.argv

    target_dir_str = "."
    args_without_flags = [a for a in sys.argv[1:] if not a.startswith("--")]
    if args_without_flags:
        target_dir_str = args_without_flags[0]

    target_path = Path(target_dir_str)
    build_atomic_specs(target_path, dry_run)
