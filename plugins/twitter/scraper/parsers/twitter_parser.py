# plugins/twitter/scraper/parsers/twitter_parser.py (100行以下)
import json, re
from typing import Any, Dict, List, Optional
from .base_parser import BaseParser


class TwitterParser(BaseParser):
    """Twitter / X 特化型抽出エンジン (SPEC-PLUGIN-001)"""
    URL_PATTERN = re.compile(r"(?:https?://)?(?:[a-zA-Z0-9_.\-]+\.)?(?:twitter\.com|x\.com)/([a-zA-Z0-9_]+)(?:/status(?:es)?/(\d+))?")
    RESERVED = {"i", "search", "home", "explore", "settings", "intent", "hashtag", "share", "api", "1", "oauth", "account"}

    def detect_platform_and_account(self, uri: str) -> Optional[Dict[str, str]]:
        m = self.URL_PATTERN.search(uri)
        if m and m.group(1).lower() not in self.RESERVED:
            return {"platform": "twitter", "account": m.group(1), "status_id": m.group(2) or ""}
        return None

    def parse_record(self, raw_data: Any, uri: str) -> Optional[Dict[str, Any]]:
        if isinstance(raw_data, (str, bytes)):
            try:
                data = json.loads(raw_data)
                if isinstance(data, dict): return self._parse_json(data, uri)
            except Exception: pass
            html = raw_data.decode("utf-8", errors="ignore") if isinstance(raw_data, bytes) else str(raw_data)
            return self._parse_html(html, uri)
        elif isinstance(raw_data, dict):
            return self._parse_json(raw_data, uri)
        return None

    def _parse_json(self, data: Dict[str, Any], uri: str) -> Optional[Dict[str, Any]]:
        tweet = data.get("tweet") or data.get("data") or data
        user = data.get("user") or (data.get("includes", {}).get("users", [{}])[0] if "includes" in data else {})
        t_id = str(tweet.get("id_str") or tweet.get("id") or "")
        u_name = user.get("screen_name") or user.get("username") or ""

        if not t_id or not u_name:
            det = self.detect_platform_and_account(uri)
            if det:
                u_name, t_id = u_name or det["account"], t_id or det.get("status_id", "")
        if not t_id: return None

        media_list: List[Dict[str, Any]] = []
        entities = tweet.get("extended_entities", {}) or tweet.get("entities", {})
        for m in entities.get("media", []):
            url = m.get("media_url_https") or m.get("media_url") or m.get("url")
            if url:
                sizes = m.get("sizes", {}).get("large", {})
                media_list.append({"url": url, "type": m.get("type", "image"),
                                   "width": sizes.get("w", 0), "height": sizes.get("h", 0)})

        urls_list = []
        for u in entities.get("urls", []):
            s_url = u.get("url")
            e_url = u.get("expanded_url") or u.get("unwound", {}).get("url") or s_url
            if s_url and e_url: urls_list.append({"short_url": s_url, "expanded_url": e_url})

        return {
            "platform": "twitter",
            "account": {
                "numeric_id": str(user.get("id_str") or user.get("id") or u_name),
                "username": u_name,
                "display_name": user.get("name") or u_name,
                "avatar_url": user.get("profile_image_url_https") or user.get("profile_image_url") or "",
            },
            "post": {
                "id": t_id,
                "conversation_id": str(tweet.get("conversation_id_str") or t_id),
                "reply_to_tweet_id": tweet.get("in_reply_to_status_id_str"),
                "reply_to_handle": tweet.get("in_reply_to_screen_name"),
                "created_at": tweet.get("created_at") or "",
                "full_text": tweet.get("full_text") or tweet.get("text") or "",
                "wayback_url": uri,
                "urls": urls_list,
            },
            "media": media_list,
        }

    def _parse_html(self, html_text: str, uri: str) -> Optional[Dict[str, Any]]:
        det = self.detect_platform_and_account(uri)
        if not det: return None
        p_id = det.get("status_id") or (re.search(r'status(?:es)?/(\d+)', html_text) or re.search(r'data-tweet-id="(\d+)"', html_text) or [None, ""])[1]
        if not p_id: return None

        m_desc = re.search(r'<meta property="og:description" content="([^"]+)"', html_text) or re.search(r'<title>([^<]+)</title>', html_text)
        txt = m_desc.group(1) if m_desc else f"Archived post from {uri}"
        return {
            "platform": "twitter",
            "account": {"numeric_id": det["account"], "username": det["account"], "display_name": det["account"], "avatar_url": ""},
            "post": {"id": p_id, "conversation_id": p_id, "full_text": txt, "wayback_url": uri, "urls": []},
            "media": [],
        }
