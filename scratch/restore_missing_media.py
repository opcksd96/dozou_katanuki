import sqlite3

backup_db = "backups/database/archive_20260825_125117.db"
main_db = "archive.db"

b_conn = sqlite3.connect(backup_db)
b_cur = b_conn.cursor()

m_conn = sqlite3.connect(main_db, timeout=60.0)
m_cur = m_conn.cursor()

b_cols = [c[1] for c in b_cur.execute("PRAGMA table_info(media)").fetchall()]
b_media = b_cur.execute(f"SELECT {', '.join(b_cols)} FROM media").fetchall()
print(f"[*] Total media in backup: {len(b_media)}")

inserted = 0
for row in b_media:
    d = dict(zip(b_cols, row))
    mid = d.get("media_id")
    aid = d.get("article_id")
    art_row = m_cur.execute("SELECT account_id FROM articles WHERE id = ?", (aid,)).fetchone()
    acc_id = art_row[0] if art_row else None
    
    exists = m_cur.execute("SELECT 1 FROM media WHERE media_id = ?", (mid,)).fetchone()
    if not exists and acc_id:
        m_cur.execute("""
        INSERT INTO media (media_id, article_id, account_id, type, download_url, width, height, download_status)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        """, (mid, aid, acc_id, d.get("type", "video"), d.get("download_url", ""), d.get("width", 0), d.get("height", 0), d.get("download_status", "COMPLETED")))
        inserted += 1

m_conn.commit()
print(f"[+] Restored {inserted} missing media records from backup into archive.db")

b_conn.close()
m_conn.close()
