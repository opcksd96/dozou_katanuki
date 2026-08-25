"""
sotwe_extractors.py (SPEC-PLUGIN-001 / 100行以下)
Sotwe HTMLからメトリクス、メディア、アカウント情報、Snowflake IDを抽出するヘルパーモジュール
"""
from datetime import datetime, timezone
import re
from typing import Any, Dict, List, Optional
from bs4 import Tag

TWITTER_EPOCH = 1288834974657

def parse_iso_datetime(dt_str: Optional[str]) -> str:
    """time[datetime] の文字列を標準ISO8601 (UTC) に正規化"""
    if not dt_str:
        return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    try:
        dt = datetime.fromisoformat(dt_str.replace(".000Z", "+00:00").replace("Z", "+00:00"))
        return dt.astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    except Exception:
        return dt_str

def generate_snowflake_id(dt_str: Optional[str]) -> str:
    """投稿日時から19桁のTwitter Snowflake IDを正確に算出"""
    if not dt_str:
        epoch_ms = int(datetime.now(timezone.utc).timestamp() * 1000)
    else:
        try:
            dt = datetime.fromisoformat(dt_str.replace(".000Z", "+00:00").replace("Z", "+00:00"))
            epoch_ms = int(dt.timestamp() * 1000)
        except Exception:
            epoch_ms = int(datetime.now(timezone.utc).timestamp() * 1000)
    
    snowflake = (epoch_ms - TWITTER_EPOCH) << 22
    return str(max(snowflake, epoch_ms))

def extract_metrics(card: Tag) -> Dict[str, int]:
    """いいね、RT、リプライ、ブックマーク、閲覧数を抽出"""
    metrics = {"replies": 0, "likes": 0, "retweets": 0, "bookmarks": 0, "views": 0}
    stats = card.select(".tweet-stats-item")
    for st in stats:
        txt = st.get_text(strip=True)
        num = int(re.sub(r"[^\d]", "", txt)) if re.sub(r"[^\d]", "", txt) else 0
        if st.select_one(".fa-comment"): metrics["replies"] = num
        elif st.select_one(".fa-heart"): metrics["likes"] = num
        elif st.select_one(".fa-retweet"): metrics["retweets"] = num
        elif st.select_one(".fa-bookmark"): metrics["bookmarks"] = num
        elif st.select_one(".fa-chart-bar"): metrics["views"] = num
    return metrics

def extract_card_media(card: Tag) -> List[Dict[str, Any]]:
    """カルーセル画像、動画MP4、ポスターサムネイルを完全抽出"""
    media_list = []
    for img in card.select(".media-carousel img[src], .media-carousel-image img[src]"):
        src = img.get("src", "")
        if src and "profile_images" not in src and not any(m["url"] == src for m in media_list):
            media_list.append({"url": src, "type": "image", "width": 0, "height": 0})

    for vid in card.select("video.video-player source[type='video/mp4']"):
        src = vid.get("src", "")
        if src and not any(m["url"] == src for m in media_list):
            media_list.append({"url": src, "type": "video", "width": 0, "height": 0})

    for vid_p in card.select("video.video-player[poster]"):
        poster = vid_p.get("poster", "")
        if poster and not any(m["url"] == poster for m in media_list):
            media_list.append({"url": poster, "type": "image", "width": 0, "height": 0})
    return media_list

def extract_status_id_from_media(media_list: List[Dict[str, Any]], created_at_iso: str) -> str:
    """メディアURL内のStatus IDを探索、存在しない場合はSnowflake IDを算出"""
    for m in media_list:
        match = re.search(r"ext_tw_video(?:_thumb)?/(\d{15,20})", m["url"])
        if match:
            return match.group(1)
    return generate_snowflake_id(created_at_iso)
