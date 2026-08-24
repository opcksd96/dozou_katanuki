# plugins/twitter/scraper/sources/wayback_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys
from typing import Any, Callable, Dict, List, Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource
from plugins.base.scraper.core.base_scraper import BaseScraper


class WaybackSource(BaseSource):
    """Wayback Machine CDX API 魚拓探索ソース"""

    def __init__(self, platform: str = "twitter", output_dir: str = "backups/dumps", priority: int = 50):
        super().__init__(name="wayback", priority=priority)
        self.scraper = BaseScraper(platform=platform, output_dir=output_dir)

    def fetch_account(
        self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        self.scraper.get_cdx_url_pattern = lambda acc: f"twitter.com/{clean_acc}/status/*"
        snapshots = self.scraper.search_cdx(clean_acc, limit=limit, log_fn=log_fn)
        records = []
        for s in snapshots:
            orig, ts = s.get("original", ""), s.get("timestamp", "")
            raw = self.scraper.fetch_snapshot(ts, orig, account=clean_acc)
            if raw:
                records.append({"source": "wayback", "uri": orig, "timestamp": ts, "raw_data": raw})
        return records

    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip() or "*"
        target_url = f"twitter.com/{clean_acc}/status/{post_id}"
        snapshots = self.scraper.search_cdx(clean_acc, limit=1, log_fn=log_fn)
        if not snapshots: return None
        s = snapshots[0]
        raw = self.scraper.fetch_snapshot(s.get("timestamp", ""), s.get("original", target_url), account=clean_acc)
        return {"source": "wayback", "uri": s.get("original", target_url), "timestamp": s.get("timestamp", ""), "raw_data": raw} if raw else None
