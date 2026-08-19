# 第11章：技術カタログ＆開発・規約リファレンス
**プロジェクト名** : dozou_katanuki (Pluggable UI & Multi-Format Local Archival System "土蔵・型抜き") Pluggable UI & Multi-Format Archival System)  
**ドキュメントID** : SPEC-REFERENCE-001  
**バージョン** : 3.6.0  
**作成日** : 2026-08-17  
**ステータス** : 正式仕様（パラダイム・規約・言語・モジュール・外部アプリ・Webサービス一元技術カタログ化）  

**Navigation** : [← 前の章: 第10章：システム管理・設定・保守](10_admin_board_and_governance_v5) | [📚 目次 (Home)](Home) | [DocWiki ポータルへ戻る →](Home)

--------------------------------------------------------------------------------

## 11.1 プログラミングパラダイム (Programming Paradigms)

本システムは、ローカル環境ならではのゼロ遅延（ゼロレイテンシ）表示 [1.2] と、外部SNSの凍結やAPI閉鎖に対する「極限のデータ永続性」 [1.1] を両立させるため、以下の4つの高度なプログラミングパラダイムをソースコード全体に強制適用します。

### 1. 宣言型UI (Declarative UI) [3.2.1]
*   **概念** : 「どうDOMを操作するか（命令型）」ではなく、「データがこの状態であれば、画面はこう描画されるべきである（宣言型）」という関係性を記述します [3.2.1, 5.1.1]。
*   **適用** : フロントエンド（Vue 3.5 / `:5173`）のすべての View コンポーネントは、ミドルウェアから受領した `RenderTree` [5.3] をそのままPropsバインドして描画するだけの **Dumb Component（愚かなUI）** とします [3.2.1, 5.1.1]。
*   **効果** : 表示条件の複雑なSNSフィードであっても、コードが100%予測可能になり、描画ラグやバグを根絶します [5.1.1, 5.5.2]。

### 2. 単一データフロー (Unidirectional Data Flow / UDF) [3.2.2]
*   **概念** : データが一方向（Top-to-Bottom）にのみ流れるようにデータフローを制限します [3.2.2, 5.1.3]。
*   **適用** : **[Go Middleware Hub (:5175)] ➔ [Composable (State) :5173] ➔ [View (Props)]** の流れを絶対に遵守します [5.1.3]。下流コンポーネントが Props やグローバル状態を直接書き換える「双方向データバインディングによる暗黙の状態破壊」は厳格に禁止されます [3.2.1, 5.1.3]。
*   **効果** : 状態変更のトリガー（イベント）が常に明確になり、大規模なタイムライン展開でもバグの追跡が容易になります [5.1.3, 139]。

### 3. リアクティブシグナル (Signals-based Reactivity) [5.1.2]
*   **概念** : データの変更を極めて細粒度な「値そのもの」で追跡し、変更が発生した対象UIエレメントのみをピンポイントで高速に自動更新します [5.1.2, 138]。
*   **適用** : Vue 3 Composition API の `ref` / `reactive` / `computed` をシグナルとして活用。タイムラインで「いいね」を押した際、関係するカード（コンポーネント）の状態のみをミリ秒でピンポイント描画します [5.1.2, 140]。
*   **効果** : 画面全体の不必要な再描画や、仮想DOMの差分計算によるCPU負荷を完全に回避します [138, 140]。

### 4. Unixデバイスドライバ思想 (Unix Device Driver Paradigm) [9.1.2]
*   **概念** : ハードウェア（永続ストレージ）に低レベルI/Oを行うドライバと、それらを統治・オーケストレーションする制御層（カーネル/ミドルウェア）の責務を100%分離します [9.1.2]。
*   **適用** :
    *   **Driver層（Core Backend :5176）** : SQLite3 への GORM アクセスおよび Stashapp の UUID バインドに特化した「純粋なデバイスドライバ」 [9.1.2, 9.3]。
    *   **Middleware層（Go :5175）** : Pythonサイドカーの直接起動、進捗監視、アセット解決を担う「ファームウェア/カーネル」 [9.1.2]。
*   **効果** : データベース構造と表示用スキンが互いに汚染し合うのを物理的に防ぎます [9.1.2, 9.3]。

--------------------------------------------------------------------------------

## 11.2 厳格なコーディング規約 (Coding Standards & Constraints)

AIエージェントによるハルシネーション（暴走・デグレード）を100%防止し、人間開発者との円滑な協調を実現するため、以下の 5大規約をソースコードに対して強制適用します。

1.  **「1ファイル 100行以下」ルール (Strict 100-Line Limit)** [3.1]
    *   Go, TypeScript, Vue SFC, Python のすべてのファイルは、空行を含めて **100行以下** を絶対の上限制約とします [3.1]。
    *   これを超過しそうな場合、機械的に「CSSの Tailwind インライン化」「純粋関数の utils 逃がし」「Composable への状態分離」「Dumb UIコンポーネント分割」を徹底します [3.1.2]。
2.  **同階層3点セット原則 (Same Source, Same Flow / SSOT)** [3.3.1]
    *   起動・運用時は、ルート直下の **「実行バイナリ (dozo_katanuki.exe)」「実DB (archive.db)」「一元設定 (config.json)」** の3点セットのみを唯一のマスター（SSOT）とします [1.4, 3.3.1]。絶対パス依存やデータベースの散逸を禁止します [3.3.1]。
3.  **公開・閉塞IPバインドの完全分離 (IP Security Binding)** [3.3.2]
    *   外部（スマートフォン等）からの接続を受け付ける `:5173` (UI) および `:5175` (Middleware) は `0.0.0.0`（public）でバインド [3.3.2]。
    *   データベースCRUDを担う `:5176` (Core API) および `:9998` (Proxy)、`:9999` (Stashapp) は `127.0.0.1`（loopback）に強制バインドし、外部からの不正操作を物理的に遮断します [3.3.2]。
4.  **安全データパージ・ゴミ箱退避原則 (Safe Data Purging & File Removal)** [3.3.4]
    *   開発および運用時、不要ファイル・アセットの削除時は `rm` などの完全破壊的コマンドを禁止。OSのゴミ箱への移動を仲介するか、`.bak` などの一時退避を徹底します [3.3.4]。
5.  **start.bat 一元起動によるゾンビプロセス撲滅** [3.3.5]
    *   個別の手動起動を禁止。必ず `start.bat` をキックし、前回の残存ゾンビプロセスを強制キル（ゾンビキル）したのち、クリーンに各サービスを一括起動させます [3.3.5]。

--------------------------------------------------------------------------------

## 11.3 採用プログラミング言語 (Programming Languages)

本システムは、各レイヤーの特性と動作スタック（同期/常時稼働/非同期）の最適性を踏まえ、以下の3言語を適材適所で組み合わせて構築します。

| 言語 | 採用レイヤー | 動作特性 | 採用の技術的意義 |
| ------ | ------ | ------ | ------ |
| **Go 1.22+** | Middleware (:5175) <br>Core Backend (:5176) | 同期・常時稼働 <br>単一バイナリ | シングルバイナリパッキング、極めて高速な並行処理（Goroutine）、GORMによるDBMSの完全カプセル化、および低消費メモリ [1.2, 2.5]。 |
| **TypeScript / JS (ES6+)** | Frontend (:5173) <br>Skin Controllers | 宣言型UI <br>シグナルベース | 静的な型安全性の担保、Vue 3.5 との完全密着、およびプラガブル UI における `SkinController` インターフェースの整合契約 [1.2, 5.4]。 |
| **Python 3.10+** | Salvage Sidecar | オンデマンド起動 <br>非常駐スクリプト | 豊富なスクレイピング（CDX/WARC）ライブラリ、Stashインジェクション用 GraphQL API 処理の記述容易性、およびバッチバースト特性の最大化 [1.2, 2.4.2]。 |
| **YAML** | Skin Layouts (`layout.yaml`) | 設定・定義宣言 | JSON Schema (Draft 7) に完全準拠した、人間・AIエージェント双方に最も可読性の高いレイアウト設計図の宣言形式 [9.3]。 |
| **CSS3 (Tailwind)** | Skin Styles (`design.css`) | 宣言的装飾 | コンポーネントに個別の `<style>` を持たせないための、グローバルかつスコーププレフィックス付きスタイル定義 [3.1.2, 9.2]。 |

--------------------------------------------------------------------------------

## 11.4 採用モジュール・ライブラリ (Key Modules & Libraries)

「1ファイル100行以下」のルールを守りつつ、自作すると複雑化する処理を完全に外部委託するために、厳選された以下のコアライブラリをバンドル・ロードします。

### 1. Go (Middleware / Core API 側)
*   **GORM v2** (`gorm.io/gorm`, `gorm.io/driver/sqlite`) [2.5]
    *   生SQLの作成を100%禁止し、型安全なデータ構造（GORM Structs）から SQLite3 スキーマへの自動展開（AutoMigrate）を管理 [2.5, 7.1.1]。
*   **Go net/http & httputil**
    *   ミドルウェア内の「Stash Side Loader (:9998)」における、CORS許可ヘッダーをオーバーライド注入する透過プロキシハンドラの構築。

### 2. Frontend (Vue 3 側)
*   **hls.js** & **plyr** [5.7.1]
    *   Stashapp（:9999）がトランスコードする HLS ストリーム（m3u8）をブラウザ側CORSをバイパスしてゼロラグ再生する、最強のカスタムビデオプレイヤーコンポーネント。
*   **lucide-vue-next** [5.7.3]
    *   Vue の Tree Shaking 機能を100%効かせ、バンドルサイズを極限まで軽量化しつつ、SVGアイコン（Camera, Play, Loop 等）をインラインで Dumb 描画。
*   **lodash-es** (`get`) [5.7.4]
    *   `layout.yaml` に基づいて `RenderTree` からネストされたデータを安全に（ゼロポインタ・バグを回避して）動的走査・バインドするための必須ヘルパー。
*   **js-yaml** [5.7.4]
    *   ミドルウェアから配信された `layout.yaml` のテキストを、フロント側で瞬時に型安全な JSON オブジェクトへパース・デシリアライズする軽量ライブラリ。

### 3. Python (Salvage Sidecar 側)
*   **warcio** [2.2.3]
    *   Wayback Machineからのフェッチ時、HTTPリクエスト/レスポンスパケットを1ビットも改ざんせずに、一期一会の原本証明として圧縮コンテナ `.warc.gz` へストリーム同時キャプチャ保存。
*   **requests**
    *   Wayback CDX APIのスキャン、および3段階メディア確保における最初期の HTTP 直接 GET ストリーミングアタック。
*   **gql** (Python GraphQL Client) [2.4.2]
    *   破損チェック（ハッシュ整合）をクリアした本編メディア実ファイルを Stashapp メディアサーバーにインジェクションし、UUID（ stash_id ）を回収するための GraphQL 通信処理。

--------------------------------------------------------------------------------

## 11.5 外部連携アプリケーション (External Applications)

本システムは、データの重複排除や大容量トランスコード、マルチセッションの高速ダウンロードを自作して車輪の再発明（肥大化・バグ）を招くのを防ぐため、以下の2つの実績ある外部オープンソースアプリを完全に「物理接続エンジン（ヘッドレス）」として協調させます。

### 1. Stashapp Media Server (Port 9999 / Loopback 閉塞) [2.4, 3.3.2]
*   **役割** : 大容量の動画（H.264/H.265/VP9等）や静止画アセットの重複排除（OSHash / Perceptual Hash 照合） [2.4.1]、マルチビットレート HLS への自動トランスコード、および HLS ストリーミング配信 [2.4.1]。
*   **境界** : UIは一切露出させず、Python（Downloader）からの GraphQL API 経由でのインジェクションおよび、ミドルウェア（Stash Side Loader :9998）経由での透過プロキシ配信にのみ徹せさせます [2.4, 6.2]。

### 2. Motrix / Aria2 JSON-RPC (Port 6800 / Loopback 閉塞) [2.3, 3.3.2]
*   **役割** : 3段階メディア確保における「外部外注（OUTSOURCED）」の受け皿。コマンドライン型マルチセッション・高速ダウンローダ Aria2 を内蔵。
*   **境界** : 127.0.0.1:6800 の WebSocket/HTTP JSON-RPC を通じて `aria2.addUri` 委託タスクを発行。マルチスレッドでの爆速並列ロード、および途中で中断されたダウンロードの自動レジューム（再開）を外部に完全オフロードします [2.3.1]。

--------------------------------------------------------------------------------

## 11.6 外部連携Webサービス (Web Services)

1.  **Internet Archive Wayback Machine** [2.2]
    *   **CDX Server API** : 対象アカウントのすべての投稿ID（PostID）と過去のアーカイブスナップショット「タイムスタンプ」を  $O(N)$  で一括走査 [2.2.1]。
    *   **Memento Protocol (RFC 7089)** : 投稿日時に最も近い（時刻誤差が最小の）メディアアセットやHTMLキャッシュを特定するための時間折協（Time Negotiation）プロトコル [2.2.2]。
2.  **Google / DeepL 翻訳 API** [6.3.3]
    *   **適用** : インポート時の Python Mutator フェーズにおいて、1ポストごとに1.5秒以上のスリープ（`time.sleep`）と指数バックオフを伴う厳格な throttling（流量制御）を効かせながら「礼儀正しく」1回限りキック [10.2.1, 10.3.1]。
    *   **役割** : データベースの `articles` テーブルに `full_text_ja/en/zh` の 3大主要言語テキストを事前キャッシュ保存し、フロントエンドでの完全オフライン・ミリ秒言語トグル表示を成立させます [4.2, 5.6.3]。

--------------------------------------------------------------------------------

## 11.7 公式ドキュメント・リンク集 (Official Documentation)

> プロジェクトで採用した技術の公式ドキュメントを一覧にまとめています。
> バージョン固有の情報はリンク先で最新情報を確認してください。

### 🟦 プログラミング言語 (Languages)

| 技術 | バージョン | 公式リンク |
| ------ | ------ | ------ |
| **Go** | 1.22+ | [go.dev/doc](https://go.dev/doc/) |
| **Go — net/http パッケージ** | stdlib | [pkg.go.dev/net/http](https://pkg.go.dev/net/http) |
| **TypeScript** | 5.x | [typescriptlang.org/docs](https://www.typescriptlang.org/docs/) |
| **Python** | 3.10+ | [docs.python.org/3](https://docs.python.org/3/) |
| **SQLite3** | 3.x | [sqlite.org/docs](https://www.sqlite.org/docs.html) |

---

### 🟩 フロントエンド (Frontend)

| 技術 | バージョン | 公式リンク |
| ------ | ------ | ------ |
| **Vue 3** | 3.5.x | [vuejs.org/guide](https://vuejs.org/guide/introduction.html) |
| **Vue — Composition API** | 3.x | [vuejs.org/api/composition-api](https://vuejs.org/api/composition-api-setup.html) |
| **Vue Router** | 4.x | [router.vuejs.org](https://router.vuejs.org/) |
| **Vite** | 8.x | [vite.dev/guide](https://vite.dev/guide/) |
| **Tailwind CSS** | 3.x | [tailwindcss.com/docs](https://tailwindcss.com/docs/installation) |
| **hls.js** | latest | [github.com/video-dev/hls.js](https://github.com/video-dev/hls.js) |
| **Plyr** | latest | [plyr.io](https://plyr.io/) / [github.com/sampotts/plyr](https://github.com/sampotts/plyr) |
| **Lucide Vue Next** | latest | [lucide.dev/guide/packages/lucide-vue-next](https://lucide.dev/guide/packages/lucide-vue-next) |
| **Lodash-es** (`get`) | latest | [lodash.com/docs#get](https://lodash.com/docs/4.17.15#get) |
| **js-yaml** | latest | [github.com/nodeca/js-yaml](https://github.com/nodeca/js-yaml) |

---

### 🟧 バックエンド / Go ライブラリ (Backend / Go)

| 技術 | バージョン | 公式リンク |
| ------ | ------ | ------ |
| **GORM v2** | v2.x | [gorm.io/docs](https://gorm.io/docs/) |
| **GORM — SQLite ドライバ** | v2.x | [gorm.io/driver/sqlite](https://gorm.io/docs/connecting_to_the_database.html#SQLite) |
| **go-sqlite3** | latest | [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) |

---

### 🟨 Python サイドカー ライブラリ (Python Sidecar)

| 技術 | バージョン | 公式リンク |
| ------ | ------ | ------ |
| **warcio** | latest | [github.com/webrecorder/warcio](https://github.com/webrecorder/warcio) |
| **requests** | 2.x | [docs.python-requests.org](https://docs.python-requests.org/en/latest/) |
| **gql** (Python GraphQL Client) | 3.x | [gql.readthedocs.io](https://gql.readthedocs.io/en/latest/) |

---

### 📦 外部連携アプリケーション (External Applications)

| アプリ | 役割 | 公式リンク |
| ------ | ------ | ------ |
| **Stashapp** | メディアサーバー（HLS配信・重複排除） | [stashapp.cc](https://stashapp.cc/) / [github.com/stashapp/stash](https://github.com/stashapp/stash) |
| **Stashapp — GraphQL API** | インジェクション・UUID回収 | [github.com/stashapp/stash/blob/develop/graphql/schema](https://github.com/stashapp/stash/tree/develop/graphql/schema) |
| **Motrix** | GUI ダウンロードマネージャー（Aria2 内蔵） | [motrix.app](https://motrix.app/) / [github.com/agalwood/Motrix](https://github.com/agalwood/Motrix) |
| **Aria2** | マルチセッション高速ダウンローダー CLI | [aria2.github.io](https://aria2.github.io/) |
| **Aria2 — JSON-RPC** | RPC プロトコル仕様（:6800） | [aria2.github.io/manual/en/html/aria2c.html#rpc-interface](https://aria2.github.io/manual/en/html/aria2c.html#rpc-interface) |

---

### 🌐 外部 Web サービス / API (Web Services)

| サービス | 概要 | 公式リンク |
| ------ | ------ | ------ |
| **Internet Archive Wayback Machine** | 過去スナップショット検索・取得 | [archive.org/web](https://archive.org/web/) |
| **Wayback CDX Server API** | スナップショットインデックス高速走査 | [github.com/internetarchive/wayback/tree/master/wayback-cdx-server](https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server) |
| **Wayback CDX API ドキュメント** | クエリパラメータ詳細 | [web.archive.org/web/cdx/search](https://web.archive.org/cdx/search/cdx?url=example.com&output=json&limit=1) |
| **Memento Protocol (RFC 7089)** | 時間折協プロトコル仕様 | [RFC 7089 — IETF](https://datatracker.ietf.org/doc/html/rfc7089) / [timetravel.mementoweb.org](http://timetravel.mementoweb.org/) |
| **Memento Timemap API** | タイムマップ取得エンドポイント | [mementoweb.org/guide/api](http://mementoweb.org/guide/api/) |
| **DeepL API** | 機械翻訳（ja/en/zh 事前キャッシュ） | [developers.deepl.com/docs](https://developers.deepl.com/docs) |
| **DeepL API — Python SDK** | Python クライアントライブラリ | [github.com/DeepLcom/deepl-python](https://github.com/DeepLcom/deepl-python) |

---

### 📐 標準規格 / プロトコル (Standards & Protocols)

| 規格 | 概要 | 公式リンク |
| ------ | ------ | ------ |
| **WARC (ISO 28500)** | Webアーカイブコンテナ形式 | [iipc.github.io/warc-specifications](https://iipc.github.io/warc-specifications/) |
| **HLS — HTTP Live Streaming** | Apple HLS 仕様（RFC 8216） | [RFC 8216 — IETF](https://datatracker.ietf.org/doc/html/rfc8216) |
| **JSON Schema Draft 7** | config.json スキーマ検証規格 | [json-schema.org/draft-07](https://json-schema.org/specification-links.html#draft-7) |
| **GraphQL** | APIクエリ言語仕様 | [graphql.org/learn](https://graphql.org/learn/) |
| **SQLite WAL モード** | Write-Ahead Logging 仕様 | [sqlite.org/wal](https://www.sqlite.org/wal.html) |

---

**Navigation** : [← 前の章: 第10章：システム管理・設定・保守](10_admin_board_and_governance_v5) | [📚 目次 (Home)](Home) | [DocWiki ポータルへ戻る →](Home)
