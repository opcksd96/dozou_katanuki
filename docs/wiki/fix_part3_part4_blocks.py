#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
dozou_katanuki 【Part 3 & Part 4】ヘッダー/フッターブロック完全置換スクリプト
- 先頭のタイトル/メタデータ/ナビゲーション領域を完全なブロックとして丸ごと差し替え
- 末尾のフッターナビゲーションも完全一致で同期
"""

import sys
import re
from pathlib import Path

DOCS_DIR = Path(".")

# 各ファイルの「完全ヘッダーブロック」および「フッターナビゲーション」定義
TARGET_BLOCKS = {
    # 1. 第3編 第1章：データベース設計
    "part3_01_database_design.md": {
        "header": """[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]

# 第3編 第1章：データベース設計と仮想ストレージプール (Database Design & Virtual Storage Pool)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-DATABASE-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（GORMモデル完全同期・SQLite3 DDL・10大インデックス最適化・WALモード競合制御定義）""",
        "footer": "[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]"
    },

    # 2. 第3編 第2章：フロントエンド層概論
    "part3_02_pure_dumb_frontend.md": {
        "header": """[[← 前の章: 第3編第1章：データベース設計と仮想ストレージプール|part3_01_database_design]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第3章：ミドルウェア層インデックス →|part3_03_0_middleware_index]]

# 第3編 第2章：フロントエンド層概論（Pure Dumb UI Framework）

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-FRONTEND-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（Dumb UI一般化・シグナル・UDF原則、設定管理系コンポーネント、SVGプレースホルダーフィラー出し分け）""",
        "footer": "[[← 前の章: 第3編第1章：データベース設計と仮想ストレージプール|part3_01_database_design]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第3章：ミドルウェア層インデックス →|part3_03_0_middleware_index]]"
    },

    # 3. 第3編 第3章：ミドルウェア層インデックス
    "part3_03_0_middleware_index.md": {
        "header": """[[← 前の章: 第3編第2章：フロントエンド層概論|part3_02_pure_dumb_frontend]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第4章：ドライバー層 →|part3_04_backend_driver]]

# 第3編 第3章：ミドルウェア層インデックス (Intelligent Hub Architecture)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-MIDDLEWARE-000  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（Wails v2 キメラアーキテクチャ準拠・プロキシ層委譲完了・サブファイル細分化）""",
        "footer": "[[← 前の章: 第3編第2章：フロントエンド層概論|part3_02_pure_dumb_frontend]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第4章：ドライバー層 →|part3_04_backend_driver]]"
    },

    # 4. 第3編 第3章.1：Middleware Core Components
    "part3_03_1_middleware_core.md": {
        "header": """[[← インデックス: 第3編第3章 ミドルウェア層 ポータル|part3_03_0_middleware_index]] | [[📚 目次 (Home)|Home]] | [[次の節: 3.2 Data & Skin Decorator →|part3_03_2_data_decorator]]

# 第3編 第3章.1：Middleware Core Components（要求オーケストレーションとリスト統治）

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-MIDDLEWARE-001-1  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（Wails v2 キメラアーキテクチャ・要求終端・中間JSONリスト統治・フォールトトレラント反復キュー純化）""",
        "footer": "[[← インデックス: 第3編第3章 ミドルウェア層 ポータル|part3_03_0_middleware_index]] | [[📚 目次 (Home)|Home]] | [[次の節: 3.2 Data & Skin Decorator →|part3_03_2_data_decorator]]"
    },

    # 5. 第3編 第3章.2：Data & Skin Decorator
    "part3_03_2_data_decorator.md": {
        "header": """[[← 前の節: 3.1 Middleware Core Components|part3_03_1_middleware_core]] | [[📚 目次 (Home)|Home]] | [[次の節: 3.3 Job & Process Orchestrator →|part3_03_3_job_orchestrator]]

# 第3編 第3章.2：Data & Skin Decorator（描画データ装飾とスキン配信）

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-MIDDLEWARE-001-2  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（Wails v2 キメラアーキテクチャ・プラグインインターフェース・Skin配信統合・ゼロレイテンシ純化）""",
        "footer": "[[← 前の節: 3.1 Middleware Core Components|part3_03_1_middleware_core]] | [[📚 目次 (Home)|Home]] | [[次の節: 3.3 Job & Process Orchestrator →|part3_03_3_job_orchestrator]]"
    },

    # 6. 第3編 第3章.3：Job & Process Orchestrator
    "part3_03_3_job_orchestrator.md": {
        "header": """[[← 前の節: 3.2 Data & Skin Decorator|part3_03_2_data_decorator]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第4章：ドライバー層 →|part3_04_backend_driver]]

# 第3編 第3章.3：Job & Process Orchestrator（管理コンソール向け非同期ジョブ制御）

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-MIDDLEWARE-001-3  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（Wails v2 キメラアーキテクチャ・管理コンソール専用・非同期サブプロセス制御純化）""",
        "footer": "[[← 前の節: 3.2 Data & Skin Decorator|part3_03_2_data_decorator]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第4章：ドライバー層 →|part3_04_backend_driver]]"
    },

    # 7. 第3編 第4章：ドライバー層
    "part3_04_backend_driver.md": {
        "header": """[[← 前の章: 第3編第3章：ミドルウェア層インデックス|part3_03_0_middleware_index]] | [[📚 目次 (Home)|Home]] | [[次の章: 第4編第1章：参考資料・技術文献・型定義カタログ・公式リンク集 →|part4_01_references_and_literature]]

# 第3編 第4章：ドライバー層（Core Backend Driver & Data Abstraction）

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-DRIVER-001  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（GORMモデル完全一般化・多言語事前翻訳カラム・純粋ストレージドライバ純化）""",
        "footer": "[[← 前の章: 第3編第3章：ミドルウェア層インデックス|part3_03_0_middleware_index]] | [[📚 目次 (Home)|Home]] | [[次の章: 第4編第1章：参考資料・技術文献・型定義カタログ・公式リンク集 →|part4_01_references_and_literature]]"
    },

    # 8. 第4編 第1章：参考資料・技術文献・型定義カタログ
    "part4_01_references_and_literature.md": {
        "header": """[[← 前の章: 第3編第4章：ドライバー層|part3_04_backend_driver]] | [[📚 目次 (Home)|Home]] | [[DocWiki ポータルへ戻る →|Home]]

# 第4編 第1章：参考資料・技術文献・型定義カタログ・公式リンク集

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-REFERENCE-001  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（パラダイム・規約・言語・モジュール・外部アプリ・Webサービス一元技術カタログ化）""",
        "footer": "[[← 前の章: 第3編第4章：ドライバー層|part3_04_backend_driver]] | [[📚 目次 (Home)|Home]] | [[DocWiki ポータルへ戻る →|Home]]"
    },
}

def is_nav_or_sep(line: str) -> bool:
    s = line.strip()
    return (
        s.startswith("[[←")
        or s.startswith("**Navigation**")
        or s.startswith("Navigation :")
        or s.startswith("[← 前の章")
        or s == "---"
        or s.startswith("---")
        or s.startswith("----")
    )

def apply_clean_block(file_path: Path, defs: dict, dry_run: bool):
    try:
        content = file_path.read_text(encoding="utf-8")
    except Exception as e:
        print(f"❌ 読み込みエラー ({file_path.name}): {e}")
        return

    lines = content.splitlines()

    # 1. 最初のセクション見出し（## 1. または ## 2. または ## 4. など）を探す
    body_start_idx = -1
    for idx, line in enumerate(lines):
        if line.strip().startswith("## ") or line.strip().startswith("### 📦") or line.strip().startswith("## 概要"):
            body_start_idx = idx
            break

    if body_start_idx != -1:
        raw_body_lines = lines[body_start_idx:]
    else:
        # 見つからない場合は最初の区切り線より後
        dash_indices = [i for i, l in enumerate(lines) if l.strip().startswith("---") or l.strip().startswith("----")]
        raw_body_lines = lines[dash_indices[0] + 1:] if dash_indices else lines

    # 2. 末尾のフッター残骸を除去
    while raw_body_lines and (raw_body_lines[-1].strip() == "" or is_nav_or_sep(raw_body_lines[-1])):
        raw_body_lines.pop()

    # 3. 正しいヘッダーブロック + 本文 + フッターナビゲーション を完全結合
    new_full_content = (
        defs["header"].strip()
        + "\n\n---\n\n"
        + "\n".join(raw_body_lines).strip()
        + "\n\n---\n\n"
        + defs["footer"].strip()
        + "\n"
    )

    if not dry_run:
        file_path.write_text(new_full_content, encoding="utf-8")
    print(f"  ✨ [完全置換完了] {file_path.name}")

def main():
    dry_run = "--run" not in sys.argv
    print("=" * 70)
    print("  dozou_katanuki Part 3 & Part 4 ヘッダー/フッター完全置換")
    print(f"  モード: {'【 Dry-Run (事前確認) 】' if dry_run else '【 Real-Run (本番反映) 】'}")
    print("=" * 70)

    for filename, defs in TARGET_BLOCKS.items():
        file_path = DOCS_DIR / filename
        if file_path.exists():
            apply_clean_block(file_path, defs, dry_run)
        else:
            print(f"  ⚠️ ファイルが見つかりません: {filename}")

    print("\n" + "=" * 70)
    if dry_run:
        print("💡 Dry-Run完了。本番反映するには `--run` を付けて実行してください:")
        print("   > python fix_part3_part4_blocks.py --run")
    else:
        print("🎉 Part 3 & Part 4 の全ファイルのブロック完全置換が完了いたしました、先輩！")
    print("=" * 70)

if __name__ == "__main__":
    main()
