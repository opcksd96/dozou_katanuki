import sqlite3
import sys

sys.path.insert(0, "plugins/base/scraper")
from core.stash_client import StashClient

stash = StashClient()
res = stash.query(
    "query { allScenes { id title files { path } } allImages { id title files {"
    " path } } }"
)
scenes = res.get("allScenes", [])
images = res.get("allImages", [])

print(f"Stash Scenes sample (total {len(scenes)}):")
for s in scenes[:5]:
  print(" ", s)

conn = sqlite3.connect("archive.db")
cur = conn.cursor()
video_media = cur.execute(
    "SELECT media_id, article_id, download_url, stash_scene_id FROM media"
    " WHERE type = 'video' LIMIT 10"
).fetchall()
print(f"DB Video Media sample:")
for v in video_media:
  print(" ", v)
