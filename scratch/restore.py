import sqlite3

db_main = sqlite3.connect('archive.db')
cur_main = db_main.cursor()
cur_main.execute("SELECT * FROM articles WHERE id = '1914625467719831829'")
row = cur_main.fetchone()
print('Main DB row:', row)

db_backup = sqlite3.connect('backups/database/archive_20260831_082956.db')
cur_backup = db_backup.cursor()
cur_backup.execute("SELECT * FROM articles WHERE id = '1914625467719831829'")
backup_row = cur_backup.fetchone()
print('Backup DB row exists:', backup_row is not None)

if row:
    print('Row exists in main, updating full_text')
    # index 6 is full_text
    cur_main.execute("UPDATE articles SET full_text = ?, is_trash = 0 WHERE id = ?", (backup_row[6], '1914625467719831829'))
    db_main.commit()
    print('Updated.')
elif backup_row:
    print('Row does not exist in main, inserting from backup')
    placeholders = ', '.join(['?'] * len(backup_row))
    cur_main.execute(f"INSERT INTO articles VALUES ({placeholders})", backup_row)
    db_main.commit()
    print('Inserted.')
else:
    print("Not found in backup DB either.")
