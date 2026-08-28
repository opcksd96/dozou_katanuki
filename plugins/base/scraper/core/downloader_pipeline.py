# plugins/base/scraper/core/downloader_pipeline.py (SPEC-PLUGIN-001 / 100行以下)
import os, sqlite3
from typing import Any, Callable, List, Optional
from .thunder_client import ThunderClient


class DownloaderPipelineHelper:
    """メディア回収パイプラインの高度な統合・フォールバック・Thunder連携支援 (100行以下)"""
    def __init__(self, downloader: Any = None):
        self.d = downloader
        self.thunder = ThunderClient()

    def escalate_dead_media(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        with self.d._get_conn() as conn:
            records = conn.cursor().execute("SELECT m.media_id, m.download_url, m.type, a.wayback_url, ac.username FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id WHERE m.download_status = 'DEAD_404'").fetchall()
        total, outsourced, existing = len(records), 0, self.d.aria2.get_queued_filenames()
        if log_fn: log_fn(0, max(total, 1), f"Found {total} dead_404 media. Motrix cached: {len(existing)} items.")
        for idx, (m_id, url, m_type, _wb, user) in enumerate(records, start=1):
            base_name = m_id.split(":")[0]
            if base_name in existing or m_id in existing:
                self.d._update_status(m_id, "OUTSOURCED", "Motrix既存キュー確認 (送信スキップ)"); continue
            u, dest_dir = self.d.resolve_media_url(m_id, url, m_type), os.path.dirname(self.d.get_target_path(user or "unknown", m_id, m_type))
            fallback_urls = [cand for cand, _ in self.d._build_fallback_urls(u)]
            gid = self.d.aria2.add_uri(fallback_urls, dest_dir, base_name) if fallback_urls else None
            self.d._update_status(m_id, "OUTSOURCED" if gid else "RETAINED", f"Motrix外注 (GID: {gid})" if gid else "Aria2 offline")
            outsourced += (1 if gid else 0); (log_fn and log_fn(idx, total, f"{m_id} -> {'OUTSOURCED' if gid else 'RETAINED'}"))
        return outsourced

    def clean_failed_outsourced(self, log_fn: Optional[Callable[[str], None]] = None) -> int:
        """Motrixでエラーまたはキューから消失したタスクを検知し、DBステータスを確実にRETAINEDへ移行"""
        if not hasattr(self.d, "aria2"): return 0
        self.d.aria2.purge_failed_tasks()
        queued_in_motrix = self.d.aria2.get_queued_filenames()
        with self.d._get_conn() as conn:
            cur = conn.cursor()
            rows = cur.execute("SELECT m.media_id, ac.username, m.type FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id WHERE m.download_status = 'OUTSOURCED'").fetchall()
            to_retained = []
            for m_id, user, m_type in rows:
                target_path = self.d.get_target_path(user or "unknown", m_id, m_type)
                if not (os.path.exists(target_path) and os.path.getsize(target_path) > 0) and m_id not in queued_in_motrix:
                    to_retained.append(m_id)
            if to_retained:
                cur.executemany("UPDATE media SET download_status = 'RETAINED', failed_reason = 'Motrix未完了・キュー不在 (404/Timeout)' WHERE media_id = ?", [(x,) for x in to_retained])
                conn.commit()
        if log_fn: log_fn(f"[CLEAN] Reverted {len(to_retained)} failed/orphaned Motrix tasks to RETAINED.")
        return len(to_retained)

    def escalate_to_thunder(self, log_fn: Optional[Callable[[int, int, str], None]] = None, max_batch: int = 50) -> int:
        if not self.thunder.is_available():
            if log_fn: log_fn(0, 0, "Thunder.exe not found on system."); return 0
        with self.d._get_conn() as conn:
            records = conn.cursor().execute("SELECT m.media_id, m.download_url, m.type FROM media m WHERE m.download_status = 'RETAINED' LIMIT ?", (max_batch,)).fetchall()
        total, urls_to_send = len(records), []
        for m_id, url, m_type in records:
            u = self.d.resolve_media_url(m_id, url, m_type)
            wb_media_url = f"https://web.archive.org/web/2id_/{u}" if u and not u.startswith("http://web.archive") and not u.startswith("https://web.archive") else u
            urls_to_send.append(wb_media_url or u)
        sent = self.thunder.add_batch_tasks([u for u in urls_to_send if u], max_limit=max_batch)
        if log_fn: log_fn(sent, total, f"Sent {sent}/{total} direct media URLs to Thunder.exe")
        return sent

    def run_smart_recovery(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> dict:
        r1 = self.d.process_queued_media(log_fn=log_fn)
        r2 = self.d.escalate_dead_media(log_fn=log_fn)
        r3 = self.d.poll_outsourced_media()
        r4 = self.clean_failed_outsourced(log_fn=lambda m: (log_fn(0, 0, m) if log_fn else None))
        return {"stage1_salvaged": r1, "stage2_outsourced": r2, "stage3_reconciled": r3, "failed_cleaned": r4}
