# plugins/twitter/scraper/parsers/x_media.py (SPEC-PLUGIN-001 / 100行以下)
import re
from typing import Any, Dict, List
try:
    from plugins.twitter.scraper.parsers.x_extractors import get_filename_from_url
except ImportError:
    from parsers.x_extractors import get_filename_from_url

def build_image_variant_url(url: str, name: str = "orig") -> str:
    """pbs.twimg.comの画像URLを指定解像度 (orig / large) に変換"""
    if "pbs.twimg.com/media/" not in url: return url
    clean = re.sub(r":(large|orig|small|medium|thumb)$", "", url)
    if "?" in clean:
        base, qs = clean.split("?", 1)
        params = [p for p in qs.split("&") if not p.startswith("name=")]
        params.append(f"name={name}")
        return f"{base}?{'&'.join(params)}"
    return f"{clean}?name={name}"

def build_orig_image_url(url: str) -> str: return build_image_variant_url(url, "orig")

def extract_media_from_syndication(details: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """Syndication APIの mediaDetails からメディアリストを抽出"""
    media_list: List[Dict[str, Any]] = []
    for item in details or []:
        m_type = "image" if item.get("type") == "photo" else item.get("type", "image")
        orig_info = item.get("original_info") or {}
        w, h = orig_info.get("width", 0), orig_info.get("height", 0)
        base_url = item.get("media_url_https") or item.get("url") or ""
        high_res = build_orig_image_url(base_url) if m_type == "image" else base_url
        v_info = item.get("video_info") or {}
        variants_raw = v_info.get("variants") or []
        mp4_vars = sorted([v for v in variants_raw if v.get("content_type") == "video/mp4" and "url" in v], key=lambda x: x.get("bitrate", x.get("bit_rate", 0)), reverse=True)
        if m_type == "image":
            variants = [{"url": high_res, "bit_rate": 10000, "content_type": "image/jpeg"},
                        {"url": build_image_variant_url(base_url, "large"), "bit_rate": 5000, "content_type": "image/jpeg"}]
        else:
            variants = [{"url": v["url"], "bit_rate": v.get("bitrate", v.get("bit_rate", 0)), "content_type": "video/mp4"} for v in mp4_vars]
        dl_url = mp4_vars[0]["url"] if (m_type in ["video", "animated_gif"] and mp4_vars) else high_res
        fn = get_filename_from_url(dl_url or high_res)
        m_id = str(item.get("id_str") or item.get("id") or fn.split(".")[0])
        media_list.append({
            "media_id": m_id, "url": high_res, "download_url": dl_url, "type": m_type,
            "width": w, "height": h, "bitrate": mp4_vars[0].get("bitrate", 0) if mp4_vars else 0,
            "duration": float(v_info.get("duration_millis", 0)) / 1000.0, "filename": fn,
            "thumbnail_url": base_url, "variants": variants
        })
    return media_list

def extract_media_from_html_card(card_soup: Any) -> List[Dict[str, Any]]:
    """HTML article 要素から画像・動画メディアを抽出"""
    media_list, seen = [], set()
    for img in card_soup.find_all("img"):
        src = img.get("src", "")
        if "pbs.twimg.com/media/" in src:
            orig = build_orig_image_url(src); fn = get_filename_from_url(orig); base_id = fn.split(".")[0]
            if base_id not in seen:
                seen.add(base_id)
                media_list.append({
                    "media_id": base_id, "url": orig, "download_url": orig, "type": "image",
                    "width": 0, "height": 0, "bitrate": 0, "duration": 0.0,
                    "filename": fn, "thumbnail_url": src,
                    "variants": [{"url": orig, "bit_rate": 10000, "content_type": "image/jpeg"},
                                 {"url": build_image_variant_url(src, "large"), "bit_rate": 5000, "content_type": "image/jpeg"}]
                })
    for vid in card_soup.find_all("video"):
        src = vid.get("src", "") or ""
        if not src:
            for s in vid.find_all("source"):
                if s.get("src"): src = s["src"]; break
        if src and src not in seen:
            seen.add(src); fn = get_filename_from_url(src)
            media_list.append({
                "media_id": fn.split(".")[0], "url": src, "download_url": src, "type": "video",
                "width": 0, "height": 0, "bitrate": 0, "duration": 0.0,
                "filename": fn, "thumbnail_url": vid.get("poster", ""), "variants": [{"url": src, "bit_rate": 0, "content_type": "video/mp4"}]
            })
    return media_list
