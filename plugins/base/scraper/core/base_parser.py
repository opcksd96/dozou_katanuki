# plugins/base/scraper/core/base_parser.py (SPEC-PLUGIN-001 / 100行以下)
from abc import ABC, abstractmethod
from typing import Any, Dict, Optional


class BaseParser(ABC):
    """プラットフォーム固有の生レコードを共通中間表現へパースする抽象基底クラス"""

    @abstractmethod
    def parse_record(self, raw_content: str, url: str) -> Optional[Dict[str, Any]]:
        """
        生テキスト/パケットと取得元URLから共通構造化JSONを生成して返す
        返却フォーマット:
        {
            "post": {
                "id": str, "numeric_id": str, "account_id": str, "account_username": str,
                "full_text": str, "created_at": str, "wayback_url": str,
                "conversation_id": str, "reply_to_id": str, "reply_to_handle": str,
                "urls": [{"short_url": str, "expanded_url": str}],
                "hashtags": [str], "mentions": [str]
            },
            "account": {
                "numeric_id": str, "username": str, "display_name": str,
                "avatar_url": str, "description": str, "observed_at": str
            },
            "media": [
                {
                    "media_id": str, "type": "image"|"video"|"animated_gif",
                    "original_url": str, "download_url": str, "local_path": str,
                    "width": int, "height": int, "bitrate": int, "duration": float
                }
            ]
        }
        """
        pass
