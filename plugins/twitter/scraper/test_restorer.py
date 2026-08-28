# plugins/twitter/scraper/test_restorer.py (100行以下)
import json, os, shutil, sqlite3, sys, tempfile
try:
    from plugins.twitter.scraper.core.restorer import Restorer
except ImportError:
    from core.restorer import Restorer


def init_test_db(db_path: str):
    conn = sqlite3.connect(db_path)
    cur = conn.cursor()
    cur.execute("CREATE TABLE accounts (numeric_id TEXT PRIMARY KEY, username TEXT NOT NULL, display_name TEXT, avatar_url TEXT, avatar_base64 TEXT, description TEXT, group_name TEXT DEFAULT '', alias_of TEXT DEFAULT '', updated_at DATETIME);")
    cur.execute("CREATE TABLE account_profile_histories (id INTEGER PRIMARY KEY AUTOINCREMENT, account_id TEXT NOT NULL, display_name TEXT NOT NULL, description TEXT DEFAULT '', avatar_original_url TEXT NOT NULL, avatar_seq INTEGER NOT NULL, avatar_virtual_key TEXT NOT NULL, avatar_base64 TEXT, observed_at DATETIME NOT NULL);")
    cur.execute("CREATE TABLE whitelists (id INTEGER PRIMARY KEY, type TEXT, value TEXT NOT NULL UNIQUE, is_active INTEGER);")
    cur.execute("CREATE TABLE articles (id TEXT PRIMARY KEY, account_id TEXT NOT NULL, conversation_id TEXT, reply_to_id TEXT, reply_to_handle TEXT, created_at DATETIME, full_text TEXT, lang TEXT DEFAULT 'ja', full_text_ja TEXT, full_text_en TEXT, full_text_zh TEXT, via TEXT DEFAULT 'twitter', is_repost INTEGER DEFAULT 0, is_liked INTEGER DEFAULT 0, wayback_url TEXT, source_name TEXT, source_domain TEXT, original_url TEXT, sotwe_url TEXT, nitter_url TEXT, twistalker_url TEXT, is_trash BOOLEAN DEFAULT 0, trashed_by TEXT, trash_reason TEXT, trashed_at DATETIME);")
    cur.execute("CREATE TABLE media (media_id TEXT PRIMARY KEY, article_id TEXT NOT NULL, type TEXT, download_url TEXT, width INTEGER, height INTEGER, download_status TEXT, failed_reason TEXT, stash_scene_id TEXT, stash_image_id TEXT, media_quality TEXT);")
    cur.execute("CREATE TABLE url_redirects (short_url TEXT PRIMARY KEY, expanded_url TEXT, article_id TEXT);")
    conn.commit(); conn.close()


def test_restorer():
    temp_dir = tempfile.mkdtemp()
    db_path = os.path.join(temp_dir, "test_restore.db")
    dumps_dir, avatar_dir, blobs_dir = os.path.join(temp_dir, "dumps"), os.path.join(temp_dir, "avatars"), os.path.join(temp_dir, "blobs")

    try:
        init_test_db(db_path)
        post_dir = os.path.join(dumps_dir, "twitter", "mash_kyrielight", "123456789")
        os.makedirs(os.path.join(post_dir, "avatars"), exist_ok=True)
        with open(os.path.join(post_dir, "avatars", "mash_avatar.jpg"), "w") as f: f.write("mock_avatar")

        sample_meta = {
            "account": {"numeric_id": "99999", "username": "mash_kyrielight", "display_name": "マシュ", "avatar_url": "https://pbs.twimg.com/mash.jpg"},
            "post": {"id": "123456789", "created_at": "2026-08-20 00:00:00", "full_text": "先輩、準備完了です！ https://t.co/xyz",
                     "urls": [{"short_url": "https://t.co/xyz", "expanded_url": "https://example.com/chaldea"}], "wayback_url": "https://web.archive.org/web/..."},
            "media": [{"url": "https://pbs.twimg.com/media/shield.jpg", "media_id": "shield.jpg", "type": "image", "width": 1200, "height": 800}]
        }
        with open(os.path.join(post_dir, "metadata.json"), "w", encoding="utf-8") as f: json.dump(sample_meta, f, ensure_ascii=False)

        os.makedirs(blobs_dir, exist_ok=True)
        with open(os.path.join(blobs_dir, "shield.jpg"), "wb") as f: f.write(b"mock_image_bytes")

        restorer = Restorer(dumps_dir=dumps_dir, db_path=db_path, storage_dir=blobs_dir, avatar_dir=avatar_dir)
        from unittest.mock import patch
        with patch.object(restorer.downloader.reconciler, "register_media", return_value="img-123"):
            stats = restorer.run_restore()

        assert stats["articles"] == 1 and stats["avatars"] == 1
        assert os.path.exists(os.path.join(avatar_dir, "mash_avatar.jpg"))

        conn = sqlite3.connect(db_path); cur = conn.cursor()
        acc = cur.execute("SELECT username FROM accounts WHERE numeric_id = '99999'").fetchone()
        assert acc == ("mash_kyrielight",)
        art = cur.execute("SELECT full_text FROM articles WHERE id = '123456789'").fetchone()
        assert "https://example.com/chaldea" in art[0]
        med = cur.execute("SELECT download_status FROM media WHERE media_id = 'shield.jpg'").fetchone()
        assert med[0] == "COMPLETED"
        conn.close()
    finally:
        shutil.rmtree(temp_dir, ignore_errors=True)


if __name__ == "__main__":
    test_restorer()
