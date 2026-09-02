# plugins/base/scraper/core/base_mutator.py (SPEC-PLUGIN-001 / 100行以下)
import json, sqlite3, time, uuid
from typing import Any, Dict, List, Optional, Tuple
from .translator import Translator

DEFAULT_AVATAR_SVG = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>"

class BaseMutator:
    """クリティカルセクション極小化・マルチソース統合・Tri-Node対等リレーションミューテーター"""
    def __init__(self, db_path: str = "archive.db", platform: str = "base", translator: Optional[Translator] = None, enable_translation: bool = True):
        self.db_path, self.platform, self.enable_translation = db_path, platform, enable_translation
        self.translator, self._migrated = translator or (Translator() if enable_translation else None), False

    def _get_conn(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self.db_path, timeout=60.0)
        conn.execute("PRAGMA journal_mode = WAL;"); conn.execute("PRAGMA synchronous = NORMAL;")
        conn.execute("PRAGMA foreign_keys = ON;"); conn.execute("PRAGMA busy_timeout = 60000;")
        if not self._migrated: self._ensure_schema(conn); self._migrated = True
        return conn

    def _ensure_schema(self, conn: sqlite3.Connection) -> None:
        cur = conn.cursor()
        cols = {row[1] for row in cur.execute("PRAGMA table_info(articles)").fetchall()}
        for col, ctype in [("source_name", "TEXT"), ("source_domain", "TEXT"), ("original_url", "TEXT"), ("sotwe_url", "TEXT"), ("nitter_url", "TEXT"), ("twistalker_url", "TEXT"), ("is_trash", "BOOLEAN DEFAULT 0"), ("trashed_by", "TEXT"), ("trash_reason", "TEXT"), ("trashed_at", "DATETIME")]:
            if col not in cols: cur.execute(f"ALTER TABLE articles ADD COLUMN {col} {ctype}")
        acc_cols = {row[1] for row in cur.execute("PRAGMA table_info(accounts)").fetchall()}
        for col, ctype in [("description", "TEXT"), ("group_name", "TEXT DEFAULT ''"), ("alias_of", "TEXT DEFAULT ''"), ("avatar_base64", "TEXT")]:
            if col not in acc_cols: cur.execute(f"ALTER TABLE accounts ADD COLUMN {col} {ctype}")
        med_cols = {row[1] for row in cur.execute("PRAGMA table_info(media)").fetchall()}
        for col, ctype in [("thumbnail_url", "TEXT"), ("tweet_urls", "TEXT"), ("media_quality", "TEXT"), ("account_id", "TEXT")]:
            if col not in med_cols: cur.execute(f"ALTER TABLE media ADD COLUMN {col} {ctype}")
        cur.execute("CREATE TABLE IF NOT EXISTS media_excluded (media_id TEXT PRIMARY KEY, article_id TEXT, account_id TEXT, type TEXT, download_url TEXT, width INTEGER, height INTEGER, download_status TEXT, failed_reason TEXT, tweet_urls TEXT, thumbnail_url TEXT, quarantined_at DATETIME DEFAULT CURRENT_TIMESTAMP, quarantine_reason TEXT)")
        cur.execute("CREATE TABLE IF NOT EXISTS media_variants (variant_hash TEXT PRIMARY KEY, media_id TEXT, article_id TEXT, download_url TEXT, bit_rate INTEGER)")
        conn.commit()

    def _get_active_whitelist(self, conn: sqlite3.Connection) -> set:
        try: return {r[0].lower() for r in conn.cursor().execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall() if r[0]}
        except Exception: return set()

    def _prepare_in_memory(self, records: List[Dict[str, Any]], wl: set = None) -> Tuple[List[tuple], List[tuple], List[tuple], List[tuple], List[tuple], List[tuple], List[tuple], int]:
        now_ts, accounts, hists, redirects, articles, media, excluded, valid_count = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime()), [], [], [], [], [], [], 0
        media_variants = []
        wl_set = wl or set()
        for data in records:
            acc, post, m_list = data.get("account", {}), data.get("post", {}), data.get("media", [])
            if not post.get("id") or not acc.get("username"): continue
            raw_id, u_name = str(acc.get("numeric_id") or "").strip(), str(acc["username"]).strip().lower()
            acc_id = raw_id if (raw_id and (raw_id.isdigit() or (len(raw_id) == 36 and "-" in raw_id))) else str(uuid.uuid5(uuid.NAMESPACE_DNS, f"{self.platform}_{u_name}"))
            d_name, desc = acc.get("display_name") or acc["username"], acc.get("description", "")
            av_url, av_b64, post_id, c_at = acc.get("avatar_url", ""), acc.get("avatar_base64") or DEFAULT_AVATAR_SVG, str(post["id"]), post.get("created_at") or now_ts
            accounts.append((acc_id, u_name, d_name, av_url, desc, av_b64, now_ts)); (acc.get("avatar_original_url") and hists.append((acc_id, d_name, desc, acc.get("avatar_original_url") or av_url, 1, f"{u_name}_avatar_001", av_b64, c_at)))
            ftext = post.get("full_text", "")
            for u in post.get("urls", []):
                s_u, e_u = u.get("short_url"), u.get("expanded_url")
                if s_u and e_u: (redirects.append((s_u, e_u, post_id)), (ftext := ftext.replace(s_u, e_u) if s_u in ftext else ftext))
            lang, ja, en, zh = post.get("lang") or "ja", post.get("full_text_ja"), post.get("full_text_en"), post.get("full_text_zh")
            if self.translator and ja is None:
                t = self.translator.translate_article(ftext); lang, ja, en, zh = t.get("lang", lang), t.get("ja"), t.get("en"), t.get("zh")
            via, sname = post.get("via") or self.platform, post.get("source_name") or (post.get("via") or self.platform).lower()
            sdom, raw_wb = post.get("source_domain") or ("sotwe.com" if "sotwe" in sname else ("archive.org" if "wayback" in sname else "x.com")), post.get("wayback_url", "")
            wb_u, orig_url = (raw_wb if "web.archive.org" in raw_wb else ""), (post.get("original_url") or (raw_wb if ("x.com" in raw_wb or "twitter.com" in raw_wb) else f"https://x.com/{u_name}/status/{post_id}"))
            sotwe_u = post.get("sotwe_url") or (f"https://www.sotwe.com/tweet/{post_id}" if "sotwe" in sname else None)
            articles.append((post_id, acc_id, str(post.get("conversation_id") or post_id), post.get("reply_to_id") or post.get("reply_to_tweet_id"), post.get("reply_to_handle"), c_at, ftext, lang, ja or ftext, en, zh, via, 0, 0, wb_u, sname, sdom, orig_url, sotwe_u, post.get("nitter_url"), post.get("twistalker_url")))
            t_urls_json = json.dumps([u for u in [orig_url, sotwe_u, wb_u] if u], ensure_ascii=False)
            is_wl = (not wl_set) or (u_name in wl_set)
            for m in m_list:
                m_url = m.get("download_url") or m.get("url") or ""
                if not m_url: continue
                fn = m_url.split("?")[0].split("/")[-1]
                for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: fn = fn[:-len(sfx)] if fn.endswith(sfx) else fn
                ext = fn.split(".")[-1].lower() if "." in fn else ("jpg" if "format=jpg" in m_url or m.get("type") == "image" else "mp4")
                clean_mid = fn if ("." in fn and not fn.split(".")[0].isdigit()) else (m.get("media_id") or fn).split("?")[0]
                for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: clean_mid = clean_mid[:-len(sfx)] if clean_mid.endswith(sfx) else clean_mid
                if "." not in clean_mid: clean_mid = f"{clean_mid}.{ext}"
                if clean_mid.split(".")[0].isdigit() and len(clean_mid.split(".")[0]) >= 15 and "." in fn and not fn.split(".")[0].isdigit(): clean_mid = fn
                if is_wl:
                    media.append((clean_mid, post_id, acc_id, m.get("type", "image"), m_url, m.get("width", 0), m.get("height", 0), m.get("thumbnail_url"), t_urls_json, m.get("quality", ""), "QUEUED", None))
                    for var in m.get("variants", []):
                        v_url = var.get("url") or ""
                        if not v_url: continue
                        v_fn = v_url.split("?")[0].split("/")[-1]
                        for sfx in [":large", ":orig", ":small", ":medium", ":thumb"]: v_fn = v_fn[:-len(sfx)] if v_fn.endswith(sfx) else v_fn
                        v_hash = v_fn.split(".")[0]
                        media_variants.append((v_hash, clean_mid, post_id, v_url, var.get("bit_rate", 0)))
                else:
                    excluded.append((clean_mid, post_id, acc_id, m.get("type", "image"), m_url, m.get("width", 0), m.get("height", 0), "EXCLUDED", "Whitelist外", t_urls_json, m.get("thumbnail_url"), f"Non-whitelist account: {u_name}"))
            valid_count += 1
        return accounts, hists, redirects, articles, media, media_variants, excluded, valid_count

    def upsert_batch(self, records: List[Dict[str, Any]]) -> Tuple[int, int, int]:
        if not records: return (0, 0, 0)
        with self._get_conn() as conn:
            wl = self._get_active_whitelist(conn)
            accs, hists, redirs, arts, meds, med_vars, excls, valid_count = self._prepare_in_memory(records, wl)
            if valid_count == 0: return (0, 0, 0)
            cur = conn.cursor()
            existing_ids = {r[0] for r in cur.execute(f"SELECT id FROM articles WHERE id IN ({','.join(['?']*len(arts))})", [a[0] for a in arts]).fetchall()}
            new_cnt = sum(1 for a in arts if a[0] not in existing_ids)
            cur.executemany("INSERT INTO accounts (numeric_id, username, display_name, avatar_url, description, avatar_base64, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT(numeric_id) DO UPDATE SET display_name=coalesce(excluded.display_name, accounts.display_name), avatar_url=coalesce(excluded.avatar_url, accounts.avatar_url), description=case when excluded.description != '' then excluded.description else accounts.description end, avatar_base64=coalesce(accounts.avatar_base64, excluded.avatar_base64), updated_at=excluded.updated_at", accs)
            cur.executemany("INSERT INTO account_profile_histories (account_id, display_name, description, avatar_original_url, avatar_seq, avatar_virtual_key, avatar_base64, observed_at) SELECT ?, ?, ?, ?, ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM account_profile_histories WHERE account_id = ?)", [(h[0], h[1], h[2], h[3], h[4], h[5], h[6], h[7], h[0]) for h in hists])
            cur.executemany("""INSERT INTO articles (id, account_id, conversation_id, reply_to_id, reply_to_handle, created_at, full_text, lang, full_text_ja, full_text_en, full_text_zh, via, is_repost, is_liked, wayback_url, source_name, source_domain, original_url, sotwe_url, nitter_url, twistalker_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                ON CONFLICT(id) DO UPDATE SET account_id=coalesce(excluded.account_id, articles.account_id), full_text=case when length(excluded.full_text) > length(articles.full_text) then excluded.full_text else articles.full_text end, lang=excluded.lang, full_text_ja=coalesce(excluded.full_text_ja, articles.full_text_ja), full_text_en=coalesce(excluded.full_text_en, articles.full_text_en), full_text_zh=coalesce(excluded.full_text_zh, articles.full_text_zh), wayback_url=case when excluded.wayback_url != '' then excluded.wayback_url else articles.wayback_url end, source_name=coalesce(excluded.source_name, articles.source_name), source_domain=coalesce(excluded.source_domain, articles.source_domain), original_url=coalesce(excluded.original_url, articles.original_url), sotwe_url=coalesce(excluded.sotwe_url, articles.sotwe_url), nitter_url=coalesce(excluded.nitter_url, articles.nitter_url), twistalker_url=coalesce(excluded.twistalker_url, articles.twistalker_url)""", arts)
            cur.executemany("INSERT OR REPLACE INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)", redirs)
            if meds:
                cur.executemany("""INSERT INTO media (media_id, article_id, account_id, type, download_url, width, height, thumbnail_url, tweet_urls, media_quality, download_status, failed_reason) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    ON CONFLICT(media_id) DO UPDATE SET account_id=coalesce(excluded.account_id, media.account_id), download_url=case when (media.download_url IS NULL OR media.download_url = '' OR media.download_status = 'DEAD_404') and excluded.download_url != '' then excluded.download_url else media.download_url end, width=case when excluded.width > 0 then excluded.width else media.width end, height=case when excluded.height > 0 then excluded.height else media.height end, thumbnail_url=coalesce(excluded.thumbnail_url, media.thumbnail_url), tweet_urls=coalesce(excluded.tweet_urls, media.tweet_urls), media_quality=coalesce(excluded.media_quality, media_quality), download_status=case when media.download_status = 'DEAD_404' and excluded.download_url != '' then excluded.download_status else media.download_status end""", meds)
            if med_vars:
                cur.executemany("INSERT OR IGNORE INTO media_variants (variant_hash, media_id, article_id, download_url, bit_rate) VALUES (?, ?, ?, ?, ?)", med_vars)
            if excls:
                cur.executemany("INSERT OR IGNORE INTO media_excluded (media_id, article_id, account_id, type, download_url, width, height, download_status, failed_reason, tweet_urls, thumbnail_url, quarantine_reason) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", excls)
            conn.commit()
        return (new_cnt, len(arts) - new_cnt, valid_count)

    def upsert_record(self, data: Dict[str, Any]) -> bool:
        res = self.upsert_batch([data]); return (res[2] if isinstance(res, tuple) else res) == 1
