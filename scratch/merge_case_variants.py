import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 1. Yike_Luo (f35dad07-7066-5405-bb36-7e50e2d94434) を yike_luo (2cda0ead-b020-efbb-be02-6faba155be4b) に統合
parent_id = "2cda0ead-b020-efbb-be02-6faba155be4b"
child_id = "f35dad07-7066-5405-bb36-7e50e2d94434"

cur.execute(
    "UPDATE articles SET account_id = ? WHERE account_id = ?",
    (parent_id, child_id),
)
cur.execute(
    "UPDATE account_profile_histories SET account_id = ? WHERE account_id = ?",
    (parent_id, child_id),
)
cur.execute("DELETE FROM accounts WHERE numeric_id = ?", (child_id,))
print("Merged Yike_Luo into yike_luo.")

# 2. whitelists の大文字小文字名寄せ・重複削除
# 重複を削除: id 17 (Sayapom4), 267 (MsLuo14), 506 (Yike_Luo)
cur.execute("DELETE FROM whitelists WHERE id IN (17, 267, 506)")

# 全whitelistレコードの value を小文字化
cur.execute("UPDATE whitelists SET value = lower(value)")
cur.execute("UPDATE whitelists SET alias_of = '' WHERE alias_of LIKE '%Yike%'")

conn.commit()

print("\nCleaned Whitelist records:")
for r in cur.execute(
    "SELECT id, type, value, group_name, alias_of, is_active FROM whitelists"
).fetchall():
  print(" ", r)

print("\nCleaned Accounts records:")
for r in cur.execute(
    "SELECT numeric_id, username, display_name, group_name FROM accounts"
).fetchall():
  print(" ", r)
