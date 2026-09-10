# scratch/verify_fixed_timeline.py
import urllib.request, json, sqlite3

def main():
    conn = sqlite3.connect('archive.db')
    cur = conn.cursor()

    cur.execute("SELECT count(*) FROM articles WHERE account_id = 'ext_subyike_YIKE' AND is_repost = 0")
    print('subyike non-RT count in DB:', cur.fetchone()[0])

    print('\n--- API: all timeline (top 5) ---')
    with urllib.request.urlopen('http://localhost:5173/api/timeline?platform=twitter&account_id=all&filter=all&limit=5') as res:
        arts = json.loads(res.read().decode())
        for a in arts:
            print(' ', a.get('author', {}).get('handle'), a.get('id'), a.get('full_text', '')[:40])

    print('\n--- API: subyike timeline (top 5) ---')
    with urllib.request.urlopen('http://localhost:5173/api/timeline?platform=twitter&account_id=subyike&filter=all&limit=5') as res:
        arts = json.loads(res.read().decode())
        for a in arts:
            print(' ', a.get('author', {}).get('handle'), a.get('id'), a.get('full_text', '')[:40])

    print('\n--- API: subyike reposts timeline (checking the 2 fixed posts) ---')
    with urllib.request.urlopen('http://localhost:5173/api/timeline?platform=twitter&account_id=subyike&filter=reposts&limit=10') as res:
        arts = json.loads(res.read().decode())
        for a in arts:
            if a.get('id') in ('2028816874549707202', '2012531768629342472'):
                rt = a.get('retweeted_article')
                print('  Found Fixed RT:', a.get('id'))
                print('    rt author:', rt.get('author', {}).get('handle') if rt else None)
                print('    rt text:', rt.get('full_text', '')[:40] if rt else None)
                print('    rt media count:', len(rt.get('media', [])) if rt else 0)
                if rt and rt.get('media'):
                    print('    rt media preview URL:', rt.get('media')[0].get('preview_url')[:60])

if __name__ == '__main__':
    main()
