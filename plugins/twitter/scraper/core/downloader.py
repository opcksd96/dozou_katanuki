# plugins/twitter/scraper/core/downloader.py (100行以下)
import os, sqlite3, time
from typing import Optional, Tuple
import requests
from .aria2_client import Aria2Client
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
        self.aria2 = Aria2Client()

    def _get_target_path(self, username: str, media_id: str) -> str:
        target_dir = os.path.join(self.storage_dir, "Twitter", username) if "Influencers" in self.storage_dir else self.storage_dir
        os.makedirs(target_dir, exist_ok=True)
        return os.path.join(target_dir, media_id)

    def process_queued_media(self, article_id: Optional[str] = None, media_id: Optional[str] = None) -> int:
        """QUEUED メディアを順次ダウンロード (第1段階 -> 404時第2段階委託)"""
        conn = sqlite3.connect(self.db_path)
        try:
            query = ("SELECT m.media_id, m.download_url, m.type, ac.username FROM media m "
                     "JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id "
                     "WHERE m.download_status = 'QUEUED'")
            params = []
            if media_id:
                query += " AND m.media_id = ?"; params.append(media_id)
            elif article_id:
                query += " AND m.article_id = ?"; params.append(article_id)
            else:
                query += " LIMIT 20"
            records = conn.cursor().execute(query, params).fetchall()
        finally:
            conn.close()

        success = 0
        for m_id, url, m_type, user in records:
            dest = self._get_target_path(user or "unknown", m_id)
            status, reason, img_id, scn_id = self._try_download_and_escalate(m_id, url, m_type, dest)
            self._update_status(m_id, status, reason, img_id, scn_id)
            if status == "COMPLETED": success += 1
        return success

    def _try_download_and_escalate(self, media_id: str, url: str, m_type: str, dest: str) -> Tuple[str, Optional[str], Optional[str], Optional[str]]:
        """第1段階: requests 直接取得 ➔ 404時第2段階 Motrix/Aria2 委託"""
        for attempt in range(3):
            try:
                resp = self.session.get(url, stream=True, timeout=15)
                if resp.status_code == 200:
                    with open(dest, "wb") as f:
                        for chunk in resp.iter_content(65536):
                            if chunk: f.write(chunk)
                    return "COMPLETED", None, (self.stash.find_image_by_path(dest) if m_type == "image" else None), (self.stash.find_scene_by_path(dest) if m_type != "image" else None)
                elif resp.status_code == 404:
                    gid = self.aria2.add_uri([url], os.path.dirname(dest), os.path.basename(dest))
                    return ("OUTSOURCED", f"Delegated to Motrix (GID: {gid})", None, None) if gid else ("DEAD_404", "404 Not Found & Aria2 offline", None, None)
                time.sleep(0.5 * (2 ** attempt))
            except Exception:
                time.sleep(0.5 * (2 ** attempt))
        return "QUEUED", "Temporary network error", None, None

    def poll_outsourced_media(self) -> int:
        """第3段階: OUTSOURCED/RETAINED 実ファイルのポーリング検知・Stash回収"""
        conn = sqlite3.connect(self.db_path)
        try:
            records = conn.cursor().execute(
                "SELECT m.media_id, m.type, ac.username FROM media m "
                "JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id "
                "WHERE m.download_status IN ('OUTSOURCED', 'RETAINED')"
            ).fetchall()
        finally:
            conn.close()

        salvaged = 0
        for m_id, m_type, user in records:
            dest = self._get_target_path(user or "unknown", m_id)
            if os.path.exists(dest) and os.path.getsize(dest) > 0:
                img_id = self.stash.find_image_by_path(dest) if m_type == "image" else None
                scn_id = self.stash.find_scene_by_path(dest) if m_type != "image" else None
                self._update_status(m_id, "COMPLETED", None, img_id, scn_id)
                salvaged += 1
        return salvaged

    def _update_status(self, media_id: str, status: str, reason: Optional[str], img_id: Optional[str], scn_id: Optional[str]) -> None:
        conn = sqlite3.connect(self.db_path)
        try:
            conn.execute(
                "UPDATE media SET download_status = ?, failed_reason = ?, "
                "stash_image_id = coalesce(?, stash_image_id), stash_scene_id = coalesce(?, stash_scene_id) "
                "WHERE media_id = ?", (status, reason, img_id, scn_id, media_id)
            )
            conn.commit()
        finally:
            conn.close()
