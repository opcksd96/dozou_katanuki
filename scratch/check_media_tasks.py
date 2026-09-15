import sqlite3

conn = sqlite3.connect("archive.db")
c = conn.cursor()

c.execute("""
    SELECT m.media_id, count(t.id) as task_count
    FROM media m
    JOIN thunder_tasks t ON m.media_id = t.media_id
    WHERE m.download_status = 'ESCALATED'
    GROUP BY m.media_id
    ORDER BY task_count DESC
    LIMIT 5
""")
rows = c.fetchall()
print("Top 5 escalated media with thunder_tasks:")
for r in rows:
    print(f"MediaID: {r[0]}, Task count: {r[1]}")
