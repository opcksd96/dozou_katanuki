# plugins/twitter/scraper/test_downloader.py (100行以下)
import gc, os, sqlite3, tempfile, unittest
from unittest.mock import MagicMock, patch
from core.downloader import Downloader


class TestDownloaderPipeline(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.TemporaryDirectory()
        self.db_path, self.storage_dir = os.path.join(self.temp_dir.name, "test.db"), os.path.join(self.temp_dir.name, "storage")
        conn = sqlite3.connect(self.db_path)
        try:
            conn.execute("CREATE TABLE whitelists (id INTEGER PRIMARY KEY, type TEXT, value TEXT, is_active INTEGER);")
            conn.execute("CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT, display_name TEXT);")
            conn.execute("CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT, wayback_url TEXT, full_text TEXT, full_text_ja TEXT, created_at TEXT);")
            conn.execute("CREATE TABLE media (media_id TEXT PRIMARY KEY, article_id TEXT, type TEXT, download_url TEXT, download_status TEXT DEFAULT 'QUEUED', failed_reason TEXT, stash_scene_id TEXT, stash_image_id TEXT);")
            conn.execute("INSERT INTO whitelists VALUES (1, 'account', 'alice', 1);")
            conn.execute("INSERT INTO accounts VALUES ('1001', 'alice', 'アリス'), ('1002', 'bob', 'ボブ');")
            conn.execute("INSERT INTO articles VALUES ('post_1', '1001', '', 'Hello world', 'こんにちは', '2025-01-22T10:00:00Z'), ('post_2', '1002', '', 'Secret', '', '2025-01-23T10:00:00Z');")
            conn.commit()
        finally:
            conn.close()
        self.dl = Downloader(db_path=self.db_path, storage_dir=self.storage_dir)

    def tearDown(self):
        self.dl.session.close(); del self.dl; gc.collect()
        try: self.temp_dir.cleanup()
        except Exception: pass

    def test_stage1_success_direct_download(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) VALUES ('img1.jpg', 'post_1', 'image', 'https://example.com/img1.jpg', 'QUEUED');")
        m_resp = MagicMock(status_code=200, iter_content=MagicMock(return_value=[b"fake_bytes"]))
        with patch.object(self.dl.session, "get", return_value=m_resp), patch.object(self.dl.stash, "register_media", return_value="stash-img-123"):
            self.assertEqual(self.dl.process_queued_media(article_id="post_1"), 1)
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT download_status, stash_image_id FROM media WHERE media_id = 'img1.jpg'").fetchone()
            self.assertEqual(row, ("COMPLETED", "stash-img-123"))

    def test_stage2_escalation_to_aria2_on_403_and_404(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) VALUES ('vid1.mp4', 'post_1', 'video', 'https://example.com/vid1.mp4', 'QUEUED');")
        m_resp = MagicMock(status_code=403)
        with patch.object(self.dl.session, "get", return_value=m_resp), patch.object(self.dl.aria2, "add_uri", return_value="gid-98765"):
            self.dl.process_queued_media(article_id="post_1")
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT download_status, failed_reason FROM media WHERE media_id = 'vid1.mp4'").fetchone()
            self.assertEqual(row[0], "OUTSOURCED")

    def test_non_whitelist_liveness_only_marked_dead404(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) VALUES ('ext_vid.mp4', 'post_2', 'video', 'https://video.twimg.com/ext.mp4', 'QUEUED');")
        m_head = MagicMock(status_code=403)
        with patch.object(self.dl.session, "head", return_value=m_head), patch.object(self.dl.aria2, "add_uri") as mock_aria:
            self.assertEqual(self.dl.process_queued_media(article_id="post_2"), 0)
            mock_aria.assert_not_called()
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT download_status, failed_reason FROM media WHERE media_id = 'ext_vid.mp4'").fetchone()
            self.assertEqual(row[0], "EXCLUDED")
            self.assertIn("Whitelist外", row[1])

    def test_stage3_polling_recovery(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_url, download_status) VALUES ('vid2.mp4', 'post_1', 'video', 'https://example.com/vid2.mp4', 'OUTSOURCED');")
        target = self.dl._get_target_path("alice", "vid2.mp4")
        with open(target, "wb") as f: f.write(b"salvaged")
        with patch.object(self.dl.stash, "register_media", return_value="scene-777"):
            self.assertEqual(self.dl.poll_outsourced_media(), 1)
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT download_status, stash_scene_id FROM media WHERE media_id = 'vid2.mp4'").fetchone()
            self.assertEqual(row, ("COMPLETED", "scene-777"))

    def test_stash_reconciliation_by_standard_title(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_status) VALUES ('unbound.jpg', 'post_1', 'image', 'QUEUED');")
        mock_data = {
            "allScenes": [],
            "allImages": [{"id": "img-999", "title": "X (@alice): Tweet post_1"}]
        }
        with patch.object(self.dl.stash, "query", return_value=mock_data):
            bound = self.dl.stash.reconcile_to_db(self.db_path)
            self.assertEqual(bound, 1)
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT stash_image_id, download_status FROM media WHERE media_id = 'unbound.jpg'").fetchone()
            self.assertEqual(row, ("img-999", "COMPLETED"))


if __name__ == "__main__":
    unittest.main()
