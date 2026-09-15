import urllib.request
import json
import websocket # if available, or just http

print("=== 迅雷 CDP (port 9222) 疎通確認 ===")
try:
    with urllib.request.urlopen('http://127.0.0.1:9222/json', timeout=2) as r:
        targets = json.loads(r.read())
        print(f"CDP Targets: {len(targets)} 件検出")
        for t in targets:
            print(f"  - Title: {t.get('title')} | URL: {t.get('url')[:60]}")
except Exception as e:
    print(f"❌ CDP 接続エラー: {e}")
