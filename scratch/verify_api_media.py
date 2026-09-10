import requests

res = requests.get('http://127.0.0.1:5173/api/media?limit=24&offset=24')
print('Status:', res.status_code)
if res.status_code == 200:
    items = res.json().get('items', [])
    for it in items:
        if it.get('article_id') == '1873298482121023707':
            print('media_id:', it.get('media_id'), 'status:', it.get('download_status'), 'stash_image_id:', it.get('stash_image_id'), 'stash_url:', it.get('stash_url'))
