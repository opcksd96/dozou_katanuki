# plugins/twitter/scraper/core/warc_importer.py (100行以下)
import os, re
from typing import Any, Callable, Dict, Optional
from warcio.archiveiterator import ArchiveIterator
from .mutator import Mutator
from .downloader import Downloader
try:
    from parsers.twitter_parser import TwitterParser
except ImportError:
    from ..parsers.twitter_parser import TwitterParser


class WarcImporter:
    """手動WARCファイルのオフライン自動監査・逆引きパース・メディア抽出 (SPEC-PLUGIN-001)"""
    def __init__(self, warc_path: str, db_path: str = "archive.db", storage_dir: Optional[str] = None, offline: bool = True):
        self.warc_path = warc_path
        self.offline = offline
        self.mutator = Mutator(db_path=db_path)
        self.downloader = Downloader(db_path=db_path, storage_dir=storage_dir)
        self.parser = TwitterParser()

    def audit_warc(self) -> Dict[str, Any]:
        """WARCコンテナ内の通信レコードからプラットフォームとアカウント名を自動監査逆引き"""
        if not os.path.exists(self.warc_path):
            return {"platform": "unknown", "account": "", "records": 0}
        detected = {"platform": "twitter", "account": "", "records": 0, "accounts": set()}
        with open(self.warc_path, "rb") as s:
            for r in ArchiveIterator(s):
                detected["records"] += 1
                uri = r.rec_headers.get_header("WARC-Target-URI") or ""
                det = self.parser.detect_platform_and_account(uri)
                if det and det.get("account"):
                    detected["accounts"].add(det["account"])
                    if not detected["account"]: detected["account"] = det["account"]
        detected["accounts"] = list(detected["accounts"])
        return detected

    def run_import(self, progress_callback: Optional[Callable[[int, int, str], None]] = None) -> int:
        if not os.path.exists(self.warc_path):
            if progress_callback: progress_callback(1, 1, f"WARC not found: {self.warc_path}")
            return 0

        audit = self.audit_warc()
        account = audit.get("account") or "unknown"
        tot = max(audit.get("records", 1), 1)
        if progress_callback:
            progress_callback(0, tot, f"Audited WARC: platform={audit['platform']}, account=@{account} (records: {tot})")

        saved, extracted = 0, 0
        with open(self.warc_path, "rb") as s:
            for idx, r in enumerate(ArchiveIterator(s), start=1):
                if r.rec_type != "response": continue
                uri = r.rec_headers.get_header("WARC-Target-URI") or ""
                c_type = (r.http_headers.get_header("Content-Type") if r.http_headers else "") or ""

                if "json" in c_type or "html" in c_type or "status" in uri:
                    raw = r.raw_stream.read()
                    parsed = self.parser.parse_record(raw, uri)
                    if parsed and self.mutator.upsert_record(parsed):
                        saved += 1
                        u = parsed.get("account", {}).get("username") or account
                        if u and u != "unknown": account = u
                        if progress_callback:
                            progress_callback(idx, tot, f"Saved post {parsed['post']['id']} (@{u})")
                elif any(d in uri for d in ["twimg.com", "video.twimg.com"]):
                    media_id = os.path.basename(uri.split("?")[0])
                    m_type = "video" if any(media_id.endswith(ext) for ext in [".mp4", ".webm", ".m3u8"]) else "image"
                    dest = self.downloader._get_target_path(account, media_id, m_type)
                    with open(dest, "wb") as mf:
                        mf.write(r.raw_stream.read())
                    img_id = self.downloader.stash.register_media(dest, "image") if m_type == "image" else None
                    scn_id = self.downloader.stash.register_media(dest, "video") if m_type != "image" else None
                    status = "COMPLETED" if (img_id or scn_id or self.offline) else "RETAINED"
                    reason = None if (img_id or scn_id or self.offline) else "Extracted to disk, awaiting Stash index"
                    self.downloader._update_status(media_id, status, reason, img_id, scn_id)
                    extracted += 1
                    if progress_callback:
                        progress_callback(idx, tot, f"Extracted media {media_id} -> {status}")

        if progress_callback:
            progress_callback(tot, tot, f"Offline import completed: {saved} posts, {extracted} media extracted.")
        return saved
