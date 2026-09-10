import sqlite3

# 復元対象のアカウントデータ定義
trashed_accounts = [
    {
        "numeric_id": "f35dad07-7066-5405-bb36-7e50e2d94434",
        "username": "yike_luo",
        "display_name": "yike_luo",
        "avatar_url": "",
        "updated_at": "2026-08-30 08:29:56.9903197+09:00",
        "description": "",
        "avatar_base64": "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>",
        "group_name": "yike_luo",
        "alias_of": "ext_yike_luo_yike_luo",
        "is_whitelist": 0,
        "post_count": 0,
        "is_trash": 1,
        "trashed_by": "admin_ui",
        "trash_reason": "名寄せ後処理",
        "trashed_at": "2026-08-28 16:37:31.849353+09:00"
    },
    {
        "numeric_id": "937fdf1b-524d-53b5-8b88-3714ddd275e0",
        "username": "yike2024",
        "display_name": "yike2024",
        "avatar_url": "",
        "updated_at": "2026-08-30 08:29:56.9839804+09:00",
        "description": "",
        "avatar_base64": "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>",
        "group_name": "yike_luo",
        "alias_of": "1753507071260315648_yike2024_三好学生yike",
        "is_whitelist": 0,
        "post_count": 0,
        "is_trash": 1,
        "trashed_by": "admin_ui",
        "trash_reason": "名寄せ後処理",
        "trashed_at": "2026-08-28 19:28:19.8920169+09:00"
    },
    {
        "numeric_id": "24d22341-2a66-5b39-9268-5542c824fcaf",
        "username": "no14_coco",
        "display_name": "no14_coco",
        "avatar_url": "",
        "updated_at": "2026-08-30 08:29:56.9886142+09:00",
        "description": "",
        "avatar_base64": "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>",
        "group_name": "yike_luo",
        "alias_of": "1749477300754878464_no14_coco_英语罗老师",
        "is_whitelist": 0,
        "post_count": 0,
        "is_trash": 1,
        "trashed_by": "admin_ui",
        "trash_reason": "名寄せ後処理",
        "trashed_at": "2026-08-28 19:30:04.0745872+09:00"
    },
    {
        "numeric_id": "73a4ee3d-f852-5097-89fc-e10cf2269d21",
        "username": "subyike",
        "display_name": "subyike",
        "avatar_url": "",
        "updated_at": "2026-09-07 20:44:39.8430771+09:00",
        "description": "",
        "avatar_base64": "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>",
        "group_name": "yike_luo",
        "alias_of": "ext_subyike_YIKE",
        "is_whitelist": 0,
        "post_count": 0,
        "is_trash": 1,
        "trashed_by": "admin_ui",
        "trash_reason": "名寄せ後処理",
        "trashed_at": "2026-09-07 20:50:00.000000+09:00"
    },
    {
        "numeric_id": "bd8c2c8e-5e45-52fb-897b-9026ff76c4b0",
        "username": "nenne1001",
        "display_name": "Nenne1001",
        "avatar_url": "",
        "updated_at": "2026-09-07 09:01:59",
        "description": "",
        "avatar_base64": "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>",
        "group_name": "",
        "alias_of": "1553745049057300481_nenne1001_ねんね",
        "is_whitelist": 0,
        "post_count": 0,
        "is_trash": 1,
        "trashed_by": "admin_ui",
        "trash_reason": "名寄せ後処理",
        "trashed_at": "2026-09-07 09:05:00.000000+09:00"
    }
]

conn = sqlite3.connect("archive.db")
c = conn.cursor()

inserted = 0
for acc in trashed_accounts:
    # 既に存在しているかチェック
    c.execute("SELECT 1 FROM accounts WHERE numeric_id = ?", (acc["numeric_id"],))
    if c.fetchone():
        print(f"Account {acc['numeric_id']} already exists, updating trash status...")
        c.execute("""
            UPDATE accounts
            SET is_trash = 1, trash_reason = ?, trashed_by = ?, trashed_at = ?, alias_of = ?
            WHERE numeric_id = ?
        """, (acc["trash_reason"], acc["trashed_by"], acc["trashed_at"], acc["alias_of"], acc["numeric_id"]))
    else:
        print(f"Restoring trashed account: {acc['username']} ({acc['numeric_id']}) -> {acc['trash_reason']}")
        c.execute("""
            INSERT INTO accounts (
                numeric_id, username, display_name, avatar_url, updated_at,
                description, avatar_base64, group_name, alias_of, is_whitelist,
                post_count, is_trash, trashed_by, trash_reason, trashed_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            acc["numeric_id"], acc["username"], acc["display_name"], acc["avatar_url"], acc["updated_at"],
            acc["description"], acc["avatar_base64"], acc["group_name"], acc["alias_of"], acc["is_whitelist"],
            acc["post_count"], acc["is_trash"], acc["trashed_by"], acc["trash_reason"], acc["trashed_at"]
        ))
        inserted += 1

conn.commit()
print(f"Restoration complete! Inserted {inserted} trashed accounts.")

# 検証
c.execute("SELECT numeric_id, username, is_trash, trash_reason, alias_of FROM accounts WHERE is_trash = 1")
print("Current trashed accounts in archive.db:")
for r in c.fetchall():
    print(" ", r)
