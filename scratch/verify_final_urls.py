import sqlite3

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

print("=== VERIFYING twistalker_url & nitter_url IN DATABASE ===")
rows = cur.execute("""
    SELECT id, via, twistalker_url, nitter_url, original_url
    FROM articles 
    WHERE via IN ('TwStalker', 'Nitter')
""").fetchall()

for r in rows:
    print(f"Post {r[0]} | via: {r[1]}")
    print(f"   twistalker_url : {r[2]}")
    print(f"   nitter_url     : {r[3]}")
    print(f"   original_url   : {r[4]}")
