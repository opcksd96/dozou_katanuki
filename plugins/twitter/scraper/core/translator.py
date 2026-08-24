# plugins/twitter/scraper/core/translator.py (共通翻訳エンジン再エクスポート / 100行以下)
import os, sys

_BASE_CORE = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../base/scraper/core"))
if _BASE_CORE not in sys.path: sys.path.insert(0, _BASE_CORE)

try:
    from plugins.base.scraper.core.translator import Translator
except ImportError:
    from translator import Translator

__all__ = ["Translator"]
