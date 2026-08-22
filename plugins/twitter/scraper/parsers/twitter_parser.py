# plugins/twitter/scraper/parsers/twitter_parser.py (100行以下)
import datetime, json, re
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
        t_id, u_name = str(tweet.get("id_str") or tweet.get("id") or ""), user.get("screen_name") or user.get("username") or ""
        if not t_id or not u_name:
            det = self.detect_platform_and_account(uri)
            if det: u_name, t_id = u_name or det["account"], t_id or det.get("status_id", "")
        if not t_id: return None

        media_list = [{"url": m.get("media_url_https") or m.get("media_url") or m.get("url"), "type": m.get("type", "image"),
                       "width": m.get("sizes", {}).get("large", {}).get("w", 0), "height": m.get("sizes", {}).get("large", {}).get("h", 0)}
                      for m in (tweet.get("extended_entities", {}) or tweet.get("entities", {})).get("media", []) if m.get("media_url_https") or m.get("media_url") or m.get("url")]
        urls_list = [{"short_url": u.get("url"), "expanded_url": u.get("expanded_url") or u.get("unwound", {}).get("url") or u.get("url")}
                     for u in (tweet.get("extended_entities", {}) or tweet.get("entities", {})).get("urls", []) if u.get("url")]

        return {
            "platform": "twitter",
            "account": {"numeric_id": str(user.get("id_str") or user.get("id") or u_name), "username": u_name,
                        "display_name": user.get("name") or u_name, "avatar_url": user.get("profile_image_url_https") or user.get("profile_image_url") or ""},
            "post": {"id": t_id, "conversation_id": str(tweet.get("conversation_id_str") or t_id), "reply_to_tweet_id": tweet.get("in_reply_to_status_id_str"),
                     "reply_to_handle": tweet.get("in_reply_to_screen_name"), "created_at": tweet.get("created_at") or "",
                     "full_text": tweet.get("full_text") or tweet.get("text") or "", "wayback_url": uri, "urls": urls_list},
            "media": media_list,
        }

    def _parse_html(self, html_text: str, uri: str) -> Optional[Dict[str, Any]]:
        det = self.detect_platform_and_account(uri)
        if not det: return None
        p_id = det.get("status_id") or (re.search(r'status(?:es)?/(\d+)', html_text) or re.search(r'data-tweet-id="(\d+)"', html_text) or [None, ""])[1]
        if not p_id: return None

        m_desc = re.search(r'<meta property="og:description" content="([^"]+)"', html_text) or re.search(r'<title>([^<]+)</title>', html_text)
        txt = m_desc.group(1) if m_desc else f"Archived post from {uri}"

        m_av = re.search(r'<img[^>]+(?:class="[^"]*(?:avatar|ProfileAvatar-image)[^"]*"|src="([^"]*profile_images[^"]*)")[^>]*src="?([^" >]+)', html_text)
        avatar_url = (m_av.group(2) if m_av and m_av.group(2) else (m_av.group(1) if m_av else "")) or ""
        m_name = re.search(r'<strong class="fullname[^"]*">([^<]+)</strong>', html_text) or re.search(r'<meta property="og:title" content="([^"（(]+)', html_text)
        display_name = m_name.group(1).strip() if m_name else det["account"]

        media_list, seen = [], set()
        for pat, m_type in [(r'(https?://pbs\.twimg\.com/media/[a-zA-Z0-9_\-]+(?:\.[a-zA-Z0-9]+|\?[^"\'\s<>]*)?)', "image"),
                            (r'(https?://video\.twimg\.com/[^"\'\s<>]+\.(?:mp4|m3u8|webm))', "video")]:
            for u in re.findall(pat, html_text):
                u_c = u.replace("&amp;", "&"); base = u_c.split("?")[0].split(":")[0]
                if base not in seen: seen.add(base); media_list.append({"url": u_c, "type": m_type, "width": 0, "height": 0})

        if not media_list:
            for u in re.findall(r'<meta property="og:image" content="([^"]+)"', html_text):
                if "profile_images" not in u and u not in seen: seen.add(u); media_list.append({"url": u.replace("&amp;", "&"), "type": "image", "width": 0, "height": 0})

        created_at = ""
        m_time = re.search(r'data-time="(\d+)"', html_text)
        if m_time:
            try: created_at = datetime.datetime.fromtimestamp(int(m_time.group(1)), tz=datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
            except Exception: pass
        if not created_at and p_id.isdigit() and int(p_id) > 30000000000:
            try: created_at = datetime.datetime.fromtimestamp(((int(p_id) >> 22) + 1288834974657) / 1000, tz=datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
            except Exception: pass
        if not created_at:
            m_wb = re.search(r'/web/(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})(\d{2})', uri)
            if m_wb: created_at = f"{m_wb.group(1)}-{m_wb.group(2)}-{m_wb.group(3)} {m_wb.group(4)}:{m_wb.group(5)}:{m_wb.group(6)}"

        return {
            "platform": "twitter",
            "account": {"numeric_id": det["account"], "username": det["account"], "display_name": display_name, "avatar_url": avatar_url},
            "post": {"id": p_id, "conversation_id": p_id, "created_at": created_at, "full_text": txt, "wayback_url": uri, "urls": []},
            "media": media_list,
        }
