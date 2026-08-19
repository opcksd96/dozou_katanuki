# plugins/twitter/scraper/parsers/twitter_parser.py (100行以下)
import json
import re
from typing import Any, Dict, List, Optional
from .base_parser import BaseParser


class TwitterParser(BaseParser):
    """Twitter / X 特化型抽出エンジン (SPEC-PLUGIN-001)"""

    URL_PATTERN = re.compile(r"(?:https?://)?(?:twitter\.com|x\.com)/([a-zA-Z0-9_]+)(?:/status(?:es)?/(\d+))?")

    def detect_platform_and_account(self, uri: str) -> Optional[Dict[str, str]]:
        m = self.URL_PATTERN.search(uri)
        if m:
            account = m.group(1)
            status_id = m.group(2) or ""
            return {"platform": "twitter", "account": account, "status_id": status_id}
        return None

    def parse_record(self, raw_data: Any, uri: str) -> Optional[Dict[str, Any]]:
        if isinstance(raw_data, (str, bytes)):
            try:
                data = json.loads(raw_data)
                if isinstance(data, dict):
                    return self._parse_json(data, uri)
            except Exception:
                pass
            return self._parse_html(raw_data.decode("utf-8", errors="ignore") if isinstance(raw_data, bytes) else str(raw_data), uri)
        elif isinstance(raw_data, dict):
            return self._parse_json(raw_data, uri)
        return None

    def _parse_json(self, data: Dict[str, Any], uri: str) -> Optional[Dict[str, Any]]:
        tweet = data.get("tweet") or data.get("data") or data
        user = data.get("user") or (data.get("includes", {}).get("users", [{}])[0] if "includes" in data else {})
        tweet_id = str(tweet.get("id_str") or tweet.get("id") or "")
        username = user.get("screen_name") or user.get("username") or ""

        if not tweet_id or not username:
            detect = self.detect_platform_and_account(uri)
            if detect:
                username = username or detect["account"]
                tweet_id = tweet_id or detect.get("status_id", "")
        if not tweet_id:
            return None

        media_list: List[Dict[str, Any]] = []
        entities = tweet.get("extended_entities", {}) or tweet.get("entities", {})
        for m in entities.get("media", []):
            m_url = m.get("media_url_https") or m.get("media_url") or m.get("url")
            if m_url:
                media_list.append({
                    "url": m_url,
                    "type": m.get("type", "image"),
                    "width": m.get("sizes", {}).get("large", {}).get("w", 0),
                    "height": m.get("sizes", {}).get("large", {}).get("h", 0),
                })

        return {
            "platform": "twitter",
            "account": {
                "numeric_id": str(user.get("id_str") or user.get("id") or username),
                "username": username,
                "display_name": user.get("name") or username,
                "avatar_url": user.get("profile_image_url_https") or user.get("profile_image_url") or "",
            },
            "post": {
                "id": tweet_id,
                "conversation_id": str(tweet.get("conversation_id_str") or tweet_id),
                "reply_to_tweet_id": tweet.get("in_reply_to_status_id_str"),
                "reply_to_handle": tweet.get("in_reply_to_screen_name"),
                "created_at": tweet.get("created_at") or "",
                "full_text": tweet.get("full_text") or tweet.get("text") or "",
                "wayback_url": uri,
            },
            "media": media_list,
        }

    def _parse_html(self, html_text: str, uri: str) -> Optional[Dict[str, Any]]:
        detect = self.detect_platform_and_account(uri)
        if not detect:
            return None
        post_id = detect.get("status_id")
        if not post_id:
            id_m = re.search(r'status(?:es)?/(\d+)', html_text) or re.search(r'data-tweet-id="(\d+)"', html_text)
            post_id = id_m.group(1) if id_m else ""
        if not post_id:
            return None

        text_m = re.search(r'<meta property="og:description" content="([^"]+)"', html_text) or re.search(r'<title>([^<]+)</title>', html_text)
        full_text = text_m.group(1) if text_m else f"Archived post from {uri}"
        return {
            "platform": "twitter",
            "account": {"numeric_id": detect["account"], "username": detect["account"], "display_name": detect["account"], "avatar_url": ""},
            "post": {"id": post_id, "conversation_id": post_id, "full_text": full_text, "wayback_url": uri},
            "media": [],
        }
