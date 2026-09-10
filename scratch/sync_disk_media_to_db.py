# scratch/sync_disk_media_to_db.py
import sqlite3, os, glob, shutil

DB_PATH = 'archive.db'
BACKUP_PATH = 'backups/database/archive_pre_media_sync.db'
BASE_DIR = 'G:/Media_Storage/Influencers'

def main():
    print('=== STEP 1: BACKUP DATABASE ===')
    os.makedirs('backups/database', exist_ok=True)
    shutil.copyfile(DB_PATH, BACKUP_PATH)
    print(f'Backed up to {BACKUP_PATH}')

    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    print('\n=== STEP 2: REMOVE DUPLICATE MEDIA RECORDS ===')
    # Check duplicate media in the same article where one has extension and one does not
    cur.execute("SELECT media_id, article_id FROM media WHERE media_id LIKE '%.jpg' OR media_id LIKE '%.png' OR media_id LIKE '%.mp4'")
    ext_records = cur.fetchall()
    removed_dupes = 0
    for mid, aid in ext_records:
        clean_id = mid.rsplit('.', 1)[0]
        cur.execute("SELECT 1 FROM media WHERE article_id = ? AND media_id = ?", (aid, clean_id))
        if cur.fetchone():
            # Duplicate found, delete the extension-appended record
            cur.execute("DELETE FROM media WHERE article_id = ? AND media_id = ?", (aid, mid))
            removed_dupes += 1
            print(f'  Removed duplicate: {mid} for article {aid} (clean {clean_id} exists)')
    conn.commit()
    print(f'Total duplicate records removed: {removed_dupes}')

    print('\n=== STEP 3: SYNC DISK FILES TO COMPLETED ===')
    cur.execute("SELECT media_id, article_id, type, download_status FROM media WHERE download_status != 'COMPLETED'")
    uncompleted = cur.fetchall()
    print(f'Uncompleted media to check: {len(uncompleted)}')

    synced_count = 0
    for mid, aid, mtype, status in uncompleted:
        clean_id = mid.rsplit('.', 1)[0] if '.' in mid else mid
        # Search for matching files on disk
        matches = glob.glob(f'{BASE_DIR}/**/{clean_id}.*', recursive=True)
        if not matches:
            matches = glob.glob(f'{BASE_DIR}/**/{clean_id}', recursive=True)

        if matches and os.path.getsize(matches[0]) > 0:
            cur.execute("""
                UPDATE media
                SET download_status = 'COMPLETED',
                    failed_reason = NULL
                WHERE media_id = ?
            """, (mid,))
            synced_count += 1
            print(f'  [SYNCED] {mid} ({aid}) -> COMPLETED (found {os.path.basename(matches[0])}, {os.path.getsize(matches[0])} bytes)')

    conn.commit()
    print(f'\nTotal media synced to COMPLETED: {synced_count}')

    print('\n=== STEP 4: VERIFY TARGET ARTICLE 1890346081348636973 ===')
    cur.execute("SELECT media_id, article_id, type, download_status, download_url FROM media WHERE article_id = '1890346081348636973'")
    rows = cur.fetchall()
    print(f'Article 1890346081348636973 media count: {len(rows)}')
    for r in rows:
        print(' ', r)

    conn.close()

if __name__ == '__main__':
    main()
