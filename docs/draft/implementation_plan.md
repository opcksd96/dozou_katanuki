# フロントエンド描画エンジン再構築 導入計画

## 🎯 目標とパラダイム
旧システムから抽出した要件（リッチなメディア表示、Stash連携、統計情報）を満たす新フロントエンドを構築します。
以下の**絶対要件**に準拠します：
1. **1ファイル100行以下**: コンポーネントおよびロジックの極限までの細分化。
2. **関数型プログラミング基盤**:
   - **アトミック性 / 普遍性**: UIコンポーネントは状態を持たない（Stateless / Pure）プレーンな描画関数として設計。
   - **透過性**: データ変換や計算は純粋関数（Pure Functions）として `utils/` に分離。
   - **副作用の部分的許容**: API通信やDOM操作（ライフサイクル）などの副作用（Side Effects）は、Vue 3のコンポーザブル（`composables/`）にのみカプセル化して隔離。

## ⚠️ User Review Required
> [!IMPORTANT]  
> 1ファイル100行以下の制約を守るため、Vueコンポーネント（SFC）は「ロジック（script）」「構造（template）」「装飾（style）」のうち、スタイルをグローバルのTailwind/CSSに逃がす、あるいはロジックをコンポーザブルに完全分離する必要があります。この粒度（コンポーネントの細分化）で進めてよろしいでしょうか？

## 🏗️ アーキテクチャ設計とディレクトリ構成

### 1. ドメイン・モデル層 (副作用なし)
- `[NEW]` `frontend/src/models/RenderTree.ts`: `RenderTree`, `MediaItem`, `Author` などの型定義のみ。
- `[NEW]` `frontend/src/models/StashData.ts`: Stash連携用の型定義。

### 2. ユーティリティ層 (純粋関数 / 純度100%)
- `[NEW]` `frontend/src/utils/formatters.ts`: 日付や数値のフォーマット（純粋関数）。
- `[NEW]` `frontend/src/utils/parser.ts`: ツイート本文のリンク化・パース処理（純粋関数）。

### 3. インフラ・副作用層 (API通信)
- `[NEW]` `frontend/src/api/client.ts`: `fetch` をラップする純粋な非同期関数。
- `[NEW]` `frontend/src/api/timeline.ts`: タイムライン取得APIリクエスト。

### 4. 状態管理・作用層 (Composables)
- `[NEW]` `frontend/src/composables/useTimeline.ts`: データのFetch、Loading状態、Error状態を管理するリアクティブ副作用。
- `[NEW]` `frontend/src/composables/useMediaOverlay.ts`: Lightboxの開閉状態管理。

### 5. プレゼンテーション層 (UIコンポーネント / 100行制約)
VueのSFCを極小化し、単一責任の原則を貫きます。

#### 構造・レイアウト
- `[MODIFY]` `frontend/src/App.vue`: アプリのルート。
- `[NEW]` `frontend/src/components/layout/AppSidebar.vue`: ナビゲーション。
- `[NEW]` `frontend/src/components/layout/StatusBanner.vue`: サーバー状態表示。

#### ツイートカードコンポーネント群
- `[NEW]` `frontend/src/components/tweet/TweetCard.vue`: 各サブコンポーネントを束ねる親。
- `[NEW]` `frontend/src/components/tweet/TweetAuthor.vue`: アバターと名前表示。
- `[NEW]` `frontend/src/components/tweet/TweetBody.vue`: テキスト表示と翻訳トグル。
- `[NEW]` `frontend/src/components/tweet/TweetStats.vue`: いいね・RTのアクションバー。
- `[NEW]` `frontend/src/components/tweet/MediaGrid.vue`: 1〜4枚のメディアを自動レイアウト。

#### メディア・Stash連携コンポーネント群
- `[NEW]` `frontend/src/components/media/MediaOverlay.vue`: Lightbox本体。
- `[NEW]` `frontend/src/components/media/StashPlayer.vue`: 動画プレイヤー部分。
- `[NEW]` `frontend/src/components/media/StashMetaSidebar.vue`: メディアの技術スペックや出演者表示。

## 🚀 検証プラン (Verification Plan)

### 実装ステップの自動検証
1. **静的解析**: 各ファイルが100行を超えていないことを `wc -l` 等でスクリプト検証。
2. **純粋関数のテスト**: データのフォーマッタやパーサーが、外部状態に依存せず同一入力に対し同一出力を返すか確認。
3. **E2Eマニュアル確認**: 
   - `UNKNOWN_msluo14` の実データを読み込み、カード、メディアグリッド、Stashメタデータが正常に表示されるかブラウザ上で確認。
   - サイドエフェクト（APIローディング等）がコンポーネントツリーを汚染せず透過的に動作しているかVue DevToolsで確認。
