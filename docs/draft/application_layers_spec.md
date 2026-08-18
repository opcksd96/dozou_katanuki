# アプリケーション層 アーキテクチャ仕様書 (Application Layers Specification)

**プロジェクト名**: x_timeline_app (SNS Timeline & Media Archival System)  
**ドキュメントID**: SPEC-LAYER-001  
**バージョン**: 2.0.0  
**作成日**: 2026-08-17  
**ステータス**: 正式仕様（実装照合・改定完了）

---

## 1. 全体アーキテクチャ概要

本システムは、Twitter (X) をはじめとするSNSの投稿メタデータおよびメディア（静止画・動画）をアーカイブ・統合管理し、柔軟なタイムラインUIとしてレンダリングするためのマルチレイヤー・アーキテクチャを採用しています。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. プレゼンテーション層 (Presentation Layer)                                 │
│    - Vue 3.5 + TypeScript + Tailwind CSS                                    │
│    - Timeline UI / TweetCard / MediaGrid / MediaOverlay / FilterTabs        │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │ HTTP / REST API (JSON)
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ 2. コンテンツディスパッチャー層 (Content Dispatcher / Middleware Layer)      │
│    - Go HTTP Middleware (:5175) & Reverse Proxy (:5174)                     │
│    - フレームワーク変換 (RendererPlugin: Twitter, etc.)                      │
│    - 外部サービス再接続 (Wayback Machine, Stash Media URL, Translation)     │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │ 内部 REST / GORM
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ 3. ドライバー層 (Driver / Core Backend Layer)                                │
│    - Go Core Backend (:5176) + GORM                                         │
│    - SQLite3 (archive.db / chronoarchive.db) & Stash メディア統合           │
│    - テキスト単体投稿とメディア統合エンティティの抽象化                    │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │ IPC / CLI 呼び出し
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ 4. プラグイン層 (Plugin / Downloader Sidecar Layer)                         │
│    - Python サルベージ / スクレイピングスクリプト群 (非常駐サイドカー)      │
│    - Wayback CDX / X API / HTML解析 / Stash Injector                        │
│    - メディアダウンロード & メタデータSQLite登録                             │
└─────────────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. コントローラー層 (Controller / Administration & Configuration Layer)     │
│    - システム設定管理 (settings.yaml / config.json / Settings API)          │
│    - UI挙動設定・評価システム (Bookmarks / Pins)・アセット管理              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. 各層の仕様定義・実装照合・現状仕様

---

### 第1層：プレゼンテーション層（Presentation Layer / Vue）

#### 1. 設計定義と責務
Webブラウザへのレンダリングを統括する層。
- **レンダリングエンジン切り替え**: バックエンド/ミドルウェア側のレンダリングプラグイン（RendererPlugin）と連携し、各SNS用の表示モード（デフォルトはTwitter/X互換モード）をコントロール。
- **統合タイムライン表示**: 投稿の文字情報、画像、動画メディア、各種統計・エンゲージメント情報を統合表示。
- **コア共通インフラ**:
  - タイムライン形式表示（`TweetCard`, `TweetBody`, `MediaGrid` 等）
  - 単一投稿の個別・詳細表示
  - メディアリストギャラリーおよびメディア表示オーバーレイ（`MediaOverlay`）
  - フィルタリング・アカウント切り替えUI（`FilterTabs`, `UserSelector`, `HeroHeader`）
- **拡張機能プラグイン**:
  - ツイート翻訳機能（`translate_api` 連携）
  - 動画の字幕対応、画像ポートフォリオ表示機能

#### 2. 実装との照合（一致・不一致分析）

| 項目 | 設計定義 | 現状の実装コード | 判定・差異 |
| :--- | :--- | :--- | :--- |
| **ベースフレームワーク** | Vue 3 + TypeScript | `frontend/src` (Vue 3.5.40, Vite, Tailwind CSS 3.4.19) | **一致** |
| **SNSプラグイン連携** | プラグイン形式で表示層切替（デフォルトTwitter） | `middleware/plugins/twitter.go` から生成される `RenderTree` を受領してレンダリング | **一致**（Goミドルウェア主導で変換） |
| **コアUI機能** | タイムライン、メディアオーバーレイ、個別表示 | `App.vue`, `TweetCard.vue`, `MediaGrid.vue`, `MediaOverlay.vue`, `FilterTabs.vue` 等で実装完了 | **一致** |
| **翻訳機能** | ツイート翻訳機能プラグイン | `frontend/src/components/tweet/TweetBody.vue` から `/api/translate` を呼び出し可能 | **一致**（Composable分離の最適化中） |
| **字幕・ポートフォリオ** | 動画字幕対応、画像ポートフォリオ表示 | メディアグリッドやプレイヤースタブは存在するが、字幕トラック解析・専用ポートフォリオビューは未完全 | **一部不一致（拡張機能が部分実装）** |

#### 3. 現状の確定仕様
- フロントエンドは Vue 3 Composition API + TypeScript で構築され、コンポーネントはMVVMパターン（View: `components/`, ViewModel: `composables/useTimeline.ts`）に準拠。
- バックエンドから受領した `RenderTree`（Author, Content, Media, Meta）を元に、Twitter (X) 互換スタイルのリッチなタイムラインとメディアモーダルを展開する。
- **Twitter互換 URLルーティング体系（Vue Router / HTML5 History API）**:
  - `/:platform`（例: `/twitter` または `/`）: 統合タイムライン
  - `/:platform/:username/`（例: `/twitter/msluo14/`）: 特定アカウントの個別タイムライン
  - `/:platform/:username/status/:id`（例: `/twitter/msluo14/status/1879382757924868404`）: 単一ツイート個別詳細・スレッド表示
  - `/settings`: システム管理・DB統計画面

---

### 第2層：コンテンツディスパッチャー層 (Content Dispatcher / Middleware Layer)

#### 1. 設計定義と責務
フロントエンドからのリクエストを受け付け、適切なデータ形式で返却するバックエンドの受け口・ハブ。
- **データ配信**: 投稿（ツイート）を1単位として、指定アカウントやフィルタに応じたタイムライン・スレッドデータをJSON形式で返却。
- **外部リンク・メディア再接続**:
  - Wayback Machine アーカイブURL（`wayback_url`）への解決・再接続
  - Stash メディアURL（Scene ID / Image ID）へのリバースプロキシ連携（`/stash-proxy/...`）
  - 投稿内短縮URLの解決およびメディアライブラリへのアクセスポイント提供
- **フレームワーク変換**: プラグイン（`RendererPlugin`）により、生データを各SNSの共通レンダリングツリー構造に変換。
- **SPA Fallback ルーティング**: ブラウザによるURL直打ち・リロード時、存在しない静的ファイルパスに対して `index.html` を返却し、フロントエンドのルーティングを完全に保護。
- **動的ホスト対応（Dynamic Host Resolution）**: `localhost` のハードコードを全廃し、スマホやLAN内別端末（`0.0.0.0` listen）からのアクセスでも破綻しない相対パス（Relative Path）プロキシを提供。

#### 2. 実装との照合（一致・不一致分析）

| 項目 | 設計定義 | 現状の実装コード | 判定・差異 |
| :--- | :--- | :--- | :--- |
| **プロトコル形式** | REST API 準拠 | `middleware/main.go` (:5175) および `backend/main.go` (:5176) の REST API (`/api/render`, `/api/posts` 等) | **一致** |
| **データ単位と返却** | 投稿1単位・スレッド返却 | `RawPost` から `RenderTree` への変換。スレッドID (`conversation_id`, `reply_to_tweet_id`) はDBに保持されフラット返却 | **一致** |
| **Wayback / Stash 再接続** | 相対パスプロキシによる配信 | ポート5174の `NewProxyHandler` および 5175 の `/stash-proxy/*` 連携 | **改善改定（localhostハードコード撤廃・相対パス化）** |
| **フレームワークプラグイン** | SNS別レンダリングエンジン | `middleware/plugins/plugin.go` の `RendererPlugin` インターフェースおよび `twitter.go` | **一致** |
| **SPA Fallback** | URL直打ち・リロード保護 | Go HTTPルーターにおける NotFound ハンドラで `index.html` を返却 | **新仕様策定（完全直打ち対応）** |

#### 3. 現状の確定仕様
- ポート **5175** で稼働する Go 製 Middleware がディスパッチャーの中核を担う。
- エンドポイント `GET /api/render?account_id=XXX&filter=XXX&framework=twitter` を提供し、Core Backend (5176) から取得した生データを `RenderTree` に構造化してフロントエンドに配信。
- ポート **5174** のリバースプロキシおよびミドルウェアの相対ルーティングを通じて、Stash（`:9999`）のメディアストリームを完全相対パスで供給する。

---

### 第3層：ドライバー層 (Driver / Core Backend Layer)

#### 1. 設計定義と責務
Stash および Twitter互換 SQLite データベースの統合インターフェースを定義し、永続化データへのアクセスを提供する層。
- **データ構造の統合**: Stash 単体では管理不可能な「メディアを含まないテキスト単体の投稿（テキストツイート）」を補完し、投稿単位のマスターレコードを形成。
- **ORM & データベース抽象化**: SQLite3（`archive.db` / `chronoarchive.db`）に対して GORM を用いたマッピングを実施。
- **統一URLインターフェース（Unified Media URLs Interface）**:
  - `StreamURL` (`/stash-proxy/scene/{id}/stream` または `/scene/{id}/m3u8`)
  - `ImageURL` (`/stash-proxy/image/{id}/image`)
  - `ThumbnailURL` (`/stash-proxy/scene/{id}/preview`, `/stash-proxy/image/{id}/thumbnail`)
  - `OriginalURL` (外部ダウンロード元/Wayback URL)
  メディア種別に応じた完全なURL群をドライバー層で構造化して応答し、上位レイヤーでのURL組み立てスパゲッティを根絶する。
- **CRUD インターフェース**: アカウント管理、投稿取得、フィルタリング（All, Retweets, Media, Bookmarks）、ページネーション（Limit/Offset）の提供。

#### 2. 実装との照合（一致・不一致分析）

| 項目 | 設計定義 | 現状の実装コード | 判定・差異 |
| :--- | :--- | :--- | :--- |
| **統一URLインターフェース** | 構造化された相対URL群を返却 | `ResponseMediaItem` に `urls` 辞書（stream, image, thumbnail, original）を持たせる | **新仕様策定（TwiDB側での一元化）** |
| **非メディア投稿の補完** | Stashの弱点を補いテキスト単体も完全管理 | `models.Tweet` テーブルで全ツイートを管理し、`models.Media` と1対多で関連付け | **完全一致** |
| **データソース統合** | SQLite + Stashメタデータ | `archive.db` に `stash_scene_id`, `stash_image_id` を保持し、Stashと連携 | **完全一致** |
| **フィルタ・クエリ機能** | アカウント別・種別別抽出 | `GetPostsByAccount` 内で `is_retweet`, `media`, `bookmarks (is_liked)` 条件分岐 | **完全一致** |

#### 3. 現状の確定仕様
- ポート **5176** で稼働する Go 製 Backend API。
- データベース構造（GORMモデル）:
  - `Account` (`accounts`): アカウント基本情報（`numeric_id`, `username`, `avatar_local_path` 等）
  - `AccountProfileHistory` (`account_profile_history`): 表示名やバイオの変更履歴
  - `Tweet` (`tweets`): 投稿本体（`tweet_id`, `conversation_id`, `full_text`, `wayback_url`, `is_liked` 等）
  - `Media` (`media`): メディア情報（`media_id`, `type`, `download_url`, `stash_scene_id`, `stash_image_id` 等）
- テキストのみの投稿とメディア付き投稿を同一の `Tweet` エンティティとして透過的に扱い、StashのメディアIDとリレーションを結ぶ。

---

### 第4層：プラグイン層 (Plugin / Downloader Sidecar Layer)

#### 1. 設計定義と責務
ユーザーや管理者の要求に基づき、各種外部サービス（Wayback Machine, X, Instagram等）からコンテンツを収集・保全する層。
- **データ収集・スクレイピング**: 各SNSやWeb Archiveに対するスクレイピングおよびAPI通信。
- **メディアとメタデータの分散保存**:
  - メディアファイル（画像・動画）: ダウンロード後、Stash へインジェクション。
  - テキスト・メタデータ: SQLite3 データベースへ独自スキーマで格納。
- **非常駐サイドカー方式**: バックエンドプロセスを圧迫しないよう、必要なバッチやトリガー時のみ実行されるオンデマンド構成。

#### 2. 実装との照合（一致・不一致分析）

| 項目 | 設計定義 | 現状の実装コード | 判定・差異 |
| :--- | :--- | :--- | :--- |
| **実装言語・形態** | Python プラグイン形式（非常駐） | `wayback_tweet_rescure/` および `orz/` 内の Python スクリプト群（`salvage_pipeline.py`, `Wayback_CDX.py`, `stash_injector.py` 等） | **一致** |
| **対象サービス** | Wayback, X, Instagram 等 | Wayback Machine (CDX API / Timemap) および Twitter HTML/JSON サルベージが中心 | **概ね一致**（Instagram等は将来拡張枠） |
| **Stash / DB 保存** | メディア→Stash、文字→SQLite | `stash_injector.py` でStash APIに登録、`database_manager.py` で `tweets.db` / `archive.db` に保存 | **完全一致** |
| **サイドカー呼び出し** | Goからの動的キック・連携 | 現状はCLI/バッチ実行スクリプト（`build.bat`, `full_pipeline_runner.py`）連携が主 | **一部不一致（Goバックエンドからの常駐プロセス管理RPC化は未実装）** |

#### 3. 現状の確定仕様
- Python ベースのモジュール群が独立したツールセットとして機能。
- Wayback CDX APIからスナップショットを取得し、消滅したツイートのテキストおよびメディアURLを復元。
- 復元されたメディアは `stash_injector.py` を介してローカル Stash にインジェクションされ、返却された Scene/Image ID が SQLite3 の `media` テーブルに記録される。

---

### 第5層：コントローラー層 (Controller / Administration & Configuration Layer)

#### 1. 設計定義と責務
システム全体の設定管理、メディアのエンコード・整合性チェック、管理機能を提供する層。
- **設定ファイルの永続化**: `config.json`（または `settings.yaml`）によるキー・バリュー形式の設定保持。
- **UI・動作設定管理**: 配色、フォント、表示モード、ポート番号、ダウンローダーのパス設定。
- **メディア整合性・評価機能**:
  - ブックマーク（`is_liked`）やピン留めなどの評価フラグ管理。
  - データベースやメディアファイル（RDVB/独自フォーマット等）の整合性定期チェック機構。

#### 2. 実装との照合（一致・不一致分析）

| 項目 | 設計定義 | 現状の実装コード | 判定・差異 |
| :--- | :--- | :--- | :--- |
| **設定ファイル保存** | `config.json` (JSON) | `middleware/settings/settings.go` にて `settings.yaml` (YAML) として永続化 | **差異（フォーマットがYAML）** |
| **ポート・モード設定** | デフォルトSNS、ポート番号管理 | `default_framework`, `stash_port`, `backend_port` を設定可能 | **一致** |
| **評価システム** | ブックマーク・ピン留め | `models.Tweet.IsLiked`（Bookmarksタブ）としてフロント・バックエンド連携完了 | **一致** |
| **データ整合性チェック** | 定期ファイル/DBチェック機能 | `wayback_tweet_rescure/` に `cleanup_extraneous_media.py`, `cleanup_orphaned_files.py` 等が存在 | **概ね一致**（スクリプト主導で実行） |
| **Admin Board UI** | 管理画面による集中コントロール | `orz/update_admin_2.py` 等の初期UI実装が存在するが、フル機能Web Adminは開発中 | **一部不一致（Admin UIの統合推進中）** |

#### 3. 現状の確定仕様
- 設定は `settings.yaml`（`middleware/settings`）によりロード・セーブが行われ、起動時および実行時のポート・フレームワーク解決を行う。
- データベースのクリーンアップや孤立メディアの整合性確認は、専用のメンテナンススクリプト群によって保全される。

---

## 3. 現行システムのネットワーク・ポート構成

| コンポーネント | ポート | プロトコル | 役割・提供機能 |
| :--- | :--- | :--- | :--- |
| **Frontend (Vue/Vite)** | `5173` | HTTP | タイムラインUI、メディアオーバーレイ表示 |
| **Middleware (Hub)** | `5175` | HTTP (REST) | データ集約、フレームワーク変換 (`/api/render`)、翻訳プロキシ |
| **Core Backend API** | `5176` | HTTP (REST) | SQLite GORM CRUD (`/api/posts`, `/api/accounts`, `/api/translate`) |
| **Stash Reverse Proxy** | `5174` | HTTP (Proxy) | Stashメディアサーバー (`:9999`) へのCORS対応透過リバースプロキシ |
| **Stash Engine** | `9999` | GraphQL / HTTP | メディア実ファイル（動画・静止画）の管理・配信エンジン |

---

## 4. 今後の開発・改定ロードマップ

1. **プレゼンテーション層**:
   - `TweetBody.vue` からの翻訳ロジックの Composable (`useTweetTranslation.ts`) 完全分離。
   - 動画字幕（VTT/SRT）およびポートフォリオグリッドビューの追加実装。
2. **コンテンツディスパッチャー層 & コントローラー層**:
   - `settings.yaml` と `config.json` の仕様統一、および Web Admin Board からの動的設定変更 UI の統合。
   - スレッド表示（`conversation_id` に基づくツリー構築）の API 最適化。
3. **プラグイン層**:
   - Go バックエンドからの Python サルベージエンジンのプロセス起動・制御 API（サブプロセス / Job キュー）の整備。