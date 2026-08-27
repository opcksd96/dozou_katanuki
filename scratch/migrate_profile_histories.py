import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 1. description カラムの存在チェックと追加
cols = [
    c[1]
    for c in cur.execute("PRAGMA table_info(account_profile_histories)").fetchall()
]
if "description" not in cols:
  cur.execute(
      "ALTER TABLE account_profile_histories ADD COLUMN description TEXT"
      " DEFAULT ''"
  )
  print("Added 'description' column to account_profile_histories.")

# 2. 全アカウントごとに observed_at 昇順で avatar_seq を 1, 2, 3... とリナンバリング
accounts = [
    r[0]
    for r in cur.execute(
        "SELECT DISTINCT account_id FROM account_profile_histories"
    ).fetchall()
]
for acc_id in accounts:
  histories = cur.execute(
      "SELECT id, observed_at, avatar_virtual_key FROM"
      " account_profile_histories WHERE account_id = ? ORDER BY observed_at"
      " ASC, id ASC",
      (acc_id,),
  ).fetchall()
  for seq, (hid, obs, vkey) in enumerate(histories, start=1):
    cur.execute(
        "UPDATE account_profile_histories SET avatar_seq = ? WHERE id = ?",
        (seq, hid),
    )

conn.commit()
print(f"Renumbered profile history sequences across {len(accounts)} accounts.")

# yike_luo の結果を確認
for r in cur.execute(
    "SELECT id, account_id, avatar_seq, display_name, description, observed_at"
    " FROM account_profile_histories WHERE account_id ="
    " '2cda0ead-b020-efbb-be02-6faba155be4b' ORDER BY avatar_seq ASC"
).fetchall():
  print(" ", r)
