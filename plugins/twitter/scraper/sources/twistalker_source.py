# plugins/twitter/scraper/sources/twistalker_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, re, sys
from typing import Any, Callable, Dict, List, Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource


class TwistalkerSource(BaseSource):
    """Twistalker (twistalker.com) Twitterアーカイブ探索ソース"""
    BASE_URL = "https://twistalker.com"

    def __init__(self, priority: int = 30):
        super().__init__(name="twistalker", priority=priority)

    def fetch_account(
        self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        records: List[Dict[str, Any]] = []
        try:
            url = f"{self.BASE_URL}/{clean_acc}"
            resp = self.session.get(url, timeout=self.timeout)
            if resp.status_code == 200:
                html_text = resp.text
                articles = re.findall(r'(<div[^>]+class="[^"]*post[^"]*"[^>]*>.*?</div>\s*</div>)', html_text, re.DOTALL)
                for art in articles:
                    m_id = re.search(r'data-id="(\d+)"', art) or re.search(r'/status/(\d+)', art)
                    p_id = m_id.group(1) if m_id else ""
                    uri = f"https://twitter.com/{clean_acc}/status/{p_id}" if p_id else f"{self.BASE_URL}/{clean_acc}"
                    records.append({"source": "twistalker", "uri": uri, "raw_data": art, "id": p_id})
                    if limit > 0 and len(records) >= limit: break
        except Exception as e:
            if log_fn: log_fn(f"[TwistalkerSource] Failed to fetch @{clean_acc}: {e}")
        return records

    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip() or "i"
        try:
            url = f"{self.BASE_URL}/{clean_acc}/status/{post_id}"
            resp = self.session.get(url, timeout=self.timeout)
            if resp.status_code == 200:
                uri = f"https://twitter.com/{clean_acc}/status/{post_id}"
                return {"source": "twistalker", "uri": uri, "raw_data": resp.text, "id": post_id}
        except Exception as e:
            if log_fn: log_fn(f"[TwistalkerSource] Failed to fetch post {post_id}: {e}")
        return None
