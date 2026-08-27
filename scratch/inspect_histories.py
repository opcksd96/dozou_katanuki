import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

print("Schema of account_profile_histories:")
for col in cur.execute("PRAGMA table_info(account_profile_histories)").fetchall():
  print(" ", col)

print("\nHistories for yike_luo (2cda0ead-b020-efbb-be02-6faba155be4b):")
for r in cur.execute(
    "SELECT * FROM account_profile_histories WHERE account_id ="
    " '2cda0ead-b020-efbb-be02-6faba155be4b'"
).fetchall():
  print(" ", r)

print("\nAll account_profile_histories:")
for r in cur.execute(
    "SELECT account_id, avatar_seq, display_name, description, avatar_url,"
    " observed_at FROM account_profile_histories"
).fetchall():
  print(" ", r)
