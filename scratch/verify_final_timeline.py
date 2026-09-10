import urllib.request, json

url = 'http://localhost:5173/api/timeline?platform=twitter&account_id=ext_subyike_YIKE&filter=all&limit=5&offset=0'
with urllib.request.urlopen(url) as res:
    data = json.loads(res.read().decode())
    print(f"Timeline items count: {len(data)}")
    for i, it in enumerate(data):
        c = it.get('content', {})
        t = c.get('original', '')
        print(f"Item {i+1}: ID {it.get('id')} | Date: {it.get('created_at')} | Text: {t[:60]}")
