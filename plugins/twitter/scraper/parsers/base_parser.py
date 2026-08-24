# plugins/twitter/scraper/parsers/base_parser.py (SPEC-PLUGIN-001 / 100行以下)
import os, sys

# プラグインベースのパスを解決
_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_parser import BaseParser

__all__ = ["BaseParser"]
