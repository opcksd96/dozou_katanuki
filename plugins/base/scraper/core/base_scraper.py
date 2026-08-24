# plugins/base/scraper/core/base_scraper.py (SPEC-PLUGIN-001 / 100行以下)
import io, os, re, time
from typing import Any, Callable, Dict, List, Optional
import requests
from warcio.archiveiterator import ArchiveIterator
from warcio.statusandheaders import StatusAndHeaders
from warcio.warcwriter import WARCWriter


class BaseScraper:
    """Wayback CDX走査 & WARCキャッシュ優先読み込み & HTTPフェッチ & 原本WARC保存基底クラス"""
    CDX_API_URL = "https://web.archive.org/cdx/search/cdx"

    def __init__(self, platform: str = "base", output_dir: str = "backups/dumps"):
        self.platform, self.output_dir = platform, output_dir
        self.session = requests.Session()
        self.session.headers.update({"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", "Accept": "*/*"})

    def _get_warc_path(self, original_url: str, timestamp: str, account: str = "") -> str:
        m = re.search(r'status(?:es)?/(\d+)', original_url)
        post_id = m.group(1) if m else timestamp
        m_user = re.search(r'(?:twitter\.com|x\.com)/([a-zA-Z0-9_]+)', original_url)
        username = (account or (m_user.group(1) if m_user else "unknown")).lower()
        return os.path.join(self.output_dir, self.platform, username, post_id, "snapshot.warc.gz")

    def search_cdx(self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None) -> List[Dict[str, str]]:
        def _log(msg: str) -> None:
            if log_fn: log_fn(msg)
            print(f"[CDX:{self.platform}] {msg}", flush=True)

        clean_acc = account.lstrip("@").strip()
        params = {"url": f"twitter.com/{clean_acc}", "matchType": "prefix", "output": "json",
                  "fl": "timestamp,original,mimetype,statuscode,digest", "filter": "!statuscode:[45]..", "collapse": "urlkey"}
        if limit > 0: params["limit"] = str(limit)
        _log(f"CDX Query: target='twitter.com/{clean_acc}' (matchType=prefix), limit={limit or 'unlimited'}")

        for attempt in range(1, 4):
            t0 = time.time()
            try:
                _log(f"Attempt {attempt}/3: Fetching CDX index from web.archive.org...")
                resp = self.session.get(self.CDX_API_URL, params=params, timeout=40)
                ms = int((time.time() - t0) * 1000)
                _log(f"Attempt {attempt}/3 Response: HTTP {resp.status_code} ({ms}ms), size={len(resp.content)} bytes")
                if resp.status_code == 200:
                    rows = resp.json()
                    raw_items = [dict(zip(rows[0], r)) for r in rows[1:]] if len(rows) > 1 else []
                    acc_pat = re.compile(rf'^(?:https?://)?(?:[a-zA-Z0-9_-]+\.)?(?:twitter\.com|x\.com)/{re.escape(clean_acc)}(?:[/?#]|$)', re.IGNORECASE)
                    matching_items = [i for i in raw_items if acc_pat.match(i.get("original", ""))]
                    status_items = [i for i in matching_items if "status" in i.get("original", "").lower()]
                    final_items = status_items if status_items else matching_items
                    _log(f"CDX OK: Found {len(final_items)} post snapshots for @{clean_acc} (total scanned: {len(raw_items)}, excluded prefix overlap: {len(raw_items) - len(matching_items)})")
                    return final_items
                elif resp.status_code in (429, 503):
                    w = attempt * 3; _log(f"Rate-limited / Busy (HTTP {resp.status_code}). Waiting {w}s..."); time.sleep(w)
                else: _log(f"HTTP {resp.status_code}: {resp.text[:100]}")
            except requests.exceptions.Timeout: _log(f"Attempt {attempt}/3 Timed out after 40s. Retrying...")
            except Exception as e: _log(f"Attempt {attempt}/3 Error: {type(e).__name__}: {e}")
            if attempt < 3: time.sleep(1.5)
        _log(f"CDX Failed: Could not retrieve snapshots for @{clean_acc} after 3 attempts.")
        return []

    def fetch_snapshot(self, timestamp: str, original_url: str, account: str = "") -> Optional[str]:
        warc_path = self._get_warc_path(original_url, timestamp, account)
        if os.path.exists(warc_path) and os.path.getsize(warc_path) > 0:
            try:
                with open(warc_path, "rb") as stream:
                    for record in ArchiveIterator(stream):
                        if record.rec_type == "response": return record.raw_stream.read().decode("utf-8", errors="ignore")
            except Exception as e: print(f"[BaseScraper] WARC cache read error ({warc_path}): {e}", flush=True)

        wayback_url = f"https://web.archive.org/web/{timestamp}id_/{original_url}"
        for attempt in range(1, 3):
            try:
                resp = self.session.get(wayback_url, timeout=12.0)
                if resp.status_code == 200:
                    self._save_warc_dump(warc_path, original_url, timestamp, resp)
                    return resp.text
                elif resp.status_code in (429, 503): time.sleep(1.0 * attempt)
            except Exception as e:
                if attempt == 2: print(f"[BaseScraper] Fetch failed for {wayback_url}: {e}", flush=True)
                time.sleep(0.5 * attempt)
        return None

    def _save_warc_dump(self, warc_file: str, original_url: str, timestamp: str, resp: requests.Response) -> None:
        try:
            os.makedirs(os.path.dirname(warc_file), exist_ok=True)
            with open(warc_file, "wb") as f:
                writer = WARCWriter(f, gzip=True)
                headers_list = [(k, v) for k, v in resp.headers.items()]
                http_headers = StatusAndHeaders(f"{resp.status_code} OK", headers_list, protocol="HTTP/1.1")
                writer.write_record(writer.create_warc_record(
                    uri=original_url, record_type="response", payload=io.BytesIO(resp.content),
                    http_headers=http_headers, warc_content_type=resp.headers.get("Content-Type", "text/html")
                ))
        except Exception as e: print(f"[BaseScraper] WARC dump error: {e}", flush=True)
