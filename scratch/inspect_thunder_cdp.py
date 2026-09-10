import json, requests, websocket

res = requests.get('http://127.0.0.1:9222/json').json()
ws_url = None
for t in res:
    if 'main-renderer' in t.get('url', ''):
        ws_url = t.get('webSocketDebuggerUrl')
        break

print('Connecting to:', ws_url)
ws = websocket.create_connection(ws_url, timeout=5)

script = """(() => {
	const fmt = (b) => {
		if (!b || b <= 0) return '0B';
		const k = 1024, s = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(b) / Math.log(k));
		return (b / Math.pow(k, i)).toFixed(2) + s[i];
	};
	let sizeMap = {};
	const tb = document.querySelector('.xly-download-tab__operate');
	if (tb && tb.__vue__ && tb.__vue__.taskBaseMap) {
		const m = tb.__vue__.taskBaseMap;
		for (let k in m) { if (m[k] && m[k].taskName) sizeMap[m[k].taskName] = fmt(m[k].fileSize); }
	}
	const items = Array.from(document.querySelectorAll('.xly-side-item'));
	return items.filter(el => !el.querySelector('.xly-icon-restore')).map(el => {
		let text = el.innerText || '';
		const titleEl = el.querySelector('.xly-file-name__ad, .xly-file-name, .xly-side-title');
		const title = titleEl ? titleEl.innerText.trim() : '';
		if (title && sizeMap[title]) { text = text + '\\n' + sizeMap[title]; }
		return text;
	}).filter(t => t && (t.includes('.jpg') || t.includes('.mp4') || t.includes('.png') || t.includes('.webp')));
})()"""

req = {
    "id": 1,
    "method": "Runtime.evaluate",
    "params": {
        "expression": script,
        "returnByValue": True
    }
}
ws.send(json.dumps(req))
result = json.loads(ws.recv())
ws.close()

val = result.get('result', {}).get('result', {}).get('value', [])
print(f"Total extracted blocks: {len(val)}")
for i, block in enumerate(val):
    print(f"--- BLOCK {i} ---")
    print(repr(block))
