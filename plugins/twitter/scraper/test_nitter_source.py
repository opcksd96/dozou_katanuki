# plugins/twitter/scraper/test_nitter_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, unittest
from unittest.mock import MagicMock, patch

_CUR = os.path.dirname(os.path.abspath(__file__))
_ROOT = os.path.abspath(os.path.join(_CUR, "../../.."))
for d in [_CUR, _ROOT]:
    if d not in sys.path: sys.path.insert(0, d)

from plugins.twitter.scraper.sources.nitter_source import NitterSource
from plugins.twitter.scraper.parsers.nitter_parser import parse_nitter_card, parse_nitter_html_tweets
from plugins.twitter.scraper.parsers.nitter_media import normalize_twitter_image_url

SAMPLE_NITTER_HTML = """
<div class="timeline-item">
  <div class="tweet-header">
    <a class="tweet-avatar"><img src="https://nitter.space/pic/profile_images%2F123%2Fmash_bigger.jpg" /></a>
    <a class="fullname">Mash Kyrielight</a>
    <a class="username">@mash_retro</a>
    <span class="tweet-date"><a href="/mash_retro/status/2096779662659629247#m" title="Sep 7, 2026 · 1:56 AM UTC"></a></span>
  </div>
  <div class="tweet-content media-body">Senpai, NES APU sound is decoded!<br/>Authentic 2A03 registers.</div>
  <div class="attachments">
    <div class="attachment image"><a class="still-image" href="https://nitter.space/pic/orig/media%2Ffamicom_board.jpg"></a></div>
    <div class="attachment video"><video poster="https://nitter.space/pic/video_thumb.jpg"><source src="https://video.twimg.com/ext_tw_video/123/pu/vid/1280x720/apu_demo.mp4" /></video></div>
  </div>
  <div class="tweet-stats">
    <span class="tweet-stat"><div class="icon-comment"></div> 12</span>
    <span class="tweet-stat"><div class="icon-retweet"></div> 45</span>
    <span class="tweet-stat"><div class="icon-heart"></div> 2.5K</span>
  </div>
</div>
"""

class TestNitterSource(unittest.TestCase):
    def test_clone_url_and_mock(self):
        src = NitterSource()
        self.assertEqual(src.get_clone_url(), "https://nitter.space")
        mocked_src = NitterSource(clone_url="https://custom-nitter.internal")
        self.assertEqual(mocked_src.get_clone_url(), "https://custom-nitter.internal")

    def test_image_url_normalization(self):
        url1 = "https://pbs.twimg.com/media/sample.jpg?format=jpg&name=medium"
        self.assertEqual(normalize_twitter_image_url(url1), "https://pbs.twimg.com/media/sample.jpg?name=orig")
        url2 = "https://nitter.space/pic/orig/media%2Ftest_img.png"
        self.assertEqual(normalize_twitter_image_url(url2), "https://pbs.twimg.com/media/test_img.png?name=orig")

    def test_parse_html_tweets(self):
        posts = parse_nitter_html_tweets(SAMPLE_NITTER_HTML, default_account="mash_retro")
        self.assertEqual(len(posts), 1)
        p = posts[0]
        self.assertEqual(p["post"]["id"], "2096779662659629247")
        self.assertEqual(p["account"]["username"], "mash_retro")
        self.assertEqual(p["account"]["display_name"], "Mash Kyrielight")
        self.assertIn("Senpai, NES APU sound", p["post"]["full_text"])
        self.assertEqual(p["post"]["metrics"]["likes"], 2500)
        self.assertEqual(len(p["media"]), 2)
        v = next(m for m in p["media"] if m["type"] == "video")
        self.assertEqual(v["width"], 1280); self.assertEqual(v["height"], 720)

    @patch("plugins.twitter.scraper.sources.nitter_source.SB")
    def test_fetch_account_mock(self, mock_sb_cls):
        mock_sb = MagicMock()
        mock_sb.get_page_source.return_value = SAMPLE_NITTER_HTML
        mock_sb_cls.return_value.__enter__.return_value = mock_sb
        src = NitterSource(clone_url="https://mock.nitter.test")
        results = src.fetch_account("mash_retro", limit=1)
        self.assertEqual(len(results), 1)
        self.assertEqual(results[0]["post"]["id"], "2096779662659629247")

if __name__ == "__main__":
    unittest.main()
