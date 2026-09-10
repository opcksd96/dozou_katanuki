# scratch/get_media_schema.py
import sqlite3

def main():
    conn = sqlite3.connect('archive.db')
    cur = conn.cursor()
    cur.execute("SELECT sql FROM sqlite_master WHERE type = 'table' AND name = 'media'")
    print(cur.fetchone()[0])

if __name__ == '__main__':
    main()
