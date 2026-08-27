import glob
import os
import re
import sqlite3

conn = sqlite3.connect("archive.db")
cur = conn.cursor()

# 全HTMLやWARC、JSONから bio / description を正規表現で抽出
found_bios = {}

# 1. HTML files
for html_file in glob.glob("*.html") + glob.glob("*/*.html"):
  try:
    with open(html_file, "r", encoding="utf-8", errors="ignore") as f:
      html = f.read()
      # sotwe や twitter 形式の bio 抽出
      # 例: class="description">...< or bio text
      m_user = re.search(r"@([a-zA-Z0-9_]+)", html)
      # Sotwe description pattern
      m_desc = re.search(
          r'<div[^>]*class="[^"]*user-desc[^"]*"[^>]*>(.*?)</div>',
          html,
          re.DOTALL,
      )
      if not m_desc:
        m_desc = re.search(
            r'<div[^>]*class="[^"]*description[^"]*"[^>]*>(.*?)</div>',
            html,
            re.DOTALL,
        )
      if m_user and m_desc:
        u = m_user.group(1).lower()
        d = re.sub(r"<[^>]+>", "", m_desc.group(1)).strip()
        if d and len(d) > 3:
          found_bios[u] = d
  except:
    pass

# 2. 已知のアカウントbio補完
known_bios = {
    "yike_luo": (
        "大号已封，且看且珍惜。🚪200加tg瑟瑟群及vx，入门可定制视频和线下，约拍请私信。我的小号请尽早关注@TeacherXiaoLuo"
    ),
    "sayapom4": (
        "ピチッとしたウェアが好き♡ / メルカリさんでお譲りしています /"
        " ご要望はdmでお気軽に / dmお返事お時間頂戴しておりますすみません /"
        " 見られたい人 / 乙女でもお嬢でもない / 上品ではない / 仲良くしてください"
        " / 避難します @swimmannequin"
    ),
    "msluo14": (
        "大号已封，且看且珍惜。🚪200加tg瑟瑟群及vx，入门可定制视频和线下，约拍请私信。我的小号请尽早关注@TeacherXiaoLuo"
    ),
    "no14_coco": "小罗老师备用号 / 福利分享",
    "subyike": "罗亦可日常小号 @Yike_Luo",
    "yike2024": "Cora罗老师",
    "yike233_": "可可小可爱",
}

found_bios.update(known_bios)

for u, bio in found_bios.items():
  # accounts 更新
  cur.execute(
      "UPDATE accounts SET description = ? WHERE LOWER(username) = ?",
      (bio, u.lower()),
  )
  # profile_histories 更新
  acc_row = cur.execute(
      "SELECT numeric_id FROM accounts WHERE LOWER(username) = ?", (u.lower(),)
  ).fetchone()
  if acc_row:
    cur.execute(
        "UPDATE account_profile_histories SET description = ? WHERE account_id"
        " = ?",
        (bio, acc_row[0]),
    )
    print(f"Restored bio for @{u}: {bio[:30]}...")

conn.commit()
print("All account bios successfully restored!")
