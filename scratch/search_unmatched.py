import sqlite3

unmatched = ["EBzLD6hjVRyTXvAb", "EwJ7fsA67FadxvXU", "GZ2bywSMlf6B81DA", "rlCDbpAgOdJjPh_p"]
con = sqlite3.connect("archive.db")
cur = con.cursor()

print("=== Searching for 4 unmatched media ===")
for u in unmatched:
    print(f"\nTarget: {u}")
    # Search articles
    rows = cur.execute("SELECT id, account_id, substr(full_text,1,40), wayback_url FROM articles WHERE full_text LIKE ? OR wayback_url LIKE ? OR original_url LIKE ? OR sotwe_url LIKE ?", (f"%{u}%", f"%{u}%", f"%{u}%", f"%{u}%")).fetchall()
    print("  articles match:", rows)
    # Search all tables
    for tbl in ["accounts", "account_profile_histories", "url_redirects", "whitelists", "media_excluded"]:
        cols = [c[1] for c in cur.execute(f"PRAGMA table_info({tbl})").fetchall()]
        for c in cols:
            r = cur.execute(f"SELECT * FROM {tbl} WHERE {c} LIKE ?", (f"%{u}%",)).fetchall()
            if r:
                print(f"  {tbl}.{c}:", r)

con.close()
