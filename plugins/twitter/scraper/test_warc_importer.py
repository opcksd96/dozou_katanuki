# plugins/twitter/scraper/test_warc_importer.py
import io, json, os, shutil, sqlite3, subprocess, sys, tempfile, unittest

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
if CURRENT_DIR not in sys.path:
    sys.path.insert(0, CURRENT_DIR)

from warcio.warcwriter import WARCWriter
from warcio.statusandheaders import StatusAndHeaders
from core.warc_importer import WarcImporter


class TestWarcImporter(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.mkdtemp()
        self.db_path = os.path.join(self.temp_dir, "test_warc.db")
        self.blobs_dir = os.path.join(self.temp_dir, "blobs")
        self.warc_path = os.path.join(self.temp_dir, "sample.warc.gz")
        self._init_db()
        self._create_dummy_warc()

    def tearDown(self):
        shutil.rmtree(self.temp_dir, ignore_errors=True)

    def _init_db(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.executescript("""
                CREATE TABLE accounts (
                    numeric_id TEXT PRIMARY KEY, username TEXT NOT NULL, display_name TEXT, avatar_url TEXT, updated_at DATETIME
                );
                CREATE TABLE articles (
                    id TEXT PRIMARY KEY, account_id TEXT NOT NULL, conversation_id TEXT, reply_to_id TEXT, reply_to_handle TEXT,
                    created_at DATETIME, full_text TEXT, lang TEXT DEFAULT 'ja', full_text_ja TEXT, full_text_en TEXT,
                    full_text_zh TEXT, via TEXT DEFAULT 'twitter', is_repost INTEGER DEFAULT 0, is_liked INTEGER DEFAULT 0,
                    wayback_url TEXT, FOREIGN KEY (account_id) REFERENCES accounts(numeric_id)
                );
                CREATE TABLE media (
                    media_id TEXT PRIMARY KEY, article_id TEXT NOT NULL, type TEXT, download_url TEXT, width INTEGER,
                    height INTEGER, download_status TEXT, failed_reason TEXT, stash_scene_id TEXT, stash_image_id TEXT,
                    FOREIGN KEY (article_id) REFERENCES articles(id)
                );
                CREATE TABLE url_redirects (
                    short_url TEXT PRIMARY KEY, expanded_url TEXT, article_id TEXT
                );
            """)

    def _create_dummy_warc(self):
        with open(self.warc_path, "wb") as f:
            writer = WARCWriter(f, gzip=True)
            # Record 1: API JSON Tweet
            tweet_payload = json.dumps({
                "tweet": {
                    "id_str": "1888999000111",
                    "full_text": "先輩！WARCからの完全オフライン抽出テスト投稿です！ https://t.co/famicom",
                    "created_at": "2026-08-21 02:00:00",
                    "extended_entities": {
                        "media": [{"media_url_https": "https://pbs.twimg.com/media/chaldea_shield.jpg", "type": "image", "sizes": {"large": {"w": 1920, "h": 1080}}}],
                        "urls": [{"url": "https://t.co/famicom", "expanded_url": "https://example.com/nes_apu"}]
                    }
                },
                "user": {"id_str": "777001", "screen_name": "mash_kyrielight", "name": "マシュ・キリエライト"}
            }).encode("utf-8")
            h1 = StatusAndHeaders("200 OK", [("Content-Type", "application/json; charset=utf-8")])
            r1 = writer.create_warc_record("https://api.twitter.com/1.1/statuses/show.json?id=1888999000111", "response",
                                          warc_headers_dict={"WARC-Target-URI": "https://twitter.com/mash_kyrielight/status/1888999000111"},
                                          http_headers=h1, payload=io.BytesIO(tweet_payload))
            writer.write_record(r1)

            # Record 2: Media Binary
            media_bytes = b"MOCK_PNG_IMAGE_BINARY_PAYLOAD_FOR_OFFLINE_STASH"
            h2 = StatusAndHeaders("200 OK", [("Content-Type", "image/jpeg")])
            r2 = writer.create_warc_record("https://pbs.twimg.com/media/chaldea_shield.jpg", "response",
                                          warc_headers_dict={"WARC-Target-URI": "https://pbs.twimg.com/media/chaldea_shield.jpg"},
                                          http_headers=h2, payload=io.BytesIO(media_bytes))
            writer.write_record(r2)

            # Record 3: HTML Snapshot
            html_payload = """<!DOCTYPE html><html><head><title>Twitter Post</title>
            <meta property="og:description" content="マシュ・キリエライト: 先輩、2件目のHTMLアーカイブ記事です。">
            </head><body><div data-tweet-id="1888999000222"></div></body></html>""".encode("utf-8")
            h3 = StatusAndHeaders("200 OK", [("Content-Type", "text/html; charset=utf-8")])
            r3 = writer.create_warc_record("https://twitter.com/mash_kyrielight/status/1888999000222", "response",
                                          warc_headers_dict={"WARC-Target-URI": "https://twitter.com/mash_kyrielight/status/1888999000222"},
                                          http_headers=h3, payload=io.BytesIO(html_payload))
            writer.write_record(r3)

    def test_audit_warc(self):
        importer = WarcImporter(self.warc_path, db_path=self.db_path, storage_dir=self.blobs_dir, offline=True)
        audit = importer.audit_warc()
        self.assertEqual(audit["platform"], "twitter")
        self.assertEqual(audit["account"], "mash_kyrielight")
        self.assertGreaterEqual(audit["records"], 3)

    def test_run_import(self):
        logs = []
        importer = WarcImporter(self.warc_path, db_path=self.db_path, storage_dir=self.blobs_dir, offline=True)
        saved = importer.run_import(progress_callback=lambda c, t, m: logs.append(f"{c}/{t}: {m}"))

        self.assertEqual(saved, 2)
        self.assertTrue(any("Audited WARC" in l for l in logs))

        with sqlite3.connect(self.db_path) as conn:
            cur = conn.cursor()
            acc = cur.execute("SELECT numeric_id, username, display_name FROM accounts WHERE username = 'mash_kyrielight'").fetchone()
            self.assertIsNotNone(acc)
            self.assertEqual(acc[1], "mash_kyrielight")

            art1 = cur.execute("SELECT id, full_text FROM articles WHERE id = '1888999000111'").fetchone()
            self.assertIsNotNone(art1)
            self.assertIn("https://example.com/nes_apu", art1[1])

            art2 = cur.execute("SELECT id FROM articles WHERE id = '1888999000222'").fetchone()
            self.assertIsNotNone(art2)

            med = cur.execute("SELECT media_id, download_status FROM media WHERE media_id = 'chaldea_shield.jpg'").fetchone()
            self.assertIsNotNone(med)
            self.assertEqual(med[1], "COMPLETED")

        # Verify extracted binary file exists in blobs_dir
        extracted_file = os.path.join(self.blobs_dir, "chaldea_shield.jpg")
        self.assertTrue(os.path.exists(extracted_file))
        with open(extracted_file, "rb") as f:
            self.assertEqual(f.read(), b"MOCK_PNG_IMAGE_BINARY_PAYLOAD_FOR_OFFLINE_STASH")

    def test_cli_execution(self):
        main_py = os.path.join(CURRENT_DIR, "main.py")
        cmd = [
            sys.executable, main_py,
            "--mode", "manual",
            "--warc-path", self.warc_path,
            "--db-path", self.db_path,
            "--storage-dir", self.blobs_dir,
            "--offline"
        ]
        proc = subprocess.run(cmd, capture_output=True, text=True, check=True)
        self.assertIn("PROGRESS:", proc.stdout)
        self.assertIn("Offline import completed", proc.stdout)


if __name__ == "__main__":
    unittest.main()
