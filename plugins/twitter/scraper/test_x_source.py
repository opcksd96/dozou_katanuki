# plugins/twitter/scraper/test_x_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, unittest
from bs4 import BeautifulSoup
from unittest.mock import MagicMock, patch

_CUR = os.path.dirname(os.path.abspath(__file__))
_ROOT = os.path.abspath(os.path.join(_CUR, "../../.."))
for d in [_CUR, _ROOT]:
    if d not in sys.path: sys.path.insert(0, d)

from plugins.twitter.scraper.sources.official_source import OfficialSource
from plugins.twitter.scraper.parsers.x_parser import XParser

SAMPLE_SYNDICATION_JSON = {
    "id_str": "2096779662659629247",
    "text": "NES 2A03 APU Sheet Music reproduction video published! #nes",
    "created_at": "Mon Sep 07 01:56:42 +0000 2026",
    "user": {"name": "Senpai", "screen_name": "senpai_retro", "profile_image_url_https": "https://pbs.twimg.com/avatar_normal.jpg"},
    "favorite_count": 1500, "retweet_count": 320, "reply_count": 42,
    "mediaDetails": [{
        "type": "video", "media_url_https": "https://pbs.twimg.com/media/thumb.jpg",
        "video_info": {"variants": [
            {"content_type": "video/mp4", "bitrate": 800000, "url": "https://video.twimg.com/low.mp4"},
            {"content_type": "video/mp4", "bitrate": 2100000, "url": "https://video.twimg.com/high.mp4"}
        ]}
    }]
}

SAMPLE_ARTICLE_HTML = """
<article>
  <a href="/senpai_retro/status/2096779662659629247"></a>
  <div dir="auto" class="css-146c3p1 r-bcqeeo r-1ttztb7 r-qvutc0 font-chirp">Testing NES triangle wave!</div>
  <a aria-label="1.2万 いいね"><span class="tabular-nums">1.2万</span></a>
  <img src="https://pbs.twimg.com/media/nes_chip.jpg:large" />
</article>
"""

class TestXSource(unittest.TestCase):
    def setUp(self):
        self.parser = XParser()

    def test_cookie_injection_mock(self):
        mock_cookies = {"auth_token": "mocked_auth_12345", "ct0": "mocked_csrf_67890"}
        src = OfficialSource(cookies=mock_cookies)
        self.assertEqual(src.cookies, mock_cookies)
        mock_sb = MagicMock()
        src._inject_cookies(mock_sb)
        self.assertEqual(mock_sb.add_cookie.call_count, 2)

    def test_parse_syndication_json(self):
        res = self.parser.parse_syndication_json(SAMPLE_SYNDICATION_JSON, default_acc="senpai_retro")
        self.assertIsNotNone(res)
        self.assertEqual(res["post"]["id"], "2096779662659629247")
        self.assertEqual(res["account"]["username"], "senpai_retro")
        self.assertEqual(res["post"]["metrics"]["likes"], 1500)
        self.assertEqual(len(res["media"]), 1)
        self.assertEqual(res["media"][0]["type"], "video")
        self.assertEqual(res["media"][0]["download_url"], "https://video.twimg.com/high.mp4")

    def test_parse_html_card(self):
        soup = BeautifulSoup(SAMPLE_ARTICLE_HTML, "html.parser")
        res = self.parser.parse_html_card(soup.find("article"), default_acc="senpai_retro")
        self.assertIsNotNone(res)
        self.assertEqual(res["post"]["id"], "2096779662659629247")
        self.assertIn("Testing NES triangle wave!", res["post"]["full_text"])
        self.assertEqual(res["post"]["metrics"]["likes"], 12000)
        self.assertEqual(len(res["media"]), 1)
        self.assertIn("name=orig", res["media"][0]["url"])

    @patch("requests.Session.get")
    def test_fetch_post_syndication_mock(self, mock_get):
        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.json.return_value = SAMPLE_SYNDICATION_JSON
        mock_get.return_value = mock_resp
        src = OfficialSource()
        p = src.fetch_post("2096779662659629247", account="senpai_retro")
        self.assertIsNotNone(p)
        self.assertEqual(p["post"]["id"], "2096779662659629247")

if __name__ == "__main__":
    unittest.main()
