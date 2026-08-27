# plugins/base/scraper/core/thunder_client.py (SPEC-PLUGIN-001 / 100行以下)
import base64, os, shutil, subprocess, time
from typing import List, Optional


class ThunderClient:
    """Thunder (迅雷) P2SP ダウンロード連携クライアント (SPEC-PLUGIN-001 / 100行以下)"""
    DEFAULT_PATHS = [
        r"C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe",
        r"C:\Program Files\Thunder Network\Thunder\Program\Thunder.exe",
        os.path.expandvars(r"%LOCALAPPDATA%\Programs\Thunder\Thunder.exe"),
    ]

    def __init__(self, executable_path: Optional[str] = None):
        self.exe_path = executable_path or self._detect_executable()

    def _detect_executable(self) -> Optional[str]:
        for p in self.DEFAULT_PATHS:
            if p and os.path.exists(p): return p
        return shutil.which("Thunder.exe") or shutil.which("Thunder")

    def is_available(self) -> bool: return bool(self.exe_path and os.path.exists(self.exe_path))

    @staticmethod
    def encode_thunder_url(url: str) -> str:
        if not url or url.startswith("thunder://"): return url or ""
        return f"thunder://{base64.b64encode(f'AA{url}ZZ'.encode('utf-8')).decode('ascii')}"

    def add_download_task(self, url: str) -> bool:
        if not self.is_available() or not url: return False
        try:
            t_url = self.encode_thunder_url(url)
            subprocess.Popen([self.exe_path, t_url], creationflags=0x00000008 if os.name == "nt" else 0)
            return True
        except Exception: return False

    def add_batch_tasks(self, urls: List[str], max_limit: int = 50) -> int:
        if not self.is_available() or not urls: return 0
        added = 0
        for u in urls[:max_limit]:
            if u and self.add_download_task(u):
                added += 1
                time.sleep(0.15)
        return added
