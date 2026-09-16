import sys, os
_ROOT = r"d:\Projects\10_tools\dozou_katanuki\dozou_katanuki"
if _ROOT not in sys.path: sys.path.insert(0, _ROOT)

from plugins.twitter.scraper.parsers.sotwe_parser import parse_sotwe_vue_tweets, parse_sotwe_html_tweets
from plugins.twitter.scraper.parsers.sotwe_extractors import normalize_vue_tweet

# Simulate Vue object with 4 images from Twitter/Sotwe
vue_multi = {
    "id": "1907698917766213981",
    "createdAt": 1743665839000,
    "text": "我是发情母狗，不是痴傻少女！！！\n@Daguidiyi @Xiaoxiaofoer https://x.com/MsLuo14/status/1907698917766213981/photo/1",
    "user": {"screenName": "MsLuo14", "name": "小罗老师"},
    "extendedEntities": {
        "media": [
            {"type": "image", "mediaURL": "https://pbs.twimg.com/media/GnmCCzebYAAkg9-.jpg"},
            {"type": "image", "mediaURL": "https://pbs.twimg.com/media/GnmCCzqbUAALdDB.jpg"},
            {"type": "image", "mediaURL": "https://pbs.twimg.com/media/GnmCCzgbQAAkIG2.jpg"},
            {"type": "image", "mediaURL": "https://pbs.twimg.com/media/GnmCCzlaMAIPwFr.jpg"}
        ]
    }
}

parsed_vue = normalize_vue_tweet(vue_multi, "MsLuo14")
print("=== PARSED VUE MULTI-IMAGE ===")
print("Media Count:", len(parsed_vue["media"]))
for m in parsed_vue["media"]:
    print(" ", m["media_id"], m["download_url"])

# Simulate HTML card with 4 images in carousel
html_multi = """
<div class="tweet-card">
    <div class="tweet-profile"><a href="/MsLuo14">小罗老师</a></div>
    <div class="tweet-text"><div class="dynamic-link-content">我是发情母狗，不是痴傻少女！！！</div></div>
    <div class="media-carousel">
        <div class="v-window-item"><img class="img-content" src="https://pbs.twimg.com/media/GnmCCzebYAAkg9-.jpg"></div>
        <div class="v-window-item" style="background-image: url('https://pbs.twimg.com/media/GnmCCzqbUAALdDB.jpg')"></div>
        <div class="v-window-item"><a href="https://pbs.twimg.com/media/GnmCCzgbQAAkIG2.jpg">link</a></div>
        <div class="v-window-item" data-src="https://pbs.twimg.com/media/GnmCCzlaMAIPwFr.jpg"></div>
    </div>
</div>
"""
parsed_html = parse_sotwe_html_tweets(html_multi, "MsLuo14")
print("\n=== PARSED HTML CAROUSEL MULTI-IMAGE ===")
print("Media Count:", len(parsed_html[0]["media"]))
for m in parsed_html[0]["media"]:
    print(" ", m["media_id"], m["download_url"])
