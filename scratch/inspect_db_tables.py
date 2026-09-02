import sqlite3

con = sqlite3.connect("archive.db")
cur = con.cursor()
tables = [r[0] for r in cur.execute("SELECT name FROM sqlite_master WHERE type='table'").fetchall()]
print("Tables:", tables)

for t in ["articles", "media", "accounts", "download_tasks", "thunder_tasks"]:
    if t in tables:
        print(f"=== {t} info ===")
        cols = cur.execute(f"PRAGMA table_info({t})").fetchall()
        for c in cols:
            print(f"  {c[1]} ({c[2]})")
