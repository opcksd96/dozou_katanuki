# scratch/validate_stash_image_bindings.py
import urllib.request, json, sqlite3

def main():
    conn = sqlite3.connect('archive.db')
    cur = conn.cursor()
    cur.execute("SELECT media_id, article_id, stash_image_id, download_status FROM media WHERE stash_image_id IS NOT NULL AND stash_image_id != ''")
    db_rows = cur.fetchall()
    print(f'Total media with stash_image_id in DB: {len(db_rows)}')

    req = urllib.request.Request(
        'http://127.0.0.1:9999/graphql',
        data=json.dumps({'query': '{ allImages { id files { basename } } }'}).encode(),
        headers={'Content-Type': 'application/json'}
    )
    with urllib.request.urlopen(req) as res:
        data = json.loads(res.read().decode())
        stash_images = {img['id']: [f['basename'] for f in img.get('files', [])] for img in data.get('data', {}).get('allImages', [])}

    matched, mismatched = [], []
    for mid, aid, s_id, status in db_rows:
        files = stash_images.get(s_id, [])
        is_match = any(mid == fb or mid.split('.')[0] == fb.split('.')[0] for fb in files)
        if is_match:
            matched.append((mid, aid, s_id))
        else:
            mismatched.append((mid, aid, s_id, files))

    print(f'Correctly matched: {len(matched)}')
    print(f'Mismatched (wrongly bound): {len(mismatched)}')
    print('\nFirst 10 mismatched:')
    for m in mismatched[:10]:
        print(f'  DB MediaID: {m[0]} (Article: {m[1]}) -> Stash #{m[2]} actually has files: {m[3]}')

if __name__ == '__main__':
    main()
