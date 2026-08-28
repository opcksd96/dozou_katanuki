# plugins/twitter/scraper/test_sotwe_source.py (SPEC-PLUGIN-001 / 100行以下)
import os, pytest
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

def test_normalize_vue_tweet_real_snowflake_id():
    rec = normalize_vue_tweet(SAMPLE_VUE_TWEET, "MsLuo14")
    assert rec["post"]["id"] == "1912895648481116616"
    assert rec["post"]["source_name"] == "sotwe"
    assert rec["post"]["sotwe_url"] == "https://www.sotwe.com/tweet/1912895648481116616"
    assert rec["post"]["original_url"] == "https://x.com/MsLuo14/status/1912895648481116616"
    assert rec["account"]["username"] == "MsLuo14"
    assert "_normal." not in rec["account"]["avatar_url"]

def test_streamsaver_and_video_variants():
    rec = normalize_vue_tweet(SAMPLE_VUE_TWEET, "MsLuo14")
    media = rec["media"]
    assert len(media) == 1
    m0 = media[0]
    assert m0["type"] == "video"
    assert m0["filename"] == "RQQWYlj86uaacmh2.mp4"
    assert m0["streamsaver_url"] == "https://jimmywarting.github.io/StreamSaver.js/www.sotwe.com/RQQWYlj86uaacmh2.mp4"
    assert len(m0["stream_variants"]) == 2
    assert any(v["content_type"] == "application/x-mpegURL" for v in m0["stream_variants"])
    assert any(v["bitrate"] == 2176000 for v in m0["stream_variants"])

def test_parse_sotwe_vue_tweets_dedup():
    records = parse_sotwe_vue_tweets([SAMPLE_VUE_TWEET, SAMPLE_VUE_TWEET], "MsLuo14")
    assert len(records) == 1
    assert records[0]["post"]["id"] == "1912895648481116616"

def test_parse_sotwe_html_fallback():
    html_sample = '<div class="tweet-card"><div class="tweet-text"><div class="dynamic-link-content">Hello</div></div></div>'
    res = parse_sotwe_html_tweets(html_sample, "test_user")
    assert len(res) == 1
    assert res[0]["post"]["full_text"] == "Hello"
    assert res[0]["account"]["username"] == "test_user"
