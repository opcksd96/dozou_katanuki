import sqlite3
import re

con = sqlite3.connect('archive.db')
cur = con.cursor()

# Check all articles of subyike
cur.execute("""
    SELECT id, account_id, is_repost, full_text, created_at FROM articles 
    WHERE account_id LIKE '%subyike%'
    ORDER BY created_at DESC
""")
subyike_all = cur.fetchall()

print(f"Total articles under subyike: {len(subyike_all)}")

# Find all RTs
rts = {}
for aid, acc, is_rp, txt, cat in subyike_all:
    if txt.startswith("RT @"):
        m = re.match(r"RT @([a-zA-Z0-9_]+):\s*(.*)", txt, re.DOTALL)
        if m:
            orig_user = m.group(1)
            orig_txt = m.group(2).strip()
            rts[orig_txt] = (aid, orig_user, cat)

print(f"Subyike RT count: {len(rts)}")

phantoms = []
for aid, acc, is_rp, txt, cat in subyike_all:
    if not txt.startswith("RT @"):
        stripped = txt.strip()
        if stripped in rts:
            rt_id, orig_user, rt_cat = rts[stripped]
            phantoms.append((aid, orig_user, rt_id, stripped[:40]))

print(f"Subyike Phantom Non-RT articles (originated from external users): {len(phantoms)}")
for ph_id, orig_user, rt_id, txt_snippet in phantoms:
    print(f"  Phantom: {ph_id} (Real Author: @{orig_user}) (RT: {rt_id}) -> Text: {repr(txt_snippet)}")
