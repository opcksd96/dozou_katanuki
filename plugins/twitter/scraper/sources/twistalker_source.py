# plugins/twitter/scraper/sources/twistalker_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, time
from typing import Any, Callable, Dict, List, Optional
from seleniumbase import SB

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource

try:
    from plugins.twitter.scraper.parsers.twistalker_parser import parse_twistalker_api_tweets, parse_twistalker_html_tweets
except ImportError:
    from parsers.twistalker_parser import parse_twistalker_api_tweets, parse_twistalker_html_tweets

class TwistalkerSource(BaseSource):
    """TwStalker (twstalker.com) SeleniumBase & 内部API探索ソース"""
    BASE_URL = "https://twstalker.com"

    def __init__(self, priority: int = 25, timeout: float = 25.0):
        super().__init__(name="twistalker", priority=priority, timeout=timeout)

    def fetch_account(self, account: str, limit: int = 50, log_fn: Optional[Callable[[str], None]] = None) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        url = f"{self.BASE_URL}/{clean_acc}"
        if log_fn: log_fn(f"[TwistalkerSource:SB] Opening Web UI: {url}")
        results, seen_ids = [], set()
        try:
            with SB(uc=True, headless=True) as sb:
                sb.open(url); time.sleep(3.5)
                for p in parse_twistalker_html_tweets(sb.get_page_source(), clean_acc):
                    tid = p.get("post", {}).get("id")
                    if tid and tid not in seen_ids: seen_ids.add(tid); results.append(p)
                if log_fn: log_fn(f"[TwistalkerSource:SB] Initial page captured {len(results)} posts")

                page_idx = 1
                while (limit <= 0 or len(results) < limit) and page_idx < 20:
                    api_data = self._fetch_next_page(sb)
                    if not api_data or not api_data.get("tweets"): break
                    added = 0
                    for p in parse_twistalker_api_tweets(api_data.get("tweets", {}), clean_acc):
                        tid = p.get("post", {}).get("id")
                        if tid and tid not in seen_ids: seen_ids.add(tid); results.append(p); added += 1
                    if log_fn: log_fn(f"[TwistalkerSource:API] Page {page_idx+1}: +{added} posts (Total: {len(results)})")
                    if added == 0 or not api_data.get("cursor"): break
                    page_idx += 1; time.sleep(1.0)

            result = results[:limit] if limit > 0 else results
            if result:
                try:
                    from plugins.twitter.scraper.core.warc_archiver import WarcArchiver
                    c = WarcArchiver().archive_posts(result, platform="twitter")
                    if log_fn: log_fn(f"[TwistalkerSource:WARC] Dumped & enriched {c} posts to backups/dumps/twitter/{clean_acc}/")
                except Exception as we:
                    if log_fn: log_fn(f"[TwistalkerSource:WARC] Archival Warning: {we}")
            return result
        except Exception as e:
            if log_fn: log_fn(f"[TwistalkerSource:SB] Error: {e}")
            return results

    def _fetch_next_page(self, sb) -> Optional[Dict[str, Any]]:
        js_fetch = """
        const done = arguments[0];
        const btn = document.querySelector('.add-nw-event');
        if (!btn || !btn.getAttribute('data-cursor')) return done(null);
        const fd = new URLSearchParams();
        fd.append('page', btn.getAttribute('data-ec') || '2');
        fd.append('cursor', btn.getAttribute('data-cursor'));
        fd.append('data', btn.getAttribute('data-query'));
        fd.append('action', 'profile');

        fetch('/service/api', {
            method: 'POST', body: fd,
            headers: {'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8', 'X-Requested-With': 'XMLHttpRequest'}
        }).then(r => r.json()).then(res => {
            if (res.cursor) {
                btn.setAttribute('data-cursor', res.cursor);
                btn.setAttribute('data-ec', String(Number(btn.getAttribute('data-ec') || 2) + 1));
            }
            done(res);
        }).catch(err => done(null));
        """
        try: return sb.execute_async_script(js_fetch)
        except Exception: return None

    def fetch_post(self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None) -> Optional[Dict[str, Any]]:
        posts = self.fetch_account(account or "i", limit=20, log_fn=log_fn)
        return next((p for p in posts if p.get("post", {}).get("id") == post_id), None)
