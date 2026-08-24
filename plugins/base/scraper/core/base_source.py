# plugins/base/scraper/core/base_source.py (SPEC-PLUGIN-001 / 100行以下)
from abc import ABC, abstractmethod
from typing import Any, Callable, Dict, List, Optional
import requests


class BaseSource(ABC):
    """マルチソース・スクレイパー基底インターフェース"""

    def __init__(self, name: str, priority: int = 100, timeout: float = 15.0):
        self.name = name
        self.priority = priority
        self.timeout = timeout
        self.session = requests.Session()
        self.session.headers.update({
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "Accept": "*/*",
        })

    @abstractmethod
    def fetch_account(
        self, account: str, limit: int = 0, log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        """指定アカウントのタイムライン・スレッド生レコード一覧を取得"""
        pass

    @abstractmethod
    def fetch_post(
        self, post_id: str, account: str = "", log_fn: Optional[Callable[[str], None]] = None
    ) -> Optional[Dict[str, Any]]:
        """指定投稿（ID/URL）の生レコードを取得"""
        pass

    def is_available(self) -> bool:
        """ソースが現在利用可能か（ヘルスチェック）"""
        return True
