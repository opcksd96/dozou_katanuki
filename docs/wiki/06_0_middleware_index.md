# 第6章：ミドルウェア層（Middleware Hub）インデックス
**プロジェクト名** : dozou_katanuki (Pluggable UI & Multi-Format Local Archival System "土蔵・型抜き")  
**ドキュメントID** : SPEC-MIDDLEWARE-000  
**バージョン** : 4.0.0  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（Wails v2 キメラアーキテクチャ準拠・プロキシ層委譲完了・サブファイル細分化）  

**Navigation** : [← 前の章: 第5章：フロントエンド層](05_foolish_frontend_and_declarative_ui_v4) | [📚 目次 (Home)](Home) | [次の章: 第7章：バックエンド層 →](07_robust_backend_driver_and_api_v4)

---

## 概要：純化されたミドルウェアの真の責務

v4.0.0のアーキテクチャ大改修により、ミドルウェア層は劇的な進化を遂げました。
従来のポート開放（`:5175`, `:9998` など）に依存したリバースプロキシやSPAフォールバックルーティングの責務は、【Part 2】の Wails `AssetHandler` および管理層へと完全に委譲されています。

現在のミドルウェア層は、フロントエンド（Vue 3）とバックエンド（Core Driver）の間に立つ「純粋なデータ変換・統治ハブ」として、以下の3つの中核機能のみに100%専念します。AIによるハルシネーションを防ぐため、本章は機能ごとに以下の3つのサブファイルに厳密に細分化されています。

### 📦 1. [6.1 Middleware Core Components](06_1_middleware_core)
フロントエンドからのAPIリクエストを受け止め、バックエンドへのアクセスを整理・キューイングします。レスポンスを監査し、フロントエンドへとレンダリング情報を流す心臓部です。

### 🎨 2. [6.2 Data Decorator](06_2_data_decorator)
アバターの仮想解決、URLデコレーション、完全相対パス化、そしてフロントエンドに対する「プレースホルダー表示」の指示など、生データをUI向けに装飾する変換器です。

### 🔌 3. [6.3 Plugin Orchestrator](06_3_plugin_orchestrator)
レンダリングとスクレイピングを統合した「中核プラグイン」がインジェクトされ、最終的な表示形式（レイアウト・デザイン）が決定される拡張領域です。

---
**Navigation** : [← 前の章: 第5章：フロントエンド層](05_foolish_frontend_and_declarative_ui_v4) | [📚 目次 (Home)](Home) | [次の章: 第7章：バックエンド層 →](07_robust_backend_driver_and_api_v4)