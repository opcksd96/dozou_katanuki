import os, sys, time

_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.core.restorer import Restorer

print("=== STARTING FULL WARC RESTORE ===")
start_time = time.time()

def progress(current, total, msg):
    if current % 500 == 0 or current == total or current == 0:
        print(f"[{current}/{total}] ({time.time() - start_time:.1f}s) {msg}")

restorer = Restorer(dumps_dir="backups/dumps/twitter", db_path="archive.db", max_workers=8)
stats = restorer.run_restore(progress_callback=progress)

elapsed = time.time() - start_time
print(f"\n=== RESTORE COMPLETED in {elapsed:.2f}s ===")
print("Stats:", stats)
