# plugins/twitter/scraper/test_translator.py (100行以下)
import os, sqlite3, tempfile, unittest, gc, time
from unittest.mock import MagicMock, patch
from plugins.twitter.scraper.core.translator import Translator
from plugins.twitter.scraper.core.mutator import Mutator
from plugins.twitter.scraper.main import run_batch_translate


class TestTranslatorPipeline(unittest.TestCase):
    def setUp(self):
        self.temp_db = tempfile.NamedTemporaryFile(delete=False, suffix=".db")
        self.temp_db.close(); self.db_path = self.temp_db.name
        conn = sqlite3.connect(self.db_path)
        with open("archive_schema.sql", "r", encoding="utf-8") as f:
            conn.executescript(";".join([s for s in f.read().split(";") if "sqlite_sequence" not in s]))
        conn.close()

    def tearDown(self):
        gc.collect(); time.sleep(0.05)
        if os.path.exists(self.db_path):
            try: os.remove(self.db_path)
            except OSError: pass

    def test_detect_lang(self):
        t = Translator()
        self.assertEqual(t.detect_lang("先輩、おはようございます！"), "ja")
        self.assertEqual(t.detect_lang("NES APU register emulation test"), "en")
        self.assertEqual(t.detect_lang("红白机音频测试"), "zh")

    def test_translate_fallback_without_keys(self):
        t = Translator(provider="none")
        res = t.translate_article("先輩、おはようございます！")
        self.assertEqual(res["lang"], "ja")
        self.assertEqual(res["ja"], "先輩、おはようございます！")
        self.assertIsNone(res["en"])
        self.assertIsNone(res["zh"])

    @patch("requests.Session.post")
    def test_deepl_mock_translation(self, mock_post):
        mock_resp = MagicMock(); mock_resp.status_code = 200
        mock_resp.json.return_value = {"translations": [{"text": "Good morning, Senpai!"}]}
        mock_post.return_value = mock_resp
        with patch.dict(os.environ, {"DEEPL_API_KEY": "test-key:fx"}):
            t = Translator(delay_sec=0)
            self.assertEqual(t.translate("先輩、おはようございます！", "en", "ja"), "Good morning, Senpai!")

    def test_mutator_saves_translated_article(self):
        mock_trans = Translator()
        mock_trans.translate_article = MagicMock(return_value={"lang": "en", "ja": "ファミコン", "en": "NES", "zh": "红白机"})
        mut = Mutator(db_path=self.db_path, translator=mock_trans)
        record = {
            "account": {"numeric_id": "1001", "username": "senpai", "display_name": "Senpai"},
            "post": {"id": "art_2001", "full_text": "NES", "created_at": "2026-08-21 12:00:00"}, "media": []
        }
        self.assertTrue(mut.upsert_record(record))
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT lang, full_text_ja, full_text_zh FROM articles WHERE id = 'art_2001'").fetchone()
            self.assertEqual(row, ("en", "ファミコン", "红白机"))

    def test_run_batch_translate(self):
        with sqlite3.connect(self.db_path) as conn:
            conn.cursor().execute("""
                INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at) VALUES ('1002', 'mash', 'Mash', 'avatar.jpg', '2026-08-21 12:00:00');
            """)
            conn.cursor().execute("""
                INSERT INTO articles (id, account_id, conversation_id, created_at, full_text, lang, full_text_ja, via, wayback_url)
                VALUES ('art_3001', '1002', 'art_3001', '2026-08-21 12:00:00', 'Chaldea base test', 'en', NULL, 'twitter', 'https://web.archive.org/...');
            """)
            conn.commit()
        run_batch_translate(self.db_path, article_id="art_3001", overwrite=True)
        with sqlite3.connect(self.db_path) as conn:
            row = conn.cursor().execute("SELECT lang, full_text_en FROM articles WHERE id = 'art_3001'").fetchone()
            self.assertEqual(row[0], "en")
            self.assertEqual(row[1], "Chaldea base test")


if __name__ == "__main__":
    unittest.main()
