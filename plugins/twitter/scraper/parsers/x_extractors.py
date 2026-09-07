# plugins/twitter/scraper/parsers/x_extractors.py (SPEC-PLUGIN-001 / 100行以下)
import datetime, re
from typing import Any, Optional
from urllib.parse import parse_qs, urlparse

def snowflake_to_datetime(snowflake_id: str) -> Optional[str]:
    """Twitter Snowflake IDからミリ秒精度のUTC日時文字列 (YYYY-MM-DD HH:MM:SS) を復元"""
    try:
        s_id = int(snowflake_id)
        if s_id <= 30000000000: return None
        ms = (s_id >> 22) + 1288834974657
        dt = datetime.datetime.fromtimestamp(ms / 1000.0, tz=datetime.timezone.utc)
        return dt.strftime("%Y-%m-%d %H:%M:%S")
    except Exception: return None

def parse_metric_number(val: Any) -> int:
    """万、億、K、M表記を含むメトリクス数値を整数へ変換"""
    if val is None: return 0
    if isinstance(val, (int, float)): return int(val)
    txt = str(val).strip().replace(",", "").replace(" ", "")
    if not txt: return 0
    try:
        if txt.endswith("万"): return int(float(txt[:-1]) * 10000)
        if txt.endswith("億"): return int(float(txt[:-1]) * 100000000)
        if txt.endswith(("K", "k")): return int(float(txt[:-1]) * 1000)
        if txt.endswith(("M", "m")): return int(float(txt[:-1]) * 1000000)
        return int(float(txt))
    except Exception:
        m = re.search(r"(\d+(?:\.\d+)?)", txt)
        return int(float(m.group(1))) if m else 0

def get_filename_from_url(url: str) -> str:
    """メディアURLから拡張子付きファイル名を抽出"""
    if not url: return ""
    parsed = urlparse(url)
    fn = parsed.path.split("/")[-1]
    for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]:
        if fn.endswith(sfx): fn = fn[:-len(sfx)]
    if "." not in fn and parsed.query:
        qs = parse_qs(parsed.query)
        fmt = qs.get("format", [""])[0]
        if fmt: fn = f"{fn}.{fmt}"
    return fn
