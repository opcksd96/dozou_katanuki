import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("SELECT * FROM articles WHERE id = '1907698917766213981'")
col_names = [d[0] for d in c.description]
row = c.fetchone()
for k, v in zip(col_names, row):
    print(f"{k}: {v}")
