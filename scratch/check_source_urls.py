import sqlite3

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

rows = cur.execute("""
    SELECT id, via, substr(original_url, 1, 45), nitter_url, twistalker_url, sotwe_url
    FROM articles 
    WHERE via IN ('TwStalker', 'Nitter', 'twistalker', 'nitter')
       OR nitter_url IS NOT NULL
       OR twistalker_url IS NOT NULL
""").fetchall()

print(f"Total found: {len(rows)}")
for r in rows:
    print(r)
