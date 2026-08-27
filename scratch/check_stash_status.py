import os
import sqlite3
import sys

sys.path.insert(0, "plugins/twitter/scraper")
from core.stash_client import StashClient

stash = StashClient()
scenes = stash.find_all_scenes()
images = stash.find_all_images()
print(f"Stash inventory: {len(scenes)} scenes, {len(images)} images")

conn = sqlite3.connect("archive.db")
cur = conn.cursor()
completed = cur.execute(
    "SELECT COUNT(*) FROM media WHERE download_status = 'COMPLETED'"
).fetchone()[0]
bound_scenes = cur.execute(
    "SELECT COUNT(*) FROM media WHERE stash_scene_id IS NOT NULL AND"
    " stash_scene_id != ''"
).fetchone()[0]
bound_images = cur.execute(
    "SELECT COUNT(*) FROM media WHERE stash_image_id IS NOT NULL AND"
    " stash_image_id != ''"
).fetchone()[0]
print(
    f"DB: {completed} COMPLETED, {bound_scenes} bound scenes, {bound_images}"
    " bound images"
)

# 一部バインドされていないCOMPLETEDメディアの調査
unbound = cur.execute(
    "SELECT m.media_id, m.type, ac.username FROM media m JOIN articles a ON"
    " m.article_id = a.id JOIN accounts ac ON a.account_id = ac.numeric_id"
    " WHERE m.download_status = 'COMPLETED' AND (m.stash_image_id IS NULL OR"
    " m.stash_image_id = '') AND (m.stash_scene_id IS NULL OR m.stash_scene_id"
    " = '') LIMIT 10"
).fetchall()
print(f"Unbound COMPLETED samples (total {completed - bound_images - bound_scenes}):")
for u in unbound:
  print(" ", u)
