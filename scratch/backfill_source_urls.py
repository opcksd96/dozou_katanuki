import sqlite3, re

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

def extract_handle_and_id(orig_url, pid):
    if not orig_url: return None, pid
    m = re.search(r'(?:twitter\.com|x\.com|twstalker\.com|nitter\.[a-z0-9.\-]+)/([a-zA-Z0-9_]+)(?:/status(?:es)?/(\d+))?', orig_url)
    if m:
        return m.group(1), (m.group(2) or pid)
    return None, pid

print("=== BACKFILLING twistalker_url and nitter_url ===")

# 1. Backfill twistalker_url
rows_tw = cur.execute("""
    SELECT id, original_url, twistalker_url
    FROM articles
    WHERE via IN ('TwStalker', 'twistalker')
       OR source_name IN ('TwStalker', 'twistalker')
       OR original_url LIKE '%twistalker%'
""").fetchall()

cnt_tw = 0
for pid, orig_url, current_tws in rows_tw:
    handle, status_id = extract_handle_and_id(orig_url, pid)
    handle = handle or "unknown"
    tws_url = f"https://twstalker.com/{handle}/status/{status_id}"
    cur.execute("UPDATE articles SET twistalker_url = ? WHERE id = ?", (tws_url, pid))
    cnt_tw += 1
    print(f"Updated TwStalker: {pid} -> {tws_url}")

# 2. Backfill nitter_url
rows_nit = cur.execute("""
    SELECT id, original_url, nitter_url
    FROM articles
    WHERE via IN ('Nitter', 'nitter')
       OR source_name IN ('Nitter', 'nitter')
       OR original_url LIKE '%nitter%'
""").fetchall()

cnt_nit = 0
for pid, orig_url, current_nit in rows_nit:
    handle, status_id = extract_handle_and_id(orig_url, pid)
    handle = handle or "unknown"
    nit_url = f"https://nitter.space/{handle}/status/{status_id}"
    cur.execute("UPDATE articles SET nitter_url = ? WHERE id = ?", (nit_url, pid))
    cnt_nit += 1
    print(f"Updated Nitter: {pid} -> {nit_url}")

conn.commit()

print(f"\nCompleted! Backfilled {cnt_tw} twistalker_url and {cnt_nit} nitter_url.")

# Verification
cnt_tw_filled = cur.execute("SELECT COUNT(*) FROM articles WHERE twistalker_url IS NOT NULL AND twistalker_url != ''").fetchone()[0]
cnt_nit_filled = cur.execute("SELECT COUNT(*) FROM articles WHERE nitter_url IS NOT NULL AND nitter_url != ''").fetchone()[0]
print(f"Verification: twistalker_url non-empty count = {cnt_tw_filled}")
print(f"Verification: nitter_url non-empty count = {cnt_nit_filled}")

conn.close()
