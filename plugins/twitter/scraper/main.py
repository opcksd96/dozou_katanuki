# plugins/twitter/scraper/main.py (100行以下)
import argparse, concurrent.futures, json, os, sqlite3, sys, time
CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.abspath(os.path.join(CURRENT_DIR, "../../.."))
for d in [PROJECT_ROOT, CURRENT_DIR]:
    if d not in sys.path: sys.path.insert(0, d)
from core.downloader import Downloader; from core.mutator import Mutator; from core.restorer import Restorer
from core.scraper import Scraper; from core.translator import Translator; from core.warc_importer import WarcImporter
from parsers.twitter_parser import TwitterParser


def emit_progress(current: int, total: int, message: str) -> None:
    print(f"PROGRESS: {current}/{total} | {message}", flush=True)


def run_batch_translate(db_path: str, article_id: str = "", account: str = "", overwrite: bool = False, limit: int = 500, dry_run: bool = False) -> None:
    trans = Translator()
    with sqlite3.connect(db_path) as conn:
        cur = conn.cursor()
        if article_id: rows = cur.execute("SELECT id, full_text FROM articles WHERE id = ?", (article_id,)).fetchall()
        else:
            q, p = "SELECT id, full_text FROM articles WHERE 1=1", []
            if account: q += " AND account_id = (SELECT numeric_id FROM accounts WHERE username = ? OR numeric_id = ?)"; p.extend([account, account])
            if not overwrite: q += " AND (full_text_ja IS NULL OR full_text_ja = '' OR (full_text_ja = full_text AND full_text_en = full_text AND lang != 'ja'))"
            q += " LIMIT ?"; p.append(limit); rows = cur.execute(q, p).fetchall()
        if dry_run and article_id and rows: return print(f"JSON:{json.dumps(trans.translate_article(rows[0][1]), ensure_ascii=False)}")
        total = len(rows); emit_progress(0, max(total, 1), f"[PHASE-5:TRANSLATE] Translating {total} articles...")
        for idx, (aid, ftext) in enumerate(rows, start=1):
            t = trans.translate_article(ftext)
            cur.execute("UPDATE articles SET lang=?, full_text_ja=?, full_text_en=?, full_text_zh=? WHERE id=?", (t["lang"], t["ja"], t["en"], t["zh"], aid))
            conn.commit(); emit_progress(idx, total, f"[PHASE-5:TRANSLATE] [{idx}/{total}] article {aid}")
        emit_progress(total, total, f"[PHASE-5:TRANSLATE] Completed translation for {total} articles.")


def run_auto_salvage(platform: str, account: str, limit: int, db_path: str, storage_dir: str, enable_trans: bool, source_filter: str = "all", max_workers: int = 3, chunk_size: int = 50) -> None:
    max_step = 100 if limit <= 0 else limit
    emit_progress(0, max_step, f"[PHASE-1:ORCHESTRATE] Starting multi-source salvage for @{account} on {platform} (filter={source_filter})...")
    scraper, parser, mutator, trans = Scraper(platform=platform), TwitterParser(), Mutator(db_path=db_path, enable_translation=False), (Translator() if enable_trans else None)
    raw_records = scraper.collect_multi_source(account=account, limit=limit, source_filter=source_filter, log_fn=lambda m: emit_progress(0, max_step, f"[PHASE-1:FETCH] {m}"))
    if not raw_records: return emit_progress(100, 100, f"[PHASE-1:FETCH] Finished: No records retrieved for @{account}.")
    tot = len(raw_records); emit_progress(0, tot, f"[PHASE-2:PARSE] Collected {tot} raw records. Parsing...")

    def _parse_and_translate(rec):
        t0 = time.time(); uri = rec.get("uri") or rec.get("original", ""); raw = rec.get("raw_data") or rec
        p = parser.parse_record(raw, uri)
        if not p: return None, {"status": "PARSE_FAILED", "url": uri, "ms": int((time.time() - t0)*1000)}
        post = p.get("post", {}); ftext = post.get("full_text", "")
        for u in post.get("urls", []):
            if u.get("short_url") and u.get("expanded_url") and u["short_url"] in ftext: ftext = ftext.replace(u["short_url"], u["expanded_url"])
        post["full_text"] = ftext
        if trans:
            t = trans.translate_article(ftext); post["lang"], post["full_text_ja"], post["full_text_en"], post["full_text_zh"] = t.get("lang", "ja"), t.get("ja"), t.get("en"), t.get("zh")
        else: post["lang"], post["full_text_ja"] = "ja", ftext
        return p, {"status": "OK", "url": uri, "id": post.get("id"), "media": len(p.get("media", [])), "ms": int((time.time() - t0)*1000)}

    buffer, journal, saved, done = [], [], 0, 0
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as ex:
        for fut in concurrent.futures.as_completed([ex.submit(_parse_and_translate, r) for r in raw_records]):
            done += 1; parsed_json, log_info = fut.result(); journal.append(log_info)
            if parsed_json: buffer.append(parsed_json)
            if len(buffer) >= chunk_size or (done == tot and buffer):
                t_db = time.time(); c_saved = mutator.upsert_batch(buffer); saved += c_saved
                emit_progress(done, tot, f"[PHASE-6:DRIVER_COMMIT] Saved {c_saved} articles ({saved}/{tot}) in {int((time.time() - t_db)*1000)}ms"); buffer.clear()
            elif done % 5 == 0 or done == tot: emit_progress(done, tot, f"[PHASE-3:PROCESS] [{done}/{tot}] buffer={len(buffer)} status={log_info.get('status')}")
    s_ok = sum(1 for j in journal if j.get("status") == "OK")
    emit_progress(100, 100, f"[SALVAGE_SUMMARY] Completed. Target={tot} | Parsed={s_ok} | Errors={tot - s_ok} | DB_Saved={saved}")


def main():
    p = argparse.ArgumentParser(description="Twitter Multi-Source Scraper Sidecar")
    p.add_argument("-m", "--mode", choices=["auto", "manual", "download", "escalate", "poll", "restore", "translate", "smart_recovery", "thunder"], default="auto")
    p.add_argument("-p", "--platform", default="twitter"); p.add_argument("-a", "--account", default=""); p.add_argument("-l", "--limit", type=int, default=0)
    p.add_argument("-s", "--source", default="all", choices=["all", "wayback", "sotwe", "twistalker", "nitter", "official"])
    p.add_argument("-w", "--warc-path", default=""); p.add_argument("--dumps-dir", default="backups/dumps"); p.add_argument("--avatar-dir", default="assets/avatars")
    p.add_argument("--media-id", default=""); p.add_argument("--article-id", default=""); p.add_argument("--offline", action="store_true")
    p.add_argument("--no-translate", action="store_true"); p.add_argument("--overwrite", action="store_true"); p.add_argument("--dry-run", action="store_true")
    p.add_argument("--db-path", default="archive.db"); p.add_argument("--storage-dir", default=""); p.add_argument("-j", "--threads", type=int, default=3); p.add_argument("--chunk-size", type=int, default=50)
    args = p.parse_args()

    if args.mode == "translate": run_batch_translate(args.db_path, article_id=args.article_id, account=args.account, overwrite=args.overwrite, limit=args.limit, dry_run=args.dry_run)
    elif args.mode == "auto": run_auto_salvage(args.platform, args.account, args.limit, args.db_path, args.storage_dir, not args.no_translate and not args.offline, source_filter=args.source, max_workers=args.threads, chunk_size=args.chunk_size)
    elif args.mode == "manual": WarcImporter(args.warc_path, db_path=args.db_path, storage_dir=args.storage_dir, offline=args.offline).run_import(emit_progress)
    elif args.mode == "download": Downloader(db_path=args.db_path, storage_dir=args.storage_dir).process_queued_media(args.article_id or None, args.media_id or None, log_fn=lambda c, t, m: emit_progress(c, t, f"[PHASE-DL] {m}"))
    elif args.mode == "escalate": Downloader(db_path=args.db_path, storage_dir=args.storage_dir).escalate_dead_media(log_fn=lambda c, t, m: emit_progress(c, t, f"[PHASE-ESCALATE] {m}"))
    elif args.mode == "poll": Downloader(db_path=args.db_path, storage_dir=args.storage_dir).poll_outsourced_media(log_fn=lambda m: emit_progress(1, 1, f"[PHASE-POLL] {m}"))
    elif args.mode == "smart_recovery": Downloader(db_path=args.db_path, storage_dir=args.storage_dir).smart_recovery_pipeline(log_fn=lambda c, t, m: emit_progress(c, t, f"[SMART-RECOVERY] {m}"))
    elif args.mode == "thunder": Downloader(db_path=args.db_path, storage_dir=args.storage_dir).escalate_to_thunder(log_fn=lambda c, t, m: emit_progress(c, t, f"[THUNDER] {m}"))
    elif args.mode == "restore": Restorer(dumps_dir=args.dumps_dir, db_path=args.db_path, storage_dir=args.storage_dir, avatar_dir=args.avatar_dir, max_workers=args.threads).run_restore(emit_progress)


if __name__ == "__main__": main()
