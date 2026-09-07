# plugins/twitter/scraper/parsers/twistalker_media.py (SPEC-PLUGIN-001 / 100行以下)
from typing import Any, Dict, List
try:
    from plugins.twitter.scraper.parsers.twistalker_extractors import get_filename_from_url
except ImportError:
    from parsers.twistalker_extractors import get_filename_from_url

def extract_media_from_api(entities: Dict[str, Any]) -> List[Dict[str, Any]]:
    """APIツイートのextended_entitiesからメディア一覧を抽出"""
    media_list = []
    raw_media = entities.get("media") or []
    for m in raw_media:
        m_type = m.get("type", "image")
        v_info = m.get("video_info") or {}
        v_list = v_info.get("variants") or []
        mp4_vars = [v for v in v_list if "mp4" in (v.get("content_type") or "")]
        direct_vid, best_bitrate = "", 0
        variants = []
        for v in mp4_vars:
            v_url = (v.get("url") or "").replace("video-s.twimg.com", "pbs.twimg.com")
            b_rate = int(v.get("bitrate") or 0)
            variants.append({"content_type": "video/mp4", "bit_rate": b_rate, "url": v_url, "filename": get_filename_from_url(v_url)})
            if b_rate >= best_bitrate:
                best_bitrate, direct_vid = b_rate, v_url

        thumb = m.get("media_url_https") or m.get("media_url") or ""
        target_url = direct_vid if (m_type in ("video", "animated_gif") and direct_vid) else thumb
        fn = get_filename_from_url(target_url)
        if "." not in fn and target_url:
            fn = f"{fn}.{'mp4' if m_type in ('video', 'animated_gif') else 'jpg'}"

        sizes = (m.get("sizes") or {}).get("large") or {}
        media_list.append({
            "media_id": fn, "type": "video" if m_type in ("video", "animated_gif") else "image",
            "original_url": target_url, "download_url": target_url, "filename": fn,
            "thumbnail_url": thumb, "width": sizes.get("w", 0), "height": sizes.get("h", 0),
            "variants": sorted(variants, key=lambda x: x["bit_rate"], reverse=True)
        })
    return media_list

def extract_media_from_html(card_soup) -> List[Dict[str, Any]]:
    """HTMLカード要素から画像および動画メディアを抽出"""
    media_list, seen_urls = [], set()
    for vid in card_soup.select("video"):
        src = vid.get("src")
        src_tag = vid.select_one("source[src]")
        if src_tag: src = src_tag.get("src") or src
        if not src or src in seen_urls: continue
        seen_urls.add(src)
        fn = get_filename_from_url(src)
        if "." not in fn: fn = f"{fn}.mp4"
        thumb = vid.get("poster", "")
        media_list.append({
            "media_id": fn, "type": "video", "original_url": src, "download_url": src,
            "filename": fn, "thumbnail_url": thumb, "width": 0, "height": 0,
            "variants": [{"content_type": "video/mp4", "bit_rate": 0, "url": src, "filename": fn}]
        })

    for a in card_soup.select(".card-carousel a[data-image], .activity-descp a.thumbnail[data-image], img[src]"):
        u = a.get("data-image") or a.get("src") or a.get("href")
        if not u or u.startswith("#") or "profile_images" in u or u in seen_urls: continue
        seen_urls.add(u)
        fn = get_filename_from_url(u)
        if "." not in fn: fn = f"{fn}.jpg"
        media_list.append({
            "media_id": fn, "type": "image", "original_url": u, "download_url": u,
            "filename": fn, "thumbnail_url": u, "width": 0, "height": 0, "variants": []
        })
    return media_list
