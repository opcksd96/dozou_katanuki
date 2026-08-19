# plugins/twitter/scraper/parsers/__init__.py
from .base_parser import BaseParser
from .twitter_parser import TwitterParser

__all__ = ["BaseParser", "TwitterParser"]
