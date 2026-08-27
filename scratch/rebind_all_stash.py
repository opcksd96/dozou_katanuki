import os
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
print(f"Stash has {len(scenes)} scenes and {len(images)} images.")

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 1. 画像のバインド
bound_img = 0
for img in images:
  img_id = str(img["id"])
  for f in img.get("files", []):
    p = f.get("path", "")
    bn = os.path.basename(p)
    if bn:
      cur.execute(
          "UPDATE media SET stash_image_id = ?, download_status = 'COMPLETED'"
          " WHERE (media_id = ? OR media_id = ? OR media_id LIKE ?)",
          (img_id, bn, os.path.splitext(bn)[0], f"%{bn}%"),
      )
      bound_img += cur.rowcount

# 2. 動画（シーン）のバインド
bound_scn = 0
for scn in scenes:
  scn_id = str(scn["id"])
  for f in scn.get("files", []):
    p = f.get("path", "")
    bn = os.path.basename(p)
    if bn:
      cur.execute(
          "UPDATE media SET stash_scene_id = ?, download_status = 'COMPLETED'"
          " WHERE (media_id = ? OR media_id = ? OR media_id LIKE ?)",
          (scn_id, bn, os.path.splitext(bn)[0], f"%{bn}%"),
      )
      bound_scn += cur.rowcount

conn.commit()
print(f"Bound updates - Images: {bound_img}, Scenes: {bound_scn}")
b_scenes = cur.execute(
    "SELECT COUNT(*) FROM media WHERE stash_scene_id IS NOT NULL AND"
    " stash_scene_id != ''"
).fetchone()[0]
b_images = cur.execute(
    "SELECT COUNT(*) FROM media WHERE stash_image_id IS NOT NULL AND"
    " stash_image_id != ''"
).fetchone()[0]
completed = cur.execute(
    "SELECT COUNT(*) FROM media WHERE download_status = 'COMPLETED'"
).fetchone()[0]
print(
    f"DB totals - COMPLETED: {completed}, Bound Scenes: {b_scenes}, Bound"
    f" Images: {b_images}"
)
