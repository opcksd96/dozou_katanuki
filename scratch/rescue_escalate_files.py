import os, re, shutil, sqlite3

esc_dir = r"G:\Media_Storage\Influencers\_escalate"
root_dir = r"G:\Media_Storage\Influencers"

files = os.listdir(esc_dir) if os.path.exists(esc_dir) else []
print(f"[*] Files in _escalate ({len(files)}): {files}")

con = sqlite3.connect("archive.db")
cur = con.cursor()

def clean_name(fn):
    base = os.path.splitext(fn)[0]
    base = re.sub(r'\(\d+\)$', '', base).strip()
    base = re.sub(r'(_[01])+$', '', base)
    for sfx in ["_motrix", "_requests", "_thunder", "_plain", "_orig", "_large", "_wayback_orig", "_wayback_plain", "_wayback"]:
        if base.endswith(sfx):
            base = base[:-len(sfx)]
            break
    return base

rescued = 0
for fn in files:
    src_path = os.path.join(esc_dir, fn)
    ext = os.path.splitext(fn)[1].lower()
    cid = clean_name(fn)
    
    # 1. media.account_id 直接逆引き
    row = cur.execute("""
        SELECT ac.username FROM media m
        JOIN accounts ac ON (ac.numeric_id = m.account_id OR ac.username = m.account_id)
        WHERE m.media_id = ? OR m.media_id = ? OR m.media_id LIKE ? OR m.download_url LIKE ?
        LIMIT 1
    """, (fn, cid + ext, cid + "%", "%" + cid + "%")).fetchone()
    
    # 2. articles 経由
    if not row:
        row = cur.execute("""
            SELECT ac.username FROM articles a
            JOIN accounts ac ON (ac.numeric_id = a.account_id OR ac.username = a.account_id)
            WHERE a.full_text LIKE ? OR a.wayback_url LIKE ? OR a.original_url LIKE ?
            LIMIT 1
        """, (f"%{cid}%", f"%{cid}%", f"%{cid}%")).fetchone()

    if row and row[0]:
        owner = row[0]
        dest_dir = os.path.join(root_dir, owner, "X(Twitter)", "_assets")
        os.makedirs(dest_dir, exist_ok=True)
        # 保存先ファイル名はクリーンなメディア名 + 拡張子
        dest_fn = cid + ext
        dest_path = os.path.join(dest_dir, dest_fn)
        
        shutil.copy2(src_path, dest_path)
        os.remove(src_path)
        print(f"  [+] Rescued: {fn} -> {owner}/X(Twitter)/_assets/{dest_fn}")
        rescued += 1
    else:
        print(f"  [-] Unresolved: {fn} (remains in _escalate for now)")

con.close()
print(f"[*] Rescued {rescued} files from _escalate.")
