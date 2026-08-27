import os
import sqlite3
import sys

sys.path.insert(0, "plugins/twitter/scraper")
from core.downloader import Downloader

dl = Downloader(db_path="archive.db")
conn = sqlite3.connect("archive.db")
cur = conn.cursor()

rows = cur.execute(
    "SELECT m.media_id, ac.username, m.type, m.download_status FROM media m JOIN articles a ON m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id"
).fetchall()

valid_images = 0
html_fakes = 0
not_found = 0

for m_id, user, m_type, st in rows:
  p = dl.get_target_path(user or "unknown", m_id, m_type)
  if not (os.path.exists(p) and os.path.getsize(p) > 0):
    not_found += 1
    if st == "COMPLETED":
      cur.execute(
          "UPDATE media SET download_status = 'RETAINED', failed_reason ="
          " '実体ファイル不在 (404)' WHERE media_id = ?",
          (m_id,),
      )
    continue

  try:
    with open(p, "rb") as f:
      head = f.read(32)

    # HTML / JSON / XML エラーページの判定
    if (
        head.startswith(b"<!DOCTYPE")
        or head.startswith(b"<html")
        or head.startswith(b"<HTML")
        or head.startswith(b"{\n")
        or head.startswith(b'{"')
        or head.startswith(b"<?xml")
    ):
      html_fakes += 1
      try:
        os.remove(p)
      except Exception:
        pass
      cur.execute(
          "UPDATE media SET download_status = 'RETAINED', failed_reason ="
          " 'Wayback HTMLエラーページ (404)', stash_image_id = NULL,"
          " stash_scene_id = NULL WHERE media_id = ?",
          (m_id,),
      )
    else:
      valid_images += 1
      cur.execute(
          "UPDATE media SET download_status = 'COMPLETED' WHERE media_id = ?",
          (m_id,),
      )
  except Exception:
    pass

conn.commit()
print(f"Valid binary files (COMPLETED): {valid_images}")
print(f"HTML fake error pages purged (RETAINED): {html_fakes}")
print(f"Not found files (RETAINED): {not_found}")
print(
    "Actual DB Status counts:",
    cur.execute(
        "SELECT download_status, COUNT(*) FROM media GROUP BY download_status"
    ).fetchall(),
)
