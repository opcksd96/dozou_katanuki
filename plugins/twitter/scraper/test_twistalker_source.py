# plugins/twitter/scraper/test_twistalker_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, unittest

_CUR = os.path.dirname(os.path.abspath(__file__))
_ROOT = os.path.abspath(os.path.join(_CUR, "../../.."))
for d in [_CUR, _ROOT]:
    if d not in sys.path: sys.path.insert(0, d)

from plugins.twitter.scraper.sources.twistalker_source import TwistalkerSource
from plugins.twitter.scraper.parsers.twistalker_parser import parse_twistalker_api_tweets, parse_twistalker_html_tweets

SAMPLE_API_TWEET = {
    "conversation_id_str": "2096099235095560667",
    "created_at": "Mon Sep 07 01:56:42 +0000 2026",
    "full_text": "CyberCab demonstration! #robotaxi",
    "reply_count": 10, "favorite_count": 200, "retweet_count": 50, "bookmark_count": 12, "view_count": 50000,
    "core": {
        "rest_id": "44196397", "name": "Elon Musk", "screen_name": "elonmusk",
        "profile_image_url_https": "https://pbs.twimg.com/profile_images/1/avatar_normal.jpg"
    },
    "extended_entities": {
        "media": [{
            "type": "video",
            "media_url_https": "https://pbs.twimg.com/thumb.jpg",
            "video_info": {"variants": [
                {"content_type": "video/mp4", "bitrate": 500000, "url": "https://video-s.twimg.com/low.mp4"},
                {"content_type": "video/mp4", "bitrate": 2000000, "url": "https://video-s.twimg.com/high.mp4"}
            ]}
        }]
    }
}

SAMPLE_HTML = """
<div class="activity-posts">
  <div class="main-user-dts1">
    <a href="/senpai_retro"><img src="https://pbs.twimg.com/profile_images/1/senpai_normal.jpg" /></a>
    <div class="user-text3">
      <h4>Senpai Dev<span> @senpai_retro</span></h4>
      <span><a href="/senpai_retro/status/2096779662659629247">2 hours ago</a></span>
    </div>
  </div>
  <div class="activity-descp"><p>Famicom APU sound is awesome!</p></div>
  <div class="like-comment-view">
    <div class="left-comments">
      <a class="like-item"><i class="fa fa-comment"></i><span>5</span></a>
      <a class="like-item lc-left"><i class="fa fa-heart"></i><span>100</span></a>
    </div>
  </div>
</div>
"""

class TestTwistalkerSource(unittest.TestCase):
    def test_source_initialization(self):
        src = TwistalkerSource()
        self.assertEqual(src.name, "twistalker")
        self.assertTrue(src.is_available())

    def test_parse_api_tweets(self):
        records = parse_twistalker_api_tweets({"tw1": SAMPLE_API_TWEET}, "elonmusk")
        self.assertEqual(len(records), 1)
        r = records[0]
        self.assertEqual(r["post"]["id"], "2096099235095560667")
        self.assertEqual(r["account"]["username"], "elonmusk")
        self.assertNotIn("_normal.", r["account"]["avatar_original_url"])
        self.assertEqual(r["post"]["metrics"]["likes"], 200)
        self.assertEqual(len(r["media"]), 1)
        self.assertEqual(r["media"][0]["type"], "video")
        self.assertEqual(r["media"][0]["variants"][0]["bit_rate"], 2000000)

    def test_parse_html_tweets(self):
        records = parse_twistalker_html_tweets(SAMPLE_HTML, "senpai_retro")
        self.assertEqual(len(records), 1)
        r = records[0]
        self.assertEqual(r["post"]["id"], "2096779662659629247")
        self.assertEqual(r["account"]["username"], "senpai_retro")
        self.assertIn("Famicom APU sound", r["post"]["full_text"])
        self.assertEqual(r["post"]["created_at"], "2026-09-07 01:56:42")
        self.assertEqual(r["post"]["metrics"]["likes"], 100)

if __name__ == "__main__":
    unittest.main()
