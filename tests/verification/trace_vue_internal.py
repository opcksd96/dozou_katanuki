"""
trace_vue_internal.py
Vueコンポーネントツリー内部のリアクティブ状態 (articles, loading, error等) を抽出・診断
"""
import time, os, json
from seleniumbase import SB

def main():
    with SB(uc=False, headless=True) as sb:
        sb.open("http://127.0.0.1:5175/")
        time.sleep(3)

        info = sb.execute_script("""
            // window.__VUE_DEVTOOLS_GLOBAL_HOOK__ または DOM から Vue の状態を調査
            const appEl = document.querySelector('#app');
            const vueApp = appEl?.__vue_app__;
            
            return {
                appExists: !!appEl,
                hasVueApp: !!vueApp,
                bodyHtml: document.querySelector('main')?.innerHTML?.substring(0, 500)
            };
        """)
        print("[*] Vue App Info:", info)

if __name__ == "__main__":
    main()
