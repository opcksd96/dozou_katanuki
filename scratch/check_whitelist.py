import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

print('=== WHITELISTS ===')
c.execute("SELECT * FROM whitelists")
for r in c.fetchall():
    print(r)

print('\n=== MEDIA_EXCLUDED ===')
c.execute("SELECT * FROM media_excluded WHERE article_id LIKE '%1907698917766213981%' OR media_id LIKE '%GnmCC%'")
for r in c.fetchall():
    print(r)
