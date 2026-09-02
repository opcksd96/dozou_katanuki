import os, re, sqlite3

files = [
    "EBzLD6hjVRyTXvAb_0.mp4",
    "EwJ7fsA67FadxvXU_0_0.mp4",
    "fEoWjbe4PHI11zGS.mp4",
    "Ga4tQg0a8AA0tQt_0.mp4",
    "GEFsflnacAAOXck.jpg",
    "Gf9LCaIaQAAMsSV(5).jpg",
    "GNioajjaUAArFS6.jpg",
    "GNioajjaUAArFS6(5).jpg",
    "GZ2bywSMlf6B81DA_0.mp4",
    "P5HCLz4CxtcJZRkV.mp4",
    "rlCDbpAgOdJjPh_p.mp4"
]

con = sqlite3.connect("archive.db")
cur = con.cursor()

def clean_name(fn):
    base = os.path.splitext(fn)[0]
    # Remove (1), (2), (5) etc
    base = re.sub(r'\(\d+\)$', '', base)
    # Remove _0, _0_0 etc
    base = re.sub(r'(_0)+$', '', base)
    # Remove resolution suffixes
    for sfx in ["_motrix", "_requests", "_thunder", "_plain", "_orig", "_large", "_wayback_orig", "_wayback_plain", "_wayback"]:
        if base.endswith(sfx):
            base = base[:-len(sfx)]
            break
    return base

print(f"{'Filename':<25} | {'CleanID':<20} | {'Media Match':<20} | {'Account Match'}")
print("-" * 85)

for fn in files:
    cid = clean_name(fn)
    # Media search
    m_row = cur.execute("SELECT media_id, account_id, article_id FROM media WHERE media_id LIKE ? OR download_url LIKE ?", (cid + "%", "%" + cid + "%")).fetchall()
    
    # Check accounts
    acc_name = "NONE"
    if m_row:
        m_id, acc_id, art_id = m_row[0]
        a_row = cur.execute("SELECT username FROM accounts WHERE numeric_id = ? OR username = ?", (acc_id, acc_id)).fetchone()
        acc_name = a_row[0] if a_row else f"UNKNOWN_ACC({acc_id})"
    else:
        # Check articles or tweet_urls directly
        a_row = cur.execute("SELECT a.account_id, ac.username FROM articles a JOIN accounts ac ON (ac.numeric_id = a.account_id OR ac.username = a.account_id) WHERE a.full_text LIKE ? OR a.wayback_url LIKE ? OR a.original_url LIKE ?", (f"%{cid}%", f"%{cid}%", f"%{cid}%")).fetchone()
        if a_row:
            acc_name = f"via article: {a_row[1]}"

    print(f"{fn:<25} | {cid:<20} | {len(m_row):<20} | {acc_name}")

con.close()
