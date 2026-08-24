# plugins/base/scraper/core/stash_client.py (SPEC-PLUGIN-001 / 100行以下)
import json, os, requests
from typing import Any, Dict, List, Optional


class StashClient:
    """Stashapp GraphQL API (:9999) 連携基本クライアント (SPEC-STASH-DB-001)"""
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
        studios = (r.get("findStudios") or {}).get("studios") or [] if r else []
        for s in studios:
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
        performers = (r.get("findPerformers") or {}).get("performers") or [] if r else []
        for p in performers:
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
