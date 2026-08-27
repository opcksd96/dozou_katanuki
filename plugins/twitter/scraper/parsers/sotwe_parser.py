"""
sotwe_parser.py (SPEC-PLUGIN-001 / 100行以下)
Sotwe HTMLからアカウントBio、アバター、リツイート、全メディアを完全抽出するパーサー。
"""
from typing import Any, Dict, List
from bs4 import BeautifulSoup
try:
    from plugins.twitter.scraper.parsers.sotwe_extractors import (
        parse_iso_datetime, extract_metrics, extract_card_media, extract_status_id_from_media
    )
except ImportError:
    from parsers.sotwe_extractors import (
        parse_iso_datetime, extract_metrics, extract_card_media, extract_status_id_from_media
    )

def parse_sotwe_html_tweets(html_str: str, default_account: str) -> List[Dict[str, Any]]:
    """SotweのHTMLからアカウントBioおよびツイート一覧を抽出して標準スキーマ化"""
    if not html_str: return []
    soup = BeautifulSoup(html_str, "html.parser")

    # 1. ページ共通プロフィール（Bio・表示名・アバター）
    bio_el = soup.select_one(".break-word .dynamic-link-content, div[data-v-3cd48e78] .dynamic-link-content")
    page_bio = bio_el.get_text(separator="\n").strip() if bio_el else ""
    header_name_el = soup.select_one(".profile-name, .v-card__title .font-weight-bold")
    page_name = header_name_el.get_text(strip=True) if header_name_el else default_account

    header_av_el = soup.select_one(".profile-avatar img, .v-avatar img")
    page_av = header_av_el["src"].replace("_normal.", ".") if header_av_el and header_av_el.get("src") else ""

    cards = soup.select("div.tweet-card")
    results = []

    for card in cards:
        # 2. リツイート判定 & ピン留め判定
        is_repost = bool(card.select_one(".v-card__title .fa-retweet"))
        is_pinned = bool(card.select_one(".pinned-text, .pinned-icon"))
        rt_user_el = card.select_one(".v-card__title a[href^='/']")
        retweeted_by = rt_user_el["href"].strip("/").split("/")[0] if (is_repost and rt_user_el and rt_user_el.get("href")) else ""

        # 3. 投稿者アカウント情報の解決
        profile_link = card.select_one(".tweet-profile a[href^='/']")
        author_user = profile_link["href"].strip("/").split("/")[0] if (profile_link and profile_link.get("href")) else default_account
        name_el = card.select_one(".tweet-profile--text span.font-weight-medium")
        display_name = name_el.get_text(strip=True) if name_el else (page_name if author_user == default_account else author_user)

        av_img = card.select_one(".v-avatar img, .tweet-profile img")
        raw_av = av_img["src"] if (av_img and av_img.get("src")) else (page_av if author_user == default_account else "")
        avatar_url = raw_av.replace("_normal.", ".") if "_normal." in raw_av else raw_av
        author_bio = page_bio if author_user.lower() == default_account.lower() else ""

        # 4. 投稿日時（実際の投稿日時を確実に取得）
        time_el = card.select_one("time[datetime]")
        created_at_iso = parse_iso_datetime(time_el.get("datetime") if time_el else None)

        # 5. 本文・メディア・メトリクス・ID
        text_el = card.select_one(".tweet-text .dynamic-link-content")
        full_text = text_el.get_text(separator="\n").strip() if text_el else ""
        media_list = extract_card_media(card)
        metrics = extract_metrics(card)
        post_id = extract_status_id_from_media(media_list, created_at_iso)

        results.append({
            "platform": "twitter",
            "account": {
                "numeric_id": f"ext_{author_user}",
                "username": author_user,
                "display_name": display_name,
                "avatar_url": avatar_url,
                "avatar_original_url": avatar_url,
                "description": author_bio,
                "profile_history": [{
                    "display_name": display_name,
                    "avatar_original_url": avatar_url,
                    "observed_at": created_at_iso
                }] if avatar_url else []
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
                "is_pinned": is_pinned,
                "retweeted_by": retweeted_by,
                "wayback_url": f"https://sotwe.com/{author_user}/status/{post_id}",
                "metrics": metrics,
                "urls": []
            },
            "media": media_list
        })

    return results
