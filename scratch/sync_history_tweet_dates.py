import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 全 profile_histories の observed_at を、紐づくツイートの最古投稿日に更新
# (もし複数世代ある場合は、その世代の文脈に合わせた期間)
histories = cur.execute(
    "SELECT id, account_id, avatar_seq, display_name FROM"
    " account_profile_histories ORDER BY account_id, avatar_seq ASC"
).fetchall()

for hid, acc_id, seq, dname in histories:
  # 当該アカウントのツイート取得
  tweets = cur.execute(
      "SELECT created_at FROM articles WHERE account_id = ? ORDER BY"
      " created_at ASC",
      (acc_id,),
  ).fetchall()
  if tweets:
    first_date = tweets[0][0]
    last_date = tweets[-1][0]
    # 複数世代がある場合は世代順に按分または最古/最新を割り当て
    total_seq = cur.execute(
        "SELECT COUNT(*) FROM account_profile_histories WHERE account_id = ?",
        (acc_id,),
    ).fetchone()[0]
    if total_seq > 1:
      idx = min(
          int((seq - 1) * len(tweets) / total_seq), len(tweets) - 1
      )  # 世代に応じたツイート代表日
      target_date = tweets[idx][0]
    else:
      target_date = first_date

    cur.execute(
        "UPDATE account_profile_histories SET observed_at = ? WHERE id = ?",
        (target_date, hid),
    )
    print(
        f"History ID {hid} (acc: {acc_id[:8]}.., seq: {seq}) updated"
        f" observed_at -> {target_date}"
    )

conn.commit()
print("Finished syncing observed_at with actual tweet dates.")
