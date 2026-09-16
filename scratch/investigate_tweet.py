import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

print('=== MEDIA SCHEMA ===')
c.execute("PRAGMA table_info(media)")
for r in c.fetchall():
    print(r)

print('\n=== MEDIA FOR ARTICLE ===')
c.execute("SELECT * FROM media WHERE article_id LIKE '%1907698917766213981%'")
col_names = [d[0] for d in c.description]
for r in c.fetchall():
    print(dict(zip(col_names, r)))

print('\n=== THUNDER_TASKS ===')
c.execute("PRAGMA table_info(thunder_tasks)")
for r in c.fetchall():
    print(r)

print('\n=== THUNDER_TASKS RECORDS ===')
c.execute("SELECT * FROM thunder_tasks WHERE file_name LIKE '%GnmCC%' OR media_id IN (SELECT media_id FROM media WHERE article_id LIKE '%1907698917766213981%')")
col_names = [d[0] for d in c.description]
for r in c.fetchall():
    print(dict(zip(col_names, r)))
