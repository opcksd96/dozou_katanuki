# plugins/twitter/scraper/core/stash_client.py (100行以下)
import json, os, re, sqlite3, time
from typing import Any, Dict, List, Optional
import requests


class StashClient:
    """Stashapp GraphQL API (:9999) 連携・完全メタデータ注入エンジン (SPEC-STASH-DB-001 / SPEC-STORAGE-001)"""
    TITLE_PATTERN = re.compile(r"^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([A-Za-z0-9_]+)$")

    def __init__(self, endpoint: str = "http://127.0.0.1:9999/graphql"):
        self.endpoint, self.session = endpoint, requests.Session()
        self._studio_cache: Dict[str, str] = {}
        self._performer_cache: Dict[str, str] = {}

    def query(self, query_str: str, variables: Optional[Dict[str, Any]] = None) -> Optional[Dict[str, Any]]:
        try:
            resp = self.session.post(self.endpoint, json={"query": query_str, "variables": variables or {}}, timeout=5)
            if resp.status_code == 200:
                d = resp.json()
                if "data" in d and not d.get("errors"): return d["data"]
        except Exception: pass
        return None

    def find_or_create_studio(self, name: str) -> Optional[str]:
        if not name: return None
        if name in self._studio_cache: return self._studio_cache[name]
        r = self.query("query($q: String!) { findStudios(filter: { q: $q, per_page: 1 }) { studios { id name } } }", {"q": name})
        studios = (r.get("findStudios") or {}).get("studios") if r else []
        for s in (studios or []):
            if s.get("name", "").lower() == name.lower():
                self._studio_cache[name] = str(s["id"]); return self._studio_cache[name]
        c_res = self.query("mutation($in: StudioCreateInput!) { studioCreate(input: $in) { id } }", {"in": {"name": name}})
        if c_res and (c_res.get("studioCreate") or {}).get("id"):
            self._studio_cache[name] = str(c_res["studioCreate"]["id"]); return self._studio_cache[name]
        return None

    def find_or_create_performer(self, name: str, disambiguation: str = "") -> Optional[str]:
        if not name: return None
        if name in self._performer_cache: return self._performer_cache[name]
        r = self.query("query($q: String!) { findPerformers(filter: { q: $q, per_page: 1 }) { performers { id name } } }", {"q": name})
        performers = (r.get("findPerformers") or {}).get("performers") if r else []
        for p in (performers or []):
            if p.get("name", "").lower() == name.lower():
                self._performer_cache[name] = str(p["id"]); return self._performer_cache[name]
        inp: Dict[str, Any] = {"name": name}
        if disambiguation: inp["disambiguation"] = disambiguation
        c_res = self.query("mutation($in: PerformerCreateInput!) { performerCreate(input: $in) { id } }", {"in": inp})
        if c_res and (c_res.get("performerCreate") or {}).get("id"):
            self._performer_cache[name] = str(c_res["performerCreate"]["id"]); return self._performer_cache[name]
        return None

    def find_image_by_path(self, file_path: str) -> Optional[str]:
        bn = os.path.basename(file_path)
        q = """query($p: String!, $bn: String!) {
            findImages(image_filter: { path: { value: $p, modifier: EQUALS } }, filter: { per_page: 1 }) { images { id } }
            f2: findImages(image_filter: { path: { value: $bn, modifier: INCLUDES } }, filter: { per_page: 1 }) { images { id } }
        }"""
        r = self.query(q, {"p": file_path.replace("\\", "/"), "bn": bn})
        if r:
            imgs = r.get("findImages", {}).get("images") or r.get("f2", {}).get("images")
            if imgs: return str(imgs[0]["id"])
        return None

    def find_scene_by_path(self, file_path: str) -> Optional[str]:
        bn = os.path.basename(file_path)
        q = """query($p: String!, $bn: String!) {
            findScenes(scene_filter: { path: { value: $p, modifier: EQUALS } }, filter: { per_page: 1 }) { scenes { id } }
            f2: findScenes(scene_filter: { path: { value: $bn, modifier: INCLUDES } }, filter: { per_page: 1 }) { scenes { id } }
        }"""
        r = self.query(q, {"p": file_path.replace("\\", "/"), "bn": bn})
        if r:
            scns = r.get("findScenes", {}).get("scenes") or r.get("f2", {}).get("scenes")
            if scns: return str(scns[0]["id"])
        return None

    def update_image(self, img_id: str, title: str, details: str = "", urls: Optional[List[str]] = None, date: str = "", studio_id: Optional[str] = None, performer_ids: Optional[List[str]] = None, custom_fields: Optional[Dict[str, Any]] = None) -> bool:
        inp: Dict[str, Any] = {"id": img_id, "title": title}
        if details: inp["details"] = details
        if urls: inp["urls"] = urls; inp["url"] = urls[0]
        if date: inp["date"] = date[:10]
        if studio_id: inp["studio_id"] = studio_id
        if performer_ids: inp["performer_ids"] = performer_ids
        if custom_fields:
            fmt_map = {k: json.dumps(v, ensure_ascii=False) if isinstance(v, (list, dict)) else str(v) for k, v in custom_fields.items()}
            inp["custom_fields"] = {"partial": fmt_map}
        return bool(self.query("mutation($in: ImageUpdateInput!) { imageUpdate(input: $in) { id } }", {"in": inp}))

    def update_scene(self, scn_id: str, title: str, details: str = "", urls: Optional[List[str]] = None, date: str = "", studio_id: Optional[str] = None, performer_ids: Optional[List[str]] = None, custom_fields: Optional[Dict[str, Any]] = None) -> bool:
        inp: Dict[str, Any] = {"id": scn_id, "title": title}
        if details: inp["details"] = details
        if urls: inp["urls"] = urls; inp["url"] = urls[0]
        if date: inp["date"] = date[:10]
        if studio_id: inp["studio_id"] = studio_id
        if performer_ids: inp["performer_ids"] = performer_ids
        if custom_fields:
            fmt_map = {k: json.dumps(v, ensure_ascii=False) if isinstance(v, (list, dict)) else str(v) for k, v in custom_fields.items()}
            inp["custom_fields"] = {"partial": fmt_map}
        return bool(self.query("mutation($in: SceneUpdateInput!) { sceneUpdate(input: $in) { id } }", {"in": inp}))

    def trigger_scan(self, paths: Optional[List[str]] = None) -> bool:
        inp: Dict[str, Any] = {"rescan": False}
        if paths: inp["paths"] = [p.replace("\\", "/") for p in paths]
        return bool(self.query("mutation($in: ScanMetadataInput!) { metadataScan(input: $in) }", {"in": inp}))

    def register_media(self, file_path: str, media_type: str = "image", title: str = "", details: str = "", urls: Optional[List[str]] = None, date: str = "", username: str = "", display_name: str = "", custom_fields: Optional[Dict[str, Any]] = None, max_wait: float = 10.0) -> Optional[str]:
        find_fn = self.find_image_by_path if media_type == "image" else self.find_scene_by_path
        m_id = find_fn(file_path)
        if not m_id:
            scan_target = os.path.dirname(file_path) or file_path
            self.trigger_scan([scan_target]); end_t = time.time() + max_wait
            while time.time() < end_t and not m_id: time.sleep(0.3); m_id = find_fn(file_path)
        if m_id:
            s_id = self.find_or_create_studio(username) if username else None
            p_id = self.find_or_create_performer(username, disambiguation=display_name) if username else None
            p_ids = [p_id] if p_id else None
            (self.update_image if media_type == "image" else self.update_scene)(
                m_id, title=title, details=details, urls=urls, date=date,
                studio_id=s_id, performer_ids=p_ids, custom_fields=custom_fields
            )
        return m_id

    def reconcile_to_db(self, db_path: str) -> int:
        res = self.query("query { allScenes { id title details files { path } } allImages { id title details files { path } } }")
        if not res: return 0
        bound = 0
        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()
            for is_scn, items in [(True, res.get("allScenes", [])), (False, res.get("allImages", []))]:
                col = "stash_scene_id" if is_scn else "stash_image_id"
                for item in items:
                    s_id, title, details, files, matched = str(item.get("id", "")), item.get("title", ""), item.get("details", ""), item.get("files", []), False
                    art_id = None
                    for f in files:
                        bn = os.path.basename(f.get("path", ""))
                        if bn:
                            cur.execute(f"SELECT article_id FROM media WHERE (media_id = ? OR media_id = ? OR download_url LIKE ?)", (bn, os.path.splitext(bn)[0], f"%/{bn}"))
                            row = cur.fetchone()
                            if row: art_id = row[0]
                            cur.execute(f"UPDATE media SET {col} = ?, download_status = 'COMPLETED' WHERE (media_id = ? OR media_id = ? OR media_id LIKE ? OR download_url LIKE ? OR download_url LIKE ?) AND ({col} IS NULL OR {col} = '')",
                                        (s_id, bn, os.path.splitext(bn)[0], f"%_{bn}", f"%/{bn}", f"%/{bn}?%"))
                            if cur.rowcount > 0: bound += cur.rowcount; matched = True
                    if not matched:
                        m = self.TITLE_PATTERN.match(title)
                        if m:
                            art_id = m.group(4)
                            cur.execute(f"UPDATE media SET {col} = ?, download_status = 'COMPLETED' WHERE media_id IN (SELECT media_id FROM media WHERE article_id = ? AND {col} IS NULL LIMIT 1)", (s_id, art_id))
                            bound += cur.rowcount

                    # 親ツイートメタデータを Stash に完全注入・同期
                    if art_id:
                        cur.execute("SELECT a.full_text, a.full_text_ja, a.created_at, a.wayback_url, ac.username, ac.display_name FROM articles a JOIN accounts ac ON a.account_id = ac.numeric_id WHERE a.id = ?", (art_id,))
                        art_row = cur.fetchone()
                        if art_row:
                            ft, ft_ja, created_at, wb_url, uname, dname = art_row[0] or "", art_row[1] or "", str(art_row[2] or ""), art_row[3] or "", art_row[4] or "", art_row[5] or ""
                            txt = f"{ft_ja}\n\n{ft}".strip() if ft_ja and ft_ja != ft else (ft or ft_ja)
                            orig_url = f"https://twitter.com/{uname}/status/{art_id}"
                            urls_list = [orig_url]
                            if wb_url: urls_list.append(wb_url)
                            urls_list.append(f"http://localhost:9999/plugin/x-timeline-middleware/index.html?view=x-timeline&performer={uname}&jump_to_tweet={art_id}")
                            s_id_obj = self.find_or_create_studio(uname) if uname else None
                            p_id_obj = self.find_or_create_performer(uname, disambiguation=dname) if uname else None
                            c_fields = {
                                "tweet_id": art_id,
                                "original_name": dname or uname,
                                "source_system": "X_Wayback",
                                "wayback_url": [wb_url] if wb_url else [],
                                "dead_media": []
                            }
                            m_ts = re.search(r'/web/(\d{14})', wb_url) if wb_url else None
                            if m_ts: c_fields["wayback_timestamp"] = m_ts.group(1)

                            (self.update_scene if is_scn else self.update_image)(
                                s_id, title=title, details=txt, urls=urls_list, date=created_at[:10],
                                studio_id=s_id_obj, performer_ids=[p_id_obj] if p_id_obj else None, custom_fields=c_fields
                            )
            conn.commit()
        return bound
