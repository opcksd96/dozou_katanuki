# 📚 dozou_katanuki DocWiki 総合ポータル

ようこそ、 **dozou_katanuki** 公式技術ドキュメント Wiki へ！ 本ドキュメント群は、Webの深淵から失われた投稿・メディアをサルベージし、ローカルで動態保存するためのアーキテクチャ・仕様・実装ガイドラインを網羅しています。

本ドキュメントは、システム全体の設計思想とレイヤーごとの責務を明確にするため、「4つの編（Part）」にカテゴライズされたモジュラー構造を採用しています。特定の章番号に依存しない「Single Source of Truth」としての情報管理を徹底しており、各トピックは独立したコンポーネントとして機能します。

---

## 🚩 【総合ステータス】 実装状況と今後のロードマップ
* **[🚀 実装状況と今後の実装計画書 (Status & Roadmap)](Status_and_Roadmap)**
  現在までの全5層（フロント/ミドル/ドライバ/DB/プラグイン）の実装達成度（95%達成）と、今後のマルチプラットフォーム（Bluesky/Instagram）拡張計画。

---

## 🛡️ 【Part 1】 アーキテクチャ・遵守事項 (Architecture & Principles)
システムが存在する意義と、AIおよび開発者が絶対に守るべき全体ルール・制約・ポリシーを定義します。

* **[技術仕様とバックボーン](part1_01_technical_specs)**
  サルベージ動態保存の崇高な目的、技術スタック、ポート非開放のキメラアーキテクチャの基本構造。
* **[外部サービスの概要とサルベージ技術](part1_02_external_services)**
  Wayback CDX API、Aria2/Motrix P2P並列通信、Stashapp、公式リンク集。
* **[実装規約・制約原則（宣言型UI・UDF・AI駆動開発）](part1_04_implementation_principles)**
  「1ファイル100行以下」厳格ルール、宣言型UI＋単一データフロー（UDF）＋シグナル原則、スクリプト隔離原則。
* **[ローカルストレージ保全とメディアポリシー](part1_03_storage_persistence)**
  URL BaseName命名原則、仮想アバターリゾルバ（3桁ナンバリング世代管理）、Stashメディアプールとの完全分離ポリシー。

---

## 📦 【Part 2】 Wails App Shell ＆ 運用基盤 (Shell & Governance)
システム全体を包み込む「外殻（Wails）」の起動・プロセスライフサイクルと、統合管理機構・拡張機構を定義します。

* **[Wails App Shell ＆ 統合管理基盤](part2_01_00_index)**
  `wails.json` を軸としたプロセスライフサイクル制御（`taskkill`道連れによるゾンビ化防止）と、統合設定ファイル `config.json` を操作するGUI基盤「Admin Board」の仕様。
* **[プラグインアーキテクチャとサイドカー](part2_02_plugin_architecture)**
  レンダリングプラグイン（Go）、Python非常駐サルベージサイドカーパイプライン、Stashインジェクション。

---

## ⚙️ 【Part 3】 内部システム設計・データモデル (Internal System & Data Flow)
Wailsの内部（内臓）で稼働する、UIからデータベースまでの具体的なデータフローと三層アーキテクチャ（フロントエンド・ミドルウェア・バックエンド）の設計を定義します。

* **[データベース設計と仮想ストレージプール](part3_01_database_design)**
  SQLite3 実稼働DDL完全版、最適化インデックス、コンフリクト監査結果、WALモード運用。
* **[フロントエンド層（Foolish Frontend & 宣言型UI）](part3_02_pure_dumb_frontend)**
  Dumb Component (Stateless Pure View)、Vue 3シグナル活用、アーカイブパスルーティング、主要型定義（RenderTree）。
* **[ミドルウェア層（Middleware Hub）](part3_03_0_middleware_index)**
  生データのRenderTree完全構造化、ページネーション制御、アバター仮想URLの解決と世代管理、多言語翻訳テキストのバインドとDOMリンク化。
* **[バックエンド層（Core Backend API & Driver）](part3_04_backend_driver)**
  Go Core Backend + GORM、ArchiveDB 統一メディアURLインターフェース、ローレベルクエリ処理。

---

## 📚 【Part 4】 リファレンス・付録 (Appendix)

* **[参考資料・技術文献・型定義カタログ](part4_01_references_and_literature)**
  RFC 7089 Memento、各ライブラリ公式ドキュメント、リファレンスリンク。

---