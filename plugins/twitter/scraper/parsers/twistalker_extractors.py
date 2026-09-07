# plugins/twitter/scraper/parsers/twistalker_extractors.py (SPEC-PLUGIN-001 / 100行以下)
import re
from datetime import datetime, timezone
from typing import Any, Optional
from urllib.parse import urlparse

def get_filename_from_url(url: str) -> str:
    if not url: return ""
    parsed = urlparse(url)
    fn = parsed.path.split("/")[-1]
    for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]:
        if fn.endswith(sfx): fn = fn[:-len(sfx)]
    if "." not in fn:
        import urllib.parse
        qs = urllib.parse.parse_qs(parsed.query)
        if "format" in qs: fn = f"{fn}.{qs['format'][0]}"
    return fn

def parse_twitter_datetime(dt_val: Any, tweet_id: Optional[str] = None) -> str:
    """Twitter日時形式またはSnowflake IDからUTC日時文字列(YYYY-MM-DD HH:MM:SS)を生成"""
    if dt_val:
        for fmt in ("%a %b %d %H:%M:%S %z %Y", "%Y-%m-%d %H:%M:%S"):
            try:
                dt = datetime.strptime(str(dt_val), fmt)
                return dt.astimezone(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
            except Exception: pass
        try:
            dt = datetime.fromisoformat(str(dt_val).replace(".000Z", "+00:00").replace("Z", "+00:00"))
            return dt.astimezone(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
        except Exception: pass

    # Snowflake IDからのミリ秒計算: (ID >> 22) + 1288834974657
    if tweet_id and str(tweet_id).isdigit():
        tid = int(tweet_id)
        if tid > 100000000000:
            ms = (tid >> 22) + 1288834974657
            return datetime.fromtimestamp(ms / 1000.0, tz=timezone.utc).strftime("%Y-%m-%d %H:%M:%S")

    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")

def parse_metric_val(val: Any) -> int:
    """'14K', '2.5M', '670' などの数値を整数に変換"""
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
    """HTMLタグを除去してプレーンテキスト化"""
    if not text: return ""
    clean = re.sub(r"<br\s*/?>", "\n", text, flags=re.IGNORECASE)
    clean = re.sub(r"<[^>]+>", "", clean)
    return clean.strip()
