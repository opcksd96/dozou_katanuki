# plugins/twitter/scraper/core/downloader.py (100行以下)
import os, re, sqlite3, time
from typing import Optional, Tuple
import requests
from .aria2_client import Aria2Client
from .stash_client import StashClient


class Downloader:
    """メディア原本ストリーム取得 & Stash登録 & Motrix委託エンジン (SPEC-PLUGIN-001)"""
    DEFAULT_STORAGE = "G:/Media_Storage/Influencers" if os.path.exists("G:/Media_Storage/Influencers") else "blobs"

    def __init__(self, db_path: str = "archive.db", storage_dir: Optional[str] = None):
        self.db_path = db_path
        self.storage_dir = storage_dir or self.DEFAULT_STORAGE
        os.makedirs(self.storage_dir, exist_ok=True)
        self.session, self.stash, self.aria2 = requests.Session(), StashClient(), Aria2Client()
        self.session.headers.update({"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"})

    def _get_target_path(self, username: str, media_id: str, media_type: str = "image", platform: str = "twitter") -> str:
        if "Influencers" in self.storage_dir: td = os.path.join(self.storage_dir, username, "X(Twitter)", "_assets")
        elif os.path.basename(os.path.normpath(self.storage_dir)) == "blobs": td = self.storage_dir
        else: td = os.path.join(self.storage_dir, "scenes" if media_type == "video" or media_id.endswith((".mp4", ".webm", ".m3u8")) else "images", platform, username)
        os.makedirs(td, exist_ok=True)
        return os.path.join(td, media_id)

    def process_queued_media(self, article_id: Optional[str] = None, media_id: Optional[str] = None) -> int:
        with sqlite3.connect(self.db_path) as conn:
            wl = set()
            try: wl = {r[0].lower() for r in conn.cursor().execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall() if r[0]}
            except Exception: pass
            p = []
            where_clause = "WHERE m.download_status = 'QUEUED'"
            if media_id: where_clause += " AND m.media_id = ?"; p.append(media_id)
            elif article_id: where_clause += " AND m.article_id = ?"; p.append(article_id)
            try:
                q = f"SELECT m.media_id, m.download_url, m.type, ac.username, a.wayback_url FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id {where_clause}"
                records = [(r[0], r[1], r[2], r[3], r[4]) for r in conn.cursor().execute(q, p).fetchall()]
            except Exception:
                q = f"SELECT m.media_id, m.download_url, m.type, ac.username FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id {where_clause}"
                records = [(r[0], r[1], r[2], r[3], "") for r in conn.cursor().execute(q, p).fetchall()]

        success = 0
        for m_id, url, m_type, user, wb_url in records:
            is_wl = bool(not wl or (user and user.lower() in wl))
            dest = self._get_target_path(user or "unknown", m_id, m_type)
            m_ts = re.search(r'/web/(\d{14})', wb_url) if wb_url else None
            st, reason, img, scn = self._try_download_and_escalate(m_id, url, m_type, dest, is_wl, article_id or "", user or "", m_ts.group(1) if m_ts else "")
            self._update_status(m_id, st, reason, img, scn)
            if st == "COMPLETED": success += 1
        return success

    def _try_download_and_escalate(self, media_id: str, url: str, m_type: str, dest: str, is_wl: bool = True, article_id: str = "", username: str = "", wayback_ts: str = "") -> Tuple[str, Optional[str], Optional[str], Optional[str]]:
        if not is_wl: return "EXCLUDED", "Whitelist外 (ダウンロード対象外)", None, None
        u = url or (f"https://pbs.twimg.com/media/{media_id}" if m_type == "image" else f"https://video.twimg.com/ext_tw_video/{media_id}")
        t_title = f"X (@{username}): Tweet {article_id}" if username and article_id else ""
        if os.path.exists(dest) and os.path.getsize(dest) > 0:
            img_id = self.stash.register_media(dest, "image", title=t_title, url=u) if m_type == "image" else None
            scn_id = self.stash.register_media(dest, "video", title=t_title, url=u) if m_type != "image" else None
            if img_id or scn_id:
                return "COMPLETED", None, img_id, scn_id
            return "RETAINED", "Saved to disk, awaiting Stash index", None, None

        base_u = u.rsplit(":", 1)[0] if any(u.endswith(s) for s in [":large", ":orig", ":small", ":medium"]) else u
        http_u = base_u.replace("https://", "http://")
        targets = []
        if wayback_ts: targets.extend([f"https://web.archive.org/web/{wayback_ts}im_/{base_u}", f"https://web.archive.org/web/{wayback_ts}im_/{http_u}"])
        if "web.archive.org" in u: targets.append(u)
        else: targets.extend([u, f"https://web.archive.org/web/2/{base_u}", f"https://web.archive.org/web/2/{http_u}"])

        for t_url in targets:
            for attempt in range(2):
                try:
                    resp = self.session.get(t_url, stream=True, timeout=6, allow_redirects=True)
                    if resp.status_code == 200:
                        ct = resp.headers.get("Content-Type", "")
                        if "html" in ct and "text" in ct: break
                        with open(dest, "wb") as f:
                            for chunk in resp.iter_content(65536):
                                if chunk: f.write(chunk)
                        if os.path.exists(dest) and os.path.getsize(dest) > 0:
                            img_id = self.stash.register_media(dest, "image", title=t_title, url=u) if m_type == "image" else None
                            scn_id = self.stash.register_media(dest, "video", title=t_title, url=u) if m_type != "image" else None
                            if img_id or scn_id:
                                return "COMPLETED", None, img_id, scn_id
                            return "RETAINED", "Saved to disk, awaiting Stash index", None, None
                    elif resp.status_code in (403, 404, 410): break
                    elif resp.status_code == 429: time.sleep(1.0)
                except Exception: time.sleep(0.2)

        gid = self.aria2.add_uri(targets, os.path.dirname(dest), os.path.basename(dest))
        return ("OUTSOURCED", f"Delegated (GID: {gid})", None, None) if gid else ("DEAD_404", "404 & Wayback missed & Aria2 offline", None, None)

    def poll_outsourced_media(self) -> int:
        with sqlite3.connect(self.db_path) as conn:
            records = conn.cursor().execute("SELECT m.media_id, m.type, ac.username FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id WHERE m.download_status IN ('OUTSOURCED', 'RETAINED')").fetchall()
        salvaged = 0
        for m_id, m_type, user in records:
            dest = self._get_target_path(user or "unknown", m_id, m_type)
            if os.path.exists(dest) and os.path.getsize(dest) > 0:
                img = self.stash.register_media(dest, "image") if m_type == "image" else None
                scn = self.stash.register_media(dest, "video") if m_type != "image" else None
                if img or scn:
                    self._update_status(m_id, "COMPLETED", None, img, scn)
                    salvaged += 1
        return salvaged

    def _update_status(self, media_id: str, status: str, reason: Optional[str], img_id: Optional[str], scn_id: Optional[str]) -> None:
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("UPDATE media SET download_status = ?, failed_reason = ?, stash_image_id = coalesce(?, stash_image_id), stash_scene_id = coalesce(?, stash_scene_id) WHERE media_id = ?",
                         (status, reason, img_id, scn_id, media_id))
            conn.commit()
