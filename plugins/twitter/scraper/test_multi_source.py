# plugins/twitter/scraper/test_multi_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, unittest

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.abspath(os.path.join(CURRENT_DIR, "../../.."))
for d in [PROJECT_ROOT, CURRENT_DIR]:
    if d not in sys.path: sys.path.insert(0, d)

from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from plugins.twitter.scraper.core.scraper import Scraper
from plugins.twitter.scraper.parsers.x_parser import XParser
from plugins.twitter.scraper.parsers.nitter_parser import parse_nitter_html_tweets

class TestMultiSourceScraper(unittest.TestCase):
    def setUp(self):
        self.parser = TwitterParser()
        self.scraper = Scraper(platform="twitter")

    def test_orchestrator_registration(self):
        sources = self.scraper.orchestrator.list_sources()
        for src in ["official", "sotwe", "twistalker", "nitter", "wayback"]:
            self.assertIn(src, sources)

    def test_sotwe_json_parsing(self):
        sample_sotwe = {
            "id": "1800000000000000000", "conversation_id_str": "1800000000000000000",
            "text": "Hello from Sotwe mirror! https://t.co/xyz",
            "user": {"screen_name": "mash_test", "name": "Mash Kyrielight", "avatar": "https://pbs.twimg.com/avatar.jpg"},
            "mediaEntities": [{"url": "https://pbs.twimg.com/media/sample.jpg", "type": "image"}],
            "urls": [{"url": "https://t.co/xyz", "expanded_url": "https://example.com/dest"}]
        }
        res = self.parser.parse_record(sample_sotwe, "https://twitter.com/mash_test/status/1800000000000000000")
        self.assertIsNotNone(res)
        self.assertEqual(res["account"]["username"], "mash_test")
        self.assertEqual(res["post"]["id"], "1800000000000000000")
        self.assertEqual(len(res["media"]), 1)

    def test_twistalker_html_parsing(self):
        sample_html = '''<div class="post" data-id="1700000000000000000"><strong class="fullname">Senpai Dev</strong><span class="username">@senpai_retro</span><div class="post-text">APU Sound reproduction is amazing!</div><img src="https://pbs.twimg.com/media/nes_sound.png" /></div>'''
        res = self.parser.parse_record(sample_html, "https://twistalker.com/senpai_retro/status/1700000000000000000")
        self.assertIsNotNone(res)
        self.assertEqual(res["post"]["id"], "1700000000000000000")
        self.assertIn("APU Sound", res["post"]["full_text"])

    def test_nitter_html_parsing(self):
        sample_html = '''<div class="timeline-item"><a class="fullname">Mash</a><a class="username">@mash</a><div class="tweet-content">Shield activated.</div><span class="tweet-date"><a href="/mash/status/1600000000000000000#m">Dec 31, 2025</a></span></div>'''
        posts = parse_nitter_html_tweets(sample_html, default_account="mash")
        self.assertEqual(len(posts), 1)
        self.assertEqual(posts[0]["post"]["id"], "1600000000000000000")
        self.assertIn("Shield activated.", posts[0]["post"]["full_text"])

    def test_x_syndication_parsing(self):
        sample_x = {"id_str": "1900000000000000000", "text": "X test tweet", "user": {"screen_name": "x_dev", "name": "X Dev"}}
        res = XParser().parse_syndication_json(sample_x, default_acc="x_dev")
        self.assertIsNotNone(res)
        self.assertEqual(res["post"]["id"], "1900000000000000000")
        self.assertEqual(res["account"]["username"], "x_dev")

if __name__ == "__main__":
    unittest.main()
