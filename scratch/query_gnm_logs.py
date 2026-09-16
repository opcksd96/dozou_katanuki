import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

try:
    c.execute("SELECT id, created_at, category, level, message FROM pipeline_logs WHERE message LIKE '%GnmCC%' ORDER BY id ASC")
    rows = c.fetchall()
    print("Found logs:", len(rows))
    for r in rows:
        print(r)
except Exception as e:
    print(e)
