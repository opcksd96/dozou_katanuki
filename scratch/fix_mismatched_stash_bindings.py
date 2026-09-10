# scratch/fix_mismatched_stash_bindings.py
import urllib.request, json, sqlite3, os, glob, shutil

DB_PATH = 'archive.db'
BACKUP_PATH = 'backups/database/archive_pre_fix_stash_mismatch.db'
BASE_DIR = 'G:/Media_Storage/Influencers'

def main():
    print('=== STEP 1: BACKUP DATABASE ===')
    os.makedirs('backups/database', exist_ok=True)
    shutil.copyfile(DB_PATH, BACKUP_PATH)
    print(f'Backed up to {BACKUP_PATH}')

    print('\n=== STEP 2: QUERY STASH FOR ALL IMAGES & FILES ===')
    req = urllib.request.Request(
        'http://127.0.0.1:9999/graphql',
        data=json.dumps({'query': '{ allImages { id files { basename } } }'}).encode(),
        headers={'Content-Type': 'application/json'}
    )
    with urllib.request.urlopen(req) as res:
        data = json.loads(res.read().decode())
        images = data.get('data', {}).get('allImages', [])

    stash_images = {} # stash_id -> list of basenames
    file_to_stash = {} # clean_basename -> stash_id
    for img in images:
        s_id = str(img['id'])
        basenames = [f['basename'] for f in img.get('files', []) if f.get('basename')]
        stash_images[s_id] = basenames
        for bn in basenames:
            file_to_stash[bn] = s_id
            file_to_stash[bn.split('.')[0]] = s_id
    print(f'Loaded {len(images)} Stash images with {len(file_to_stash)} file mappings.')

    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    cur.execute("SELECT media_id, article_id, stash_image_id, download_status FROM media WHERE stash_image_id IS NOT NULL AND stash_image_id != ''")
    db_rows = cur.fetchall()

    rebound = 0
    cleared = 0

    for mid, aid, s_id, status in db_rows:
        files = stash_images.get(s_id, [])
        clean_mid = mid.split('.')[0]
        is_match = any(mid == fb or clean_mid == fb.split('.')[0] for fb in files)

        if not is_match:
            # Check if this mid actually has a different stash_id
            real_stash_id = file_to_stash.get(mid) or file_to_stash.get(clean_mid)
            if real_stash_id:
                cur.execute("UPDATE media SET stash_image_id = ?, download_status = 'COMPLETED' WHERE media_id = ?", (real_stash_id, mid))
                rebound += 1
                print(f'  [RE-BOUND] {mid} in article {aid}: wrong #{s_id} -> real #{real_stash_id}')
            else:
                # Does not exist in Stash! Check if file exists on disk
                disk_files = glob.glob(f'{BASE_DIR}/**/{clean_mid}.*', recursive=True)
                if disk_files and os.path.getsize(disk_files[0]) > 0:
                    # Disk file exists, just clear stash_image_id so it serves from local
                    cur.execute("UPDATE media SET stash_image_id = NULL, download_status = 'COMPLETED' WHERE media_id = ?", (mid,))
                    print(f'  [LOCAL ONLY] {mid} in article {aid}: cleared wrong #{s_id}, serves from local disk')
                else:
                    # Does not exist on disk either! Revert to ESCALATED
                    cur.execute("""
                        UPDATE media
                        SET stash_image_id = NULL,
                            download_status = 'ESCALATED',
                            failed_reason = '迅雷投入中 (未取得)'
                        WHERE media_id = ?
                    """, (mid,))
                    cleared += 1
                    print(f'  [CLEARED] {mid} in article {aid}: cleared wrong #{s_id} -> ESCALATED (not on disk)')

    conn.commit()
    print(f'\nSummary: Re-bound to correct Stash ID = {rebound}, Cleared mismatch = {cleared}')

    print('\n=== STEP 3: VERIFY ARTICLE 1873298482121023707 ===')
    cur.execute("SELECT media_id, article_id, stash_image_id, download_status, failed_reason FROM media WHERE article_id = '1873298482121023707'")
    for r in cur.fetchall():
        print(' ', r)

    conn.close()
    print('\nFix completed successfully!')

if __name__ == '__main__':
    main()
