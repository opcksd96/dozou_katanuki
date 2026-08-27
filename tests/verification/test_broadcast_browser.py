"""
test_broadcast_browser.py
ブラウザで http://127.0.0.1:5175/ を開き、コンソールエラー、DOMツリー、ネットワークエラーを調査・スクリーンショット撮影
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
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "broadcast_5175.png"))

        # コンソールログの取得
        try:
            logs = sb.driver.get_log("browser")
            print("[*] Browser Console Logs:")
            for log in logs:
                print("  ", log)
        except Exception as e:
            print("[-] Could not get console logs:", e)

        # DOM状態の診断
        diag = sb.execute_script("""
            return {
                title: document.title,
                bodyText: document.body.innerText.substring(0, 500),
                totalDivs: document.querySelectorAll('div').length,
                articles: document.querySelectorAll('article').length,
                timelineItems: document.querySelectorAll('.tweet-card, .article-card, div[data-article-id]').length,
                images: Array.from(document.querySelectorAll('img')).map(img => ({
                    src: img.src.substring(0, 80),
                    display: window.getComputedStyle(img).display,
                    naturalWidth: img.naturalWidth
                }))
            };
        """)
        print("[*] DOM State:")
        print(json.dumps(diag, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    main()
