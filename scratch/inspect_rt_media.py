import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

print("Article 1919718358611853450 (Original):")
print(
    cur.execute(
        "SELECT id, is_repost, created_at, full_text FROM articles WHERE id ="
        " '1919718358611853450'"
    ).fetchall()
)
print(
    "Media for 1919718358611853450:",
    cur.execute(
        "SELECT media_id, download_status, stash_scene_id, stash_image_id FROM"
        " media WHERE article_id = '1919718358611853450'"
    ).fetchall(),
)

print("\nArticle 1920037769507676376 (RT):")
print(
    cur.execute(
        "SELECT id, is_repost, created_at, full_text FROM articles WHERE id ="
        " '1920037769507676376'"
    ).fetchall()
)
print(
    "Media for 1920037769507676376:",
    cur.execute(
        "SELECT media_id, download_status, stash_scene_id, stash_image_id FROM"
        " media WHERE article_id = '1920037769507676376'"
    ).fetchall(),
)
