import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# RT記事に紐づいているメディアを、元ツイート記事（親）に再バインド
rows = cur.execute("""
    SELECT m.media_id, a_orig.id, a_rt.id
    FROM media m
    JOIN articles a_rt ON m.article_id = a_rt.id
    JOIN articles a_orig ON a_rt.full_text LIKE '%' || a_orig.id || '%'
    WHERE a_rt.is_repost = 1 AND a_orig.is_repost = 0 AND a_orig.id != a_rt.id
""").fetchall()

rebound = 0
for m_id, orig_id, rt_id in rows:
  cur.execute(
      "UPDATE media SET article_id = ? WHERE media_id = ?", (orig_id, m_id)
  )
  rebound += 1

conn.commit()
print(
    f"Rebound {rebound} media items from RT posts to original parent posts."
)
