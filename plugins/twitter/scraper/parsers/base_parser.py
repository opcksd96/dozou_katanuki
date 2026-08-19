# plugins/twitter/scraper/parsers/base_parser.py (100行以下)
from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional


class BaseParser(ABC):
    """共通パーサー抽象基底クラス (SPEC-PLUGIN-001)"""

    @abstractmethod
    def parse_record(self, raw_data: Any, uri: str) -> Optional[Dict[str, Any]]:
        """生データ（JSON辞書またはHTML）から共通中間表現辞書を抽出します。"""
        pass

    @abstractmethod
    def detect_platform_and_account(self, uri: str) -> Optional[Dict[str, str]]:
        """URLパターンからプラットフォームおよびアカウント名を自動検出します。"""
        pass
