# scratch/find_unlinked_rts.py
import sqlite3

def main():
    conn = sqlite3.connect('archive.db')
    cur = conn.cursor()

    cur.execute("SELECT id, full_text FROM articles WHERE account_id = 'ext_subyike_YIKE' AND is_repost = 1 AND (reply_to_id IS NULL OR reply_to_id = '')")
    unlinked = cur.fetchall()
    print(f'Unlinked RTs count: {len(unlinked)}')
    for u in unlinked:
        print(f'  RT: {u[0]} | {u[1][:70]}')

    cur.execute("SELECT id, full_text FROM articles WHERE account_id = 'ext_subyike_YIKE' AND is_repost = 0")
    non_rts = cur.fetchall()
    print(f'\nTotal non-RTs: {len(non_rts)}')
    for n in non_rts:
        if n[0] in ('2028459277254602771', '2012506627216453640'):
            print(f'  [TARGET] {n[0]} | {n[1][:70]}')

if __name__ == '__main__':
    main()
