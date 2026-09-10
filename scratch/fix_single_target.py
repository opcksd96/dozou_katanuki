# scratch/fix_single_target.py
import sqlite3

def main():
    conn = sqlite3.connect('archive.db', timeout=30.0)
    cur = conn.cursor()
    cur.execute("""
        UPDATE media
        SET stash_image_id = NULL,
            download_status = 'ESCALATED',
            failed_reason = 'Escalated to Thunder'
        WHERE media_id = 'Gf9LCaHaMAAWP3B.jpg'
    """)
    conn.commit()
    cur.execute("SELECT media_id, stash_image_id, download_status FROM media WHERE media_id = 'Gf9LCaHaMAAWP3B.jpg'")
    print('Updated row:', cur.fetchone())
    conn.close()

if __name__ == '__main__':
    main()
