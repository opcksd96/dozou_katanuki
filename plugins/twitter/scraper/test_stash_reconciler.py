# plugins/twitter/scraper/test_stash_reconciler.py (SPEC-PLUGIN-001 / 100行以下)
import os, shutil, sqlite3, tempfile, unittest
from unittest.mock import MagicMock, patch
from plugins.base.scraper.core.stash_reconciler import StashReconciler

class TestStashReconciler(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.mkdtemp()
        self.db_path = os.path.join(self.temp_dir, "test.db")
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT, display_name TEXT);")
            conn.execute("CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT, wayback_url TEXT, full_text TEXT, full_text_ja TEXT, created_at TEXT);")
            conn.execute("CREATE TABLE media (media_id TEXT PRIMARY KEY, article_id TEXT, type TEXT, download_url TEXT, download_status TEXT DEFAULT 'QUEUED', stash_scene_id TEXT, stash_image_id TEXT, thumbnail_url TEXT);")
            conn.execute("CREATE UNIQUE INDEX idx_media_stash_scene ON media(stash_scene_id) WHERE stash_scene_id IS NOT NULL;")
            conn.execute("INSERT INTO accounts VALUES ('1001', 'alice', 'アリス'), ('1002', 'bob', 'ボブ');")
            conn.execute("INSERT INTO articles VALUES ('post_1', '1001', '', 'Hello', '', '2025-01-22'), ('post_2', '1002', '', 'Hi', '', '2025-01-23');")
        self.reconciler = StashReconciler(stash=MagicMock())

    def tearDown(self):
        shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_reconcile_skips_already_bound_scene_without_error(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, stash_scene_id, download_status) VALUES ('vid1.mp4', 'post_1', 'video', 'scene-101', 'COMPLETED');")
            conn.execute("INSERT INTO media (media_id, article_id, type, download_status) VALUES ('vid2.mp4', 'post_1', 'video', 'QUEUED');")
        mock_data = {
            "allScenes": [{"id": "scene-101", "title": "X (@alice): Tweet post_1", "details": "", "files": []}],
            "allImages": []
        }
        with patch.object(self.reconciler.stash, "query", return_value=mock_data):
            bound = self.reconciler.reconcile_to_db(self.db_path)
            self.assertEqual(bound, 0)
        with sqlite3.connect(self.db_path) as conn:
            row1 = conn.cursor().execute("SELECT stash_scene_id FROM media WHERE media_id = 'vid1.mp4'").fetchone()
            row2 = conn.cursor().execute("SELECT stash_scene_id FROM media WHERE media_id = 'vid2.mp4'").fetchone()
            self.assertEqual(row1[0], "scene-101")
            self.assertIsNone(row2[0])

    def test_reconcile_binds_new_scene_by_title_pattern(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("INSERT INTO media (media_id, article_id, type, download_status) VALUES ('vid3.mp4', 'post_2', 'video', 'QUEUED');")
        mock_data = {
            "allScenes": [{"id": "scene-202", "title": "X (@bob): Tweet post_2", "details": "", "files": []}],
            "allImages": []
        }
        with patch.object(self.reconciler.stash, "query", return_value=mock_data):
            bound = self.reconciler.reconcile_to_db(self.db_path)
            self.assertEqual(bound, 1)
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT stash_scene_id, download_status FROM media WHERE media_id = 'vid3.mp4'").fetchone()
            self.assertEqual(row, ("scene-202", "COMPLETED"))

if __name__ == "__main__":
    unittest.main()
