import json
import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 1. tweets_master.json を読み込み
with open(
    "bin/stash/plugins/x-timeline-middleware/tweets_master.json",
    "r",
    encoding="utf-8",
) as f:
  data = json.load(f)

# アカウントごとの bio 収集
author_bios = {}
for entry in data:
  author = entry.get("author", {})
  handle = author.get("handle", "") or author.get("username", "")
  bio = (
      author.get("bio", "")
      or author.get("description", "")
      or author.get("about", "")
  )
  dname = author.get("display_name", "") or author.get("name", "")
  if handle and bio:
    author_bios[handle.lower()] = {
        "bio": bio,
        "display_name": dname,
        "handle": handle,
    }

# sotwe.html から Yike_Luo の bio も抽出補完
yike_bio = (
    "大号已封，且看且珍惜。🚪200加tg瑟瑟群及vx，入门可定制视频和线下，约拍请私信。我的小号请尽早关注@TeacherXiaoLuo"
)
author_bios["yike_luo"] = {
    "bio": yike_bio,
    "display_name": "罗亦可",
    "handle": "yike_luo",
}

print(f"Collected {len(author_bios)} bios from master sources:")
for h, info in author_bios.items():
  print(f"  @{h}: {info['bio'][:40]}...")

# 2. accounts テーブルの description を更新
updated_accounts = 0
for handle_lower, info in author_bios.items():
  cur.execute(
      """
        UPDATE accounts
        SET description = ?
        WHERE LOWER(username) = ?
    """,
      (info["bio"], handle_lower),
  )
  if cur.rowcount > 0:
    updated_accounts += cur.rowcount
    print(f"  Updated account @{handle_lower}")

# 3. account_profile_histories テーブルの description を更新
updated_histories = 0
for handle_lower, info in author_bios.items():
  # account_id を特定
  acc_row = cur.execute(
      "SELECT numeric_id FROM accounts WHERE LOWER(username) = ?",
      (handle_lower,),
  ).fetchone()
  if acc_row:
    acc_id = acc_row[0]
    cur.execute(
        """
            UPDATE account_profile_histories
            SET description = ?
            WHERE account_id = ?
        """,
        (info["bio"], acc_id),
    )
    updated_histories += cur.rowcount

conn.commit()
print(
    f"\nDone! Updated {updated_accounts} accounts and {updated_histories}"
    " profile histories with authentic bio/descriptions."
)

# 確認
print("\n--- Current accounts and descriptions ---")
for r in cur.execute(
    "SELECT numeric_id, username, display_name, description FROM accounts"
).fetchall():
  print(" ", r)
