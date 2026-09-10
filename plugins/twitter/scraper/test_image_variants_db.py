# plugins/twitter/scraper/test_image_variants_db.py (SPEC-PLUGIN-001 / 100行以下)
import gc, os, sqlite3, time, unittest
from plugins.base.scraper.core.base_mutator import BaseMutator
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from plugins.twitter.scraper.parsers.x_media import extract_media_from_syndication
from plugins.twitter.scraper.parsers.twistalker_media import extract_media_from_api
from plugins.twitter.scraper.parsers.sotwe_extractors import extract_media_entities

class TestImageVariantsDb(unittest.TestCase):
    def setUp(self):
        self.db_path = f"test_img_var_{self._testMethodName}.db"
        if os.path.exists(self.db_path): os.remove(self.db_path)
        conn = sqlite3.connect(self.db_path)
        conn.execute("CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT, conversation_id TEXT, reply_to_id TEXT, reply_to_handle TEXT, created_at TEXT, full_text TEXT, lang TEXT, full_text_ja TEXT, full_text_en TEXT, full_text_zh TEXT, via TEXT, is_repost INTEGER, is_liked INTEGER, wayback_url TEXT, source_name TEXT, source_domain TEXT, original_url TEXT, sotwe_url TEXT, nitter_url TEXT, twistalker_url TEXT)")
        conn.execute("CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT, display_name TEXT, avatar_url TEXT, description TEXT, avatar_base64 TEXT, updated_at TEXT)")
        conn.execute("CREATE TABLE account_profile_histories (account_id TEXT, display_name TEXT, description TEXT, avatar_original_url TEXT, avatar_seq INTEGER, avatar_virtual_key TEXT, avatar_base64 TEXT, observed_at TEXT)")
        conn.execute("CREATE TABLE url_redirects (short_url TEXT PRIMARY KEY, expanded_url TEXT, article_id TEXT)")
        conn.execute("CREATE TABLE media (media_id TEXT PRIMARY KEY, article_id TEXT, account_id TEXT, type TEXT, download_url TEXT, width INTEGER, height INTEGER, thumbnail_url TEXT, tweet_urls TEXT, media_quality TEXT, download_status TEXT, failed_reason TEXT)")
        conn.commit(); conn.close()
        self.mutator = BaseMutator(db_path=self.db_path, platform="twitter", enable_translation=False)

    def tearDown(self):
        del self.mutator; gc.collect(); time.sleep(0.05)
        for p in [self.db_path, f"{self.db_path}-wal", f"{self.db_path}-shm"]:
            if os.path.exists(p):
                try: os.remove(p)
                except Exception: pass

    def test_twitter_parser_and_db_persistence(self):
        record = {
            "tweet": {"id_str": "12345", "user": {"id_str": "999", "screen_name": "mash_kyrielight", "name": "Mash"},
                      "extended_entities": {"media": [{"type": "photo", "media_url_https": "https://pbs.twimg.com/media/GnmCCzebYAAkg9-.jpg"}]}}
        }
        parsed = TwitterParser().parse_record(record, "https://x.com/mash_kyrielight/status/12345")
        self.mutator.upsert_batch([parsed])
        with sqlite3.connect(self.db_path) as conn:
            rows = conn.execute("SELECT variant_hash, media_id, download_url, bit_rate, content_type FROM media_variants ORDER BY bit_rate DESC").fetchall()
            self.assertEqual(len(rows), 2)
            self.assertEqual(rows[0][3], 10000)
            self.assertIn("name=orig", rows[0][2])
            self.assertEqual(rows[0][4], "image/jpeg")
            self.assertEqual(rows[1][3], 5000)
            self.assertIn("name=large", rows[1][2])
            self.assertEqual(rows[1][4], "image/jpeg")

    def test_sotwe_image_variants(self):
        sotwe_raw = {"mediaEntities": [{"type": "image", "mediaURL": "https://pbs.twimg.com/media/test_sotwe.jpg"}]}
        media = extract_media_entities(sotwe_raw)
        self.assertEqual(len(media[0]["variants"]), 2)
        self.assertEqual(media[0]["variants"][0]["bit_rate"], 10000)
        self.assertEqual(media[0]["variants"][1]["bit_rate"], 5000)

    def test_x_media_syndication_variants(self):
        details = [{"type": "photo", "media_url_https": "https://pbs.twimg.com/media/test_x.jpg"}]
        media = extract_media_from_syndication(details)
        self.assertEqual(len(media[0]["variants"]), 2)
        self.assertEqual(media[0]["variants"][0]["bit_rate"], 10000)
        self.assertEqual(media[0]["variants"][1]["bit_rate"], 5000)

    def test_twistalker_image_variants(self):
        entities = {"media": [{"type": "image", "media_url_https": "https://pbs.twimg.com/media/test_tws.jpg"}]}
        media = extract_media_from_api(entities)
        self.assertEqual(len(media[0]["variants"]), 2)
        self.assertEqual(media[0]["variants"][0]["bit_rate"], 10000)
        self.assertEqual(media[0]["variants"][1]["bit_rate"], 5000)

if __name__ == "__main__":
    unittest.main()
