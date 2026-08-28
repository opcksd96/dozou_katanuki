# plugins/base/scraper/core/media_path_resolver.py
import os
from typing import Optional


class MediaPathResolver:
    """メディアの保存先ディレクトリおよびローカルファイルパスを解決する単一責務モジュール"""
    DEFAULT_STORAGE = "G:/Media_Storage/Influencers" if os.path.exists("G:/Media_Storage/Influencers") else "blobs"

    def __init__(self, storage_dir: Optional[str] = None, platform: str = "twitter"):
        self.storage_dir = storage_dir or self.DEFAULT_STORAGE
        self.platform = platform
        os.makedirs(self.storage_dir, exist_ok=True)

    def get_target_path(self, username: str, media_id: str, media_type: str = "image") -> str:
        plat = "X(Twitter)" if self.platform.lower() in ("twitter", "base", "x") else self.platform.capitalize()
        base_name = media_id.split(":")[0]

        if "Influencers" in self.storage_dir:
            target_dir = os.path.join(self.storage_dir, username, plat, "_assets")
        elif os.path.basename(os.path.normpath(self.storage_dir)) == "blobs":
            target_dir = self.storage_dir
        else:
            category = "scenes" if media_type == "video" or media_id.endswith((".mp4", ".webm", ".m3u8")) else "images"
            target_dir = os.path.join(self.storage_dir, category, plat, username)

        os.makedirs(target_dir, exist_ok=True)
        return os.path.join(target_dir, base_name)
