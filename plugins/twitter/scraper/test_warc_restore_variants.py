# plugins/twitter/scraper/test_warc_restore_variants.py (SPEC-PLUGIN-001 / 100行以下)
import gc, io, json, os, shutil, sqlite3, tempfile, time, unittest
from warcio.warcwriter import WARCWriter
from warcio.statusandheaders import StatusAndHeaders
from plugins.twitter.scraper.core.restorer import Restorer

class TestWarcRestoreVariants(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.mkdtemp()
        self.db = os.path.join(self.tmp, "archive.db")
        self.dumps = os.path.join(self.tmp, "dumps")
        self.post_dir = os.path.join(self.dumps, "twitter", "mash_kyrielight", "1888999000111")
        os.makedirs(self.post_dir, exist_ok=True)
        conn = sqlite3.connect(self.db)
        conn.execute("CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT, conversation_id TEXT, reply_to_id TEXT, reply_to_handle TEXT, created_at TEXT, full_text TEXT, lang TEXT, full_text_ja TEXT, full_text_en TEXT, full_text_zh TEXT, via TEXT, is_repost INTEGER, is_liked INTEGER, wayback_url TEXT, source_name TEXT, source_domain TEXT, original_url TEXT, sotwe_url TEXT, nitter_url TEXT, twistalker_url TEXT)")
        conn.execute("CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT, display_name TEXT, avatar_url TEXT, description TEXT, avatar_base64 TEXT, updated_at TEXT)")
        conn.execute("CREATE TABLE account_profile_histories (account_id TEXT, display_name TEXT, description TEXT, avatar_original_url TEXT, avatar_seq INTEGER, avatar_virtual_key TEXT, avatar_base64 TEXT, observed_at TEXT)")
        conn.execute("CREATE TABLE url_redirects (short_url TEXT PRIMARY KEY, expanded_url TEXT, article_id TEXT)")
        conn.execute("CREATE TABLE media (media_id TEXT PRIMARY KEY, article_id TEXT, account_id TEXT, type TEXT, download_url TEXT, width INTEGER, height INTEGER, thumbnail_url TEXT, tweet_urls TEXT, media_quality TEXT, download_status TEXT, failed_reason TEXT)")
        conn.commit(); conn.close()
        self._write_warc()

    def tearDown(self):
        gc.collect(); time.sleep(0.05); shutil.rmtree(self.tmp, ignore_errors=True)

    def _write_warc(self):
        t_json = json.dumps({"tweet": {"id_str": "1888999000111", "user": {"id_str": "999", "screen_name": "mash_kyrielight", "name": "Mash"},
                                       "extended_entities": {"media": [{"type": "photo", "media_url_https": "https://pbs.twimg.com/media/chaldea_shield.jpg"}]}}}).encode("utf-8")
        w_path = os.path.join(self.post_dir, "snapshot.warc.gz")
        with open(w_path, "wb") as f:
            w = WARCWriter(f, gzip=True)
            h = StatusAndHeaders("200 OK", [("Content-Type", "application/json; charset=utf-8")], protocol="HTTP/1.1")
            rec = w.create_warc_record("https://api.twitter.com/1.1/statuses/show.json", "response", warc_headers_dict={"WARC-Target-URI": "https://twitter.com/mash_kyrielight/status/1888999000111"}, http_headers=h, payload=io.BytesIO(t_json))
            w.write_record(rec)

    def test_restore_searches_warc_and_populates_variants(self):
        restorer = Restorer(dumps_dir=self.dumps, db_path=self.db, storage_dir=os.path.join(self.tmp, "blobs"), avatar_dir=os.path.join(self.tmp, "avatars"))
        stats = restorer.run_restore()
        self.assertEqual(stats["articles"], 1)
        with sqlite3.connect(self.db) as conn:
            cur = conn.cursor()
            rows = cur.execute("SELECT variant_hash, media_id, download_url, bit_rate, content_type FROM media_variants ORDER BY bit_rate DESC").fetchall()
            self.assertEqual(len(rows), 2)
            self.assertEqual(rows[0][3], 10000)
            self.assertIn("name=orig", rows[0][2])
            self.assertEqual(rows[0][4], "image/jpeg")
            self.assertEqual(rows[1][3], 5000)
            self.assertIn("name=large", rows[1][2])
            self.assertEqual(rows[1][4], "image/jpeg")

if __name__ == "__main__":
    unittest.main()
