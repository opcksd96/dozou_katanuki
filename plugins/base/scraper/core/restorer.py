# plugins/base/scraper/core/restorer.py (SPEC-PLUGIN-001 / 100行以下)
import concurrent.futures, json, os, shutil
from typing import Callable, Dict, Optional
from warcio.archiveiterator import ArchiveIterator
from .base_downloader import BaseDownloader
from .base_mutator import BaseMutator
from .base_parser import BaseParser


class Restorer:
    """Layer 2 (生データダンプ) からマルチスレッド並列で高速オフライン DB 再構築エンジン"""
    def __init__(self, dumps_dir: str = "backups/dumps", db_path: str = "archive.db", storage_dir: Optional[str] = None, avatar_dir: str = "assets/avatars", max_workers: int = 8, parser: Optional[BaseParser] = None, platform: str = "base"):
        self.dumps_dir, self.db_path, self.avatar_dir = dumps_dir, db_path, avatar_dir
        self.max_workers = max_workers
        self.mutator = BaseMutator(db_path=db_path, platform=platform, enable_translation=False)
        self.downloader = BaseDownloader(db_path=db_path, storage_dir=storage_dir, platform=platform)
        self.parser = parser

    def _sync_media_files(self) -> int:
        with self.downloader._get_conn() as conn:
            records = conn.cursor().execute("SELECT m.media_id, m.type, ac.username FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id WHERE m.download_status != 'COMPLETED'").fetchall()
        updates, synced = [], 0
        for m_id, m_type, user in records:
            dest = self.downloader.get_target_path(user or "unknown", m_id, m_type)
            if os.path.exists(dest) and os.path.getsize(dest) > 0:
                reg_id = self.downloader.reconciler.register_media(dest, m_type)
                if reg_id: updates.append(("COMPLETED", None, reg_id if m_type == "image" else None, reg_id if m_type != "image" else None, m_id)); synced += 1
                else: updates.append(("RETAINED", "Saved on disk", None, None, m_id))
        if updates:
            with self.downloader._get_conn() as conn:
                conn.cursor().executemany("UPDATE media SET download_status = ?, failed_reason = ?, stash_image_id = coalesce(?, stash_image_id), stash_scene_id = coalesce(?, stash_scene_id) WHERE media_id = ?", updates)
                conn.commit()
        return synced

    def _process_dir(self, pdir: str) -> Dict[str, int]:
        res = {"articles": 0, "avatars": 0}
        meta_path, warc_path = os.path.join(pdir, "metadata.json"), os.path.join(pdir, "snapshot.warc.gz")
        if os.path.exists(meta_path):
            try:
                with open(meta_path, "r", encoding="utf-8") as f:
                    if self.mutator.upsert_record(json.load(f)): res["articles"] += 1
            except Exception: pass
        elif os.path.exists(warc_path) and self.parser:
            try:
                with open(warc_path, "rb") as s:
                    for r in ArchiveIterator(s):
                        if r.rec_type != "response": continue
                        uri = r.rec_headers.get_header("WARC-Target-URI") or ""
                        c_type = (r.http_headers.get_header("Content-Type") if r.http_headers else "") or ""
                        if "json" in c_type or "html" in c_type or "status" in uri:
                            parsed = self.parser.parse_record(r.raw_stream.read(), uri)
                            if parsed and self.mutator.upsert_record(parsed): res["articles"] += 1; break
            except Exception: pass
        avatar_subdir = os.path.join(pdir, "avatars")
        if os.path.exists(avatar_subdir):
            for af in os.listdir(avatar_subdir):
                src_av = os.path.join(avatar_subdir, af)
                if os.path.isfile(src_av): shutil.copy2(src_av, os.path.join(self.avatar_dir, af)); res["avatars"] += 1
        return res

    def run_restore(self, progress_callback: Optional[Callable[[int, int, str], None]] = None) -> Dict[str, int]:
        if not os.path.exists(self.dumps_dir): return {"articles": 0, "media": 0, "avatars": 0}
        target_dirs = [root for root, _, files in os.walk(self.dumps_dir) if "metadata.json" in files or "snapshot.warc.gz" in files]
        total = len(target_dirs)
        if total == 0: return {"articles": 0, "media": 0, "avatars": 0}
        os.makedirs(self.avatar_dir, exist_ok=True)
        stats = {"articles": 0, "media": 0, "avatars": 0}
        if progress_callback: progress_callback(0, total, f"Restoring {total} dumps ({self.max_workers} threads)...")
        with concurrent.futures.ThreadPoolExecutor(max_workers=self.max_workers) as executor:
            futures = {executor.submit(self._process_dir, pdir): pdir for pdir in target_dirs}
            completed = 0
            for fut in concurrent.futures.as_completed(futures):
                completed += 1
                try:
                    r = fut.result(); stats["articles"] += r["articles"]; stats["avatars"] += r["avatars"]
                except Exception: pass
                if progress_callback and (completed % 200 == 0 or completed == total):
                    progress_callback(completed, total, f"Restored [{completed}/{total}] entries ({stats['articles']} articles)")
        stats["media"] = self._sync_media_files()
        if progress_callback: progress_callback(total, total, f"Restore done: {stats['articles']} arts, {stats['avatars']} avs, {stats['media']} media")
        return stats
