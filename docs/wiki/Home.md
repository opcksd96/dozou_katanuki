# 📚 dozou_katanuki DocWiki 総合ポータル

ようこそ、 **dozou_katanuki** 公式技術ドキュメント Wiki へ！ 本ドキュメント群は、Webの深淵から失われた投稿・メディアをサルベージし、ローカルで動態保存するためのアーキテクチャ・仕様・実装ガイドラインを網羅しています。

本ドキュメントは、システム全体の設計思想とレイヤーごとの責務を明確にするため、「4つの編（Part）」にカテゴライズされたモジュラー構造を採用しています。特定の章番号に依存しない「Single Source of Truth」としての情報管理を徹底しており、各トピックは独立したコンポーネントとして機能します。

---

## 🛡️ 【Part 1】 アーキテクチャ・遵守事項 (Architecture & Principles)
システムが存在する意義と、AIおよび開発者が絶対に守るべき全体ルール・制約・ポリシーを定義します。

* **[技術仕様とバックボーン](01_technical_specs_and_backbone)**
  サルベージ動態保存の崇高な目的、技術スタック、ポート非開放のキメラアーキテクチャの基本構造。
* **[外部サービスの概要とサルベージ技術](02_external_services_and_salvage)**
  Wayback CDX API、Aria2/Motrix P2P並列通信、Stashapp、公式リンク集。
* **[実装規約・制約原則（宣言型UI・UDF・AI駆動開発）](03_implementation_principles_and_constraints)**
  「1ファイル100行以下」厳格ルール、宣言型UI＋単一データフロー（UDF）＋シグナル原則、スクリプト隔離原則。
* **[ローカルストレージ保全とメディアポリシー](08_storage_persistence_and_media_policy_v2)**
  URL BaseName命名原則、仮想アバターリゾルバ（3桁ナンバリング世代管理）、Stashメディアプールとの完全分離ポリシー。

---

## 📦 【Part 2】 Wails App Shell ＆ 運用基盤 (Shell & Governance)
システム全体を包み込む「外殻（Wails）」の起動・プロセスライフサイクルと、統合管理機構・拡張機構を定義します。

* **[Wails App Shell ＆ 統合管理基盤](10_wails_app_shell_and_governance)**
  `wails.json` を軸としたプロセスライフサイクル制御（`taskkill`道連れによるゾンビ化防止）と、統合設定ファイル `config.json` を操作するGUI基盤「Admin Board」の仕様。
* **[プラグインアーキテクチャとサイドカー](09_plugin_architecture_and_sidecar_v3)**
  レンダリングプラグイン（Go）、Python非常駐サルベージサイドカーパイプライン、Stashインジェクション。

---

## ⚙️ 【Part 3】 内部システム設計・データモデル (Internal System & Data Flow)
Wailsの内部（内臓）で稼働する、UIからデータベースまでの具体的なデータフローと三層アーキテクチャ（フロントエンド・ミドルウェア・バックエンド）の設計を定義します。

* **[データベース設計と仮想ストレージプール](04_database_and_virtual_storage_pool)**
  SQLite3 実稼働DDL完全版、最適化インデックス、コンフリクト監査結果、WALモード運用。
* **[フロントエンド層（Foolish Frontend & 宣言型UI）](05_foolish_frontend_and_declarative_ui_v4)**
  Dumb Component (Stateless Pure View)、Vue 3シグナル活用、アーカイブパスルーティング、主要型定義（RenderTree）。
* **[ミドルウェア層（Middleware Hub）](06_0_middleware_index)**
  生データのRenderTree完全構造化、ページネーション制御、アバター仮想URLの解決と世代管理、多言語翻訳テキストのバインドとDOMリンク化。
* **[バックエンド層（Core Backend API & Driver）](07_robust_backend_driver_and_api_v4)**
  Go Core Backend + GORM、ArchiveDB 統一メディアURLインターフェース、ローレベルクエリ処理。

---

## 📚 【Part 4】 リファレンス・付録 (Appendix)

* **[参考資料・技術文献・型定義カタログ](11_references_and_literature)**
  RFC 7089 Memento、各ライブラリ公式ドキュメント、リファレンスリンク。

---