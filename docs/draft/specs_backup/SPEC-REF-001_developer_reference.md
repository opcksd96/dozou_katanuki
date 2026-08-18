# x_timeline_app 総合ドキュメントリファレンス & 技術仕様サマリ (Master Reference)

**プロジェクト名**: x_timeline_app  
**ドキュメントID**: DOC-REF-001  
**バージョン**: 2.0.0  
**作成日**: 2026-08-17  
**目的**: AI Agent および開発者がコードベース全体を走査することなく、本リファレンスを起点として最小限のトークン消費で高精度な実装・改修を行うための統一インデックス＆技術サマリ。

---

## 📑 1. ドキュメント体系カタログ (Documentation Catalog)

プロジェクト内の設計書・レポートは以下の体系で整理されています。実装目的に応じて必要なドキュメントを参照してください。

```
docs/
├── reference.md (DOC-REF-001)                ★ 本ドキュメント（総合リファレンス＆インデックス）
├── STASH_MEDIA_AND_DATABASE_SCHEME_SPEC.md (SPEC-STASH-DB-001) ★ Stashメディア＆Twitter互換DB連携仕様書
├── architecture_general_overview.md (SPEC-WIKI-001) ★ 総合技術概論＆全11章Wikiマスター
├── application_layers_spec.md (SPEC-LAYER-001) アプリケーション5層構造の仕様＆実装照合書
├── project_summary.md                          プロジェクト概要・サービス構成基本情報
├── INFRA-002_基礎インフラストラクチャ仕様.md     外部サービス連携・インフラ仕様書
├── ARCH-001_基礎アーキテクチャ詳細レポート.md     MVVMモデル・レイヤー別責務・権限定義
├── STORAGE-001_1.1_SNSストレージプールスキーマ設計書.md データベーススキーマ・テーブル詳細定義
├── MEDIA-001_1.1_メディア保存ストラテジレポート.md メディア保存・Stashインジェクション・フォーマット設計
├── frontend_requirements_report.md             フロントエンド要件・コンポーネント構成仕様
├── implementation_plan.md                      フロントエンド再構築計画＆コーディング規約
└── 20260816_アーキテクチャ権限違反分析レポート.md  レイヤー権限違反の分析とリファクタリング方針
```

### ドキュメント別 参照ガイド

| ドキュメント | 主な内容 | 参照すべき作業・場面 |
| :--- | :--- | :--- |
| **`reference.md`** | 全体インデックス、API/DBサマリ、コーディング規約 | **全開発作業の開始時、実装方針確認** |
| **`STASH_MEDIA_AND_DATABASE_SCHEME_SPEC.md`** | Stash GraphQLメタデータ（Scene, Image, Performer）とSQLiteの相互連携仕様 | **メディア再生、Stashインジェクション、DBモデル開発** |
| **`architecture_general_overview.md`** | 全11章の技術仕様、外部サービス、DB、UI、ミドルウェア概論 | **システム全体の包括的理解・Wikiドキュメント** |
| **`application_layers_spec.md`** | 5層モデル（Presentation, Dispatcher, Driver, Plugin, Controller）の仕様と照合 | アーキテクチャの変更、層間通信設計 |
| **`STORAGE-001_1.1_...`** | SQLiteテーブル定義、フィールド型、リレーション、インデックス | DBモデル追加・クエリ変更・マイグレーション |
| **`MEDIA-001_1.1_...`** | Stash連携、ダウンロードステータス、メディアURL解決 | 動画/画像プレイヤー改修、Stashプロキシ連携 |
| **`INFRA-002_...`** | ポート設定、Wayback API、Stash API、外部ツール連携 | サーバー起動設定、外部API連携 |
| **`implementation_plan.md`** | 1ファイル100行以下制約、MVVM+関数型設計 | フロントエンド新規作成・コンポーネント分割 |
| **`20260816_...分析レポート.md`**| Composable分離、API呼び出し権限違反の是正ガイド | フロントエンドのリファクタリング |

---

## 🗺️ 2. システム配置マップ (Where is What)

| 領域 | 主要ファイル・パス | 役割・技術 |
| :--- | :--- | :--- |
| **Frontend (View)** | `frontend/src/components/tweet/` | ツイート表示コンポーネント (`TweetCard`, `TweetBody`, `MediaGrid` 等) |
| | `frontend/src/components/media/` | メディア表示 (`MediaOverlay`, `StashPlayer`) |
| | `frontend/src/components/layout/` | レイアウト (`AppSidebar`, `StatusBanner`, `HeroHeader`) |
| **Frontend (ViewModel)** | `frontend/src/composables/` | 状態管理・副作用 (`useTimeline.ts`, `useMediaOverlay.ts`) |
| **Frontend (Utils/Model)**| `frontend/src/utils/` | 純粋関数 (`formatters.ts`, `parser.ts`) |
| | `frontend/src/models/` | 型定義 (`RenderTree.ts`, `StashData.ts`) |
| **Middleware Hub** | `middleware/main.go` (:5175) | 集約配信・`/api/render`・翻訳プロキシ |
| | `middleware/plugins/` | ★ **プラグイン配置位置** (`twitter.go`, `plugin.go` 等) |
| | `middleware/assets/{framework}/` | ★ **プラグイン用アセット配置位置** (`assets/twitter/` 等) |
| | `middleware/settings/` | 設定永続化 (`settings.yaml`) |
| **Core Backend API** | `backend/main.go` (:5176) | エントリポイント・ルーティング |
| | `backend/api/` | ハンドラ (`post_api.go`, `account_api.go`, `translate_api.go`) |
| | `backend/crud/` | GORM クエリ (`post.go`, `account.go`) |
| | `backend/models/` | DBモデル定義 (`models.go`) |
| | `backend/db/` | SQLite3 接続初期化 (`db.go`) |
| **Stash Reverse Proxy** | `backend/main.go` (:5174) | Stash (:9999) 向け CORS プロキシ |
| **Downloader Sidecar** | `wayback_tweet_rescure/` | Wayback サルベージ、Stash インジェクタ、クリーンアップ |
| **非推奨・バックアップ** | `orz/` | ⚠️ **【将来削除予定】破損バージョン（参照・依存禁止）** |
| | `x_timeline_app_202608142050/` | 📦 **【バックアップ】2026/08/14時点のスナップショット** |

---

## ⚡ 3. 技術仕様クイックサマリ

### 3.1 ネットワーク & API エンドポイント

```
[Browser:5173] 
   ──> [Middleware:5175] 
          ├──> GET /api/render?account_id={id}&filter={all|retweets|media|bookmarks}&framework=twitter
          ├──> GET /api/accounts
          ├──> GET /api/account?id={id}
          ├──> POST /api/translate
          │
          └───> [Core Backend:5176] (Internal REST)
                   ├──> GET /api/posts?account_id={id}&filter={filter}&limit={n}&offset={n}
                   ├──> GET /api/accounts
                   ├──> GET /api/account?id={id}
                   └──> POST /api/translate
[Browser:5173] 
   ──> [Stash Proxy:5174] ──> [Stash:9999] (Media Streams / Static Assets)
```

### 3.2 データベーススキーマ要約 (`archive.db` / `chronoarchive.db`)

```sql
-- accounts
CREATE TABLE accounts (
    numeric_id VARCHAR PRIMARY KEY,
    username VARCHAR,
    avatar_local_path VARCHAR,
    avatar_base64 VARCHAR,
    custom_header_path VARCHAR,
    followers_count INTEGER,
    following_count INTEGER
);

-- account_profile_history
CREATE TABLE account_profile_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    numeric_id VARCHAR,
    display_name VARCHAR,
    description TEXT,
    observed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- tweets
CREATE TABLE tweets (
    tweet_id VARCHAR PRIMARY KEY,
    numeric_id VARCHAR,
    conversation_id VARCHAR,
    created_at DATETIME,
    full_text TEXT,
    reply_to_tweet_id VARCHAR,
    is_retweet BOOLEAN,
    retweet_target_id VARCHAR,
    reply_count INTEGER,
    retweet_count INTEGER,
    like_count INTEGER,
    bookmark_count INTEGER,
    view_count INTEGER,
    source_type VARCHAR,
    is_liked BOOLEAN,
    wayback_url VARCHAR,
    status VARCHAR,
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);

-- media
CREATE TABLE media (
    media_id VARCHAR PRIMARY KEY,
    tweet_id VARCHAR,
    type VARCHAR, -- 'image', 'video', 'animated_gif'
    source_platform VARCHAR,
    download_url VARCHAR,
    download_status VARCHAR DEFAULT 'PENDING',
    stash_scene_id INTEGER,
    stash_image_id INTEGER,
    width INTEGER,
    height INTEGER,
    FOREIGN KEY(tweet_id) REFERENCES tweets(tweet_id)
);
```

### 3.3 主要フロントエンド型定義 (`RenderTree`)

```typescript
export interface RenderTree {
  id: string;
  framework: string;
  author: RenderAuthor;
  content: string;
  created_at: string;
  media: RenderMedia[];
  meta: RenderMeta;
  stashData?: StashData;
}

export interface RenderAuthor {
  name: string;
  avatar_url: string;
  handle: string;
}

export interface RenderMedia {
  id: string;
  type: string;
  url: string; // 主ストリーム/原画URL (相対パス)
  urls?: {
    stream?: string;    // "/stash-proxy/scene/123/stream"
    m3u8?: string;      // "/stash-proxy/scene/123/m3u8"
    image?: string;     // "/stash-proxy/image/456/image"
    thumbnail?: string; // "/stash-proxy/image/456/thumbnail"
    original?: string;  // "https://..."
  };
  stash_id?: number;
  stash_type?: 'scene' | 'image';
}

export interface RenderMeta {
  style_class?: string;
  is_retweet?: boolean;
  wayback_url?: string;
  translated_text?: string;
}
```

### 3.4 フロントエンド & メディア URL ルーティング規約

- **SPA フロントエンド**:
  - `/:platform`（`/twitter` / `/`）: 統合タイムライン
  - `/:platform/:username/`（`/twitter/msluo14/`）: 個別タイムライン
  - `/:platform/:username/status/:id`: ツイート個別詳細
  - `/settings`: 管理画面
  - **SPA Fallback**: Goサーバーは404時に `index.html` を返却し、URL直打ち・F5リロードを保護。
- **メディアプロキシ**:
  - すべてホスト名を含まない**完全相対パス（`/stash-proxy/...`, `/assets/...`）**で配信し、LAN・スマホ（`0.0.0.0` listen）でゼロコンフィグ動作。

---

## 📏 4. 実装計画 & コーディング規約 (Coding Standards)

AI Agent はコードを生成・修正する際、以下の原則を厳格に遵守してください。

### 1. 「1ファイル 100行以下」の原則 (Strict Rule)
- コンポーネントおよびモジュールは極限まで単一責任に細分化する。
- 100行を超えそうな場合：
  - スタイルは Tailwind または共通 CSS に逃がす。
  - 計算・文字列操作は `frontend/src/utils/` の純粋関数に切り出す。
  - 状態・API 通信は `frontend/src/composables/` に切り出す。
  - サブコンポーネントに分割する（例: `TweetCard` → `TweetAuthor`, `TweetBody`, `TweetStats`, `MediaGrid`）。

### 2. レイヤー別 責務境界 (MVVM & 関数型原則)

```
┌────────────────────────────────────────────────────────────┐
│ Presentation (components/*.vue)                            │
│  - UI描画とイベント発行のみ (Stateless / Pure)             │
│  - 禁止: API直接呼出、グローバル状態の保持                │
└───────────────────────────┬────────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────────┐
│ ViewModel (composables/*.ts)                               │
│  - 状態管理 (ref/reactive)、API呼出、ライフサイクル        │
│  - 禁止: DOM直接操作、HTMLマークアップの混入              │
└───────────────────────────┬────────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────────┐
│ Utility & Domain (utils/*.ts, models/*.ts)                 │
│  - 純粋関数 (同一入力に対して常に同一出力、副作用ゼロ)    │
│  - 型定義 (TypeScript Interface / Type)                    │
└────────────────────────────────────────────────────────────┘
```

### 3. データフローとファイル配置の黄金律 (Same Source, Same Flow)
- **同階層3点セット（Single Source of Truth）**:
  - `x_timeline_app.exe`、`archive.db`、`config.json` は常に**プロジェクトルートの同一階層に1つだけ配置**し、開発時もリリース時も同じ実データを読み込む。
- **スクリプト隔離原則**:
  - テストスクリプト、検証用コード、仮スクリプトは**絶対にルートに放置せず、必ず `./scripts/` 配下に保存する**。
- **一括ビルド＆起動の徹底（個別手動起動の禁止）**:
  - AIエージェントが勝手に複数ターミナルで `npm run dev` や `go run` を個別起動し、ポートを専有してゾンビ化させる行為を厳禁とする。
  - ビルドは必ず `build.bat`、起動・デバッグは必ず `start.bat` を通じて行う。

---

## 🎯 5. AI Agent 向け作業別クイックナビゲーション

| 作業タスク | まず確認すべきドキュメント | 編集対象ファイル群 | 注意事項 |
| :--- | :--- | :--- | :--- |
| **タイムラインUIの改修・追加** | `docs/implementation_plan.md` | `frontend/src/components/tweet/`<br>`frontend/src/composables/` | 100行制限、UI内に直接API呼び出しを書かない |
| **メディア表示・プレイヤー改修** | `docs/MEDIA-001_1.1_...` | `frontend/src/components/media/`<br>`frontend/src/composables/useMediaOverlay.ts` | Stashプロキシ（`:5174`）経由のURL解決 |
| **新規APIエンドポイント追加** | `docs/application_layers_spec.md` | `backend/api/`, `backend/crud/`<br>`middleware/main.go` | CORSヘッダー付与、Core APIとMiddleware Hubの役割分担 |
| **DBテーブル・カラム拡張** | `docs/STORAGE-001_1.1_...` | `backend/models/models.go`<br>`backend/crud/` | GORMタグの正確な指定、リレーション定義 |
| **サルベージ・バッチ処理改修** | `docs/INFRA-002_...` | `wayback_tweet_rescure/` | 非常駐サイドカー構成、DBトランザクション安全化 |
