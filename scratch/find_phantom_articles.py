import sqlite3
import re

con = sqlite3.connect('archive.db')
cur = con.cursor()

# Get all retweets (RT @username: ...)
cur.execute("SELECT id, account_id, full_text, created_at FROM articles WHERE full_text LIKE 'RT @%'")
rts = cur.fetchall()
print(f"Total RT articles in DB: {len(rts)}")

matched_pairs = []
for rt_id, acc_id, rt_text, rt_created_at in rts:
    m = re.match(r"RT @([a-zA-Z0-9_]+):\s*(.*)", rt_text, re.DOTALL)
    if not m:
        continue
    orig_author, orig_content = m.group(1), m.group(2).strip()
    
    # 1. Exact text match under same account_id
    cur.execute("""
        SELECT id, account_id, is_repost, full_text, created_at FROM articles 
        WHERE account_id = ? AND is_repost = 0 AND full_text = ?
    """, (acc_id, orig_content))
    matches = cur.fetchall()
    for match in matches:
        matched_pairs.append({
            "rt_id": rt_id,
            "orig_author": orig_author,
            "phantom_id": match[0],
            "acc_id": acc_id,
            "phantom_created_at": match[4],
            "text": match[3]
        })

print(f"Matched phantom original articles under whitelist account: {len(matched_pairs)}")
for item in matched_pairs[:25]:
    print(f"  RT: {item['rt_id']} (@{item['orig_author']}) -> Phantom: {item['phantom_id']} falsely under {item['acc_id']}")
    print(f"     Text: {repr(item['text'][:40])}")
