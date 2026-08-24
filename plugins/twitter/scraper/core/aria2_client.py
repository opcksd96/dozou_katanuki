# plugins/twitter/scraper/core/aria2_client.py (共通クライアント再エクスポート / 100行以下)
import os, sys

_BASE_CORE = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../base/scraper/core"))
if _BASE_CORE not in sys.path: sys.path.insert(0, _BASE_CORE)

try:
    from plugins.base.scraper.core.aria2_client import Aria2Client
except ImportError:
    from aria2_client import Aria2Client

__all__ = ["Aria2Client"]
