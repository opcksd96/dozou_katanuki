#!/usr/bin/env python3
"""
combine_wiki.py (100行以下)
土蔵・型抜き Wiki の全Markdownドキュメントを順序付けして1つの巨大Markdownファイルに統合するスクリプト
"""
import os, sys, glob
from datetime import datetime

SCRIPT_DIR = os.path.abspath(os.path.dirname(__file__))
PROJECT_ROOT = os.path.abspath(os.path.join(SCRIPT_DIR, "..", ".."))
WIKI_DIR = os.path.join(PROJECT_ROOT, "docs", "wiki")
OUTPUT_FILE = os.path.join(PROJECT_ROOT, "docs", "COMBINED_WIKI.md")


# 優先順位定義（この順序で連結し、残りはアルファベット順）
ORDER_PREFIXES = [
    "Home.md",
    "Status_and_Roadmap.md",
    "SPEC_TASK3_TASK2_WHITELIST_EXTERNAL_AND_SCRAPERS.md",
    "PLAN_TASK3_TASK2_WHITELIST_EXTERNAL_AND_SCRAPERS.md",
    "part1_",
    "part2_01_",
    "part2_02_",
    "part3_01_",
    "part3_02_",
    "part3_03_",
    "part3_04_",
    "part4_",
]

def sort_key(filename: str) -> tuple:
    base = os.path.basename(filename)
    for idx, prefix in enumerate(ORDER_PREFIXES):
        if base == prefix or base.startswith(prefix):
            return (idx, base)
    return (len(ORDER_PREFIXES), base)

def main():
    md_files = [
        f for f in glob.glob(os.path.join(WIKI_DIR, "*.md"))
        if os.path.basename(f) != "COMBINED_WIKI.md"
    ]
    md_files.sort(key=sort_key)

    print(f"[*] Found {len(md_files)} markdown files in Wiki.")

    combined_lines = [
        "# 🏯 土蔵・型抜き (dozou_katanuki) 完全統合ドキュメント仕様書 (COMBINED WIKI)\n\n",
        f"> **生成日時**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}  \n",
        f"> **対象ファイル数**: {len(md_files)} ファイル  \n",
        f"> **リポジトリ**: `opcksd96/dozou_katanuki.wiki`\n\n",
        "---\n\n",
        "## 📑 総合目次 (Table of Contents)\n\n",
    ]

    for idx, fpath in enumerate(md_files, 1):
        fname = os.path.basename(fpath)
        combined_lines.append(f"{idx}. [{fname}](#{fname.lower().replace('.', '').replace('_', '-')})\n")

    combined_lines.append("\n---\n\n")

    for idx, fpath in enumerate(md_files, 1):
        fname = os.path.basename(fpath)
        print(f"  [{idx}/{len(md_files)}] Combining: {fname}")
        with open(fpath, "r", encoding="utf-8") as f:
            content = f.read().strip()

        combined_lines.extend([
            f"\n\n<!-- ================================================================= -->\n",
            f"<!-- SECTION {idx}: {fname} -->\n",
            f"<!-- ================================================================= -->\n\n",
            f'<div id="{fname.lower().replace(".", "").replace("_", "-")}"></div>\n\n',
            f"# 📄 [{idx}] {fname}\n\n",
            f"> *Source File: `{fname}`*\n\n",
            content,
            "\n\n---\n\n",
        ])

    out_path = sys.argv[1] if len(sys.argv) > 1 else OUTPUT_FILE
    with open(out_path, "w", encoding="utf-8") as f:
        f.writelines(combined_lines)

    size_kb = os.path.getsize(out_path) / 1024
    print(f"\n[+] Successfully generated combined markdown: {out_path} ({size_kb:.1f} KB)")

if __name__ == "__main__":
    main()
