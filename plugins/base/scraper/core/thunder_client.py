# plugins/base/scraper/core/thunder_client.py (SPEC-PLUGIN-001 / 100行以下)
import base64, json, os, shutil, subprocess
from typing import Any, Dict, List, Optional


class ThunderClient:
    """Thunder (迅雷) P2SP ダウンロード連携クライアント (COM & URL Scheme)"""
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

    def is_available(self) -> bool:
        return bool(self.exe_path and os.path.exists(self.exe_path))

    @staticmethod
    def encode_thunder_url(url: str) -> str:
        if not url or url.startswith("thunder://"): return url or ""
        return f"thunder://{base64.b64encode(f'AA{url}ZZ'.encode('utf-8')).decode('ascii')}"

    def add_download_task(self, url: str, file_name: str = "", dest_dir: str = "") -> bool:
        return self.add_batch_tasks([{"url": url, "file_name": file_name, "dest_dir": dest_dir}]) > 0

    def add_batch_tasks(self, tasks: List[Dict[str, Any]], max_limit: int = 50) -> int:
        valid = [t for t in tasks[:max_limit] if t.get("url")]
        if not valid: return 0
        json_data = json.dumps(valid, ensure_ascii=False)
        ps_script = f"""
$tasks = ConvertFrom-Json @'
{json_data}
'@
$ids = @('ThunderAgent.Agent64.1','ThunderAgent.Agent64','ThunderAgent.Agent.1','ThunderAgent.Agent')
$agent = $null
foreach ($id in $ids) {{
    try {{ $a = New-Object -ComObject $id; if ($a) {{ $agent = $a; break }} }} catch {{}}
}}
if ($agent) {{
    foreach ($t in $tasks) {{
        $agent.AddTask($t.url, ($t.file_name -as [string]), ($t.dest_dir -as [string]), '', '', 1, 0, 5)
    }}
    $agent.CommitTasks()
    exit 0
}}
exit 1
"""
        try:
            res = subprocess.run(["powershell.exe", "-NoProfile", "-NonInteractive", "-Command", ps_script], capture_output=True, timeout=15)
            if res.returncode == 0: return len(valid)
        except Exception: pass

        # フォールバック: URL Scheme 起動
        if self.is_available():
            for t in valid:
                try: subprocess.Popen(["rundll32.exe", "url.dll,FileProtocolHandler", self.encode_thunder_url(t["url"])])
                except Exception: pass
            return len(valid)
        return 0

