# plugins/twitter/scraper/parsers/twistalker_parser.py (SPEC-PLUGIN-001 / 100行以下)
import re
from typing import Any, Dict, List
from bs4 import BeautifulSoup
try:
    from plugins.twitter.scraper.parsers.twistalker_extractors import clean_html_text, parse_metric_val, parse_twitter_datetime
    from plugins.twitter.scraper.parsers.twistalker_media import extract_media_from_api, extract_media_from_html
except ImportError:
    from parsers.twistalker_extractors import clean_html_text, parse_metric_val, parse_twitter_datetime
    from parsers.twistalker_media import extract_media_from_api, extract_media_from_html

def parse_twistalker_api_tweets(tweets_data: Any, default_account: str) -> List[Dict[str, Any]]:
    results = []
    tweets_list = tweets_data.values() if isinstance(tweets_data, dict) else (tweets_data or [])
    for tw in tweets_list:
        p_id = str(tw.get("conversation_id_str") or tw.get("id_str") or "")
        core = tw.get("core") or {}
        u_name = core.get("screen_name") or default_account
        created_at = parse_twitter_datetime(tw.get("created_at"), p_id)
        media = extract_media_from_api(tw.get("extended_entities") or {})
        results.append({
            "platform": "twitter", "source_name": "twistalker",
            "account": {
                "numeric_id": str(core.get("rest_id") or f"ext_{u_name}"),
                "username": u_name, "display_name": core.get("name") or u_name,
                "avatar_url": core.get("profile_image_url_https", ""),
                "avatar_original_url": (core.get("profile_image_url_https") or "").replace("_normal.", "."),
                "description": core.get("description") or "", "profile_history": []
            },
            "post": {
                "id": p_id, "conversation_id": p_id, "reply_to_tweet_id": None, "reply_to_handle": None,
                "created_at": created_at, "full_text": clean_html_text(tw.get("full_text", "")),
                "via": "TwStalker", "source_name": "twistalker", "source_domain": "twstalker.com",
                "is_repost": bool(tw.get("is_retweet")), "is_pinned": False,
                "retweeted_by": default_account if tw.get("is_retweet") else "", "wayback_url": "",
                "original_url": f"https://x.com/{u_name}/status/{p_id}" if p_id else "",
                "metrics": {
                    "replies": int(tw.get("reply_count") or 0), "likes": int(tw.get("favorite_count") or 0),
                    "retweets": int(tw.get("retweet_count") or 0), "bookmarks": int(tw.get("bookmark_count") or 0),
                    "views": int(tw.get("view_count") or 0)
                }, "urls": []
            }, "media": media
        })
    return results

def parse_twistalker_html_tweets(html_str: str, default_account: str) -> List[Dict[str, Any]]:
    if not html_str: return []
    soup = BeautifulSoup(html_str, "html.parser")
    cards = soup.select(".activity-posts, .post")
    results = []
    for card in cards:
        status_link = card.select_one(".user-text3 span a[href*='/status/'], a[href*='/status/']")
        href = status_link.get("href", "") if status_link else ""
        m = re.search(r'/status(?:es)?/(\d+)', href) or re.search(r'data-id="(\d+)"', str(card))
        p_id = m.group(1) if m else card.get("data-id", "")
        if not p_id: continue

        user_link = card.select_one(".main-user-dts1 a[href^='/'], .username")
        u_name = user_link.get_text(strip=True).lstrip("@") if user_link and "@" in user_link.get_text() else (user_link.get("href", "").strip("/").split("/")[0] if user_link and user_link.get("href") else default_account)
        name_el = card.select_one(".user-text3 h4, .fullname")
        d_name = name_el.get_text(strip=True) if name_el else u_name
        av_img = card.select_one(".main-user-dts1 img[src], img.avatar, img")
        av_url = av_img.get("src", "") if av_img and "profile_images" in av_img.get("src", "") else ""

        desc_el = card.select_one(".activity-descp p, .post-text, .tweet-text")
        full_text = desc_el.get_text("\n", strip=True) if desc_el else ""
        created_at = parse_twitter_datetime(status_link.get_text(strip=True) if status_link else "", p_id)
        is_rt = bool(card.select_one(".fa-retweet"))
        media = extract_media_from_html(card)

        metrics = {"replies": 0, "likes": 0, "retweets": 0, "bookmarks": 0, "views": 0}
        for it in card.select(".like-comment-view .left-comments a.like-item"):
            val = parse_metric_val(it.select_one("span").get_text(strip=True) if it.select_one("span") else "")
            if it.select_one(".fa-comment"): metrics["replies"] = val
            elif it.select_one(".fa-retweet"): metrics["retweets"] = val
            elif it.select_one(".fa-heart"): metrics["likes"] = val
            elif it.select_one(".fa-chart-simple"): metrics["views"] = val
            elif it.select_one(".fa-bookmark"): metrics["bookmarks"] = val

        results.append({
            "platform": "twitter", "source_name": "twistalker",
            "account": {
                "numeric_id": f"ext_{u_name}", "username": u_name, "display_name": d_name,
                "avatar_url": av_url, "avatar_original_url": av_url.replace("_normal.", "."),
                "description": "", "profile_history": []
            },
            "post": {
                "id": p_id, "conversation_id": p_id, "reply_to_tweet_id": None, "reply_to_handle": None,
                "created_at": created_at, "full_text": full_text, "via": "TwStalker", "source_name": "twistalker",
                "source_domain": "twstalker.com", "is_repost": is_rt, "is_pinned": False,
                "retweeted_by": default_account if is_rt else "", "wayback_url": "",
                "original_url": f"https://x.com/{u_name}/status/{p_id}", "metrics": metrics, "urls": []
            }, "media": media
        })
    return results
