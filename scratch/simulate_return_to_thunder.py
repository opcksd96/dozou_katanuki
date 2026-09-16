import os, sqlite3, shutil

conn = sqlite3.connect('archive.db')
c = conn.cursor()

media_id = "GnmCCzebYAAkg9-"
clean_id = media_id
name = "GnmCCzebYAAkg9-.jpg"

# Test 1: GetMediaOwnerUsername logic
sql = """
SELECT accounts.username FROM media
JOIN articles ON articles.id = media.article_id
JOIN accounts ON (accounts.numeric_id = articles.account_id OR accounts.username = articles.account_id)
WHERE media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?
LIMIT 1
"""
c.execute(sql, (media_id, clean_id, clean_id + "%"))
row = c.fetchone()
print("Owner check result:", row)
owner = row[0] if row else ""

# Test 2: Destination path
dest_root = r"G:\Media_Storage\Influencers"
if owner:
    target_sub_dir = os.path.join(dest_root, owner, "X(Twitter)", "_assets")
else:
    target_sub_dir = os.path.join(dest_root, "_escalate")
dest_path = os.path.join(target_sub_dir, name)
print("Target destination path:", dest_path)

# Test 3: RegisterCompletedMediaFile matching logic
base = name
clean_id = os.path.splitext(base)[0]
sql_update = """
SELECT media_id, download_status FROM media
WHERE media_id = ? OR media_id = ? OR download_url LIKE ?
"""
c.execute(sql_update, (base, clean_id, "%/" + base))
matches = c.fetchall()
print("Media rows that will be updated to COMPLETED:", matches)
