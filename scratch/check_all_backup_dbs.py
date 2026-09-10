import sqlite3, glob

for db_path in glob.glob('backups/**/*.db', recursive=True):
    try:
        conn = sqlite3.connect(db_path)
        c = conn.cursor()
        c.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='accounts'")
        if not c.fetchone(): continue
        c.execute("SELECT count(*) FROM accounts")
        cnt = c.fetchone()[0]
        c.execute("PRAGMA table_info(accounts)")
        cols = [col[1] for col in c.fetchall()]
        trash_info = 'no is_trash'
        if 'is_trash' in cols:
            c.execute("SELECT count(*) FROM accounts WHERE is_trash = 1")
            t_cnt = c.fetchone()[0]
            trash_info = f"trashed={t_cnt}"
        print(f"{db_path}: total={cnt}, {trash_info}")
    except Exception as e:
        print(f"{db_path}: err {e}")
