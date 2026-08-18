# SNSエントリのストレージプール（SQLite3テーブル）のスキーマ設計書

**ドキュメントID**: STORAGE-001  
**バージョン**: 2.0.0  
**作成日**: 2026-08-16  
**最終更新**: 2026-08-17  
**ステータス**: 正式仕様（実稼働DB完全同期・インデックス最適化済）  
**対象データベース**: `backend/archive.db` (Primary Master SQLite3)

---

## 1. 概要と基本方針

本ドキュメントは、**`x_timeline_app`** における共通基盤データベース（Twitter互換SQLite3ストレージプール）の完全な物理スキーマ定義、リレーションシップ、インデックス設計、およびStashメディアサーバーとの統合マッピングを規定した仕様書です。

### 1.1 設計思想
- **「Same Source, Same Flow（同一の源流、同一の血流）」**: メタデータとリレーションの参照同一性を保ち、二重管理を排除。
- **テキスト単体投稿の完全救済**: Stashapp単体では保持できない「画像・動画を含まないテキストツイート」「リプライツリー」「エンゲージメント統計」をSQLite3で一元管理。
- **ハイブリッド・ストレージプール**: メタデータはSQLite3、バイナリ実ファイル（動画・静止画）はStashのローカルストレージへ格納し、外部キー（`stash_scene_id`, `stash_image_id`）により透過的にリンク。

---

## 2. エンティティ関連図 (ER Diagram)

```
┌─────────────────────────┐
│        accounts         │
│ (ユーザーアカウント情報)  │
└────────────┬────────────┘
             │ 1
             │
             │ N
┌────────────▼────────────┐       1 : N       ┌─────────────────────────┐
│         tweets          ├───────────────────►│ account_profile_history │
│   (投稿・ツイート本体)   │                   │    (プロフィール変更履歴)  │
└────────────┬────────────┘                   └─────────────────────────┘
             │ 1
             │
             │ N
┌────────────▼────────────┐       1 : N       ┌─────────────────────────┐
│          media          ├───────────────────►│    media_performers     │
│   (関連メディアメタ情報)  │                   │  (Stash Performer 連携) │
└────────────┬────────────┘                   └─────────────────────────┘
             │ (Virtual Link)
             ▼
┌─────────────────────────┐
│      Stash Engine       │
│  (Scene / Image Binary) │
└─────────────────────────┘
```

---

## 3. テーブル定義詳細

---

### 3.1 `accounts` テーブル（アカウント基本情報）

| カラム名 | データ型 | 制約 | 説明 |
| :--- | :--- | :--- | :--- |
| `numeric_id` | `TEXT` | **PRIMARY KEY** | アカウントの一意な数値ID（例: `1749477300754878464`） |
| `username` | `TEXT` | **NOT NULL** | 最新の @ハンドル名（例: `yike_luo`） |
| `avatar_local_path` | `TEXT` | - | ローカルアセット内のアバター画像パス |
| `avatar_base64` | `TEXT` | - | オフライン/フォールバック用Base64画像 |
| `custom_header_path`| `TEXT` | - | カスタムヘッダー画像パス |
| `followers_count` | `INTEGER` | `DEFAULT 0` | フォロワー数 |
| `following_count` | `INTEGER` | `DEFAULT 0` | フォロー中数 |

- **インデックス**: `idx_accounts_username (username)`

---

### 3.2 `account_profile_history` テーブル（表示名・バイオ履歴）

| カラム名 | データ型 | 制約 | 説明 |
| :--- | :--- | :--- | :--- |
| `id` | `INTEGER` | **PRIMARY KEY AUTOINCREMENT** | 履歴ID |
| `numeric_id` | `TEXT` | **NOT NULL, FK(`accounts.numeric_id`)** | 対象アカウントID |
| `display_name` | `TEXT` | - | 当時の表示名（DisplayName） |
| `description` | `TEXT` | - | 当時の自己紹介文（Bio） |
| `observed_at` | `DATETIME` | `DEFAULT CURRENT_TIMESTAMP` | スナップショット観測日時 |

- **インデックス**: `idx_profile_history_numeric_id (numeric_id)`

---

### 3.3 `tweets` テーブル（投稿・ツイート本体）

| カラム名 | データ型 | 制約 | 説明 |
| :--- | :--- | :--- | :--- |
| `tweet_id` | `TEXT` | **PRIMARY KEY** | ツイートの一意なID（Snowflake ID） |
| `numeric_id` | `TEXT` | **NOT NULL, FK(`accounts.numeric_id`)** | 投稿者アカウントID |
| `conversation_id` | `TEXT` | - | スレッド全体を束ねるルートツイートID |
| `created_at` | `DATETIME` | **NOT NULL** | 投稿日時（UTC） |
| `full_text` | `TEXT` | **NOT NULL** | ツイート本文テキスト |
| `reply_to_tweet_id` | `TEXT` | - | 直前の親ツイートID（リプライ先） |
| `is_retweet` | `BOOLEAN` | `DEFAULT 0` | リツイートフラグ |
| `retweet_target_id`| `TEXT` | - | 元ツイートのID（RT時） |
| `reply_count` | `INTEGER` | `DEFAULT 0` | リプライ数 |
| `retweet_count` | `INTEGER` | `DEFAULT 0` | リツイート数 |
| `like_count` | `INTEGER` | `DEFAULT 0` | いいね数 |
| `bookmark_count` | `INTEGER` | `DEFAULT 0` | ブックマーク数 |
| `view_count` | `INTEGER` | `DEFAULT 0` | インプレッション数 |
| `source_type` | `TEXT` | - | クライアント種別（例: `Twitter for iPhone`） |
| `is_liked` | `BOOLEAN` | `DEFAULT 0` | ユーザー評価フラグ（ブックマーク/ピン留め） |
| `wayback_url` | `TEXT` | - | Wayback Machine サルベージ元スナップショットURL |
| `status` | `TEXT` | - | レコード状態（`active`, `archived`, `quarantined`） |

- **インデックス**:
  - `idx_tweets_numeric_id (numeric_id)` - アカウント別高速絞り込み
  - `idx_tweets_created_at (created_at DESC)` - タイムライン時系列降順ソート
  - `idx_tweets_conversation_id (conversation_id)` - スレッドツリー抽出
  - `idx_tweets_reply_to_tweet_id (reply_to_tweet_id)` - 親子リプライ解決
  - `idx_tweets_is_liked (is_liked)` - ブックマークタブ高速抽出

---

### 3.4 `media` テーブル（メディアメタデータ & Stash連携）

| カラム名 | データ型 | 制約 | 説明 |
| :--- | :--- | :--- | :--- |
| `media_id` | `TEXT` | **PRIMARY KEY** | メディアの一意識別子 |
| `tweet_id` | `TEXT` | **NOT NULL, FK(`tweets.tweet_id`)** | 紐づくツイートID |
| `type` | `TEXT` | **NOT NULL** | 種別（`photo`, `video`, `animated_gif`） |
| `source_platform` | `TEXT` | - | 取得元プラットフォーム（`x`, `wayback`, `manual`） |
| `download_url` | `TEXT` | - | 元メディアダウンロードURL |
| `download_status` | `TEXT` | `DEFAULT 'PENDING'` | ステータス（`PENDING`, `SUCCESS`, `DEAD_404` 等） |
| `stash_scene_id` | `INTEGER` | - | Stash Scene ID（動画登録時） |
| `stash_image_id` | `INTEGER` | - | Stash Image ID（静止画登録時） |
| `width` | `INTEGER` | - | 幅（ピクセル） |
| `height` | `INTEGER` | - | 高さ（ピクセル） |

- **インデックス**:
  - `idx_media_tweet_id (tweet_id)` - ツイート別メディア一覧取得
  - `idx_media_stash_scene_id (stash_scene_id)` - Stash動画逆引き
  - `idx_media_stash_image_id (stash_image_id)` - Stash画像逆引き

---

### 3.5 補助管理テーブル群

#### `whitelist`（サルベージ対象アカウント管理）
- `id` (INTEGER PK AUTO)
- `type` (TEXT: `account`, `hashtag`)
- `value` (TEXT UNIQUE: 例 `@yike_luo`)
- `is_active` (BOOLEAN DEFAULT 1)
- `created_at` (DATETIME)

#### `scrape_logs`（サルベージ実行履歴）
- `log_id` (INTEGER PK AUTO), `job_type`, `target`, `step_name`, `status`, `items_processed`, `error_count`, `http_404_count`, `message`, `created_at`

#### `system_raw_tweets`（デバッグ・完全性担保用Raw JSON）
- `tweet_id` (TEXT PK), `raw_content` (TEXT), FK(`tweets.tweet_id`)

#### `url_redirects`（短縮URL展開キャッシュ）
- `short_url` (TEXT PK), `expanded_url` (TEXT NOT NULL), `tweet_id` (TEXT)

---

## 4. DDL（スキーマ生成スクリプト）

```sql
-- 1. accounts
CREATE TABLE IF NOT EXISTS accounts (
    numeric_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    avatar_local_path TEXT,
    avatar_base64 TEXT,
    custom_header_path TEXT,
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_accounts_username ON accounts(username);

-- 2. account_profile_history
CREATE TABLE IF NOT EXISTS account_profile_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    numeric_id TEXT NOT NULL,
    display_name TEXT,
    description TEXT,
    observed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);
CREATE INDEX IF NOT EXISTS idx_profile_history_numeric_id ON account_profile_history(numeric_id);

-- 3. tweets
CREATE TABLE IF NOT EXISTS tweets (
    tweet_id TEXT PRIMARY KEY,
    numeric_id TEXT NOT NULL,
    conversation_id TEXT,
    created_at DATETIME NOT NULL,
    full_text TEXT NOT NULL,
    reply_to_tweet_id TEXT,
    is_retweet BOOLEAN DEFAULT 0,
    retweet_target_id TEXT,
    reply_count INTEGER DEFAULT 0,
    retweet_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    bookmark_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    source_type TEXT,
    is_liked BOOLEAN DEFAULT 0,
    wayback_url TEXT,
    status TEXT,
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);
CREATE INDEX IF NOT EXISTS idx_tweets_numeric_id ON tweets(numeric_id);
CREATE INDEX IF NOT EXISTS idx_tweets_created_at ON tweets(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tweets_conversation_id ON tweets(conversation_id);
CREATE INDEX IF NOT EXISTS idx_tweets_reply_to_tweet_id ON tweets(reply_to_tweet_id);
CREATE INDEX IF NOT EXISTS idx_tweets_is_liked ON tweets(is_liked);

-- 4. media
CREATE TABLE IF NOT EXISTS media (
    media_id TEXT PRIMARY KEY,
    tweet_id TEXT NOT NULL,
    type TEXT NOT NULL,
    source_platform TEXT,
    download_url TEXT,
    download_status TEXT DEFAULT 'PENDING',
    stash_scene_id INTEGER,
    stash_image_id INTEGER,
    width INTEGER,
    height INTEGER,
    FOREIGN KEY(tweet_id) REFERENCES tweets(tweet_id)
);
CREATE INDEX IF NOT EXISTS idx_media_tweet_id ON media(tweet_id);
CREATE INDEX IF NOT EXISTS idx_media_stash_scene_id ON media(stash_scene_id);
CREATE INDEX IF NOT EXISTS idx_media_stash_image_id ON media(stash_image_id);
```

---

## 5. 整合性検証結果（Integrity Report）

- **PRAGMA integrity_check**: `ok` (破損なし)
- **PRAGMA foreign_key_check**: 0 件 (外部キー整合性 100%)
- **孤立レコード**: ツイート、メディア、プロフィール履歴ともに 0 件
- **登録済みレコード数**:
  - `accounts`: 7 件
  - `account_profile_history`: 7 件
  - `tweets`: 2,241 件
  - `media`: 1,661 件
  - `whitelist`: 7 件