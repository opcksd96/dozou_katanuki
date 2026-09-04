[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]

# 第3編 第1章：データベース設計と仮想ストレージプール (Database Design & Virtual Storage Pool)

**プロジェクト名** : dozou_katanuki (Wails-Stash Hybrid "土蔵・型抜き" Multi-Format Local Archival System)  
**ドキュメントID** : SPEC-DATABASE-001  
**バージョン** : 4.0.0 (Wailsキメラデスクトップ統合仕様)  
**作成日** : 2026-08-18  
**ステータス** : 正式仕様（GORMモデル完全同期・SQLite3 DDL・10大インデックス最適化・WALモード競合制御定義）

---

## 1.1 概要とデータ構造の整合性
本システムの中核永続化層は、極めて軽量かつ自己完結型なリレーショナルデータベースである **SQLite3** と、大容量バイナリの重複排除・トランスコード・HLS配信を担当し、Wails(Go)に内包・プロセス管理される **Stash Server** の2つの独立した仮想ストレージプールで構成されています [2, 14]。

この2つを繋ぐ唯一の架け橋（バインド関係）は、SQLite3 の `media` テーブルに保存される Stash 側の UUID（ `stash_scene_id` , `stash_image_id` ）のみです [15, 52]。これにより、データベース全損などの致命的障害時でも、100%ローカルにダンプされた二重化ソース（原本WARCおよびメタデータJSON）から、完全自動かつ非破壊的にリレーションを対称復元（ゼロリストア）できる「データの対称性（Symmetry）」を物理的に担保します [10.5]。

---

## 1.2 実稼働 SQLite3 スキーマ定義 (DDL)
実稼働マスターデータベース（ `archive.db` ）をクリーンビルド・再構築する際の、ANSI規格およびSQLite3方言に準拠した厳格な物理DDL（データ定義言語）仕様です [4, 15]。  
本システムは **pressly/goose** を用いた完全な `.sql` マイグレーション管理へ移行しており、GORMの `AutoMigrate` は使用しません。アプリケーション起動時に `migrations/` ディレクトリ内のSQLファイル群が自動的に走査・適用され、常に最新のスキーマへと整合進化します。

```sql
-- 1. accounts テーブル（ユーザープロフィールの基本原本データ）
CREATE TABLE IF NOT EXISTS accounts (
    numeric_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    display_name TEXT NOT NULL,
    avatar_url TEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    description TEXT,
    avatar_base64 TEXT,
    group_name TEXT DEFAULT "",
    alias_of TEXT DEFAULT "",
    is_whitelist NUMERIC DEFAULT 1,
    post_count INTEGER,
    is_trash BOOLEAN NOT NULL DEFAULT 0,
    trashed_by TEXT,
    trash_reason TEXT,
    trashed_at DATETIME
);

-- 2. account_profile_histories テーブル（アバター世代管理 ＆ プロフィール変遷履歴）
CREATE TABLE IF NOT EXISTS account_profile_histories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    avatar_original_url TEXT NOT NULL,
    avatar_seq INTEGER NOT NULL,
    avatar_virtual_key TEXT NOT NULL,
    observed_at DATETIME NOT NULL,
    avatar_base64 TEXT,
    description TEXT DEFAULT "",
    CONSTRAINT fk_accounts_profile_history FOREIGN KEY (account_id) REFERENCES accounts(numeric_id) ON DELETE CASCADE
);

-- 3. articles テーブル（投稿メタデータ：スレッドツリーおよび会話構造の保持）
CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    conversation_id TEXT NOT NULL,
    reply_to_id TEXT,
    reply_to_handle TEXT,
    created_at DATETIME NOT NULL,
    full_text TEXT NOT NULL,
    lang TEXT NOT NULL DEFAULT "ja",
    full_text_ja TEXT,
    full_text_en TEXT,
    full_text_zh TEXT,
    via TEXT NOT NULL,
    is_repost BOOLEAN NOT NULL DEFAULT 0,
    is_liked BOOLEAN NOT NULL DEFAULT 0,
    wayback_url TEXT NOT NULL,
    is_trash BOOLEAN NOT NULL DEFAULT 0,
    trashed_by TEXT,
    trash_reason TEXT,
    trashed_at DATETIME,
    source_domain TEXT,
    original_url TEXT,
    sotwe_url TEXT,
    nitter_url TEXT,
    twistalker_url TEXT,
    source_name TEXT,
    CONSTRAINT fk_accounts_articles FOREIGN KEY (account_id) REFERENCES accounts(numeric_id) ON DELETE CASCADE
);

-- 4. media テーブル（記事とアセットを結びつける中核リレーション層）
CREATE TABLE IF NOT EXISTS media (
    media_id TEXT PRIMARY KEY,
    article_id TEXT NOT NULL,
    type TEXT NOT NULL,
    download_url TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    download_status TEXT NOT NULL DEFAULT "QUEUED",
    failed_reason TEXT,
    stash_scene_id TEXT,
    stash_image_id TEXT,
    is_bookmarked BOOLEAN DEFAULT 0,
    media_quality TEXT DEFAULT "",
    tweet_urls TEXT DEFAULT '[]',
    thumbnail_url TEXT,
    is_trash BOOLEAN NOT NULL DEFAULT 0,
    trashed_by TEXT,
    trash_reason TEXT,
    trashed_at DATETIME,
    account_id TEXT,
    CONSTRAINT fk_articles_media FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

-- 5. media_variants テーブル（動画など、複数解像度のURLバリアントプール）
CREATE TABLE IF NOT EXISTS media_variants (
    variant_hash TEXT PRIMARY KEY,
    media_id TEXT NOT NULL,
    article_id TEXT NOT NULL,
    download_url TEXT NOT NULL,
    bit_rate INTEGER,
    CONSTRAINT fk_media_variants FOREIGN KEY (media_id) REFERENCES media(media_id) ON DELETE CASCADE
);

-- 6. url_redirects / whitelists テーブル
CREATE TABLE IF NOT EXISTS url_redirects (
    short_url TEXT PRIMARY KEY,
    expanded_url TEXT NOT NULL,
    article_id TEXT NOT NULL,
    CONSTRAINT fk_articles_url_redirects FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS whitelists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    value TEXT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    group_name TEXT DEFAULT "",
    alias_of TEXT DEFAULT ""
);

-- 7. ダウンロードタスク・予約キュー・Thunder制御用テーブル
CREATE TABLE IF NOT EXISTS download_reserves (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    g_id TEXT,
    url TEXT NOT NULL,
    file_name TEXT,
    article_id TEXT,
    media_id TEXT,
    mirror_urls TEXT,
    status TEXT DEFAULT "reserved",
    reason TEXT,
    total_length INTEGER,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS thunder_tasks (
    id TEXT PRIMARY KEY,
    media_id TEXT NOT NULL,
    article_id TEXT,
    resolution_type TEXT NOT NULL,
    url TEXT NOT NULL,
    file_name TEXT NOT NULL,
    status TEXT DEFAULT "PENDING",
    dispatched_at DATETIME,
    completed_at DATETIME,
    reaped_at DATETIME,
    created_at DATETIME,
    updated_at DATETIME,
    summary_size TEXT,
    error_reason TEXT,
    last_attempt_at DATETIME
);

CREATE TABLE IF NOT EXISTS download_tasks (
    media_id TEXT PRIMARY KEY,
    article_id TEXT,
    url TEXT NOT NULL,
    file_name TEXT NOT NULL,
    stage TEXT DEFAULT "REQUESTS",
    status TEXT DEFAULT "PENDING",
    failed_reason TEXT,
    requests_at DATETIME,
    motrix_at DATETIME,
    thunder_at DATETIME,
    stash_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME,
    updated_at DATETIME
);

-- 8. アカウント相関・Graph関連テーブル (PLAN-09対応)
CREATE TABLE IF NOT EXISTS account_relations (
    id TEXT PRIMARY KEY,
    source_account_id TEXT NOT NULL,
    target_account_id TEXT NOT NULL,
    target_handle TEXT,
    relation_type TEXT NOT NULL,
    direction TEXT NOT NULL,
    weight REAL NOT NULL DEFAULT 1.0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT fk_account_relations_source FOREIGN KEY (source_account_id) REFERENCES accounts(numeric_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS article_relation_evidences (
    id TEXT PRIMARY KEY,
    relation_id TEXT NOT NULL,
    source_article_id TEXT NOT NULL,
    target_article_id TEXT,
    evidence_type TEXT NOT NULL,
    context_snippet TEXT,
    media_id TEXT,
    is_salvaged BOOLEAN NOT NULL DEFAULT 0,
    observed_at DATETIME NOT NULL,
    CONSTRAINT fk_article_relation_evidences_rel FOREIGN KEY (relation_id) REFERENCES account_relations(id) ON DELETE CASCADE
);
```

---

## 1.3 高速化インデックス定義 (Optimizations)
ローカルでの無限スクロール描画、会話スレッドの逆引き走査、およびStashappのUUIDとのバインド突合処理において、 **ミリ秒単位のゼロレイテンシレスポンス** を100%保証するため、以下のインデックスを厳格に適用します [16, 58]。

```sql
-- 【タイムライン・記事検索 高速化】
CREATE INDEX IF NOT EXISTS idx_articles_is_liked_created ON articles(is_liked, created_at DESC) WHERE is_liked = 1;
CREATE INDEX IF NOT EXISTS idx_articles_account_created ON articles(account_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_conversation ON articles(conversation_id, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_articles_reply_to ON articles(reply_to_id);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_is_trash ON articles(is_trash);
CREATE INDEX IF NOT EXISTS idx_articles_trashed_by ON articles(trashed_by);

-- 【アカウント・プロフィール履歴 高速化】
CREATE INDEX IF NOT EXISTS idx_accounts_username ON accounts(username);
CREATE INDEX IF NOT EXISTS idx_history_lookup ON account_profile_histories(account_id, avatar_seq DESC);
CREATE INDEX IF NOT EXISTS idx_accounts_is_trash ON accounts(is_trash);

-- 【メディア紐付け（Reconciliation）自動高速化】
CREATE INDEX IF NOT EXISTS idx_media_article ON media(article_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_media_stash_scene ON media(stash_scene_id) WHERE stash_scene_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_media_stash_image ON media(stash_image_id) WHERE stash_image_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_media_status_type ON media(download_status, type);
CREATE INDEX IF NOT EXISTS idx_media_account ON media(account_id);
CREATE INDEX IF NOT EXISTS idx_media_is_trash ON media(is_trash);
CREATE INDEX IF NOT EXISTS idx_media_variants_media_id ON media_variants(media_id);
CREATE INDEX IF NOT EXISTS idx_media_variants_article_id ON media_variants(article_id);

-- 【ダウンロードキュー検索 高速化】
CREATE INDEX IF NOT EXISTS idx_download_reserves_status ON download_reserves(status);
CREATE INDEX IF NOT EXISTS idx_download_reserves_media_id ON download_reserves(media_id);
CREATE INDEX IF NOT EXISTS idx_thunder_tasks_status ON thunder_tasks(status);
CREATE INDEX IF NOT EXISTS idx_thunder_tasks_media_id ON thunder_tasks(media_id);
CREATE INDEX IF NOT EXISTS idx_download_tasks_status ON download_tasks(status);
CREATE INDEX IF NOT EXISTS idx_download_tasks_stage ON download_tasks(stage);

-- 【Graphリレーション検索 高速化】
CREATE INDEX IF NOT EXISTS idx_account_relations_source ON account_relations(source_account_id);
CREATE INDEX IF NOT EXISTS idx_account_relations_target ON account_relations(target_account_id);
CREATE INDEX IF NOT EXISTS idx_evidences_relation_id ON article_relation_evidences(relation_id);
```

---

## 1.4 WAL (Write-Ahead Logging) モードと同時実行制御仕様
本システムは、データベース接続の開始時に必ず **WAL（Write-Ahead Logging）モード** を有効化します [16]。これは、ローカルのマルチプロセス・マルチスレッド環境における「非ブロッキングUDFデータフロー」を成立させる上での絶対的な要件です [16, 43]。

##### 1. なぜ WAL モードなのか？（非ブロッキング読み込み）
*   **読み込み・書き込みの完全非衝突** [16]： Python非常駐サイドカー（Mutator）がバックグラウンドで何万行もの投稿・メディアデータをデータベースに一括インジェクションして排他ロックをかけている最中であっても、フロントエンド（Vue 3）からのタイムラインフェッチ（Go Bind呼び出し）は、 **1ミリ秒のブロッキング（遅延）もなく同時に実行可能** です [16]。
*   **R/W スループットの飛躍的向上**： 実データベースファイル（ `archive.db` ）を直接ロックして書き換えるのではなく、高速な追記専用のログファイル（ `archive.db-wal` ）を仲介するため、HDD/SSD等のI/O負荷を極限まで低減し、インポート処理の時間を 1/10 以下に圧縮します。

##### 2. Go（GORM）におけるWAL接続初期化コード規約
Wails バックエンドのデータベース接続プール初期化時、必ず以下の `PRAGMA` 設定を実行することを義務付けます。

```go
package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func InitDatabase(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to get database handle: %v", err)
	}

	// SQLite3 WALモードと整合性制約を強制有効化
	// busy_timeout = 5000 (ロック競合時の待機時間を5秒に設定し、エラー発生を防ぐ)
	pragmaQueries := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, query := range pragmaQueries {
		if _, err := sqlDB.Exec(query); err != nil {
			log.Fatalf("[FATAL] Failed to apply database pragma '%s': %v", query, err)
		}
	}

	log.Println("[INFO] SQLite3 Database connected in WAL mode with Foreign Keys enabled.")
	return db
}
```

---

## 1.5 仮想アバター ＆ Stash メディアプール物理分離規約のDB的意義
本システムにおける「仮想アバターリゾルバ」および「Stashとアバターの完全物理分離（Avatar Isolation Policy）」は、データベース構造的にも極めてクリーンな状態を保つために設計されました [15, 69]。

1.  **アバター情報の「オリジナル原本」と「表示用仮想キー」の同居** [15, 67]：
    *   `accounts.avatar_url` および `account_profile_history.avatar_url` には、将来の監査や原本証明のために、 `https://pbs.twimg.com/...` などの生URLを **100%不変の基礎データ** としてそのままDB保存します [15, 67]。
    *   フロントエンドに値を返す際、Wails Goバックエンドは自動で最新の世代履歴（ `avatar_seq` ）を結合し、表示用プロパティ `avatar_url` にのみ仮想アバターキー（例: `msluo14_avatar_001` ）をセットした RenderTree を構築して返却します [56, 68]。
2.  **Stashのデータベース（Scene / Image テーブル）のクリーン性保護** [15, 69]：
    *   StashにインジェクションされるメディアのUUIDは、すべて `media` テーブルの `stash_scene_id` 、 `stash_image_id` としてSQLite3にのみ登録され、アバター（プロフィール画像）がStash側に登録されることは1件もありません [15, 69]。
    *   これにより、Stash本来 of 画像ビューアに、アバター画像がゴミ・ノイズとして紛れ込んでUI表示が崩れるバグを、物理的（ディレクトリ構成）かつリレーショナルデータベースレベル（外部キー/NULL制限）で100%遮断します [15, 69]。

---

[[← 前の章: 第2編第2章：プラグインアーキテクチャとサイドカー|part2_02_plugin_architecture]] | [[📚 目次 (Home)|Home]] | [[次の章: 第3編第2章：フロントエンド層概論 →|part3_02_pure_dumb_frontend]]
