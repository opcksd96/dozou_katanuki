import sqlite3

conn = sqlite3.connect("backups/database/archive_20260831_082956.db")
conn.row_factory = sqlite3.Row
c = conn.cursor()

c.execute("SELECT * FROM accounts WHERE is_trash = 1 OR numeric_id IN ('24d22341-2a66-5b39-9268-5542c824fcaf', '937fdf1b-524d-53b5-8b88-3714ddd275e0', 'f35dad07-7066-5405-bb36-7e50e2d94434')")
rows = [dict(r) for r in c.fetchall()]
print(f"Found {len(rows)} trashed in Aug 31 db:")
for r in rows:
    print(r)

conn2 = sqlite3.connect("backups/database/archive_snapshot_20260907_205625.db")
conn2.row_factory = sqlite3.Row
c2 = conn2.cursor()
c2.execute("SELECT * FROM accounts WHERE numeric_id IN ('73a4ee3d-f852-5097-89fc-e10cf2269d21', 'bd8c2c8e-5e45-52fb-897b-9026ff76c4b0')")
rows2 = [dict(r) for r in c2.fetchall()]
print(f"Found {len(rows2)} in Sep 7 db:")
for r in rows2:
    print(r)
