import gzip

warc_path = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki\backups\dumps\twitter\msluo14\1907698917766213981\snapshot.warc.gz"

with gzip.open(warc_path, 'rb') as f:
    content = f.read()
    print("WARC Total Size:", len(content))
    # print first 2000 chars decoded
    text = content.decode('utf-8', errors='ignore')
    print("=== FIRST 3000 CHARS ===")
    print(text[:3000])
