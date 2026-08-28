# plugins/twitter/scraper/sources/sotwe_source.py (SPEC-PLUGIN-001 / 100行以下)
import time
from typing import Any, Callable, Dict, List, Optional
from seleniumbase import SB

try:
    from plugins.twitter.scraper.parsers.sotwe_parser import parse_sotwe_vue_tweets, parse_sotwe_html_tweets
    from plugins.twitter.scraper.parsers.sotwe_extractors import VUE_EXTRACT_JS
except ImportError:
    from parsers.sotwe_parser import parse_sotwe_vue_tweets, parse_sotwe_html_tweets
    from parsers.sotwe_extractors import VUE_EXTRACT_JS


class SotweSource:
    def __init__(self, name: str = "sotwe", priority: int = 20, timeout: int = 25):
        self.name, self.priority, self.timeout = name, priority, timeout

    def is_available(self) -> bool:
        return True

    def fetch_account(
        self, account: str, limit: int = 50, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        url = f"https://www.sotwe.com/{account}"
        if log_fn: log_fn(f"[SotweSource:SB] Opening Web UI (Infinite Scroll): {url}")

        try:
            with SB(uc=True, headless=True) as sb:
                sb.open(url)
                sb.wait_for_element_present("div.tweet-card", timeout=12)
                
                total_posts = self._detect_total_posts(sb)
                if total_posts and log_fn:
                    log_fn(f"[SotweSource:SB] Header Total Posts detected: {total_posts}")

                target = min(limit, total_posts) if (limit > 0 and total_posts) else (limit or total_posts or 0)
                last_count, stagnation = 0, 0
                max_scrolls = max(5, (target // 10) + 3) if target > 0 else 30
                
                for scroll_idx in range(max_scrolls):
                    sb.execute_script("window.scrollTo(0, document.body.scrollHeight);")
                    time.sleep(1.8)
                    
                    cards = sb.find_elements("div.tweet-card")
                    current_count = len(cards)
                    if log_fn:
                        tot_disp = target or total_posts or 'all'
                        log_fn(f"[SotweSource:SB] Scroll {scroll_idx+1}/{max_scrolls}: {current_count} cards (Target: {tot_disp})")
                    
                    if target > 0 and current_count >= target:
                        break
                    if current_count == last_count:
                        stagnation += 1
                        if stagnation >= 2:
                            if log_fn: log_fn("[SotweSource:SB] End of timeline detected (no new cards).")
                            break
                    else:
                        stagnation = 0
                    last_count = current_count

                raw_vue = sb.execute_script(VUE_EXTRACT_JS) or []
                if raw_vue:
                    records = parse_sotwe_vue_tweets(raw_vue, account)
                else:
                    records = parse_sotwe_html_tweets(sb.get_page_source(), account)

            result = records[:limit] if limit > 0 else records
            if log_fn:
                if total_posts:
                    exp = min(total_posts, limit if limit > 0 else total_posts)
                    status = "SUCCESS" if len(result) >= exp else "PARTIAL"
                    log_fn(f"[SotweSource:SB] Verification [{status}]: {len(result)}/{total_posts} posts captured.")
                else:
                    log_fn(f"[SotweSource:SB] Completed: {len(result)} posts captured.")
            return result
        except Exception as e:
            if log_fn: log_fn(f"[SotweSource:SB] Error: {e}")
            return []

    def _detect_total_posts(self, sb) -> Optional[int]:
        try:
            val = sb.execute_script("""
                const xpath = "//div[contains(@class, 'heading')]//span[normalize-space()='Posts']/preceding-sibling::span[1]";
                const node = document.evaluate(xpath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
                if (!node) return null;
                const txt = node.textContent.trim().replace(/,/g, '');
                if (txt.endsWith('K') || txt.endsWith('k')) return Math.round(parseFloat(txt) * 1000);
                if (txt.endsWith('M') || txt.endsWith('m')) return Math.round(parseFloat(txt) * 1000000);
                const n = parseInt(txt, 10);
                return isNaN(n) ? null : n;
            """)
            return int(val) if val is not None else None
        except Exception:
            return None

    def fetch_post(self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None) -> Optional[Dict[str, Any]]:
        posts = self.fetch_account(account or "i", limit=20, log_fn=log_fn)
        return next((p for p in posts if p.get("post", {}).get("id") == post_id), None)
