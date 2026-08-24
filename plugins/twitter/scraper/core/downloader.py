# plugins/twitter/scraper/core/downloader.py (Twitter特化具象ダウンローダー / 100行以下)
import os, re, sys
from typing import Any, Dict, List, Optional, Tuple

_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)
from plugins.base.scraper.core.base_downloader import BaseDownloader


class Downloader(BaseDownloader):
    """Twitter / X 特化型メディアダウンロード＆Stash注入パイプライン"""
    def __init__(self, db_path: str = "archive.db", storage_dir: Optional[str] = None):
        super().__init__(db_path=db_path, storage_dir=storage_dir, platform="twitter")

    def resolve_media_url(self, media_id: str, download_url: str, media_type: str) -> str:
        u = download_url or (f"https://pbs.twimg.com/media/{media_id}" if media_type == "image" else f"https://video.twimg.com/ext_tw_video/{media_id}")
        return u.rsplit(":", 1)[0] if any(u.endswith(s) for s in [":large", ":orig", ":small", ":medium"]) else u

    def build_metadata(self, article_id: str, username: str, display_name: str, created_at: str, wayback_url: str, full_text: str, full_text_ja: str) -> Tuple[str, str, List[str], Dict[str, Any]]:
        t_title = f"X (@{username}): Tweet {article_id}" if username and article_id else ""
        txt = f"{full_text_ja}\n\n{full_text}".strip() if full_text_ja and full_text_ja != full_text else (full_text or full_text_ja or "")
        urls = [f"https://twitter.com/{username}/status/{article_id}"] if username and article_id else []
        if wayback_url: urls.append(wayback_url)
        if article_id and username: urls.append(f"http://localhost:9999/plugin/x-timeline-middleware/index.html?view=x-timeline&performer={username}&jump_to_tweet={article_id}")
        m_ts = re.search(r'/web/(\d{14})', wayback_url) if wayback_url else None
        cf = {"tweet_id": article_id, "original_name": display_name or username, "source_system": "X_Wayback", "wayback_url": [wayback_url] if wayback_url else [], "dead_media": []}
        if m_ts: cf["wayback_timestamp"] = m_ts.group(1)
        return t_title, txt, urls, cf
