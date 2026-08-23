# plugins/twitter/scraper/core/mutator.py (100行以下)
import sqlite3, time
from typing import Any, Dict, List, Optional, Tuple
from .translator import Translator


class Mutator:
    """クリティカルセクション極小化・純粋インメモリ構築＆executemany一括コミットエンジン (SPEC-PLUGIN-001)"""
    def __init__(self, db_path: str = "archive.db", translator: Optional[Translator] = None, enable_translation: bool = True):
        self.db_path = db_path
        self.enable_translation = enable_translation
        self.translator = translator or (Translator() if enable_translation else None)

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=60.0)
        conn.execute("PRAGMA journal_mode = WAL;"); conn.execute("PRAGMA synchronous = NORMAL;")
        conn.execute("PRAGMA foreign_keys = ON;"); conn.execute("PRAGMA busy_timeout = 60000;")
        return conn

    def _prepare_in_memory(self, records: List[Dict[str, Any]]) -> Tuple[List[tuple], List[tuple], List[tuple], List[tuple], List[tuple], int]:
        now_ts = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())
        accounts, whitelists, redirects, articles, media = [], [], [], [], []
        valid_count = 0

        for data in records:
            acc, post, m_list = data.get("account", {}), data.get("post", {}), data.get("media", [])
            if not post.get("id") or not acc.get("username"): continue
            valid_count += 1
            acc_id, u_name = str(acc.get("numeric_id") or acc.get("username")), acc["username"]

            accounts.append((acc_id, u_name, acc.get("display_name", u_name), acc.get("avatar_url", ""), now_ts))
            whitelists.append((u_name,))

            post_id, c_at, ftext = str(post["id"]), post.get("created_at") or now_ts, post.get("full_text", "")
            for u in post.get("urls", []):
                s_u, e_u = u.get("short_url"), u.get("expanded_url")
                if s_u and e_u:
                    redirects.append((s_u, e_u, post_id))
                    if s_u in ftext: ftext = ftext.replace(s_u, e_u)

            lang, ja, en, zh = post.get("lang") or "ja", post.get("full_text_ja"), post.get("full_text_en"), post.get("full_text_zh")
            if self.translator and ja is None:
                t = self.translator.translate_article(ftext)
                lang, ja, en, zh = t.get("lang", lang), t.get("ja"), t.get("en"), t.get("zh")
            if ja is None: ja = ftext

            articles.append((post_id, acc_id, str(post.get("conversation_id") or post_id), post.get("reply_to_tweet_id"),
                             post.get("reply_to_handle"), c_at, ftext, lang, ja, en, zh, "twitter", 0, 0, post.get("wayback_url", "")))

            for m in m_list:
                m_url = m.get("url")
                if not m_url: continue
                fn = m_url.split("?")[0].split("/")[-1]
                ext = fn.split(".")[-1].lower() if "." in fn else ("jpg" if "format=jpg" in m_url or m.get("type") == "image" else "mp4")
                m_id = m.get("media_id") or (fn if "." in fn else f"{fn}.{ext}")
                media.append((m_id, post_id, m.get("type", "image"), m_url, m.get("width", 0), m.get("height", 0)))

        return accounts, whitelists, redirects, articles, media, valid_count

    def upsert_batch(self, records: List[Dict[str, Any]]) -> int:
        if not records: return 0
        accs, wls, redirs, arts, meds, valid_count = self._prepare_in_memory(records)
        if valid_count == 0: return 0

        # クリティカルセクション（純粋な一括SQL投入のみ・安易なtry-except不使用）
        with self._get_conn() as conn:
            cur = conn.cursor()
            cur.executemany("""
                INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at) VALUES (?, ?, ?, ?, ?)
                ON CONFLICT(numeric_id) DO UPDATE SET display_name=coalesce(excluded.display_name, accounts.display_name),
                    avatar_url=coalesce(excluded.avatar_url, accounts.avatar_url), updated_at=excluded.updated_at
            """, accs)
            cur.executemany("INSERT INTO whitelists (type, value, is_active) VALUES ('account', ?, 1) ON CONFLICT(value) DO UPDATE SET is_active=1", wls)
            cur.executemany("INSERT OR REPLACE INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)", redirs)
            cur.executemany("""
                INSERT INTO articles (id, account_id, conversation_id, reply_to_id, reply_to_handle, created_at, full_text,
                    lang, full_text_ja, full_text_en, full_text_zh, via, is_repost, is_liked, wayback_url)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                ON CONFLICT(id) DO UPDATE SET full_text=excluded.full_text, lang=excluded.lang,
                    full_text_ja=coalesce(excluded.full_text_ja, articles.full_text_ja),
                    full_text_en=coalesce(excluded.full_text_en, articles.full_text_en),
                    full_text_zh=coalesce(excluded.full_text_zh, articles.full_text_zh),
                    wayback_url=coalesce(excluded.wayback_url, articles.wayback_url)
            """, arts)
            cur.executemany("INSERT OR IGNORE INTO media (media_id, article_id, type, download_url, width, height, download_status) VALUES (?, ?, ?, ?, ?, ?, 'QUEUED')", meds)
            conn.commit()
        return valid_count

    def upsert_record(self, data: Dict[str, Any]) -> bool:
        return self.upsert_batch([data]) == 1
