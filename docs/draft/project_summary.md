# プロジェクト概要: x_timeline_app

## 概要
x_timeline_app は Twitter のデータとメディアコンテンツを可視化するためのタイムラインアプリケーションです。Vue.js フロントエンド、Go 製バックエンドのマイクロサービスアーキテクチャ、データ処理・プロキシを担うミドルウェアレイヤーで構成されています。

## アーキテクチャ

### フロントエンド
- **フレームワーク**: Vue.js 3.5（TypeScript 使用）
- **ビルドツール**: Vite
- **スタイル**: Tailwind CSS
- **主要コンポーネント**:
  - `App.vue` - サイドバー、ユーザーセレクタ、フィルタタブを備えたメインアプリケーションシェル
  - `TweetCard.vue` - 個別のツイート表示コンポーネント
  - `MediaOverlay.vue` - メディア閲覧用オーバーレイ
  - `AppSidebar.vue` - ナビゲーションサイドバー
  - `StatusBanner.vue` - ローディング・エラー状態の表示
  - `UserSelector.vue` - アカウント選択
  - `FilterTabs.vue` - タイムラインフィルタリング
  - `HeroHeader.vue` - プロフィールヘッダーセクション
- **Composables**:
  - `useTimeline.ts` - タイムラインデータの取得と状態管理
  - `useHeartbeat.ts` - 定期的なヘルスチェック
  - `useMediaOverlay.ts` - メディアオーバーレイの管理
- **API 連携**: Axios による REST API でバックエンドと通信

### バックエンド（Go）
- **ポート**: 5176（CRUD API）、5174（Stash プロキシ）
- **フレームワーク**: 標準ライブラリ HTTP + GORM ORM
- **データベース**: SQLite（GORM 使用）
- **主要テーブル**:
  - `accounts` - ユーザーアカウント情報
  - `tweets` - 完全なメタデータを持つ Twitter 投稿（ID、タイムスタンプ、エンゲージメント指標）
  - `media` - ダウンロード状況を追跡する関連メディアファイル
- **エンドポイント**:
  - `GET /api/posts?account_id=XXX&filter=XXX&limit=XXX&offset=XXX` - アカウント別投稿取得
  - `GET /api/accounts` - 全アカウント一覧
  - `GET /api/account?id=XXX` - 特定アカウントの詳細情報取得
  - `GET /api/translate` - 翻訳エンドポイント

### ミドルウェア（Go）
- **ポート**: 5175
- **役割**: データ集約とプラグイン処理
- **機能**:
  - Stash 連携のためのリバースプロキシ（ポート 5174）
  - プラグインによる Twitter データ変換
  - バックエンドからのアカウントデータ集約
  - 翻訳 API のプロキシ
  - フレームワーク固有のレンダリングプラグイン

### データベーススキーマ
- **chronoarchive.db** - メインアプリケーションデータベース
- **archive.db** - アーカイブ保存用のセカンダリデータベース
- モデル: `Account`、`AccountProfileHistory`、`Tweet`、`Media`

## 主要機能
- アカウントフィルタリングによるタイムライン可視化
- メディアの表示とオーバーレイサポート
- アカウント管理と切り替え
- プラグインベースのフレームワークレンダリング（Twitter など）
- アーカイブコンテンツの Wayback Machine 統合
- CORS 対応 API エンドポイント
- リアルタイムのローディング・エラー状態ハンドリング

## 開発環境
- **最近の活動**: 活発に開発中（2026-08-14 のコミット）
- **ツール**: Playwright（テスト用）、Firecrawl（ウェブスクレイピング用）
- **設定**: `.kilo/kilo.json` でエージェント管理を設定
- **パッケージ管理**: npm（フロントエンド）、Go modules（バックエンド）

## ポート設定
| サービス | ポート | 説明 |
|---------|------|------|
| フロントエンド開発サーバー | 5173（デフォルト） | Vue/Vite 開発サーバー |
| バックエンド API | 5176 | CRUD 操作とデータ配信 |
| ミドルウェアプロキシ | 5175 | データ集約とプラグイン処理 |
| Stash リバースプロキシ | 5174 | 外部サービスへのプロキシ |

## 依存関係
- **フロントエンド**: vue@3.5.40、vite@8.2.0、vue-tsc、tailwindcss@3.4.19
- **バックエンド**: go-sqlite3、gorm、github.com/glebarez/go-sqlite
- **テスト**: backend/api、backend/crud ディレクトリに各種テストファイルを配置