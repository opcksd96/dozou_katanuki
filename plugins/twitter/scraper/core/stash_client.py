# plugins/twitter/scraper/core/stash_client.py (100行以下)
import json
from typing import Any, Dict, List, Optional
import requests


class StashClient:
    """Stashapp GraphQL API (:9999) 連携クライアント (SPEC-STASH-DB-001)"""

    def __init__(self, endpoint: str = "http://127.0.0.1:9999/graphql"):
        self.endpoint = endpoint
        self.session = requests.Session()

    def query(self, query_str: str, variables: Optional[Dict[str, Any]] = None) -> Optional[Dict[str, Any]]:
        try:
            payload = {"query": query_str, "variables": variables or {}}
            resp = self.session.post(self.endpoint, json=payload, timeout=5)
            if resp.status_code == 200:
                data = resp.json()
                if "errors" in data:
                    return None
                return data.get("data")
        except Exception:
            pass
        return None

    def find_image_by_path(self, file_path: str) -> Optional[str]:
        """物理ファイルパスから Stash Image ID を検索"""
        q = """
        query FindImageByPath($path: String!) {
            findImages(image_filter: { path: { value: $path, modifier: EQUALS } }, filter: { per_page: 1 }) {
                images { id }
            }
        }
        """
        res = self.query(q, {"path": file_path.replace("\\", "/")})
        if res and res.get("findImages", {}).get("images"):
            return str(res["findImages"]["images"][0]["id"])
        return None

    def find_scene_by_path(self, file_path: str) -> Optional[str]:
        """物理ファイルパスから Stash Scene ID を検索"""
        q = """
        query FindSceneByPath($path: String!) {
            findScenes(scene_filter: { path: { value: $path, modifier: EQUALS } }, filter: { per_page: 1 }) {
                scenes { id }
            }
        }
        """
        res = self.query(q, {"path": file_path.replace("\\", "/")})
        if res and res.get("findScenes", {}).get("scenes"):
            return str(res["findScenes"]["scenes"][0]["id"])
        return None

    def trigger_scan(self, paths: List[str]) -> bool:
        """指定パスの Stash スキャンタスクを起動"""
        q = """
        mutation MetadataScan($input: ScanMetadataInput!) {
            metadataScan(input: $input)
        }
        """
        res = self.query(q, {"input": {"paths": paths, "rescan": False}})
        return bool(res and res.get("metadataScan"))
