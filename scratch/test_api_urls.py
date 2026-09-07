import urllib.request, json

url = 'http://localhost:5175/api/timeline?platform=twitter&account_id=1553745049057300481&filter=all&limit=5&offset=0'
resp = urllib.request.urlopen(url)
data = json.loads(resp.read().decode('utf-8'))
items = data if isinstance(data, list) else data.get('items', [])

print(f"API Returned {len(items)} items:")
for it in items:
    print(f"  Post {it.get('id')} | nitter_url: {it.get('nitter_url')} | twistalker_url: {it.get('twistalker_url')}")
