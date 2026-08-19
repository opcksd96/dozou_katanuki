# plugins/twitter/scraper/parsers/twitter_parser.py (100行以下)
import json
import re
from typing import Any, Dict, List, Optional
from .base_parser import BaseParser


class TwitterParser(BaseParser):
    """Twitter / X 特化型抽出エンジン (SPEC-PLUGIN-001)"""

    URL_PATTERN = re.compile(
        r"(?:https?://)?(?:twitter\.com|x\.com)/([a-zA-Z0-9_]+)(?:/status/(\d+))?"
    )

    def detect_platform_and_account(self, uri: str) -> Optional[Dict[str, str]]:
        m = self.URL_PATTERN.search(uri)
        if m:
            return {"platform": "twitter", "account": m.group(1), "status_id": m.group(2) or ""}
        return None

    def parse_record(self, raw_data: Any, uri: str) -> Optional[Dict[str, Any]]:
        if isinstance(raw_data, (str, bytes)):
            try:
                data = json.loads(raw_data)
            except Exception:
                return self._parse_html(str(raw_data), uri)
        else:
            data = raw_data

        if isinstance(data, dict):
            return self._parse_json(data, uri)
        return None

    def _parse_json(self, data: Dict[str, Any], uri: str) -> Optional[Dict[str, Any]]:
        # Twitter API / Syndication / Tweet Result 互換構造の抽出
        tweet = data.get("tweet") or data.get("data") or data
        user = data.get("user") or data.get("includes", {}).get("users", [{}])[0] or {}

        tweet_id = str(tweet.get("id_str") or tweet.get("id") or "")
        username = user.get("screen_name") or user.get("username") or ""

        if not tweet_id or not username:
            detect = self.detect_platform_and_account(uri)
            if detect:
                username = username or detect["account"]
                tweet_id = tweet_id or detect.get("status_id", "")

        if not tweet_id:
            return None

        # メディア抽出
        media_list: List[Dict[str, Any]] = []
        extended_entities = tweet.get("extended_entities", {}) or tweet.get("entities", {})
        for m in extended_entities.get("media", []):
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
                "numeric_id": str(user.get("id_str") or user.get("id") or ""),
                "username": username,
                "display_name": user.get("name") or username,
                "avatar_url": user.get("profile_image_url_https") or user.get("profile_image_url") or "",
                "description": user.get("description") or "",
            },
            "post": {
                "id": tweet_id,
                "conversation_id": str(tweet.get("conversation_id_str") or tweet_id),
                "reply_to_tweet_id": tweet.get("in_reply_to_status_id_str"),
                "created_at": tweet.get("created_at") or "",
                "full_text": tweet.get("full_text") or tweet.get("text") or "",
                "wayback_url": uri,
            },
            "media": media_list,
        }

    def _parse_html(self, html_text: str, uri: str) -> Optional[Dict[str, Any]]:
        # 簡易 HTML 抽出フォールバック
        detect = self.detect_platform_and_account(uri)
        if not detect or not detect.get("status_id"):
            return None
        return {
            "platform": "twitter",
            "account": {
                "numeric_id": "",
                "username": detect["account"],
                "display_name": detect["account"],
                "avatar_url": "",
                "description": "",
            },
            "post": {
                "id": detect["status_id"],
                "conversation_id": detect["status_id"],
                "reply_to_tweet_id": None,
                "created_at": "",
                "full_text": f"Archived post from {uri}",
                "wayback_url": uri,
            },
            "media": [],
        }
