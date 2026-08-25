"""
inspect_avatar_dom.py
ブラウザ上でアバター要素のDOM状態、src属性、エラー内容を精密調査
"""
import time, os, sqlite3, json
from seleniumbase import SB

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
    yike = next(a for a in accounts if "Yike_Luo" in a["username"])
    print("[*] Yike_Luo DB record:", dict(yike))

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

        # 1. 設定・管理ボタンをクリック
        sb.click("button:contains('設定・管理')")
        time.sleep(2)

        # 2. アカウントタブをクリック
        sb.execute_script("""
            const buttons = Array.from(document.querySelectorAll('aside button'));
            const accBtn = buttons.find(b => b.textContent.includes('アカウント'));
            if (accBtn) accBtn.click();
        """)
        time.sleep(2)

        # 3. Yike_Luo をクリック
        sb.execute_script("""
            const rows = Array.from(document.querySelectorAll('tbody tr, .overflow-y-auto tr'));
            const yike = rows.find(r => r.textContent.includes('Yike_Luo'));
            if (yike) yike.click();
        """)
        time.sleep(2)

        # 4. DOMの状態とコンソールエラーを取得
        diag = sb.execute_script("""
            const imgs = Array.from(document.querySelectorAll('img, svg')).map(el => ({
                tag: el.tagName,
                src: el.getAttribute('src') || '',
                srcPrefix: (el.getAttribute('src') || '').substring(0, 50),
                display: window.getComputedStyle(el).display,
                className: el.className,
                naturalWidth: el.naturalWidth || 0,
                naturalHeight: el.naturalHeight || 0
            }));
            return {
                images: imgs,
                selectedDetail: document.querySelector('.space-y-4')?.innerHTML?.substring(0, 300)
            };
        """)
        print("[*] DOM Diagnostics:")
        print(json.dumps(diag, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    main()
