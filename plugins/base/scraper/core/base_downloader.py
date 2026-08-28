# plugins/base/scraper/core/base_downloader.py (SPEC-PLUGIN-001)
import os
import sqlite3
import requests
from typing import Any, Callable, Dict, List, Optional, Tuple

from .aria2_client import Aria2Client
from .stash_client import StashClient
from .stash_reconciler import StashReconciler
from .downloader_pipeline import DownloaderPipelineHelper
from .media_path_resolver import MediaPathResolver
from .media_fetcher import MediaFetcher
from .media_queue_processor import MediaQueueProcessor


class BaseDownloader:
    """メディア原本取得 & Stash登録 & Motrix委託 & 品質追跡の統合ファサード"""

    def __init__(self, db_path: str = "archive.db", storage_dir: Optional[str] = None, platform: str = "twitter"):
        self.db_path = db_path
        self.platform = platform
        self.path_resolver = MediaPathResolver(storage_dir, platform)
        self.storage_dir = self.path_resolver.storage_dir

        self.session = requests.Session()
        self.stash = StashClient()
        self.aria2 = Aria2Client()
        self.reconciler = StashReconciler(self.stash)
        self.fetcher = MediaFetcher(self.session, self.reconciler)
        self.pipeline_helper = DownloaderPipelineHelper(self)
        self.queue_processor = MediaQueueProcessor(self)

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=60.0)
        conn.execute("PRAGMA journal_mode = WAL;")
        conn.execute("PRAGMA synchronous = NORMAL;")
        conn.execute("PRAGMA foreign_keys = ON;")
        conn.execute("PRAGMA busy_timeout = 60000;")
        return conn

    def get_target_path(self, username: str, media_id: str, media_type: str = "image") -> str:
        return self.path_resolver.get_target_path(username, media_id, media_type)

    def resolve_media_url(self, media_id: str, download_url: str, media_type: str) -> str:
        return download_url

    def build_metadata(
        self, article_id: str, username: str, display_name: str,
        created_at: str, wayback_url: str, full_text: str, full_text_ja: str
    ) -> Tuple[str, str, List[str], Dict[str, Any]]:
        t_title = f"{self.platform.capitalize()} (@{username}): Post {article_id}" if username and article_id else ""
        txt = f"{full_text_ja}\n\n{full_text}".strip() if full_text_ja and full_text_ja != full_text else (full_text or full_text_ja or "")
        urls = [wayback_url] if wayback_url else []
        custom_fields = {
            "post_id": article_id, "original_name": display_name or username,
            "source_system": self.platform, "wayback_url": urls, "dead_media": []
        }
        return t_title, txt, urls, custom_fields

    def _update_status(
        self, media_id: str, status: str, reason: Optional[str] = None,
        img_id: Optional[str] = None, scn_id: Optional[str] = None, quality: Optional[str] = None
    ) -> None:
        with self._get_conn() as conn:
            conn.cursor().execute(
                """UPDATE media SET download_status = ?, failed_reason = coalesce(?, failed_reason),
                   stash_image_id = coalesce(?, stash_image_id), stash_scene_id = coalesce(?, stash_scene_id),
                   media_quality = coalesce(?, media_quality) WHERE media_id = ?""",
                (status, reason, img_id, scn_id, quality, media_id)
            )
            conn.commit()

    def process_queued_media(
        self, article_id: Optional[str] = None, media_id: Optional[str] = None,
        log_fn: Optional[Callable[[int, int, str], None]] = None
    ) -> int:
        return self.queue_processor.process_queued(article_id, media_id, log_fn=log_fn)

    def escalate_dead_media(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        return self.pipeline_helper.escalate_dead_media(log_fn=log_fn)

    def poll_outsourced_media(self, log_fn: Optional[Callable[[str], None]] = None) -> int:
        return self.reconciler.reconcile_to_db(self.db_path, log_fn=log_fn)

    def smart_recovery_pipeline(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> dict:
        return self.pipeline_helper.run_smart_recovery(log_fn=log_fn)

    def escalate_to_thunder(self, log_fn: Optional[Callable[[int, int, str], None]] = None) -> int:
        return self.pipeline_helper.escalate_to_thunder(log_fn=log_fn)

    def clean_failed_outsourced(self, log_fn: Optional[Callable[[str], None]] = None) -> int:
        return self.pipeline_helper.clean_failed_outsourced(log_fn=log_fn)
