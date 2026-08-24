# plugins/twitter/scraper/sources/sotwe_source.py (SPEC-PLUGIN-001 / 100行以下)
import json, os, sys
from typing import Any, Callable, Dict, List, Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource


class SotweSource(BaseSource):
    """Sotwe (sotwe.com) Twitterミラー・アーカイブ探索ソース"""
    API_BASE = "https://api.sotwe.com/api/v2"

    def __init__(self, priority: int = 20):
        super().__init__(name="sotwe", priority=priority)

    def fetch_account(
        self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        records: List[Dict[str, Any]] = []
        try:
            url = f"{self.API_BASE}/user/{clean_acc}"
            resp = self.session.get(url, timeout=self.timeout)
            if resp.status_code == 200:
                data = resp.json()
                data_list = data.get("data", []) if isinstance(data, dict) else (data if isinstance(data, list) else [])
                for item in data_list:
                    p_id = str(item.get("id") or item.get("id_str") or "")
                    uri = f"https://twitter.com/{clean_acc}/status/{p_id}" if p_id else f"https://sotwe.com/{clean_acc}"
                    records.append({"source": "sotwe", "uri": uri, "raw_data": item, "id": p_id})
                    if limit > 0 and len(records) >= limit: break
        except Exception as e:
            if log_fn: log_fn(f"[SotweSource] Failed to fetch @{clean_acc}: {e}")
        return records

    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        try:
            url = f"{self.API_BASE}/post/{post_id}"
            resp = self.session.get(url, timeout=self.timeout)
            if resp.status_code == 200:
                data = resp.json()
                if data:
                    uri = f"https://twitter.com/{account or 'unknown'}/status/{post_id}"
                    return {"source": "sotwe", "uri": uri, "raw_data": data, "id": post_id}
        except Exception as e:
            if log_fn: log_fn(f"[SotweSource] Failed to fetch post {post_id}: {e}")
        return None
