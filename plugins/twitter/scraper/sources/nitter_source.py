# plugins/twitter/scraper/sources/nitter_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, re, sys
from typing import Any, Callable, Dict, List, Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource


class NitterSource(BaseSource):
    """Nitter インスタンス分散ミラー探索ソース"""
    DEFAULT_INSTANCES = ["https://nitter.net", "https://nitter.poast.org", "https://nitter.privacydev.net"]

    def __init__(self, instances: Optional[List[str]] = None, priority: int = 40):
        super().__init__(name="nitter", priority=priority)
        self.instances = instances or self.DEFAULT_INSTANCES

    def fetch_account(
        self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        records: List[Dict[str, Any]] = []
        for inst in self.instances:
            try:
                url = f"{inst.rstrip('/')}/{clean_acc}"
                resp = self.session.get(url, timeout=self.timeout)
                if resp.status_code == 200:
                    html_text = resp.text
                    items = re.findall(r'(<div[^>]+class="[^"]*timeline-item[^"]*"[^>]*>.*?</div>\s*</div>\s*</div>)', html_text, re.DOTALL)
                    for it in items:
                        m_id = re.search(r'/status/(\d+)', it)
                        p_id = m_id.group(1) if m_id else ""
                        uri = f"https://twitter.com/{clean_acc}/status/{p_id}" if p_id else f"{inst}/{clean_acc}"
                        records.append({"source": "nitter", "uri": uri, "raw_data": it, "id": p_id})
                        if limit > 0 and len(records) >= limit: break
                    if records: break
            except Exception as e:
                if log_fn: log_fn(f"[NitterSource] Instance {inst} failed: {e}")
        return records

    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip() or "i"
        for inst in self.instances:
            try:
                url = f"{inst.rstrip('/')}/{clean_acc}/status/{post_id}"
                resp = self.session.get(url, timeout=self.timeout)
                if resp.status_code == 200:
                    uri = f"https://twitter.com/{clean_acc}/status/{post_id}"
                    return {"source": "nitter", "uri": uri, "raw_data": resp.text, "id": post_id}
            except Exception as e:
                if log_fn: log_fn(f"[NitterSource] Instance {inst} fetch failed for post {post_id}: {e}")
        return None
