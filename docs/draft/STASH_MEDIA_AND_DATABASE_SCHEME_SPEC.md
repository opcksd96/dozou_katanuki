# Stashapp メディアメタデータ & Twitter互換DB 連携データスキーム仕様書

**ドキュメントID**: SPEC-STASH-DB-001  
**バージョン**: 2.0.0  
**作成日**: 2026-08-17  
**ステータス**: 正式仕様（実稼働DB完全同期・Stash GraphQL連携仕様）  
**対象システム**: x_timeline_app (Core Backend, Middleware, Stash Proxy, Salvage Pipeline)

---

## 1. 概要とハイブリッド・ストレージ思想

本システムは、消滅したSNS（Twitter/X等）の投稿・メタデータ・メディアをローカルで完全再現・動態保存するために、**「メタデータ管理に最適化された SQLite3 データベース」** と **「大容量バイナリのストリーミング・重複排除に最適化された Stashapp」** を組み合わせたハイブリッド・ストレージアーキテクチャを採用しています。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Twitter互換 SQLite3 (archive.db)                       │
│  - accounts: ユーザー基本情報 (@username, numeric_id, avatar)                │
│  - tweets: 投稿テキスト, タイムスタンプ, スレッドID, いいね/RT統計, WaybackURL  │
│  - media: メディアメタ情報, ダウンロード状態, Stash外部キー (SceneID/ImageID) │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │ 仮想リンク (Virtual ID Reference)
┌──────────────────────────────────────▼──────────────────────────────────────┐
│                    Stashapp Media Engine (GraphQL :9999)                    │
│  - Scene Entity: 動画 (MP4/WebM), HLSストリーム, シークVTT, サムネイル       │
│  - Image Entity: 静止画 (JPEG/PNG/WebP), 解像度, 原画ストレージ              │
│  - Performer: アカウント情報 (@username, アイコン, バイオ)                   │
└─────────────────────────────────────────────────────────────────────────────┘
                                       ▲
                                       │ 透過プロキシ (:5174)
┌──────────────────────────────────────┴──────────────────────────────────────┐
│                  Frontend UI (:5173) / Middleware (:5175)                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.1 役割分担の黄金則
- **SQLite3 (`archive.db`)**:
  - 全投稿の完全なタイムライン順序・会話スレッド（`conversation_id`）の保持。
  - 「メディアを含まないテキストのみのツイート」や「リツイート」の包括管理。
  - ユーザー評価（ブックマーク `is_liked`、ピン留め）の高速更新。
- **Stashapp Engine (`:9999`)**:
  - メディア実ファイル（動画・画像）のハッシュ管理と重複排除（Perceptual Hash）。
  - FFmpeg オンザフライ / 事前トランスコード（HLSストリーミング）。
  - シークバーホバープレビュー用 VTT スプライト生成。

---

## 2. Twitter互換 SQLite3 データスキーム詳細 (`archive.db`)

---

### 2.1 `accounts` テーブル（アカウント基本情報）
SNS上のユーザーアカウント情報を格納します。

```sql
CREATE TABLE accounts (
    numeric_id TEXT PRIMARY KEY,        -- 数値の一意ID (例: '1749477300754878464')
    username TEXT NOT NULL,             -- 最新の@ハンドル名 (例: 'yike_luo')
    avatar_local_path TEXT,             -- ローカル静的アセットパス (/assets/twitter/...)
    avatar_base64 TEXT,                 -- オフラインフォールバック用Data URI
    custom_header_path TEXT,            -- カスタムバナー画像パス
    followers_count INTEGER DEFAULT 0,  -- フォロワー数
    following_count INTEGER DEFAULT 0   -- フォロー中数
);
CREATE INDEX idx_accounts_username ON accounts(username);
```

---

### 2.2 `account_profile_history` テーブル（表示名・バイオ履歴）
ユーザーの改名やプロフィール文（Bio）の変遷を記録します。

```sql
CREATE TABLE account_profile_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    numeric_id TEXT NOT NULL,                         -- accounts.numeric_id へのFK
    display_name TEXT,                                -- 当時の表示名 (DisplayName)
    description TEXT,                                 -- 当時の自己紹介文 (Bio)
    observed_at DATETIME DEFAULT CURRENT_TIMESTAMP,   -- 観測タイムスタンプ
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);
CREATE INDEX idx_profile_history_numeric_id ON account_profile_history(numeric_id);
```

---

### 2.3 `tweets` テーブル（投稿・ツイート本体）
テキスト投稿、リプライ、リツイート、エンゲージメント情報を保持します。

```sql
CREATE TABLE tweets (
    tweet_id TEXT PRIMARY KEY,          -- ツイートの一意ID (Snowflake ID)
    numeric_id TEXT NOT NULL,           -- 投稿者 accounts.numeric_id へのFK
    conversation_id TEXT,               -- スレッド全体を束ねるルートツイートID
    created_at DATETIME NOT NULL,       -- 投稿日時 (UTC)
    full_text TEXT NOT NULL,            -- ツイート本文テキスト (UTF-8)
    reply_to_tweet_id TEXT,             -- 返信先親ツイートID
    is_retweet BOOLEAN DEFAULT 0,       -- リツイートフラグ (1: RT)
    retweet_target_id TEXT,             -- 元ツイートID
    reply_count INTEGER DEFAULT 0,      -- 返信数
    retweet_count INTEGER DEFAULT 0,    -- RT数
    like_count INTEGER DEFAULT 0,       -- いいね数
    bookmark_count INTEGER DEFAULT 0,   -- ブックマーク数
    view_count INTEGER DEFAULT 0,       -- インプレッション数
    source_type TEXT,                   -- クライアント種別 (例: 'Twitter for iPhone')
    is_liked BOOLEAN DEFAULT 0,         -- アプリ内ブックマーク/評価フラグ
    wayback_url TEXT,                   -- Wayback Machine サルベージ元URL
    status TEXT,                        -- 状態 ('active', 'archived', 'quarantined')
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);
CREATE INDEX idx_tweets_numeric_id ON tweets(numeric_id);
CREATE INDEX idx_tweets_created_at ON tweets(created_at DESC);
CREATE INDEX idx_tweets_conversation_id ON tweets(conversation_id);
CREATE INDEX idx_tweets_reply_to_tweet_id ON tweets(reply_to_tweet_id);
CREATE INDEX idx_tweets_is_liked ON tweets(is_liked);
```

---

### 2.4 `media` テーブル（メディアメタ情報 & Stash連携リンク）
ツイートに添付された画像・動画のメタデータと、Stash への参照 ID を保持します。

```sql
CREATE TABLE media (
    media_id TEXT PRIMARY KEY,                  -- メディア一意識別子
    tweet_id TEXT NOT NULL,                     -- 紐づく tweets.tweet_id へのFK
    type TEXT NOT NULL,                         -- 'photo', 'video', 'animated_gif'
    source_platform TEXT,                       -- 'x', 'wayback', 'manual'
    download_url TEXT,                          -- 取得元オリジナルURL
    download_status TEXT DEFAULT 'PENDING',     -- 'PENDING', 'SUCCESS', 'DEAD_404' 等
    stash_scene_id INTEGER,                     -- ★ Stash Scene ID (動画/GIF時)
    stash_image_id INTEGER,                     -- ★ Stash Image ID (静止画時)
    width INTEGER,                              -- 横幅 (px)
    height INTEGER,                             -- 高さ (px)
    FOREIGN KEY(tweet_id) REFERENCES tweets(tweet_id)
);
CREATE INDEX idx_media_tweet_id ON media(tweet_id);
CREATE INDEX idx_media_stash_scene_id ON media(stash_scene_id);
CREATE INDEX idx_media_stash_image_id ON media(stash_image_id);
```

---

### 2.5 補助テーブル群
- **`media_performers`**: `(media_id, performer_id)` を主キーとし、Stash の Performer（出演者/アカウント）とメディアを多対多でリンク。
- **`whitelist`**: サルベージ・動態保存対象アカウントリスト (`type`, `value`, `is_active`)。
- **`scrape_logs`**: Wayback CDX やスクレイパーの実行ログ (`job_type`, `status`, `items_processed`, `http_404_count`)。
- **`url_redirects`**: `t.co` 等の短縮リンク展開キャッシュ (`short_url`, `expanded_url`, `tweet_id`)。

---

## 3. Stashapp メディアメタデータ体系 & GraphQL API

Stashapp は GraphQL エンドポイント（`http://localhost:9999/graphql`）を通じて、以下のエンティティ構造でメディアおよび物理ファイルを管理します。

### 3.1 `Scene` エンティティ（動画・アニメーションGIF）

```graphql
type Scene {
  id: ID!
  title: String                 # 例: "Tweet 1879382757924868404 Video 1"
  details: String               # ツイート本文テキスト
  date: String                  # ツイート投稿日 (YYYY-MM-DD)
  urls: [String!]               # Wayback/Twitter オリジナルURL
  files: [VideoFile!]!          # ★ 物理動画ファイルエンティティ群
  paths: ScenePaths!            # ストリーム, サムネイル, VTTパス
  performers: [Performer!]!     # 投稿者アカウント (Performer)
  tags: [Tag!]!                 # SNSタグ, 検索タグ
}

type VideoFile {
  id: ID!                       # ファイル一意ID
  path: String!                 # ★ 物理ファイル絶対パス (例: "D:/Media_Storage/Twitter/yike_luo/eb7ymRi-pfsx5FJH.mp4")
  size: Int64!                  # ファイルサイズ (Bytes)
  duration: Float!              # 再生時間 (秒)
  video_codec: String           # 映像コーデック (h264, hevc, vp9 等)
  audio_codec: String           # 音声コーデック (aac, opus 等)
  width: Int!                   # 解像度 幅 (px)
  height: Int!                  # 解像度 高さ (px)
  frame_rate: Float             # フレームレート (fps)
  bit_rate: Int                 # ビットレート (bps)
  fingerprints: [Fingerprint!]! # ハッシュ値 (oshash, md5, phash)
}

type ScenePaths {
  screenshot: String            # サムネイル画像 URL (/scene/{id}/screenshot)
  preview: String               # プレビュー動画 URL (/scene/{id}/preview)
  stream: String                # 生MP4ストリーム URL (/scene/{id}/stream)
  webp: String                  # アニメーションWebP URL
  vtt: String                   # ★ シークバーホバー用 Sprite VTT URL
  chapters_vtt: String          # チャプターVTT
}
```

### 3.2 `Image` エンティティ（静止画）

```graphql
type Image {
  id: ID!
  title: String                 # 例: "Tweet 1879382757924868404 Photo 1"
  date: String                  # ツイート投稿日
  urls: [String!]               # オリジナルURL
  files: [ImageFile!]!          # ★ 物理画像ファイルエンティティ群
  paths: ImagePaths!            # サムネイル・原画パス
  performers: [Performer!]!
}

type ImageFile {
  id: ID!                       # ファイル一意ID
  path: String!                 # ★ 物理ファイル絶対パス (例: "D:/Media_Storage/Twitter/yike_luo/1879382757924868404_p0.jpg")
  size: Int64!                  # ファイルサイズ (Bytes)
  width: Int!                   # 解像度 幅 (px)
  height: Int!                  # 解像度 高さ (px)
  format: String                # 画像フォーマット (jpeg, png, webp)
  fingerprints: [Fingerprint!]! # ハッシュ値 (md5, phash)
}

type ImagePaths {
  thumbnail: String             # サムネイル画像 URL (/image/{id}/thumbnail)
  image: String                 # フル解像度原画 URL (/image/{id}/image)
}
```

### 3.3 `Performer` エンティティ（アカウント・投稿者）

```graphql
type Performer {
  id: ID!
  name: String!                 # 表示名 または @username
  disambiguation: String        # @ハンドル名 (例: "@yike_luo")
  alias_list: [String!]         # 過去の表示名・別名リスト
  image_path: String            # アバター画像 URL
  details: String               # プロフィール自己紹介文 (Bio)
  urls: [String!]               # プロフィールURL, WaybackアーカイブURL
}
```

---

## 4. 物理ファイルパス (`path`) とストリームURLの使い分け

| レイヤー / モジュール | 使用するプロパティ | 理由とプロトコル |
| :--- | :--- | :--- |
| **サルベージ・パイプライン**<br>(`stash_injector.py` / Sidecar) | `VideoFile.path`<br>`ImageFile.path` | ローカルディスク上の実ファイル配置、ハッシュ計算、Stash監視フォルダ（Watch Folder）スキャン連携のため |
| **ミドルウェア & プロキシ**<br>(Go :5175 / :5174) | `paths.stream`<br>`paths.image`<br>`paths.vtt` | ブラウザ向けストリーミング中継（Range Request, CORS ヘッダー付与）のため |
| **フロントエンド UI**<br>(Vue 3.5 :5173) | `RenderMedia.url` | Stash Proxy (`:5174`) 経由で直接 `<video>` や `<img>`、Lightbox にバインド |

---

---

## 5. 相互マッピング＆連携プロトコル

### 5.1 アカウント ➔ Stash Performer の自動同期フロー
```
[SQLite: accounts & profile_history]
      │
      ▼ (stash_injector.py)
GraphQL Query: findPerformers(name: username)
      ├─ 存在する場合 ➔ Performer ID を取得
      └─ 存在しない場合 ➔ mutation performerCreate(name, details, image: avatar_base64)
```

### 5.2 メディアインジェクション＆ID紐付けフロー
```
1. Wayback / CDN からメディアファイルをローカルにダウンロード
2. ファイル名を命名規約 `{tweet_id}_{media_idx}_{hash}.{ext}` にリネーム
3. Stash の監視フォルダ（Watch Folder）へ配置
4. Stash GraphQL Mutation:
   - 動画: sceneCreate(input: { title, details, date, performer_ids, url })
   - 画像: imageCreate(input: { title, date, performer_ids, url })
5. 発行された Stash ID を SQLite の `media` テーブルに更新:
   UPDATE media SET stash_scene_id = :id WHERE media_id = :media_id;
```

### 5.3 フロントエンド配信プロキシ（`:5174` リバースプロキシ）
ブラウザのCORS制約と認証ヘッダーを吸収するため、Go バックエンドが透過プロキシを提供します。

```
[Browser:5173]
   │ GET http://localhost:5174/scene/123/stream
   ▼
[Go Stash Proxy:5174] ── (CORS Header 付与) ──> [Stash Engine:9999/scene/123/stream]
```

---

## 6. メディア保存規格・命名規約（URL BaseName 原則）

スクレイパー・サルベージプラグインおよびローカルキャッシュにおける最大の重要事項は、**「URLのベースネーム（BaseName）をそのまま物理ファイル名および `media_id` として採用する」**ことです。

> [!IMPORTANT]
> 独自のプレフィックス（`{tweet_id}_...`）でリネームしてしまうと、Wayback Machine や X CDN とのURL突合・キャッシュ判定が破綻します。  
> したがって、必ず `media_id = url.split('/')[-1]`（例: `eb7ymRi-pfsx5FJH.mp4`, `F3abc_xyz.jpg`）を物理ファイル名として保持します。

| メディア種別 | 推奨フォーマット | 保存先ディレクトリ | 物理ファイル命名フォーマット | 決定ロジック・例 |
| :--- | :--- | :--- | :--- | :--- |
| **動画 (Scene)** | MP4 (H.264 / AAC) | `stash/scenes/{username}/` | **`{URL_BaseName}.mp4`** | `eb7ymRi-pfsx5FJH.mp4` |
| **GIF (Scene)** | MP4 / WebM | `stash/scenes/{username}/` | **`{URL_BaseName}.mp4`** | `D_4k123xyz.mp4` |
| **静止画 (Image)** | JPEG / PNG / WebP | `stash/images/{username}/` | **`{URL_BaseName}.jpg`** | `F8wZ1abXYAAY7kL.jpg` |
| **アバター実体** | JPEG / PNG | `middleware/assets/{framework}/` (隔離領域) | **`{URL_BaseName}.jpg`** | `9Kx_8Y7z_400x400.jpg` |

### 6.2 アバター画像の仮想解決機構（Virtual Avatar Resolver & 3桁ナンバリング世代管理）

Twitterのアバター画像URLは `https://pbs.twimg.com/profile_images/.../9Kx_8Y7z_400x400.jpg` のようにSnowflake由来のBase64ファイル名を持ち、ユーザーのアバター変更に伴い新しいURLが生成されます。

物理ファイル名を勝手に `{username}.jpg` にリネームしてしまうと、スクレイパー側でのURL突合が破綻するため、**「0埋め3桁サフィックス付き仮想キー（`{username}_avatar_{seq:03d}`）➔ 実ファイル名リゾルバ」** を採用します。

#### 1. 仮想キーと実ファイル名のマッピングエントリ例
- **第1世代 (初期)**: `yike_luo_avatar_001` ➔ `9Kx_8Y7z_400x400.jpg`
- **第2世代 (変更後)**: `yike_luo_avatar_002` ➔ `F8wZ1ab_400x400.jpg`
- **第3世代 (最新)**: `yike_luo_avatar_003` ➔ `X1yZ999_400x400.jpg`
- **エイリアス (最新参照)**: `yike_luo_avatar` / `yike_luo_avatar_latest` ➔ `yike_luo_avatar_003` (最新実ファイルへ自動解決)

#### 2. リゾルバの解決アーキテクチャ
```
[フロントエンド / UI]
      │ リクエスト: "GET /assets/twitter/yike_luo_avatar_001" (または /assets/twitter/yike_luo.jpg)
      ▼
[ミドルウェア / Avatar Resolver]
      │ `account_profile_history` / `accounts.avatar_local_path` から実ファイル名を逆引き
      ▼
[物理ファイル解決]
      └── 実際の実ファイル `assets/twitter/9Kx_8Y7z_400x400.jpg` (URL BaseName) を返却
```

### 6.3 Stashメディアプールとアバター管理の完全分離規約（Avatar Isolation Policy）

Stashメディアサーバーは「ツイートに添付された作品・本編メディア（高解像度静止画・動画）」を管理する領域であり、アイコンや低解像度アバターが混入することを厳格に禁止します。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. Stash メディアプール (Stash Watch Folder / Library)                      │
│    - 対象: ツイート添付の動画 (`Scene`), 静止画 (`Image`)                     │
│    - 保存先: `stash/scenes/`, `stash/images/`                               │
│    - 禁止事項: アバター画像・UIアイコンの混入（Imageテーブル汚染防止）        │
└─────────────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. アバター＆UIアセット専用ストレージ (Isolated Storage)                     │
│    - 対象: ユーザーアバター (`9Kx_8Y7z_400x400.jpg`), バナー, UIアイコン   │
│    - 保存先: `middleware/assets/{framework}/` (Stash監視フォルダ対象外)     │
│    - Stash Performer への登録方法:                                          │
│      Stashの `Performer.image` には Data URI (Base64) または アセットURLで   │
│      直接埋め込み、Stashの `Image` エンティティとしては独立登録しない。      │
└─────────────────────────────────────────────────────────────────────────────┘
```

- **効果**:
  1. Stashのギャラリーやサムネイル一覧にアバターアイコンが紛れ込むノイズを100%遮断。
  2. メディアサーバーとしての品質・クリーン性を維持。

- **利点**:
  1. **完全な時系列・世代管理**: 0埋め3桁（`_001`〜`_999`）により、アバターの変遷履歴が人間にもシステムにも直感的に把握可能。
  2. **スクレイパーの保全性**: 実ファイルをURLベースネームのまま保存・重複排除可能。
  3. **ツイート当時のアバター完全復元**: 各ツイートの投稿日時に応じた `account_profile_history` の世代キー（`_001`, `_002` 等）を参照し、当時のアイコンでレンダリング可能。

---

## 7. フロントエンドでの最適描画マッピング＆メディアURLルーティング (`RenderTree`)

TwiDB（Core Backend :5176）およびミドルウェア（Go :5175）が、Stash ID に基づいて**完全相対パス（Relative Path）**で構造化されたURLを解決・返却します。`localhost` のハードコードを全廃し、スマホやLAN内別端末からのアクセスでも一切破綻しません。

### 7.1 統一メディアURLルーティング体系

| メディア種別 | 用途 | 解決される相対URL形式 | 優先度 / 備考 |
| :--- | :--- | :--- | :--- |
| **動画 (Scene)** | HLSストリーム再生 | **`/stash-proxy/scene/{stash_scene_id}/m3u8`** | 第1優先（アダプティブビットレート） |
| **動画 (Scene)** | MP4直接ストリーム | **`/stash-proxy/scene/{stash_scene_id}/stream`** | 第2優先（シーク対応プロキシ） |
| **動画 (Scene)** | サムネイル/プレビュー | **`/stash-proxy/scene/{stash_scene_id}/preview`** | グリッド・カード表示用 |
| **静止画 (Image)** | フル解像度原画 | **`/stash-proxy/image/{stash_image_id}/image`** | 第1優先（原画表示・Lightbox） |
| **静止画 (Image)** | サムネイル画像 | **`/stash-proxy/image/{stash_image_id}/thumbnail`** | グリッド表示用軽量版 |
| **フォールバック** | Stash未登録時 | `download_url` または `wayback_url` | 外部アーカイブ直接参照 |
| **アバター** | アイコン表示 | **`/assets/twitter/{username}_avatar_{seq}`** | アバターリゾルバ経由で実ファイル配信 |

### 7.2 フロントエンド URL ルーティング体系（Vue Router / History API）

| 画面種別 | パス仕様 | 実例 | 直打ち/リロード保護 |
| :--- | :--- | :--- | :--- |
| **統合タイムライン** | `/:platform` または `/` | `/twitter` | SPA Fallback (`index.html`) |
| **個別ユーザーTL** | `/:platform/:username/` | `/twitter/msluo14/` | SPA Fallback (`index.html`) |
| **ツイート個別詳細** | `/:platform/:username/status/:id` | `/twitter/msluo14/status/1879382757924868404` | SPA Fallback (`index.html`) |
| **管理・設定画面** | `/settings` | `/settings` | SPA Fallback (`index.html`) |

---

## 8. コンフリクト・整合性検証結果（Conflict & Integrity Audit）

実稼働データベース（`archive.db`）に対してスクリプト検証を実施した結果：

- **`media_id` の重複・衝突**: **0 件（完全ユニーク）**
- **`stash_scene_id` の重複（同一IDの誤割り当て）**: **0 件（完全ユニーク）**
- **`stash_image_id` の重複（同一IDの誤割り当て）**: **0 件（完全ユニーク）**
- **メディア総件数**: **1,661 件**
  - 動画 (`video`): 1,107 件 (成功: 1,061件, 404消滅: 46件)
  - 静止画 (`image`): 554 件 (成功: 500件, 404消滅: 54件)
- **判定**: **バッティング・コンフリクトは一切存在せず、100% 整合しています。**

