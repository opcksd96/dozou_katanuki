# plugins/base/scraper/core/base_downloader.py (SPEC-PLUGIN-001 / 100行以下)
import os, sqlite3, requests
from typing import Any, Callable, Dict, List, Optional, Tuple
from .aria2_client import Aria2Client
from .stash_client import StashClient
from .stash_reconciler import StashReconciler
from .downloader_pipeline import DownloaderPipelineHelper

class BaseDownloader:
    """メディア原本取得 & Stash登録 & Motrix委託 & 品質追跡パイプライン (100行以下)"""
    DEFAULT_STORAGE = "G:/Media_Storage/Influencers" if os.path.exists("G:/Media_Storage/Influencers") else "blobs"

    def __init__(self, db_path: str = "archive.db", storage_dir: Optional[str] = None, platform: str = "twitter"):
        self.db_path, self.platform = db_path, platform
        self.storage_dir = storage_dir or self.DEFAULT_STORAGE
        os.makedirs(self.storage_dir, exist_ok=True); self.session = requests.Session()
        self.session.headers.update({"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"})
        self.stash, self.aria2 = StashClient(), Aria2Client()
        self.reconciler, self.pipeline_helper = StashReconciler(self.stash), DownloaderPipelineHelper(self)

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=60.0)
        conn.execute("PRAGMA journal_mode = WAL;"); conn.execute("PRAGMA synchronous = NORMAL;")
        conn.execute("PRAGMA foreign_keys = ON;"); conn.execute("PRAGMA busy_timeout = 60000;"); return conn

    def get_target_path(self, username: str, media_id: str, media_type: str = "image") -> str:
        plat, base_name = ("X(Twitter)" if self.platform.lower() in ("twitter", "base", "x") else self.platform.capitalize()), media_id.split(":")[0]
        td = os.path.join(self.storage_dir, username, plat, "_assets") if "Influencers" in self.storage_dir else (self.storage_dir if os.path.basename(os.path.normpath(self.storage_dir)) == "blobs" else os.path.join(self.storage_dir, "scenes" if media_type == "video" or media_id.endswith((".mp4", ".webm", ".m3u8")) else "images", plat, username))
        os.makedirs(td, exist_ok=True); return os.path.join(td, base_name)
    def resolve_media_url(self, media_id: str, download_url: str, media_type: str) -> str: return download_url
    def build_metadata(self, article_id: str, username: str, display_name: str, created_at: str, wayback_url: str, full_text: str, full_text_ja: str) -> Tuple[str, str, List[str], Dict[str, Any]]:
        t_title, txt = (f"{self.platform.capitalize()} (@{username}): Post {article_id}" if username and article_id else ""), (f"{full_text_ja}\n\n{full_text}".strip() if full_text_ja and full_text_ja != full_text else (full_text or full_text_ja or ""))
        urls = [wayback_url] if wayback_url else []
        return t_title, txt, urls, {"post_id": article_id, "original_name": display_name or username, "source_system": self.platform, "wayback_url": urls, "dead_media": []}
    def _build_fallback_urls(self, u: str) -> List[Tuple[str, str]]:
        if not u: return []
        base_u = u.split("?")[0]
        for sfx in [":orig", ":large", ":medium", ":small", ":thumb", ":tiny"]:
            if base_u.endswith(sfx): base_u = base_u[:-len(sfx)]; break
        base_no_ext, ext = os.path.splitext(base_u); fmt = ext.lstrip(".").lower() or "jpg"
        cands = [(u, "standard"), (f"https://web.archive.org/web/2id_/{base_no_ext}?format={fmt}&name=orig", "orig"), (f"https://web.archive.org/web/2id_/{base_u}:orig", "orig"), (f"https://web.archive.org/web/2id_/{base_u}", "standard")]
        seen, res = set(), []; [res.append((c, q)) for c, q in cands if c not in seen and not seen.add(c)]; return res
    def _is_valid_binary(self, p: str) -> bool:
        if not (os.path.exists(p) and os.path.getsize(p) > 0): return False
        try:
            with open(p, "rb") as f: h = f.read(32)
            return not (h.startswith(b"<!DOCTYPE") or h.startswith(b"<html") or h.startswith(b"<HTML") or h.startswith(b"{\n") or h.startswith(b'{"') or h.startswith(b"<?xml"))
        except Exception: return False

    def process_queued_media(self, article_id: Optional[str] = None, media_id: Optional[str] = None, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        with self._get_conn() as conn:
            wl = {r[0].lower() for r in conn.cursor().execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall() if r[0]}
            p, where = [], "WHERE m.download_status = 'QUEUED'"
            if media_id: where += " AND m.media_id = ?"; p.append(media_id)
            elif article_id: where += " AND m.article_id = ?"; p.append(article_id)
            records = conn.cursor().execute(f"SELECT m.media_id, m.download_url, m.type, ac.username, ac.display_name, a.wayback_url, a.id, a.full_text, a.full_text_ja, a.created_at, m.thumbnail_url FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id {where}", p).fetchall()
        total, success = len(records), 0; (log_fn and log_fn(0, max(total, 1), f"Found {total} queued media items."))
        for idx, (m_id, url, m_type, user, dname, wb_url, art_id, f_text, f_text_ja, cr_at, thumb_url) in enumerate(records, start=1):
            dest = self.get_target_path(user or "unknown", m_id, m_type)
            st, reason, img, scn, q = self._try_download_and_escalate(m_id, url, m_type, dest, not wl or (user and user.lower() in wl), art_id or "", user or "", dname or "", str(cr_at or ""), wb_url or "", f_text or "", f_text_ja or "", thumb_url or "")
            self._update_status(m_id, st, reason, img, scn, q); success += (1 if st == "COMPLETED" else 0); (log_fn and log_fn(idx, total, f"Media {m_id} -> {st} ({q or 'none'})"))
        return success

    def _try_download_and_escalate(self, media_id: str, url: str, m_type: str, dest: str, is_wl: bool = True, article_id: str = "", username: str = "", display_name: str = "", created_at: str = "", wayback_url: str = "", full_text: str = "", full_text_ja: str = "", thumbnail_url: str = "") -> Tuple[str, Optional[str], Optional[str], Optional[str], Optional[str]]:
        if not is_wl: return "EXCLUDED", "Whitelist外", None, None, None
        u, (t_title, txt, urls_list, c_fields) = self.resolve_media_url(media_id, url, m_type), self.build_metadata(article_id, username, display_name, created_at, wayback_url, full_text, full_text_ja)
        if self._is_valid_binary(dest):
            reg_id = self.reconciler.register_media(dest, m_type, title=t_title, details=txt, urls=urls_list, date=created_at[:10], username=username, display_name=display_name, thumbnail_url=thumbnail_url, custom_fields=c_fields)
            return ("COMPLETED", None, reg_id if m_type == "image" else None, reg_id if m_type != "image" else None, "local") if reg_id else ("RETAINED", "Saved to disk", None, None, "local")
        for cand_url, q_name in self._build_fallback_urls(u):
            try:
                resp = self.session.get(cand_url, stream=True, timeout=8 if m_type == "image" else 20, allow_redirects=True)
                if resp.status_code == 200 and not ("html" in resp.headers.get("Content-Type", "") and "text" in resp.headers.get("Content-Type", "")):
                    with open(dest, "wb") as f:
                        for chunk in resp.iter_content(65536): (f.write(chunk) if chunk else None)
                    if self._is_valid_binary(dest):
                        reg_id = self.reconciler.register_media(dest, m_type, title=t_title, details=txt, urls=urls_list, date=created_at[:10], username=username, display_name=display_name, thumbnail_url=thumbnail_url, custom_fields=c_fields)
                        return ("COMPLETED", None, reg_id if m_type == "image" else None, reg_id if m_type != "image" else None, q_name) if reg_id else ("RETAINED", "Saved to disk", None, None, q_name)
                    elif os.path.exists(dest): os.remove(dest)
            except Exception: pass
        esc_st, esc_r, esc_img, esc_scn = self._escalate_single(media_id, u, m_type, username=username)
        return esc_st, esc_r, esc_img, esc_scn, None

    def _escalate_single(self, media_id: str, u: str, m_type: str, username: str = "unknown") -> Tuple[str, Optional[str], Optional[str], Optional[str]]:
        dest_dir, base_name = os.path.dirname(self.get_target_path(username or "unknown", media_id, m_type)), media_id.split(":")[0]
        fallback_urls = [cand for cand, _ in self._build_fallback_urls(u)]
        gid = self.aria2.add_uri(fallback_urls, dest_dir, base_name) if fallback_urls else None
        return ("OUTSOURCED", f"Motrix外注 (GID: {gid})", None, None) if gid else ("RETAINED", "原本消失・外注待機 (404)", None, None)

    def _update_status(self, media_id: str, status: str, reason: Optional[str] = None, img_id: Optional[str] = None, scn_id: Optional[str] = None, quality: Optional[str] = None) -> None:
        with self._get_conn() as conn:
            conn.cursor().execute("UPDATE media SET download_status = ?, failed_reason = coalesce(?, failed_reason), stash_image_id = coalesce(?, stash_image_id), stash_scene_id = coalesce(?, stash_scene_id), media_quality = coalesce(?, media_quality) WHERE media_id = ?", (status, reason, img_id, scn_id, quality, media_id)); conn.commit()
    def escalate_dead_media(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int: return self.pipeline_helper.escalate_dead_media(log_fn=log_fn)
    def poll_outsourced_media(self, log_fn: Optional[Callable[[str], None]] = None) -> int: return self.reconciler.reconcile_to_db(self.db_path, log_fn=log_fn)
    def smart_recovery_pipeline(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> dict: return self.pipeline_helper.run_smart_recovery(log_fn=log_fn)
    def escalate_to_thunder(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int: return self.pipeline_helper.escalate_to_thunder(log_fn=log_fn)
    def clean_failed_outsourced(self, log_fn: Optional[Callable[[str], None]] = None) -> int: return self.pipeline_helper.clean_failed_outsourced(log_fn=log_fn)
