# plugins/twitter/scraper/core/aria2_client.py (100行以下)
import json
import uuid
from typing import Any, Dict, List, Optional
import requests


class Aria2Client:
    """Motrix / Aria2 JSON-RPC (:6800) 連携クライアント (SPEC-PLUGIN-001)"""

    def __init__(self, endpoint: str = "http://127.0.0.1:6800/jsonrpc", secret: str = ""):
        self.endpoint = endpoint
        self.secret = secret
        self.session = requests.Session()

    def _call(self, method: str, params: Optional[List[Any]] = None) -> Optional[Dict[str, Any]]:
        rpc_params: List[Any] = []
        if self.secret:
            rpc_params.append(f"token:{self.secret}")
        if params:
            rpc_params.extend(params)

        payload = {
            "jsonrpc": "2.0",
            "id": str(uuid.uuid4()),
            "method": method,
            "params": rpc_params,
        }
        try:
            resp = self.session.post(self.endpoint, json=payload, timeout=3)
            if resp.status_code == 200:
                data = resp.json()
                if "result" in data:
                    return data["result"]
        except Exception:
            pass
        return None

    def add_uri(self, urls: List[str], download_dir: str, filename: str) -> Optional[str]:
        """ダウンロードタスクを Aria2/Motrix に委託し GID を返却"""
        options = {
            "dir": download_dir.replace("\\", "/"),
            "out": filename,
            "allow-overwrite": "true",
            "auto-file-renaming": "false",
        }
        res = self._call("aria2.addUri", [urls, options])
        return str(res) if res else None

    def tell_status(self, gid: str) -> Optional[Dict[str, Any]]:
        """タスクの進行状況 (active, waiting, paused, error, complete) を取得"""
        return self._call("aria2.tellStatus", [gid, ["status", "totalLength", "completedLength"]])

    def is_alive(self) -> bool:
        """Motrix/Aria2 RPC が稼働中か確認"""
        res = self._call("aria2.getVersion")
        return res is not None
