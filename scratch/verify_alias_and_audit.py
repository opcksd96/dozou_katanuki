import sqlite3

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

def fetch_articles_by_account(acc_id):
    sql = """
        SELECT id, account_id, via, substr(original_url, 1, 60), substr(full_text, 1, 30)
        FROM articles
        WHERE account_id = ? OR account_id IN (
            SELECT numeric_id FROM accounts
            WHERE numeric_id = ? OR username = ? OR alias_of = ?
               OR alias_of IN (SELECT username FROM accounts WHERE numeric_id = ? OR username = ?)
               OR username IN (SELECT alias_of FROM accounts WHERE (numeric_id = ? OR username = ?) AND alias_of != '')
        )
    """
    return cur.execute(sql, (acc_id, acc_id, acc_id, acc_id, acc_id, acc_id, acc_id, acc_id)).fetchall()

print("=== 1. AUDIT: NENNE1001 (by numeric_id: 1553745049057300481) ===")
rows_nenne = fetch_articles_by_account('1553745049057300481')
print(f"Total articles retrieved: {len(rows_nenne)} (Expected: 3)")
for r in rows_nenne:
    print(f"  Post {r[0]} | account_id: {r[1]} | via: {r[2]} | {r[4]}")

print("\n=== 2. AUDIT: NENNE1001 (by username: nenne1001) ===")
rows_nenne_u = fetch_articles_by_account('nenne1001')
print(f"Total articles retrieved by username: {len(rows_nenne_u)} (Expected: 3)")

print("\n=== 3. AUDIT: MSLUO14 (by numeric_id: 4fa935dc-eda5-4217-4b5b-c9c9ea0fb491) ===")
rows_msluo = fetch_articles_by_account('4fa935dc-eda5-4217-4b5b-c9c9ea0fb491')
print(f"Total articles retrieved: {len(rows_msluo)} (Expected: 163)")
# Check if any fake handle is in the retrieved posts
fakes_in_msluo = [r for r in rows_msluo if any(f in (r[3] or '') for f in ['MsLuo1433', 'MsLuo14b', 'MsLuo14d', 'MsLuo14hh'])]
print(f"Fake posts in msluo14 timeline: {len(fakes_in_msluo)} (Expected: 0)")

print("\n=== 4. AUDIT: FAKE ACCOUNT MsLuo1433 (by numeric_id: 2039493607183249408) ===")
rows_fake = fetch_articles_by_account('2039493607183249408')
print(f"Total articles for MsLuo1433: {len(rows_fake)} (Expected: 1)")
for r in rows_fake:
    print(f"  Post {r[0]} | account_id: {r[1]} | {r[4]}")

print("\n=== 5. ACCOUNTS TABLE SUMMARY ===")
for r in cur.execute("SELECT numeric_id, username, display_name, alias_of, is_whitelist FROM accounts WHERE username LIKE '%luo%' OR username LIKE '%nenne%'").fetchall():
    print(f"  {r[0]} | @{r[1]} | {r[2]} | alias_of: '{r[3]}' | whitelist: {r[4]}")
