"""
test_sotwe_source.py
SotweSourceおよびSotweParserの動作を検証する単体テスト。
モックデータを用いてオフラインで実行可能。
"""
from unittest.mock import MagicMock, patch
import pytest

from parsers.sotwe_parser import extract_media_entities, parse_sotwe_post, parse_timestamp
from sources.sotwe_source import SotweSource


MOCK_USER = {
    "id": "123456789",
    "screenName": "target_user",
    "name": "Target User",
    "profileImage": "https://example.com/avatar.jpg",
    "description": "Test profile description"
}

MOCK_POST_ITEM = {
    "id": "987654321098765432",
    "conversationId": "987654321098765432",
    "inReplyToStatusId": "111222333444555666",
    "inReplyToScreenName": "replied_user",
    "createdAt": 1724500000000,
    "text": "Hello world from Sotwe! @replied_user",
    "mediaEntities": [
        {
            "type": "photo",
            "url": "https://example.com/image.jpg",
            "width": 1200,
            "height": 800
        }
    ],
    "user": MOCK_USER
}


def test_parse_timestamp():
    ts = parse_timestamp(1724500000000)
    assert ts == "2024-08-24T11:46:40Z"


def test_extract_media_entities():
    raw_media = MOCK_POST_ITEM["mediaEntities"]
    extracted = extract_media_entities(raw_media)
    assert len(extracted) == 1
    assert extracted[0]["url"] == "https://example.com/image.jpg"
    assert extracted[0]["type"] == "image"
    assert extracted[0]["width"] == 1200


def test_parse_sotwe_post():
    parsed = parse_sotwe_post(MOCK_POST_ITEM)
    assert parsed is not None
    assert parsed["platform"] == "twitter"
    assert parsed["account"]["username"] == "target_user"
    assert parsed["post"]["id"] == "987654321098765432"
    assert parsed["post"]["reply_to_tweet_id"] == "111222333444555666"
    assert len(parsed["media"]) == 1


@patch("requests.Session.get")
def test_fetch_account(mock_get):
    mock_resp = MagicMock()
    mock_resp.status_code = 200
    mock_resp.json.return_value = {
        "data": [MOCK_POST_ITEM],
        "user": MOCK_USER
    }
    mock_get.return_value = mock_resp

    source = SotweSource()
    posts = source.fetch_account("target_user", limit=10)

    assert len(posts) == 1
    assert posts[0]["account"]["username"] == "target_user"
    assert posts[0]["post"]["id"] == "987654321098765432"


@patch("requests.Session.get")
def test_fetch_post(mock_get):
    mock_resp = MagicMock()
    mock_resp.status_code = 200
    mock_resp.json.return_value = MOCK_POST_ITEM
    mock_get.return_value = mock_resp

    source = SotweSource()
    post = source.fetch_post("987654321098765432")

    assert post is not None
    assert post["post"]["id"] == "987654321098765432"
