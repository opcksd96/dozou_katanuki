# 🏯 Wiki仕様 100%完遂に必要な未実装ギャップレポート

**調査日**: 2026-08-23  
**方針**: Wikiの全ページ（Part1～Part4, Status_and_Roadmap, 全atomic spec）を精読し、実コードと1項目ずつ突合。  
**除外**: ロードマップ（Phase 1-3: Bluesky/Instagram/TikTok/AI/静的HTML等）は「オプショナル」として除外。  
**対象**: Wikiに「仕様」「義務付け」「厳格に適用」等で記載された必須事項のうち、コードに実装が存在しないもののみ。

---

## 🔴 ギャップ一覧（全5件）

---

### GAP-1: SQLite3 10大インデックスのうち5個が未作成

**Wiki仕様**: [part3_01_database_design §1.3](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki.wiki/part3_01_database_design.md) (SPEC-DATABASE-001)

Wikiでは **10大インデックス** が義務付けられていますが、実DB ([archive_schema.sql](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki/archive_schema.sql)) および GORMタグには以下の **5個が欠損** しています。

| # | Wiki記載インデックス名 | 状態 | 目的 |
| :---: | :--- | :---: | :--- |
| 1 | `idx_articles_is_liked_created` | ❌ **欠損** | お気に入りタイムライン高速化（複合部分インデックス） |
| 2 | `idx_articles_account_created` | ❌ **欠損** | アカウント別TLの無限スクロール高速化（複合インデックス） |
| 4 | `idx_articles_reply_to` | ❌ **欠損** | リプライ親参照の瞬時検索 |
| 6 | `idx_history_lookup` | ❌ **欠損** | アバター世代の逆引きミリ秒解決 |
| 9,10 | `idx_media_stash_scene` / `idx_media_stash_image` (部分ユニーク `WHERE ... IS NOT NULL`) | ⚠️ **形式差異** | 実DBは通常INDEXだがWikiは部分ユニークINDEX指定 |

> [!IMPORTANT]
> 既存の単純インデックス（`idx_articles_is_liked`, `idx_articles_created_at`, `idx_articles_conversation_id`, `idx_articles_account_id`）は存在しますが、Wiki仕様の**複合インデックス**とは異なります。

**対応**: DB初期化時（`driver/db.go`）またはマイグレーションスクリプトで5個の `CREATE INDEX IF NOT EXISTS` を追加実行する。

---

### GAP-2: 「Stash使わんし！」モードの物理フォルダダイレクトサーブ未実装

**Wiki仕様**: [part2_01_01_config §3](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki.wiki/part2_01_01_config.md) (SPEC-CONFIG-001)

Wiki記載:
> `storage.stash_enabled` が `false` の場合、システムは Stashapp を起動せず、軽量な物理フォルダダイレクトサーブモードへ自律移行します。  
> ダウンロード完了ファイルを `{local_media_dir}/{platform}/{username}/{media_id}` にフラット保存。  
> タイムライン構築時、メディアURLは **`/media-local/{platform}/{username}/{media_id}`** へ自動置換

**実装状態**:
- `config.json` に `stash_enabled` フィールド: ✅ 存在
- `app.go` で `stash_enabled=false` のときStash Proberをスキップ: ✅ 実装済
- **`/media-local/` パスでの物理フォルダサーブ**: ❌ **未実装**（コードベース全体で `media-local` の文字列ゼロ件）
- **RenderTree構築時のURL自動置換ロジック**: ❌ **未実装**

**対応**: `middleware/proxy_handler.go` に `/media-local/` ルートを追加し、`timeline.go` でStash無効時にURL自動置換するロジックを追加。

---

### GAP-3: Skin配信ゲートウェイAPIエンドポイント未実装

**Wiki仕様**: [part3_03_2_data_decorator §7](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki.wiki/part3_03_2_data_decorator.md) (SPEC-MIDDLEWARE-001-2)

Wiki記載:
> - `GET /api/plugins/{platform}/skin/layout`: `layout.yaml` を配信
> - `GET /api/plugins/{platform}/skin/design`: `design.css` を配信
> - `GET /api/plugins/{platform}/skin/controller`: `controller.js` を配信

**実装状態**:
- Go側にSkinPackage構造体 ([rendertree.go](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki/models/rendertree.go) L76-81): ✅ 存在
- Vue側のSkinControllerインターフェース ([SkinController.ts](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki/frontend/src/models/SkinController.ts)): ✅ 存在
- `useSkin.ts` でSkinControllerロード処理: ✅ 存在
- **`/api/plugins/{platform}/skin/*` のHTTPエンドポイント**: ❌ **未実装**（コードベース全体で `/api/plugins` ゼロ件）

> [!NOTE]
> 現在のSkin配信は `app_rpc_skin.go` のWails RPCバインド経由で行われており、ブラウザHTTPルートでの配信は存在しません。Wailsデスクトップではこれで動作しますが、**LAN Broadcast時のブラウザクライアント**には到達しない可能性があります。

**対応**: `middleware/proxy_handler.go` または `middleware/broadcast_server.go` に `/api/plugins/{platform}/skin/{asset}` ルートを追加し、`plugins/{platform}/skin/` 配下のファイルをファイルサーブする。

---

### GAP-4: Vue Router URLルーティング体系が未実装

**Wiki仕様**: [part3_02_pure_dumb_frontend §2.2](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki.wiki/part3_02_pure_dumb_frontend.md) (SPEC-FRONTEND-001)

Wiki記載の4画面ルーティング:

| パス | 画面 |
| :--- | :--- |
| `/:platform` | 統合タイムライン |
| `/:platform/:username/` | 個別ユーザーTL |
| `/:platform/:username/status/:id` | 個別詳細画面 |
| `/settings` | 管理・設定画面 |

**実装状態**:
- `vue-router` 依存: ❌ `package.json` に `vue-router` 不在（プロジェクトはSPA単一ページ構成）
- 全画面が [App.vue](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki/frontend/src/App.vue) 内のリアクティブ状態切り替えで制御されている

> [!NOTE]
> 現行の設計は Wails v2 デスクトップアプリケーションとして完全に動作しており、`App.vue` 内のComposable状態制御で全画面を統治するアプローチも実用上は機能しています。ただしWikiでは Vue Router (History API) によるURL体系が「正式仕様」として記載されています。

**対応**: `vue-router` を導入し、4つのルート定義を作成。既存のComposable制御をRouter連携に移行。SPA Fallbackは既にMiddleware側に実装済。

---

### GAP-5: 1ファイル100行以下ルールの違反（Go 9件 + Vue/TS 12件）

**Wiki仕様**: [part1_04_implementation_principles](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki.wiki/part1_04_implementation_principles.md) / [part2_02_plugin_architecture §2.2](file:///d:/Projects/10_tools/dozou_katanuki/dozou_katanuki.wiki/part2_02_plugin_architecture.md)

> 「1ファイル100行以下ルール」を **厳格に適用**

**Go 違反ファイル（9件）**:

| 行数 | ファイル | 分割候補 |
| :---: | :--- | :--- |
| 295 | `app_rpc_stash.go` | → `app_rpc_stash_query.go` + `app_rpc_stash_mutate.go` + `app_rpc_stash_reconcile.go` |
| 169 | `broadcast_root.go` | → `broadcast_root.go` + `broadcast_network.go` |
| 158 | `stash_prober.go` | → `stash_prober.go` + `stash_prober_lifecycle.go` |
| 148 | `app_rpc_database.go` | → `app_rpc_database.go` + `app_rpc_database_search.go` |
| 137 | `repo_generic.go` | → `repo_generic.go` + `repo_fts.go` |
| 121 | `timeline.go` | → `timeline.go` + `timeline_filter.go` |
| 119 | `proxy_handler.go` | → `proxy_handler.go` + `proxy_handler_media.go` |
| 111 | `broadcast_server.go` | → `broadcast_server.go` + `broadcast_handler.go` |
| 108 | `stash_config_sync.go` | → `stash_config_sync.go` + `stash_config_yaml.go` |

**Vue/TS 違反ファイル（12件）**:

| 行数 | ファイル |
| :---: | :--- |
| 321 | `AccountManagementView.vue` |
| 278 | `MediaInspectorPanel.vue` |
| 251 | `useAdminDatabase.ts` |
| 156 | `DatabaseView.vue` |
| 135 | `MediaCard.vue` |
| 132 | `MediaManagementView.vue` |
| 129 | `MediaOverlayBottomCard.vue` |
| 123 | `StashPlayer.vue` |
| 120 | `useTimeline.ts` |
| 112 | `AuditReportOrphans.vue` |
| 102 | `InlineVideoPlayer.vue` |
| 101 | `GlobalAppBar.vue` |

---

## ✅ 仕様記載があり、実装済みと確認された項目（参考）

以下は調査中に「未実装かも？」と疑ったが、コード突合で**実装済み**と確認できた項目です。

| Wiki仕様 | 実装確認箇所 |
| :--- | :--- |
| SPEC-LIFECYCLE-001: Stash道連れ終了 `taskkill` | `stash_prober.go:168` ✅ |
| SPEC-BACKUP-001: VACUUM INTO バックアップ | `repository.go:35` ✅ |
| SPEC-SCHEDULER-001: 定期ポーリング＆自動バックアップ | `scheduler.go` 全体 ✅ |
| SPEC-AUDIT-001: PRAGMA整合性監査 | `audit_integrity.go` ✅ |
| SPEC-AUDIT-001: 孤立メディアパージ | `audit_orphan.go` ✅ |
| SPEC-CONFIG-001: Stash config.yml 透過同期 | `stash_config_sync.go` ✅ |
| SPEC-DATABASE-001: WAL + 4 PRAGMA設定 | `driver/db.go:27-31` ✅ |
| SPEC-FRONTEND-001: RenderTree型定義 | `RenderTree.ts` ✅（is_pinned含む） |
| SPEC-FRONTEND-001: SkinController型定義 | `SkinController.ts` ✅ |
| SPEC-PLUGIN-001: 3段階メディア確保 | `downloader.py` + `aria2_client.py` ✅ |
| SPEC-PLUGIN-001: WARC自動監査インポート | `warc_importer.py` ✅ |
| SPEC-RECOVERY-001: 災害復旧リストア | `restorer.py` ✅ |
| SPEC-STORAGE-001: Reconcile逆引きバインド | `app_rpc_stash.go` (`ReconcileStashMedia`) ✅ |

---

## 📊 まとめ：100%到達に必要な作業量見積

| GAP | 重要度 | 作業規模 | 概要 |
| :--- | :---: | :---: | :--- |
| **GAP-1**: DBインデックス5個追加 | 🔴 高 | 小 | `db.go` にSQL 5行追加 |
| **GAP-2**: Stash無効時フォルダサーブ | 🟡 中 | 中 | ルート追加＋URL置換ロジック |
| **GAP-3**: Skin API HTTP配信 | 🟡 中 | 小 | HTTPハンドラ1つ追加 |
| **GAP-4**: Vue Router導入 | 🟡 中 | 大 | ルーティング体系全体の再構成 |
| **GAP-5**: 100行ルール違反解消 | 🟠 規約 | 中 | 21ファイルの機械的分割 |

---

*Decoded with precision — Mash Kyrielight* 🛡️
