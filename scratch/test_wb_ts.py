import sqlite3, requests, re

conn = sqlite3.connect('archive.db')
row = conn.cursor().execute("SELECT a.wayback_url, m.download_url FROM articles a JOIN media m ON a.id = m.article_id WHERE m.media_id LIKE '%E1gGLDLVEAIQC2E%'").fetchone()
print('Row:', row)
if row:
    wb_url, dl_url = row
    m = re.search(r'/web/(\d{14})', wb_url) if wb_url else None
    ts = m.group(1) if m else ''
    print('Timestamp:', ts)
    s = requests.Session()
    s.headers.update({"User-Agent": "Mozilla/5.0"})
    urls_to_test = [
        dl_url,
        f"https://web.archive.org/web/{ts}im_/{dl_url}",
        f"https://web.archive.org/web/{ts}oe_/{dl_url}",
        f"https://web.archive.org/web/{ts}/{dl_url}",
        f"https://web.archive.org/web/{ts}im_/{dl_url.replace('https://', 'http://')}",
        f"https://web.archive.org/web/2/{dl_url}",
        f"https://web.archive.org/web/2/{dl_url.replace('https://', 'http://')}"
    ]
    for u in urls_to_test:
        try:
            r = s.get(u, timeout=5, allow_redirects=True)
            print(u, r.status_code, len(r.content), r.headers.get('Content-Type'))
        except Exception as e:
            print(u, e)
