# plugins/base/scraper/core/downloader_pipeline.py (SPEC-PLUGIN-001 / 100行以下)
import os, sqlite3
from typing import Any, Callable, List, Optional
from .thunder_client import ThunderClient

class DownloaderPipelineHelper:
    """メディア回収パイプライン・エスカレーション・Thunder連携支援 (100行以下)"""
    def __init__(self, downloader: Any = None):
        self.d, self.thunder = downloader, ThunderClient()

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
            gids = []
            if m_type == "image":
                mirrors = [u]
                if not u.startswith("http://web.archive") and not u.startswith("https://web.archive"):
                    mirrors.append(f"https://web.archive.org/web/2id_/{u}")
                gid = self.d.aria2.add_uri(mirrors, dest_dir, base_name)
                if gid: gids.append(gid)
            else:
                variants = conn.cursor().execute("SELECT download_url FROM media_variants WHERE media_id = ? ORDER BY bit_rate DESC", (m_id,)).fetchall()
                f_urls = [v[0] for v in variants if v[0]] or ([u] if u else [])
                for v_url in f_urls:
                    v_fn = v_url.split("?")[0].split("/")[-1]
                    for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: v_fn = v_fn[:-len(sfx)] if v_fn.endswith(sfx) else v_fn
                    gid = self.d.aria2.add_uri([v_url], dest_dir, v_fn.split(".")[0] if "." in v_fn else v_fn)
                    if gid: gids.append(gid)
            self.d._update_status(m_id, "OUTSOURCED" if gids else "RETAINED", f"Motrix外注 (GID: {gids[0]})" if (m_type == "image" and gids) else (f"Motrix個別外注 (GIDs: {','.join(gids[:2])})" if gids else "Aria2 offline/Empty"))
            outsourced += (1 if gids else 0); (log_fn and log_fn(idx, total, f"{m_id} -> {'OUTSOURCED' if gids else 'RETAINED'}"))
        return outsourced

    def clean_failed_outsourced(self, log_fn: Optional[Callable[[str], None]] = None) -> int:
        if not hasattr(self.d, "aria2"): return 0
        self.d.aria2.purge_failed_tasks(); queued = self.d.aria2.get_queued_filenames()
        with self.d._get_conn() as conn:
            cur = conn.cursor(); rows = cur.execute("SELECT m.media_id, ac.username, m.type FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id WHERE m.download_status = 'OUTSOURCED'").fetchall()
            to_ret = [m_id for m_id, user, m_type in rows if not (os.path.exists(self.d.get_target_path(user or "unknown", m_id, m_type)) and os.path.getsize(self.d.get_target_path(user or "unknown", m_id, m_type)) > 0) and m_id not in queued]
            if to_ret:
                cur.executemany("UPDATE media SET download_status = 'RETAINED', failed_reason = 'Motrix未完了・キュー不在 (404/Timeout)' WHERE media_id = ?", [(x,) for x in to_ret])
                conn.commit()
        if log_fn: log_fn(f"[CLEAN] Reverted {len(to_ret)} failed/orphaned Motrix tasks to RETAINED.")
        return len(to_ret)

    def escalate_to_thunder(self, log_fn: Optional[Callable[[int, int, str], None]] = None, max_batch: int = 50) -> int:
        if not self.thunder.is_available(): (log_fn and log_fn(0, 0, "Thunder.exe not found on system.")); return 0
        with self.d._get_conn() as conn:
            cur = conn.cursor(); records = cur.execute("SELECT m.media_id, m.download_url, m.type, ac.username FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id WHERE m.download_status = 'RETAINED' LIMIT ?", (max_batch,)).fetchall()
            total, tasks, m_ids = len(records), [], []
            for m_id, url, m_type, user in records:
                u, dest_dir = self.d.resolve_media_url(m_id, url, m_type), os.path.dirname(self.d.get_target_path(user or "unknown", m_id, m_type))
                if m_type == "image":
                    wb_u = f"https://web.archive.org/web/2id_/{u}" if not u.startswith("http://web.archive") and not u.startswith("https://web.archive") else u
                    tasks.append({"url": wb_u, "file_name": m_id.split(":")[0], "dest_dir": dest_dir})
                else:
                    variants = cur.execute("SELECT download_url FROM media_variants WHERE media_id = ? ORDER BY bit_rate DESC", (m_id,)).fetchall()
                    for v_u in ([v[0] for v in variants if v[0]] or ([u] if u else [])):
                        v_fn = v_u.split("?")[0].split("/")[-1]
                        for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: v_fn = v_fn[:-len(sfx)] if v_fn.endswith(sfx) else v_fn
                        tasks.append({"url": f"https://web.archive.org/web/2id_/{v_u}" if not v_u.startswith("http://web.archive") and not v_u.startswith("https://web.archive") else v_u, "file_name": v_fn.split(".")[0] if "." in v_fn else v_fn, "dest_dir": dest_dir})
                m_ids.append(m_id)
            sent = self.thunder.add_batch_tasks(tasks, max_limit=max_batch)
            if sent > 0 and m_ids:
                cur.executemany("UPDATE media SET download_status = 'ESCALATED', failed_reason = 'Thunder P2SP エスカレーション投入' WHERE media_id = ?", [(mid,) for mid in m_ids]); conn.commit()
        if log_fn: log_fn(sent, max(total, 1), f"Escalated {sent}/{total} media items to Thunder (ESCALATED).")
        return sent

    def run_smart_recovery(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> dict:
        return {"stage1_salvaged": self.d.process_queued_media(log_fn=log_fn), "stage2_outsourced": self.d.escalate_dead_media(log_fn=log_fn), "stage3_reconciled": self.d.poll_outsourced_media(), "failed_cleaned": self.clean_failed_outsourced(log_fn=lambda m: (log_fn(0, 0, m) if log_fn else None))}
