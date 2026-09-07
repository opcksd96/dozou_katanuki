# plugins/twitter/scraper/parsers/nitter_parser.py (SPEC-PLUGIN-001 / 100行以下)
import re
from typing import Any, Dict, List, Optional
from bs4 import BeautifulSoup
try:
    from plugins.twitter.scraper.parsers.nitter_extractors import clean_html_text, parse_metric_val, parse_nitter_datetime
    from plugins.twitter.scraper.parsers.nitter_media import extract_media_from_nitter_card
except ImportError:
    from parsers.nitter_extractors import clean_html_text, parse_metric_val, parse_nitter_datetime
    from parsers.nitter_media import extract_media_from_nitter_card

def parse_nitter_card(card, default_account: str = "", fallback_id: str = "") -> Optional[Dict[str, Any]]:
    link_el = card.select_one("a.tweet-link") or card.select_one(".tweet-date a")
    href = link_el.get("href", "") if link_el else ""
    m = re.search(r'/status/(\d+)', href)
    p_id = m.group(1) if m else (fallback_id or "")
    if not p_id: return None

    user_el = card.select_one(".tweet-header .username")
    u_name = user_el.get_text(strip=True).lstrip("@") if user_el else default_account
    name_el = card.select_one(".tweet-header .fullname")
    d_name = name_el.get_text(strip=True) if name_el else u_name
    av_img = card.select_one(".tweet-avatar img")
    av_url = av_img.get("src", "") if av_img else ""

    content_el = card.select_one(".tweet-content")
    full_text = clean_html_text(content_el.decode_contents() if content_el else "")
    date_el = card.select_one(".tweet-date a")
    raw_date = date_el.get("title", "") if date_el else ""
    created_at = parse_nitter_datetime(raw_date, p_id)

    is_pinned = bool(card.select_one(".pinned"))
    rt_header = card.select_one(".retweet-header")
    is_rt = bool(rt_header)
    retweeted_by = default_account if is_rt else ""

    metrics = {"replies": 0, "likes": 0, "retweets": 0, "bookmarks": 0, "views": 0}
    for stat in card.select(".tweet-stat"):
        txt = stat.get_text(strip=True)
        val = parse_metric_val(txt)
        if stat.select_one(".icon-comment"): metrics["replies"] = val
        elif stat.select_one(".icon-retweet"): metrics["retweets"] = val
        elif stat.select_one(".icon-heart"): metrics["likes"] = val
        elif stat.select_one(".icon-views"): metrics["views"] = val
        elif stat.select_one(".icon-bookmark"): metrics["bookmarks"] = val

    media = extract_media_from_nitter_card(card)
    return {
        "platform": "twitter", "source_name": "nitter",
        "account": {
            "numeric_id": f"ext_{u_name}", "username": u_name, "display_name": d_name,
            "avatar_url": av_url, "avatar_original_url": av_url.replace("_bigger.", "."),
            "description": "", "profile_history": []
        },
        "post": {
            "id": p_id, "conversation_id": p_id, "reply_to_tweet_id": None, "reply_to_handle": None,
            "created_at": created_at, "full_text": full_text, "via": "Nitter",
            "source_name": "nitter", "source_domain": "nitter.space", "is_repost": is_rt,
            "is_pinned": is_pinned, "retweeted_by": retweeted_by, "wayback_url": "",
            "original_url": f"https://x.com/{u_name}/status/{p_id}", "metrics": metrics, "urls": []
        },
        "media": media
    }

def parse_nitter_html_tweets(html_str: str, default_account: str = "") -> List[Dict[str, Any]]:
    if not html_str: return []
    soup = BeautifulSoup(html_str, "html.parser")
    cards = soup.select(".timeline-item")
    results = []
    for c in cards:
        p = parse_nitter_card(c, default_account=default_account)
        if p: results.append(p)
    return results
