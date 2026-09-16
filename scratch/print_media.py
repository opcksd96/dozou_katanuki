import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("PRAGMA table_info(media)")
cols = [r[1] for r in c.fetchall()]
print("Cols:", cols)

c.execute("SELECT * FROM media WHERE article_id LIKE '%1907698917766213981%'")
for r in c.fetchall():
    print(dict(zip(cols, r)))
