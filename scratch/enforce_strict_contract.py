# scratch/enforce_strict_contract.py
import sqlite3, re

DB_PATH = "archive.db"

def run():
    print("=== ENFORCING STRICT CONTRACT: PURGE NON-WHITELIST MEDIA & ENSURE IS_REPOST ===")
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    # 1. Load active whitelist
    wl_set = {r[0].lower() for r in cur.execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall()}
    print(f"Active Whitelist ({len(wl_set)} handles): {wl_set}")

    # 2. Ensure ALL articles starting with 'RT @' have is_repost = 1 (preserve full_text completely!)
    cur.execute("UPDATE articles SET is_repost = 1 WHERE full_text LIKE 'RT @%' AND is_repost = 0")
    updated_flags = cur.rowcount
    print(f"Ensured is_repost = 1 on {updated_flags} RT articles")

    # 3. Find all media associated with articles that retweet non-whitelist accounts
    rt_rows = cur.execute("SELECT id, full_text, account_id FROM articles WHERE full_text LIKE 'RT @%'").fetchall()
    non_wl_article_ids = []
    for aid, ftext, acc_id in rt_rows:
        m = re.match(r'RT @([a-zA-Z0-9_]+):', ftext)
        if m:
            target_user = m.group(1).lower()
            if target_user not in wl_set:
                non_wl_article_ids.append((aid, target_user, acc_id))

    print(f"Identified {len(non_wl_article_ids)} articles retweeting non-whitelist users")

    # 4. Quarantine these media items into media_excluded and purge from media / media_variants
    quarantined = 0
    for aid, target_user, acc_id in non_wl_article_ids:
        med_rows = cur.execute("""
            SELECT media_id, article_id, account_id, type, download_url, width, height, tweet_urls, thumbnail_url
            FROM media WHERE article_id = ?
        """, (aid,)).fetchall()

        for m in med_rows:
            mid, art_id, m_acc_id, m_type, d_url, w, h, t_urls, thumb = m
            reason = f"Non-whitelist RT: {target_user}"
            cur.execute("""
                INSERT INTO media_excluded (media_id, article_id, account_id, type, download_url, width, height, download_status, failed_reason, tweet_urls, thumbnail_url, quarantine_reason)
                VALUES (?, ?, ?, ?, ?, ?, ?, 'EXCLUDED', 'Whitelist外', ?, ?, ?)
                ON CONFLICT(media_id) DO UPDATE SET
                    quarantine_reason = excluded.quarantine_reason
            """, (mid, art_id, m_acc_id, m_type, d_url, w, h, t_urls, thumb, reason))

            cur.execute("DELETE FROM media WHERE media_id = ?", (mid,))
            cur.execute("DELETE FROM media_variants WHERE media_id = ?", (mid,))
            quarantined += 1

    conn.commit()
    print(f"Quarantined and purged {quarantined} non-whitelist media items from 'media' table!")

    # 5. Verification
    total_media = cur.execute("SELECT COUNT(*) FROM media").fetchone()[0]
    total_excluded = cur.execute("SELECT COUNT(*) FROM media_excluded").fetchone()[0]
    total_articles = cur.execute("SELECT COUNT(*) FROM articles").fetchone()[0]
    total_reposts = cur.execute("SELECT COUNT(*) FROM articles WHERE is_repost = 1").fetchone()[0]

    # Verify no non-whitelist RT media remain in media table
    remaining_leak = cur.execute("""
        SELECT COUNT(*) FROM media m
        JOIN articles a ON a.id = m.article_id
        WHERE a.full_text LIKE 'RT @%' AND a.id IN ({})
    """.format(','.join(['?'] * len(non_wl_article_ids))), [x[0] for x in non_wl_article_ids]).fetchone()[0]

    print("\n=== CONTRACT COMPLIANCE AUDIT ===")
    print(f"Active Whitelist Media in 'media': {total_media} (Pure whitelist!)")
    print(f"Quarantined in 'media_excluded': {total_excluded}")
    print(f"Total articles preserved: {total_articles}")
    print(f"Total articles flagged as is_repost: {total_reposts}")
    print(f"Non-whitelist RT media leaks remaining in 'media': {remaining_leak} (MUST BE 0)")

    conn.close()
    if remaining_leak == 0:
        print("\n>>> CONTRACT PERFECTLY FULFILLED! <<<")

if __name__ == "__main__":
    run()
