# plugins/twitter/scraper/test_sotwe_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys, unittest

_CUR = os.path.dirname(os.path.abspath(__file__))
_ROOT = os.path.abspath(os.path.join(_CUR, "../../.."))
for d in [_CUR, _ROOT]:
    if d not in sys.path: sys.path.insert(0, d)

from parsers.sotwe_parser import parse_sotwe_html_tweets, parse_sotwe_vue_tweets
from parsers.sotwe_extractors import normalize_vue_tweet, build_streamsaver_url

SAMPLE_VUE_TWEET = {
    "id": "1912895648481116616", "id_str": "1912895648481116616", "createdAt": 1744937237000,
    "text": "学业很忙，所以不再开通线下服务。", "replyCount": 0, "retweetCount": 2, "favoriteCount": 38, "viewCount": 110260,
    "user": {"screenName": "MsLuo14", "name": "小罗老师", "profileImage": "https://pbs.twimg.com/profile_images/1877640503430438912/vNUa5of2_normal.jpg"},
    "mediaEntities": [{
        "id": "1919718174796480512", "type": "video", "videoURL": "https://video-s.twimg.com/amplify_video/1919718174796480512/vid/avc1/720x1280/RQQWYlj86uaacmh2.mp4?tag=14",
        "mediaURL": "https://pbs.twimg.com/amplify_video_thumb/1919718174796480512/img/ndZc6DJcXtuYM3Zs.jpg",
        "videoInfo": {"variants": [
            {"type": "application/x-mpegURL", "bitrate": 0, "url": "https://video-s.twimg.com/amplify_video/1919718174796480512/pl/NvfLVZvz9xiwBSZ9.m3u8?tag=14"},
            {"type": "video/mp4", "bitrate": 2176000, "url": "https://video-s.twimg.com/amplify_video/1919718174796480512/vid/avc1/720x1280/RQQWYlj86uaacmh2.mp4?tag=14"}
        ]}
    }]
}

SAMPLE_UNRELATED_TWEET = {
    "id": "1999999999999999999", "id_str": "1999999999999999999", "createdAt": 1744937237000,
    "text": "Trending ad tweet from sidebar", "replyCount": 1, "retweetCount": 10, "favoriteCount": 50,
    "user": {"screenName": "random_sidebar_user", "name": "Random", "profileImage": ""},
    "mediaEntities": [{"id": "m_unrelated", "type": "image", "mediaURL": "https://pbs.twimg.com/media/unrelated.jpg"}]
}

class TestSotweSource(unittest.TestCase):
    def test_normalize_vue_tweet_real_snowflake_id(self):
        rec = normalize_vue_tweet(SAMPLE_VUE_TWEET, "MsLuo14")
        self.assertEqual(rec["post"]["id"], "1912895648481116616")
        self.assertEqual(rec["post"]["source_name"], "sotwe")
        self.assertEqual(rec["post"]["sotwe_url"], "https://www.sotwe.com/tweet/1912895648481116616")
        self.assertEqual(rec["account"]["username"], "MsLuo14")
        self.assertNotIn("_normal.", rec["account"]["avatar_url"])

    def test_streamsaver_and_video_variants(self):
        rec = normalize_vue_tweet(SAMPLE_VUE_TWEET, "MsLuo14")
        media = rec["media"]
        self.assertEqual(len(media), 1)
        self.assertEqual(media[0]["type"], "video")
        self.assertEqual(media[0]["filename"], "RQQWYlj86uaacmh2.mp4")
        self.assertEqual(media[0]["streamsaver_url"], "https://jimmywarting.github.io/StreamSaver.js/www.sotwe.com/RQQWYlj86uaacmh2.mp4")

    def test_parse_sotwe_vue_tweets_filters_unrelated_accounts(self):
        records = parse_sotwe_vue_tweets([SAMPLE_VUE_TWEET, SAMPLE_UNRELATED_TWEET], "MsLuo14")
        self.assertEqual(len(records), 1)
        self.assertEqual(records[0]["post"]["id"], "1912895648481116616")

    def test_parse_sotwe_vue_tweets_allows_whitelisted_accounts(self):
        records = parse_sotwe_vue_tweets([SAMPLE_VUE_TWEET, SAMPLE_UNRELATED_TWEET], "MsLuo14", whitelist={"random_sidebar_user"})
        self.assertEqual(len(records), 2)

    def test_parse_sotwe_html_fallback(self):
        html_sample = '<div class="tweet-card"><div class="tweet-profile"><a href="/test_user"></a></div><div class="tweet-text"><div class="dynamic-link-content">Hello</div></div></div><div class="tweet-card"><div class="tweet-profile"><a href="/stranger"></a></div></div>'
        res = parse_sotwe_html_tweets(html_sample, "test_user")
        self.assertEqual(len(res), 1)
        self.assertEqual(res[0]["post"]["full_text"], "Hello")
        self.assertEqual(res[0]["account"]["username"], "test_user")

if __name__ == "__main__":
    unittest.main()

