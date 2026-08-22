# plugins/twitter/scraper/core/aria2_client.py (100行以下)
import json, os, shutil, subprocess, time, uuid
from typing import Any, Dict, List, Optional
import requests


class Aria2Client:
    """Motrix (:16800) / Aria2 (:6800) JSON-RPC 連携クライアント (SPEC-PLUGIN-001)"""
    PORTS = [16800, 6800]

    def __init__(self, endpoint: str = "", secret: str = ""):
        self.endpoint, self.secret, self.session = endpoint, secret, requests.Session()

    def _call(self, method: str, params: Optional[List[Any]] = None) -> Optional[Dict[str, Any]]:
        rpc_params = [f"token:{self.secret}"] if self.secret else []
        if params: rpc_params.extend(params)
        endpoints = [self.endpoint] if self.endpoint else [f"http://127.0.0.1:{p}/jsonrpc" for p in self.PORTS]
        for ep in endpoints:
            try:
                resp = self.session.post(ep, json={"jsonrpc": "2.0", "id": str(uuid.uuid4()), "method": method, "params": rpc_params}, timeout=2)
                if resp.status_code == 200:
                    data = resp.json()
                    if "result" in data:
                        self.endpoint = ep
                        return data["result"]
            except Exception: pass
        return None

    def is_alive(self) -> bool:
        return self._call("aria2.getVersion") is not None

    def launch_and_wait(self, timeout_sec: float = 6.0) -> bool:
        """Motrix が未起動の場合に bin/ やシステムから自動起動して RPC 待機"""
        if self.is_alive(): return True
        cur = os.path.dirname(os.path.abspath(__file__))
        root = os.path.abspath(os.path.join(cur, "../../../.."))
        candidates = [
            os.path.join(root, "bin", "Motrix", "Motrix.exe"),
            os.path.join(root, "bin", "Motrix", "resources", "extra", "win32", "x64", "aria2c.exe"),
            os.path.join(root, "bin", "motrix", "Motrix.exe"), os.path.join(root, "bin", "Motrix.exe"),
            shutil.which("Motrix") or shutil.which("Motrix.exe"),
            os.path.expandvars(r"%LOCALAPPDATA%\Programs\Motrix\Motrix.exe"), os.path.expandvars(r"%ProgramFiles%\Motrix\Motrix.exe"),
        ]
        target = next((p for p in candidates if p and os.path.exists(p)), None)
        try:
            if target:
                if "aria2c" in target.lower():
                    subprocess.Popen([target, "--enable-rpc", "--rpc-listen-port=6800", "--rpc-listen-all=false"], creationflags=0x08000000 if os.name == "nt" else 0)
                else:
                    subprocess.Popen([target], creationflags=0x00000008 if os.name == "nt" else 0)
                print(f"[Aria2Client] Launched ({target}), waiting for RPC...", flush=True)
            else:
                subprocess.Popen(["cmd.exe", "/c", "start", "", "Motrix"], shell=True)
        except Exception as e:
            print(f"[Aria2Client] Auto-launch failed: {e}", flush=True)

        start_t = time.time()
        while time.time() - start_t < timeout_sec:
            time.sleep(0.5)
            if self.is_alive():
                print(f"[Aria2Client] Motrix RPC is ready at {self.endpoint}!", flush=True)
                return True
        return False

    def add_uri(self, urls: List[str], download_dir: str, filename: str) -> Optional[str]:
        """ダウンロードタスクを Motrix (DHT/Torrent/Thunder追跡) に委託し GID を返却"""
        if not self.is_alive() and not self.launch_and_wait(): return None
        options = {"dir": download_dir.replace("\\", "/"), "out": filename, "allow-overwrite": "true", "auto-file-renaming": "false"}
        res = self._call("aria2.addUri", [urls, options])
        return str(res) if res else None

    def tell_status(self, gid: str) -> Optional[Dict[str, Any]]:
        return self._call("aria2.tellStatus", [gid, ["status", "totalLength", "completedLength"]])
