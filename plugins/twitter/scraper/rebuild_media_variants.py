# plugins/twitter/scraper/rebuild_media_variants.py (SPEC-PLUGIN-001 / 100行以下)
import concurrent.futures, os, sqlite3, sys, time
sys.path.insert(0, os.path.abspath("."))
from warcio.archiveiterator import ArchiveIterator
from plugins.base.scraper.core.base_mutator import BaseMutator
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser

def rebuild_variants(dumps_dir: str = "backups/dumps", db_path: str = "archive.db", max_workers: int = 8):
    if not os.path.exists(dumps_dir):
        print(f"Dumps dir {dumps_dir} not found."); return
    mutator = BaseMutator(db_path=db_path, platform="twitter", enable_translation=False)
    parser = TwitterParser()
    warc_files = [os.path.join(root, f) for root, _, files in os.walk(dumps_dir) for f in files if f.endswith(".warc.gz") or f.endswith(".warc")]
    total = len(warc_files)
    print(f"=== Starting media_variants Rebuild from {total} WARCs (threads: {max_workers}) ===")

    def process_warc(w_path: str):
        saved = 0
        try:
            with open(w_path, "rb") as s:
                for r in ArchiveIterator(s):
                    if r.rec_type != "response": continue
                    uri = r.rec_headers.get_header("WARC-Target-URI") or ""
                    ct = (r.http_headers.get_header("Content-Type") if r.http_headers else "") or ""
                    if "json" in ct or "html" in ct or "status" in uri or "tweet" in uri:
                        parsed = parser.parse_record(r.raw_stream.read(), uri)
                        if parsed and mutator.upsert_record(parsed): saved += 1
        except Exception: pass
        return saved

    done, total_saved, t0 = 0, 0, time.time()
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = {executor.submit(process_warc, p): p for p in warc_files}
        for fut in concurrent.futures.as_completed(futures):
            done += 1; total_saved += fut.result()
            if done % 100 == 0 or done == total:
                elapsed = time.time() - t0
                print(f"Progress: [{done}/{total}] WARCs scanned | {total_saved} records upserted ({elapsed:.1f}s)")

    with sqlite3.connect(db_path) as conn:
        cur = conn.cursor()
        v_cnt = cur.execute("SELECT COUNT(*) FROM media_variants").fetchone()[0]
        m_cnt = cur.execute("SELECT COUNT(*) FROM media").fetchone()[0]
        img_v = cur.execute("SELECT COUNT(*) FROM media_variants WHERE bit_rate IN (10000, 5000) OR content_type LIKE 'image%'").fetchone()[0]
        vid_v = cur.execute("SELECT COUNT(*) FROM media_variants WHERE content_type LIKE 'video%' OR (bit_rate != 10000 AND bit_rate != 5000)").fetchone()[0]
    print(f"\n=== REBUILD COMPLETED in {time.time() - t0:.1f}s ===")
    print(f"Total media records in DB: {m_cnt}")
    print(f"Total media_variants registered: {v_cnt} (Images: {img_v}, Videos: {vid_v})")

if __name__ == "__main__":
    dumps = sys.argv[1] if len(sys.argv) > 1 else "backups/dumps"
    db = sys.argv[2] if len(sys.argv) > 2 else "archive.db"
    rebuild_variants(dumps_dir=dumps, db_path=db)
