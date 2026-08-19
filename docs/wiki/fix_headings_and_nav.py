#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
dozou_katanuki 見出しレベル・節番号・ナビゲーション完全修復スクリプト
"""

import sys
import re
from pathlib import Path

DOCS_DIR = Path(".")

# 1. 各ファイルごとの正規化ナビゲーション定義 (完全な数珠つなぎ)
NAV_DEFINITIONS = {
    "part1_01_technical_specs.md": {
        "title": "# 第1編 第1章：技術仕様とバックボーン (Technical Specs & Backbone)",
        "header_nav": "[[← DocWiki ポータル|Home]] | [[📚 目次 (Home)|Home]] | [[次の章: 第1編第2章：外部サービスの概要とサルベージ技術 →|part1_02_external_services]]",
        "footer_nav": "[[← DocWiki ポータル|Home]] | [[📚 目次 (Home)|Home]] | [[次の章: 第1編第2章：外部サービスの概要とサルベージ技術 →|part1_02_external_services]]",
        "sec_prefix": "1."
    },
    "part1_02_external_services.md": {
        "title": "# 第1編 第2章：外部サービスの概要とサルベージ技術 (External Services & Salvage Technologies)",
        "header_nav": "[[← 前の章: 第1編第1章：技術仕様とバックボーン|part1_01_technical_specs]] | [[📚 目次 (Home)|Home]] | [[次の章: 第1編第3章：ローカルストレージ保全とメディアポリシー →|part1_03_storage_persistence]]",
        "footer_nav": "[[← 前の章: 第1編第1章：技術仕様とバックボーン|part1_01_technical_specs]] | [[📚 目次 (Home)|Home]] | [[次の章: 第1編第3章：ローカルストレージ保全とメディアポリシー →|part1_03_storage_persistence]]",
        "sec_prefix": "2."
    },
    "part1_03_storage_persistence.md": {
        "title": "# 第1編 第3章：ローカルストレージ保全とメディアポリシー (Storage Persistence & Media Policy)",
        "header_nav": "[[← 前の章: 第1編第2章：外部サービスの概要とサルベージ技術|part1_02_external_services]] | [[📚 目次 (Home)|Home]] | [[次の章: 第1編第4章：実装規約・制約原則 →|part1_04_implementation_principles]]",
        "footer_nav": "[[← 前の章: 第1編第2章：外部サービスの概要とサルベージ技術|part1_02_external_services]] | [[📚 目次 (Home)|Home]] | [[次の章: 第1編第4章：実装規約・制約原則 →|part1_04_implementation_principles]]",
        "sec_old": "8.",
        "sec_new": "3."
    },
    "part1_04_implementation_principles.md": {
        "title": "# 第1編 第4章：実装規約・制約原則 (Implementation Principles & Constraints)",
        "header_nav": "[[← 前の章: 第1編第3章：ローカルストレージ保全とメディアポリシー|part1_03_storage_persistence]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第1章：管理・設定・ディザスタリカバリ運用 →|part2_01_administration_and_recovery]]",
        "footer_nav": "[[← 前の章: 第1編第3章：ローカルストレージ保全とメディアポリシー|part1_03_storage_persistence]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第1章：管理・設定・ディザスタリカバリ運用 →|part2_01_administration_and_recovery]]",
        "sec_old": "3.",
        "sec_new": "4."
    },
    "part2_01_administration_and_recovery.md": {
        "title": "# 第2編 第1章：管理・設定・ディザスタリカバリ運用 (Administration & Governance)",
        "header_nav": "[[← 前の章: 第1編第4章：実装規約・制約原則|part1_04_implementation_principles]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第2章：プラグインアーキテクチャとサイドカー →|part2_02_plugin_architecture]]",
        "footer_nav": "[[← 前の章: 第1編第4章：実装規約・制約原則|part1_04_implementation_principles]] | [[📚 目次 (Home)|Home]] | [[次の章: 第2編第2章：プラグインアーキテクチャとサイドカー →|part2_02_plugin_architecture]]",
        "sec_old": "10.",
        "sec_new": "1."
    },
    "part2_02_plugin_architecture.md": {
        "title": "# 第2編 第2章：プラグインアーキテクチャとサイドカー (Plugin Architecture & Sidecar)",
        "header_nav": "[[← 前の章: 第2編第1章：管理・設定・ディザスタリカバリ運用|part2_01_administration_and_recovery]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第1章：データベース設計と仮想ストレージプール →|part3_01_database_design]]",
        "footer_nav": "[[← 前の章: 第2編第1章：管理・設定・ディザスタリカバリ運用|part2_01_administration_and_recovery]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第1章：データベース設計と仮想ストレージプール →|part3_01_database_design]]",
        "sec_old": "9.",
        "sec_new": "2."
    },
    "part3_01_database_design.md": {
        "title": "# 第3編 第1章：データベース設計と仮想ストレージプール (Database Design & Virtual Storage Pool)",
        "header_nav": "[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]",
        "footer_nav": "[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]",
        "sec_old": "4.",
        "sec_new": "1."
    },
    "part3_02_pure_dumb_frontend.md": {
        "title": "# 第3編 第2章：フロントエンド層概論（Pure Dumb UI Framework）",
        "header_nav": "[[← 前の章: 第3編第1章：データベース設計と仮想ストレージプール|part3_01_database_design]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第3章：ミドルウェア層インデックス →|part3_03_0_middleware_index]]",
        "footer_nav": "[[← 前の章: 第3編第1章：データベース設計と仮想ストレージプール|part3_01_database_design]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第3章：ミドルウェア層インデックス →|part3_03_0_middleware_index]]",
        "sec_old": "5.",
        "sec_new": "2."
    },
    "part3_03_0_middleware_index.md": {
        "title": "# 第3編 第3章：ミドルウェア層インデックス (Intelligent Hub Architecture)",
        "header_nav": "[[← 前の章: 第3編第2章：フロントエンド層概論|part3_02_pure_dumb_frontend]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第4章：ドライバー層 →|part3_04_backend_driver]]",
        "footer_nav": "[[← 前の章: 第3編第2章：フロントエンド層概論|part3_02_pure_dumb_frontend]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第4章：ドライバー層 →|part3_04_backend_driver]]",
    },
    "part3_04_backend_driver.md": {
        "title": "# 第3編 第4章：ドライバー層（Core Backend Driver & Data Abstraction）",
        "header_nav": "[[← 前の章: 第3編第3章：ミドルウェア層インデックス|part3_03_0_middleware_index]] | [[📚 目次 (Home)|Home]] | [[次の章: 第4編第1章：参考資料・技術文献・型定義カタログ・公式リンク集 →|part4_01_references_and_literature]]",
        "footer_nav": "[[← 前の章: 第3編第3章：ミドルウェア層インデックス|part3_03_0_middleware_index]] | [[📚 目次 (Home)|Home]] | [[次の章: 第4編第1章：参考資料・技術文献・型定義カタログ・公式リンク集 →|part4_01_references_and_literature]]",
        "sec_old": "7.",
        "sec_new": "4."
    },
    "part4_01_references_and_literature.md": {
        "title": "# 第4編 第1章：参考資料・技術文献・型定義カタログ・公式リンク集",
        "header_nav": "[[← 前の章: 第3編第4章：ドライバー層|part3_04_backend_driver]] | [[📚 目次 (Home)|Home]] | [[DocWiki ポータルへ戻る →|Home]]",
        "footer_nav": "[[← 前の章: 第3編第4章：ドライバー層|part3_04_backend_driver]] | [[📚 目次 (Home)|Home]] | [[DocWiki ポータルへ戻る →|Home]]",
        "sec_old": "11.",
        "sec_new": "1."
    }
}

def clean_and_fix_file(file_path: Path, defs: dict, dry_run: bool):
    content = file_path.read_text(encoding="utf-8")
    lines = content.splitlines()

    # 1. 見出しレベル（#### -> ##, ### -> ##）の正規化
    new_lines = []
    for line in lines:
        # #### X.Y または ## X.Y を ## X.Y に統一
        line = re.sub(r'^(?:####|##)\s+(\d+\.\d+)', r'## \1', line)
        # #### X.Y.Z または ### X.Y.Z を ### X.Y.Z に統一
        line = re.sub(r'^(?:####|###)\s+(\d+\.\d+\.\d+)', r'### \1', line)
        new_lines.append(line)
    
    content = "\n".join(new_lines)

    # 2. 節番号の置換 (例: 8.1 -> 3.1)
    if "sec_old" in defs and "sec_new" in defs:
        old_p = defs["sec_old"]
        new_p = defs["sec_new"]
        content = re.sub(rf'## {re.escape(old_p)}(\d+)', rf'## {new_p}\1', content)
        content = re.sub(rf'### {re.escape(old_p)}(\d+)', rf'### {new_p}\1', content)

    # 3. タイトル行 (H1) の正規化
    content = re.sub(r'^(?:#{1,3})\s+(?:第\d+編\s+)?第\d+章.*$', defs["title"], content, flags=re.MULTILINE)

    # 4. ヘッダーナビゲーションの置換
    # [[← ... ]] で始まる1行目またはNavigation行
    lines = content.splitlines()
    if lines:
        if lines[0].startswith("[[←") or lines[0].startswith("**Navigation**"):
            lines[0] = defs["header_nav"]
        elif len(lines) > 1 and (lines[1].startswith("[[←") or lines[1].startswith("**Navigation**")):
            lines[1] = defs["header_nav"]

        # フッターナビゲーションの置換（末尾のナビゲーション）
        for i in range(len(lines) - 1, -1, -1):
            if lines[i].startswith("[[←") or lines[i].startswith("**Navigation**"):
                lines[i] = defs["footer_nav"]
                break

    new_content = "\n".join(lines) + "\n"

    if not dry_run:
        file_path.write_text(new_content, encoding="utf-8")
    print(f"  ✨ 修復完了: {file_path.name}")

def main():
    dry_run = "--run" not in sys.argv
    print("=" * 70)
    print("  dozou_katanuki 見出しレベル・節番号・ナビゲーション完全修復")
    print(f"  モード: {'【 Dry-Run (事前プレビュー) 】' if dry_run else '【 Real-Run (本番反映) 】'}")
    print("=" * 70)

    for filename, defs in NAV_DEFINITIONS.items():
        file_path = DOCS_DIR / filename
        if file_path.exists():
            clean_and_fix_file(file_path, defs, dry_run)
        else:
            print(f"  ⚠️ ファイルが見つかりません: {filename}")

    print("\n" + "=" * 70)
    if dry_run:
        print("💡 Dry-Run完了。本番反映するには引数に `--run` を付けてください:")
        print("   > python fix_headings_and_nav.py --run")
    else:
        print("🎉 すべてのファイルの見出し・節番号・ナビゲーションが完全に修復されました、先輩！")
    print("=" * 70)

if __name__ == "__main__":
    main()
