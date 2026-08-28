# plugins/twitter/scraper/parsers/sotwe_extractors.py (SPEC-PLUGIN-001 / 100行以下)
import os, re, time
from datetime import datetime, timezone
from typing import Any, Dict, List, Optional
from urllib.parse import urlparse

VUE_EXTRACT_JS = """(() => {
    const all = Array.from(document.querySelectorAll('*')), found = [], seen = new Set();
    function inspect(vm) {
        if (!vm) return;
        const cands = [vm.tweet, vm.post, vm.item, vm.data, vm.$props?.tweet, vm.$props?.post, vm.$props?.item, vm.$props?.data, vm.$data?.tweet, vm.$data?.post, vm.$data?.item, vm.$data?.data];
        for (const c of cands) {
            if (c && typeof c === 'object' && (c.id || c.id_str) && (c.text !== undefined || c.full_text !== undefined || c.mediaEntities)) {
                const tid = String(c.id || c.id_str);
                if (tid && !seen.has(tid)) { seen.add(tid); found.push(c); }
            }
        }
        if (vm.$children && Array.isArray(vm.$children)) for (const ch of vm.$children) inspect(ch);
    }
    for (const el of all) { if (el.__vue__) inspect(el.__vue__); }
    return found;
})()"""

def get_filename_from_url(url: str) -> str:
    if not url: return ""
    fn = urlparse(url.split("?")[0]).path.split("/")[-1]
    for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: fn = fn[:-len(sfx)] if fn.endswith(sfx) else fn
    return fn

def build_streamsaver_url(filename: str) -> str:
    return f"https://jimmywarting.github.io/StreamSaver.js/www.sotwe.com/{filename}" if filename else ""

def parse_iso_datetime(dt_val: Any) -> str:
    if not dt_val: return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
    if isinstance(dt_val, (int, float)):
        ts = dt_val / 1000.0 if dt_val > 1e11 else dt_val
        return datetime.fromtimestamp(ts, tz=timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
    try:
        dt = datetime.fromisoformat(str(dt_val).replace(".000Z", "+00:00").replace("Z", "+00:00"))
        return dt.astimezone(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
    except Exception: return str(dt_val)

def extract_metrics_dict(raw: Dict[str, Any]) -> Dict[str, int]:
    return {"replies": int(raw.get("replyCount", 0)), "likes": int(raw.get("favoriteCount", 0)),
            "retweets": int(raw.get("retweetCount", 0)), "bookmarks": 0, "views": int(raw.get("viewCount", 0))}

def extract_media_entities(raw: Dict[str, Any]) -> List[Dict[str, Any]]:
    media_list = []
    raw_media = raw.get("mediaEntities") or raw.get("extended_entities", {}).get("media") or raw.get("entities", {}).get("media") or []
    for m in raw_media:
        m_type = m.get("type", "image")
        direct_vid = m.get("videoURL") or ""
        v_list = (m.get("videoInfo") or {}).get("variants", [])
        mp4_vars = [v for v in v_list if "mp4" in (v.get("type") or v.get("content_type", "") or v.get("url", ""))]
        if mp4_vars:
            best_mp4 = max(mp4_vars, key=lambda x: int(x.get("bitrate") or 0))
            direct_vid = best_mp4.get("url") or direct_vid
        thumb = m.get("mediaURL") or m.get("media_url_https") or m.get("url") or ""
        target_u = direct_vid if direct_vid else thumb
        fn = get_filename_from_url(target_u)
        if "." not in fn and target_u:
            ext = "mp4" if direct_vid or m_type == "video" else "jpg"
            fn = f"{fn}.{ext}"
        ss_url = build_streamsaver_url(fn)
        variants = []
        for v in v_list:
            v_url = v.get("url")
            if v_url:
                v_fn = get_filename_from_url(v_url)
                variants.append({"content_type": v.get("type") or v.get("content_type", ""), "bitrate": v.get("bitrate", 0), "filename": v_fn, "streamsaver_url": build_streamsaver_url(v_fn), "direct_url": v_url})
        media_list.append({"media_id": fn, "url": target_u, "download_url": target_u, "type": "video" if direct_vid or m_type == "video" else "image",
                           "width": (m.get("sizes", {}).get("large", {}) or {}).get("w", 0), "height": (m.get("sizes", {}).get("large", {}) or {}).get("h", 0),
                           "filename": fn, "streamsaver_url": ss_url, "thumbnail_url": thumb, "stream_variants": variants})
    return media_list

def normalize_vue_tweet(raw: Dict[str, Any], default_account: str) -> Dict[str, Any]:
    t_id = str(raw.get("id") or raw.get("id_str") or "")
    u_info = raw.get("user") or {}
    u_name = u_info.get("screenName") or u_info.get("screen_name") or default_account
    d_name = u_info.get("name") or raw.get("fullname") or u_name
    c_at = parse_iso_datetime(raw.get("createdAt"))
    f_text = raw.get("text") or raw.get("full_text") or ""
    media = extract_media_entities(raw)
    return {
        "platform": "twitter", "source_name": "sotwe",
        "account": {"numeric_id": str(u_info.get("id") or u_info.get("id_str") or f"ext_{u_name}"), "username": u_name, "display_name": d_name,
                    "avatar_url": (u_info.get("profileImage") or u_info.get("profile_image_url_https") or "").replace("_normal.", "."),
                    "avatar_original_url": (u_info.get("profileImage") or u_info.get("profile_image_url_https") or "").replace("_normal.", "."),
                    "description": u_info.get("description", ""), "profile_history": []},
        "post": {"id": t_id, "conversation_id": t_id, "reply_to_tweet_id": None, "reply_to_handle": None, "created_at": c_at, "full_text": f_text,
                 "via": "Sotwe", "source_name": "sotwe", "source_domain": "sotwe.com", "is_repost": False, "is_pinned": bool(raw.get("isPinned")),
                 "retweeted_by": "", "wayback_url": "",
                 "sotwe_url": f"https://www.sotwe.com/tweet/{t_id}" if t_id else "", "original_url": f"https://x.com/{u_name}/status/{t_id}" if t_id else "",
                 "metrics": extract_metrics_dict(raw), "urls": []},
        "media": media
    }
