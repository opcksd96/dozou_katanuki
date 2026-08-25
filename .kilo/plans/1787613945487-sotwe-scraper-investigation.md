# Sotwe Scraper 実装改善計画

**事実確認**: ドキュメント (`COMBINED_WIKI.md`, `SPEC-SCRAPER-EXTERNAL-001`) に記載の「Sotwe JSON API (`api.sotwe.com/api/v2`)」は**存在しない**。  
**正規実装**: `sotwe_source.py` + `sotwe_parser.py` の現状（SeleniumBase UC Mode + Web UI DOMスクレイピング）を正とする。

---

## 1. ユーザー確認済み事実

### アバターURL
- `_400x400.` サフィックス付与 → **404**（存在しない）
- `_normal.` を削除 → **オリジナルサイズのURLが取得可能**（確認済み）

### 動画要素
存在するパターン：
```html
<div class="video-player absolute fill-height w-full">
  <video aria-label="Video player" preload="none" 
         poster="https://pbs.twimg.com/ext_tw_video_thumb/1749123239694266369/pu/img/DKQeLTOG6J9TmCL7.jpg"
         controls="controls" playsinline="true">
    <source src="https://video.twimg.com/ext_tw_video/1749123239694266369/pu/vid/avc1/720x1280/x7vQf7QsgEONy7TF.mp4?tag=12"
            type="video/mp4">
    <p>Yike_Luo's tweet video.</p>
  </video>
</div>
```

---

## 2. 問題①：アバター第1世代の画像データ伝搬

### 現状の問題
- `avatar_url`: `_normal.` → `_400x400.` 変換 → **404**
- `avatar_original_url`: `avatar_url` と同じ値 → **404**

### 修正方針
`_400x400.` は404が確認されているため使用不可。`_normal.` を削除したオリジナルサイズURLを両方に設定する。

| フィールド | 値 | 根拠 |
|:---|:---|:---|
| `avatar_url` | `_normal.` を削除したURL | 実在するURL |
| `avatar_original_url` | `_normal.` を削除したURL | 実在するURL（originalと同じ） |

**実装修正** (`sotwe_parser.py`):
```python
if "_normal." in raw_av:
    avatar_url = raw_av.replace("_normal.", "")
    avatar_original_url = raw_av.replace("_normal.", "")
else:
    avatar_url = raw_av
    avatar_original_url = raw_av
```

---

## 3. 問題②：メディアURLのスクレイピング精度

### 3.1 画像メディア
- 現在: `.media-carousel img[src], .media-carousel-image img[src]` で抽出
- 画像URL自体は正しく抽出できる
- `width`/`height` は常に 0（HTMLにサイズ情報がないため）
- 修正不要（制限事項として明記）

### 3.2 動画メディア（新規抽出対象）

**現在の実装**: `video source[src]` ロジックが存在するが、親要素の特定が不正確な可能性がある。

**修正方針**:
1. `video.video-player` または `.video-player video` を探す
2. `<source type="video/mp4">` の `src` を抽出
3. `video[poster]` からサムネイル画像を抽出（動画メディアとして記録）
4. `?tag=12` などのクエリパラメータは URL から除去せずそのまま保持（挙動不明なため）

**抽出ルール**:
- 動画URL: `video source[type="video/mp4"]` の `src`
- サムネイルURL: `video[poster]` の `poster` 属性値
- 重複排除: URLベース

### 3.3 メディアサイズ
- HTML DOM にサイズ情報がないため `width`/`height` は常に `0`
- 制限事項として明記

---

## 4. 実装タスク

### Task 1: `sotwe_parser.py` アバター修正
- `_normal.` を削除したURLを `avatar_url` と `avatar_original_url` に設定
- `_400x400.` 変換を削除

### Task 2: `sotwe_parser.py` 動画メディア抽出の修正
- `video.video-player` を親セレクタとして `source[type="video/mp4"]` の `src` を抽出
- `video[poster]` からサムネイル画像を抽出
- サムネイルは `type: "image"` として media_list に追加

### Task 3: ドキュメント修正 (`docs/COMBINED_WIKI.md`, `SPEC_*.md`, `PLAN_*.md`)
- Sotwe JSON API (`api.sotwe.com/api/v2`) の記述を削除
- Web UI DOMスクレイピング方式に修正
- SeleniumBase UC Mode 使用を明記

### Task 4: 検証
- `sotwe.html` を使った単体テスト作成・実行
- 修正後のアバターURL、メディアURL、メディア件数を確認

---

## 5. 制限事項（Web UI スクレイピングの限界）

- `width`/`height`: HTML DOM にサイズ情報がないため `0` のまま
- `conversation_id`, `reply_to_id`: HTML からは取得困難（未実装）
- 動画メディア: `?tag=12` などのクエリパラメータの挙動は未確認（そのまま保持）

---

*本計画はユーザーの指示（「現実装を正とする」）および追加情報（`_400x400.` は404、`_normal.` 削除でoriginal取得可能、`video-player` 要素の実例）に基づき更新された。*
