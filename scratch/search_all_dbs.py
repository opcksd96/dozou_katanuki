import sqlite3, glob

targets = ["EBzLD6hjVRyTXvAb", "EwJ7fsA67FadxvXU", "GZ2bywSMlf6B81DA", "rlCDbpAgOdJjPh_p"]

dbs = glob.glob("*.db") + glob.glob("backups/database/*.db")
print("Found DBs:", dbs)

for db in dbs:
    try:
        con = sqlite3.connect(db)
        cur = con.cursor()
        for t in targets:
            # Check media
            try:
                r = cur.execute("SELECT media_id, article_id, download_url FROM media WHERE media_id LIKE ? OR download_url LIKE ?", (f"%{t}%", f"%{t}%")).fetchall()
                if r:
                    print(f"[{db}] media match for {t}: {r}")
            except Exception: pass
            # Check articles
            try:
                r = cur.execute("SELECT id, account_id, created_at FROM articles WHERE full_text LIKE ? OR wayback_url LIKE ? OR original_url LIKE ?", (f"%{t}%", f"%{t}%", f"%{t}%")).fetchall()
                if r:
                    print(f"[{db}] articles match for {t}: {r}")
            except Exception: pass
        con.close()
    except Exception as e:
        print(f"Error {db}: {e}")
