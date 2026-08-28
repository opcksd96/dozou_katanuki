# plugins/twitter/scraper/core/warc_archiver.py (100行以下 - SPEC-PRINCIPLE-001)
import json, os
from typing import Any, Dict, List


class WarcArchiver:
    """各投稿の原本ダング (backups/dumps/{platform}/{user}/{id}/) への WARC & メタデータ充足エンジン"""
    def __init__(self, base_dir: str = "backups/dumps"):
        self.base_dir = base_dir

    def build_standalone_html(self, post: Dict[str, Any]) -> str:
        p, a = post.get("post", {}), post.get("account", {})
        post_id, handle = p.get("id", "0"), a.get("username", "unknown")
        text = p.get("full_text", "")
        media_html = "".join([f'<img src="{m.get("url")}" style="max-width:100%;margin-top:8px;border-radius:8px;" />' for m in post.get("media", [])])
        return f"""<!DOCTYPE html>
<html lang="ja"><head><meta charset="utf-8"><title>Tweet by @{handle} ({post_id})</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>body{{font-family:sans-serif;background:#0f172a;color:#f8fafc;padding:20px;max-width:600px;margin:auto;}}
.card{{background:#1e293b;border:1px solid #334155;border-radius:12px;padding:16px;}}
.meta{{color:#94a3b8;font-size:12px;margin-bottom:8px;}}
.text{{white-space:pre-wrap;line-height:1.5;}}</style></head>
<body><div class="card"><div class="meta"><strong>@{handle}</strong> · {p.get("created_at", "")}</div>
<div class="text">{text}</div>{media_html}
<div class="meta" style="margin-top:12px;">Archived Tweet ID: {post_id}</div></div></body></html>"""

    def _merge_metadata(self, dest_path: str, new_data: Dict[str, Any]) -> None:
        merged = new_data
        if os.path.exists(dest_path):
            try:
                with open(dest_path, "r", encoding="utf-8") as f:
                    old_data = json.load(f)
                if isinstance(old_data, dict):
                    # 既存のメディアや不足項目をマージ
                    old_media = {m.get("url"): m for m in old_data.get("media", []) if m.get("url")}
                    for m in new_data.get("media", []):
                        if m.get("url"): old_media[m["url"]] = m
                    merged = dict(old_data)
                    merged.update(new_data)
                    merged["media"] = list(old_media.values())
            except Exception: pass
        with open(dest_path, "w", encoding="utf-8") as f:
            json.dump(merged, f, ensure_ascii=False, indent=2)

    def archive_posts(self, posts: List[Dict[str, Any]], platform: str = "twitter") -> int:
        if not posts: return 0
        warcwriter_cls, status_cls = None, None
        try:
            from warcio.warcwriter import WARCWriter
            from warcio.statusandheaders import StatusAndHeaders
            warcwriter_cls, status_cls = WARCWriter, StatusAndHeaders
        except ImportError: pass

        archived = 0
        for post in posts:
            p, a = post.get("post", {}), post.get("account", {})
            post_id, handle = str(p.get("id", "")), a.get("username", "unknown")
            if not post_id or not handle: continue
            
            post_dir = os.path.join(self.base_dir, platform, handle, post_id)
            os.makedirs(post_dir, exist_ok=True)
            self._merge_metadata(os.path.join(post_dir, "metadata.json"), post)

            if warcwriter_cls and status_cls:
                warc_path = os.path.join(post_dir, "snapshot.warc.gz")
                canonical_url = f"https://twitter.com/{handle}/status/{post_id}"
                html_bytes = self.build_standalone_html(post).encode("utf-8")
                json_bytes = json.dumps(post, ensure_ascii=False).encode("utf-8")
                
                with open(warc_path, "wb") as output:
                    writer = warcwriter_cls(output, gzip=True)
                    rec_html = writer.create_warc_record(
                        uri=canonical_url, record_type="response",
                        http_headers=status_cls("200 OK", [("Content-Type", "text/html; charset=utf-8"), ("Content-Length", str(len(html_bytes))), ("X-Archive-Source", "sotwe")], protocol="HTTP/1.1")
                    )
                    rec_html.raw_stream.write(html_bytes)
                    writer.write_record(rec_html)
                    rec_meta = writer.create_warc_record(
                        uri=f"{canonical_url}#metadata", record_type="metadata",
                        http_headers=status_cls("200 OK", [("Content-Type", "application/json; charset=utf-8"), ("Content-Length", str(len(json_bytes)))], protocol="HTTP/1.1")
                    )
                    rec_meta.raw_stream.write(json_bytes)
                    writer.write_record(rec_meta)
            archived += 1
        return archived
