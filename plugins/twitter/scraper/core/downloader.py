# plugins/twitter/scraper/core/downloader.py (100行以下)
import os
import sqlite3
import time
from typing import Optional, Tuple
import requests
from .stash_client import StashClient


class Downloader:
    """3段階メディア保全 ＆ Stash/Localインジェクション (SPEC-PLUGIN-001)"""

    DEFAULT_STORAGE = r"G:\Media_Storage\Influencers"

    def __init__(self, db_path: str = "archive.db", storage_dir: Optional[str] = None):
        self.db_path = db_path
        self.storage_dir = storage_dir or (self.DEFAULT_STORAGE if os.path.exists(r"G:\Media_Storage") else "blobs")
        os.makedirs(self.storage_dir, exist_ok=True)
        self.session = requests.Session()
        self.stash = StashClient()

    def process_queued_media(self, article_id: Optional[str] = None) -> int:
        """QUEUED メディアを順次ダウンロードし Stash 連携・DB 更新"""
        with sqlite3.connect(self.db_path) as conn:
            cur = conn.cursor()
            query = """
                SELECT m.media_id, m.download_url, m.type, a.account_id, ac.username
                FROM media m
                JOIN articles a ON m.article_id = a.id
                JOIN accounts ac ON a.account_id = ac.numeric_id
                WHERE m.download_status = 'QUEUED'
            """
            if article_id:
                query += " AND m.article_id = ?"
                cur.execute(query, (article_id,))
            else:
                query += " LIMIT 20"
                cur.execute(query)
            records = cur.fetchall()

        success_count = 0
        for media_id, url, m_type, _, username in records:
            status, img_id, scn_id = self._download_and_inject(media_id, url, m_type, username or "unknown")
            self._update_status(media_id, status, img_id, scn_id)
            if status == "COMPLETED":
                success_count += 1
        return success_count

    def _download_and_inject(self, media_id: str, url: str, m_type: str, username: str) -> Tuple[str, Optional[str], Optional[str]]:
        """第1段階: 直接ダウンロード ＆ Stash 登録"""
        target_dir = os.path.join(self.storage_dir, "Twitter", username) if "Influencers" in self.storage_dir else self.storage_dir
        os.makedirs(target_dir, exist_ok=True)
        dest_path = os.path.join(target_dir, media_id)

        try:
            resp = self.session.get(url, stream=True, timeout=15)
            if resp.status_code == 200:
                with open(dest_path, "wb") as f:
                    for chunk in resp.iter_content(65536):
                        if chunk:
                            f.write(chunk)
                # Stash 照合
                img_id = self.stash.find_image_by_path(dest_path) if m_type == "image" else None
                scn_id = self.stash.find_scene_by_path(dest_path) if m_type != "image" else None
                return "COMPLETED", img_id, scn_id
            elif resp.status_code == 404:
                return "DEAD_404", None, None
        except Exception:
            pass
        return "QUEUED", None, None

    def _update_status(self, media_id: str, status: str, img_id: Optional[str], scn_id: Optional[str]) -> None:
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("""
                UPDATE media
                SET download_status = ?, stash_image_id = coalesce(?, stash_image_id),
                    stash_scene_id = coalesce(?, stash_scene_id)
                WHERE media_id = ?
            """, (status, img_id, scn_id, media_id))
            conn.commit()
