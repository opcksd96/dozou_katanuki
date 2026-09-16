import sqlite3, os

for db in ['archive.db', 'dozou_katanuki.db', 'katanuki.db']:
    if os.path.exists(db):
        conn = sqlite3.connect(db)
        c = conn.cursor()
        c.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='pipeline_logs'")
        if c.fetchone():
            print('Found pipeline_logs in', db)
            c.execute("SELECT id, created_at, category, level, message FROM pipeline_logs WHERE message LIKE '%GnmCC%' ORDER BY id ASC")
            for r in c.fetchall():
                print(r)
