# plugins/twitter/scraper/test_multi_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, unittest
from unittest.mock import MagicMock, patch

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.abspath(os.path.join(CURRENT_DIR, "../../.."))
for d in [PROJECT_ROOT, CURRENT_DIR]:
    if d not in sys.path: sys.path.insert(0, d)

from plugins.base.scraper.core.source_orchestrator import SourceOrchestrator
from plugins.twitter.scraper.sources import OfficialSource, SotweSource, TwistalkerSource, NitterSource, WaybackSource
from plugins.twitter.scraper.parsers.twitter_parser import TwitterParser
from plugins.twitter.scraper.core.scraper import Scraper


class TestMultiSourceScraper(unittest.TestCase):
    def setUp(self):
        self.parser = TwitterParser()
        self.scraper = Scraper(platform="twitter")

    def test_orchestrator_registration(self):
        sources = self.scraper.orchestrator.list_sources()
        self.assertIn("official", sources)
        self.assertIn("sotwe", sources)
        self.assertIn("twistalker", sources)
        self.assertIn("nitter", sources)
        self.assertIn("wayback", sources)

    def test_sotwe_json_parsing(self):
        sample_sotwe = {
            "id": "1800000000000000000",
            "conversation_id_str": "1800000000000000000",
            "text": "Hello from Sotwe mirror! https://t.co/xyz",
            "user": {"screen_name": "mash_test", "name": "Mash Kyrielight", "avatar": "https://pbs.twimg.com/profile_images/1/avatar.jpg"},
            "mediaEntities": [{"url": "https://pbs.twimg.com/media/sample.jpg", "type": "image"}],
            "urls": [{"url": "https://t.co/xyz", "expanded_url": "https://example.com/dest"}]
        }
        res = self.parser.parse_record(sample_sotwe, "https://twitter.com/mash_test/status/1800000000000000000")
        self.assertIsNotNone(res)
        self.assertEqual(res["account"]["username"], "mash_test")
        self.assertEqual(res["post"]["id"], "1800000000000000000")
        self.assertEqual(res["post"]["conversation_id"], "1800000000000000000")
        self.assertEqual(len(res["media"]), 1)
        self.assertEqual(res["media"][0]["url"], "https://pbs.twimg.com/media/sample.jpg")

    def test_twistalker_html_parsing(self):
        sample_html = '''
        <div class="post" data-id="1700000000000000000">
            <strong class="fullname">Senpai Dev</strong>
            <span class="username">@senpai_retro</span>
            <div class="post-text">APU Sound reproduction is amazing!</div>
            <img src="https://pbs.twimg.com/media/nes_sound.png" />
        </div>
        '''
        res = self.parser.parse_record(sample_html, "https://twistalker.com/senpai_retro/status/1700000000000000000")
        self.assertIsNotNone(res)
        self.assertEqual(res["post"]["id"], "1700000000000000000")
        self.assertEqual(res["account"]["username"], "senpai_retro")
        self.assertIn("APU Sound", res["post"]["full_text"])
        self.assertEqual(len(res["media"]), 1)

    def test_nitter_html_parsing(self):
        sample_html = '''
        <div class="timeline-item">
            <a class="fullname" title="Mash">Mash</a>
            <a class="username" title="@mash">@mash</a>
            <div class="tweet-content media-body">Shield activated.</div>
            <a href="/mash/status/1600000000000000000#m">Dec 31, 2025</a>
        </div>
        '''
        res = self.parser.parse_record(sample_html, "https://nitter.net/mash/status/1600000000000000000")
        self.assertIsNotNone(res)
        self.assertEqual(res["post"]["id"], "1600000000000000000")
        self.assertEqual(res["account"]["username"], "mash")
        self.assertIn("Shield activated.", res["post"]["full_text"])


if __name__ == "__main__":
    unittest.main()
