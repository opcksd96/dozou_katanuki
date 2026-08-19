# plugins/twitter/scraper/test_downloader.py
import gc
import os
import sqlite3
import tempfile
import unittest
from unittest.mock import MagicMock, patch

from core.downloader import Downloader


class TestDownloaderPipeline(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.TemporaryDirectory()
        self.db_path = os.path.join(self.temp_dir.name, "test_archive.db")
        self.storage_dir = os.path.join(self.temp_dir.name, "storage")

        # Setup test schema
        conn = sqlite3.connect(self.db_path)
        try:
            conn.execute("CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT);")
            conn.execute("CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT);")
            conn.execute("""
                CREATE TABLE media (
                    media_id TEXT PRIMARY KEY,
                    article_id TEXT,
                    type TEXT,
                    download_url TEXT,
                    download_status TEXT DEFAULT 'QUEUED',
                    failed_reason TEXT,
                    stash_scene_id TEXT,
                    stash_image_id TEXT
                );
            """)
            conn.execute("INSERT INTO accounts VALUES ('1001', 'alice');")
            conn.execute("INSERT INTO articles VALUES ('post_1', '1001');")
            conn.commit()
        finally:
            conn.close()

        self.dl = Downloader(db_path=self.db_path, storage_dir=self.storage_dir)

    def tearDown(self):
        self.dl.session.close()
        del self.dl
        gc.collect()
        try:
            self.temp_dir.cleanup()
        except Exception:
            pass

    def test_stage1_success_direct_download(self):
        """第1段階: requests 直接取得成功 -> COMPLETED"""
        conn = sqlite3.connect(self.db_path)
        try:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) "
                         "VALUES ('img1.jpg', 'post_1', 'image', 'https://example.com/img1.jpg', 'QUEUED');")
            conn.commit()
        finally:
            conn.close()

        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.iter_content.return_value = [b"fake_image_bytes"]

        with patch.object(self.dl.session, "get", return_value=mock_resp):
            with patch.object(self.dl.stash, "find_image_by_path", return_value="stash-img-123"):
                success_count = self.dl.process_queued_media(article_id="post_1")
                self.assertEqual(success_count, 1)

        conn = sqlite3.connect(self.db_path)
        try:
            row = conn.cursor().execute("SELECT download_status, stash_image_id FROM media WHERE media_id = 'img1.jpg'").fetchone()
            self.assertEqual(row[0], "COMPLETED")
            self.assertEqual(row[1], "stash-img-123")
        finally:
            conn.close()

    def test_stage2_escalation_to_aria2_on_404(self):
        """第2段階: 404原本消失 -> Motrix/Aria2委託 -> OUTSOURCED"""
        conn = sqlite3.connect(self.db_path)
        try:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) "
                         "VALUES ('vid1.mp4', 'post_1', 'video', 'https://example.com/vid1.mp4', 'QUEUED');")
            conn.commit()
        finally:
            conn.close()

        mock_resp = MagicMock()
        mock_resp.status_code = 404

        with patch.object(self.dl.session, "get", return_value=mock_resp):
            with patch.object(self.dl.aria2, "add_uri", return_value="gid-98765"):
                self.dl.process_queued_media(article_id="post_1")

        conn = sqlite3.connect(self.db_path)
        try:
            row = conn.cursor().execute("SELECT download_status, failed_reason FROM media WHERE media_id = 'vid1.mp4'").fetchone()
            self.assertEqual(row[0], "OUTSOURCED")
            self.assertIn("gid-98765", row[1])
        finally:
            conn.close()

    def test_stage3_polling_recovery(self):
        """第3段階: OUTSOURCED 実ファイル検知 -> COMPLETED"""
        conn = sqlite3.connect(self.db_path)
        try:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) "
                         "VALUES ('vid2.mp4', 'post_1', 'video', 'https://example.com/vid2.mp4', 'OUTSOURCED');")
            conn.commit()
        finally:
            conn.close()

        target_path = self.dl._get_target_path("alice", "vid2.mp4")
        with open(target_path, "wb") as f:
            f.write(b"salvaged_video_data")

        with patch.object(self.dl.stash, "find_scene_by_path", return_value="scene-777"):
            salvaged = self.dl.poll_outsourced_media()
            self.assertEqual(salvaged, 1)

        conn = sqlite3.connect(self.db_path)
        try:
            row = conn.cursor().execute("SELECT download_status, stash_scene_id FROM media WHERE media_id = 'vid2.mp4'").fetchone()
            self.assertEqual(row[0], "COMPLETED")
            self.assertEqual(row[1], "scene-777")
        finally:
            conn.close()


if __name__ == "__main__":
    unittest.main()
