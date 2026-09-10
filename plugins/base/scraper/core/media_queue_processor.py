# plugins/base/scraper/core/media_queue_processor.py (SPEC-PLUGIN-001 / 100行以下)
import os, sqlite3
from typing import Any, Callable, Optional

class MediaQueueProcessor:
    """DBキュー走査・Image/Scene完全分離外注・Stash登録を実行するプロセッサ"""
    def __init__(self, downloader: Any):
        self.d = downloader

    def process_queued(self, article_id: Optional[str] = None, media_id: Optional[str] = None, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        with self.d._get_conn() as conn:
            cur = conn.cursor()
            wl = {r[0].lower() for r in cur.execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall() if r[0]}
            params, where = [], "WHERE m.download_status = 'QUEUED'"
            if media_id: where += " AND m.media_id = ?"; params.append(media_id)
            elif article_id: where += " AND m.article_id = ?"; params.append(article_id)
            sql = f"SELECT m.media_id, m.download_url, m.type, ac.username, ac.display_name, a.wayback_url, a.id, a.full_text, a.full_text_ja, a.created_at, m.thumbnail_url FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id {where}"
            records = cur.execute(sql, params).fetchall()

        total, success = len(records), 0
        if log_fn: log_fn(0, max(total, 1), f"Found {total} queued media items.")
        for idx, row in enumerate(records, start=1):
            m_id, url, m_type, user, dname, wb_url, art_id, f_text, f_text_ja, cr_at, thumb_url = row
            if wl and (not user or user.lower() not in wl):
                self.d._update_status(m_id, "EXCLUDED", "Whitelist外", None, None, None); continue
            dest_path = self.d.get_target_path(user or "unknown", m_id, m_type)
            u = self.d.resolve_media_url(m_id, url, m_type)
            meta = self.d.build_metadata(art_id or "", user or "", dname or "", str(cr_at or ""), wb_url or "", f_text or "", f_text_ja or "")

            variants = cur.execute("SELECT download_url FROM media_variants WHERE media_id = ? ORDER BY bit_rate DESC", (m_id,)).fetchall()
            v_urls = [v[0] for v in variants if v[0]] or ([u] if u else [])
            fallback_urls = []
            for v_u in v_urls:
                fallback_urls.append(v_u)
                if not v_u.startswith("http://web.archive") and not v_u.startswith("https://web.archive"):
                    fallback_urls.append(f"https://web.archive.org/web/2id_/{v_u}")

            primary_u = v_urls[0] if v_urls else u
            ok, st, img_id, scn_id, qual = self.d.fetcher.try_fetch_direct(
                dest_path, m_id, primary_u, m_type, meta, thumbnail_url=thumb_url or "", created_at=str(cr_at or ""),
                username=user or "", display_name=dname or "", fallback_urls=fallback_urls
            )
            if ok:
                self.d._update_status(m_id, "COMPLETED", None, img_id, scn_id, qual); success += 1
            else:
                dest_dir, base_name = os.path.dirname(dest_path), m_id.split(":")[0]
                gids = []
                if m_type == "image":
                    for v_u in v_urls:
                        v_fn = v_u.split("?")[0].split("/")[-1]
                        for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: v_fn = v_fn[:-len(sfx)] if v_fn.endswith(sfx) else v_fn
                        mirrors = [v_u] + ([f"https://web.archive.org/web/2id_/{v_u}"] if not v_u.startswith("http://web.archive") else [])
                        gid = self.d.aria2.add_uri(mirrors, dest_dir, base_name)
                        if gid: gids.append(gid)
                else:
                    for v_u in v_urls:
                        v_fn = v_u.split("?")[0].split("/")[-1]
                        for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: v_fn = v_fn[:-len(sfx)] if v_fn.endswith(sfx) else v_fn
                        gid = self.d.aria2.add_uri([v_u], dest_dir, v_fn.split(".")[0] if "." in v_fn else v_fn)
                        if gid: gids.append(gid)
                status = "OUTSOURCED" if gids else "RETAINED"
                reason = f"Motrix外注 (GID: {gids[0]})" if gids else "原本消失・外注待機 (404)"
                self.d._update_status(m_id, status, reason, None, None, None)
            if log_fn: log_fn(idx, total, f"Media {m_id} -> {'COMPLETED' if ok else status}")
        return success
