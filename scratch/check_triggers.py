import sqlite3

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

cur.execute("SELECT media_id, article_id FROM media WHERE media_id LIKE '%GnmCC%'")
print("All GnmCC rows:", cur.fetchall())

# Let's check table definition / triggers
cur.execute("SELECT sql FROM sqlite_master WHERE name = 'media'")
print("\nMedia Table SQL:\n", cur.fetchone()[0])

cur.execute("SELECT name, sql FROM sqlite_master WHERE type = 'trigger' AND tbl_name = 'media'")
print("\nMedia Triggers:\n", cur.fetchall())
