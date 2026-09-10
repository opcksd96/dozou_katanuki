import json, requests, websocket

res = requests.get('http://127.0.0.1:9222/json').json()
ws_url = None
for t in res:
    if 'main-renderer' in t.get('url', ''):
        ws_url = t.get('webSocketDebuggerUrl')
        break

ws = websocket.create_connection(ws_url, timeout=5)

script = """(() => {
    const tb = document.querySelector('.xly-download-tab__operate');
    if (!tb || !tb.__vue__ || !tb.__vue__.taskBaseMap) return { error: 'no taskBaseMap' };
    const m = tb.__vue__.taskBaseMap;
    const res = {};
    for (let k in m) {
        if (m[k] && m[k].taskName && m[k].taskName.includes('x7vQf7QsgEONy7TF')) {
            res[k] = m[k];
        }
    }
    return res;
})()"""

req = {'id': 1, 'method': 'Runtime.evaluate', 'params': {'expression': script, 'returnByValue': True}}
ws.send(json.dumps(req))
result = json.loads(ws.recv())
ws.close()

val = result.get('result', {}).get('result', {}).get('value', {})
print("Vue taskBaseMap match:")
print(json.dumps(val, indent=2, ensure_ascii=False))
