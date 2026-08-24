# plugins/base/scraper/core/stash_reconciler.py (SPEC-PLUGIN-001 / 100行以下)
import os, re, sqlite3, time
from typing import Any, Callable, Dict, List, Optional
from .stash_client import StashClient

TITLE_PATTERN = re.compile(r"^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([A-Za-z0-9_]+)$")


class StashReconciler:
    """Stashapp メディア登録・DB双方向照合エンジン (SPEC-STASH-DB-001)"""
    def __init__(self, stash: Optional[StashClient] = None):
        self.stash = stash or StashClient()

    def register_media(self, file_path: str, media_type: str = "image", title: str = "", details: str = "", urls: Optional[List[str]] = None, date: str = "", username: str = "", display_name: str = "", custom_fields: Optional[Dict[str, Any]] = None, max_wait: float = 10.0) -> Optional[str]:
        find_fn = self.stash.find_image_by_path if media_type == "image" else self.stash.find_scene_by_path
        m_id = find_fn(file_path)
        if not m_id:
            scan_target = os.path.dirname(file_path) or file_path
            self.stash.trigger_scan([scan_target]); end_t = time.time() + max_wait
            while time.time() < end_t and not m_id: time.sleep(0.3); m_id = find_fn(file_path)
        if m_id:
            s_id = self.stash.find_or_create_studio(username) if username else None
            p_id = self.stash.find_or_create_performer(username, disambiguation=display_name) if username else None
            (self.stash.update_image if media_type == "image" else self.stash.update_scene)(
                m_id, title=title, details=details, urls=urls, date=date,
                studio_id=s_id, performer_ids=[p_id] if p_id else None, custom_fields=custom_fields
            )
        return m_id

    def reconcile_to_db(self, db_path: str, log_fn: Optional[Callable[[str], None]] = None) -> int:
        def _log(msg: str) -> None:
            if log_fn: log_fn(msg)
            print(f"[RECONCILE] {msg}", flush=True)

        _log("Querying Stash GraphQL (:9999) for all scenes and images...")
        res = self.stash.query("query { allScenes { id title details files { path } } allImages { id title details files { path } } }")
        if not res: _log("Stash query failed or Stash service is not reachable (:9999)."); return 0

        scenes, images = res.get("allScenes", []), res.get("allImages", [])
        _log(f"Stash inventory: Found {len(scenes)} scenes, {len(images)} images.")
        bound = 0
        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()
            for is_scn, items in [(True, scenes), (False, images)]:
                col = "stash_scene_id" if is_scn else "stash_image_id"
                for item in items:
                    s_id, title, files, matched, art_id = str(item.get("id", "")), item.get("title", ""), item.get("files", []), False, None
                    for f in files:
                        bn = os.path.basename(f.get("path", ""))
                        if bn:
                            cur.execute("SELECT article_id FROM media WHERE (media_id = ? OR media_id = ? OR download_url LIKE ?)", (bn, os.path.splitext(bn)[0], f"%/{bn}"))
                            row = cur.fetchone()
                            if row: art_id = row[0]
                            cur.execute(f"UPDATE media SET {col} = ?, download_status = 'COMPLETED' WHERE (media_id = ? OR media_id = ? OR media_id LIKE ? OR download_url LIKE ? OR download_url LIKE ?) AND ({col} IS NULL OR {col} = '')",
                                        (s_id, bn, os.path.splitext(bn)[0], f"%_{bn}", f"%/{bn}", f"%/{bn}?%"))
                            if cur.rowcount > 0:
                                bound += cur.rowcount; matched = True
                                _log(f"Bound media '{bn}' -> Stash {'Scene' if is_scn else 'Image'} #{s_id} (Parent: {art_id or 'none'})")
                    if not matched:
                        m = TITLE_PATTERN.match(title)
                        if m:
                            art_id = m.group(4)
                            cur.execute(f"UPDATE media SET {col} = ?, download_status = 'COMPLETED' WHERE media_id IN (SELECT media_id FROM media WHERE article_id = ? AND {col} IS NULL LIMIT 1)", (s_id, art_id))
                            if cur.rowcount > 0:
                                bound += cur.rowcount
                                _log(f"Bound title pattern '{title}' -> Stash #{s_id} (Tweet {art_id})")
                    if art_id: self._sync_parent_article(cur, art_id, s_id, title, is_scn)
            conn.commit()
        _log(f"Reconciliation completed: Bound {bound} assets into archive.db.")
        return bound

    def _sync_parent_article(self, cur: sqlite3.Cursor, art_id: str, s_id: str, title: str, is_scn: bool) -> None:
        cur.execute("SELECT a.full_text, a.full_text_ja, a.created_at, a.wayback_url, ac.username, ac.display_name FROM articles a JOIN accounts ac ON a.account_id = ac.numeric_id WHERE a.id = ?", (art_id,))
        row = cur.fetchone()
        if not row: return
        ft, ft_ja, cr_at, wb_url, uname, dname = row[0] or "", row[1] or "", str(row[2] or ""), row[3] or "", row[4] or "", row[5] or ""
        txt = f"{ft_ja}\n\n{ft}".strip() if ft_ja and ft_ja != ft else (ft or ft_ja)
        urls = [f"https://twitter.com/{uname}/status/{art_id}"] if uname and art_id else []
        if wb_url: urls.append(wb_url)
        s_obj = self.stash.find_or_create_studio(uname) if uname else None
        p_obj = self.stash.find_or_create_performer(uname, disambiguation=dname) if uname else None
        cf = {"tweet_id": art_id, "original_name": dname or uname, "source_system": "X_Wayback", "wayback_url": [wb_url] if wb_url else [], "dead_media": []}
        m_ts = re.search(r'/web/(\d{14})', wb_url) if wb_url else None
        if m_ts: cf["wayback_timestamp"] = m_ts.group(1)
        (self.stash.update_scene if is_scn else self.stash.update_image)(
            s_id, title=title, details=txt, urls=urls, date=cr_at[:10], studio_id=s_obj, performer_ids=[p_obj] if p_obj else None, custom_fields=cf
        )
