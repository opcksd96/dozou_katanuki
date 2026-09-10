# scratch/check_media_extension_stats.py
import sqlite3

def main():
    conn = sqlite3.connect('archive.db')
    cur = conn.cursor()

    cur.execute("SELECT count(*) FROM media WHERE type = 'image' AND media_id LIKE '%.%'")
    c_ext = cur.fetchone()[0]
    cur.execute("SELECT count(*) FROM media WHERE type = 'image' AND media_id NOT LIKE '%.%'")
    c_noext = cur.fetchone()[0]
    print(f'Images in media table: with ext = {c_ext}, without ext = {c_noext}')

    cur.execute("SELECT count(*) FROM media WHERE type = 'video' AND media_id LIKE '%.%'")
    v_ext = cur.fetchone()[0]
    cur.execute("SELECT count(*) FROM media WHERE type = 'video' AND media_id NOT LIKE '%.%'")
    v_noext = cur.fetchone()[0]
    print(f'Videos in media table: with ext = {v_ext}, without ext = {v_noext}')

    # Sample without ext
    cur.execute("SELECT media_id, download_url FROM media WHERE type = 'image' AND media_id NOT LIKE '%.%' LIMIT 5")
    print('\nSample images without ext:')
    for r in cur.fetchall():
        print(' ', r)

    # Sample with ext
    cur.execute("SELECT media_id, download_url FROM media WHERE type = 'image' AND media_id LIKE '%.%' LIMIT 5")
    print('\nSample images with ext:')
    for r in cur.fetchall():
        print(' ', r)

if __name__ == '__main__':
    main()
