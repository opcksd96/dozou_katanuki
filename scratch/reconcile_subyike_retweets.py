# scratch/reconcile_subyike_retweets.py
import os, sys, glob, json, re, shutil, sqlite3

DB_PATH = 'archive.db'
BACKUP_PATH = 'backups/database/archive_pre_subyike_reconcile.db'
DUMPS_DIR = 'backups/dumps/twitter/subyike'

def main():
    print('=== STEP 1: BACKUP DATABASE ===')
    os.makedirs('backups/database', exist_ok=True)
    shutil.copyfile(DB_PATH, BACKUP_PATH)
    print(f'Backed up {DB_PATH} -> {BACKUP_PATH}')

    print('=== STEP 2: LOAD DUMP METADATA ===')
    dump_map = {}
    meta_files = glob.glob(os.path.join(DUMPS_DIR, '*', 'metadata.json'))
    for mf in meta_files:
        try:
            with open(mf, 'r', encoding='utf-8') as f:
                data = json.load(f)
                tid = data.get('post', {}).get('id')
                if tid: dump_map[tid] = data
        except Exception: pass
    print(f'Loaded {len(dump_map)} dump metadata entries')

    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    print('=== STEP 3: ANALYZE SUBYIKE POSTS ===')
    cur.execute("SELECT id, created_at, full_text, is_repost, reply_to_id, account_id FROM articles WHERE account_id = 'ext_subyike_YIKE' ORDER BY created_at DESC")
    subyike_articles = cur.fetchall()
    print(f'Total subyike articles: {len(subyike_articles)}')

    rt_articles = [a for a in subyike_articles if a[3] == 1 or a[2].startswith('RT @')]
    non_rt_articles = [a for a in subyike_articles if a[3] == 0 and not a[2].startswith('RT @')]
    print(f'RT articles: {len(rt_articles)}, Non-RT candidate articles: {len(non_rt_articles)}')

    matched_pairs = []
    used_orig_ids = set()

    for rt in rt_articles:
        rt_id, rt_c_at, rt_text, _, _, _ = rt
        m = re.match(r'RT @([a-zA-Z0-9_]+):\s*(.*)', rt_text, re.DOTALL)
        if not m: continue
        target_author = m.group(1).strip()
        body_text = m.group(2).strip()

        best_orig = None
        for orig in non_rt_articles:
            orig_id = orig[0]
            if orig_id in used_orig_ids: continue
            orig_text = orig[2].strip()
            if orig_text == body_text or (len(orig_text) >= 15 and (body_text.startswith(orig_text[:20]) or orig_text.startswith(body_text[:20]))):
                best_orig = orig
                break

        if best_orig:
            used_orig_ids.add(best_orig[0])
            matched_pairs.append((rt, best_orig, target_author))

    print(f'Successfully paired {len(matched_pairs)} RTs with original author posts')

    print('=== STEP 4: RECONCILE IDENTITIES & LINKAGES ===')
    reassigned_origs, linked_rts, restored_excluded_media = 0, 0, 0

    for rt, orig, author_name in matched_pairs:
        rt_id, orig_id = rt[0], orig[0]
        canonical_account_id = f'ext_{author_name}'

        cur.execute('SELECT numeric_id FROM accounts WHERE numeric_id = ? OR username = ?', (canonical_account_id, author_name))
        if not cur.fetchone():
            cur.execute("INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at, is_whitelist, post_count) VALUES (?, ?, ?, '', CURRENT_TIMESTAMP, 0, 0)", (canonical_account_id, author_name, author_name))

        cur.execute("UPDATE articles SET account_id = ?, is_repost = 0 WHERE id = ?", (canonical_account_id, orig_id))
        reassigned_origs += 1

        cur.execute("UPDATE articles SET reply_to_id = ?, is_repost = 1 WHERE id = ?", (orig_id, rt_id))
        linked_rts += 1

        meta = dump_map.get(orig_id) or dump_map.get(rt_id)
        if meta and meta.get('media'):
            for med in meta['media']:
                med_id = med.get('media_id') or med.get('filename')
                if not med_id: continue
                d_url = med.get('download_url') or med.get('url') or ''
                thumb_url = med.get('thumbnail_url') or ''
                m_type = med.get('type') or ('video' if 'mp4' in d_url else 'image')
                w = med.get('width') or 0
                h = med.get('height') or 0
                q_reason = f'Retweet Non-Whitelist: {author_name}'

                cur.execute("""
                    INSERT INTO media_excluded (
                        media_id, article_id, account_id, type, download_url, width, height,
                        download_status, failed_reason, tweet_urls, thumbnail_url, quarantined_at, quarantine_reason
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, 'EXCLUDED', 'Retweet External Media', ?, ?, CURRENT_TIMESTAMP, ?)
                    ON CONFLICT(media_id) DO UPDATE SET
                        article_id = excluded.article_id,
                        account_id = excluded.account_id,
                        download_url = excluded.download_url,
                        download_status = excluded.download_status,
                        quarantine_reason = excluded.quarantine_reason
                """, (med_id, orig_id, canonical_account_id, m_type, d_url, w, h, d_url, thumb_url, q_reason))
                restored_excluded_media += 1

    conn.commit()

    print('=== STEP 5: VERIFY FINAL COUNTS ===')
    rem_subyike_orig = cur.execute("SELECT COUNT(*) FROM articles WHERE account_id = 'ext_subyike_YIKE' AND is_repost = 0").fetchone()[0]
    rem_subyike_rt = cur.execute("SELECT COUNT(*) FROM articles WHERE account_id = 'ext_subyike_YIKE' AND is_repost = 1").fetchone()[0]
    total_excluded = cur.execute("SELECT COUNT(*) FROM media_excluded").fetchone()[0]
    total_active_media = cur.execute("SELECT COUNT(*) FROM media").fetchone()[0]

    print(f'Remaining true self posts for subyike (is_repost=0): {rem_subyike_orig}')
    print(f'Retweet posts for subyike (is_repost=1): {rem_subyike_rt}')
    print(f'Original author posts reassigned: {reassigned_origs}')
    print(f'RTs linked with reply_to_id: {linked_rts}')
    print(f'Media in media_excluded (External Stream URLs): {total_excluded}')
    print(f'Media in active media table (UNTOUCHED): {total_active_media}')

    conn.close()
    print('=== FINISHED SUCCESSFULLY ===')

if __name__ == '__main__':
    main()