import sqlite3

conn = sqlite3.connect("archive.db", timeout=60.0)
cur = conn.cursor()

print("--- accounts ---")
for r in cur.execute(
    "SELECT numeric_id, username, display_name FROM accounts"
).fetchall():
  print(" ", r)

print("\n--- account_profile_histories ---")
for r in cur.execute(
    "SELECT id, account_id, display_name, avatar_seq, avatar_virtual_key,"
    " observed_at FROM account_profile_histories"
).fetchall():
  print(" ", r)
