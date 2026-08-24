# plugins/twitter/scraper/sources/official_source.py (SPEC-PLUGIN-001 / 100行以下)
import json, os, re, sys
from typing import Any, Callable, Dict, List, Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource


class OfficialSource(BaseSource):
    """X / Twitter 本家・Syndication API 公開取得ソース"""
    SYNDICATION_URL = "https://cdn.syndication.twimg.com/tweet-result"

    def __init__(self, priority: int = 10):
        super().__init__(name="official", priority=priority)

    def fetch_account(
        self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        # 本家Webはログイン/トークンなしではアカウント全体のリスト取得が制限されるため、フォールバック網の一部として稼働
        return []

    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        try:
            params = {"id": post_id, "token": "x"}
            resp = self.session.get(self.SYNDICATION_URL, params=params, timeout=self.timeout)
            if resp.status_code == 200:
                data = resp.json()
                if data and isinstance(data, dict) and "id_str" in data:
                    uri = f"https://twitter.com/{account or data.get('user', {}).get('screen_name', 'unknown')}/status/{post_id}"
                    return {"source": "official", "uri": uri, "raw_data": data, "id": post_id}
        except Exception as e:
            if log_fn: log_fn(f"[OfficialSource] Fetch failed for post {post_id}: {e}")
        return None
