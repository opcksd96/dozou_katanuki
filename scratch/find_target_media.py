import sqlite3

con = sqlite3.connect("archive.db", timeout=60.0)
cur = con.cursor()
target = "GEFsflnacAAOXck"

print(f"[*] Searching for target: {target}")

for tbl in ["media", "media_excluded", "thunder_tasks", "download_tasks"]:
    try:
        cols = [c[1] for c in cur.execute(f"PRAGMA table_info({tbl})").fetchall()]
        print(f"=== {tbl} ===")
        found = False
        for c in cols:
            rows = cur.execute(f"SELECT * FROM {tbl} WHERE {c} LIKE ?", (f"%{target}%",)).fetchall()
            if rows:
                print(f"  Matched in column [{c}]: {rows}")
                found = True
        if not found:
            print("  No match.")
    except Exception as e:
        print(f"  Error checking {tbl}: {e}")

print("=== articles search ===")
for c in ["id", "full_text", "wayback_url", "original_url", "sotwe_url"]:
    try:
        rows = cur.execute(f"SELECT id, account_id, created_at FROM articles WHERE {c} LIKE ?", (f"%{target}%",)).fetchall()
        if rows:
            print(f"  Matched articles in [{c}]: {rows}")
    except Exception as e:
        print(f"  Error on articles {c}: {e}")

print("=== Check url_redirects ===")
try:
    rows = cur.execute("SELECT * FROM url_redirects WHERE short_url LIKE ? OR expanded_url LIKE ?", (f"%{target}%", f"%{target}%")).fetchall()
    print("  url_redirects:", rows)
except Exception as e:
    print(f"  url_redirects error: {e}")

con.close()
