# plugins/twitter/scraper/core/mutator.py (100行以下)
import os
import sqlite3
import time
from typing import Any, Dict, List, Optional


class Mutator:
    """共通正規化辞書の SQLite3 (archive.db) Upsert エンジン (SPEC-PLUGIN-001)"""

    def __init__(self, db_path: str = "archive.db"):
        self.db_path = db_path

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=10.0)
        conn.execute("PRAGMA journal_mode = WAL;")
        conn.execute("PRAGMA foreign_keys = ON;")
        conn.execute("PRAGMA busy_timeout = 5000;")
        return conn

    def upsert_record(self, data: Dict[str, Any]) -> bool:
        """アカウント・記事・メディアのリレーショナル Upsert トランザクション"""
        account = data.get("account", {})
        post = data.get("post", {})
        media_list = data.get("media", [])

        if not post.get("id") or not account.get("username"):
            return False

        account_id = str(account.get("numeric_id") or account.get("username"))
        now_ts = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())

        with self._get_conn() as conn:
            cur = conn.cursor()
            # 1. accounts Upsert
            cur.execute("""
                INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at)
                VALUES (?, ?, ?, ?, ?)
                ON CONFLICT(numeric_id) DO UPDATE SET
                    display_name = coalesce(excluded.display_name, accounts.display_name),
                    avatar_url = coalesce(excluded.avatar_url, accounts.avatar_url),
                    updated_at = excluded.updated_at
            """, (account_id, account["username"], account.get("display_name", account["username"]),
                  account.get("avatar_url", ""), now_ts))

            # 2. articles Upsert
            post_id = str(post["id"])
            created_at = post.get("created_at") or now_ts
            full_text = post.get("full_text", "")
            urls = post.get("urls", [])

            # 短縮URLの保存 & 本文置換 (SPEC-MIDDLEWARE-001-2)
            for u in urls:
                s_url = u.get("short_url")
                e_url = u.get("expanded_url")
                if s_url and e_url:
                    cur.execute("""
                        INSERT OR REPLACE INTO url_redirects (short_url, expanded_url, article_id)
                        VALUES (?, ?, ?)
                    """, (s_url, e_url, post_id))
                    if s_url in full_text:
                        full_text = full_text.replace(s_url, e_url)

            cur.execute("""
                INSERT INTO articles (
                    id, account_id, conversation_id, reply_to_id, reply_to_handle,
                    created_at, full_text, lang, full_text_ja, full_text_en, full_text_zh,
                    via, is_repost, is_liked, wayback_url
                ) VALUES (?, ?, ?, ?, ?, ?, ?, 'ja', ?, ?, ?, 'twitter', 0, 0, ?)
                ON CONFLICT(id) DO UPDATE SET
                    full_text = excluded.full_text,
                    wayback_url = coalesce(excluded.wayback_url, articles.wayback_url)
            """, (post_id, account_id, str(post.get("conversation_id") or post_id),
                  post.get("reply_to_tweet_id"), post.get("reply_to_handle"),
                  created_at, full_text, full_text, full_text, full_text,
                  post.get("wayback_url", "")))

            # 3. media Insert (初期状態 QUEUED)
            for m in media_list:
                m_url = m.get("url")
                if not m_url:
                    continue
                media_id = m.get("media_id") or os.path.basename(m_url.split("?")[0])
                cur.execute("""
                    INSERT OR IGNORE INTO media (
                        media_id, article_id, type, download_url, width, height,
                        download_status, failed_reason, stash_scene_id, stash_image_id
                    ) VALUES (?, ?, ?, ?, ?, ?, 'QUEUED', NULL, NULL, NULL)
                """, (media_id, post_id, m.get("type", "image"), m_url,
                      m.get("width", 0), m.get("height", 0)))

            conn.commit()
            return True
