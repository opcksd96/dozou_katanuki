# plugins/twitter/scraper/core/restorer.py (Twitter特化DRリストア / 100行以下)
import os, sys
from typing import Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.base.scraper.core.restorer import Restorer as BaseRestorer
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser


class Restorer(BaseRestorer):
    """Twitter / X 特化型高速オフライン DB 再構築エンジン"""
    def __init__(self, dumps_dir: str = "backups/dumps", db_path: str = "archive.db", storage_dir: Optional[str] = None, avatar_dir: str = "assets/avatars", max_workers: int = 8):
        super().__init__(dumps_dir=dumps_dir, db_path=db_path, storage_dir=storage_dir, avatar_dir=avatar_dir, max_workers=max_workers, parser=TwitterParser(), platform="twitter")
