"""
test_sotwe_source.py (SPEC-PLUGIN-001 / 100行以下)
SotweSourceおよびSotweParserの動作を検証する単体テスト。
sotwe.htmlを用いてオフラインで実行可能。
"""
import os
from parsers.sotwe_parser import parse_sotwe_html_tweets

SOTWE_HTML_PATH = os.path.join(os.path.dirname(__file__), "..", "..", "..", "sotwe.html")

def load_sotwe_html():
    with open(SOTWE_HTML_PATH, "r", encoding="utf-8") as f:
        return f.read()

def test_parse_sotwe_html_tweets_basic():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    assert len(results) > 0

def test_avatar_url_and_profile_history():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    for tweet in results:
        acc = tweet["account"]
        av = acc["avatar_url"]
        assert "_normal." not in av, f"_normal. should be removed: {av}"
        assert "_400x400." not in av, f"_400x400. should not be present: {av}"
        if av:
            assert len(acc["profile_history"]) > 0
            assert acc["profile_history"][0]["avatar_original_url"] == av

def test_post_datetime_not_now():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    # 最初のツイートは 2023年のピン留めツイート
    first_dt = results[0]["post"]["created_at"]
    assert "2023-12-28" in first_dt, f"Expected 2023-12-28, got {first_dt}"
    assert results[0]["post"]["is_pinned"] is True

def test_snowflake_id_length():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    for tweet in results:
        post_id = tweet["post"]["id"]
        assert len(post_id) >= 18, f"Tweet ID should be a 18-19 digit Snowflake ID, got {post_id}"

def test_engagement_metrics():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    # ピン留めツイートのメトリクス検証 (likes: 117, retweets: 4)
    first_metrics = results[0]["post"]["metrics"]
    assert first_metrics["likes"] == 117
    assert first_metrics["retweets"] == 4
    assert "replies" in first_metrics

def test_media_extraction_and_query_preservation():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    video_tweets = [t for t in results if any(m["type"] == "video" for m in t["media"])]
    assert len(video_tweets) > 0
    for tweet in video_tweets:
        for m in tweet["media"]:
            if m["type"] == "video":
                assert "video.twimg.com" in m["url"]
                assert "?tag=" in m["url"]
