# plugins/twitter/scraper/test_restorer.py
import json
import os
import shutil
import sqlite3
import tempfile
from core.restorer import Restorer


def init_test_db(db_path: str):
    conn = sqlite3.connect(db_path)
    cur = conn.cursor()
    cur.execute("""
        CREATE TABLE accounts (
            numeric_id TEXT PRIMARY KEY,
            username TEXT NOT NULL,
            display_name TEXT,
            avatar_url TEXT,
            updated_at DATETIME
        );
    """)
    cur.execute("""
        CREATE TABLE articles (
            id TEXT PRIMARY KEY,
            account_id TEXT NOT NULL,
            conversation_id TEXT,
            reply_to_id TEXT,
            reply_to_handle TEXT,
            created_at DATETIME,
            full_text TEXT,
            lang TEXT DEFAULT 'ja',
            full_text_ja TEXT,
            full_text_en TEXT,
            full_text_zh TEXT,
            via TEXT DEFAULT 'twitter',
            is_repost INTEGER DEFAULT 0,
            is_liked INTEGER DEFAULT 0,
            wayback_url TEXT,
            FOREIGN KEY (account_id) REFERENCES accounts(numeric_id)
        );
    """)
    cur.execute("""
        CREATE TABLE media (
            media_id TEXT PRIMARY KEY,
            article_id TEXT NOT NULL,
            type TEXT,
            download_url TEXT,
            width INTEGER,
            height INTEGER,
            download_status TEXT,
            failed_reason TEXT,
            stash_scene_id TEXT,
            stash_image_id TEXT,
            FOREIGN KEY (article_id) REFERENCES articles(id)
        );
    """)
    cur.execute("""
        CREATE TABLE url_redirects (
            short_url TEXT PRIMARY KEY,
            expanded_url TEXT,
            article_id TEXT
        );
    """)
    conn.commit()
    conn.close()


def test_restorer():
    temp_dir = tempfile.mkdtemp()
    db_path = os.path.join(temp_dir, "test_restore.db")
    dumps_dir = os.path.join(temp_dir, "dumps")
    avatar_dir = os.path.join(temp_dir, "avatars")
    blobs_dir = os.path.join(temp_dir, "blobs")

    try:
        init_test_db(db_path)

        # テスト用の dump データ作成
        post_dir = os.path.join(dumps_dir, "twitter", "mash_kyrielight", "123456789")
        os.makedirs(post_dir, exist_ok=True)
        avatar_sub = os.path.join(post_dir, "avatars")
        os.makedirs(avatar_sub, exist_ok=True)

        with open(os.path.join(avatar_sub, "mash_avatar.jpg"), "w") as f:
            f.write("mock_avatar_binary_data")

        sample_meta = {
            "account": {
                "numeric_id": "99999",
                "username": "mash_kyrielight",
                "display_name": "マシュ・キリエライト",
                "avatar_url": "https://pbs.twimg.com/profile_images/mash.jpg"
            },
            "post": {
                "id": "123456789",
                "created_at": "2026-08-20 00:00:00",
                "full_text": "先輩、カルデアより出撃準備完了です！ https://t.co/xyz",
                "urls": [{"short_url": "https://t.co/xyz", "expanded_url": "https://example.com/chaldea"}],
                "wayback_url": "https://web.archive.org/web/20260820000000/https://twitter.com/mash/status/123456789"
            },
            "media": [
                {
                    "url": "https://pbs.twimg.com/media/shield.jpg",
                    "media_id": "shield.jpg",
                    "type": "image",
                    "width": 1200,
                    "height": 800
                }
            ]
        }

        with open(os.path.join(post_dir, "metadata.json"), "w", encoding="utf-8") as f:
            json.dump(sample_meta, f, ensure_ascii=False)

        # 疑似メディア実体を blobs に配置
        os.makedirs(blobs_dir, exist_ok=True)
        with open(os.path.join(blobs_dir, "shield.jpg"), "wb") as f:
            f.write(b"mock_image_bytes")

        # Restorer 実行
        logs = []
        def callback(cur, tot, msg):
            logs.append(f"{cur}/{tot}: {msg}")

        restorer = Restorer(dumps_dir=dumps_dir, db_path=db_path, storage_dir=blobs_dir, avatar_dir=avatar_dir)
        stats = restorer.run_restore(progress_callback=callback)

        assert stats["articles"] == 1, f"Expected 1 article, got {stats['articles']}"
        assert stats["avatars"] == 1, f"Expected 1 avatar, got {stats['avatars']}"
        assert os.path.exists(os.path.join(avatar_dir, "mash_avatar.jpg")), "Avatar file was not copied"

        # DB 検証
        conn = sqlite3.connect(db_path)
        cur = conn.cursor()
        acc = cur.execute("SELECT username, display_name FROM accounts WHERE numeric_id = '99999'").fetchone()
        assert acc == ("mash_kyrielight", "マシュ・キリエライト"), f"Account mismatch: {acc}"

        art = cur.execute("SELECT id, full_text FROM articles WHERE id = '123456789'").fetchone()
        assert art[0] == "123456789"
        assert "https://example.com/chaldea" in art[1], f"URL expanded mismatch: {art[1]}"

        med = cur.execute("SELECT media_id, download_status FROM media WHERE media_id = 'shield.jpg'").fetchone()
        assert med[0] == "shield.jpg"
        assert med[1] == "COMPLETED", f"Media status expected COMPLETED, got {med[1]}"
        conn.close()

        print("[TEST PASSED] Restorer offline DB reconstruction verified successfully!")
    finally:
        shutil.rmtree(temp_dir, ignore_errors=True)


if __name__ == "__main__":
    test_restorer()
