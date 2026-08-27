"""
seed_default_avatars.py
DB内の avatar_base64 が NULL または空のアカウント・履歴レコードに、
デフォルト人型SVG Data URIを一括設定するスクリプト
"""
import sqlite3

DB_PATH = "archive.db"
DEFAULT_AVATAR_SVG = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>"

def main():
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()
    
    # 1. accounts テーブルの更新
    cur.execute("""
        UPDATE accounts 
        SET avatar_base64 = ? 
        WHERE avatar_base64 IS NULL OR avatar_base64 = '' OR avatar_base64 NOT LIKE 'data:image/%'
    """, (DEFAULT_AVATAR_SVG,))
    print(f"[+] Updated accounts table ({cur.rowcount} rows).")

    # 2. account_profile_histories テーブルの更新
    cur.execute("""
        UPDATE account_profile_histories 
        SET avatar_base64 = ? 
        WHERE avatar_base64 IS NULL OR avatar_base64 = '' OR avatar_base64 NOT LIKE 'data:image/%'
    """, (DEFAULT_AVATAR_SVG,))
    print(f"[+] Updated account_profile_histories table ({cur.rowcount} rows).")

    conn.commit()
    conn.close()
    print("[+] Done seeding default human avatars!")

if __name__ == "__main__":
    main()
