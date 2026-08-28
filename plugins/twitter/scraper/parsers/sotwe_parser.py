# plugins/twitter/scraper/parsers/sotwe_parser.py (SPEC-PLUGIN-001 / 100行以下)
from typing import Any, Dict, List
from bs4 import BeautifulSoup
try:
    from plugins.twitter.scraper.parsers.sotwe_extractors import (
        normalize_vue_tweet, parse_iso_datetime, get_filename_from_url, build_streamsaver_url
    )
except ImportError:
    from parsers.sotwe_extractors import (
        normalize_vue_tweet, parse_iso_datetime, get_filename_from_url, build_streamsaver_url
    )

def parse_sotwe_vue_tweets(vue_records: List[Dict[str, Any]], default_account: str) -> List[Dict[str, Any]]:
    """ブラウザ内のVueオブジェクト配列を直接標準スキーマへ変換"""
    results, seen_ids = [], set()
    for item in (vue_records or []):
        t_id = str(item.get("id") or item.get("id_str") or "")
        if t_id and t_id not in seen_ids:
            seen_ids.add(t_id)
            results.append(normalize_vue_tweet(item, default_account))
    return results

def parse_sotwe_html_tweets(html_str: str, default_account: str) -> List[Dict[str, Any]]:
    """Sotwe HTMLから全カードを抽出（HTMLフォールバック用）"""
    if not html_str: return []
    soup = BeautifulSoup(html_str, "html.parser")
    bio_el = soup.select_one(".break-word .dynamic-link-content")
    page_bio = bio_el.get_text(separator="\n").strip() if bio_el else ""
    header_name_el = soup.select_one(".profile-name, .v-card__title .font-weight-bold")
    page_name = header_name_el.get_text(strip=True) if header_name_el else default_account
    header_av_el = soup.select_one(".profile-avatar img, .v-avatar img")
    page_av = header_av_el["src"].replace("_normal.", ".") if header_av_el and header_av_el.get("src") else ""

    cards, results, seen_ids = soup.select("div.tweet-card"), [], set()
    for idx, card in enumerate(cards):
        is_repost, is_pinned = bool(card.select_one(".v-card__title .fa-retweet")), bool(card.select_one(".pinned-text, .pinned-icon"))
        profile_link = card.select_one(".tweet-profile a[href^='/']")
        author_user = profile_link["href"].strip("/").split("/")[0] if (profile_link and profile_link.get("href")) else default_account
        name_el = card.select_one(".tweet-profile--text span.font-weight-medium")
        display_name = name_el.get_text(strip=True) if name_el else page_name
        av_img = card.select_one(".v-avatar img, .tweet-profile img")
        raw_av = av_img["src"] if (av_img and av_img.get("src")) else page_av
        avatar_url = raw_av.replace("_normal.", ".") if "_normal." in raw_av else raw_av

        time_el = card.select_one("time[datetime]")
        created_at_str = parse_iso_datetime(time_el.get("datetime") if time_el else None)
        text_el = card.select_one(".tweet-text .dynamic-link-content")
        full_text = text_el.get_text(separator="\n").strip() if text_el else ""

        media_list = []
        for img in card.select(".media-carousel img[src], .media-carousel-image img[src]"):
            u = img.get("src", "")
            if u and "profile_images" not in u and not any(m["url"] == u for m in media_list):
                fn = get_filename_from_url(u)
                media_list.append({"media_id": fn, "url": u, "download_url": u, "type": "image", "width": 0, "height": 0, "filename": fn, "streamsaver_url": build_streamsaver_url(fn)})
        for vid in card.select("video.video-player source[type='video/mp4']"):
            u = vid.get("src", "")
            if u and not any(m["url"] == u for m in media_list):
                fn = get_filename_from_url(u)
                media_list.append({"media_id": fn, "url": u, "download_url": u, "type": "video", "width": 0, "height": 0, "filename": fn, "streamsaver_url": build_streamsaver_url(fn)})

        post_id = f"sotwe_{author_user}_{idx+1}"
        results.append({
            "platform": "twitter", "source_name": "sotwe",
            "account": {"numeric_id": f"ext_{author_user}", "username": author_user, "display_name": display_name,
                        "avatar_url": avatar_url, "avatar_original_url": avatar_url, "description": page_bio,
                        "profile_history": [{"display_name": display_name, "avatar_original_url": avatar_url, "observed_at": created_at_str}] if avatar_url else []},
            "post": {"id": post_id, "conversation_id": post_id, "reply_to_tweet_id": None, "reply_to_handle": None,
                     "created_at": created_at_str, "full_text": full_text, "via": "Sotwe", "source_name": "sotwe", "source_domain": "sotwe.com",
                     "is_repost": is_repost, "is_pinned": is_pinned, "retweeted_by": "", "wayback_url": f"https://x.com/{author_user}/status/{post_id}",
                     "sotwe_url": f"https://www.sotwe.com/{author_user}", "original_url": f"https://x.com/{author_user}/status/{post_id}",
                     "metrics": {"replies": 0, "likes": 0, "retweets": 0, "bookmarks": 0, "views": 0}, "urls": []},
            "media": media_list
        })
    return results
