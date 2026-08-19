# plugins/twitter/scraper/core/scraper.py (100行以下)
import io
import os
import urllib.parse
from typing import Any, Dict, List, Optional
import requests
from warcio.statusandheaders import StatusAndHeaders
from warcio.warcwriter import WARCWriter


class Scraper:
    """Wayback CDX走査 & HTTPフェッチ & 原本WARCストリーム保存 (SPEC-PLUGIN-001)"""

    CDX_API_URL = "https://web.archive.org/cdx/search/cdx"

    def __init__(self, platform: str = "twitter", output_dir: str = "backups/dumps"):
        self.platform = platform
        self.output_dir = output_dir
        self.session = requests.Session()
        self.session.headers.update({
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 dozou_katanuki/1.0"
        })

    def search_cdx(self, account: str, limit: int = 50) -> List[Dict[str, str]]:
        """Wayback CDX API を走査してスナップショット一覧を取得"""
        url_match = f"twitter.com/{account}/status/*"
        params = {
            "url": url_match,
            "output": "json",
            "fl": "timestamp,original,mimetype,statuscode,digest",
            "filter": "statuscode:200",
            "limit": str(limit),
            "collapse": "urlkey",
        }
        try:
            resp = self.session.get(self.CDX_API_URL, params=params, timeout=15)
            if resp.status_code == 200:
                rows = resp.json()
                if len(rows) > 1:
                    headers = rows[0]
                    return [dict(zip(headers, row)) for row in rows[1:]]
        except Exception as e:
            print(f"[Scraper] CDX Search error: {e}")
        return []

    def fetch_snapshot(self, timestamp: str, original_url: str, account: str = "") -> Optional[str]:
        """Wayback スナップショットを取得し原本を snapshot.warc.gz にストリーム保存"""
        wayback_url = f"https://web.archive.org/web/{timestamp}id_/{original_url}"
        try:
            resp = self.session.get(wayback_url, timeout=20)
            if resp.status_code == 200:
                self._save_warc_dump(original_url, timestamp, resp, account)
                return resp.text
        except Exception as e:
            print(f"[Scraper] Snapshot fetch error ({wayback_url}): {e}")
        return None

    def _save_warc_dump(self, original_url: str, timestamp: str, resp: requests.Response, account: str) -> None:
        """原本 WARC (snapshot.warc.gz) のストリーム保存"""
        try:
            parsed = urllib.parse.urlparse(original_url)
            path_parts = [p for p in parsed.path.split("/") if p]
            post_id = path_parts[-1] if len(path_parts) >= 3 else timestamp
            username = account or (path_parts[0] if path_parts else "unknown")

            dump_dir = os.path.join(self.output_dir, self.platform, username, post_id)
            os.makedirs(dump_dir, exist_ok=True)
            warc_file = os.path.join(dump_dir, "snapshot.warc.gz")

            with open(warc_file, "ab") as f:
                writer = WARCWriter(f, gzip=True)
                headers_list = [(k, v) for k, v in resp.headers.items()]
                http_headers = StatusAndHeaders(f"{resp.status_code} OK", headers_list, protocol="HTTP/1.1")
                record = writer.create_warc_record(
                    uri=original_url,
                    record_type="response",
                    payload=io.BytesIO(resp.content),
                    http_headers=http_headers,
                    warc_content_type=resp.headers.get("Content-Type", "text/html")
                )
                writer.write_record(record)
        except Exception as e:
            print(f"[Scraper] WARC dump error: {e}")
