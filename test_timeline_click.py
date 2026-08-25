"""
test_timeline_click.py
ブラウザ上で「🌐 全てのアカウント」ボタンをクリックしてタイムラインが描画されるかを検証
"""
import time, os
from seleniumbase import SB

SCREENSHOT_DIR = os.path.join(os.path.dirname(__file__), "screenshots")

def main():
    with SB(uc=False, headless=True) as sb:
        sb.open("http://127.0.0.1:5175/")
        time.sleep(3)

        # 1. 「全てのアカウント」ボタンをクリック
        print("[*] Clicking '全てのアカウント' button...")
        sb.click("button:contains('全てのアカウント')")
        time.sleep(3)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "after_click_all.png"))

        # 2. 記事カードの数を確認
        count = sb.execute_script("return document.querySelectorAll('article').length;")
        print(f"[*] Articles count after click: {count}")

        # 3. Yike_Luo をクリック
        sb.click("button:contains('yike_luo')")
        time.sleep(3)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "after_click_yike.png"))
        count_yike = sb.execute_script("return document.querySelectorAll('article').length;")
        print(f"[*] Articles count for yike_luo: {count_yike}")

if __name__ == "__main__":
    main()
