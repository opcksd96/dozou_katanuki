import urllib.request, json
url = 'http://localhost:5173/api/media?limit=10'
with urllib.request.urlopen(url) as res:
    data = json.loads(res.read().decode())
    items = data.get('items', data.get('data', []))
    for it in items[:6]:
        print(f"Media: {it.get('media_id')} | Status: {it.get('download_status')} | TweetID: {it.get('article_id')} | AccountID: {it.get('account_id')}")
        tt = (it.get('tweet_text') or '')[:50]
        print(f"  Tweet: {tt}")
        print(f"  Account: @{it.get('username')} ({it.get('display_name')})")
