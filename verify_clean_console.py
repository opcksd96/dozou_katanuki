"""
verify_clean_console.py
http://127.0.0.1:5175/ をブラウザで開き、コンソールエラーの有無とスキンCSSのロードを検証
"""
import time, os, json
from seleniumbase import SB

SCREENSHOT_DIR = os.path.join(os.path.dirname(__file__), "screenshots")

def main():
    print("[*] Opening http://127.0.0.1:5175/ in browser...")
    with SB(uc=False, headless=True) as sb:
        sb.open("http://127.0.0.1:5175/")
        time.sleep(5)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "broadcast_clean_console.png"))

        # コンソールログの取得
        logs = sb.driver.get_log("browser")
        print(f"[*] Console Logs Count: {len(logs)}")
        twimg_errors = [l for l in logs if "twimg.com" in l.get("message", "")]
        css_errors = [l for l in logs if "design.css" in l.get("message", "")]
        print(f"[+] twimg.com 404 errors: {len(twimg_errors)}")
        print(f"[+] design.css 404 errors: {len(css_errors)}")

        for log in logs:
            if log.get("level") in ("SEVERE", "WARNING"):
                print("  [WARN/ERR]", log.get("message"))

        # アバター画像のsrc属性調査
        img_srcs = sb.execute_script("""
            return Array.from(document.querySelectorAll('img')).map(img => ({
                srcPrefix: img.src.substring(0, 60),
                isDataUri: img.src.startsWith('data:image/'),
                isExternal: img.src.startsWith('http') && !img.src.includes('127.0.0.1') && !img.src.includes('localhost'),
                naturalWidth: img.naturalWidth
            }));
        """)
        print("[*] Rendered Images Summary:")
        external_count = sum(1 for img in img_srcs if img["isExternal"])
        print(f"  Total Images: {len(img_srcs)}")
        print(f"  External Image URLs: {external_count}")

if __name__ == "__main__":
    main()
