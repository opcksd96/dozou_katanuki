import sqlite3

conn = sqlite3.connect('archive.db')
cur = conn.cursor()

print('nitter_url not empty:', cur.execute('SELECT COUNT(*) FROM articles WHERE nitter_url IS NOT NULL AND nitter_url != ""').fetchone()[0])
print('twistalker_url not empty:', cur.execute('SELECT COUNT(*) FROM articles WHERE twistalker_url IS NOT NULL AND twistalker_url != ""').fetchone()[0])
print('sotwe_url not empty:', cur.execute('SELECT COUNT(*) FROM articles WHERE sotwe_url IS NOT NULL AND sotwe_url != ""').fetchone()[0])
