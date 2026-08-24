# plugins/twitter/scraper/parsers/twitter_parser.py (100行以下)
import datetime, html, json, re
from typing import Any, Dict, List, Optional
from .base_parser import BaseParser


class TwitterParser(BaseParser):
    """Twitter / X 特化型抽出エンジン（公式JSON / Sotwe / Twistalker / Nitter / 魚拓対応）"""
    URL_PATTERN = re.compile(r"(?:https?://)?(?:[a-zA-Z0-9_.\-]+\.)?(?:twitter\.com|x\.com|sotwe\.com|twistalker\.com|nitter\.[a-z0-9.\-]+)/([a-zA-Z0-9_]+)(?:/status(?:es)?/(\d+))?")
    RESERVED = {"i", "search", "home", "explore", "settings", "intent", "hashtag", "share", "api", "1", "oauth", "account"}

    def detect_platform_and_account(self, uri: str) -> Optional[Dict[str, str]]:
        m = self.URL_PATTERN.search(uri)
        return {"platform": "twitter", "account": m.group(1), "status_id": m.group(2) or ""} if (m and m.group(1).lower() not in self.RESERVED) else None

    def parse_record(self, raw_data: Any, uri: str) -> Optional[Dict[str, Any]]:
        if isinstance(raw_data, dict) and "post" in raw_data and "account" in raw_data:
            return raw_data
        if isinstance(raw_data, (str, bytes)):
            try:
                data = json.loads(raw_data)
                if isinstance(data, dict): return self._parse_json(data, uri)
            except Exception: pass
            html_text = raw_data.decode("utf-8", errors="ignore") if isinstance(raw_data, bytes) else str(raw_data)
            return self._parse_html(html_text, uri)
        elif isinstance(raw_data, dict): return self._parse_json(raw_data, uri)
        return None

    def _parse_json(self, data: Dict[str, Any], uri: str) -> Optional[Dict[str, Any]]:
        tweet = data.get("tweet") or data.get("data") or data
        user = data.get("user") or (data.get("includes", {}).get("users", [{}])[0] if "includes" in data else {}) or tweet.get("user", {})
        t_id, u_name = str(tweet.get("id_str") or tweet.get("id") or ""), user.get("screen_name") or user.get("username") or user.get("name") or ""
        if not t_id or not u_name:
            det = self.detect_platform_and_account(uri)
            if det: u_name, t_id = u_name or det["account"], t_id or det.get("status_id", "")
        if not t_id: return None
        raw_media = (tweet.get("extended_entities", {}) or tweet.get("entities", {})).get("media", []) or tweet.get("mediaEntities", []) or tweet.get("photos", []) or []
        media_list = [{"url": m.get("media_url_https") or m.get("media_url") or m.get("url"), "type": m.get("type", "image"),
                       "width": m.get("sizes", {}).get("large", {}).get("w", 0), "height": m.get("sizes", {}).get("large", {}).get("h", 0)}
                      for m in raw_media if m.get("media_url_https") or m.get("media_url") or m.get("url")]
        raw_urls = (tweet.get("extended_entities", {}) or tweet.get("entities", {})).get("urls", []) or tweet.get("urls", [])
        urls_list = [{"short_url": u.get("url"), "expanded_url": u.get("expanded_url") or u.get("unwound", {}).get("url") or u.get("url")}
                     for u in raw_urls if u.get("url")]
        return {
            "platform": "twitter",
            "account": {"numeric_id": str(user.get("id_str") or user.get("id") or u_name), "username": u_name,
                        "display_name": user.get("name") or user.get("screen_name") or u_name,
                        "avatar_url": user.get("profile_image_url_https") or user.get("profile_image_url") or user.get("avatar") or ""},
            "post": {"id": t_id, "conversation_id": str(tweet.get("conversation_id_str") or tweet.get("conversation_id") or t_id),
                     "reply_to_tweet_id": tweet.get("in_reply_to_status_id_str") or tweet.get("in_reply_to_status_id"),
                     "reply_to_handle": tweet.get("in_reply_to_screen_name") or tweet.get("in_reply_to_username"),
                     "created_at": tweet.get("created_at") or tweet.get("createdAt") or "",
                     "full_text": tweet.get("full_text") or tweet.get("text") or "", "wayback_url": uri, "urls": urls_list},
            "media": media_list,
        }

    def _meta(self, html_text: str, prop: str) -> Optional[str]:
        for tag in re.findall(r'<meta\s+[^>]+>', html_text, re.IGNORECASE):
            if prop.lower() in tag.lower() and (m := re.search(r'content=["\'](.*?)["\']', tag, re.IGNORECASE | re.DOTALL)):
                return html.unescape(m.group(1).strip('“”" '))
        return None

    def _parse_html(self, html_text: str, uri: str) -> Optional[Dict[str, Any]]:
        det = self.detect_platform_and_account(uri)
        p_id = (det.get("status_id") if det else "") or (re.search(r'status(?:es)?/(\d+)', html_text) or re.search(r'data-id="(\d+)"', html_text) or re.search(r'data-tweet-id="(\d+)"', html_text) or [None, ""])[1]
        u_name = (det.get("account") if det else "") or (re.search(r'class="username[^"]*"[^>]*>@?([a-zA-Z0-9_]+)<', html_text) or [None, "unknown"])[1]
        if not p_id: return None
        txt = self._meta(html_text, "og:description") or self._meta(html_text, "twitter:description") or (re.search(r'<div[^>]+class="[^"]*(?:tweet-content|post-text|tweet-text)[^"]*"[^>]*>(.*?)</div>', html_text, re.DOTALL) or [None, ""])[1]
        if txt and "<" in txt: txt = re.sub(r'<[^>]+>', '', txt).strip()
        if not txt:
            m_t = re.search(r'<title[^>]*>(.*?)</title>', html_text, re.IGNORECASE | re.DOTALL)
            txt = html.unescape(m_t.group(1)).split(" on Twitter:")[ -1].split(" / Twitter")[0].strip('“”" ') if m_t else f"Archived post {p_id}"
        m_av = re.search(r'<img[^>]+(?:class="[^"]*(?:avatar|ProfileAvatar-image)[^"]*"|src="([^"]*(?:profile_images|avatar)[^"]*)")[^>]*src="?([^" >]+)', html_text)
        avatar_url = (m_av.group(2) if m_av and m_av.group(2) else (m_av.group(1) if m_av else "")) or ""
        display_name = self._meta(html_text, "og:title") or (re.search(r'<strong class="fullname[^"]*">([^<]+)</strong>', html_text) or (re.search(r'class="fullname[^"]*"[^>]*>([^<]+)<', html_text)) or [None, u_name])[1]
        if " on Twitter" in display_name: display_name = display_name.split(" on Twitter")[0].strip()

        media_list, seen = [], set()
        for pat, m_type in [(r'(https?://pbs\.twimg\.com/media/[a-zA-Z0-9_\-]+(?:\.[a-zA-Z0-9]+|\?[^"\'\s<>]*)?)', "image"),
                            (r'(https?://(?:video|video-cdn)\.twimg\.com/[^"\'\s<>]+\.(?:mp4|m3u8|webm))', "video")]:
            for u in re.findall(pat, html_text):
                u_c = u.replace("&amp;", "&"); base = re.sub(r':(large|orig|small|medium|thumb)$', '', u_c.split('?')[0])
                if base not in seen: seen.add(base); media_list.append({"url": u_c, "type": m_type, "width": 0, "height": 0})
        og_img = self._meta(html_text, "og:image")
        if og_img and "profile_images" not in og_img and not any(og_img.split("?")[0] in m["url"] for m in media_list): media_list.append({"url": og_img, "type": "image", "width": 0, "height": 0})

        created_at = self._meta(html_text, "datePublished") or ""
        m_time = re.search(r'data-time="(\d+)"', html_text)
        if not created_at and m_time:
            try: created_at = datetime.datetime.fromtimestamp(int(m_time.group(1)), tz=datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
            except Exception: pass
        if not created_at and p_id.isdigit() and int(p_id) > 30000000000:
            try: created_at = datetime.datetime.fromtimestamp(((int(p_id) >> 22) + 1288834974657) / 1000, tz=datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
            except Exception: pass
        return {
            "platform": "twitter", "account": {"numeric_id": u_name, "username": u_name, "display_name": display_name, "avatar_url": avatar_url},
            "post": {"id": p_id, "conversation_id": p_id, "created_at": created_at, "full_text": txt, "wayback_url": uri, "urls": []}, "media": media_list,
        }
