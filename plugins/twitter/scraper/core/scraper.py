# plugins/twitter/scraper/core/scraper.py (Twitter特化具象スクレイパー / 100行以下)
import os, sys
from typing import Any, Callable, Dict, List, Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_scraper import BaseScraper
from plugins.base.scraper.core.source_orchestrator import SourceOrchestrator
from plugins.twitter.scraper.sources import OfficialSource, SotweSource, TwistalkerSource, NitterSource, WaybackSource


class Scraper(BaseScraper):
    """Twitter / X 特化型 Wayback CDX走査 & マルチソース・サルベージ統括エンジン"""

    def __init__(self, platform: str = "twitter", output_dir: str = "backups/dumps"):
        super().__init__(platform=platform, output_dir=output_dir)
        self.orchestrator = SourceOrchestrator()
        self.orchestrator.register(OfficialSource())
        self.orchestrator.register(SotweSource())
        self.orchestrator.register(TwistalkerSource())
        self.orchestrator.register(NitterSource())
        self.orchestrator.register(WaybackSource(platform=platform, output_dir=output_dir))

    def get_cdx_url_pattern(self, account: str) -> str:
        clean = account.lstrip("@").strip()
        return f"twitter.com/{clean}/status/*"

    def collect_multi_source(
        self, account: str, limit: int = 0, source_filter: str = "all", log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        return self.orchestrator.collect(account=account, limit=limit, source_filter=source_filter, log_fn=log_fn)
