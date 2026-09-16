import os, sys

_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.core.restorer import Restorer

pdir = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\msluo14\1907698917766213981"
restorer = Restorer(dumps_dir="backups/dumps/twitter", db_path="archive.db")
print("Running _process_dir for:", pdir)
res = restorer._process_dir(pdir)
print("Result:", res)
