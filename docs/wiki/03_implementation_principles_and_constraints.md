[[← 前の章: 第2章：外部サービスの概要とサルベージ技術|02_external_services_and_salvage]] | [[📚 目次 (Home)|Home]] | [[次の章: 第4章：データベース設計と仮想ストレージプール →|04_database_and_virtual_storage_pool]]

### 第3章：実装規約・制約原則（宣言型UI・UDF・AI駆動開発と開発者リファレンス）

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-PRINCIPLE-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（宣言型UI・UDF・Wails Bind完全準拠、AI暴走防止規約統合）

---

#### 3.1 「1ファイル 100行以下」の絶対ルール (Strict Rule)
AIエージェントおよび開発者がコードを生成・保守する際、コンポーネントおよびモジュールは **極限まで単一責任に細分化** し、1ファイルのソースコード is 空行を含めて「100行以下」を絶対的な制約とします。

##### 1. なぜ「100行以下」なのか？
*   **コンテキスト爆発の完全回避** : AIが作業する際、コードベース全体を走査させるとトークンが枯渇し、AIのハルシネーション（暴走）を引き起こします [91]。ファイルサイズを極小に保つことで、トークン消費を抑え、1回のやり取り（極小コンテキスト）で100%正確なコードを生成できます [36, 37]。
*   **テスト容易性とバグ混入率の低下** : モジュールが単一責任に閉じているため、ユニットテストが極めて容易になり、リファクタリング時に意図しないデグレードを100%防止できます。

##### 2. 100行を超過しそうな場合の対処フロー
ファイルが100行に達しそうな、あるいは超えてしまった場合、以下のステップに沿って機械的に分割を適用しなければなりません [10, 103]：
1.  **スタイル（CSS）の排除** : Tailwind CSS ユーティリティクラスを最大限活用するか、共通の `design.css` へスタイル記述を退避させます。
2.  **純粋計算・文字列操作の切り出し** : 日付フォーマットやテキストの改行処理といったビジネスロジックは、すべて `frontend/src/utils/` 配下に副作用のない純粋関数（Pure Function）として外出しします。
3.  **状態・Wails Bind通信の切り出し** : コンポーネント内の Wails Bind 呼び出しや状態管理（`ref`, `reactive`）は、すべて `frontend/src/composables/` の Composable（状態ホルダー）へと逃がします。
4.  **UIパーツのサブコンポーネント化（コンポーネント分割）** : `ArticleCard.vue` が肥大化した場合、 `ArticleAuthor.vue` (著者ヘッダー), `ArticleBody.vue` (本文), `ArticleStats.vue` (統計), `MediaGrid.vue` (メディア枠) に細分化して結合します。

---

#### 3.2 レイヤー別 責務境界（宣言型UI ＋ 単一データフロー UDF 原則）
Wailsキメラアプリの特性に適合するよう、システムアーキテクチャのレイヤー定義と「やって良いこと、絶対にやってはならないこと」を厳格に規定します [45]。

```mermaid
graph TD
    %% 単一データフロー（UDF）の明示
    Presentation[1. Presentation Layer<br>components/*.vue<br>Dumb Pure View]
    Composable[2. State & Signal Layer<br>composables/*.ts<br>UDF State Holder]
    Utility[3. Utility Layer<br>utils/*.ts<br>Pure Functions]
    WailsGoAPI[4. Wails Go API Layer<br>Go Bind (RPC / Controller)<br>RenderTree Factory]
    Driver[5. Driver Layer<br>GORM & SQLite3 CRUD<br>Process Manager]
    Admin[6. Admin & Governance Layer<br>Wails Lifeline sync<br>Disaster Recovery]

    %% データの流れ
    Presentation -- "1. ユーザー操作 (Event / Action)" --> Composable
    Composable -- "2. Wails Bind RPC 呼出 (インメモリ)" --> WailsGoAPI
    WailsGoAPI -- "3. GORM / SQL" --> Driver
    Driver -- "4. Raw DB Record" --> WailsGoAPI
    WailsGoAPI -- "5. Props: RenderTree (UDF)" --> Composable
    Composable -- "6. Reactive Signal" --> Presentation

    style Presentation fill:#e1f5fe,stroke:#03a9f4,stroke-width:2px
    style Composable fill:#e8f5e9,stroke:#4caf50,stroke-width:2px
    style Utility fill:#f3e5f5,stroke:#9c27b0,stroke-width:2px
    style WailsGoAPI fill:#ffe0b2,stroke:#ff9800,stroke-width:2px
    style Driver fill:#ffebee,stroke:#f44336,stroke-width:2px
    style Admin fill:#eceff1,stroke:#607d8b,stroke-width:2px
```

##### 1. Presentation層（components/*.vue）- Stateless Pure View（Dumb UI） [105]
*   **責務** : Composableから受領した Props（RenderTree または状態シグナル）に基づき、画面テンプレートで描画する [105]。ユーザー操作は単にイベントとして上位にエスカレーションする [107, 110]。
*   **禁止事項** : コンポーネント内部での独自Stateの保持、Wails関数の直接呼び出し、相対パスの組み立て（アセットパス生成等）、テキストパース処理。これらはすべて下位レイヤーが解決したPropsで受領しなければならない [109]。

##### 2. Composable層（composables/*.ts）- State & Signal Layer [106]
*   **責務** : `ref`, `reactive`, `computed` を使用したシグナルベースの細粒度リアクティブ状態ホルダー [110]。一方向データフロー（UDF）に則り、Wails Go API から一方向にフェッチしたデータを格納し、Presentation層へ Props として安全に流し込む [106]。
*   **禁止事項** : DOMの直接操作、HTMLマークアップやCSSスタイルの混入 [108]。

##### 3. Utility層（utils/*.ts）- Pure Utility Layer
*   **責務** : 同一の入力引数に対して常に全く同一の戻り値を返却し、いかなる外部状態も変更しない「数学的純粋関数」のみを配置する [110]。日付変換、文字列切り出し、キー名マッピング等。
*   **禁止事項** : 内部でのグローバル変数やLocalStorage変更等の副作用の発生。

##### 4. Wails Go API層（Go Bind）- RenderTree Factory
*   **責務** : フロントエンドのロジック肥大化を防ぐ「インテリジェント・ハブ」 [109]。SQLiteの生データ（Raw Model）をフロント側が描画するだけの完成されたデータ構造 **RenderTree**（仮想アバター解決、Wails AssetHandlerへの相対パスURL解決、翻訳、テキストリンク整形等がすべて完了したフラットなデータ）へ変換して配信する [109]。
*   **禁止事項** : 永続層（SQLite3）への直接SQL呼出。DB操作は必ずDriver層経由で実行しなければならない。

##### 5. Driver層（Go Driver & Process Controller）
*   **責務** : SQLite3へのGORMを用いた型安全なCRUD操作、ローカル Stashapp への GraphQL アクセスのカプセル化、および `stash.exe` のOSプロセス起動・制御・監視 [22]。
*   **禁止事項** : UIプレゼンテーション要素の関与。タイムラインの表示レイアウト情報や画面描画メタロジックに一切関与してはならない。

##### 6. Admin & Governance層（Wails Lifeline / Administration）
*   **責務** : `config.json` に基づく設定管理 [32]、DB整合性監査、Wails終了時の `taskkill` 完全実行、および二重化バックアップ（SQLiteスナップショット ＋ WARC/JSON原本ダンプ）の統制 [32, 51]。
*   **禁止事項** : 実行制御ロジックという呼称による誤用・ルーティング記述。WebハンドラーやAPI呼び出しロジックはここに記述してはならない。

---

#### 3.3 データフローとファイル配置の黄金律 (Same Source, Same Flow)
データベースの散逸、ゾンビプロセスによる不整合、およびデータの無秩序な汚染を防ぐための絶対的な規律です [53]。

##### 1. 同階層3点セット（Single Source of Truth: SSOT）の原則
開発時および運用時を問わず、システムのルートディレクトリ直下に存在する以下の **「3点セット」のみを唯一無二のマスター** とします [4, 11, 32]：
1.  **dozou_katanuki.exe**（実行バイナリ）
2.  **archive.db**（実稼働 SQLite3 データベース）
3.  **config.json**（システム一元設定ファイル）
DBパスを散逸させることを厳密に禁止します。常に同じ相対関係で動作させ、パス解決の破綻を根絶します [4]。

##### 2. スクリプト隔離原則
開発・テスト・検証・データパージ等、一時的あるいは管理者向けに作成したすべてのスクリプトは、**絶対にプロジェクトルートに放置してはなりません** [11, 103]。必ず `./scripts/` の配下に完全に隔離・分類して保存してください [11]。

##### 3. ファイル安全削除原則
不要ファイルを削除する際は `rm` や `os.remove` 等の完全な破壊的削除を行わず、必ずOSのゴミ箱への移動を仲介するか、 `.bak` サフィックスを付与して一時退避させることで、100%の復旧可能性を担保します [11]。

##### 4. Wailsデスクトップアプリ一元起動の徹底（個別手動起動の厳禁）
個別手動による `npm run dev` や `go run`、Pythonスクリプト等のバラバラな起動は、ポートの競合やゾンビプロセスを発生させるため厳重に禁じます [104]。
*   **ビルド** : 必ずプロジェクト提供 of `wails build` スケジューラ、またはルート直下の `build.bat` を通じてフロントアセット、Stashバイナリを一括パッキングしてビルドします [4, 104]。
*   **起動・動作確認** : ルート直下の `dozou_katanuki.exe` またはデバッグ用の `start_wails_dev.bat` をキックして起動します。この仕組みは、メモリ上に残存している前回の `stash.exe` 等のゾンビプロセスをOSレベルで自動検出し、キルした上でクリーンに一元起動させます [4, 104]。

---

[[← 前の章: 第2章：外部サービスの概要とサルベージ技術|02_external_services_and_salvage]] | [[📚 目次 (Home)|Home]] | [[次の章: 第4章：データベース設計と仮想ストレージプール →|04_database_and_virtual_storage_pool]]
