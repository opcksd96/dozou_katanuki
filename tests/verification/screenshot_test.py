"""
screenshot_test.py
サイドバーの「アカウント」タブを確実にクリックしてアカウント一覧・詳細画面を撮影
"""
import time, os
from seleniumbase import SB

SCREENSHOT_DIR = os.path.join(os.path.dirname(__file__), "screenshots")
os.makedirs(SCREENSHOT_DIR, exist_ok=True)

def main():
    print("[*] Starting browser session to verify UI...")
    with SB(uc=False, headless=True) as sb:
        # 1. タイムライン画面
        sb.open("http://localhost:5173")
        time.sleep(3)

        # 2. 設定・管理ボタンをクリック
        sb.click("button:contains('設定・管理')")
        time.sleep(2)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "02_admin_modal.png"))

        # 3. アカウントタブをJSで直接クリック
        sb.execute_script("""
            const buttons = Array.from(document.querySelectorAll('aside button'));
            const accBtn = buttons.find(b => b.textContent.includes('アカウント'));
            if (accBtn) accBtn.click();
        """)
        time.sleep(3)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "03_admin_accounts_list.png"))

        # 4. Yike_Luo の行をJSでクリック
        sb.execute_script("""
            const rows = Array.from(document.querySelectorAll('tr, div, td'));
            const yikeRow = rows.find(r => r.textContent.includes('Yike_Luo') || r.textContent.includes('ext_Yike_Luo'));
            if (yikeRow) yikeRow.click();
        """)
        time.sleep(2)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "04_yike_luo_detail.png"))
        print("[+] Done capturing all states!")

if __name__ == "__main__":
    main()
