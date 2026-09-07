import sqlite3, json, os, shutil, glob
from warcio.archiveiterator import ArchiveIterator

DB_PATH = 'archive.db'
DUMPS_DIR = 'backups/dumps/twitter'

conn = sqlite3.connect(DB_PATH)
conn.execute("PRAGMA foreign_keys = ON;")
cur = conn.cursor()

print("=== STEP 1: SANITIZE MSLUO14 POSTS ===")

# Definition of the 8 fake accounts
fake_account_defs = [
    {"numeric_id": "2039101842412490752", "username": "MsLuo14d", "display_name": "MsLuo14|MsLuo14", "post_ids": ["2067940688831140183"]},
    {"numeric_id": "2039603066295484416", "username": "MsLuo14b", "display_name": "MsLuo14|MsLuo14", "post_ids": ["2066169076347547704"]},
    {"numeric_id": "2039493607183249408", "username": "MsLuo1433", "display_name": "MsLuo14|MsLuo14", "post_ids": ["2065000003857190959"]},
    {"numeric_id": "1968899858187837440", "username": "MsLuo14hh", "display_name": "MsLuo14@MsLuo14", "post_ids": ["2003915292322550072"]},
    {"numeric_id": "1971766511225475072", "username": "MsLuo14kk", "display_name": "Youngest NFT Artist.", "post_ids": ["1973109291830641052"]},
    {"numeric_id": "1952056111953063936", "username": "MsLuo1422", "display_name": "Carol Andrews", "post_ids": ["1952489134834012473"]},
    {"numeric_id": "1932688313250660352", "username": "MsLuo14_", "display_name": "MsLuo14_MsLuo14", "post_ids": ["1943288913461686461"]},
    {"numeric_id": "1909575892352790528", "username": "MsLuo141", "display_name": "小罗老师_MsLuo14", "post_ids": ["1925870551215567172"]},
]

# Other 2 posts that belonged to existing legitimate accounts
existing_reassign = [
    {"post_id": "1850548533390385289", "target_numeric_id": "1749477300754878464", "target_username": "no14_coco"},
    {"post_id": "1744989650513789153", "target_numeric_id": "2cda0ead-b020-efbb-be02-6faba155be4b", "target_username": "yike_luo"}
]

# 1. Register fake accounts in `accounts` table
for fa in fake_account_defs:
    cur.execute("""
        INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at, is_whitelist, group_name, alias_of)
        VALUES (?, ?, ?, '', CURRENT_TIMESTAMP, 0, '', '')
        ON CONFLICT(numeric_id) DO UPDATE SET
            username = excluded.username,
            display_name = excluded.display_name
    """, (fa["numeric_id"], fa["username"], fa["display_name"]))
    print(f"Registered/Verified fake account: {fa['username']} ({fa['numeric_id']})")

# 2. Reassign articles and media to fake accounts
for fa in fake_account_defs:
    for pid in fa["post_ids"]:
        cur.execute("UPDATE articles SET account_id = ? WHERE id = ?", (fa["numeric_id"], pid))
        cur.execute("UPDATE media SET account_id = ? WHERE article_id = ?", (fa["numeric_id"], pid))
        print(f"Reassigned post {pid} -> account {fa['username']} ({fa['numeric_id']})")

        # Move WARC in dumps pool if exists
        src_dir = os.path.join(DUMPS_DIR, "msluo14", pid)
        dst_dir = os.path.join(DUMPS_DIR, fa["username"], pid)
        if os.path.exists(src_dir):
            os.makedirs(os.path.dirname(dst_dir), exist_ok=True)
            if os.path.exists(dst_dir):
                shutil.rmtree(dst_dir)
            shutil.move(src_dir, dst_dir)
            print(f"Moved WARC pool: {src_dir} -> {dst_dir}")

# 3. Reassign existing accounts' posts
for ea in existing_reassign:
    pid = ea["post_id"]
    cur.execute("UPDATE articles SET account_id = ? WHERE id = ?", (ea["target_numeric_id"], pid))
    cur.execute("UPDATE media SET account_id = ? WHERE article_id = ?", (ea["target_numeric_id"], pid))
    print(f"Reassigned post {pid} -> legitimate account {ea['target_username']} ({ea['target_numeric_id']})")

    src_dir = os.path.join(DUMPS_DIR, "msluo14", pid)
    dst_dir = os.path.join(DUMPS_DIR, ea["target_username"], pid)
    if os.path.exists(src_dir):
        os.makedirs(os.path.dirname(dst_dir), exist_ok=True)
        if os.path.exists(dst_dir):
            shutil.rmtree(dst_dir)
        shutil.move(src_dir, dst_dir)
        print(f"Moved WARC pool: {src_dir} -> {dst_dir}")

print("\n=== STEP 2: NENNE1001 ALIAS_OF NORMALIZATION ===")
# Ensure TwStalker account exists for nenne1001 with alias_of = 'nenne1001'
twstalker_uuid = "bd8c2c8e-5e45-52fb-897b-9026ff76c4b0"
cur.execute("""
    INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at, is_whitelist, alias_of)
    VALUES (?, 'nenne1001', 'ねんね', '', CURRENT_TIMESTAMP, 0, 'nenne1001')
    ON CONFLICT(numeric_id) DO UPDATE SET
        username = 'nenne1001',
        alias_of = 'nenne1001'
""", (twstalker_uuid,))
print(f"Registered TwStalker account for nenne1001 with alias_of='nenne1001': {twstalker_uuid}")

# Reassign TwStalker articles to twstalker_uuid
cur.execute("UPDATE articles SET account_id = ? WHERE via = 'TwStalker' AND (original_url LIKE '%nenne1001%' OR twistalker_url LIKE '%nenne1001%')", (twstalker_uuid,))
cur.execute("UPDATE media SET account_id = ? WHERE article_id IN (SELECT id FROM articles WHERE account_id = ?)", (twstalker_uuid, twstalker_uuid))
print(f"Reassigned TwStalker articles for nenne1001 to {twstalker_uuid}")

conn.commit()

# Check foreign keys
fk_errors = cur.execute("PRAGMA foreign_key_check;").fetchall()
print("\n=== STEP 3: AUDIT & INTEGRITY CHECK ===")
if fk_errors:
    print("WARNING: Foreign key errors detected:", fk_errors)
else:
    print("SUCCESS: PRAGMA foreign_key_check passed with 0 errors!")

# Verify msluo14 count
cnt_msluo = cur.execute("SELECT COUNT(*) FROM articles WHERE account_id = '4fa935dc-eda5-4217-4b5b-c9c9ea0fb491'").fetchone()[0]
print(f"Cleaned msluo14 article count: {cnt_msluo} (Expected: 163)")

# Verify nenne counts
cnt_nenne_official = cur.execute("SELECT COUNT(*) FROM articles WHERE account_id = '1553745049057300481'").fetchone()[0]
cnt_nenne_twstalker = cur.execute("SELECT COUNT(*) FROM articles WHERE account_id = 'bd8c2c8e-5e45-52fb-897b-9026ff76c4b0'").fetchone()[0]
print(f"Nenne official articles: {cnt_nenne_official}, TwStalker articles: {cnt_nenne_twstalker}")

conn.close()
