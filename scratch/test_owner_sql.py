import sqlite3, os, re

conn = sqlite3.connect('archive.db')
c = conn.cursor()

# Go code in repo_media_destination.go:
# cleanID := strings.TrimSuffix(mediaID, filepath.Ext(mediaID))
# for _, sfx := range []string{"_orig", "_large", "_wayback_orig", "_wayback"} {
#     cleanID = strings.TrimSuffix(cleanID, sfx)
# }
# r.db.Table("media").
#     Select("accounts.username").
#     Joins("JOIN articles ON articles.id = media.article_id").
#     Joins("JOIN accounts ON (accounts.numeric_id = articles.account_id OR accounts.username = articles.account_id)").
#     Where("media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?", mediaID, cleanID, cleanID+"%").
#     First(&result)

# But wait! In processDirectoryFiles:
# name := entry.Name()
# What was the original downloaded name in D:\迅雷下载 ?
# Could it have been:
# GnmCCzebYAAkg9-_large.jpg or GnmCCzebYAAkg9-_plain.jpg or GnmCCzebYAAkg9-.jpg ?

# Let's test different name variations in SQL!
names_to_test = [
    'GnmCCzebYAAkg9-.jpg',
    'GnmCCzebYAAkg9-',
    'GnmCCzebYAAkg9-_large.jpg',
    'GnmCCzebYAAkg9-_large',
    'GnmCCzebYAAkg9-_plain.jpg',
    'GnmCCzebYAAkg9-_plain',
]

for n in names_to_test:
    sql = """
    SELECT accounts.username, articles.id, media.media_id
    FROM media
    JOIN articles ON articles.id = media.article_id
    JOIN accounts ON (accounts.numeric_id = articles.account_id OR accounts.username = articles.account_id)
    WHERE media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?
    """
    c.execute(sql, (n, n, n + '%'))
    rows = c.fetchall()
    print(f"Query for '{n}': {rows}")

# Now wait! Check what articles.account_id actually is in DB:
c.execute("SELECT id, account_id FROM articles WHERE id = '1907698917766213981'")
art = c.fetchone()
print("\nArticle:", art)

# And what are the accounts:
c.execute("SELECT numeric_id, username FROM accounts WHERE numeric_id = ? OR username = ?", (art[1], 'msluo14'))
accs = c.fetchall()
print("Accounts matching article.account_id:", accs)
