import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("SELECT media_id, article_id, type, download_status, download_url FROM media WHERE article_id LIKE '%1907698917766213981%'")
rows = c.fetchall()
print(f"Total media for 1907698917766213981: {len(rows)}")
for r in rows:
    print(" ", r)

print("\nTotal media in entire media table:")
c.execute("SELECT COUNT(*) FROM media")
print("Total rows:", c.fetchone()[0])
