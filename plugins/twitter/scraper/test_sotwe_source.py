"""
test_sotwe_source.py
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


def test_avatar_url_normal_removed():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    for tweet in results:
        av = tweet["account"]["avatar_url"]
        assert "_normal." not in av, f"_normal. should be removed from avatar_url: {av}"
        assert "_400x400." not in av, f"_400x400. should not be present: {av}"


def test_video_media_extraction():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    video_tweets = [t for t in results if any(m["type"] == "video" for m in t["media"])]
    assert len(video_tweets) > 0, "At least one tweet should contain video media"
    for tweet in video_tweets:
        video_media = [m for m in tweet["media"] if m["type"] == "video"]
        assert len(video_media) > 0
        for vm in video_media:
            assert "video.twimg.com" in vm["url"]
            assert vm["width"] == 0
            assert vm["height"] == 0


def test_video_poster_extraction():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    video_tweets = [t for t in results if any(m["type"] == "video" for m in t["media"])]
    assert len(video_tweets) > 0
    for tweet in video_tweets:
        poster_media = [m for m in tweet["media"] if m["type"] == "image" and "ext_tw_video_thumb" in m["url"]]
        assert len(poster_media) > 0, "Video tweets should include poster thumbnail images"
        for pm in poster_media:
            assert pm["width"] == 0
            assert pm["height"] == 0


def test_media_url_preserved_without_strip():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    for tweet in results:
        for m in tweet["media"]:
            if m["type"] == "video":
                assert "?tag=" in m["url"], "Video query parameters should be preserved"


def test_image_media_extraction():
    html = load_sotwe_html()
    results = parse_sotwe_html_tweets(html, "Yike_Luo")
    image_tweets = [t for t in results if any(m["type"] == "image" for m in t["media"])]
    assert len(image_tweets) > 0, "At least one tweet should contain image media"
