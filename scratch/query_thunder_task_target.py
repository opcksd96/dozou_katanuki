import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()
c.execute("SELECT * FROM thunder_tasks WHERE file_name LIKE '%x7vQf7QsgEONy7TF%'")
rows = c.fetchall()
cols = [col[0] for col in c.description]
for r in rows:
    print(dict(zip(cols, r)))
