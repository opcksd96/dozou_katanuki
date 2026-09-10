# scratch/repair_article_translations.py
import argparse
import json
import os
import shutil
import sqlite3
import sys
import time
from collections import Counter
from typing import Dict, Optional, Tuple

# プロジェクトルートを sys.path に追加
PROJECT_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
if PROJECT_ROOT not in sys.path:
    sys.path.insert(0, PROJECT_ROOT)

from plugins.base.scraper.core.translator import Translator


def load_config_api_keys():
    config_path = os.path.join(PROJECT_ROOT, "config.json")
    if os.path.exists(config_path):
        try:
            with open(config_path, "r", encoding="utf-8") as f:
                data = json.load(f)
                trans = data.get("translation", {})
                deepl = trans.get("deepl_api_key", "")
                google = trans.get("google_translate_api_key", "")
                if deepl and not os.environ.get("DEEPL_API_KEY"):
                    os.environ["DEEPL_API_KEY"] = deepl
                if google and not os.environ.get("GOOGLE_TRANSLATE_API_KEY"):
                    os.environ["GOOGLE_TRANSLATE_API_KEY"] = google
        except Exception:
            pass


def backup_database(db_path: str) -> str:
    backup_dir = os.path.join(PROJECT_ROOT, "backups", "database")
    os.makedirs(backup_dir, exist_ok=True)
    ts = time.strftime("%Y%m%d_%H%M%S")
    backup_path = os.path.join(backup_dir, f"archive_pre_repair_translations_{ts}.db")
    shutil.copy2(db_path, backup_path)
    print(f"[BACKUP] Created backup at: {backup_path}")
    return backup_path


def run_repair(db_path: str = "archive.db", dry_run: bool = True, skip_backup: bool = False,
               deepl_key: str = "", google_key: str = "", limit: int = 0):
    print("=" * 60)
    print(f"=== REPAIR ARTICLE TRANSLATIONS & PURGE INVALID FULL_TEXT_JA ===")
    print(f"Target DB: {db_path} | DryRun: {dry_run} | Limit: {limit or 'ALL'}")
    print("=" * 60)

    if not os.path.exists(db_path):
        print(f"[ERROR] Database file not found: {db_path}")
        return

    # キー設定の優先順位: 引数 > 環境変数 > config.json
    if deepl_key:
        os.environ["DEEPL_API_KEY"] = deepl_key
    if google_key:
        os.environ["GOOGLE_TRANSLATE_API_KEY"] = google_key
    load_config_api_keys()

    translator = Translator()
    has_api_key = translator.provider != "none"
    print(f"[TRANSLATOR] Provider: {translator.provider} (API Available: {has_api_key})")
    if not has_api_key:
        print("[WARNING] No DeepL or Google Translate API Key detected!")
        print("          Foreign texts (ZH/EN) will have invalid JA texts cleared and original texts saved,")
        print("          but automatic external translation will be skipped until an API key is provided.")

    if not dry_run and not skip_backup:
        backup_database(db_path)

    conn = sqlite3.connect(db_path, timeout=60.0)
    conn.execute("PRAGMA journal_mode = WAL;")
    cur = conn.cursor()

    query = "SELECT id, full_text, lang, full_text_ja, full_text_en, full_text_zh FROM articles"
    if limit > 0:
        query += f" LIMIT {limit}"
    rows = cur.execute(query).fetchall()

    total = len(rows)
    print(f"[SCAN] Processing {total} articles from database...")

    stats_old_lang = Counter()
    stats_new_lang = Counter()
    purged_ja_count = 0
    saved_orig_zh_count = 0
    saved_orig_en_count = 0
    saved_orig_ja_count = 0
    translated_count = 0
    updated_rows = []

    for idx, (aid, ftext, old_lang, old_ja, old_en, old_zh) in enumerate(rows, start=1):
        stats_old_lang[old_lang or "UNKNOWN"] += 1
        if not ftext or not ftext.strip():
            continue

        new_lang = translator.detect_lang(ftext)
        stats_new_lang[new_lang] += 1

        final_lang = new_lang
        final_ja = old_ja
        final_en = old_en
        final_zh = old_zh

        if new_lang == "ja":
            final_lang = "ja"
            final_ja = ftext
            saved_orig_ja_count += 1
            if final_en == ftext: final_en = None
            if final_zh == ftext: final_zh = None

        elif new_lang == "zh":
            final_lang = "zh"
            final_zh = ftext
            saved_orig_zh_count += 1
            if final_en == ftext: final_en = None

            if has_api_key:
                need_ja = not final_ja or final_ja == ftext
                need_en = not final_en or final_en == ftext
                targets = []
                if need_ja: targets.append("ja")
                if need_en: targets.append("en")
                if targets:
                    res = translator.translate_article(ftext, source_lang="zh", targets=targets)
                    if res.get("ja"):
                        final_ja = res["ja"]
                        translated_count += 1
                    if res.get("en"):
                        final_en = res["en"]
            else:
                if old_ja is not None:
                    purged_ja_count += 1
                final_ja = None

        elif new_lang == "en":
            final_lang = "en"
            final_en = ftext
            saved_orig_en_count += 1
            if final_zh == ftext: final_zh = None

            if has_api_key:
                need_ja = not final_ja or final_ja == ftext
                need_zh = not final_zh or final_zh == ftext
                targets = []
                if need_ja: targets.append("ja")
                if need_zh: targets.append("zh")
                if targets:
                    res = translator.translate_article(ftext, source_lang="en", targets=targets)
                    if res.get("ja"):
                        final_ja = res["ja"]
                        translated_count += 1
                    if res.get("zh"):
                        final_zh = res["zh"]
            else:
                if old_ja is not None:
                    purged_ja_count += 1
                final_ja = None

        else:
            final_lang = new_lang
            if final_en == ftext: final_en = None
            if final_zh == ftext: final_zh = None
            if has_api_key:
                targets = translator.determine_targets(new_lang)
                res = translator.translate_article(ftext, source_lang=new_lang, targets=targets)
                if res.get("ja"):
                    final_ja = res["ja"]
                    translated_count += 1
            else:
                if old_ja is not None:
                    purged_ja_count += 1
                final_ja = None

        if (final_lang != old_lang or final_ja != old_ja or
                final_en != old_en or final_zh != old_zh):
            updated_rows.append((final_lang, final_ja, final_en, final_zh, aid))

        if idx % 500 == 0 or idx == total:
            print(f"[PROGRESS] [{idx}/{total}] Processed. Pending updates: {len(updated_rows)}")

    print("-" * 60)
    print("=== SUMMARY STATS ===")
    print(f"Old Lang Distribution: {dict(stats_old_lang)}")
    print(f"New Lang Distribution: {dict(stats_new_lang)}")
    print(f"Japanese (ja) orig saved to full_text_ja: {saved_orig_ja_count}")
    print(f"Chinese (zh) orig saved to full_text_zh  : {saved_orig_zh_count}")
    print(f"English (en) orig saved to full_text_en  : {saved_orig_en_count}")
    print(f"Invalid full_text_ja purged (cleared)   : {purged_ja_count}")
    if has_api_key:
        print(f"Translations performed                  : {translated_count}")
    print(f"Articles to update                      : {len(updated_rows)} / {total}")
    print("-" * 60)

    if dry_run:
        print("[DRY-RUN] No changes were written to database. Run with --apply to commit changes.")
    else:
        print(f"[APPLY] Committing {len(updated_rows)} updates to database...")
        cur.executemany("""
            UPDATE articles
            SET lang = ?, full_text_ja = ?, full_text_en = ?, full_text_zh = ?
            WHERE id = ?
        """, updated_rows)
        conn.commit()
        print("[APPLY] Successfully updated database!")

    conn.close()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Repair article language tags and translation fields in ArchiveDB")
    parser.add_argument("--db-path", default="archive.db", help="Path to archive.db")
    parser.add_argument("--apply", action="store_true", help="Apply updates to database (default is dry-run)")
    parser.add_argument("--skip-backup", action="store_true", help="Skip automatic DB backup before apply")
    parser.add_argument("--deepl-key", default="", help="DeepL API Key (free or pro)")
    parser.add_argument("--google-key", default="", help="Google Translate API Key")
    parser.add_argument("--limit", type=int, default=0, help="Limit number of articles to process (for test)")
    args = parser.parse_args()

    run_repair(
        db_path=args.db_path,
        dry_run=not args.apply,
        skip_backup=args.skip_backup,
        deepl_key=args.deepl_key,
        google_key=args.google_key,
        limit=args.limit
    )
