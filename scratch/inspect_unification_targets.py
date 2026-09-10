# scratch/inspect_unification_targets.py
import sqlite3

def main():
    conn = sqlite3.connect('archive.db')
    cur = conn.cursor()

    cur.execute("""
        SELECT count(*) FROM media_variants mv
        JOIN media m ON mv.media_id = m.media_id
        WHERE m.type = 'image' AND m.media_id NOT LIKE '%.%'
    """)
    print('media_variants count for no-ext media:', cur.fetchone()[0])

    cur.execute("SELECT count(*) FROM media_excluded WHERE media_id NOT LIKE '%.%'")
    print('media_excluded without ext count:', cur.fetchone()[0])

    cur.execute('PRAGMA foreign_key_list(media_variants)')
    print('media_variants FKs:', cur.fetchall())

    # Check if any new target media_id (with .jpg) would collide with existing media
    cur.execute("""
        SELECT m.media_id, m.media_id || '.jpg'
        FROM media m
        WHERE m.type = 'image' AND m.media_id NOT LIKE '%.%'
          AND EXISTS (SELECT 1 FROM media m2 WHERE m2.media_id = m.media_id || '.jpg')
    """)
    collisions = cur.fetchall()
    print('Collisions with existing .jpg media:', len(collisions))
    for c in collisions:
        print('  Collision:', c)

if __name__ == '__main__':
    main()
