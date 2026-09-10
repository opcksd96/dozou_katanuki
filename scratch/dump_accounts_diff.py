import sqlite3

def dump_accounts(db_path):
    print(f"=== {db_path} ===")
    conn = sqlite3.connect(db_path)
    c = conn.cursor()
    c.execute("PRAGMA table_info(accounts)")
    cols = [col[1] for col in c.fetchall()]
    print("Columns:", cols)
    has_trash = 'is_trash' in cols
    query = f"SELECT numeric_id, username, display_name, {'is_trash, trash_reason' if has_trash else '0, NULL'}, alias_of FROM accounts"
    c.execute(query)
    for r in c.fetchall():
        print(r)

dump_accounts("backups/database/archive_20260831_082956.db")
dump_accounts("backups/database/archive_snapshot_20260907_205625.db")
dump_accounts("archive.db")
