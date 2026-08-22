# plugins/twitter/scraper/core/restorer.py (100行以下)
import glob
import json
import os
import shutil
import sqlite3
from typing import Any, Callable, Dict, Optional
from .downloader import Downloader
from .mutator import Mutator
from .warc_importer import WarcImporter


class Restorer:
    """Layer 2 (生データダンプ) から完全オフラインで DB を100%自動再構築 (SPEC-RECOVERY-001)"""

    def __init__(self, dumps_dir: str = "backups/dumps", db_path: str = "archive.db", storage_dir: Optional[str] = None, avatar_dir: str = "assets/avatars"):
        self.dumps_dir = dumps_dir
        self.db_path = db_path
        self.avatar_dir = avatar_dir
        self.mutator = Mutator(db_path=db_path)
        self.downloader = Downloader(db_path=db_path, storage_dir=storage_dir)

    def _sync_media_files(self) -> int:
        """ストレージまたは dump ディレクトリ内の既存メディアファイルを検知し DB を COMPLETED へ同期"""
        conn = sqlite3.connect(self.db_path)
        try:
            records = conn.cursor().execute(
                "SELECT m.media_id, m.type, ac.username FROM media m "
                "JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id "
                "WHERE m.download_status != 'COMPLETED'"
            ).fetchall()
        finally:
            conn.close()

        synced = 0
        for m_id, m_type, user in records:
            dest = self.downloader._get_target_path(user or "unknown", m_id, m_type)
            if os.path.exists(dest) and os.path.getsize(dest) > 0:
                img_id = self.downloader.stash.register_media(dest, "image") if m_type == "image" else None
                scn_id = self.downloader.stash.register_media(dest, "video") if m_type != "image" else None
                if img_id or scn_id:
                    self.downloader._update_status(m_id, "COMPLETED", None, img_id, scn_id)
                    synced += 1
                else:
                    self.downloader._update_status(m_id, "RETAINED", "Saved on disk, awaiting Stash index", None, None)
        return synced

    def run_restore(self, progress_callback: Optional[Callable[[int, int, str], None]] = None) -> Dict[str, int]:
        if not os.path.exists(self.dumps_dir):
            if progress_callback: progress_callback(1, 1, f"Dumps directory not found: {self.dumps_dir}")
            return {"articles": 0, "media": 0, "avatars": 0}

        target_dirs = {root for root, _, files in os.walk(self.dumps_dir) if "metadata.json" in files or "snapshot.warc.gz" in files}
        total = len(target_dirs)
        if total == 0:
            if progress_callback: progress_callback(1, 1, "No valid dumps found in directory.")
            return {"articles": 0, "media": 0, "avatars": 0}

        if progress_callback: progress_callback(0, total, f"Found {total} dump entries. Starting offline restore...")

        os.makedirs(self.avatar_dir, exist_ok=True)
        stats = {"articles": 0, "media": 0, "avatars": 0}

        for idx, pdir in enumerate(sorted(target_dirs), start=1):
            post_id = os.path.basename(pdir)
            meta_path = os.path.join(pdir, "metadata.json")
            warc_path = os.path.join(pdir, "snapshot.warc.gz")
            restored = False

            if os.path.exists(meta_path):
                try:
                    with open(meta_path, "r", encoding="utf-8") as f:
                        if self.mutator.upsert_record(json.load(f)):
                            restored = True
                            stats["articles"] += 1
                except Exception as e:
                    print(f"[Restorer] Failed to parse {meta_path}: {e}")

            if not restored and os.path.exists(warc_path):
                count = WarcImporter(warc_path, db_path=self.db_path, storage_dir=self.downloader.storage_dir).run_import()
                if count > 0:
                    stats["articles"] += count
                    restored = True

            avatar_subdir = os.path.join(pdir, "avatars")
            if os.path.exists(avatar_subdir):
                for af in os.listdir(avatar_subdir):
                    src_av = os.path.join(avatar_subdir, af)
                    if os.path.isfile(src_av):
                        shutil.copy2(src_av, os.path.join(self.avatar_dir, af))
                        stats["avatars"] += 1

            if progress_callback: progress_callback(idx, total, f"Restoring [{idx}/{total}] Post: {post_id}")

        stats["media"] = self._sync_media_files()
        if progress_callback:
            msg = f"Restore completed! Articles: {stats['articles']}, Avatars: {stats['avatars']}, Media: {stats['media']}"
            progress_callback(total, total, msg)

        return stats
