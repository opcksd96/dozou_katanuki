# plugins/twitter/scraper/core/scraper.py (100行以下)
import io, os, re, time, urllib.parse
from typing import Any, Dict, List, Optional
import requests
from warcio.archiveiterator import ArchiveIterator
from warcio.statusandheaders import StatusAndHeaders
from warcio.warcwriter import WARCWriter


class Scraper:
    """Wayback CDX走査 & WARCキャッシュ優先読み込み & HTTPフェッチ & 原本WARC保存 (SPEC-PLUGIN-001)"""
    CDX_API_URL = "https://web.archive.org/cdx/search/cdx"

    def __init__(self, platform: str = "twitter", output_dir: str = "backups/dumps"):
        self.platform = platform
        self.output_dir = output_dir
        self.session = requests.Session()
        self.session.headers.update({
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
            "Accept-Language": "ja,en-US;q=0.9,en;q=0.8",
            "Sec-Ch-Ua": '"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"',
            "Sec-Ch-Ua-Platform": '"Windows"'
        })

    def _get_warc_path(self, original_url: str, timestamp: str, account: str = "") -> str:
        m = re.search(r'status(?:es)?/(\d+)', original_url)
        post_id = m.group(1) if m else timestamp
        m_user = re.search(r'(?:twitter\.com|x\.com)/([a-zA-Z0-9_]+)', original_url)
        username = (account or (m_user.group(1) if m_user else "unknown")).lower()
        return os.path.join(self.output_dir, self.platform, username, post_id, "snapshot.warc.gz")

    def search_cdx(self, account: str, limit: int = 0) -> List[Dict[str, str]]:
        """Wayback CDX API を走査してスナップショット一覧を取得 (limit<=0 は無制限全件)"""
        params = {"url": f"twitter.com/{account}/status/*", "output": "json",
                  "fl": "timestamp,original,mimetype,statuscode,digest", "filter": "statuscode:200", "collapse": "urlkey"}
        if limit > 0: params["limit"] = str(limit)
        for attempt in range(2):
            try:
                resp = self.session.get(self.CDX_API_URL, params=params, timeout=25)
                if resp.status_code == 200:
                    rows = resp.json()
                    if len(rows) > 1: return [dict(zip(rows[0], row)) for row in rows[1:]]
                elif resp.status_code == 429: time.sleep(1.5 * (attempt + 1))
            except Exception as e:
                print(f"[Scraper] CDX Search attempt {attempt+1} error: {e}")
                time.sleep(1.0)
        return []

    def fetch_snapshot(self, timestamp: str, original_url: str, account: str = "") -> Optional[str]:
        """原本WARCが存在すればローカルから即時返却。未存在時のみWayback取得 (非同期耐性/高速フォールバック)"""
        warc_path = self._get_warc_path(original_url, timestamp, account)
        if os.path.exists(warc_path) and os.path.getsize(warc_path) > 0:
            try:
                with open(warc_path, "rb") as stream:
                    for record in ArchiveIterator(stream):
                        if record.rec_type == "response":
                            return record.raw_stream.read().decode("utf-8", errors="ignore")
            except Exception as e:
                print(f"[Scraper] WARC cache read error ({warc_path}): {e}")

        wayback_url = f"https://web.archive.org/web/{timestamp}id_/{original_url}"
        try:
            resp = requests.get(wayback_url, headers=self.session.headers, timeout=4.0)
            if resp.status_code == 200:
                self._save_warc_dump(warc_path, original_url, timestamp, resp)
                return resp.text
        except Exception:
            pass
        return None

    def _save_warc_dump(self, warc_file: str, original_url: str, timestamp: str, resp: requests.Response) -> None:
        """原本 WARC (snapshot.warc.gz) の新規/上書き保存"""
        try:
            os.makedirs(os.path.dirname(warc_file), exist_ok=True)
            with open(warc_file, "wb") as f:
                writer = WARCWriter(f, gzip=True)
                headers_list = [(k, v) for k, v in resp.headers.items()]
                http_headers = StatusAndHeaders(f"{resp.status_code} OK", headers_list, protocol="HTTP/1.1")
                record = writer.create_warc_record(
                    uri=original_url, record_type="response", payload=io.BytesIO(resp.content),
                    http_headers=http_headers, warc_content_type=resp.headers.get("Content-Type", "text/html")
                )
                writer.write_record(record)
        except Exception as e:
            print(f"[Scraper] WARC dump error: {e}")
