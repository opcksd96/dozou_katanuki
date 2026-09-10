import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()
c.execute("SELECT created_at, level, message FROM pipeline_logs WHERE message LIKE '%x7vQf7QsgEONy7TF%' ORDER BY id DESC LIMIT 10")
for r in c.fetchall():
    print(r)
