import sqlite3
import time

conn = sqlite3.connect("archive.db", timeout=60.0)
conn.execute("PRAGMA busy_timeout = 60000")
cur = conn.cursor()

# 1. accounts テーブルに is_whitelist カラムを追加
cols = [c[1] for c in cur.execute("PRAGMA table_info(accounts)").fetchall()]
if "is_whitelist" not in cols:
  cur.execute(
      "ALTER TABLE accounts ADD COLUMN is_whitelist INTEGER DEFAULT 1"
  )
  print("Added 'is_whitelist' column to accounts table.")

# 2. whitelists テーブルの登録状態を反映
whitelisted_handles = [
    r[0].lower()
    for r in cur.execute(
        "SELECT value FROM whitelists WHERE is_active = 1"
    ).fetchall()
]
print("Currently whitelisted handles:", whitelisted_handles)

# 全 accounts の username を小文字に正規化し、is_whitelist をセット
for aid, u in cur.execute("SELECT numeric_id, username FROM accounts").fetchall():
  lower_u = u.lower() if u else ""
  is_wl = 1 if lower_u in whitelisted_handles else 0
  cur.execute(
      """
        UPDATE accounts
        SET username = ?, is_whitelist = ?
        WHERE numeric_id = ?
    """,
      (lower_u, is_wl, aid),
  )

# alias_of も小文字に正規化
cur.execute("""
    UPDATE accounts
    SET alias_of = LOWER(alias_of)
    WHERE alias_of IS NOT NULL AND alias_of != ''
""")

# whitelists の value も小文字に正規化
cur.execute("""
    UPDATE whitelists
    SET value = LOWER(value)
    WHERE value IS NOT NULL
""")

conn.commit()
print("Successfully normalized accounts and set is_whitelist flags.")

print("\n--- Current accounts list ---")
for r in cur.execute(
    "SELECT numeric_id, username, display_name, is_whitelist, alias_of FROM"
    " accounts"
).fetchall():
  print(" ", r)
