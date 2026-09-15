import urllib.request, json, websocket

res = urllib.request.urlopen('http://127.0.0.1:9222/json')
targets = json.loads(res.read().decode())
main_t = [t for t in targets if 'main-renderer' in t.get('url', '')][0]
ws = websocket.create_connection(main_t['webSocketDebuggerUrl'], timeout=5)

expr = """
(() => {
    let res = {};
    for (let i = 0; i < localStorage.length; i++) {
        let k = localStorage.key(i);
        if (k.toLowerCase().includes('path') || k.toLowerCase().includes('dir') || k.toLowerCase().includes('download')) {
            res[k] = localStorage.getItem(k);
        }
    }
    return res;
})()
"""
ws.send(json.dumps({'id': 1, 'method': 'Runtime.evaluate', 'params': {'expression': expr, 'returnByValue': True}}))
r = json.loads(ws.recv())
print("Thunder localStorage:", json.dumps(r.get("result", {}).get("result", {}).get("value"), indent=2, ensure_ascii=False))
