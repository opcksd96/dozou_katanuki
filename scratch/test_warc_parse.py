import os, sys, gzip, json
from warcio.archiveiterator import ArchiveIterator

_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser

warc_path = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\msluo14\1907698917766213981\snapshot.warc.gz"

parser = TwitterParser()

with open(warc_path, "rb") as s:
    for idx, r in enumerate(ArchiveIterator(s), start=1):
        uri = r.rec_headers.get_header("WARC-Target-URI") or ""
        print(f"Record {idx}: Type={r.rec_type}, URI={uri}")
        if r.rec_type == "response":
            c_type = (r.http_headers.get_header("Content-Type") if r.http_headers else "") or ""
            print(f"  Content-Type: {c_type}")
            raw = r.raw_stream.read()
            print(f"  Raw len: {len(raw)}")
            parsed = parser.parse_record(raw, uri)
            print("  Parsed result:")
            if parsed:
                print("    Account:", parsed.get("account"))
                print("    Post:", parsed.get("post", {}).get("id"), parsed.get("post", {}).get("full_text")[:30])
                print("    Media count:", len(parsed.get("media", [])))
                for m in parsed.get("media", []):
                    print("      Media:", m)
            else:
                print("    Parsed is None!")
