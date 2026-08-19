# plugins/twitter/scraper/main.py (100行以下)
import argparse
import os
import sys
import time

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
if CURRENT_DIR not in sys.path:
    sys.path.insert(0, CURRENT_DIR)

from core.downloader import Downloader
from core.mutator import Mutator
from core.restorer import Restorer
from core.scraper import Scraper
from core.warc_importer import WarcImporter
from parsers.twitter_parser import TwitterParser


def emit_progress(current: int, total: int, message: str) -> None:
    """Go Middleware 互換の進捗 stdout フラッシュ出力"""
    print(f"PROGRESS: {current}/{total} | {message}", flush=True)


def run_auto_salvage(platform: str, account: str, limit: int, db_path: str, storage_dir: str) -> None:
    emit_progress(0, limit, f"Starting auto salvage for @{account} on {platform}...")
    scraper = Scraper(platform=platform)
    parser = TwitterParser()
    mutator = Mutator(db_path=db_path)
    downloader = Downloader(db_path=db_path, storage_dir=storage_dir)

    snapshots = scraper.search_cdx(account=account, limit=limit)
    total_snaps = len(snapshots)
    if total_snaps == 0:
        emit_progress(limit, limit, f"No CDX snapshots found for @{account}.")
        return

    emit_progress(0, total_snaps, f"Found {total_snaps} snapshots. Starting extraction...")
    success_count = 0
    for idx, snap in enumerate(snapshots, start=1):
        ts = snap.get("timestamp", "")
        orig = snap.get("original", "")
        emit_progress(idx, total_snaps, f"Fetching [{idx}/{total_snaps}] ({ts})...")
        raw_text = scraper.fetch_snapshot(ts, orig, account=account)
        if raw_text:
            parsed = parser.parse_record(raw_text, orig)
            if parsed and mutator.upsert_record(parsed):
                downloader.process_queued_media(parsed["post"]["id"])
                success_count += 1
        time.sleep(0.05)
    emit_progress(total_snaps, total_snaps, f"Completed: Saved {success_count}/{total_snaps} posts.")


def main():
    parser = argparse.ArgumentParser(description="dozou_katanuki Twitter Scraper Sidecar")
    parser.add_argument("-m", "--mode", choices=["auto", "manual", "download", "poll", "restore"], default="auto")
    parser.add_argument("-p", "--platform", default="twitter")
    parser.add_argument("-a", "--account", default="")
    parser.add_argument("-l", "--limit", type=int, default=50)
    parser.add_argument("-w", "--warc-path", default="")
    parser.add_argument("--dumps-dir", default="backups/dumps")
    parser.add_argument("--avatar-dir", default="assets/avatars")
    parser.add_argument("--media-id", default="")
    parser.add_argument("--article-id", default="")
    parser.add_argument("--offline", action="store_true")
    parser.add_argument("--db-path", default="archive.db")
    parser.add_argument("--storage-dir", default="")
    args = parser.parse_args()

    dl = Downloader(db_path=args.db_path, storage_dir=args.storage_dir)
    if args.mode == "auto":
        if not args.account:
            print("[FATAL] --account is required in auto mode", file=sys.stderr); sys.exit(1)
        run_auto_salvage(args.platform, args.account, args.limit, args.db_path, args.storage_dir)
    elif args.mode == "manual":
        if not args.warc_path:
            print("[FATAL] --warc-path is required in manual mode", file=sys.stderr); sys.exit(1)
        WarcImporter(args.warc_path, db_path=args.db_path, storage_dir=args.storage_dir).run_import(progress_callback=emit_progress)
    elif args.mode == "download":
        emit_progress(0, 1, f"Processing media download (media_id: {args.media_id or 'all_queued'})...")
        c = dl.process_queued_media(article_id=args.article_id or None, media_id=args.media_id or None)
        emit_progress(1, 1, f"Finished media download. Processed: {c}")
    elif args.mode == "poll":
        emit_progress(0, 1, "Polling outsourced media directory...")
        c = dl.poll_outsourced_media()
        emit_progress(1, 1, f"Polling completed. Salvaged: {c}")
    elif args.mode == "restore":
        res = Restorer(dumps_dir=args.dumps_dir, db_path=args.db_path, storage_dir=args.storage_dir, avatar_dir=args.avatar_dir)
        res.run_restore(progress_callback=emit_progress)


if __name__ == "__main__":
    main()
