import os, sys, shutil, sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

escalate_dir = r"G:\Media_Storage\Influencers\_escalate"
dest_root = r"G:\Media_Storage\Influencers"

if not os.path.exists(escalate_dir):
    print("No _escalate dir found")
    sys.exit(0)

files = os.listdir(escalate_dir)
print(f"Files in _escalate: {len(files)}")

for f in files:
    src = os.path.join(escalate_dir, f)
    if os.path.isdir(src): continue
    
    base = f
    clean_id = os.path.splitext(base)[0]
    for sfx in ["_orig", "_large", "_wayback_orig", "_wayback", "_plain"]:
        if clean_id.endswith(sfx): clean_id = clean_id[:-len(sfx)]
        
    # Find owner in DB
    sql = """
    SELECT accounts.username, media.media_id, media.article_id
    FROM media
    JOIN articles ON articles.id = media.article_id
    JOIN accounts ON (accounts.numeric_id = articles.account_id OR accounts.username = articles.account_id)
    WHERE media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?
    LIMIT 1
    """
    c.execute(sql, (base, clean_id, clean_id + "%"))
    row = c.fetchone()
    
    if row and row[0]:
        owner = row[0]
        med_id = row[1]
        target_dir = os.path.join(dest_root, owner, "X(Twitter)", "_assets")
        os.makedirs(target_dir, exist_ok=True)
        dst = os.path.join(target_dir, base)
        shutil.move(src, dst)
        
        # Update media status to COMPLETED
        c.execute("UPDATE media SET download_status = 'COMPLETED', failed_reason = NULL WHERE media_id = ?", (med_id,))
        print(f"✅ Moved & COMPLETED: {base} -> {owner} (Media: {med_id})")
    else:
        print(f"⚠️ Owner not found in DB for: {base}")

conn.commit()
