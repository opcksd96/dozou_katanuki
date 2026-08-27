"""
debug_timeline_vue.py
ブラウザ上で useTimeline の articles の状態や fetch の結果を直接コンソール出力して診断
"""
import time, os, json
from seleniumbase import SB

def main():
    with SB(uc=False, headless=True) as sb:
        sb.open("http://127.0.0.1:5175/")
        time.sleep(3)

        # fetch を手動実行して結果をログ
        result = sb.execute_script("""
            return (async () => {
                const res = await fetch('/api/timeline?platform=twitter&account_id=all&filter=all&limit=50&offset=0');
                const data = await res.json();
                return {
                    status: res.status,
                    isArray: Array.isArray(data),
                    itemCount: Array.isArray(data) ? data.length : 0,
                    firstItem: Array.isArray(data) && data.length > 0 ? {
                        id: data[0].id,
                        author: data[0].author,
                        body: data[0].body
                    } : null
                };
            })();
        """)
        print("[*] Direct Browser Fetch Result:")
        print(json.dumps(result, indent=2, ensure_ascii=False))

        # 再読み込みボタンをクリック
        sb.click("button:contains('再読み込み')")
        time.sleep(3)

        articles_rendered = sb.execute_script("""
            return {
                cards: document.querySelectorAll('.tweet-card, .article-card, article').length,
                innerText: document.body.innerText.substring(0, 400)
            };
        """)
        print("[*] Rendered Articles Count:", articles_rendered)

if __name__ == "__main__":
    main()
