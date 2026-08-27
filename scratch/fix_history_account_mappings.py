import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 全 accounts の username -> numeric_id マッピング (小文字比較対応)
accounts = cur.execute("SELECT numeric_id, username FROM accounts").fetchall()
user_to_id = {u.lower(): nid for nid, u in accounts}

print("Available accounts mapping:")
for u, nid in user_to_id.items():
  print(f"  @{u} -> {nid}")

# 全 account_profile_histories を走査し、avatar_virtual_key から本来の account_id を復元
histories = cur.execute(
    "SELECT id, account_id, avatar_virtual_key, display_name FROM"
    " account_profile_histories"
).fetchall()

fixed_count = 0
for hid, cur_acc_id, vkey, dname in histories:
  # vkey の形式: [username]_avatar_[seq] または [username]_gen[seq]
  clean_vkey = vkey.replace("_avatar_", "###").replace("_gen", "###")
  if "###" in clean_vkey:
    raw_user = clean_vkey.split("###")[0].lower()
  else:
    raw_user = vkey.split("_")[0].lower()

  target_acc_id = user_to_id.get(raw_user)
  if target_acc_id and target_acc_id != cur_acc_id:
    cur.execute(
        "UPDATE account_profile_histories SET account_id = ? WHERE id = ?",
        (target_acc_id, hid),
    )
    print(
        f"Fixed history ID {hid} ({vkey} / {dname}): {cur_acc_id} ->"
        f" {target_acc_id} (@{raw_user})"
    )
    fixed_count += 1

# 重複している avatar_virtual_key や無駄なレコードを整理し、各アカウントごとに avatar_seq を 1, 2... にリナンバリング
cur.execute("""
    DELETE FROM account_profile_histories
    WHERE id NOT IN (
        SELECT MIN(id)
        FROM account_profile_histories
        GROUP BY account_id, avatar_virtual_key, display_name
    )
""")

for nid in set(user_to_id.values()):
  h_rows = cur.execute(
      "SELECT id, observed_at FROM account_profile_histories WHERE account_id"
      " = ? ORDER BY observed_at ASC, id ASC",
      (nid,),
  ).fetchall()
  for seq, (hid, _) in enumerate(h_rows, start=1):
    cur.execute(
        "UPDATE account_profile_histories SET avatar_seq = ? WHERE id = ?",
        (seq, hid),
    )

conn.commit()
print(f"\nSuccessfully fixed {fixed_count} profile history mappings.")

print("\n--- Current profile histories count per account ---")
for r in cur.execute("""
    SELECT h.account_id, a.username, a.display_name, COUNT(h.id)
    FROM account_profile_histories h
    LEFT JOIN accounts a ON h.account_id = a.numeric_id
    GROUP BY h.account_id
""").fetchall():
  print(" ", r)

print("\n--- Detail for No14_coco (1749477300754878464) ---")
for r in cur.execute(
    "SELECT id, account_id, avatar_seq, display_name, avatar_virtual_key,"
    " observed_at FROM account_profile_histories WHERE account_id ="
    " '1749477300754878464'"
).fetchall():
  print(" ", r)
