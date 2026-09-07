# plugins/twitter/scraper/sources/official_source.py (SPEC-PLUGIN-001 / 100行以下)
import json, os, sys, time
from typing import Any, Callable, Dict, List, Optional
from bs4 import BeautifulSoup
from seleniumbase import SB

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource

try:
    from plugins.twitter.scraper.parsers.x_parser import XParser
except ImportError:
    from parsers.x_parser import XParser

class OfficialSource(BaseSource):
    """X / Twitter (x.com) 公式Syndication API ＆ ブラウザ巡回ハイブリッドソース"""
    SYNDICATION_URL = "https://cdn.syndication.twimg.com/tweet-result"

    def __init__(self, cookies: Optional[Dict[str, str]] = None, priority: int = 15, timeout: float = 25.0):
        super().__init__(name="official", priority=priority, timeout=timeout)
        self.parser = XParser()
        self.cookies = cookies or {}

    def set_cookies(self, cookies: Dict[str, str]) -> None:
        self.cookies = cookies or {}

    def _inject_cookies(self, sb) -> None:
        if not self.cookies: return
        for name, value in self.cookies.items():
            try: sb.add_cookie({"name": name, "value": value, "domain": ".x.com"})
            except Exception: pass

    def fetch_post(self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None) -> Optional[Dict[str, Any]]:
        clean_id = str(post_id).strip()
        try:
            resp = self.session.get(self.SYNDICATION_URL, params={"id": clean_id, "token": "x"}, timeout=self.timeout)
            if resp.status_code == 200:
                parsed = self.parser.parse_syndication_json(resp.json(), default_acc=account)
                if parsed: return parsed
        except Exception as e:
            if log_fn: log_fn(f"[OfficialSource:Syndication] Post {clean_id} error: {e}")
        return self._fetch_post_browser(clean_id, account, log_fn)

    def _fetch_post_browser(self, post_id: str, account: str, log_fn: Optional[Callable[[str], None]]) -> Optional[Dict[str, Any]]:
        acc = account.lstrip("@").strip() or "i"
        url = f"https://x.com/{acc}/status/{post_id}"
        if log_fn: log_fn(f"[OfficialSource:Browser] Opening {url}")
        try:
            with SB(uc=True, headless=True) as sb:
                sb.open("https://x.com"); self._inject_cookies(sb)
                sb.open(url); time.sleep(3.5)
                soup = BeautifulSoup(sb.get_page_source(), "html.parser")
                arts = soup.find_all("article")
                if arts: return self.parser.parse_html_card(arts[0], default_acc=account)
        except Exception as e:
            if log_fn: log_fn(f"[OfficialSource:Browser] Error: {e}")
        return None

    def fetch_account(self, account: str, limit: int = 10, log_fn: Optional[Callable[[str], None]] = None) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        url = f"https://x.com/{clean_acc}"
        if log_fn: log_fn(f"[OfficialSource] Scanning @{clean_acc}")
        results, seen_ids = [], set()
        try:
            with SB(uc=True, headless=True) as sb:
                sb.open("https://x.com"); self._inject_cookies(sb)
                sb.open(url); time.sleep(4)
                soup = BeautifulSoup(sb.get_page_source(), "html.parser")
                arts = soup.find_all("article")
                for art in arts:
                    parsed = self.parser.parse_html_card(art, default_acc=clean_acc)
                    if not parsed: continue
                    pid = parsed.get("post", {}).get("id")
                    if not pid or pid in seen_ids: continue
                    seen_ids.add(pid)
                    detailed = self.fetch_post(pid, account=clean_acc)
                    results.append(detailed if detailed else parsed)
                    if 0 < limit <= len(results): break
            if results:
                try:
                    from plugins.twitter.scraper.core.warc_archiver import WarcArchiver
                    WarcArchiver().archive_posts(results, platform="twitter")
                except Exception: pass
        except Exception as e:
            if log_fn: log_fn(f"[OfficialSource] Error scanning @{clean_acc}: {e}")
        return results
