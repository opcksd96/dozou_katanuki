import json, requests, websocket

res = requests.get('http://127.0.0.1:9222/json').json()
ws_url = None
for t in res:
    if 'main-renderer' in t.get('url', ''):
        ws_url = t.get('webSocketDebuggerUrl')
        break

ws = websocket.create_connection(ws_url, timeout=5)

script = """(() => {
    const all = Array.from(document.querySelectorAll('*'));
    return all.filter(el => (el.innerText || '').includes('5.0MB') || (el.innerText || '').includes('5.0 MB')).map(el => {
        return {
            tagName: el.tagName,
            className: el.className,
            innerText: el.innerText.substring(0, 150)
        };
    }).slice(0, 10);
})()"""

req = {'id': 1, 'method': 'Runtime.evaluate', 'params': {'expression': script, 'returnByValue': True}}
ws.send(json.dumps(req))
result = json.loads(ws.recv())
ws.close()

val = result.get('result', {}).get('result', {}).get('value', [])
print(f"Elements containing 5.0MB: {len(val)}")
for it in val:
    print(it)
