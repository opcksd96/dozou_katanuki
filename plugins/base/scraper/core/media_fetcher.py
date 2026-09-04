# plugins/base/scraper/core/media_fetcher.py
import os
import requests
from typing import Any, Dict, List, Optional, Tuple


class MediaFetcher:
    """HTTPダウンロード試行、バイナリ完全性検証、Stash登録を実行する単一責務モジュール"""

    def __init__(self, session: Optional[requests.Session] = None, reconciler: Any = None):
        self.session = session or requests.Session()
        self.session.headers.update({
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
        })
        self.reconciler = reconciler

    @staticmethod
    def is_valid_binary(file_path: str) -> bool:
        if not (os.path.exists(file_path) and os.path.getsize(file_path) > 0):
            return False
        try:
            with open(file_path, "rb") as f:
                header = f.read(32)
            invalid_prefixes = (b"<!DOCTYPE", b"<html", b"<HTML", b"{\n", b'{"', b"<?xml")
            return not any(header.startswith(pfx) for pfx in invalid_prefixes)
        except Exception:
            return False

    def try_fetch_direct(
        self,
        dest_path: str,
        media_id: str,
        resolved_url: str,
        media_type: str,
        meta_tuple: Tuple[str, str, List[str], Dict[str, Any]],
        thumbnail_url: str = "",
        created_at: str = "",
        username: str = "",
        display_name: str = "",
        fallback_urls: Optional[List[str]] = None,
    ) -> Tuple[bool, Optional[str], Optional[str], Optional[str], Optional[str]]:
        t_title, txt, urls_list, c_fields = meta_tuple

        # すでにローカルに正常バイナリが存在する場合
        if self.is_valid_binary(dest_path):
            reg_id = self.reconciler.register_media(
                dest_path, media_type, title=t_title, details=txt, urls=urls_list,
                date=created_at[:10], username=username, display_name=display_name,
                thumbnail_url=thumbnail_url, custom_fields=c_fields
            )
            return True, "COMPLETED", reg_id if media_type == "image" else None, reg_id if media_type != "image" else None, "local"

        candidates = fallback_urls if fallback_urls else [resolved_url]
        for cand_url in candidates:
            qual_name = "original" if cand_url == resolved_url else "variant"
            try:
                timeout = 8 if media_type == "image" else 20
                resp = self.session.get(cand_url, stream=True, timeout=timeout, allow_redirects=True)
                if resp.status_code == 200 and "text" not in resp.headers.get("Content-Type", ""):
                    with open(dest_path, "wb") as f:
                        for chunk in resp.iter_content(chunk_size=65536):
                            if chunk:
                                f.write(chunk)
                    if self.is_valid_binary(dest_path):
                        reg_id = self.reconciler.register_media(
                            dest_path, media_type, title=t_title, details=txt, urls=urls_list,
                            date=created_at[:10], username=username, display_name=display_name,
                            thumbnail_url=thumbnail_url, custom_fields=c_fields
                        )
                        return True, "COMPLETED", reg_id if media_type == "image" else None, reg_id if media_type != "image" else None, qual_name
                    elif os.path.exists(dest_path):
                        os.remove(dest_path)
            except Exception:
                continue

        return False, "RETAINED", None, None, None
