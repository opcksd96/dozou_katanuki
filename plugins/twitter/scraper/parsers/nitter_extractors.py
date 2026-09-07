# plugins/twitter/scraper/parsers/nitter_extractors.py (SPEC-PLUGIN-001 / 100行以下)
import re
from datetime import datetime, timezone
from typing import Any, Optional
from urllib.parse import parse_qs, urlparse

def get_filename_from_url(url: str) -> str:
    """URLからクエリパラメータや拡張子を考慮してクリーンなファイル名を生成"""
    if not url: return ""
    parsed = urlparse(url)
    fn = parsed.path.split("/")[-1]
    for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]:
        if fn.endswith(sfx): fn = fn[:-len(sfx)]
    if "." not in fn:
        qs = parse_qs(parsed.query)
        if "format" in qs:
            fn = f"{fn}.{qs['format'][0]}"
    return fn

def parse_nitter_datetime(dt_val: Any, tweet_id: Optional[str] = None) -> str:
    """Nitterの日時表示(例: 'Sep 3, 2026 · 4:13 AM UTC')またはSnowflake IDからUTC文字列生成"""
    if dt_val:
        s = str(dt_val).strip().replace("·", "").replace("  ", " ").strip()
        for fmt in [
            "%b %d, %Y %I:%M %p UTC", "%b %d, %Y %H:%M UTC",
            "%Y-%m-%d %H:%M:%S", "%Y-%m-%dT%H:%M:%S%z"
        ]:
            try:
                dt = datetime.strptime(s, fmt)
                if not dt.tzinfo: dt = dt.replace(tzinfo=timezone.utc)
                return dt.astimezone(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
            except Exception: pass

    if tweet_id and str(tweet_id).isdigit():
        tid = int(tweet_id)
        if tid > 100000000000:
            ms = (tid >> 22) + 1288834974657
            return datetime.fromtimestamp(ms / 1000.0, tz=timezone.utc).strftime("%Y-%m-%d %H:%M:%S")

    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")

def parse_metric_val(val: Any) -> int:
    """'14K', '2.5M', '6,808' などの数値を整数に変換"""
    if val is None: return 0
    if isinstance(val, (int, float)): return int(val)
    s = str(val).strip().replace(",", "")
    if not s: return 0
    mult = 1
    if s.endswith(("K", "k")): mult, s = 1000, s[:-1]
    elif s.endswith(("M", "m")): mult, s = 1000000, s[:-1]
    try: return int(float(s) * mult)
    except Exception: return 0

def clean_html_text(text: str) -> str:
    """HTMLタグを除去して改行を保持したプレーンテキスト化"""
    if not text: return ""
    clean = re.sub(r"<br\s*/?>", "\n", text, flags=re.IGNORECASE)
    clean = re.sub(r"<[^>]+>", "", clean)
    return clean.strip()
