from warcio.archiveiterator import ArchiveIterator
import re

p = r'backups/dumps\twitter\Nenne1001\2096723868694016024\snapshot.warc.gz'
with open(p, 'rb') as s:
    for r in ArchiveIterator(s):
        print(r.rec_type, r.rec_headers.get_header('WARC-Target-URI'))
        if r.rec_type == 'response':
            content = r.content_stream().read().decode('utf-8', errors='ignore')
            print('Length:', len(content))
            m_text = re.findall(r'<p[^>]*class="[^"]*tweet-text[^"]*"[^>]*>(.*?)</p>', content, re.DOTALL)
            print('Text match:', m_text)
            print('Description:', re.findall(r'<meta[^>]+description[^>]+>', content, re.I))
            print('Title:', re.findall(r'<title>(.*?)</title>', content))
