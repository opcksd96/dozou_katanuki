import os, sys
_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from warcio.archiveiterator import ArchiveIterator
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from plugins.base.scraper.core.base_mutator import BaseMutator

warc_p = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\msluo14\1907698917766213981\snapshot.warc.gz"
parser = TwitterParser()
mutator = BaseMutator(db_path="archive.db", platform="twitter")

with open(warc_p, "rb") as s:
    for r in ArchiveIterator(s):
        if r.rec_type == "response":
            parsed = parser.parse_record(r.raw_stream.read(), r.rec_headers.get_header("WARC-Target-URI"))
            print("Running upsert_batch directly...")
            res = mutator.upsert_batch([parsed])
            print("upsert_batch result:", res)
            break

import sqlite3
conn = sqlite3.connect('archive.db')
c = conn.cursor()
c.execute("SELECT media_id, article_id, type, download_status FROM media WHERE article_id LIKE '%1907698917766213981%'")
for row in c.fetchall():
    print("DB row:", row)
