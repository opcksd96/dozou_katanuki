# plugins/base/scraper/core/aria2_client.py (SPEC-PLUGIN-001 / 100行以下)
import json, os, shutil, subprocess, time, uuid
from typing import Any, Dict, List, Optional
import requests


class Aria2Client:
    """Motrix Next (:29100) / Motrix (:16800) / Aria2 (:6800) JSON-RPC 連携 (100行以下)"""
    PORTS = [29100, 16800, 6800]

    def __init__(self, endpoint: str = "", secret: str = ""):
        self.endpoint, self.secret, self.session = endpoint, secret, requests.Session()
        if not self.secret or not self.endpoint: self._detect_motrix_config()

    def _detect_motrix_config(self) -> None:
        appdata = os.environ.get("APPDATA", "")
        cfgs = [
            (os.path.join(appdata, "com.motrix.next", "config.json"), ["config", "rpcListenPort"], ["config", "rpcSecret"]),
            (os.path.join(appdata, "com.motrix.next", "system.json"), ["rpc-listen-port"], ["rpc-secret"]),
            (os.path.join(appdata, "Motrix", "settings.json"), ["engine", "rpcPort"], ["engine", "rpcSecret"]),
            (os.path.join(appdata, "Motrix", "system.json"), ["rpc-listen-port"], ["rpc-secret"]),
        ]
        for p, port_path, sec_path in cfgs:
            if os.path.exists(p):
                try:
                    with open(p, "r", encoding="utf-8") as f: d = json.load(f)
                    def _g(obj, keys):
                        for k in keys: obj = obj.get(k, {}) if isinstance(obj, dict) else None
                        return obj
                    port, sec = _g(d, port_path), _g(d, sec_path)
                    if port and int(port) not in self.PORTS: self.PORTS.insert(0, int(port))
                    if sec and not self.secret: self.secret = str(sec)
                except Exception: pass

    def _call(self, method: str, params: Optional[List[Any]] = None) -> Optional[Dict[str, Any]]:
        endpoints = [self.endpoint] if self.endpoint else [f"http://127.0.0.1:{p}/jsonrpc" for p in self.PORTS]
        for ep in endpoints:
            for sec in ([self.secret, ""] if self.secret else [""]):
                rpc_params = [f"token:{sec}"] if sec else []
                if params: rpc_params.extend(params)
                try:
                    resp = self.session.post(ep, json={"jsonrpc": "2.0", "id": str(uuid.uuid4()), "method": method, "params": rpc_params}, timeout=2)
                    if resp.status_code == 200:
                        data = resp.json()
                        if "result" in data: self.endpoint, self.secret = ep, sec; return data["result"]
                except Exception: pass
        return None

    def is_alive(self) -> bool: return self._call("aria2.getVersion") is not None

    def launch_and_wait(self, timeout_sec: float = 6.0) -> bool:
        if self.is_alive(): return True
        self._detect_motrix_config()
        cur = os.path.dirname(os.path.abspath(__file__))
        root = os.path.abspath(os.path.join(cur, "../../../.."))
        candidates = [
            r"C:\Program Files\MotrixNext\motrix-next.exe", os.path.expandvars(r"%LOCALAPPDATA%\Programs\MotrixNext\motrix-next.exe"),
            os.path.join(root, "bin", "Motrix", "Motrix.exe"), shutil.which("motrix-next") or shutil.which("Motrix.exe"),
            os.path.join(root, "bin", "Motrix", "resources", "extra", "win32", "x64", "aria2c.exe"),
        ]
        target = next((p for p in candidates if p and os.path.exists(p)), None)
        try:
            if target and "aria2c" in target.lower():
                subprocess.Popen([target, "--enable-rpc", "--rpc-listen-port=6800", "--rpc-allow-origin-all=true"], creationflags=0x08000000 if os.name == "nt" else 0)
            elif target: subprocess.Popen([target], creationflags=0x00000008 if os.name == "nt" else 0)
        except Exception: pass
        start_t = time.time()
        while time.time() - start_t < timeout_sec:
            time.sleep(0.3); self._detect_motrix_config()
            if self.is_alive(): return True
        return False

    def add_uri(self, urls: List[str], download_dir: str, filename: str) -> Optional[str]:
        if not self.is_alive() and not self.launch_and_wait(): return None
        options = {
            "dir": download_dir.replace("\\", "/"), "out": filename, "allow-overwrite": "true", "auto-file-renaming": "false",
            "check-certificate": "false", "header": ["User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", "Accept: */*"]
        }
        res = self._call("aria2.addUri", [urls, options])
        return str(res) if res else None

    def tell_status(self, gid: str) -> Optional[Dict[str, Any]]:
        return self._call("aria2.tellStatus", [gid, ["status", "totalLength", "completedLength", "errorCode", "errorMessage"]])

    def wait_for_download(self, gid: str, timeout_sec: float = 60.0, poll_sec: float = 1.0) -> bool:
        start_t = time.time()
        while time.time() - start_t < timeout_sec:
            st = self.tell_status(gid)
            if not st: return False
            if st.get("status") == "complete": return True
            if st.get("status") in ("error", "removed"): return False
            time.sleep(poll_sec)
        return False
