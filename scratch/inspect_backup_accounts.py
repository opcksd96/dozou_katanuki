import sqlite3

for db_path in ["backups/database/archive_20260831_082956.db", "backups/database/archive_snapshot_20260907_205625.db"]:
    conn = sqlite3.connect(db_path)
    c = conn.cursor()
    c.execute("SELECT * FROM accounts")
    rows = c.fetchall()
    cols = [col[0] for col in c.description]
    print(f"=== {db_path} (Total: {len(rows)}) ===")
    for r in rows:
        d = dict(zip(cols, r))
        # if it was trashed or alias or UUID
        if d.get('is_trash') or d.get('trash_reason') or '-' in d.get('numeric_id', '') or d.get('alias_of'):
            print(f"ID={d.get('numeric_id')}, User={d.get('username')}, Name={d.get('display_name')}, Trash={d.get('is_trash')}, Reason={d.get('trash_reason')}, Alias={d.get('alias_of')}")
