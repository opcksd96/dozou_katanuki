import urllib.request, json
url = 'http://localhost:5173/api/accounts'
with urllib.request.urlopen(url) as res:
    data = json.loads(res.read().decode())
    items = data if isinstance(data, list) else data.get('items', data.get('accounts', []))
    print(f'Total accounts from API: {len(items)}')
    for a in items[:10]:
        print(f"  {a.get('numeric_id')} | @{a.get('username')} ({a.get('display_name')}) | alias_of: '{a.get('alias_of')}'")
