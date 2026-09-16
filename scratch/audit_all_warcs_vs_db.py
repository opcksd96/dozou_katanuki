import os, sys, sqlite3, json, re
from warcio.archiveiterator import ArchiveIterator

_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser

conn = sqlite3.connect('archive.db')
c = conn.cursor()

# Get map of article_id -> media count in DB
c.execute("SELECT article_id, COUNT(*) FROM media GROUP BY article_id")
db_media_counts = dict(c.fetchall())

print(f"Total articles with media in DB: {len(db_media_counts)}")

parser = TwitterParser()
dumps_dir = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter"

mismatches = []
total_warcs = 0
parsed_warcs = 0

for root, dirs, files in os.walk(dumps_dir):
    if "snapshot.warc.gz" in files:
        warc_path = os.path.join(root, "snapshot.warc.gz")
        tweet_id = os.path.basename(root)
        total_warcs += 1
        
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
                            parsed_warcs += 1
                            warc_media_count = len(parsed.get("media", []))
                            db_count = db_media_counts.get(tweet_id, 0)
                            
                            if warc_media_count != db_count:
                                mismatches.append({
                                    "tweet_id": tweet_id,
                                    "account": parsed.get("account", {}).get("username"),
                                    "warc_media_count": warc_media_count,
                                    "db_media_count": db_count,
                                    "text": parsed.get("post", {}).get("full_text", "")[:40].replace("\n", " ")
                                })
                        break
        except Exception as e:
            pass

print(f"Total WARCs scanned: {total_warcs}")
print(f"Successfully parsed WARCs: {parsed_warcs}")
print(f"Mismatched media count tweets: {len(mismatches)}")

print("\n=== SAMPLE MISMATCHES (First 30) ===")
for m in mismatches[:30]:
    print(f"Tweet {m['tweet_id']} (@{m['account']}): WARC={m['warc_media_count']}枚 vs DB={m['db_media_count']}枚 | {m['text']}")
