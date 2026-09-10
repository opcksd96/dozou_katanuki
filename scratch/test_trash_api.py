import requests

# 1. 通常アカウント
r1 = requests.get('http://127.0.0.1:5175/api/accounts')
print('Normal accounts status:', r1.status_code, 'count:', len(r1.json()))

# 2. ゴミ箱アカウント
r2 = requests.get('http://127.0.0.1:5175/api/accounts?trash=true')
print('Trashed accounts status:', r2.status_code, 'count:', len(r2.json()))
for a in r2.json():
    print(f"  Trashed: {a.get('username')} ({a.get('numeric_id')}) reason={a.get('trash_reason')}, alias_of={a.get('alias_of')}")

# 3. フロントエンド Vite プロキシ経由
r3 = requests.get('http://127.0.0.1:5173/api/accounts?trash=true')
print('Vite Proxy Trashed status:', r3.status_code, 'count:', len(r3.json()))
