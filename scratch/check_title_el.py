import json, requests, websocket

res = requests.get('http://127.0.0.1:9222/json').json()
ws_url = None
for t in res:
    if 'main-renderer' in t.get('url', ''):
        ws_url = t.get('webSocketDebuggerUrl')
        break

ws = websocket.create_connection(ws_url, timeout=5)

script = """(() => {
    const items = Array.from(document.querySelectorAll('.xly-side-item'));
    return items.map((el, idx) => {
        const text = el.innerText || '';
        if (text.includes('x7vQf7QsgEONy7TF')) {
            const titleEl = el.querySelector('.xly-file-name__ad, .xly-file-name, .xly-side-title');
            const title = titleEl ? titleEl.innerText.trim() : '';
            return {
                idx: idx,
                title: title,
                text: text
            };
        }
        return null;
    }).filter(Boolean);
})()"""

req = {'id': 1, 'method': 'Runtime.evaluate', 'params': {'expression': script, 'returnByValue': True}}
ws.send(json.dumps(req))
result = json.loads(ws.recv())
ws.close()

val = result.get('result', {}).get('result', {}).get('value', [])
for it in val:
    print(it)
