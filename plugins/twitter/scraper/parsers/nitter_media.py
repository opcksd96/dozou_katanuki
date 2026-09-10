# plugins/twitter/scraper/parsers/nitter_media.py (SPEC-PLUGIN-001 / 100行以下)
import re
from typing import Any, Dict, List
from urllib.parse import unquote
try:
    from plugins.twitter.scraper.parsers.nitter_extractors import get_filename_from_url
except ImportError:
    from parsers.nitter_extractors import get_filename_from_url

def normalize_twitter_image_url(url: str, name: str = "orig") -> str:
    """pbs.twimg.com 画像URLを指定解像度(name=orig / large)に変換"""
    if not url: return ""
    if "pbs.twimg.com/media" in url: return f"{url.split('?')[0]}?name={name}"
    if "/pic/orig/media%2F" in url or "/pic/orig/media/" in url:
        clean = unquote(url.split("/pic/orig/media/")[-1].split("/pic/orig/media%2F")[-1])
        return f"https://pbs.twimg.com/media/{clean}?name={name}"
    return url

def extract_media_from_nitter_card(card_soup) -> List[Dict[str, Any]]:
    """Nitter のタイムライン/ツイート要素から画像および動画メディアを抽出"""
    media_list, seen_keys = [], set()
    attach = card_soup.select_one(".attachments")
    if not attach: return []
    for vid in attach.select("video"):
        src = (vid.select_one("source[src]").get("src") if vid.select_one("source[src]") else None) or vid.get("src") or ""
        if not src: continue
        fn = get_filename_from_url(src)
        if "." not in fn: fn = f"{fn}.mp4"
        if fn in seen_keys: continue
        seen_keys.add(fn); dim_m = re.search(r'/(\d+)x(\d+)/', src)
        w, h = (int(dim_m.group(1)), int(dim_m.group(2))) if dim_m else (0, 0)
        media_list.append({"media_id": fn, "type": "video", "original_url": src, "download_url": src, "filename": fn, "thumbnail_url": vid.get("poster", ""), "width": w, "height": h, "variants": [{"content_type": "video/mp4", "bit_rate": 0, "url": src, "filename": fn}]})
    for img_wrap in attach.select(".attachment.image, a.still-image"):
        img_tag = img_wrap.select_one("img") if img_wrap.name != "img" else img_wrap
        raw_src = img_tag.get("src") if img_tag else ""
        if not raw_src and img_wrap.name == "a": raw_src = img_wrap.get("href", "")
        if not raw_src or "profile_images" in raw_src: continue
        high_res = normalize_twitter_image_url(raw_src, "orig"); large_res = normalize_twitter_image_url(raw_src, "large")
        fn = get_filename_from_url(high_res)
        if "." not in fn: fn = f"{fn}.jpg"
        if fn in seen_keys: continue
        seen_keys.add(fn)
        variants = [{"content_type": "image/jpeg", "bit_rate": 10000, "url": high_res, "filename": fn},
                    {"content_type": "image/jpeg", "bit_rate": 5000, "url": large_res, "filename": fn}]
        media_list.append({"media_id": fn, "type": "image", "original_url": high_res, "download_url": high_res, "filename": fn, "thumbnail_url": raw_src, "width": 0, "height": 0, "variants": variants})
    return media_list
