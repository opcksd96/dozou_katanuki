# scratch/unify_media_extensions.py
import sqlite3, os, shutil

DB_PATH = 'archive.db'
BACKUP_PATH = 'backups/database/archive_pre_unify_extensions.db'

def main():
    print('=== STEP 1: BACKUP DATABASE ===')
    os.makedirs('backups/database', exist_ok=True)
    shutil.copyfile(DB_PATH, BACKUP_PATH)
    print(f'Backed up to {BACKUP_PATH}')

    conn = sqlite3.connect(DB_PATH)
    conn.execute("PRAGMA foreign_keys = OFF;")
    cur = conn.cursor()

    print('\n=== STEP 2: INSPECT NON-EXTENSION IMAGES ===')
    cur.execute("SELECT media_id, article_id FROM media WHERE type = 'image' AND media_id NOT LIKE '%.%'")
    no_ext_items = cur.fetchall()
    print(f'Total non-extension images found: {len(no_ext_items)}')

    deleted_dupes = 0
    updated_renames = 0

    for mid, aid in no_ext_items:
        target_mid = f"{mid}.jpg"
        # Check if target_mid already exists in media table
        cur.execute("SELECT article_id FROM media WHERE media_id = ?", (target_mid,))
        existing = cur.fetchone()
        if existing:
            # Already exists (e.g. on original tweet). Delete the unextended one from repost/dupe.
            cur.execute("DELETE FROM media WHERE media_id = ?", (mid,))
            deleted_dupes += 1
            print(f'  [DELETED DUPE] {mid} in article {aid} (already exists as {target_mid} in article {existing[0]})')
        else:
            # Rename media_id to target_mid
            cur.execute("UPDATE media SET media_id = ? WHERE media_id = ?", (target_mid, mid))
            cur.execute("UPDATE media_variants SET media_id = ? WHERE media_id = ?", (target_mid, mid))
            updated_renames += 1
            print(f'  [RENAMED] {mid} -> {target_mid} in article {aid}')

    print(f'\nSummary: Deleted dupes = {deleted_dupes}, Renamed to .jpg = {updated_renames}')

    print('\n=== STEP 3: UNIFY MEDIA_EXCLUDED ===')
    cur.execute("SELECT media_id FROM media_excluded WHERE media_id NOT LIKE '%.%' AND (type = 'image' OR download_url LIKE '%format=jpg%')")
    me_items = cur.fetchall()
    for row in me_items:
        old_id = row[0]
        new_id = f"{old_id}.jpg"
        cur.execute("UPDATE media_excluded SET media_id = ? WHERE media_id = ?", (new_id, old_id))
        print(f'  [MEDIA_EXCLUDED RENAMED] {old_id} -> {new_id}')

    conn.commit()

    print('\n=== STEP 4: VERIFICATION ===')
    cur.execute("SELECT count(*) FROM media WHERE type = 'image' AND media_id NOT LIKE '%.%'")
    rem_no_ext_img = cur.fetchone()[0]
    cur.execute("SELECT count(*) FROM media WHERE type = 'video' AND media_id NOT LIKE '%.%'")
    rem_no_ext_vid = cur.fetchone()[0]
    cur.execute("SELECT count(*) FROM media WHERE type = 'image' AND media_id LIKE '%.jpg'")
    total_jpg_img = cur.fetchone()[0]
    cur.execute("SELECT count(*) FROM media WHERE type = 'video' AND media_id LIKE '%.mp4'")
    total_mp4_vid = cur.fetchone()[0]

    print(f'Images without extension: {rem_no_ext_img}')
    print(f'Videos without extension: {rem_no_ext_vid}')
    print(f'Total .jpg images: {total_jpg_img}')
    print(f'Total .mp4 videos: {total_mp4_vid}')

    # Verify target article 1890346081348636973
    cur.execute("SELECT media_id, article_id, type, download_status FROM media WHERE article_id = '1890346081348636973'")
    print('\nTarget article 1890346081348636973 media:')
    for r in cur.fetchall():
        print(' ', r)

    conn.close()
    print('\nUnification completed successfully!')

if __name__ == '__main__':
    main()
