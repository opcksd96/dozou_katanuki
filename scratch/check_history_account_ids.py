import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

print("All account_profile_histories breakdown by account_id:")
for r in cur.execute("""
    SELECT h.account_id, a.username, a.display_name, COUNT(h.id)
    FROM account_profile_histories h
    LEFT JOIN accounts a ON h.account_id = a.numeric_id
    GROUP BY h.account_id
""").fetchall():
  print(" ", r)

print("\nDetail of histories for 1749477300754878464:")
for r in cur.execute(
    "SELECT id, account_id, avatar_seq, display_name, avatar_virtual_key,"
    " observed_at FROM account_profile_histories WHERE account_id ="
    " '1749477300754878464'"
).fetchall():
  print(" ", r)
