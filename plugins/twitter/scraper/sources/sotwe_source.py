"""
sotwe_source.py
SeleniumBase UC ModeでSotwe Web UIを展開し、DOMツリーから全ツイートを回収するソース。
"""
import time
from typing import Any, Callable, Dict, List, Optional
from seleniumbase import SB
from parsers.sotwe_parser import parse_sotwe_html_tweets


class SotweSource:
    def __init__(self, name: str = "sotwe", priority: int = 20, timeout: int = 25):
        self.name = name
        self.priority = priority
        self.timeout = timeout

    def is_available(self) -> bool:
        return True

    def fetch_account(
        self, account: str, limit: int = 50, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        """SeleniumBase UC ModeでSotweページを開き、DOM要素描画を待機して回収"""
        url = f"https://www.sotwe.com/{account}"
        if log_fn:
            log_fn(f"[SotweSource:SB] Opening Web UI with UC Mode: {url}")

        try:
            with SB(uc=True, headless=True) as sb:
                sb.open(url)
                # .tweet-card がDOM上に現れるまで最大10秒待機
                try:
                    sb.wait_for_element_present("div.tweet-card", timeout=10)
                except Exception:
                    time.sleep(3)

                # スクロールして下位要素・画像をトリガー
                sb.execute_script("window.scrollTo(0, 1500);")
                time.sleep(2)
                html_source = sb.get_page_source()

            records = parse_sotwe_html_tweets(html_source, account)
            if limit > 0:
                records = records[:limit]

            if log_fn:
                log_fn(f"[SotweSource:SB] Successfully parsed {len(records)} posts from DOM.")
            return records
        except Exception as e:
            if log_fn:
                log_fn(f"[SotweSource:SB] Error during extraction: {e}")
            return []

    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        target_account = account or "i"
        posts = self.fetch_account(target_account, limit=20, log_fn=log_fn)
        for p in posts:
            if p.get("post", {}).get("id") == post_id:
                return p
        return None
