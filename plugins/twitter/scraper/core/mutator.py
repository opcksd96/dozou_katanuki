# plugins/twitter/scraper/core/mutator.py (100行以下)
import os, sqlite3, time
from typing import Any, Dict, List, Optional
from .translator import Translator


class Mutator:
    """共通正規化辞書の SQLite3 (archive.db) Upsert エンジン (SPEC-PLUGIN-001)"""

    def __init__(self, db_path: str = "archive.db", translator: Optional[Translator] = None, enable_translation: bool = True):
        self.db_path = db_path
        self.enable_translation = enable_translation
        self.translator = translator or (Translator() if enable_translation else None)

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=10.0)
        conn.execute("PRAGMA journal_mode = WAL;")
        conn.execute("PRAGMA foreign_keys = ON;")
        conn.execute("PRAGMA busy_timeout = 5000;")
        return conn

    def upsert_record(self, data: Dict[str, Any]) -> bool:
        account, post, media_list = data.get("account", {}), data.get("post", {}), data.get("media", [])
        if not post.get("id") or not account.get("username"): return False
        account_id = str(account.get("numeric_id") or account.get("username"))
        now_ts = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())

        with self._get_conn() as conn:
            cur = conn.cursor()
            cur.execute("""
                INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at)
                VALUES (?, ?, ?, ?, ?)
                ON CONFLICT(numeric_id) DO UPDATE SET
                    display_name = coalesce(excluded.display_name, accounts.display_name),
                    avatar_url = coalesce(excluded.avatar_url, accounts.avatar_url),
                    updated_at = excluded.updated_at
            """, (account_id, account["username"], account.get("display_name", account["username"]),
                  account.get("avatar_url", ""), now_ts))

            post_id, created_at, full_text = str(post["id"]), post.get("created_at") or now_ts, post.get("full_text", "")
            for u in post.get("urls", []):
                s_url, e_url = u.get("short_url"), u.get("expanded_url")
                if s_url and e_url:
                    cur.execute("INSERT OR REPLACE INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)",
                                (s_url, e_url, post_id))
                    if s_url in full_text: full_text = full_text.replace(s_url, e_url)

            # 多言語翻訳キャッシュ生成 (SPEC-MIDDLEWARE-001-2)
            lang, ja, en, zh = "ja", full_text, None, None
            if self.translator:
                trans = self.translator.translate_article(full_text)
                lang = trans.get("lang", "ja")
                ja, en, zh = trans.get("ja"), trans.get("en"), trans.get("zh")

            cur.execute("""
                INSERT INTO articles (
                    id, account_id, conversation_id, reply_to_id, reply_to_handle,
                    created_at, full_text, lang, full_text_ja, full_text_en, full_text_zh,
                    via, is_repost, is_liked, wayback_url
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'twitter', 0, 0, ?)
                ON CONFLICT(id) DO UPDATE SET
                    full_text = excluded.full_text,
                    lang = excluded.lang,
                    full_text_ja = coalesce(excluded.full_text_ja, articles.full_text_ja),
                    full_text_en = coalesce(excluded.full_text_en, articles.full_text_en),
                    full_text_zh = coalesce(excluded.full_text_zh, articles.full_text_zh),
                    wayback_url = coalesce(excluded.wayback_url, articles.wayback_url)
            """, (post_id, account_id, str(post.get("conversation_id") or post_id),
                  post.get("reply_to_tweet_id"), post.get("reply_to_handle"),
                  created_at, full_text, lang, ja, en, zh, post.get("wayback_url", "")))

            for m in media_list:
                m_url = m.get("url")
                if not m_url: continue
                media_id = m.get("media_id") or os.path.basename(m_url.split("?")[0])
                cur.execute("""
                    INSERT OR IGNORE INTO media (
                        media_id, article_id, type, download_url, width, height,
                        download_status, failed_reason, stash_scene_id, stash_image_id
                    ) VALUES (?, ?, ?, ?, ?, ?, 'QUEUED', NULL, NULL, NULL)
                """, (media_id, post_id, m.get("type", "image"), m_url, m.get("width", 0), m.get("height", 0)))

            conn.commit()
            return True
