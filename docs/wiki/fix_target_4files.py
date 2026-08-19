#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
dozou_katanuki 指定4ファイル限定・ヘッダー/フッターブロック完全置換スクリプト
- 対象: part1_04, part2_01, part2_02, part3_01
- 先頭のタイトル/メタデータ/ナビゲーション領域を完全なブロックとして丸ごと差し替え
- 末尾のフッターナビゲーションも完全一致で同期
"""

import sys
import re
from pathlib import Path

DOCS_DIR = Path(".")

# 4ファイルの「完全ヘッダーブロック」および「フッターナビゲーション」定義
TARGET_BLOCKS = {
    # 1. 第1編 第4章
    "part1_04_implementation_principles.md": {
        "header": """[[← 前の章: 第1編第3章：ローカルストレージ保全とメディアポリシー|part1_03_storage_persistence]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第1章：管理・設定・ディザスタリカバリ運用 →|part2_01_administration_and_recovery]]

# 第1編 第4章：実装規約・制約原則 (Implementation Principles & Constraints)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-PRINCIPLE-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（宣言型UI・UDF・Wails Bind完全準拠、AI暴走防止規約統合）""",
        "footer": "[[← 前の章: 第1編第3章：ローカルストレージ保全とメディアポリシー|part1_03_storage_persistence]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第1章：管理・設定・ディザスタリカバリ運用 →|part2_01_administration_and_recovery]]"
    },

    # 2. 第2編 第1章
    "part2_01_administration_and_recovery.md": {
        "header": """[[← 前の章: 第1編第4章：実装規約・制約原則|part1_04_implementation_principles]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第2章：プラグインアーキテクチャとサイドカー →|part2_02_plugin_architecture]]

# 第2編 第1章：管理・設定・ディザスタリカバリ運用 (Administration & Governance)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-ADMIN-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（テキストグラフ全パージ・Mermaid完全移行・一元設定・7大管理ビュー・Scraper View・二重化バックアップ）""",
        "footer": "[[← 前の章: 第1編第4章：実装規約・制約原則|part1_04_implementation_principles]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第2章：プラグインアーキテクチャとサイドカー →|part2_02_plugin_architecture]]"
    },

    # 3. 第2編 第2章
    "part2_02_plugin_architecture.md": {
        "header": """[[← 前の章: 第2編第1章：管理・設定・ディザスタリカバリ運用|part2_01_administration_and_recovery]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第1章：データベース設計と仮想ストレージプール →|part3_01_database_design]]

# 第2編 第2章：プラグインアーキテクチャとサイドカー (Plugin Architecture & Sidecar)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-PLUGIN-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（統合プラグイン plugins/ 規格・Go製レンダラー・Python 3Arrows・3段階メディア確保ライフサイクル・ミドルウェア非同期制御・ポート9998同期）""",
        "footer": "[[← 前の章: 第2編第1章：管理・設定・ディザスタリカバリ運用|part2_01_administration_and_recovery]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第1章：データベース設計と仮想ストレージプール →|part3_01_database_design]]"
    },

    # 4. 第3編 第1章
    "part3_01_database_design.md": {
        "header": """[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]

# 第3編 第1章：データベース設計と仮想ストレージプール (Database Design & Virtual Storage Pool)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-DATABASE-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（GORMモデル完全同期・SQLite3 DDL・10大インデックス最適化・WALモード競合制御定義）""",
        "footer": "[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]"
    }
}

def is_nav_or_sep(line: str) -> bool:
    s = line.strip()
    return s.startswith("[[←") or s.startswith("**Navigation**") or s.startswith("Navigation :") or s.startswith("[← 前の章") or s == "---" or s.startswith("---") or s.startswith("----")

def apply_clean_block(file_path: Path, defs: dict, dry_run: bool):
    try:
        content = file_path.read_text(encoding="utf-8")
    except Exception as e:
        print(f"❌ 読み込みエラー ({file_path.name}): {e}")
        return

    lines = content.splitlines()

    # 1. 最初のセクション見出し（## X.Y または #### X.Y）の行を探す
    body_start_idx = -1
    for idx, line in enumerate(lines):
        if line.strip().startswith("## ") or line.strip().startswith("#### "):
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
    print("  dozou_katanuki 指定4ファイル ヘッダー/フッター完全置換")
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
        print("   > python fix_target_4files.py --run")
    else:
        print("🎉 指定4ファイルのブロック完全置換が完了いたしました、先輩！")
    print("=" * 70)

if __name__ == "__main__":
    main()
