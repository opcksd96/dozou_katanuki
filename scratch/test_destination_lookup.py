import os
import sqlite3

conn = sqlite3.connect("archive.db", timeout=60.0)
cur = conn.cursor()

def get_media_owner_username(media_id: str):
    clean_id = os.path.splitext(media_id)[0]
    for sfx in ["_motrix", "_requests", "_thunder", "_plain", "_orig", "_large", "_wayback_orig", "_wayback_plain", "_wayback"]:
        if clean_id.endswith(sfx):
            clean_id = clean_id[:-len(sfx)]
            break

    # 1. 直接所有権ルート: media.account_id -> accounts
    row = cur.execute("""
        SELECT ac.username FROM media m
        JOIN accounts ac ON (ac.numeric_id = m.account_id OR ac.username = m.account_id)
        WHERE m.media_id = ? OR m.media_id = ? OR m.media_id LIKE ?
        LIMIT 1
    """, (media_id, clean_id + os.path.splitext(media_id)[1], clean_id + "%")).fetchone()
    if row and row[0]:
        return row[0], "Route 1: Direct Ownership (media.account_id)"

    # 2. 原本記事リレーションルート: media.article_id -> articles.account_id
    row = cur.execute("""
        SELECT ac.username FROM media m
        JOIN articles a ON a.id = m.article_id
        JOIN accounts ac ON (ac.numeric_id = a.account_id OR ac.username = a.account_id)
        WHERE m.media_id = ? OR m.media_id = ? OR m.media_id LIKE ?
        LIMIT 1
    """, (media_id, clean_id + os.path.splitext(media_id)[1], clean_id + "%")).fetchone()
    if row and row[0]:
        return row[0], "Route 2: Article Relation (media.article_id)"

    # 3. 迅雷/DLタスク記録ルート: thunder_tasks
    row = cur.execute("""
        SELECT ac.username FROM thunder_tasks t
        JOIN articles a ON a.id = t.article_id
        JOIN accounts ac ON (ac.numeric_id = a.account_id OR ac.username = a.account_id)
        WHERE t.file_name = ? OR t.media_id = ? OR t.media_id LIKE ?
        LIMIT 1
    """, (media_id, clean_id, clean_id + "%")).fetchone()
    if row and row[0]:
        return row[0], "Route 3: Task Record (thunder_tasks)"

    # 4. URLセマンティック・パターンマッチ
    row = cur.execute("""
        SELECT ac.username FROM media m
        JOIN accounts ac ON (m.tweet_urls LIKE '%/' || ac.username || '/%' OR m.download_url LIKE '%/' || ac.username || '/%')
        WHERE m.media_id = ? OR m.media_id = ? OR m.media_id LIKE ?
        LIMIT 1
    """, (media_id, clean_id + os.path.splitext(media_id)[1], clean_id + "%")).fetchone()
    if row and row[0]:
        return row[0], "Route 4: URL Semantic Match"

    # 5. フォールバック: clean_id が article_id
    row = cur.execute("""
        SELECT ac.username FROM articles a
        JOIN accounts ac ON (ac.numeric_id = a.account_id OR ac.username = a.account_id)
        WHERE a.id = ?
        LIMIT 1
    """, (clean_id,)).fetchone()
    if row and row[0]:
        return row[0], "Route 5: CleanID as ArticleID"

    return None, "Not found"

# テスト実行
sample_media = cur.execute("SELECT media_id FROM media LIMIT 10").fetchall()
print("=== Reverse Lookup Verification ===")
all_pass = True
for (mid,) in sample_media:
    # 通常のmedia_id
    u, route = get_media_owner_username(mid)
    # 迅雷が保存するサフィックス付きファイル名
    clean = os.path.splitext(mid)[0]
    ext = os.path.splitext(mid)[1] or ".jpg"
    thunder_fn = f"{clean}_orig{ext}"
    u2, route2 = get_media_owner_username(thunder_fn)

    print(f"File: {thunder_fn:<30} -> Owner: {u2:<15} [{route2}]")
    if not u2:
        all_pass = False

print("\nAll sample lookups passed:", all_pass)
conn.close()
