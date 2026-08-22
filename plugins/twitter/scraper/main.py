# plugins/twitter/scraper/main.py (100行以下)
import argparse, json, os, sqlite3, sys, time
CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
if CURRENT_DIR not in sys.path: sys.path.insert(0, CURRENT_DIR)

from core.downloader import Downloader
from core.mutator import Mutator
from core.restorer import Restorer
from core.scraper import Scraper
from core.translator import Translator
from core.warc_importer import WarcImporter
from parsers.twitter_parser import TwitterParser


def emit_progress(current: int, total: int, message: str) -> None:
    print(f"PROGRESS: {current}/{total} | {message}", flush=True)


def run_batch_translate(db_path: str, article_id: str = "", account: str = "", overwrite: bool = False, limit: int = 500, dry_run: bool = False) -> None:
    trans = Translator()
    with sqlite3.connect(db_path) as conn:
        cur = conn.cursor()
        if article_id:
            rows = cur.execute("SELECT id, full_text FROM articles WHERE id = ?", (article_id,)).fetchall()
        else:
            q, params = "SELECT id, full_text FROM articles WHERE 1=1", []
            if account: q += " AND account_id = (SELECT numeric_id FROM accounts WHERE username = ? OR numeric_id = ?)"; params.extend([account, account])
            if not overwrite: q += " AND (full_text_ja IS NULL OR full_text_ja = '' OR (full_text_ja = full_text AND full_text_en = full_text AND lang != 'ja'))"
            q += " LIMIT ?"; params.append(limit)
            rows = cur.execute(q, params).fetchall()

        if dry_run and article_id and rows: return print(f"JSON:{json.dumps(trans.translate_article(rows[0][1]), ensure_ascii=False)}")
        total = len(rows); emit_progress(0, max(total, 1), f"Starting translation for {total} articles...")
        for idx, (aid, ftext) in enumerate(rows, start=1):
            t_res = trans.translate_article(ftext)
            cur.execute("UPDATE articles SET lang=?, full_text_ja=?, full_text_en=?, full_text_zh=? WHERE id=?",
                        (t_res["lang"], t_res["ja"], t_res["en"], t_res["zh"], aid))
            conn.commit(); emit_progress(idx, total, f"Translated [{idx}/{total}] article {aid}")
        emit_progress(total, total, f"Completed translation for {total} articles.")


def run_auto_salvage(platform: str, account: str, limit: int, db_path: str, storage_dir: str, enable_trans: bool) -> None:
    emit_progress(0, 100 if limit <= 0 else limit, f"Starting auto salvage for @{account} on {platform}...")
    if account:
        try:
            with sqlite3.connect(db_path) as conn:
                r = conn.cursor().execute("SELECT id FROM whitelists WHERE lower(value) = lower(?)", (account,)).fetchone()
                if r: conn.execute("UPDATE whitelists SET is_active = 1 WHERE id = ?", (r[0],))
                else: conn.execute("INSERT INTO whitelists (type, value, is_active) VALUES ('account', ?, 1)", (account,))
                conn.commit()
        except Exception: pass
    scraper, parser, mutator, dl = Scraper(platform=platform), TwitterParser(), Mutator(db_path=db_path, enable_translation=enable_trans), Downloader(db_path=db_path, storage_dir=storage_dir)
    snapshots = scraper.search_cdx(account=account, limit=limit)
    if not snapshots: return emit_progress(1, 1, f"No snapshots for @{account}.")
    tot, c = len(snapshots), 0
    emit_progress(0, tot, f"Found {tot} snapshots. Starting...")
    for idx, snap in enumerate(snapshots, start=1):
        raw = scraper.fetch_snapshot(snap.get("timestamp", ""), snap.get("original", ""), account=account)
        if raw:
            p = parser.parse_record(raw, snap.get("original", ""))
            if p and mutator.upsert_record(p):
                dl.process_queued_media(p["post"]["id"]); c += 1
        emit_progress(idx, tot, f"Processed [{idx}/{tot}] snapshots...")
        time.sleep(0.2)
    emit_progress(tot, tot, f"Completed: Processed {c} posts.")


def main():
    p = argparse.ArgumentParser(description="Twitter Scraper Sidecar")
    p.add_argument("-m", "--mode", choices=["auto", "manual", "download", "poll", "restore", "translate"], default="auto")
    p.add_argument("-p", "--platform", default="twitter")
    p.add_argument("-a", "--account", default=""); p.add_argument("-l", "--limit", type=int, default=0)
    p.add_argument("-w", "--warc-path", default=""); p.add_argument("--dumps-dir", default="backups/dumps")
    p.add_argument("--avatar-dir", default="assets/avatars"); p.add_argument("--media-id", default="")
    p.add_argument("--article-id", default=""); p.add_argument("--offline", action="store_true")
    p.add_argument("--no-translate", action="store_true"); p.add_argument("--overwrite", action="store_true")
    p.add_argument("--dry-run", action="store_true"); p.add_argument("--db-path", default="archive.db")
    p.add_argument("--storage-dir", default="")
    args = p.parse_args()

    if args.mode == "translate":
        run_batch_translate(args.db_path, article_id=args.article_id, account=args.account, overwrite=args.overwrite, limit=args.limit, dry_run=args.dry_run)
    elif args.mode == "auto":
        run_auto_salvage(args.platform, args.account, args.limit, args.db_path, args.storage_dir, not args.no_translate and not args.offline)
    elif args.mode == "manual":
        WarcImporter(args.warc_path, db_path=args.db_path, storage_dir=args.storage_dir, offline=args.offline).run_import(emit_progress)
    elif args.mode == "download":
        emit_progress(1, 1, f"Processed: {Downloader(db_path=args.db_path, storage_dir=args.storage_dir).process_queued_media(args.article_id or None, args.media_id or None)}")
    elif args.mode == "poll":
        emit_progress(1, 1, f"Salvaged: {Downloader(db_path=args.db_path, storage_dir=args.storage_dir).poll_outsourced_media()}")
    elif args.mode == "restore":
        Restorer(dumps_dir=args.dumps_dir, db_path=args.db_path, storage_dir=args.storage_dir, avatar_dir=args.avatar_dir).run_restore(emit_progress)


if __name__ == "__main__":
    main()
