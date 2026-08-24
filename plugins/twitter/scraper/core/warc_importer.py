# plugins/twitter/scraper/core/warc_importer.py (Twitter特化WARCインポーター / 100行以下)
import os, sys
from typing import Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.base.scraper.core.warc_importer import WarcImporter as BaseWarcImporter
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser


class WarcImporter(BaseWarcImporter):
    """Twitter / X 特化型手動WARCインポーター"""
    def __init__(self, warc_path: str, db_path: str = "archive.db", storage_dir: Optional[str] = None, offline: bool = True):
        super().__init__(warc_path=warc_path, db_path=db_path, storage_dir=storage_dir, offline=offline, parser=TwitterParser(), platform="twitter")
