# x_timeline_app 技術概論 & アーキテクチャ白書 (General Architecture & Technical Specification)

**ドキュメントID**: SPEC-WIKI-001  
**バージョン**: 2.0.0  
**作成日**: 2026-08-17  
**ステータス**: 正式版（Wikiマスター）  
**対象システム**: x_timeline_app (SNS Time-Capsule & Salvage Archive System)

---

## 📑 目次 (Table of Contents)

- [第1章：技術仕様とバックボーン](#第1章技術仕様とバックボーン)
- [第2章：外部サービスの概要とサルベージ技術](#第2章外部サービスの概要とサルベージ技術)
- [第3章：実装における要点・注意点（AI駆動開発と制約原則）](#第3章実装における要点注意点ai駆動開発と制約原則)
- [第4章：データベース設計と仮想ストレージプール](#第4章データベース設計と仮想ストレージプール)
- [第5章：フロントエンド概論（Foolish Frontend）](#第5章フロントエンド概論foolish-frontend)
- [第6章：分厚いミドルウェア（Middleware Hub）](#第6章分厚いミドルウェアmiddleware-hub)
- [第7章：堅牢なバックエンド（Driver & Data Abstraction）](#第7章堅牢なバックエンドdriver--data-abstraction)
- [第8章：ローカルストレージへのデータ永続化と保全ストラテジ](#第8章ローカルストレージへのデータ永続化と保全ストラテジ)
- [第9章：プラグインアーキテクチャ概説](#第9章プラグインアーキテクチャ概説)
- [第10章：Admin Board（統合管理基盤と運用設計）](#第10章admin-board統合管理基盤と運用設計)
- [第11章：参考資料・技術文献](#第11章参考資料技術文献)

---

## 第1章：技術仕様とバックボーン

### 1.1 プロジェクトの目的
公式プラットフォーム（Twitter/X, Instagram, Facebook等）におけるアカウント凍結（Ban）、投稿削除、規約変更によるAPI閉鎖、Webサービス自体の終了などにより、人類のデジタルな足跡や記憶は日々失われています。
本プロジェクト **`x_timeline_app`** は、失われたアカウントや投稿を **Wayback Machine** や **Aria2** などの外部アーカイブ技術・分散ダウンロード技術を駆使してWebの深淵からサルベージ（救出・復元）し、ローカル環境において当時のタイムラインの質感・レスポンスそのままに**「動態保存（動作可能な状態で永続化）」**することを至上命題としています。

### 1.2 使用言語とライブラリスタック
システムはマイクロサービス的なモジュール分離を行いながら、最終的に単一バイナリにパッキング可能な技術選定を行っています。

```
┌─────────────────────────────────────────────────────────────┐
│ Frontend: Vue 3.5.40 (SFC) + Vite 8.2.0 + TypeScript + Tailwind CSS
├─────────────────────────────────────────────────────────────┤
│ Middleware: Go 1.22+ (net/http, RendererPlugin Engine)
├─────────────────────────────────────────────────────────────┤
│ Backend: Go 1.22+ + GORM ORM + go-sqlite3 + Reverse Proxy
├─────────────────────────────────────────────────────────────┤
│ Salvage Sidecar: Python 3.10+ (Playwright, Requests, CDX Parser)
├─────────────────────────────────────────────────────────────┤
│ Storage / Media: SQLite3 + Stash Media Server (GraphQL / Web Streams)
└─────────────────────────────────────────────────────────────┘
```

### 1.3 データスキームの根底思想
- **「Same Source, Same Flow（同一の源流、同一の血流）」**: データの参照同一性を保ち、二重管理やアドホックな同期を排除。
- **ハイブリッド永続化**: メタデータ（テキスト、エンゲージメント、リレーション）は **SQLite3** で高速クエリ可能にし、大容量メディア（画像・動画）は **Stash** をバックエンドストレージプールとして透過的に統合。

### 1.4 コーディング規約の基本方針
- **Strict 100-Line Limit**: コンポーネントおよびロジックの極限までの単一責任化（1ファイル100行以下）。
- **Stateless & Pure**: UIコンポーネントは状態を持たない（Foolish/Pure）描画関数とし、副作用はComposablesに、純粋計算はUtilsに隔離。

---

## 第2章：外部サービスの概要とサルベージ技術

### 2.1 SNSプラットフォームと凍結アカウントの問題点
X (Twitter)、Instagram、Facebookなどの商用SNSは中央集権型のクラウドサービスであり、以下の構造的課題を抱えています：
- **アカウント凍結による全データ消滅**: アカウントがサスペンドされると、メディアファイルへのCDNリンク、リプライツリー、いいね履歴へのアクセスが完全に遮断される。
- **短縮URL（t.co 等）のリンク切れ**: 公式サーバーの停止に伴い、短縮リンクのリダイレクト解決が不可能になる。
- **メディアの暗号化・再エンコード**: プラットフォーム側の仕様変更で原画質が劣化するリスク。

### 2.2 Wayback Machine（Internet Archive）の技術的内容とその限界
- **CDX Server API**: URLスナップショットの一覧をJSON/テキストで高速取得可能 (`web.archive.org/cdx/search/cdx`)。
- **Timemap / Memento プロトコル**: 過去のスナップショット群のタイムラインを走査。
- **限界と課題**:
  - 動的JavaScriptで生成されるツイート本文やメディアURLがHTMLスナップショット内に含まれていない場合がある（JSレンダリングが必要）。
  - レート制限およびアーカイブ欠落（メディア実ファイルが保存されていないケース）。

### 2.3 Aria2 による P2P/マルチセッション通信とキャッシュサルベーション
- **高速並列ダウンロード**: 分散されたアーカイブサーバーやCDNミラーから複数接続（Multi-Connection）で高速取得。
- **フォールトトレラント**: 中断されたダウンロードの再開（Resume）およびハッシュ検証。

### 2.4 Stash によるメディア管理概論
- **Stashapp**: 大規模な画像・動画ファイルのハッシュベース管理、トランスコード、ストリーミング配信を行うオープンソースメディアサーバー。
- **メリット**: ファイルの重複排除（Perceptual Hash）、サムネイル自動生成、GraphQL APIによるメタデータ操作。

### 2.5 Motrix の RPC API
- Aria2 を内蔵した Motrix の JSON-RPC (`aria2.addUri`, `aria2.tellStatus`) を介して、サルベージパイプラインからのバックグラウンド一括ダウンロードキューを制御。

### 2.6 SQLite3 による独自データベース管理
- 外部クラウドに依存しない完全自己完結型のローカルRDBMS。WAL（Write-Ahead Logging）モードを採用し、高スループットな読み書きとポータビリティを両立。

---

## 第3章：実装における要点・注意点（AI駆動開発と制約原則）

### 3.1 Simple コーディングと関数型パラダイム
複雑性を徹底的に排除する宣言的コーディングを採用します。
- **イミュータビリティ（不変性）**: データの破壊的変更を禁止し、新しいオブジェクトを生成して返却。
- **純粋関数（Pure Functions）**: 入力引数のみに依存し、外部状態を変更しない関数を `frontend/src/utils/` に集約。

### 3.2 「1ファイル 100行以下」と「1責務 1モジュール」制約
AIエージェントのコード生成精度を最大化し、メンテナンス性を高めるための物理制約です。
- **SFC（Vue）の分割**: テンプレート、ロジック、スタイルが100行を超える場合、即座にサブコンポーネント（Author, Body, Stats, Media等）または Composable に分離する。
- **関数の細分化**: 1つの関数は20行以内を目安とし、処理フローが一眼で把握できるようにする。

### 3.3 コンテキスト爆発（Context Explosion）の回避
- AIが作業する際、コードベース全体を毎回全走査させるとトークンが枯渇し、ハルシネーションが発生します。
- 本概論書および `docs/reference.md` を起点とした「ドキュメント駆動開発（Doc-Driven Coding）」により、必要なファイルのみをピンポイントで特定・編集するワークフローを徹底します。

---

## 第4章：データベース設計と仮想ストレージプール

### 4.1 Stashapp メディア管理ツールと GraphQL API
Stash は GraphQL エンドポイント（`POST /graphql`）を提供し、`Scene`（動画）および `Image`（静止画）としてメディアを管理します。
本システムは Stash の GraphQL API を直接叩くだけでなく、メディアをプロキシ経由でストリーミング再生します。

### 4.2 SQLite3 による Twitter互換データベース
Stash 単体では「テキストのみのツイート」や「リプライ階層関係」「エンゲージメント数値」を保持できません。これを補完するのが独自 SQLite データベース（`archive.db` / `chronoarchive.db`）です。

```
┌─────────────────────────────────────────────────────────────┐
│                      Tweet Entity (SQLite)                  │
│  - TweetID, NumericID (Author), Content, CreatedAt          │
│  - Metrics (Likes, RTs, Bookmarks), WaybackURL, Status      │
└──────────────────────────────┬──────────────────────────────┘
                               │ 1 : N リレーション
┌──────────────────────────────▼──────────────────────────────┐
│                      Media Entity (SQLite)                  │
│  - MediaID, Type (image/video), URL, DownloadStatus         │
│  - StashSceneID (FK -> Stash Scene), StashImageID (FK)      │
└──────────────────────────────┬──────────────────────────────┘
                               │ 仮想リンク
┌──────────────────────────────▼──────────────────────────────┐
│                    Stashapp Media Pool                      │
│  - Scene / Image Binary, Transcoded Streams, Thumbnails     │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 仮想データベースレイヤ（Unified Data Layer）の設計
- **`Tweet` と `Media` の透過的マッピング**: GORM を介して `Tweet` 取得時に自動で `Media` を Preload し、Stash の ID が存在すれば Stash プロキシ URL、未保存なら Wayback/オリジナル URL を自動フォールバック設定。

---

## 第5章：フロントエンド概論（Foolish Frontend）

### 5.1 Foolish Frontend（愚直で無知なフロントエンド）の実現
フロントエンドは「受け取った `RenderTree` を忠実に描画すること」に専念し、ビジネスロジックやデータソースの差異を知る必要がありません。
- **Stateless Components**: コンポーネント自身は API の場所や DB スキーマを知らず、Props を受領して UI を描画するだけの純粋 View。

### 5.2 フロントエンドとバックエンドの責務分離
- **Frontend**: DOM イベント、キーボードショートカット、アニメーション、ユーザー操作のハンドリング。
- **Middleware**: データの集約、SNS プラグイン変換、外部 URL 解決。
- **Backend**: DB クエリ、トランザクション、永続化。

### 5.3 Frontend Plugin & 拡張機能
- **翻訳機能（Translation）**: `TweetBody.vue` から `/api/translate` を呼び出し、原文と翻訳文をトグル表示。
- **メディアプレイヤー（Player & Lightbox）**: `MediaOverlay.vue` による全画面プレビュー、動画のシーク・ループ再生。
- **いいね・ブックマーク**: ローカル状態および DB フラグ（`is_liked`）と連動するリアクティブな評価ボタン。

### 5.4 CSS コーディング規約
- **Tailwind CSS 3.4+**: スタイルのカプセル化と CSS ファイル肥大化防止のため、ユーティリティクラスを優先。
- **ダークモード / Glassmorphism**: SNS の世界観に没入できるダークテーマ・半透明ブラーエフェクト（`backdrop-blur`）を標準採用。

---

## 第6章：分厚いミドルウェア（Middleware Hub）

### 6.1 MVVMモデルにおけるミドルウェアの位置づけ
ミドルウェア（Go :5175）は、バックエンドの生データ（Model）をフロントエンドの表示モデル（ViewModel / `RenderTree`）へと変換・集約する**「インテリジェント・ハブ」**として機能します。

### 6.2 データディスパッチとプラグインエンジン
```go
type RendererPlugin interface {
    Identifier() string
    Transform(posts []RawPost, accountResolver func(string) RenderAuthor) []RenderTree
}
```
- クエリパラメータ `?framework=twitter`（または `instagram` 等）に応じて動的にプラグインを選択し、JSON レスポンスを構築。

### 6.3 通信規格・相対プロキシ・ルーティング
- **ポート 5174 (Stash Proxy) / 5175 (`/stash-proxy/*`)**: ブラウザの CORS 制約を回避し、ローカルの Stash (`:9999`) へのリバースプロキシを**完全相対パス**で提供。`localhost` ハードコードを全廃し、スマホやLAN環境からのアクセスでもゼロコンフィグで動作。
- **SPA Fallback ルーティング**: ブラウザによるURL直打ち（`/:platform/:username/status/:id` 等）やリロード時、存在しない静的ファイルパスに対して `index.html` を返却し、画面遷移を完全に保護。
- **ポート 5175 (Middleware Hub)**: フロントエンドからの `/api/render`, `/api/accounts`, `/api/translate` を統括。

---

## 第7章：堅牢なバックエンド（Driver & Data Abstraction）

### 7.1 ドライバー層としての Go Core API (:5176)
- **GORM による抽象化**: SQLite の低レベル操作を隠蔽し、型安全な CRUD 操作を提供。
- **高速ページネーション**: `Limit` と `Offset`、`created_at desc` インデックスを活用した高速スライス取得。
- **統一URLインターフェース（Unified Media URLs）**: `urls.stream`, `urls.image`, `urls.thumbnail`, `urls.original` をドライバー層で構造化して応答し、上位レイヤーでのURL組み立てスパゲッティを根絶。

### 7.2 統合データベース仮想レイヤ
テキストのみのツイートも、4枚の画像付きツイートも、動画付きツイートも、すべて単一の `ResponsePost` 構造体として正規化されてミドルウェアに渡されます。

---

## 第8章：ローカルストレージへのデータ永続化と保全ストラテジ

### 8.1 メディアファイルの保存規格（URL BaseName 原則）
- **動画**: MP4 (H.264 / AAC) または WebM。
- **画像**: JPEG / PNG / WebP。
- **命名規則**: **URLのベースネーム（BaseName）をそのまま物理ファイル名として採用**（例: `eb7ymRi-pfsx5FJH.mp4`, `F8wZ1abXYAAY7kL.jpg`）。  
  ※Wayback MachineやX CDNとのURL突合・キャッシュヒット判定を $O(1)$ で行うため、独自リネームは行わずURL末尾と完全一致させます。

### 8.2 Stashメディアプールとアバター管理の完全分離（Avatar Isolation）
- **Stash メディアプール**: ツイート添付の高解像度作品（動画 `Scene` / 静止画 `Image`）のみを管理。
- **アバター＆UIアセット**: `middleware/assets/{framework}/` に隔離保存し、Stashの監視フォルダから除外。Stashの `Performer.image` にはData URIで埋め込み、StashのImageテーブルへのアイコン混入（ライブラリ汚染）を100%防止。

### 8.3 Whitelist によるデータフィルタリング機構
- アーカイブ対象アカウントのホワイトリスト管理により、不要なボットデータや無関係なメディアの混入を防止。

### 8.4 バックアップとデータ保全
- SQLite の `.backup` コマンドまたはスナップショットコピー（`archive.db`）による定期バックアップ。
- メディア実ファイルは Stash のストレージディレクトリ内で RAID または外部ストレージへ退避可能。

---

## 第9章：プラグインアーキテクチャ概説

### 9.1 レンダリングプラグイン（Go）
`middleware/plugins/` 配下に新しい SNS プラグインを追加することで、UI を変更せずに他 SNS のタイムライン表示に対応。

### 9.2 スクレイパー・サルベージプラグイン（Python Sidecar）
- `Wayback_CDX.py`: Wayback Machine から URL 一覧を取得。
- `salvage_pipeline.py`: ツイート本文のパースと HTML 復元。
- `stash_injector.py`: ダウンロードしたメディアを Stash API へ登録。

### 9.3 プラグイン登録ワークフロー
1. Go 側: `plugins.Register(&MySNSPlugin{})` で初期化時に登録。
2. Python 側: CLI またはバッチからサブプロセスとしてオンデマンド実行。

---

## 第10章：Admin Board（統合管理基盤と運用設計）

### 10.1 設定管理アーキテクチャ (`settings.yaml` / `config.json`)
- ポート番号（`backend_port`, `stash_port`）、デフォルトフレームワーク（`default_framework`）の一元管理。
- 実行時の動的ロード・セーブをサポート。

### 10.2 提供される管理機能
- **テーマ・配色・レイアウト・言語設定**: UI のカスタム設定。
- **DB メンテナンスツール**: 孤立したメディアレコードのクリーンアップ、破損リンクの修復。
- **Whitelist 管理 UI**: サルベージ対象アカウントの追加・除外。
- **今後の展望**: Web ベースの統合 GUI ダッシュボードによるワンクリック・サルベージ実行。

---

## 第11章：参考資料・技術文献

1. **Wayback Machine APIs**:
   - Internet Archive CDX Server API Reference
   - Memento: Time Travel for the Web (RFC 7089)
2. **Stashapp Community Documentation**:
   - Stash GraphQL Schema & API Specification
3. **Vue.js 3 & Vite Documentation**:
   - Vue 3 Composition API Guide & RFCs
   - Tailwind CSS Utility-First Fundamentals
4. **Go Database & Concurrency**:
   - GORM Guide (Object Relational Mapping for Go)
   - Go net/http & httputil ReverseProxy Architecture
