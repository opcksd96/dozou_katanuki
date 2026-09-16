import sys, os, sqlite3
_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("SELECT id, source_name, via, original_url, sotwe_url, wayback_url FROM articles WHERE id = '1362279719387848707'")
print("Article:", c.fetchone())

c.execute("SELECT media_id, download_url, type, download_status FROM media WHERE article_id = '1362279719387848707'")
print("DB Media:", c.fetchall())

from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from warcio.archiveiterator import ArchiveIterator

warc_p = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\sayapom4\1362279719387848707\snapshot.warc.gz"
parser = TwitterParser()
with open(warc_p, "rb") as s:
    for r in ArchiveIterator(s):
        if r.rec_type == "response":
            p = parser.parse_record(r.raw_stream.read(), r.rec_headers.get_header("WARC-Target-URI"))
            print("WARC Media count:", len(p.get("media", [])))
            for m in p.get("media", []):
                print("  WARC media item:", m.get("url"))
            break
