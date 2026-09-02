import sqlite3

def migrate_and_quarantine(db_path="archive.db"):
    conn = sqlite3.connect(db_path, timeout=60.0)
    conn.execute("PRAGMA journal_mode = WAL;")
    conn.execute("PRAGMA busy_timeout = 60000;")
    cur = conn.cursor()
    print(f"[*] Starting migration on {db_path}...")

    # 1. media_excluded テーブルの作成（ホワイトリスト外・退避用）
    cur.execute("""
    CREATE TABLE IF NOT EXISTS media_excluded (
        media_id TEXT PRIMARY KEY,
        article_id TEXT NOT NULL,
        account_id TEXT,
        type TEXT NOT NULL,
        download_url TEXT NOT NULL,
        width INTEGER,
        height INTEGER,
        download_status TEXT,
        failed_reason TEXT,
        tweet_urls TEXT,
        thumbnail_url TEXT,
        quarantined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        quarantine_reason TEXT
    );
    """)

    # 2. media テーブルに account_id カラム追加
    med_cols = {r[1] for r in cur.execute("PRAGMA table_info(media)").fetchall()}
    if "account_id" not in med_cols:
        cur.execute("ALTER TABLE media ADD COLUMN account_id TEXT")
        print("  [+] Added account_id column to media table")

    # 3. 既存 media の account_id を articles からバックフィル
    cur.execute("""
    UPDATE media 
    SET account_id = (
        SELECT articles.account_id 
        FROM articles 
        WHERE articles.id = media.article_id
    )
    WHERE (account_id IS NULL OR account_id = '')
      AND EXISTS (SELECT 1 FROM articles WHERE articles.id = media.article_id);
    """)
    updated_cnt = cur.rowcount
    print(f"  [+] Backfilled account_id for {updated_cnt} media records")

    # 4. ホワイトリスト取得
    cur.execute("SELECT value FROM whitelists WHERE is_active = 1")
    wl_usernames = {r[0].lower() for r in cur.fetchall() if r[0]}
    print(f"  [*] Active whitelist: {wl_usernames}")

    # 5. ホワイトリスト外アカウントのメディアを特定して media_excluded に退避
    # (accounts.username が whitelist にないもの、または EXCLUDED ステータスのもの)
    cur.execute("""
    SELECT m.media_id, m.article_id, m.account_id, m.type, m.download_url, 
           m.width, m.height, m.download_status, m.failed_reason, m.tweet_urls, m.thumbnail_url,
           ac.username
    FROM media m
    LEFT JOIN accounts ac ON (ac.numeric_id = m.account_id OR ac.username = m.account_id)
    """)
    rows = cur.fetchall()
    quarantined_count = 0

    for r in rows:
        m_id, art_id, acc_id, m_type, dl_url, w, h, st, fr, tu, thu, uname = r
        is_wl = uname and uname.lower() in wl_usernames
        if not is_wl and (st == "EXCLUDED" or uname is None or uname.lower() not in wl_usernames):
            cur.execute("""
            INSERT OR REPLACE INTO media_excluded 
            (media_id, article_id, account_id, type, download_url, width, height, download_status, failed_reason, tweet_urls, thumbnail_url, quarantine_reason)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """, (m_id, art_id, acc_id, m_type, dl_url, w, h, st, fr, tu, thu, f"Non-whitelist account: {uname}"))
            cur.execute("DELETE FROM media WHERE media_id = ?", (m_id,))
            quarantined_count += 1

    print(f"  [+] Quarantined {quarantined_count} non-whitelist media records to media_excluded")

    # 6. インデックス作成
    cur.execute("CREATE INDEX IF NOT EXISTS idx_media_account ON media(account_id);")
    conn.commit()

    # 7. 検証確認
    total_media = cur.execute("SELECT COUNT(*) FROM media").fetchone()[0]
    with_acc = cur.execute("SELECT COUNT(*) FROM media WHERE account_id IS NOT NULL AND account_id != ''").fetchone()[0]
    total_excluded = cur.execute("SELECT COUNT(*) FROM media_excluded").fetchone()[0]
    conn.close()

    print(f"[*] Done. Active Media: {total_media} (with account_id: {with_acc}), Excluded Media: {total_excluded}")

if __name__ == "__main__":
    migrate_and_quarantine()
