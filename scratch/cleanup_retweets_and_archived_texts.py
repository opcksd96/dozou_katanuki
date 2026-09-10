# scratch/cleanup_retweets_and_archived_texts.py
import sqlite3, re

DB_PATH = "archive.db"
BACKUP_PATH = "backups/database/archive_20260831_082956.db"

def run_cleanup():
    print("=== STARTING CLEANUP: RESTORE ARCHIVED TEXTS & PURGE NON-WHITELIST RT MEDIA ===")
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()
    b_conn = sqlite3.connect(BACKUP_PATH)
    b_cur = b_conn.cursor()

    # Part 1: Restore full_text from backup
    archived_posts = cur.execute("SELECT id, full_text FROM articles WHERE full_text LIKE 'Archived post%'").fetchall()
    print(f"Found {len(archived_posts)} articles with 'Archived post...' dummy text")

    restored_count = 0
    cleared_count = 0
    for aid, dummy_txt in archived_posts:
        b_row = b_cur.execute("SELECT full_text, full_text_ja, full_text_en, full_text_zh, lang FROM articles WHERE id = ?", (aid,)).fetchone()
        if b_row and b_row[0] and not b_row[0].startswith("Archived post"):
            cur.execute("""
                UPDATE articles
                SET full_text = ?, full_text_ja = ?, full_text_en = ?, full_text_zh = ?, lang = ?
                WHERE id = ?
            """, (b_row[0], b_row[1], b_row[2], b_row[3], b_row[4], aid))
            restored_count += 1
        else:
            # Clear meaningless placeholder
            cur.execute("UPDATE articles SET full_text = '', full_text_ja = '' WHERE id = ?", (aid,))
            cleared_count += 1

    print(f"Text repair: {restored_count} articles restored from backup, {cleared_count} cleared to empty string")

    # Part 2: Purge non-whitelist retweet media into media_excluded
    wl_set = {r[0].lower() for r in cur.execute("SELECT value FROM whitelists WHERE is_active = 1").fetchall()}
    print(f"Active Whitelist: {wl_set}")

    # Find articles with RT @...
    rt_rows = cur.execute("SELECT id, full_text, account_id FROM articles WHERE full_text LIKE 'RT @%'").fetchall()
    print(f"Found {len(rt_rows)} RT articles in database")

    non_wl_rt_article_ids = []
    for aid, ftext, acc_id in rt_rows:
        m = re.match(r'RT @([a-zA-Z0-9_]+):', ftext)
        if m:
            target_user = m.group(1).lower()
            if target_user not in wl_set:
                non_wl_rt_article_ids.append((aid, target_user, acc_id))

    print(f"Found {len(non_wl_rt_article_ids)} articles retweeting non-whitelist users")

    # Move media belonging to these articles to media_excluded
    quarantined_media = 0
    for aid, target_user, acc_id in non_wl_rt_article_ids:
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

            # Delete from media & media_variants
            cur.execute("DELETE FROM media WHERE media_id = ?", (mid,))
            cur.execute("DELETE FROM media_variants WHERE media_id = ?", (mid,))
            quarantined_media += 1

    conn.commit()
    print(f"Successfully quarantined {quarantined_media} media items from non-whitelist retweets!")

    # Part 3: Verify counts
    remaining_archived = cur.execute("SELECT COUNT(*) FROM articles WHERE full_text LIKE 'Archived post%'").fetchone()[0]
    total_media = cur.execute("SELECT COUNT(*) FROM media").fetchone()[0]
    total_excluded = cur.execute("SELECT COUNT(*) FROM media_excluded").fetchone()[0]
    total_articles = cur.execute("SELECT COUNT(*) FROM articles").fetchone()[0]

    print("\n=== VERIFICATION ===")
    print(f"Remaining 'Archived post%' articles: {remaining_archived}")
    print(f"Active media in 'media': {total_media}")
    print(f"Quarantined in 'media_excluded': {total_excluded}")
    print(f"Total articles: {total_articles}")

    conn.close()
    b_conn.close()
    print("\n=== CLEANUP FINISHED SUCCESSFULLY ===")

if __name__ == "__main__":
    run_cleanup()
