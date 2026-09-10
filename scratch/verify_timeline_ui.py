import urllib.request, json
url = 'http://localhost:5173/api/timeline?platform=twitter&account_id=ext_subyike_YIKE&filter=all&limit=5&offset=0'
with urllib.request.urlopen(url) as res:
    data = json.loads(res.read().decode())
    print('Current Timeline count for subyike:', len(data))
    for it in data:
        t = it.get('original_text', '')
        print(f"Post {it.get('id')} | Date: {it.get('created_at')} | Text: {t[:60]}")
