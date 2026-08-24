# plugins/twitter/scraper/core/scraper.py (Twitter特化具象スクレイパー / 100行以下)
import os, sys

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_scraper import BaseScraper


class Scraper(BaseScraper):
    """Twitter / X 特化型 Wayback CDX走査 & HTTPフェッチ & 原本WARC保存エンジン"""
    def __init__(self, platform: str = "twitter", output_dir: str = "backups/dumps"):
        super().__init__(platform=platform, output_dir=output_dir)

    def get_cdx_url_pattern(self, account: str) -> str:
        return f"twitter.com/{account}/status/*"
