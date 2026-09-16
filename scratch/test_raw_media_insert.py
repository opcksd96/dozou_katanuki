import sqlite3

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

meds = [
  ('GnmCCzebYAAkg9-.jpg', '1907698917766213981', '1749477300754878464_msluo14_小罗老师', 'image', 'https://pbs.twimg.com/media/GnmCCzebYAAkg9-.jpg?name=orig', 1538, 2048, None, '["https://twitter.com/MsLuo14/status/1907698917766213981"]', '', 'QUEUED', None),
  ('GnmCCzqbUAALdDB.jpg', '1907698917766213981', '1749477300754878464_msluo14_小罗老师', 'image', 'https://pbs.twimg.com/media/GnmCCzqbUAALdDB.jpg?name=orig', 1538, 2048, None, '["https://twitter.com/MsLuo14/status/1907698917766213981"]', '', 'QUEUED', None),
  ('GnmCCzgbQAAkIG2.jpg', '1907698917766213981', '1749477300754878464_msluo14_小罗老师', 'image', 'https://pbs.twimg.com/media/GnmCCzgbQAAkIG2.jpg?name=orig', 1538, 2048, None, '["https://twitter.com/MsLuo14/status/1907698917766213981"]', '', 'QUEUED', None),
  ('GnmCCzlaMAIPwFr.jpg', '1907698917766213981', '1749477300754878464_msluo14_小罗老师', 'image', 'https://pbs.twimg.com/media/GnmCCzlaMAIPwFr.jpg?name=orig', 1538, 2048, None, '["https://twitter.com/MsLuo14/status/1907698917766213981"]', '', 'QUEUED', None)
]

sql = """
INSERT INTO media (media_id, article_id, account_id, type, download_url, width, height, thumbnail_url, tweet_urls, media_quality, download_status, failed_reason)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(media_id) DO UPDATE SET
  account_id=coalesce(media.account_id, excluded.account_id),
  download_url=case when (media.download_url IS NULL OR media.download_url = '' OR media.download_status = 'DEAD_404') and excluded.download_url != '' then excluded.download_url else media.download_url end,
  width=case when excluded.width > 0 then excluded.width else media.width end,
  height=case when excluded.height > 0 then excluded.height else media.height end,
  thumbnail_url=coalesce(excluded.thumbnail_url, media.thumbnail_url),
  tweet_urls=coalesce(excluded.tweet_urls, media.tweet_urls),
  media_quality=coalesce(excluded.media_quality, media_quality),
  download_status=case when media.download_status = 'DEAD_404' and excluded.download_url != '' then excluded.download_status else media.download_status end
"""

try:
    cur.executemany(sql, meds)
    print("Rowcount:", cur.rowcount)
    conn.commit()
    print("Committed successfully!")
except Exception as e:
    print("Error:", e)

cur.execute("SELECT media_id, article_id, download_status FROM media WHERE article_id = '1907698917766213981'")
rows = cur.fetchall()
print(f"Total rows now: {len(rows)}")
for r in rows:
    print(r)
