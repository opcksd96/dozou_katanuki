# scratch/reconcile_and_repair_identities.py
import os, re, sqlite3, time

DB_PATH = "archive.db"
BACKUP_PATH = "backups/database/archive_20260831_082956.db"

def clean_tag(s: str) -> str:
    cleaned = re.sub(r'[^\w\u4e00-\u9fa5\u3040-\u30ff]', '', str(s or ""))
    return re.sub(r'[\s_]+', '_', cleaned).strip('_')

def run_migration():
    print("=== STARTING TRIPLET IDENTITY MIGRATION & RECONCILIATION ===")
    conn = sqlite3.connect(DB_PATH)
    conn.execute("PRAGMA foreign_keys = OFF;") # Temporarily disable during ID updates
    cur = conn.cursor()

    # 1. Inspect current accounts
    cur.execute("SELECT numeric_id, username, display_name, avatar_url, description, avatar_base64, alias_of, is_whitelist FROM accounts")
    raw_accounts = cur.fetchall()
    print(f"Loaded {len(raw_accounts)} raw accounts from {DB_PATH}")

    # Canonical mapping definitions
    # Map old ID -> New Triplet ID
    id_map = {}
    new_accounts = [] # (new_id, username, display_name, avatar_url, desc, av_b64, alias_of, is_wl)

    # Specific well-known accounts
    well_known = {
        "msluo14": {"nid": "1749477300754878464", "u": "msluo14", "d": "小罗老师", "alias_of": ""},
        "no14_coco": {"nid": "1749477300754878464", "u": "no14_coco", "d": "英语罗老师", "alias_of": "1749477300754878464_msluo14_小罗老师"},
        "subyike": {"nid": "ext", "u": "subyike", "d": "YIKE", "alias_of": "1749477300754878464_msluo14_小罗老师"},
        "yike_luo": {"nid": "ext", "u": "yike_luo", "d": "yike_luo", "alias_of": "1749477300754878464_msluo14_小罗老师"},
        "yike2024": {"nid": "1753507071260315648", "u": "yike2024", "d": "三好学生yike", "alias_of": "1749477300754878464_msluo14_小罗老师"},
        "yike233_": {"nid": "726272140656390144", "u": "yike233_", "d": "可可小可爱", "alias_of": "1749477300754878464_msluo14_小罗老师"},
        "sayapom4": {"nid": "ext", "u": "sayapom4", "d": "さやか", "alias_of": ""},
        "nenne1001": {"nid": "1553745049057300481", "u": "nenne1001", "d": "ねんね", "alias_of": ""},
        "yike1416431": {"nid": "ext", "u": "yike1416431", "d": "半只桃", "alias_of": ""}
    }

    # Generate canonical target ID for each account
    created_new_accounts = {}
    for r in raw_accounts:
        old_id, u_name, d_name, av_url, desc, av_b64, alias_of, is_wl = r
        u_clean = u_name.strip().lower()

        if u_clean in well_known:
            info = well_known[u_clean]
            nid = info["nid"]
            target_u = info["u"]
            target_d = info["d"]
            target_alias = info["alias_of"]
        else:
            nid = old_id if (old_id and old_id.isdigit()) else "ext"
            target_u = u_clean
            target_d = clean_tag(d_name) or target_u
            target_alias = alias_of or ""

        new_id = f"{nid}_{target_u}_{target_d}"
        id_map[old_id] = new_id

        if new_id not in created_new_accounts:
            created_new_accounts[new_id] = {
                "numeric_id": new_id,
                "username": target_u,
                "display_name": d_name if d_name else target_d,
                "avatar_url": av_url or "",
                "description": desc or "",
                "avatar_base64": av_b64 or "",
                "alias_of": target_alias,
                "is_whitelist": is_wl
            }
        else:
            # merge info
            if av_url and not created_new_accounts[new_id]["avatar_url"]:
                created_new_accounts[new_id]["avatar_url"] = av_url
            if av_b64 and not created_new_accounts[new_id]["avatar_base64"]:
                created_new_accounts[new_id]["avatar_base64"] = av_b64
            if is_wl:
                created_new_accounts[new_id]["is_whitelist"] = 1

    print(f"Mapped {len(id_map)} legacy IDs to {len(created_new_accounts)} canonical triplet IDs")

    # Step 2: Separate hijacked articles under 1749477300754878464
    msluo_triplet = "1749477300754878464_msluo14_小罗老师"
    coco_triplet = "1749477300754878464_no14_coco_英语罗老师"

    # Identify msluo articles currently under 1749477300754878464 or 24d22341-2a66-5b39-9268-5542c824fcaf
    hijacked_to_msluo = []
    legit_coco = []
    for row in cur.execute("SELECT id, original_url, wayback_url, sotwe_url, nitter_url, twistalker_url, account_id FROM articles WHERE account_id IN ('1749477300754878464', '24d22341-2a66-5b39-9268-5542c824fcaf')").fetchall():
        aid = row[0]
        urls_str = ' '.join([str(u or '') for u in row[1:6]]).lower()
        if "msluo14" in urls_str:
            hijacked_to_msluo.append(aid)
        else:
            legit_coco.append(aid)

    print(f"Article separation: {len(hijacked_to_msluo)} articles -> msluo14, {len(legit_coco)} articles -> no14_coco")

    # Step 3: Update `articles` and `media`
    # First, handle all standard mappings from id_map
    for old_id, new_id in id_map.items():
        if old_id == "1749477300754878464":
            continue # handled specially above
        cur.execute("UPDATE articles SET account_id = ? WHERE account_id = ?", (new_id, old_id))
        cur.execute("UPDATE media SET account_id = ? WHERE account_id = ?", (new_id, old_id))
        cur.execute("UPDATE account_profile_histories SET account_id = ? WHERE account_id = ?", (new_id, old_id))
        cur.execute("UPDATE media_excluded SET account_id = ? WHERE account_id = ?", (new_id, old_id))

    # Now apply the separated msluo vs coco articles & media
    for aid in hijacked_to_msluo:
        cur.execute("UPDATE articles SET account_id = ? WHERE id = ?", (msluo_triplet, aid))
        cur.execute("UPDATE media SET account_id = ? WHERE article_id = ?", (msluo_triplet, aid))

    for aid in legit_coco:
        cur.execute("UPDATE articles SET account_id = ? WHERE id = ?", (coco_triplet, aid))
        cur.execute("UPDATE media SET account_id = ? WHERE article_id = ?", (coco_triplet, aid))

    # Any remaining articles under 1749477300754878464
    cur.execute("UPDATE articles SET account_id = ? WHERE account_id = '1749477300754878464'", (coco_triplet,))
    cur.execute("UPDATE media SET account_id = ? WHERE account_id = '1749477300754878464'", (coco_triplet,))

    # Step 4: Recreate `accounts` table with canonical triplet IDs
    # Clear old accounts and insert canonical accounts
    cur.execute("DELETE FROM accounts")
    now_ts = time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())
    for acc in created_new_accounts.values():
        cur.execute("""
            INSERT INTO accounts (numeric_id, username, display_name, avatar_url, description, avatar_base64, alias_of, is_whitelist, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            acc["numeric_id"], acc["username"], acc["display_name"],
            acc["avatar_url"], acc["description"], acc["avatar_base64"],
            acc["alias_of"], acc["is_whitelist"], now_ts
        ))

    # Synchronize post counts
    for acc_id in created_new_accounts.keys():
        p_cnt = cur.execute("SELECT COUNT(*) FROM articles WHERE account_id = ?", (acc_id,)).fetchone()[0]
        cur.execute("UPDATE accounts SET post_count = ? WHERE numeric_id = ?", (p_cnt, acc_id))

    conn.commit()
    conn.execute("PRAGMA foreign_keys = ON;")

    # Step 5: Verify results
    print("\n=== VERIFICATION ===")
    fk_errors = cur.execute("PRAGMA foreign_key_check;").fetchall()
    if fk_errors:
        print("WARNING: Foreign key errors detected:", fk_errors)
    else:
        print("SUCCESS: 0 Foreign Key errors!")

    print("\nCanonical Accounts and post/media counts:")
    for r in cur.execute("SELECT numeric_id, username, display_name, alias_of, post_count FROM accounts ORDER BY post_count DESC").fetchall():
        m_cnt = cur.execute("SELECT COUNT(*) FROM media WHERE account_id = ?", (r[0],)).fetchone()[0]
        print(f"  {r[0]} | @{r[1]} ({r[2]}) | posts: {r[4]}, media: {m_cnt} | alias_of: '{r[3]}'")

    conn.close()
    print("\n=== MIGRATION COMPLETED SUCCESSFULLY ===")

if __name__ == "__main__":
    run_migration()
