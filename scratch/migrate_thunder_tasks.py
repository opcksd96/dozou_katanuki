import sqlite3

conn = sqlite3.connect('archive.db')
cursor = conn.cursor()

cols = [
    ('thunder_task_id', 'INTEGER DEFAULT 0'),
    ('file_size', 'INTEGER DEFAULT 0'),
    ('download_size', 'INTEGER DEFAULT 0'),
    ('save_path', "TEXT DEFAULT ''"),
    ('task_status_code', 'INTEGER DEFAULT 0'),
    ('error_code', 'INTEGER DEFAULT 0'),
    ('detail_text', "TEXT DEFAULT ''")
]

cursor.execute('PRAGMA table_info(thunder_tasks)')
existing = [row[1] for row in cursor.fetchall()]

for name, defn in cols:
    if name not in existing:
        cursor.execute(f'ALTER TABLE thunder_tasks ADD COLUMN {name} {defn}')
        print(f'Added column: {name}')
    else:
        print(f'Column already exists: {name}')

cursor.execute('''
    UPDATE thunder_tasks 
    SET file_size = 0,
        download_size = 0,
        task_status_code = 9,
        error_code = 402,
        detail_text = CASE 
            WHEN error_reason IS NOT NULL AND error_reason != '' THEN error_reason 
            ELSE '原始资源不存在，且未找到候选资源，无法继续下载' 
        END
    WHERE status = 'REAPED'
''')
print('Updated REAPED rows:', cursor.rowcount)
conn.commit()

# 確認
cursor.execute('''
    SELECT id, file_name, status, file_size, download_size, task_status_code, error_code, detail_text 
    FROM thunder_tasks 
    WHERE status = 'REAPED' 
    LIMIT 3
''')
for r in cursor.fetchall():
    print('Sample REAPED row:', r)

conn.close()
