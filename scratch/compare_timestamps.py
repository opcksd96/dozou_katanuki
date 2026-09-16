import os, sqlite3, time

conn = sqlite3.connect('archive.db')
c = conn.cursor()

print("=== ARTICLE ===")
c.execute("SELECT id, created_at, source_name, sotwe_url, wayback_url FROM articles WHERE id = '1907698917766213981'")
print(c.fetchone())

print("\n=== THUNDER_TASKS ===")
c.execute("SELECT file_name, created_at, updated_at, dispatched_at, reaped_at FROM thunder_tasks WHERE article_id = '1907698917766213981'")
for r in c.fetchall():
    print(r)

warc_file = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\msluo14\1907698917766213981\snapshot.warc.gz"
if os.path.exists(warc_file):
    mtime = os.path.getmtime(warc_file)
    ctime = os.path.getctime(warc_file)
    print("\n=== WARC FILE TIMESTAMP ===")
    print("mtime:", time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(mtime)))
    print("ctime:", time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(ctime)))
