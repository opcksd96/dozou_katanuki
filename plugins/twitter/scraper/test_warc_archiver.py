import os, sys, tempfile, shutil, json
_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../../.."))
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.core.warc_archiver import WarcArchiver
from warcio.archiveiterator import ArchiveIterator

def test_warc_archiver():
    tmp_dir = tempfile.mkdtemp()
    try:
        archiver = WarcArchiver(base_dir=tmp_dir)
        mock_posts = [{
            "platform": "twitter",
            "account": {"numeric_id": "12345", "username": "testuser", "display_name": "Test User", "avatar_url": "https://pbs.twimg.com/avatar.jpg"},
            "post": {
                "id": "999888777", "conversation_id": "999888777",
                "created_at": "2026-08-28 12:00:00",
                "full_text": "Hello, Antigravity and Mash! #NES",
                "wayback_url": "", "urls": []
            },
            "media": [{"url": "https://pbs.twimg.com/media/test.jpg", "type": "image", "width": 800, "height": 600}]
        }]
        
        count = archiver.archive_posts(mock_posts, platform="twitter")
        assert count == 1, f"Expected 1, got {count}"
        
        dump_dir = os.path.join(tmp_dir, "twitter", "testuser", "999888777")
        assert os.path.exists(dump_dir), f"Directory {dump_dir} not created"
        
        meta_file = os.path.join(dump_dir, "metadata.json")
        warc_file = os.path.join(dump_dir, "snapshot.warc.gz")
        assert os.path.exists(meta_file), "metadata.json missing"
        assert os.path.exists(warc_file), "snapshot.warc.gz missing"
        
        with open(meta_file, "r", encoding="utf-8") as f:
            data = json.load(f)
            assert data["post"]["id"] == "999888777"
            assert data["account"]["username"] == "testuser"
            
        records_found = []
        with open(warc_file, "rb") as f:
            for r in ArchiveIterator(f):
                records_found.append((r.rec_type, r.rec_headers.get_header("WARC-Target-URI")))
                
        assert len(records_found) == 2, f"Expected 2 records (response, metadata), found {records_found}"
        assert records_found[0][0] == "response"
        assert "999888777" in records_found[0][1]
        assert records_found[1][0] == "metadata"
        assert "#metadata" in records_found[1][1]
        
        print("✅ WarcArchiver test passed successfully!")
    finally:
        shutil.rmtree(tmp_dir, ignore_errors=True)

if __name__ == "__main__":
    test_warc_archiver()
