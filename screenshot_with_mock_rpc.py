"""
screenshot_with_mock_rpc.py
ListAllAccounts および GetAccountDetail のモックを注入し、
Admin Governance Portalのアカウント詳細（アバター・自己紹介）を完璧に撮影する
"""
import time, os, sqlite3, json
from seleniumbase import SB

SCREENSHOT_DIR = os.path.join(os.path.dirname(__file__), "screenshots")
os.makedirs(SCREENSHOT_DIR, exist_ok=True)
DB_PATH = os.path.join(os.path.dirname(__file__), "archive.db")

def get_db_data():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    cur = conn.cursor()
    accounts = [dict(r) for r in cur.execute("SELECT * FROM accounts ORDER BY username ASC").fetchall()]
    histories = [dict(r) for r in cur.execute("SELECT * FROM account_profile_histories").fetchall()]
    conn.close()
    return accounts, histories

def main():
    accounts, histories = get_db_data()
    print(f"[*] Loaded {len(accounts)} accounts and {len(histories)} histories from DB.")

    with SB(uc=False, headless=True) as sb:
        sb.open("http://localhost:5173")
        time.sleep(2)

        # Wails RPC のモックを注入
        mock_script = f"""
        const mockAccounts = {json.dumps(accounts)};
        const mockHistories = {json.dumps(histories)};

        window.go = window.go || {{}};
        window.go.app = window.go.app || {{}};
        window.go.app.App = {{
            ListAllAccounts: async () => mockAccounts,
            GetAccountDetail: async (numericId) => {{
                const acc = mockAccounts.find(a => a.numeric_id === numericId) || mockAccounts[0];
                const hists = mockHistories.filter(h => h.account_id === acc.numeric_id || h.account_id === acc.username);
                return {{
                    account: acc,
                    histories: hists,
                    post_count: 29
                }};
            }},
            GetAvailableAvatars: async () => [],
            UpdateAccount: async () => true,
            UpdateAvatarBase64: async () => true,
        }};
        """
        sb.execute_script(mock_script)
        print("[+] Injected Wails RPC mock with ListAllAccounts.")

        # 1. 設定・管理ボタンをクリック
        sb.click("button:contains('設定・管理')")
        time.sleep(2)

        # 2. サイドバーのアカウントタブをクリック
        sb.execute_script("""
            const buttons = Array.from(document.querySelectorAll('aside button'));
            const accBtn = buttons.find(b => b.textContent.includes('アカウント'));
            if (accBtn) accBtn.click();
        """)
        time.sleep(2)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "03_admin_accounts_list.png"))

        # 3. Yike_Luo をクリックして詳細を表示
        sb.execute_script("""
            const rows = Array.from(document.querySelectorAll('tbody tr, .overflow-y-auto tr'));
            const yike = rows.find(r => r.textContent.includes('Yike_Luo'));
            if (yike) yike.click();
        """)
        time.sleep(3)
        sb.save_screenshot(os.path.join(SCREENSHOT_DIR, "04_yike_luo_detail.png"))
        print("[+] Done! Captured updated screenshots.")

if __name__ == "__main__":
    main()
