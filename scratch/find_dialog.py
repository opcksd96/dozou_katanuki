import urllib.request, json, websocket

res = urllib.request.urlopen('http://127.0.0.1:9222/json')
targets = json.loads(res.read().decode())
for i, t in enumerate(targets):
    ws_url = t.get('webSocketDebuggerUrl')
    if not ws_url:
        continue
    try:
        ws = websocket.create_connection(ws_url, timeout=2)
        ws.send(json.dumps({'id': 1, 'method': 'Runtime.evaluate', 'params': {'expression': 'document.body.innerText.slice(0, 200)', 'returnByValue': True}}))
        txt = json.loads(ws.recv()).get('result', {}).get('result', {}).get('value', '')
        if any(k in txt for k in ['添加', 'FDHB', '存储', '下载']):
            print(f"MATCH {i}: {t.get('title')} | {t.get('url')}")
            print("Text:", repr(txt[:100]))
    except Exception as e:
        pass
