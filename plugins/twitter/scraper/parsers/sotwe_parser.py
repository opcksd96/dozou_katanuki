"""
sotwe_parser.py
Sotwe HTMLからアカウントBio、アバター、リツイート、全メディアを完全抽出するパーサー。
"""
from datetime import datetime, timezone
import re
from typing import Any, Dict, List
from bs4 import BeautifulSoup


def parse_sotwe_html_tweets(html_str: str, default_account: str) -> List[Dict[str, Any]]:
    """SotweのHTMLからアカウントBioおよびツイート一覧を抽出して標準スキーマ化"""
    if not html_str:
        return []

    soup = BeautifulSoup(html_str, "html.parser")

    # 1. ページ全体の共通プロフィール（Bio・表示名・アバター）をヘッダーから抽出
    page_bio = ""
    bio_el = soup.select_one(".break-word .dynamic-link-content, div[data-v-3cd48e78] .dynamic-link-content")
    if bio_el:
        page_bio = bio_el.get_text(separator="\n").strip()

    header_name_el = soup.select_one(".profile-name, .v-card__title .font-weight-bold")
    page_display_name = header_name_el.get_text(strip=True) if header_name_el else default_account

    header_avatar_el = soup.select_one(".profile-avatar img, .v-avatar img")
    page_avatar_url = ""
    if header_avatar_el and header_avatar_el.get("src"):
        raw_av = header_avatar_el["src"]
        page_avatar_url = raw_av.replace("_normal.", "_400x400.") if "_normal." in raw_av else raw_av

    cards = soup.select("div.tweet-card")
    results = []

    for card in cards:
        # 2. リツイート判定
        rt_badge = card.select_one(".v-card__title .fa-retweet")
        is_repost = bool(rt_badge)
        retweeted_by = ""
        if is_repost:
            rt_user_el = card.select_one(".v-card__title a[href^='/']")
            if rt_user_el and rt_user_el.get("href"):
                retweeted_by = rt_user_el["href"].strip("/").split("/")[0]

        # 3. 投稿者アカウント情報の解決
        profile_link = card.select_one(".tweet-profile a[href^='/']")
        author_username = default_account
        if profile_link and profile_link.get("href"):
            author_username = profile_link["href"].strip("/").split("/")[0]

        name_el = card.select_one(".tweet-profile--text span.font-weight-medium")
        display_name = name_el.get_text(strip=True) if name_el else (page_display_name if author_username == default_account else author_username)

        avatar_img = card.select_one(".v-avatar img, .tweet-profile img")
        avatar_url = ""
        if avatar_img and avatar_img.get("src"):
            raw_av = avatar_img["src"]
            avatar_url = raw_av.replace("_normal.", "_400x400.") if "_normal." in raw_av else raw_av
        elif author_username == default_account:
            avatar_url = page_avatar_url

        # 本人ポストの場合はヘッダーのBioをバインド
        author_bio = page_bio if author_username.lower() == default_account.lower() else ""

        # 4. 投稿日時
        time_el = card.select_one("time[datetime]")
        created_at_iso = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        if time_el and time_el.get("datetime"):
            created_at_iso = time_el["datetime"]

        # 5. 本文
        text_el = card.select_one(".tweet-text .dynamic-link-content")
        full_text = text_el.get_text(separator="\n").strip() if text_el else ""

        # 6. カルーセルメディア抽出 (非表示スライド含む全画像・動画)
        media_list = []
        for img in card.select(".media-carousel img[src], .media-carousel-image img[src]"):
            src = img.get("src", "")
            if src and "profile_images" not in src and not any(m["url"] == src for m in media_list):
                media_list.append({"url": src, "type": "image", "width": 0, "height": 0})

        for vid_src in card.select("video source[src]"):
            src = vid_src.get("src", "")
            if src and not any(m["url"] == src for m in media_list):
                media_list.append({"url": src, "type": "video", "width": 0, "height": 0})

        # 7. ポストID
        post_id = ""
        for m in media_list:
            match = re.search(r"/(?:media|ext_tw_video(?:_thumb)?)/(\d+)", m["url"])
            if match:
                post_id = match.group(1)
                break
        if not post_id and time_el and time_el.get("datetime"):
            try:
                dt = datetime.fromisoformat(time_el["datetime"].replace("Z", "+00:00"))
                post_id = str(int(dt.timestamp() * 1000))
            except Exception:
                post_id = str(int(datetime.now(timezone.utc).timestamp()))

        results.append({
            "platform": "twitter",
            "account": {
                "numeric_id": f"ext_{author_username}",
                "username": author_username,
                "display_name": display_name,
                "avatar_url": avatar_url,
                "avatar_original_url": avatar_url,
                "description": author_bio,
            },
            "post": {
                "id": post_id,
                "conversation_id": post_id,
                "reply_to_tweet_id": None,
                "reply_to_handle": None,
                "created_at": created_at_iso,
                "full_text": full_text,
                "via": "Sotwe",
                "is_repost": is_repost,
                "retweeted_by": retweeted_by,
                "wayback_url": f"https://sotwe.com/{author_username}/status/{post_id}",
                "urls": []
            },
            "media": media_list
        })

    return results
