import sqlite3, time

ts = time.strftime('%Y%m%d_%H%M%S')
p = f'backups/database/archive_snapshot_{ts}.db'
conn = sqlite3.connect('archive.db')
conn.execute(f"VACUUM INTO '{p}'")
conn.close()
print(f'Snapshot created: {p}')
