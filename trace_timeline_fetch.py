"""
trace_timeline_fetch.py
ブラウザ上で fetchTimeline をステップ実行し、どの行で何が発生しているか精密調査
"""
import time, os, json
from seleniumbase import SB

def main():
    with SB(uc=False, headless=True) as sb:
        sb.open("http://127.0.0.1:5175/")
        time.sleep(3)

        trace = sb.execute_script("""
            return (async () => {
                const logs = [];
                try {
                    logs.push("Step 1: fetching /api/timeline...");
                    const res = await fetch('/api/timeline?platform=twitter&account_id=all&filter=all&limit=50&offset=0');
                    logs.push(`Step 2: status = ${res.status}`);
                    const json = await res.json();
                    logs.push(`Step 3: json is array? ${Array.isArray(json)}, length = ${Array.isArray(json) ? json.length : 'not array'}`);
                    
                    if (Array.isArray(json) && json.length > 0) {
                        logs.push(`Step 4: first item id = ${json[0].id}, author = ${json[0].author?.handle}`);
                    }
                } catch (e) {
                    logs.push(`ERROR: ${e.message}`);
                }
                return logs;
            })();
        """)
        print("[*] Trace Result:")
        for line in trace:
            print("  ", line)

if __name__ == "__main__":
    main()
