"""
verify_broadcast_timeline.py
http://127.0.0.1:5175/ をブラウザで開き、記事カード・アバター・タイムラインの描画を検証・撮影
"""
import time, os, json
from seleniumbase import SB

SCREENSHOT_DIR = os.path.join(os.path.dirname(__file__), "screenshots")
os.makedirs(SCREENSHOT_DIR, exist_ok=True)

def main():
    print("[*] Opening http://127.0.0.1:5175/ in browser...")
    with SB(uc=False, headless=True) as sb:
        sb.open("http://127.0.0.1:5175/")
        time.sleep(5)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "broadcast_5175_fixed.png"))

        articles = sb.execute_script("return document.querySelectorAll('article').length;")
        print(f"[+] Rendered articles count: {articles}")

        diag = sb.execute_script("""
            const arts = Array.from(document.querySelectorAll('article')).map(a => ({
                id: a.getAttribute('data-article-id') || '',
                handle: a.querySelector('.twitter-handle, a[href^=\"/\"]')?.textContent || '',
                text: a.querySelector('.article-body, .twitter-body')?.textContent?.substring(0, 50) || '',
                hasImg: a.querySelectorAll('img').length
            }));
            return {
                totalArticles: arts.length,
                firstThree: arts.slice(0, 3)
            };
        """)
        print("[*] Diagnostic data:")
        print(json.dumps(diag, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    main()
