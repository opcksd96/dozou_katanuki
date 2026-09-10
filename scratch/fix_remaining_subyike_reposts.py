# scratch/fix_remaining_subyike_reposts.py
import sqlite3, datetime

DB_PATH = 'archive.db'

def main():
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    # 1. Ensure ext_codmanconder exists in accounts
    cur.execute("SELECT numeric_id FROM accounts WHERE numeric_id = 'ext_codmanconder'")
    if not cur.fetchone():
        now = datetime.datetime.now().isoformat()
        cur.execute("""
            INSERT INTO accounts (numeric_id, username, display_name, is_whitelist, is_trash, created_at, updated_at)
            VALUES ('ext_codmanconder', 'codmanconder', '宅舞COS舞JK舞', 0, 0, ?, ?)
        """, (now, now))
        print("Created account: ext_codmanconder")

    # 2. Fix post 2028459277254602771
    cur.execute("""
        UPDATE articles
        SET account_id = 'ext_codmanconder',
            full_text = '以安 定制 顶跨  https://t.co/Bec0U784wY https://t.co/w4XcLZNeia'
        WHERE id = '2028459277254602771'
    """)
    print("Updated 2028459277254602771 -> ext_codmanconder")

    # Update RT 2028816874549707202 reply_to_id
    cur.execute("""
        UPDATE articles
        SET reply_to_id = '2028459277254602771'
        WHERE id = '2028816874549707202'
    """)
    print("Updated RT 2028816874549707202 -> reply_to_id 2028459277254602771")

    # Add media_excluded for 2028459277254602771
    cur.execute("SELECT media_id FROM media_excluded WHERE article_id = '2028459277254602771'")
    if not cur.fetchone():
        now = datetime.datetime.now().isoformat()
        cur.execute("""
            INSERT INTO media_excluded (media_id, article_id, account_id, type, download_url, quarantined_at, quarantine_reason)
            VALUES ('me_2028459277254602771_0', '2028459277254602771', 'ext_codmanconder', 'video',
                    'https://video-s.twimg.com/amplify_video/2028458925205626880/vid/avc1/720x1280/NwmnDyk3JwF9tYii.mp4?tag=14',
                    ?, 'Retweet media excluded from local storage')
        """, (now,))
        print("Added media_excluded for 2028459277254602771")

    # 3. Fix post 2012506627216453640
    cur.execute("""
        UPDATE articles
        SET account_id = 'ext_xVictorialynnx',
            full_text = '御姐!!https://t.co/WAbjg9OYBu https://t.co/DnbETq1QRB'
        WHERE id = '2012506627216453640'
    """)
    print("Updated 2012506627216453640 -> ext_xVictorialynnx")

    # Update RT 2012531768629342472 reply_to_id
    cur.execute("""
        UPDATE articles
        SET reply_to_id = '2012506627216453640'
        WHERE id = '2012531768629342472'
    """)
    print("Updated RT 2012531768629342472 -> reply_to_id 2012506627216453640")

    # Add media_excluded for 2012506627216453640
    cur.execute("SELECT media_id FROM media_excluded WHERE article_id = '2012506627216453640'")
    if not cur.fetchone():
        now = datetime.datetime.now().isoformat()
        cur.execute("""
            INSERT INTO media_excluded (media_id, article_id, account_id, type, download_url, quarantined_at, quarantine_reason)
            VALUES ('me_2012506627216453640_0', '2012506627216453640', 'ext_xVictorialynnx', 'video',
                    'https://video-s.twimg.com/amplify_video/2012506573189353472/vid/avc1/720x1280/9Y85CRPk79JCAESX.mp4?tag=14',
                    ?, 'Retweet media excluded from local storage')
        """, (now,))
        print("Added media_excluded for 2012506627216453640")

    conn.commit()
    conn.close()
    print("Fix completed successfully!")

if __name__ == '__main__':
    main()
