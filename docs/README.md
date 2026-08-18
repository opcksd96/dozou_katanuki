### 📚 dozou_katanuki DocWiki 総合ポータル
ようこそ、 **dozou_katanuki** 公式技術ドキュメント Wiki へ！ 本ドキュメント群は、Webの深淵から失われた投稿・メディアをサルベージし、ローカルで動態保存するためのアーキテクチャ・仕様・実装ガイドラインを網羅しています。

--------------------------------------------------------------------------------

#### 📑 全11章 技術概論＆完全仕様ナビゲーション
| 章 | タイトル | 概要・主要トピック |
| ------ | ------ | ------ |
| **第1章** | [第1章：技術仕様とバックボーン](01_technical_specs_and_backbone.md) | サルベージ動態保存の崇高な目的、技術スタック、ポートマップ、同階層3点セット |
| **第2章** | [第2章：外部サービスの概要とサルベージ技術](02_external_services_and_salvage.md) | Wayback CDX API、Aria2/Motrix P2P並列通信、Stashapp、公式リンク集 |
| **第3章** | [第3章：実装規約・制約原則（宣言型UI・UDF・AI駆動開発）](03_implementation_principles_and_constraints.md) | 「1ファイル100行以下」厳格ルール、宣言型UI＋単一データフロー（UDF）＋シグナル原則、スクリプト隔離原則、ファイル安全削除ルール |
| **第4章** | [第4章：データベース設計と仮想ストレージプール](04_database_and_virtual_storage_pool.md) | SQLite3 実稼働DDL完全版、10個の最適化インデックス、コンフリクト監査結果 |
| **第5章** | [第5章：プレゼンテーション層概論（Foolish Frontend & 宣言型UI）](05_foolish_frontend_and_declarative_ui_v4.md) | Dumb Component (Stateless Pure View)、単一データフロー (UDF)、シグナルの活用、Twitter互換URLルーティング、主要型定義（RenderTree） |
| **第6章** | [第6章：コンテンツディスパッチャー層（Middleware Hub & Proxy）](06_thick_middleware_and_proxy_v4.md) | Go Middleware Hub（生データのRenderTree完全構造化）、SPA Fallback（直打ち/F5保護）、完全相対パスプロキシ、動的ホスト対応 |
| **第7章** | [第7章：ドライバー層（Core Backend API）](07_robust_backend_driver_and_api_v4.md) | Go Core Backend + GORM、ArchiveDB 統一メディアURLインターフェース、ページネーション、Stash透過プロキシ |
| **第8章** | [第8章：ローカルストレージ保全とメディアポリシー](08_storage_persistence_and_media_policy_v2.md) | URL BaseName命名原則（URL末尾一致）、仮想アバターリゾルバ（3桁ナンバリング世代管理）、Stashメディアプールとの完全分離（ライブラリ汚染防止） |
| **第9章** | [第9章：プラグインアーキテクチャとサイドカー](09_plugin_architecture_and_sidecar_v3.md) | レンダリングプラグイン（Go）、Python非常駐サルベージサイドカーパイプライン、Stashインジェクション |
| **第10章** | [第10章：Admin Board（統合管理基盤と運用設計）](10_admin_board_and_governance_v5.md) | config.jsonによる設定一元化（Single Source of Truth）、DBメンテナンス、SQLiteスナップショット＆WARC/JSONダンプによる「二重化バックアップ（Disaster Recovery）」 |
| **第11章** | [第11章：参考資料・技術文献・型定義カタログ](11_references_and_literature.md) | RFC 7089 Memento、各ライブラリ公式ドキュメント、リファレンスリンク |

--------------------------------------------------------------------------------

#### 🗄️ アーカイブ・ドラフト
過去の設計草案や移行前レポートは [docs/draft/](draft/) に安全に保管されています。
