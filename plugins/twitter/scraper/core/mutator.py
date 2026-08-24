# plugins/twitter/scraper/core/mutator.py (Twitter特化具象ミューテーター / 100行以下)
import os, sys
from typing import Optional

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_mutator import BaseMutator
from plugins.base.scraper.core.translator import Translator


class Mutator(BaseMutator):
    """Twitter / X 特化型ドライバーミューテーター (GORM/SQLite3 連携)"""
    def __init__(self, db_path: str = "archive.db", translator: Optional[Translator] = None, enable_translation: bool = True):
        super().__init__(db_path=db_path, platform="twitter", translator=translator, enable_translation=enable_translation)
