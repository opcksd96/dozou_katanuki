# plugins/base/scraper/core/base_downloader.py (SPEC-PLUGIN-001 / 100行以下)
import os, sqlite3, requests
from typing import Any, Callable, Dict, List, Optional, Tuple
from .aria2_client import Aria2Client
from .stash_client import StashClient
from .stash_reconciler import StashReconciler


class BaseDownloader:
    """メディア原本ストリーム取得 & Stash登録 & Motrix委託 3段階確保共通パイプライン"""
    DEFAULT_STORAGE = "G:/Media_Storage/Influencers" if os.path.exists("G:/Media_Storage/Influencers") else "blobs"

    def __init__(self, db_path: str = "archive.db", storage_dir: Optional[str] = None, platform: str = "base"):
        self.db_path, self.platform = db_path, platform
        self.storage_dir = storage_dir or self.DEFAULT_STORAGE
        os.makedirs(self.storage_dir, exist_ok=True)
        self.session = requests.Session()
        self.session.headers.update({"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"})
        self.stash, self.aria2 = StashClient(), Aria2Client()
        self.reconciler = StashReconciler(self.stash)

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=60.0)
        conn.execute("PRAGMA journal_mode = WAL;"); conn.execute("PRAGMA synchronous = NORMAL;")
        conn.execute("PRAGMA foreign_keys = ON;"); conn.execute("PRAGMA busy_timeout = 60000;")
        return conn

    def get_target_path(self, username: str, media_id: str, media_type: str = "image") -> str:
        if "Influencers" in self.storage_dir: td = os.path.join(self.storage_dir, username, self.platform.capitalize(), "_assets")
        elif os.path.basename(os.path.normpath(self.storage_dir)) == "blobs": td = self.storage_dir
        else: td = os.path.join(self.storage_dir, "scenes" if media_type == "video" or media_id.endswith((".mp4", ".webm", ".m3u8")) else "images", self.platform, username)
        os.makedirs(td, exist_ok=True); return os.path.join(td, media_id)

    def resolve_media_url(self, media_id: str, download_url: str, media_type: str) -> str: return download_url

    def build_metadata(self, article_id: str, username: str, display_name: str, created_at: str, wayback_url: str, full_text: str, full_text_ja: str) -> Tuple[str, str, List[str], Dict[str, Any]]:
        t_title = f"{self.platform.capitalize()} (@{username}): Post {article_id}" if username and article_id else ""
        txt = f"{full_text_ja}\n\n{full_text}".strip() if full_text_ja and full_text_ja != full_text else (full_text or full_text_ja or "")
        urls = [wayback_url] if wayback_url else []
        return t_title, txt, urls, {"post_id": article_id, "original_name": display_name or username, "source_system": self.platform, "wayback_url": urls, "dead_media": []}

    def process_queued_media(self, article_id: Optional[str] = None, media_id: Optional[str] = None, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        def _log(c: int, t: int, m: str) -> None:
            if log_fn: log_fn(c, t, m)
            print(f"[DOWNLOAD] [{c}/{t}] {m}", flush=True)

        with self._get_conn() as conn:
            wl = {r[0].lower() for r in conn.cursor().execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall() if r[0]}
            p, where = [], "WHERE m.download_status = 'QUEUED'"
            if media_id: where += " AND m.media_id = ?"; p.append(media_id)
            elif article_id: where += " AND m.article_id = ?"; p.append(article_id)
            records = conn.cursor().execute(f"SELECT m.media_id, m.download_url, m.type, ac.username, ac.display_name, a.wayback_url, a.id, a.full_text, a.full_text_ja, a.created_at FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id {where}", p).fetchall()

        total = len(records); _log(0, max(total, 1), f"Found {total} queued media items awaiting download.")
        success = 0
        for idx, (m_id, url, m_type, user, dname, wb_url, art_id, f_text, f_text_ja, cr_at) in enumerate(records, start=1):
            dest = self.get_target_path(user or "unknown", m_id, m_type)
            st, reason, img, scn = self._try_download_and_escalate(m_id, url, m_type, dest, not wl or (user and user.lower() in wl), art_id or "", user or "", dname or "", str(cr_at or ""), wb_url or "", f_text or "", f_text_ja or "")
            self._update_status(m_id, st, reason, img, scn)
            if st == "COMPLETED": success += 1
            _log(idx, total, f"Media {m_id} -> {st} ({reason or 'OK' if st != 'COMPLETED' else 'Injected to Stash'})")
        _log(total, total, f"Download batch completed: {success}/{total} successfully salvaged.")
        return success

    def _try_download_and_escalate(self, media_id: str, url: str, m_type: str, dest: str, is_wl: bool = True,
                                   article_id: str = "", username: str = "", display_name: str = "", created_at: str = "",
                                   wayback_url: str = "", full_text: str = "", full_text_ja: str = "") -> Tuple[str, Optional[str], Optional[str], Optional[str]]:
        if not is_wl: return "EXCLUDED", "Whitelist外", None, None
        u = self.resolve_media_url(media_id, url, m_type)
        t_title, txt, urls_list, c_fields = self.build_metadata(article_id, username, display_name, created_at, wayback_url, full_text, full_text_ja)
        if os.path.exists(dest) and os.path.getsize(dest) > 0:
            reg_id = self.reconciler.register_media(dest, m_type, title=t_title, details=txt, urls=urls_list, date=created_at[:10], username=username, display_name=display_name, custom_fields=c_fields)
            return ("COMPLETED", None, reg_id if m_type == "image" else None, reg_id if m_type != "image" else None) if reg_id else ("RETAINED", "Saved to disk", None, None)
        try:
            to = 10 if m_type == "image" else 30
            resp = self.session.get(u, stream=True, timeout=to, allow_redirects=True)
            if resp.status_code == 200 and not ("html" in resp.headers.get("Content-Type", "") and "text" in resp.headers.get("Content-Type", "")):
                with open(dest, "wb") as f:
                    for chunk in resp.iter_content(65536): f.write(chunk) if chunk else None
                if os.path.exists(dest) and os.path.getsize(dest) > 0:
                    reg_id = self.reconciler.register_media(dest, m_type, title=t_title, details=txt, urls=urls_list, date=created_at[:10], username=username, display_name=display_name, custom_fields=c_fields)
                    return ("COMPLETED", None, reg_id if m_type == "image" else None, reg_id if m_type != "image" else None) if reg_id else ("RETAINED", "Saved to disk", None, None)
            elif resp.status_code in (403, 404): return ("DEAD_404", f"原本消失 ({resp.status_code})", None, None)
            elif resp.status_code in (408, 429, 500, 502, 503, 504): return ("QUEUED", f"一時障害リトライ待機 ({resp.status_code})", None, None)
        except Exception as e: return ("DEAD_404", f"取得例外: {e}", None, None)
        return ("QUEUED", "リトライ待機", None, None)

    def escalate_dead_media(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        def _log(c: int, t: int, m: str) -> None:
            if log_fn: log_fn(c, t, m)
            print(f"[ESCALATE] [{c}/{t}] {m}", flush=True)

        with self._get_conn() as conn:
            records = conn.cursor().execute("SELECT m.media_id, m.download_url, m.type, a.wayback_url FROM media m JOIN articles a ON m.article_id = a.id WHERE m.download_status = 'DEAD_404'").fetchall()

        total = len(records); _log(0, max(total, 1), f"Found {total} DEAD_404 media awaiting escalation.")
        outsourced = 0
        for idx, (m_id, url, m_type, _wb) in enumerate(records, start=1):
            u = self.resolve_media_url(m_id, url, m_type)
            dest_dir = os.path.dirname(self.get_target_path("_escalate", m_id, m_type))
            wb_media_url = f"https://web.archive.org/web/2id_/{u}" if u else ""
            fallback_urls = [x for x in [u, wb_media_url] if x]
            gid = self.aria2.add_uri(fallback_urls, dest_dir, m_id) if fallback_urls else None
            if gid:
                self._update_status(m_id, "OUTSOURCED", f"Motrix外注 (GID: {gid})"); outsourced += 1
                _log(idx, total, f"{m_id} -> OUTSOURCED (GID: {gid})")
            else: _log(idx, total, f"{m_id} -> Aria2 offline, DEAD_404 維持")
        _log(total, total, f"Escalation completed: {outsourced}/{total} outsourced to Motrix.")
        return outsourced

    def _update_status(self, media_id: str, status: str, reason: Optional[str] = None, img_id: Optional[str] = None, scn_id: Optional[str] = None) -> None:
        with self._get_conn() as conn:
            conn.cursor().execute("UPDATE media SET download_status = ?, failed_reason = coalesce(?, failed_reason), stash_image_id = coalesce(?, stash_image_id), stash_scene_id = coalesce(?, stash_scene_id) WHERE media_id = ?", (status, reason, img_id, scn_id, media_id))
            conn.commit()

    def poll_outsourced_media(self, log_fn: Optional[Callable[[str], None]] = None) -> int:
        return self.reconciler.reconcile_to_db(self.db_path, log_fn=log_fn)
