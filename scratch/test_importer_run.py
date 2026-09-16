import os, sys
_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.core.warc_importer import WarcImporter

warc_path = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\msluo14\1907698917766213981\snapshot.warc.gz"

importer = WarcImporter(warc_path=warc_path, db_path="archive.db", offline=True)
audit = importer.audit_warc()
print("Audit:", audit)
