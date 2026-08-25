"""
sync_real_avatar.py
Sotwe / Twitter の本物のアバター画像 (https://pbs.twimg.com/profile_images/1730257152630312960/NGzvfvD0.jpg) を
取得して Base64 にエンコードし、DBの accounts & account_profile_histories に反映するスクリプト
"""
import sqlite3, requests, base64

DB_PATH = "archive.db"
AVATAR_URL = "https://pbs.twimg.com/profile_images/1730257152630312960/NGzvfvD0.jpg"

def main():
    print(f"[*] Fetching real avatar from {AVATAR_URL}...")
    headers = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"}
    try:
        resp = requests.get(AVATAR_URL, headers=headers, timeout=10)
        if resp.status_code == 200:
            b64_str = f"data:image/jpeg;base64,{base64.b64encode(resp.content).decode('utf-8')}"
            print(f"[+] Downloaded avatar successfully ({len(resp.content)} bytes).")
            
            conn = sqlite3.connect(DB_PATH)
            cur = conn.cursor()
            cur.execute("""
                UPDATE accounts 
                SET avatar_base64 = ?, avatar_url = ?, description = ? 
                WHERE username = 'Yike_Luo' OR numeric_id = 'ext_Yike_Luo'
            """, (b64_str, AVATAR_URL, "大号已封，且看且珍惜。🚪200加tg瑟瑟群及vx，入门可定制视频和线下，约拍请私信。我的小号请尽早关注@TeacherXiaoLuo"))
            
            cur.execute("""
                UPDATE account_profile_histories 
                SET avatar_base64 = ?, avatar_original_url = ? 
                WHERE account_id LIKE '%Yike_Luo%'
            """, (b64_str, AVATAR_URL))
            
            conn.commit()
            conn.close()
            print("[+] DB updated with real avatar Base64 and description!")
        else:
            print(f"[-] HTTP error: {resp.status_code}")
    except Exception as e:
        print(f"[-] Error fetching avatar: {e}")

if __name__ == "__main__":
    main()
