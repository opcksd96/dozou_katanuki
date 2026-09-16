import sqlite3, os

conn = sqlite3.connect('archive.db')
c = conn.cursor()

files = [
    "EBzLD6hjVRyTXvAb.mp4",
    "Esudn4NjbxCp8vKe.mp4",
    "EwJ7fsA67FadxvXU.mp4",
    "fEoWjbe4PHI11zGS.mp4",
    "FWpuTj1akAAiRMv.jpg",
    "Ga4tQg0a8AA0tQt.mp4",
    "GEFsflnacAAOXck.jpg",
    "Gf9LCaIaQAAMsSV.jpg",
    "GNioajjaUAArFS6.jpg",
    "GnmCCzebYAAkg9-.jpg",
    "GZ2bywSMlf6B81DA.mp4",
    "P5HCLz4CxtcJZRkV.mp4",
    "pXswPdtFOhov-oMj.mp4",
    "rlCDbpAgOdJjPh_p.mp4",
    "tNlvehz60jPqgxdl.mp4",
    "x7vQf7QsgEONy7TF.mp4"
]

print("=== CHECKING 16 ESCALATED FILES ===")
for f in files:
    base = f
    ext = os.path.splitext(base)[1]
    clean_id = os.path.splitext(base)[0]
    for sfx in ["_orig", "_large", "_wayback_orig", "_wayback"]:
        if clean_id.endswith(sfx):
            clean_id = clean_id[:-len(sfx)]
            
    # Check media table
    c.execute("SELECT media_id, article_id, account_id FROM media WHERE media_id = ? OR media_id = ? OR media_id LIKE ?", (base, clean_id, clean_id + "%"))
    m_rows = c.fetchall()
    
    # Check full join
    sql = """
    SELECT accounts.username, articles.id, media.media_id, articles.account_id
    FROM media
    JOIN articles ON articles.id = media.article_id
    JOIN accounts ON (accounts.numeric_id = articles.account_id OR accounts.username = articles.account_id)
    WHERE media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?
    """
    c.execute(sql, (base, clean_id, clean_id + "%"))
    join_rows = c.fetchall()
    
    print(f"\nFile: {f}")
    print(f"  Media found: {m_rows}")
    print(f"  Join found: {join_rows}")
