# plugins/base/scraper/core/base_mutator.py (SPEC-PLUGIN-001 / 100行以下)
import sqlite3, time, uuid
from typing import Any, Dict, List, Optional, Tuple
from .translator import Translator

DEFAULT_AVATAR_SVG = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>"

class BaseMutator:
    """クリティカルセクション極小化・アカウント名寄せ・一括コミットミューテーター"""
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
            raw_id = str(acc.get("numeric_id") or "").strip()
            u_name = str(acc["username"]).strip().lower()
            acc_id = raw_id if (raw_id and (raw_id.isdigit() or (len(raw_id) == 36 and "-" in raw_id))) else str(uuid.uuid5(uuid.NAMESPACE_DNS, f"{self.platform}_{u_name}"))
            d_name = acc.get("display_name") or acc["username"]
            desc, av_url, av_b64 = acc.get("description", ""), acc.get("avatar_url", ""), acc.get("avatar_base64") or DEFAULT_AVATAR_SVG
            av_orig, post_id, c_at = acc.get("avatar_original_url") or av_url, str(post["id"]), post.get("created_at") or now_ts
            accounts.append((acc_id, u_name, d_name, av_url, desc, av_b64, now_ts))
            if av_orig: hists.append((acc_id, d_name, desc, av_orig, 1, f"{u_name}_avatar_001", av_b64, c_at))
            ftext = post.get("full_text", "")
            for u in post.get("urls", []):
                s_u, e_u = u.get("short_url"), u.get("expanded_url")
                if s_u and e_u: redirects.append((s_u, e_u, post_id)); ftext = ftext.replace(s_u, e_u) if s_u in ftext else ftext
            lang, ja, en, zh = post.get("lang") or "ja", post.get("full_text_ja"), post.get("full_text_en"), post.get("full_text_zh")
            if self.translator and ja is None:
                t = self.translator.translate_article(ftext)
                lang, ja, en, zh = t.get("lang", lang), t.get("ja"), t.get("en"), t.get("zh")
            articles.append((post_id, acc_id, str(post.get("conversation_id") or post_id), post.get("reply_to_id") or post.get("reply_to_tweet_id"),
                             post.get("reply_to_handle"), c_at, ftext, lang, ja or ftext, en, zh, post.get("via") or self.platform, 0, 0, post.get("wayback_url", "")))
            for m in m_list:
                m_url = m.get("download_url") or m.get("url")
                fn = m_url.split("?")[0].split("/")[-1]
                for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: fn = fn[:-len(sfx)] if fn.endswith(sfx) else fn
                ext = fn.split(".")[-1].lower() if "." in fn else ("jpg" if "format=jpg" in m_url or m.get("type") == "image" else "mp4")
                clean_mid = (m.get("media_id") or fn).split("?")[0]
                for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: clean_mid = clean_mid[:-len(sfx)] if clean_mid.endswith(sfx) else clean_mid
                if "." not in clean_mid: clean_mid = f"{clean_mid}.{ext}"
                media.append((clean_mid, post_id, m.get("type", "image"), m_url, m.get("width", 0), m.get("height", 0)))
            valid_count += 1
        return accounts, hists, redirects, articles, media, valid_count

    def upsert_batch(self, records: List[Dict[str, Any]]) -> int:
        if not records: return 0
        accs, hists, redirs, arts, meds, valid_count = self._prepare_in_memory(records)
        if valid_count == 0: return 0
        with self._get_conn() as conn:
            cur = conn.cursor()
            cur.executemany("""
                INSERT INTO accounts (numeric_id, username, display_name, avatar_url, description, avatar_base64, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?)
                ON CONFLICT(numeric_id) DO UPDATE SET display_name=coalesce(excluded.display_name, accounts.display_name),
                    avatar_url=coalesce(excluded.avatar_url, accounts.avatar_url),
                    description=case when excluded.description != '' then excluded.description else accounts.description end,
                    avatar_base64=coalesce(accounts.avatar_base64, excluded.avatar_base64), updated_at=excluded.updated_at
            """, accs)
            cur.executemany("""
                INSERT INTO account_profile_histories (account_id, display_name, description, avatar_original_url, avatar_seq, avatar_virtual_key, avatar_base64, observed_at)
                SELECT ?, ?, ?, ?, ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM account_profile_histories WHERE account_id = ?)
            """, [(h[0], h[1], h[2], h[3], h[4], h[5], h[6], h[7], h[0]) for h in hists])
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
