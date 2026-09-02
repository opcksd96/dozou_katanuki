import sqlite3

con = sqlite3.connect("archive.db", timeout=60.0)
cur = con.cursor()

print("account:", cur.execute("SELECT numeric_id, username, display_name FROM accounts WHERE numeric_id = '2cda0ead-b020-efbb-be02-6faba155be4b'").fetchall())
print("article:", cur.execute("SELECT id, account_id, created_at FROM articles WHERE id = '1747797437467815997'").fetchall())
print("all accounts matching 2cda:", cur.execute("SELECT * FROM accounts WHERE numeric_id LIKE '2cda%'").fetchall())

con.close()
