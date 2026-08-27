import glob
import sqlite3

# 1. 現行 archive.db の accounts.description 確認
conn = sqlite3.connect("archive.db")
cur = conn.cursor()

print("Current archive.db accounts.description:")
for r in cur.execute(
    "SELECT numeric_id, username, display_name, description FROM accounts"
).fetchall():
  print(" ", r)

# 2. バックアップDB群から description を探す
print("\nChecking backup DBs for description...")
for db_file in glob.glob("*.db") + glob.glob("*/*.db"):
  if db_file == "archive.db":
    continue
  try:
    bconn = sqlite3.connect(db_file)
    bcur = bconn.cursor()
    # テーブルチェック
    tbls = [
        t[0]
        for t in bcur.execute(
            "SELECT name FROM sqlite_master WHERE type='table'"
        ).fetchall()
    ]
    if "accounts" in tbls:
      rows = bcur.execute(
          "SELECT numeric_id, username, display_name, description FROM accounts"
          " WHERE description != '' AND description IS NOT NULL"
      ).fetchall()
      if rows:
        print(f"  Found {len(rows)} descriptions in {db_file}:")
        for row in rows:
          print("   ", row)
    bconn.close()
  except Exception as e:
    pass
