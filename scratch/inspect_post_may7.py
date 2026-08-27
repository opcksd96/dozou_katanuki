import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

print("Article 1920037769507676376 (5/7 RT):")
print(
    cur.execute(
        "SELECT id, account_id, full_text, created_at FROM articles WHERE id ="
        " ?",
        ("1920037769507676376",),
    ).fetchall()
)
print(
    "Media for 1920037769507676376:",
    cur.execute(
        "SELECT media_id, download_status, stash_scene_id, stash_image_id FROM"
        " media WHERE article_id = ?",
        ("1920037769507676376",),
    ).fetchall(),
)

print("\nArticle 1919718358611853450 (5/6 Original):")
print(
    cur.execute(
        "SELECT id, account_id, full_text, created_at FROM articles WHERE id ="
        " ?",
        ("1919718358611853450",),
    ).fetchall()
)
print(
    "Media for 1919718358611853450:",
    cur.execute(
        "SELECT media_id, download_status, stash_scene_id, stash_image_id FROM"
        " media WHERE article_id = ?",
        ("1919718358611853450",),
    ).fetchall(),
)

print("\nmsluo14 articles around 2025/05:")
for r in cur.execute(
    "SELECT id, account_id, created_at, substr(full_text, 1, 40) FROM articles"
    " WHERE (account_id = '4fa935dc-eda5-4217-4b5b-c9c9ea0fb491' OR"
    " full_text LIKE '%跳操%') ORDER BY created_at DESC"
).fetchall():
  print(" ", r)
