import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("SELECT * FROM articles WHERE id = '1907698913637302272'")
row = c.fetchone()
cols = [d[0] for d in c.description]
print("Article 1907698913637302272:", dict(zip(cols, row)) if row else 'Article not found')

c.execute("SELECT media_id, article_id, download_status FROM media WHERE article_id = '1907698913637302272'")
for r in c.fetchall():
    print("Media:", r)
