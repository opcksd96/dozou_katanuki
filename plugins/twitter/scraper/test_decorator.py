import os
import sqlite3
import tempfile
import unittest
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from plugins.twitter.scraper.core.mutator import Mutator


class TestDataDecorator(unittest.TestCase):
    def setUp(self):
        self.parser = TwitterParser()
        self.temp_db = tempfile.NamedTemporaryFile(delete=False, suffix=".db")
        self.temp_db.close()
        self.db_path = self.temp_db.name

        # DDL 初期化
        conn = sqlite3.connect(self.db_path)
        with open("archive_schema.sql", "r", encoding="utf-8") as f:
            ddl = f.read()
            # sqlite_sequence は自動生成されるためスキップ
            ddl_statements = [s for s in ddl.split(";") if "sqlite_sequence" not in s]
            conn.executescript(";".join(ddl_statements))
        conn.close()

        self.mutator = Mutator(db_path=self.db_path)

    def tearDown(self):
        import gc
        import time
        gc.collect()
        time.sleep(0.05)
        if os.path.exists(self.db_path):
            try:
                os.remove(self.db_path)
            except OSError:
                pass

    def test_twitter_parser_extracts_urls(self):
        raw_json = {
            "tweet": {
                "id_str": "123456789",
                "full_text": "Check out our site! https://t.co/abc123xyz #NES",
                "entities": {
                    "urls": [
                        {
                            "url": "https://t.co/abc123xyz",
                            "expanded_url": "https://example.com/nes-apu-guide",
                        }
                    ]
                },
            },
            "user": {
                "id_str": "99999",
                "screen_name": "mash_fgo",
                "name": "Mash Kyrielight",
            },
        }

        parsed = self.parser.parse_record(raw_json, "https://twitter.com/mash_fgo/status/123456789")
        self.assertIsNotNone(parsed)
        self.assertIn("urls", parsed["post"])
        self.assertEqual(len(parsed["post"]["urls"]), 1)
        self.assertEqual(parsed["post"]["urls"][0]["short_url"], "https://t.co/abc123xyz")
        self.assertEqual(parsed["post"]["urls"][0]["expanded_url"], "https://example.com/nes-apu-guide")

    def test_mutator_saves_url_redirects_and_expands_text(self):
        record = {
            "platform": "twitter",
            "account": {
                "numeric_id": "99999",
                "username": "mash_fgo",
                "display_name": "Mash Kyrielight",
                "avatar_url": "avatar.jpg",
            },
            "post": {
                "id": "123456789",
                "conversation_id": "123456789",
                "full_text": "Check out our site! https://t.co/abc123xyz #NES",
                "wayback_url": "https://web.archive.org/web/...",
                "urls": [
                    {
                        "short_url": "https://t.co/abc123xyz",
                        "expanded_url": "https://example.com/nes-apu-guide",
                    }
                ],
            },
            "media": [],
        }

        success = self.mutator.upsert_record(record)
        self.assertTrue(success)

        # DB の検証
        conn = sqlite3.connect(self.db_path)
        cur = conn.cursor()

        # 1. url_redirects に保存されているか
        cur.execute("SELECT short_url, expanded_url, article_id FROM url_redirects WHERE article_id = '123456789'")
        row = cur.fetchone()
        self.assertIsNotNone(row)
        self.assertEqual(row[0], "https://t.co/abc123xyz")
        self.assertEqual(row[1], "https://example.com/nes-apu-guide")
        self.assertEqual(row[2], "123456789")

        # 2. articles.full_text 内の短縮URLが展開されているか
        cur.execute("SELECT full_text FROM articles WHERE id = '123456789'")
        art_row = cur.fetchone()
        self.assertIsNotNone(art_row)
        self.assertIn("https://example.com/nes-apu-guide", art_row[0])
        self.assertNotIn("https://t.co/abc123xyz", art_row[0])

        conn.close()


if __name__ == "__main__":
    unittest.main()
