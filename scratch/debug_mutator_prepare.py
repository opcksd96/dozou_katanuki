import os, sys, json
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
            print("Parsed media count:", len(parsed.get("media", [])))
            print("Media items:")
            for m in parsed.get("media", []):
                print(" ", m)
            
            with mutator._get_conn() as conn:
                wl = mutator._get_active_whitelist(conn)
                print("Active whitelist count:", len(wl))
                print("msluo14 in whitelist?", "msluo14" in wl)
                
                accs, hists, redirs, arts, meds, med_vars, excls, valid_count = mutator._prepare_in_memory([parsed], wl)
                print(f"Prepared: arts={len(arts)}, meds={len(meds)}, excls={len(excls)}")
                for med in meds:
                    print("  Prepared med:", med)
                for exc in excls:
                    print("  Prepared excl:", exc)
            break
