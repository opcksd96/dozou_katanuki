# plugins/base/scraper/core/base_mutator.py (SPEC-PLUGIN-001 / 100行以下)
import sqlite3, time
from typing import Any, Dict, List, Optional, Tuple
from .translator import Translator


class BaseMutator:
    """クリティカルセクション極小化・アカウント名寄せ(Reconciliation)・一括コミットミューテーター"""
    def __init__(self, db_path: str = "archive.db", platform: str = "base", translator: Optional[Translator] = None, enable_translation: bool = True):
        self.db_path, self.platform, self.enable_translation = db_path, platform, enable_translation
        self.translator = translator or (Translator() if enable_translation else None)

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=60.0)
        conn.execute("PRAGMA journal_mode = WAL;"); conn.execute("PRAGMA synchronous = NORMAL;")
        conn.execute("PRAGMA foreign_keys = ON;"); conn.execute("PRAGMA busy_timeout = 60000;")
        return conn

    def _prepare_in_memory(self, records: List[Dict[str, Any]]) -> Tuple[List[tuple], List[tuple], List[tuple], List[tuple], List[tuple], int]:
        now_ts = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())
        accounts, hists, redirects, articles, media = [], [], [], [], []
        valid_count = 0
        for data in records:
            acc, post, m_list = data.get("account", {}), data.get("post", {}), data.get("media", [])
            if not post.get("id") or not acc.get("username"): continue
            valid_count += 1
            acc_id, u_name = str(acc.get("numeric_id") or acc.get("username")), acc["username"]
            d_name, av_url = acc.get("display_name", u_name), acc.get("avatar_url", "")
            accounts.append((acc_id, u_name, d_name, av_url, now_ts))
            hists.append((acc_id, d_name, av_url, 1, f"{u_name}_avatar_001", now_ts))
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
            articles.append((post_id, acc_id, str(post.get("conversation_id") or post_id), post.get("reply_to_id") or post.get("reply_to_tweet_id"),
                             post.get("reply_to_handle"), c_at, ftext, lang, ja, en, zh, post.get("via") or self.platform, 0, 0, post.get("wayback_url", "")))
            for m in m_list:
                m_url = m.get("download_url") or m.get("url")
                if not m_url: continue
                fn = m_url.split("?")[0].split("/")[-1]
                ext = fn.split(".")[-1].lower() if "." in fn else ("jpg" if "format=jpg" in m_url or m.get("type") == "image" else "mp4")
                m_id = m.get("media_id") or (fn if "." in fn else f"{fn}.{ext}")
                media.append((m_id, post_id, m.get("type", "image"), m_url, m.get("width", 0), m.get("height", 0)))
        return accounts, hists, redirects, articles, media, valid_count

    def _reconcile_temp_accounts(self, cur: sqlite3.Cursor, accs: List[tuple]) -> None:
        """仮アカウント(UUID等)から真の数値IDへの自動名寄せ・リレーション付け替え"""
        has_b64 = False
        try:
            cols = [c[1] for c in cur.execute("PRAGMA table_info(accounts)").fetchall()]
            has_b64 = "avatar_base64" in cols
        except Exception: pass

        for acc_id, u_name, _, _, _ in accs:
            if acc_id.isdigit():
                temps = [r[0] for r in cur.execute("SELECT numeric_id FROM accounts WHERE lower(username) = lower(?) AND numeric_id != ?", (u_name, acc_id)).fetchall()]
                for t_id in temps:
                    cur.execute("UPDATE articles SET account_id = ? WHERE account_id = ?", (acc_id, t_id))
                    if cur.execute("SELECT 1 FROM sqlite_master WHERE type='table' AND name='account_profile_histories'").fetchone():
                        cur.execute("UPDATE account_profile_histories SET account_id = ? WHERE account_id = ?", (acc_id, t_id))
                    if has_b64:
                        b64 = cur.execute("SELECT avatar_base64 FROM accounts WHERE numeric_id = ?", (t_id,)).fetchone()
                        if b64 and b64[0]: cur.execute("UPDATE accounts SET avatar_base64 = coalesce(avatar_base64, ?) WHERE numeric_id = ?", (b64[0], acc_id))
                    cur.execute("DELETE FROM accounts WHERE numeric_id = ?", (t_id,))

    def upsert_batch(self, records: List[Dict[str, Any]]) -> int:
        if not records: return 0
        accs, hists, redirs, arts, meds, valid_count = self._prepare_in_memory(records)
        if valid_count == 0: return 0
        with self._get_conn() as conn:
            cur = conn.cursor()
            self._reconcile_temp_accounts(cur, accs)
            cur.executemany("""
                INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at) VALUES (?, ?, ?, ?, ?)
                ON CONFLICT(numeric_id) DO UPDATE SET display_name=coalesce(excluded.display_name, accounts.display_name),
                    avatar_url=coalesce(excluded.avatar_url, accounts.avatar_url), updated_at=excluded.updated_at
            """, accs)
            if cur.execute("SELECT 1 FROM sqlite_master WHERE type='table' AND name='account_profile_histories'").fetchone():
                cur.executemany("""
                    INSERT INTO account_profile_histories (account_id, display_name, avatar_original_url, avatar_seq, avatar_virtual_key, observed_at)
                    SELECT ?, ?, ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM account_profile_histories WHERE account_id = ?)
                """, [(h[0], h[1], h[2], h[3], h[4], h[5], h[0]) for h in hists])
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
            cur.executemany("INSERT OR REPLACE INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)", redirs)
            cur.executemany("INSERT OR IGNORE INTO media (media_id, article_id, type, download_url, width, height, download_status) VALUES (?, ?, ?, ?, ?, ?, 'QUEUED')", meds)
            conn.commit()
        return valid_count

    def upsert_record(self, data: Dict[str, Any]) -> bool: return self.upsert_batch([data]) == 1
