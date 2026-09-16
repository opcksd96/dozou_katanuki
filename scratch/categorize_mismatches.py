import os, sys, sqlite3
from warcio.archiveiterator import ArchiveIterator

_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("SELECT article_id, COUNT(*) FROM media GROUP BY article_id")
db_media_counts = dict(c.fetchall())

parser = TwitterParser()
dumps_dir = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter"

warc_greater = 0 # WARC has more media than DB
db_zero = 0      # DB has 0 media but WARC has >=1
partial_loss = 0 # DB has some, but less than WARC (like 4 vs 1)
db_greater = 0   # DB has more media than WARC

sample_partial = []

for root, dirs, files in os.walk(dumps_dir):
    if "snapshot.warc.gz" in files:
        warc_path = os.path.join(root, "snapshot.warc.gz")
        tweet_id = os.path.basename(root)
        try:
            with open(warc_path, "rb") as s:
                for r in ArchiveIterator(s):
                    if r.rec_type != "response": continue
                    c_type = (r.http_headers.get_header("Content-Type") if r.http_headers else "") or ""
                    uri = r.rec_headers.get_header("WARC-Target-URI") or ""
                    if "json" in c_type or "status" in uri:
                        raw = r.raw_stream.read()
                        parsed = parser.parse_record(raw, uri)
                        if parsed:
                            warc_cnt = len(parsed.get("media", []))
                            db_cnt = db_media_counts.get(tweet_id, 0)
                            if warc_cnt > db_cnt:
                                warc_greater += 1
                                if db_cnt == 0:
                                    db_zero += 1
                                else:
                                    partial_loss += 1
                                    sample_partial.append((tweet_id, parsed.get("account", {}).get("username"), warc_cnt, db_cnt))
                            elif db_cnt > warc_cnt:
                                db_greater += 1
                        break
        except Exception:
            pass

print(f"Total with WARC > DB: {warc_greater}")
print(f"  - DB completely missing media (WARC >= 1, DB = 0): {db_zero}")
print(f"  - Partial loss like this tweet (WARC > DB > 0): {partial_loss}")
print(f"Total with DB > WARC: {db_greater}")

print("\nSample partial losses (like 4 vs 1):")
for s in sample_partial[:15]:
    print(f"  Tweet {s[0]} (@{s[1]}): WARC has {s[2]} media, DB only has {s[3]}")
