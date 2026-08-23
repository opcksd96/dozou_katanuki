# plugins/twitter/scraper/test_decorator.py (100行以下)
import os, sqlite3, tempfile, unittest, gc, time
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from plugins.twitter.scraper.core.mutator import Mutator


class TestDataDecorator(unittest.TestCase):
    def setUp(self):
        self.parser = TwitterParser()
        self.temp_db = tempfile.NamedTemporaryFile(delete=False, suffix=".db")
        self.temp_db.close(); self.db_path = self.temp_db.name

        conn = sqlite3.connect(self.db_path)
        with open("archive_schema.sql", "r", encoding="utf-8") as f:
            conn.executescript(";".join([s for s in f.read().split(";") if "sqlite_sequence" not in s]))
        conn.close()
        self.mutator = Mutator(db_path=self.db_path)

    def tearDown(self):
        gc.collect(); time.sleep(0.05)
        if os.path.exists(self.db_path):
            try: os.remove(self.db_path)
            except OSError: pass

    def test_twitter_parser_extracts_urls(self):
        raw = {
            "tweet": {"id_str": "123456789", "full_text": "Check out https://t.co/abc123xyz",
                      "entities": {"urls": [{"url": "https://t.co/abc123xyz", "expanded_url": "https://example.com/nes-apu"}]}},
            "user": {"id_str": "99999", "screen_name": "mash_fgo", "name": "Mash"}
        }
        parsed = self.parser.parse_record(raw, "https://twitter.com/mash_fgo/status/123456789")
        self.assertEqual(len(parsed["post"]["urls"]), 1)
        self.assertEqual(parsed["post"]["urls"][0]["expanded_url"], "https://example.com/nes-apu")

    def test_mutator_saves_url_redirects_and_expands_text(self):
        record = {
            "platform": "twitter",
            "account": {"numeric_id": "99999", "username": "mash_fgo", "display_name": "Mash", "avatar_url": "avatar.jpg"},
            "post": {
                "id": "123456789", "conversation_id": "123456789", "full_text": "Check https://t.co/abc123xyz",
                "wayback_url": "https://web.archive.org/web/...",
                "urls": [{"short_url": "https://t.co/abc123xyz", "expanded_url": "https://example.com/nes-apu"}],
            }, "media": [],
        }
        self.assertTrue(self.mutator.upsert_record(record))
        with sqlite3.connect(self.db_path) as conn:
            cur = conn.cursor()
            row = cur.execute("SELECT short_url, expanded_url FROM url_redirects WHERE article_id = '123456789'").fetchone()
            self.assertEqual(row, ("https://t.co/abc123xyz", "https://example.com/nes-apu"))
            art = cur.execute("SELECT full_text FROM articles WHERE id = '123456789'").fetchone()
            self.assertIn("https://example.com/nes-apu", art[0])


    def test_multiple_media_extraction_from_html(self):
        html = """
        <div data-tweet-id="123456789">
            <strong class="fullname">Mash</strong>
            <img src="https://pbs.twimg.com/profile_images/1/avatar.jpg" class="avatar" />
            <div class="AdaptiveMedia">
                <img src="https://pbs.twimg.com/media/photo_1.jpg:large" />
                <img src="https://pbs.twimg.com/media/photo_2.jpg:orig" />
                <img src="https://pbs.twimg.com/media/photo_3?format=png&name=small" />
                <video src="https://video.twimg.com/ext_tw_video/vid_1.mp4"></video>
            </div>
        </div>
        """
        parsed = self.parser.parse_record(html, "https://twitter.com/mash_fgo/status/123456789")
        self.assertIsNotNone(parsed)
        media = parsed.get("media", [])
        self.assertEqual(len(media), 4)
        urls = [m["url"] for m in media]
        self.assertTrue(any("photo_1" in u for u in urls))
        self.assertTrue(any("photo_2" in u for u in urls))
        self.assertTrue(any("photo_3" in u for u in urls))
        self.assertTrue(any("vid_1" in u for u in urls))


if __name__ == "__main__":
    unittest.main()
