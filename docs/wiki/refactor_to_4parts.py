#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
dozou_katanuki 全4編確定体系一括リファクタリングスクリプト (Windows / Wiki直下対応版)
"""

import sys
import os
import re
import stat
import shutil
from pathlib import Path

# Wikiリポジトリのルート直下を作業ディレクトリに指定
DOCS_DIR = Path(".")

# 1. ファイルリネームマッピング (旧名 -> 新名)
FILE_RENAME_MAP = {
    # --- 第1編 ---
    "01_technical_specs_and_backbone.md": "part1_01_technical_specs.md",
    "02_external_services_and_salvage.md": "part1_02_external_services.md",
    "08_storage_persistence_and_media_policy_v2.md": "part1_03_storage_persistence.md",
    "08_storage_persistence_and_media_policy_v2_2.md": "part1_03_storage_persistence.md",
    "03_implementation_principles_and_constraints.md": "part1_04_implementation_principles.md",
    
    # --- 第2編 ---
    "10_wails_app_shell_and_governance.md": "part2_01_administration_and_recovery.md",
    "10_administration_and_recovery_v3.md": "part2_01_administration_and_recovery.md",
    "10_admin_board_and_governance_v5.md": "part2_01_administration_and_recovery.md",
    "09_plugin_architecture_and_sidecar_v3.md": "part2_02_plugin_architecture.md",
    
    # --- 第3編 ---
    "04_database_and_virtual_storage_pool.md": "part3_01_database_design.md",
    "04_database_and_virtual_storage_pool_2.md": "part3_01_database_design.md",
    "05_foolish_frontend_and_declarative_ui_v4.md": "part3_02_pure_dumb_frontend.md",
    "05_foolish_frontend_and_declarative_ui_v4_2.md": "part3_02_pure_dumb_frontend.md",
    "05_declarative_pure_dumb_frontend_v4.md": "part3_02_pure_dumb_frontend.md",
    "06_0_middleware_index.md": "part3_03_0_middleware_index.md",
    "06_thick_middleware_and_proxy_v4.md": "part3_03_0_middleware_index.md",
    "06_1_middleware_core.md": "part3_03_1_middleware_core.md",
    "06_2_data_decorator.md": "part3_03_2_data_decorator.md",
    "06_2_data_decorator_2.md": "part3_03_2_data_decorator.md",
    "06_3_plugin_orchestrator.md": "part3_03_3_job_orchestrator.md",
    "06_3_job_orchestrator.md": "part3_03_3_job_orchestrator.md",
    "07_robust_backend_driver_and_api_v4.md": "part3_04_backend_driver.md",
    "07_robust_backend_driver_and_api_v4_2.md": "part3_04_backend_driver.md",
    
    # --- 第4編 ---
    "11_references_and_literature.md": "part4_01_references_and_literature.md",
    "11_environment_and_troubleshooting.md": "part4_01_references_and_literature.md",
    "part4_01_environment_and_troubleshooting.md": "part4_01_references_and_literature.md",
}

# 2. 本文テキスト・リンク・見出し置換ルール
TEXT_REPLACEMENTS = [
    # --- リンク・ファイル参照の置換 ---
    (r"01_technical_specs_and_backbone", "part1_01_technical_specs"),
    (r"02_external_services_and_salvage", "part1_02_external_services"),
    (r"08_storage_persistence_and_media_policy_v2(_2)?", "part1_03_storage_persistence"),
    (r"03_implementation_principles_and_constraints", "part1_04_implementation_principles"),
    (r"10_wails_app_shell_and_governance", "part2_01_administration_and_recovery"),
    (r"10_administration_and_recovery_v3", "part2_01_administration_and_recovery"),
    (r"10_admin_board_and_governance_v5", "part2_01_administration_and_recovery"),
    (r"09_plugin_architecture_and_sidecar_v3", "part2_02_plugin_architecture"),
    (r"04_database_and_virtual_storage_pool(_2)?", "part3_01_database_design"),
    (r"05_foolish_frontend_and_declarative_ui_v4(_2)?", "part3_02_pure_dumb_frontend"),
    (r"05_declarative_pure_dumb_frontend_v4", "part3_02_pure_dumb_frontend"),
    (r"06_0_middleware_index", "part3_03_0_middleware_index"),
    (r"06_thick_middleware_and_proxy_v4", "part3_03_0_middleware_index"),
    (r"06_1_middleware_core", "part3_03_1_middleware_core"),
    (r"06_2_data_decorator(_2)?", "part3_03_2_data_decorator"),
    (r"06_3_plugin_orchestrator", "part3_03_3_job_orchestrator"),
    (r"06_3_job_orchestrator", "part3_03_3_job_orchestrator"),
    (r"07_robust_backend_driver_and_api_v4(_2)?", "part3_04_backend_driver"),
    (r"11_references_and_literature", "part4_01_references_and_literature"),
    (r"11_environment_and_troubleshooting", "part4_01_references_and_literature"),

    # --- 見出し（タイトル）の置換 ---
    (r"### 第1章：技術仕様とバックボーン", "### 第1編 第1章：技術仕様とバックボーン"),
    (r"### 第2章：外部サービスの概要とサルベージ技術", "### 第1編 第2章：外部サービスの概要とサルベージ技術"),
    (r"# 第8章：ローカルストレージ保全とメディアポリシー", "# 第1編 第3章：ローカルストレージ保全とメディアポリシー"),
    (r"### 第3章：実装規約・制約原則", "### 第1編 第4章：実装規約・制約原則"),
    (r"# 第2編 第1章：.*", "# 第2編 第1章：管理・設定・ディザスタリカバリ運用"),
    (r"# 第10章：.*", "# 第2編 第1章：管理・設定・ディザスタリカバリ運用"),
    (r"# 第9章：プラグインアーキテクチャ.*", "# 第2編 第2章：プラグインアーキテクチャとサイドカー"),
    (r"### 第4章：データベース設計と仮想ストレージプール", "### 第3編 第1章：データベース設計と仮想ストレージプール"),
    (r"# 第5章：プレゼンテーション層概論", "# 第3編 第2章：フロントエンド層概論"),
    (r"# 第6章：ミドルウェア層.*インデックス", "# 第3編 第3章：ミドルウェア層インデックス"),
    (r"# 6\.1 Middleware Core Components", "# 第3編 第3章.1：Middleware Core Components"),
    (r"# 6\.2 Data & Skin Decorator", "# 第3編 第3章.2：Data & Skin Decorator"),
    (r"# 6\.3 Job & Process Orchestrator", "# 第3編 第3章.3：Job & Process Orchestrator"),
    (r"# 第7章：ドライバー層", "# 第3編 第4章：ドライバー層"),
    (r"# 第4編 第1章：.*", "# 第4編 第1章：参考資料・技術文献・型定義カタログ・公式リンク集"),
    (r"# 第11章：.*", "# 第4編 第1章：参考資料・技術文献・型定義カタログ・公式リンク集"),
]

def remove_readonly(func, path, excinfo):
    """Windowsの読み取り専用属性を強制解除して削除を再試行するハンドラ"""
    try:
        os.chmod(path, stat.S_IWRITE)
        func(path)
    except Exception as e:
        print(f"⚠️ 強制削除スキップ: {path} ({e})")

def create_safe_backup(target_dir: Path, backup_dir: Path):
    """ .git や一時フォルダを除外して .md ファイルのみを安全に退避 """
    if backup_dir.exists():
        shutil.rmtree(backup_dir, onexc=remove_readonly)
    backup_dir.mkdir(parents=True, exist_ok=True)
    
    count = 0
    for md_file in target_dir.glob("*.md"):
        shutil.copy2(md_file, backup_dir / md_file.name)
        count += 1
    print(f"📦 [Backup] Markdown {count} 件を安全にバックアップしました: {backup_dir.resolve()}\n")

def main():
    dry_run = "--run" not in sys.argv
    print("=" * 70)
    print("  dozou_katanuki 【全4編完全版】一括リファクタリング (Windows Safe)")
    print(f"  モード: {'【 Dry-Run (事前プレビュー) 】' if dry_run else '【 Real-Run (本番反映) 】'}")
    print("=" * 70)

    # 1. バックアップ作成 (本番時のみ、.md ファイルのみ対象)
    backup_dir = DOCS_DIR / "_backup_md_safe"
    if not dry_run:
        create_safe_backup(DOCS_DIR, backup_dir)

    # 2. Markdown 本文の置換
    print("--- 1. 本文テキスト・リンク置換 ---")
    md_files = [f for f in DOCS_DIR.glob("*.md") if not f.name.startswith("_")]
    for file_path in md_files:
        try:
            content = file_path.read_text(encoding="utf-8")
        except Exception as e:
            print(f"❌ 読み込み失敗 ({file_path.name}): {e}")
            continue

        new_content = content
        count = 0
        for pattern, repl in TEXT_REPLACEMENTS:
            m = len(re.findall(pattern, new_content))
            if m > 0:
                count += m
                new_content = re.sub(pattern, repl, new_content)
        if count > 0:
            if not dry_run:
                file_path.write_text(new_content, encoding="utf-8")
            print(f"  📝 {file_path.name}: {count} 箇所更新")

    # 3. ファイル名のリネーム
    print("\n--- 2. ファイル名リネーム ---")
    for old_name, new_name in FILE_RENAME_MAP.items():
        if old_name == new_name:
            continue
        old_path = DOCS_DIR / old_name
        new_path = DOCS_DIR / new_name

        if old_path.exists():
            if new_path.exists():
                print(f"  ⚠️ [スキップ] 既にリネーム先が存在します: {old_name} ➔ {new_name}")
                continue
            if dry_run:
                print(f"  🔄 [予定] {old_name} ➔ {new_name}")
            else:
                old_path.rename(new_path)
                print(f"  ✅ [完了] {old_name} ➔ {new_name}")

    print("\n" + "=" * 70)
    if dry_run:
        print("💡 Dry-Run完了。本番反映するには引数に `--run` を付けてください:")
        print("   > python refactor_to_4parts.py --run")
    else:
        print("🎉 すべての置換とリネームが正常に完了いたしました、先輩！")
    print("=" * 70)

if __name__ == "__main__":
    main()
