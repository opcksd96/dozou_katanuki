import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# yike_luo (2cda0ead-b020-efbb-be02-6faba155be4b) のツイート投稿日範囲を確認
rows = cur.execute(
    "SELECT MIN(created_at), MAX(created_at), COUNT(*) FROM articles WHERE"
    " account_id = '2cda0ead-b020-efbb-be02-6faba155be4b'"
).fetchall()
print("yike_luo articles date range:", rows)

# 全アカウントの profile_histories と、紐づくツイートの日付を確認
for r in cur.execute(
    "SELECT id, account_id, avatar_seq, display_name, observed_at FROM"
    " account_profile_histories"
).fetchall():
  # そのアカウントのツイート期間
  acc_dates = cur.execute(
      "SELECT MIN(created_at), MAX(created_at), COUNT(*) FROM articles WHERE"
      " account_id = ?",
      (r[1],),
  ).fetchone()
  print(f"History ID {r[0]} ({r[3]}): Tweet Range = {acc_dates}")
