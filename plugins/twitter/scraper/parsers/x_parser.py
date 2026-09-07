# plugins/twitter/scraper/parsers/x_parser.py (SPEC-PLUGIN-001 / 100行以下)
import re
from typing import Any, Dict, Optional
from bs4 import BeautifulSoup
try:
    from plugins.twitter.scraper.parsers.x_extractors import snowflake_to_datetime, parse_metric_number
    from plugins.twitter.scraper.parsers.x_media import extract_media_from_syndication, extract_media_from_html_card
except ImportError:
    from parsers.x_extractors import snowflake_to_datetime, parse_metric_number
    from parsers.x_media import extract_media_from_syndication, extract_media_from_html_card

class XParser:
    """X / Twitter (x.com) 独自パーサー"""

    def parse_syndication_json(self, data: Dict[str, Any], default_acc: str = "") -> Optional[Dict[str, Any]]:
        if not data or not isinstance(data, dict): return None
        t_id = str(data.get("id_str") or data.get("id") or "")
        if not t_id: return None
        u_info = data.get("user") or {}
        u_name = u_info.get("screen_name") or default_acc or "unknown"
        d_name = u_info.get("name") or u_name
        av_url = u_info.get("profile_image_url_https") or ""
        created_at = snowflake_to_datetime(t_id) or str(data.get("created_at") or "")
        f_text = data.get("text") or ""
        media_list = extract_media_from_syndication(data.get("mediaDetails") or [])
        return {
            "platform": "twitter", "source_name": "x_official",
            "account": {"numeric_id": str(u_info.get("id_str") or f"ext_{u_name}"), "username": u_name,
                        "display_name": d_name, "avatar_url": av_url, "avatar_original_url": av_url.replace("_normal.", "."),
                        "description": u_info.get("description", ""), "profile_history": []},
            "post": {"id": t_id, "conversation_id": str(data.get("conversation_id_str") or t_id),
                     "reply_to_tweet_id": data.get("in_reply_to_status_id_str"), "reply_to_handle": data.get("in_reply_to_screen_name"),
                     "created_at": created_at, "full_text": f_text, "via": "X", "source_name": "x_official",
                     "source_domain": "x.com", "is_repost": bool(data.get("retweeted_status")), "is_pinned": False,
                     "retweeted_by": "", "wayback_url": "", "sotwe_url": "", "original_url": f"https://x.com/{u_name}/status/{t_id}",
                     "metrics": {"replies": parse_metric_number(data.get("reply_count")),
                                 "likes": parse_metric_number(data.get("favorite_count")),
                                 "retweets": parse_metric_number(data.get("retweet_count")),
                                 "bookmarks": 0, "views": 0}, "urls": []},
            "media": media_list
        }

    def parse_html_card(self, card_soup: BeautifulSoup, default_acc: str = "") -> Optional[Dict[str, Any]]:
        status_a = card_soup.find("a", href=lambda h: h and "/status/" in h)
        if not status_a: return None
        m_id = re.search(r"/status/(\d+)", status_a["href"])
        if not m_id: return None
        t_id = m_id.group(1)
        u_name = default_acc or "unknown"
        for a in card_soup.find_all("a", href=True):
            h = a["href"].lstrip("/")
            if "/" not in h and h.lower() not in ["home", "explore", "search", "i"]:
                u_name = h; break
        d_name_el = card_soup.find("a", href=lambda h: h and u_name in h and "underline" in str(card_soup))
        d_name = d_name_el.get_text().strip() if d_name_el else u_name
        av_img = card_soup.find("img", src=lambda s: s and "profile_images" in s)
        av_url = av_img["src"] if av_img else ""
        is_pinned = bool(card_soup.find("svg", attrs={"data-icon": "icon-pin-fill"}))
        text_el = card_soup.find("div", attrs={"dir": "auto", "class": lambda c: c and "font-chirp" in c})
        f_text = text_el.get_text().strip() if text_el else ""
        created_at = snowflake_to_datetime(t_id) or ""
        metrics = {"replies": 0, "likes": 0, "retweets": 0, "bookmarks": 0, "views": 0}
        for btn in card_soup.find_all(["button", "a"], attrs={"aria-label": True}):
            lbl = btn["aria-label"]
            val_el = btn.find("span", attrs={"data-animated-count-visual": "true"}) or btn.find("span", class_=lambda c: c and "tabular-nums" in c)
            val = parse_metric_number(val_el.get_text() if val_el else 0)
            if "返信" in lbl or "reply" in lbl.lower(): metrics["replies"] = val
            elif "リポスト" in lbl or "retweet" in lbl.lower(): metrics["retweets"] = val
            elif "いいね" in lbl or "like" in lbl.lower(): metrics["likes"] = val
            elif "ビュー" in lbl or "view" in lbl.lower(): metrics["views"] = val
            elif "ブックマーク" in lbl or "bookmark" in lbl.lower(): metrics["bookmarks"] = val
        return {
            "platform": "twitter", "source_name": "x_official",
            "account": {"numeric_id": f"ext_{u_name}", "username": u_name, "display_name": d_name,
                        "avatar_url": av_url, "avatar_original_url": av_url.replace("_normal.", "."),
                        "description": "", "profile_history": []},
            "post": {"id": t_id, "conversation_id": t_id, "reply_to_tweet_id": None, "reply_to_handle": None,
                     "created_at": created_at, "full_text": f_text, "via": "X", "source_name": "x_official",
                     "source_domain": "x.com", "is_repost": False, "is_pinned": is_pinned,
                     "retweeted_by": "", "wayback_url": "", "sotwe_url": "", "original_url": f"https://x.com/{u_name}/status/{t_id}",
                     "metrics": metrics, "urls": []},
            "media": extract_media_from_html_card(card_soup)
        }
