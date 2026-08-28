# plugins/twitter/scraper/test_warc_importer.py (100行以下)
import io, json, os, shutil, sqlite3, subprocess, sys, tempfile, unittest
from warcio.warcwriter import WARCWriter
from warcio.statusandheaders import StatusAndHeaders
from core.warc_importer import WarcImporter


class TestWarcImporter(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.mkdtemp()
        self.db_path, self.blobs_dir, self.warc_path = os.path.join(self.temp_dir, "test_warc.db"), os.path.join(self.temp_dir, "blobs"), os.path.join(self.temp_dir, "sample.warc.gz")
        self._init_db(); self._create_dummy_warc()

    def tearDown(self): shutil.rmtree(self.temp_dir, ignore_errors=True)

    def _init_db(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.executescript("""
                CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT NOT NULL, display_name TEXT, avatar_url TEXT, avatar_base64 TEXT, description TEXT, group_name TEXT DEFAULT '', alias_of TEXT DEFAULT '', updated_at DATETIME);
                CREATE TABLE account_profile_histories (id INTEGER PRIMARY KEY AUTOINCREMENT, account_id TEXT NOT NULL, display_name TEXT NOT NULL, description TEXT DEFAULT '', avatar_original_url TEXT NOT NULL, avatar_seq INTEGER NOT NULL, avatar_virtual_key TEXT NOT NULL, avatar_base64 TEXT, observed_at DATETIME NOT NULL);
                CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT NOT NULL, conversation_id TEXT, reply_to_id TEXT, reply_to_handle TEXT, created_at DATETIME, full_text TEXT, lang TEXT DEFAULT 'ja', full_text_ja TEXT, full_text_en TEXT, full_text_zh TEXT, via TEXT DEFAULT 'twitter', is_repost INTEGER DEFAULT 0, is_liked INTEGER DEFAULT 0, wayback_url TEXT, source_name TEXT, source_domain TEXT, original_url TEXT, sotwe_url TEXT, nitter_url TEXT, twistalker_url TEXT, is_trash BOOLEAN DEFAULT 0, trashed_by TEXT, trash_reason TEXT, trashed_at DATETIME);
                CREATE TABLE media (media_id TEXT PRIMARY KEY, article_id TEXT NOT NULL, type TEXT, download_url TEXT, width INTEGER, height INTEGER, download_status TEXT, failed_reason TEXT, stash_scene_id TEXT, stash_image_id TEXT, media_quality TEXT);
                CREATE TABLE url_redirects (short_url TEXT PRIMARY KEY, expanded_url TEXT, article_id TEXT);
                CREATE TABLE whitelists (id INTEGER PRIMARY KEY AUTOINCREMENT, type TEXT NOT NULL, value TEXT NOT NULL UNIQUE, is_active BOOLEAN DEFAULT 1);
            """)

    def _create_dummy_warc(self):
        with open(self.warc_path, "wb") as f:
            w = WARCWriter(f, gzip=True)
            t_json = json.dumps({"tweet": {"id_str": "1888999000111", "full_text": "テスト https://t.co/famicom", "created_at": "2026-08-21",
                                           "extended_entities": {"media": [{"media_url_https": "https://pbs.twimg.com/media/chaldea_shield.jpg", "type": "image", "sizes": {"large": {"w": 100, "h": 100}}}],
                                                                 "urls": [{"url": "https://t.co/famicom", "expanded_url": "https://example.com/nes_apu"}]}},
                                 "user": {"id_str": "777001", "screen_name": "mash_kyrielight", "name": "マシュ"}}).encode("utf-8")
            h1 = StatusAndHeaders("200 OK", [("Content-Type", "application/json")])
            w.write_record(w.create_warc_record("https://api.twitter.com/1.1/statuses/show.json", "response", warc_headers_dict={"WARC-Target-URI": "https://twitter.com/mash_kyrielight/status/1888999000111"}, http_headers=h1, payload=io.BytesIO(t_json)))
            h2 = StatusAndHeaders("200 OK", [("Content-Type", "image/jpeg")])
            w.write_record(w.create_warc_record("https://pbs.twimg.com/media/chaldea_shield.jpg", "response", warc_headers_dict={"WARC-Target-URI": "https://pbs.twimg.com/media/chaldea_shield.jpg"}, http_headers=h2, payload=io.BytesIO(b"MOCK_PNG")))
            h3 = StatusAndHeaders("200 OK", [("Content-Type", "text/html")])
            w.write_record(w.create_warc_record("https://twitter.com/mash_kyrielight/status/1888999000222", "response", warc_headers_dict={"WARC-Target-URI": "https://twitter.com/mash_kyrielight/status/1888999000222"}, http_headers=h3, payload=io.BytesIO(b"<title>T</title><div data-tweet-id=\"1888999000222\"></div>")))

    def test_audit_warc(self):
        importer = WarcImporter(self.warc_path, db_path=self.db_path, storage_dir=self.blobs_dir, offline=True)
        audit = importer.audit_warc()
        self.assertEqual(audit["platform"], "twitter"); self.assertEqual(audit["account"], "mash_kyrielight")

    def test_run_import(self):
        importer = WarcImporter(self.warc_path, db_path=self.db_path, storage_dir=self.blobs_dir, offline=True)
        self.assertEqual(importer.run_import(), 2)
        with sqlite3.connect(self.db_path) as conn:
            self.assertIsNotNone(conn.cursor().execute("SELECT id FROM articles WHERE id = '1888999000111'").fetchone())
            self.assertEqual(conn.cursor().execute("SELECT download_status FROM media WHERE media_id = 'chaldea_shield.jpg'").fetchone()[0], "COMPLETED")
        self.assertTrue(os.path.exists(os.path.join(self.blobs_dir, "chaldea_shield.jpg")))

    def test_cli_execution(self):
        cmd = [sys.executable, os.path.join(os.path.dirname(__file__), "main.py"), "--mode", "manual", "--warc-path", self.warc_path, "--db-path", self.db_path, "--storage-dir", self.blobs_dir, "--offline"]
        proc = subprocess.run(cmd, capture_output=True, text=True, check=True)
        self.assertIn("PROGRESS:", proc.stdout)


if __name__ == "__main__":
    unittest.main()
