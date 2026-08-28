# plugins/base/scraper/core/media_queue_processor.py
import os
import sqlite3
from typing import Any, Callable, Optional
from .media_url_builder import MediaUrlBuilder


class MediaQueueProcessor:
    """DBのキューを走査してダウンロードとエスカレーションを実行する単一責務モジュール"""

    def __init__(self, downloader: Any):
        self.d = downloader

    def process_queued(
        self,
        article_id: Optional[str] = None,
        media_id: Optional[str] = None,
        log_fn: Optional[Callable[[int, int, str], None]] = None,
    ) -> int:
        with self.d._get_conn() as conn:
            cur = conn.cursor()
            wl = {r[0].lower() for r in cur.execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall() if r[0]}
            params = []
            where_clause = "WHERE m.download_status = 'QUEUED'"
            if media_id:
                where_clause += " AND m.media_id = ?"
                params.append(media_id)
            elif article_id:
                where_clause += " AND m.article_id = ?"
                params.append(article_id)

            sql = f"""
                SELECT m.media_id, m.download_url, m.type, ac.username, ac.display_name,
                       a.wayback_url, a.id, a.full_text, a.full_text_ja, a.created_at, m.thumbnail_url
                FROM media m
                JOIN articles a ON m.article_id = a.id
                JOIN accounts ac ON a.account_id = ac.numeric_id
                {where_clause}
            """
            records = cur.execute(sql, params).fetchall()

        total = len(records)
        success = 0
        if log_fn:
            log_fn(0, max(total, 1), f"Found {total} queued media items.")

        for idx, row in enumerate(records, start=1):
            m_id, url, m_type, user, dname, wb_url, art_id, f_text, f_text_ja, cr_at, thumb_url = row
            is_whitelisted = (not wl) or (user and user.lower() in wl)
            if not is_whitelisted:
                self.d._update_status(m_id, "EXCLUDED", "Whitelist外", None, None, None)
                continue

            dest_path = self.d.get_target_path(user or "unknown", m_id, m_type)
            u = self.d.resolve_media_url(m_id, url, m_type)
            meta = self.d.build_metadata(art_id or "", user or "", dname or "", str(cr_at or ""), wb_url or "", f_text or "", f_text_ja or "")

            # 1. 直接フェッチ試行 (全解像度 & Wayback)
            ok, st, img_id, scn_id, qual = self.d.fetcher.try_fetch_direct(
                dest_path, m_id, u, m_type, meta,
                thumbnail_url=thumb_url or "", created_at=str(cr_at or ""),
                username=user or "", display_name=dname or ""
            )

            if ok:
                self.d._update_status(m_id, "COMPLETED", None, img_id, scn_id, qual)
                success += 1
            else:
                # 2. Motrix (Aria2) への全候補URLマルチソース外注
                dest_dir = os.path.dirname(dest_path)
                base_name = m_id.split(":")[0]
                fallback_urls = MediaUrlBuilder.build_url_list(u)
                gid = self.d.aria2.add_uri(fallback_urls, dest_dir, base_name) if fallback_urls else None
                status = "OUTSOURCED" if gid else "RETAINED"
                reason = f"Motrix外注 (GID: {gid})" if gid else "原本消失・外注待機 (404)"
                self.d._update_status(m_id, status, reason, None, None, None)

            if log_fn:
                log_fn(idx, total, f"Media {m_id} -> {status if not ok else 'COMPLETED'}")

        return success
