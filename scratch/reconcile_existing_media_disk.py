# scratch/reconcile_existing_media_disk.py
import sqlite3, os, glob

DB_PATH = 'archive.db'
BASE_DIR = 'G:/Media_Storage/Influencers'

def main():
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    cur.execute("SELECT media_id, article_id, type, download_status FROM media WHERE download_status != 'COMPLETED'")
    uncompleted = cur.fetchall()
    print(f'Total uncompleted media in DB: {len(uncompleted)}')

    found_completed = []
    for mid, aid, mtype, status in uncompleted:
        clean_id = mid[:-4] if mid.endswith(('.jpg', '.png', '.mp4', '.webp')) else mid
        matches = glob.glob(f'{BASE_DIR}/**/{clean_id}*', recursive=True)
        if matches:
            found_completed.append((mid, aid, status, matches[0]))

    print(f'Found {len(found_completed)} files that exist on disk but not marked COMPLETED!')
    for item in found_completed:
        print(f'  MediaID: {item[0]} | ArticleID: {item[1]} | CurrentStatus: {item[2]} | File: {item[3]}')

if __name__ == '__main__':
    main()
