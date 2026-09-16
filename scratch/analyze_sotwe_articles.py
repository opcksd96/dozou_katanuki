import sqlite3

conn = sqlite3.connect('archive.db')
c = conn.cursor()

c.execute("""
SELECT a.id, a.account_id, a.sotwe_url, a.source_name, COUNT(m.media_id)
FROM articles a
LEFT JOIN media m ON a.id = m.article_id
WHERE a.sotwe_url IS NOT NULL OR a.source_name = 'sotwe'
GROUP BY a.id
ORDER BY COUNT(m.media_id) DESC
LIMIT 30
""")
rows = c.fetchall()
print(f"Sample articles from sotwe (top media counts):")
for r in rows:
    print(r)

c.execute("""
SELECT COUNT(m.media_id), COUNT(DISTINCT a.id)
FROM articles a
LEFT JOIN media m ON a.id = m.article_id
WHERE a.sotwe_url IS NOT NULL OR a.source_name = 'sotwe'
GROUP BY a.id
""")
counts = [r[0] for r in c.fetchall()]
from collections import Counter
print("Media count distribution for Sotwe articles:", Counter(counts))
