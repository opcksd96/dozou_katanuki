import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

print("==================================================")
print(" 📊 DOZOU ARCHIVE.DB 総合ステータスサマリー")
print("==================================================")

print("\n--- 1. thunder_tasks テーブル (全 883 件) ---")
for row in c.execute("SELECT status, count(*) FROM thunder_tasks GROUP BY status ORDER BY count(*) DESC"):
    print(f"  {row[0]:<12}: {row[1]:>4} 件")

print("\n--- 2. media テーブル (download_status 別) ---")
for row in c.execute("SELECT download_status, count(*) FROM media GROUP BY download_status ORDER BY count(*) DESC LIMIT 10"):
    print(f"  {str(row[0]):<16}: {row[1]:>5} 件")

print("\n--- 3. thunder_tasks: 投入中 (RUNNING / ONBOARDED / HOLDING) の最新5件 ---")
for row in c.execute("SELECT id, media_id, file_name, status, summary_size, error_reason, dispatched_at FROM thunder_tasks WHERE status IN ('RUNNING', 'ONBOARDED', 'HOLDING') ORDER BY updated_at DESC LIMIT 5"):
    print(f"  [{row[3]}] {row[2]} | Size: {row[4] or '-'} | Err: {row[5] or '-'} | Disp: {row[6]}")

print("\n--- 4. thunder_tasks: 待機中 (PENDING) の先頭5件 ---")
for row in c.execute("SELECT id, file_name, status, created_at FROM thunder_tasks WHERE status = 'PENDING' ORDER BY created_at ASC LIMIT 5"):
    print(f"  [PENDING] {row[1]} | Created: {row[3]}")
