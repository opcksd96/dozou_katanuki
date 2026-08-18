# フロントエンド描画エンジン要件定義レポート

旧実装（`orz/`ディレクトリ内）のコンポーネント群（`App.vue`, `TweetCard.vue`, `MediaOverlay.vue`, `StashPlayer.vue` 等）を解析し、新フロントエンドの描画エンジンに必要な**データスキーマ（RenderTree）**と**表示・UI機能**をまとめました。

## 1. 要求されるデータスキーマ (RenderTree)

Middlewareからフロントエンドへ供給すべき、抽象化された投稿データの構造要件です。

### 旧バックエンドから引き継ぐべきデータ構造（`orz/backend` 解析結果）
以前の実装では、バックエンドがデータベースのモデル（`types.go`）を直接JSONとしてフロントエンドに返却していました。この構造から、新アーキテクチャ（Middleware経由のプラグイン変換）の `RenderTree` が包含すべきすべての情報が明確になりました。

旧API（`/api/timeline` 等）が返却していた主要な生データ（Raw）要素：
- **統計情報**: `reply_count`, `retweet_count`, `like_count`, `bookmark_count`, `view_count`
- **投稿ステータス**: `is_retweet`, `reply_to_tweet_id`, `conversation_id`, `wayback_url`, `source_type`
- **メディアのStash直結ID**: `stash_scene_id`, `stash_image_id`, `width`, `height`

これらは新アーキテクチャにおいて、**Backend API → Middleware（Twitterプラグイン）** の過程で失われずに `RenderTree` へマッピングされる必要があります。

### 基本投稿構造 (RenderTree)
- `id` (String): 投稿のユニークID
- `framework` (String): データソースの種別（例: "twitter"）
- `author` (Object):
  - `name` (String): 表示名（Display Name）
  - `handle` (String): ユーザー名（@username）
  - `avatar_url` (String): アバター画像のURL
- `content` (String): 投稿の本文（HTMLレンダリング可能なパース済みテキスト）
- `created_at` (String): 投稿日時
- `media` (Array<MediaItem>): 添付メディアの配列（最大4つを想定）

### メディア構造 (MediaItem)
- `type` (String): メディア種別（"image", "photo", "video" など）
- `url` (String): メディアの表示・再生用URL
- `stash_id` / `media_id` (String): Stash側のIDまたはDB上の管理ID

### 統計・拡張メタデータ
- `reply_count` (Number): リプライ数
- `retweet_count` (Number): リツイート数（UI上では "Salvaged" として表示）
- `like_count` (Number): いいね数
- `view_count` (Number): 閲覧数
- `is_translated` (Boolean): 翻訳済みかどうかのフラグ
- `translated_text` (String): 日本語訳のテキスト（翻訳機能用）
- `quoted_tweet` (Object): 引用リツイート情報（存在する場合）

---

## 2. 必須となる表示機能（UIコンポーネント要件）

旧アーキテクチャが実現していたリッチなユーザー体験を維持・向上させるため、以下の機能実装が必要です。

### 📱 アプリケーション基盤 (`App.vue`)
- **サーバー生存監視（Heartbeat）**: 3秒間隔でサーバー状態を監視し、オフライン時には固定バナーとサイドバーで警告を表示。
- **レスポンシブ・ナビゲーション**: PC版では左サイドバー、モバイル版（768px以下）では画面下部のボトムナビゲーションバーとして機能。

### 🕊️ 投稿カード (`TweetCard.vue`)
- **X(Twitter) 準拠のデザイン**: 認証バッジ（チェックマーク）やブランドロゴの配置。
- **メディアグリッド**: メディアの数（1〜4つ）に応じたダイナミックなグリッドレイアウト。
- **動画プレビュー再生**: 動画の場合、サムネイル代わりにミュート状態での自動プレビュー（`playsinline`, `muted`, `autoplay`）を行い、中央にPlayアイコンをオーバーレイ配置。
- **翻訳機能（Toggle）**: ボタンクリックで日本語訳（`translation-box`）をインライン展開するトグル機能。
- **アクションバー**: リプライ、リツイート（Salvaged）、いいねのアクションアイコンとカウンター。

### 🖼️ メディアオーバーレイ (`MediaOverlay.vue`)
- **フルスクリーンLightbox**: メディアをクリックした際に、背景を暗転させて全画面表示するモーダル機能。
- **サイドバー情報**: メディアの横に、投稿者の情報、本文、統計情報を再掲するサイドパネル。
- **Stashメタデータ連携**: メディアがStashと紐付いている場合、以下の情報をサイドバー下部に表示。
  - バッジ（"🎬 Stash Scene" / "🖼️ Stash Image"）
  - Stashのお気に入り（★）状態
  - 技術スペック（解像度、再生時間、ビデオコーデック）
  - パフォーマー（出演者）のアイコンリストとリンク
  - タグリスト（リンク）

### 🎬 高機能動画プレイヤー (`StashPlayer.vue`)
- **HLSストリーミング対応**: Stashからの動画ストリームを再生するカスタムプレイヤー。
- **シークバーとサムネイルプレビュー**: タイムラインホバー時にVTT形式のSpriteサムネイル（スナップショット）を表示。
- **レジューム機能（Resume）**: 途中まで見た動画の再生位置を記憶し、再開時に「Resumed from XX:XX」というトースト通知と最初から見直す（Restart）ボタンを提示。
- **A-Bループ**: 特定の区間（Loop Start / Loop End）をリピート再生し、シークバー上にループ範囲を視覚化（Marker）する機能。
- **フォールバック**: HLS再生等のFatal Error時に「MP4 Direct Stream」へ切り替えるリトライボタン。
- **コマ送り・コマ戻し**: ホットキーやボタンによる精密なフレーム単位の制御。
