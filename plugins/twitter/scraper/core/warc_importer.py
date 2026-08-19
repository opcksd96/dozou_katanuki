# plugins/twitter/scraper/core/warc_importer.py (100行以下)
import os
import re
from typing import Callable, Optional
from warcio.archiveiterator import ArchiveIterator
from .mutator import Mutator
from .downloader import Downloader

try:
    from parsers.twitter_parser import TwitterParser
except ImportError:
    from ..parsers.twitter_parser import TwitterParser


class WarcImporter:
    """手動WARCファイルのオフライン自動監査・パース・メディア抽出 (SPEC-PLUGIN-001)"""

    def __init__(self, warc_path: str, db_path: str = "archive.db", storage_dir: Optional[str] = None):
        self.warc_path = warc_path
        self.mutator = Mutator(db_path=db_path)
        self.downloader = Downloader(db_path=db_path, storage_dir=storage_dir)
        self.parser = TwitterParser()

    def run_import(self, progress_callback: Optional[Callable[[int, int, str], None]] = None) -> int:
        if not os.path.exists(self.warc_path):
            if progress_callback:
                progress_callback(1, 1, f"WARC file not found: {self.warc_path}")
            return 0

        saved_posts = 0
        media_extracted = 0
        with open(self.warc_path, "rb") as stream:
            for idx, record in enumerate(ArchiveIterator(stream), start=1):
                if record.rec_type != "response":
                    continue

                uri = record.rec_headers.get_header("WARC-Target-URI") or ""
                content_type = record.http_headers.get_header("Content-Type") or "" if record.http_headers else ""

                if "json" in content_type or "html" in content_type or "status" in uri:
                    raw_bytes = record.raw_stream.read()
                    parsed = self.parser.parse_record(raw_bytes, uri)
                    if parsed and self.mutator.upsert_record(parsed):
                        saved_posts += 1
                        if progress_callback:
                            progress_callback(idx, idx + 10, f"Saved post {parsed['post']['id']} (@{parsed['account']['username']})")

                elif any(domain in uri for domain in ["twimg.com", "video.twimg.com"]):
                    media_id = os.path.basename(uri.split("?")[0])
                    dest_dir = self.downloader.storage_dir
                    os.makedirs(dest_dir, exist_ok=True)
                    dest_path = os.path.join(dest_dir, media_id)
                    with open(dest_path, "wb") as mf:
                        mf.write(record.raw_stream.read())
                    self.downloader._update_status(media_id, "COMPLETED", None, None)
                    media_extracted += 1

        if progress_callback:
            progress_callback(100, 100, f"Import finished: {saved_posts} posts, {media_extracted} media extracted.")
        return saved_posts
