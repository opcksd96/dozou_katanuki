import json
import sys

sys.path.insert(0, "plugins/base/scraper")
from core.stash_client import StashClient

stash = StashClient()

# 1. 1920037769507676376 のシーン
r1 = stash.query("""query {
  findScenes(filter: { q: "1920037769507676376" }) {
    scenes { id title details date urls custom_fields }
  }
}""")
print("Target Scene 1920037769507676376:")
print(json.dumps(r1, indent=2, ensure_ascii=False))

# 2. X (@...) のタイトルを持つシーン
r2 = stash.query("""query {
  findScenes(scene_filter: { title: { value: "X (@", modifier: INCLUDES } }, filter: { per_page: 3 }) {
    scenes { id title details date urls custom_fields }
  }
}""")
print("\nX Twitter Scenes sample:")
print(json.dumps(r2, indent=2, ensure_ascii=False))

# 3. X (@...) のタイトルを持つ画像
r3 = stash.query("""query {
  findImages(image_filter: { title: { value: "X (@", modifier: INCLUDES } }, filter: { per_page: 3 }) {
    images { id title details date urls custom_fields }
  }
}""")
print("\nX Twitter Images sample:")
print(json.dumps(r3, indent=2, ensure_ascii=False))
