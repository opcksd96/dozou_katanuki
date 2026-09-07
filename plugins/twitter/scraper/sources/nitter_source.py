# plugins/twitter/scraper/sources/nitter_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, time
from typing import Any, Callable, Dict, List, Optional
from bs4 import BeautifulSoup
from seleniumbase import SB

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_source import BaseSource

try:
    from plugins.twitter.scraper.parsers.nitter_parser import parse_nitter_card, parse_nitter_html_tweets
except ImportError:
    from parsers.nitter_parser import parse_nitter_card, parse_nitter_html_tweets

class NitterSource(BaseSource):
    """Nitter (nitter.space / クローンURL) ブラウザ走査 ＆ HTMLスクレイピングソース"""
    DEFAULT_CLONE_URL = "https://nitter.space"

    def __init__(self, clone_url: Optional[str] = None, priority: int = 30, timeout: float = 25.0):
        super().__init__(name="nitter", priority=priority, timeout=timeout)
        self.clone_url = (clone_url or os.environ.get("NITTER_CLONE_URL") or self.DEFAULT_CLONE_URL).rstrip("/")

    def get_clone_url(self) -> str:
        return self.clone_url

    def fetch_account(self, account: str, limit: int = 50, log_fn: Optional[Callable[[str], None]] = None) -> List[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip()
        url = f"{self.clone_url}/{clean_acc}"
        if log_fn: log_fn(f"[NitterSource:SB] Opening Clone URL: {url}")
        results, seen_ids = [], set()
        try:
            with SB(uc=True, headless=True) as sb:
                sb.open(url)
                try: sb.wait_for_element_present(".timeline-item", timeout=10)
                except Exception: sb.sleep(3.0)
                page_idx = 1
                while limit <= 0 or len(results) < limit:
                    html_src = sb.get_page_source()
                    added = 0
                    for p in parse_nitter_html_tweets(html_src, default_account=clean_acc):
                        tid = p.get("post", {}).get("id")
                        if tid and tid not in seen_ids:
                            seen_ids.add(tid); results.append(p); added += 1
                    if log_fn: log_fn(f"[NitterSource:Page {page_idx}] +{added} posts (Total: {len(results)})")
                    if added == 0 or (0 < limit <= len(results)): break

                    soup = BeautifulSoup(html_src, "html.parser")
                    cur_link = soup.select_one("a[href*='cursor=']")
                    if not cur_link: break
                    cursor = cur_link.get("href", "").split("cursor=")[-1].split("&")[0]
                    if not cursor: break
                    page_idx += 1
                    next_url = f"{self.clone_url}/{clean_acc}?cursor={cursor}"
                    if log_fn: log_fn(f"[NitterSource:Cursor] Next: {next_url}")
                    sb.open(next_url)
                    try: sb.wait_for_element_present(".timeline-item", timeout=10)
                    except Exception: sb.sleep(3.0)

            result = results[:limit] if limit > 0 else results
            if result:
                try:
                    from plugins.twitter.scraper.core.warc_archiver import WarcArchiver
                    c = WarcArchiver().archive_posts(result, platform="twitter")
                    if log_fn: log_fn(f"[NitterSource:WARC] Dumped & enriched {c} posts")
                except Exception as we:
                    if log_fn: log_fn(f"[NitterSource:WARC] Archival Warning: {we}")
            return result
        except Exception as e:
            if log_fn: log_fn(f"[NitterSource:SB] Error: {e}")
            return results

    def fetch_post(self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None) -> Optional[Dict[str, Any]]:
        clean_acc = account.lstrip("@").strip() or "i"
        url = f"{self.clone_url}/{clean_acc}/status/{post_id}"
        if log_fn: log_fn(f"[NitterSource:Post] Opening status: {url}")
        try:
            with SB(uc=True, headless=True) as sb:
                sb.open(url)
                try: sb.wait_for_element_present(".main-tweet, .timeline-item", timeout=10)
                except Exception: sb.sleep(3.0)
                soup = BeautifulSoup(sb.get_page_source(), "html.parser")
                main_tweet = soup.select_one(".main-tweet, .timeline-item")
                if main_tweet:
                    return parse_nitter_card(main_tweet, default_account=clean_acc, fallback_id=post_id)
        except Exception as e:
            if log_fn: log_fn(f"[NitterSource:Post] Error for {post_id}: {e}")
        return None
