import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()
c.execute("SELECT media_id, download_status, stash_scene_id, stash_image_id FROM media WHERE media_id LIKE '%x7vQf7QsgEONy7TF%'")
for r in c.fetchall():
    print(r)
