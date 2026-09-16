import os, shutil, sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

# 1. Move 1st image from _escalate to Msluo14/_assets
src_1 = r"G:\Media_Storage\Influencers\_escalate\GnmCCzebYAAkg9-.jpg"
dst_dir = r"G:\Media_Storage\Influencers\Msluo14\X(Twitter)\_assets"
dst_1 = os.path.join(dst_dir, "GnmCCzebYAAkg9-.jpg")

if os.path.exists(src_1):
    os.makedirs(dst_dir, exist_ok=True)
    shutil.move(src_1, dst_1)
    print(f"Moved {src_1} -> {dst_1}")
else:
    print(f"File already at {dst_1} or not in _escalate")

# 2. Re-link 2nd, 3rd, 4th media to genuine article 1907698917766213981
target_art = "1907698917766213981"
old_art = "1907698913637302272"

c.execute("""
UPDATE media 
SET article_id = ?, download_status = 'COMPLETED', failed_reason = NULL, account_id = '1749477300754878464_msluo14_小罗老师'
WHERE article_id = ? AND media_id LIKE 'GnmCC%'
""", (target_art, old_art))
print("Re-linked rows count:", c.rowcount)

# 3. Update 1st media status to COMPLETED
c.execute("""
UPDATE media 
SET download_status = 'COMPLETED', failed_reason = NULL, account_id = '1749477300754878464_msluo14_小罗老师'
WHERE article_id = ? AND media_id = 'GnmCCzebYAAkg9-.jpg'
""", (target_art,))
print("Updated 1st media count:", c.rowcount)

# 4. Clean up the ghost trashed article
c.execute("DELETE FROM articles WHERE id = ?", (old_art,))
print("Deleted ghost article:", c.rowcount)

conn.commit()

print("\n=== VERIFICATION: MEDIA FOR 1907698917766213981 ===")
c.execute("SELECT media_id, article_id, download_status, download_url FROM media WHERE article_id = ?", (target_art,))
for r in c.fetchall():
    print(" ", r)
