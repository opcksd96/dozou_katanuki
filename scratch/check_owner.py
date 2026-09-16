import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

print("=== CHECK MEDIA & ACCOUNT FOR GnmCCzebYAAkg9-.jpg ===")
c.execute("""
    SELECT m.media_id, m.download_status, m.account_id, a.account_id, a.id, acc.username
    FROM media m
    LEFT JOIN articles a ON m.article_id = a.id
    LEFT JOIN accounts acc ON a.account_id = acc.numeric_id
    WHERE m.media_id LIKE '%GnmCCzebYAAkg9-%'
""")
for r in c.fetchall():
    print(r)

print("\n=== ACCOUNTS LIKE MSLUO ===")
c.execute("SELECT numeric_id, username, display_name FROM accounts WHERE username LIKE '%msluo%' OR display_name LIKE '%罗%'")
for r in c.fetchall():
    print(r)
