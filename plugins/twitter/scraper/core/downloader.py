# plugins/twitter/scraper/core/downloader.py (100行以下)
import os
import sqlite3
import time
from typing import Optional
import requests


class Downloader:
    """3段階メディア保全 ＆ Stash/Localインジェクション (SPEC-PLUGIN-001)"""

    def __init__(self, db_path: str = "archive.db", blobs_dir: str = "blobs", stash_url: str = "http://127.0.0.1:9999"):
        self.db_path = db_path
        self.blobs_dir = blobs_dir
        self.stash_url = stash_url
        os.makedirs(self.blobs_dir, exist_ok=True)
        self.session = requests.Session()

    def process_queued_media(self, article_id: Optional[str] = None) -> int:
        """QUEUED 状態のメディアを順次ダウンロードして状態を遷移します"""
        with sqlite3.connect(self.db_path) as conn:
            cur = conn.cursor()
            if article_id:
                cur.execute("SELECT id, url FROM media WHERE article_id = ? AND download_status = 'QUEUED'", (article_id,))
            else:
                cur.execute("SELECT id, url FROM media WHERE download_status = 'QUEUED' LIMIT 20")
            records = cur.fetchall()

        success_count = 0
        for media_id, url in records:
            status = self._download_single(media_id, url)
            self._update_status(media_id, status)
            if status == "COMPLETED":
                success_count += 1
        return success_count

    def _download_single(self, media_id: str, url: str) -> str:
        """第1段階: 直接HTTP GETダウンロード"""
        try:
            resp = self.session.get(url, stream=True, timeout=15)
            if resp.status_code == 200:
                dest_path = os.path.join(self.blobs_dir, media_id)
                with open(dest_path, "wb") as f:
                    for chunk in resp.iter_content(chunk_size=65536):
                        if chunk:
                            f.write(chunk)
                return "COMPLETED"
            elif resp.status_code == 404:
                return "DEAD_404"
            else:
                return "QUEUED"
        except Exception:
            return "QUEUED"

    def _update_status(self, media_id: str, status: str) -> None:
        with sqlite3.connect(self.db_path) as conn:
            now_ts = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())
            conn.execute("""
                UPDATE media
                SET download_status = ?, updated_at = ?
                WHERE id = ?
            """, (status, now_ts, media_id))
            conn.commit()
