import sqlite3, glob

for db_path in glob.glob('backups/**/*.db', recursive=True) + ['archive.db']:
    try:
        conn = sqlite3.connect(db_path)
        c = conn.cursor()
        c.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='accounts'")
        if not c.fetchone(): continue
        c.execute("PRAGMA table_info(accounts)")
        cols = [col[1] for col in c.fetchall()]
        if 'is_trash' not in cols: continue
        c.execute("SELECT username, is_trash, trash_reason, alias_of FROM accounts WHERE is_trash = 1 OR (trash_reason IS NOT NULL AND trash_reason != '')")
        rows = c.fetchall()
        if rows:
            print(f"{db_path}: {len(rows)} trashed/reasoned accounts:")
            for r in rows:
                print('  ', r)
    except Exception as e:
        print(f"{db_path}: error {e}")
