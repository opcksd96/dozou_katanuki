import sqlite3

con = sqlite3.connect("archive.db")
cur = con.cursor()

mediaID = "GEFsflnacAAOXck"
cleanID = "GEFsflnacAAOXck"

print("--- Test Query 1: Exactly as written in repo_media_destination.go ---")
sql = """
SELECT accounts.username FROM media
JOIN accounts ON (accounts.numeric_id = media.account_id OR accounts.username = media.account_id)
WHERE media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?
"""
params = (mediaID, cleanID + "", cleanID + "%")
print("SQL:", sql)
print("Params:", params)
res = cur.execute(sql, params).fetchall()
print("Result:", res)

print("\n--- Test Direct lookup without JOIN ---")
res2 = cur.execute("SELECT media_id, account_id FROM media WHERE media_id LIKE ?", (cleanID + "%",)).fetchall()
print("media match:", res2)

if res2:
    acc_id = res2[0][1]
    print(f"account_id is: '{acc_id}'")
    acc = cur.execute("SELECT numeric_id, username FROM accounts WHERE numeric_id = ? OR username = ?", (acc_id, acc_id)).fetchall()
    print("accounts match:", acc)

con.close()
