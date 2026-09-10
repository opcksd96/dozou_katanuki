import sqlite3
import requests
import os

query = '''query { allImages { id title files { path } } }'''
res = requests.post('http://127.0.0.1:9999/graphql', json={'query': query})
stash_images = {img['id']: [os.path.basename(f['path']) for f in img['files']] for img in res.json()['data']['allImages']}

conn = sqlite3.connect('archive.db')
c = conn.cursor()
c.execute("SELECT media_id, stash_image_id, article_id FROM media WHERE stash_image_id IS NOT NULL AND stash_image_id != ''")
rows = c.fetchall()

mismatched = []
for mid, sid, aid in rows:
    if sid in stash_images:
        files = stash_images[sid]
        clean_files = [f for f in files] + [os.path.splitext(f)[0] for f in files]
        mid_base = os.path.basename(mid)
        mid_clean = os.path.splitext(mid_base)[0]
        if mid_base not in clean_files and mid_clean not in clean_files:
            mismatched.append((mid, sid, aid))

print(f"Total mismatched records found: {len(mismatched)}")
for mid, sid, aid in mismatched:
    print(f"Clearing mismatched stash_image_id={sid} from media_id={mid} (Tweet {aid})")
    c.execute("UPDATE media SET stash_image_id = NULL, download_status = 'ESCALATED' WHERE media_id = ?", (mid,))

conn.commit()
print("Cleaned up successfully.")
