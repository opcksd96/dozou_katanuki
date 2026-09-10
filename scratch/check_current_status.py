import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()
c.execute("SELECT file_name, status, summary_size, error_reason FROM thunder_tasks WHERE file_name LIKE '%x7vQf7QsgEONy7TF%'")
for r in c.fetchall():
    print(r)
