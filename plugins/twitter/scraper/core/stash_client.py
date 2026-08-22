# plugins/twitter/scraper/core/stash_client.py (100行以下)
import json, os, re, sqlite3, time
from typing import Any, Dict, List, Optional
import requests


class StashClient:
    """Stashapp GraphQL API (:9999) 連携・ミューテーションエンジン (SPEC-STASH-DB-001 / SPEC-STORAGE-001)"""
    TITLE_PATTERN = re.compile(r"^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([A-Za-z0-9_]+)$")

    def __init__(self, endpoint: str = "http://127.0.0.1:9999/graphql"):
        self.endpoint = endpoint
        self.session = requests.Session()

    def query(self, query_str: str, variables: Optional[Dict[str, Any]] = None) -> Optional[Dict[str, Any]]:
        try:
            resp = self.session.post(self.endpoint, json={"query": query_str, "variables": variables or {}}, timeout=5)
            if resp.status_code == 200:
                d = resp.json()
                if "data" in d and not d.get("errors"): return d["data"]
        except Exception: pass
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

    def update_image(self, img_id: str, title: str, details: str = "", url: str = "", date: str = "") -> bool:
        inp = {"id": img_id, "title": title}
        if details: inp["details"] = details
        if url: inp["url"] = url
        if date: inp["date"] = date
        return bool(self.query("mutation($in: ImageUpdateInput!) { imageUpdate(input: $in) { id } }", {"in": inp}))

    def update_scene(self, scn_id: str, title: str, details: str = "", url: str = "", date: str = "") -> bool:
        inp = {"id": scn_id, "title": title}
        if details: inp["details"] = details
        if url: inp["url"] = url
        if date: inp["date"] = date
        return bool(self.query("mutation($in: SceneUpdateInput!) { sceneUpdate(input: $in) { id } }", {"in": inp}))

    def trigger_scan(self, paths: Optional[List[str]] = None) -> bool:
        inp: Dict[str, Any] = {"rescan": False}
        if paths: inp["paths"] = [p.replace("\\", "/") for p in paths]
        return bool(self.query("mutation($in: ScanMetadataInput!) { metadataScan(input: $in) }", {"in": inp}))

    def register_media(self, file_path: str, media_type: str = "image", title: str = "", url: str = "", date: str = "", max_wait: float = 10.0) -> Optional[str]:
        find_fn = self.find_image_by_path if media_type == "image" else self.find_scene_by_path
        m_id = find_fn(file_path)
        if not m_id:
            scan_target = os.path.dirname(file_path) or file_path
            self.trigger_scan([scan_target])
            end_t = time.time() + max_wait
            while time.time() < end_t and not m_id:
                time.sleep(0.3)
                m_id = find_fn(file_path)
        if m_id and title:
            (self.update_image if media_type == "image" else self.update_scene)(m_id, title=title, url=url, date=date)
        return m_id

    def reconcile_to_db(self, db_path: str) -> int:
        res = self.query("query { allScenes { id title files { path } } allImages { id title files { path } } }")
        if not res: return 0
        bound = 0
        with sqlite3.connect(db_path) as conn:
            cur = conn.cursor()
            for is_scn, items in [(True, res.get("allScenes", [])), (False, res.get("allImages", []))]:
                col = "stash_scene_id" if is_scn else "stash_image_id"
                for item in items:
                    s_id, title, files = str(item.get("id", "")), item.get("title", ""), item.get("files", [])
                    m = self.TITLE_PATTERN.match(title)
                    if m:
                        cur.execute(f"UPDATE media SET {col} = ?, download_status = 'COMPLETED' WHERE article_id = ? AND {col} IS NULL", (s_id, m.group(4)))
                        bound += cur.rowcount
                    for f in files:
                        bn = os.path.basename(f.get("path", ""))
                        bn_no_ext = os.path.splitext(bn)[0]
                        if bn:
                            cur.execute(f"UPDATE media SET {col} = ?, download_status = 'COMPLETED' WHERE (media_id = ? OR media_id = ? OR media_id LIKE ? OR download_url LIKE ? OR download_url LIKE ?) AND ({col} IS NULL OR {col} = '')",
                                        (s_id, bn, bn_no_ext, f"%_{bn}", f"%/{bn}", f"%/{bn}?%"))
                            bound += cur.rowcount
            conn.commit()
        return bound
