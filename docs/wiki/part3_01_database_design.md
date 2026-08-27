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
GORMによる `AutoMigrate` [17] の実行時、このスキーマおよび各種制約（外部キー、NULL許容、ユニークインデックス）が完全自動で整合展開されます。

```sql
-- 1. accounts テーブル（ユーザープロフィールの基本原本データ）
CREATE TABLE IF NOT EXISTS accounts (
    numeric_id TEXT PRIMARY KEY,               -- SNS固有の不変な数値文字列ID (例: "1234567890123456789")
    username TEXT NOT NULL,                     -- 一意なスクリーンネーム / ハンドル (例: "msluo14", @マークなし)
    display_name TEXT NOT NULL,                -- 表示名。絵文字を含むマルチバイト文字列
    avatar_url TEXT NOT NULL,                  -- 本家SNSにおける生のアバターURL（基礎データ原本として100%不変保存）
    updated_at DATETIME NOT NULL               -- 最終同期・更新タイムスタンプ
);

-- 2. account_profile_history テーブル（0埋め3桁アバター世代管理 ＆ プロフィール変遷履歴）
CREATE TABLE IF NOT EXISTS account_profile_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,      -- 履歴シリアルID
    account_id TEXT NOT NULL,                  -- 外部キー: accounts(numeric_id)
    display_name TEXT NOT NULL,                -- その時点での表示名
    avatar_url TEXT NOT NULL,                  -- その時点でのアバターオリジナル生URL
    avatar_seq INTEGER NOT NULL,               -- アバター変更履歴を1からカウントアップする世代シリアル値 (1, 2, 3...)
    observed_at DATETIME NOT NULL,             -- 履歴観測・保存日時
    FOREIGN KEY (account_id) REFERENCES accounts(numeric_id) ON DELETE CASCADE
);

-- 3. articles テーブル（投稿メタデータ：スレッドツリーおよび会話構造の保持）
CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,                       -- 投稿ID (または記事ID)
    account_id TEXT NOT NULL,                  -- 外部キー: accounts(numeric_id)
    conversation_id TEXT NOT NULL,             -- 会話ツリーをグルーピングするためのスレッドルートID
    reply_to_status_id TEXT,                   -- 返信先（親）投稿ID (NULL許容)
    reply_to_username TEXT,                    -- 返信先（親）スクリーンネーム (NULL許容)
    created_at DATETIME NOT NULL,              -- 投稿作成タイムスタンプ
    full_text TEXT NOT NULL,                   -- 投稿本文テキスト
    via TEXT NOT NULL,                         -- 投稿ソース（クライアント名、例: "Twitter for Android"）
    is_retweet BOOLEAN NOT NULL DEFAULT 0,     -- リツイート（他者投稿の引用/転載）フラグ
    is_liked BOOLEAN NOT NULL DEFAULT 0,       -- お気に入り（ローカルブックマーク）フラグ
    wayback_url TEXT NOT NULL,                 -- Wayback Machine のキャッシュ原本URL
    FOREIGN KEY (account_id) REFERENCES accounts(numeric_id) ON DELETE CASCADE
);

-- 4. media テーブル（Stash IDとSQLite3をバインドする最重要リレーション層）
CREATE TABLE IF NOT EXISTS media (
    media_id TEXT PRIMARY KEY,                 -- オリジナルアセット取得URLの末尾（BaseName）をそのまま採用 [65]
    article_id TEXT NOT NULL,                  -- 外部キー: articles(id)
    type TEXT NOT NULL,                        -- メディア種別: "image" | "video" | "gif"
    download_url TEXT NOT NULL,                -- 本家CDNまたはWaybackのオリジナルメディアURL [52]
    width INTEGER NOT NULL,                    -- ピクセル横幅
    height INTEGER NOT NULL,                   -- ピクセル縦幅
    stash_scene_id TEXT UNIQUE,                -- Wailsが秘匿するStash内部プロセス側の動画UUID (NULL許容) [52]
    stash_image_id TEXT UNIQUE,                -- Wailsが秘匿するStash内部プロセス側の静止画UUID (NULL許容) [52]
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

-- 5. url_redirects テーブル（t.co 等のSNS短縮URL逆引きマップ）
CREATE TABLE IF NOT EXISTS url_redirects (
    short_url TEXT PRIMARY KEY,                -- 短縮URL (例: "https://t.co/eb7ymRi")
    expanded_url TEXT NOT NULL,                -- 解決済みのフルURL (例: "https://example.com/actual_destination")
    article_id TEXT NOT NULL,                  -- 外部キー: articles(id)
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

-- 6. whitelist テーブル（スパムパージ用・アーカイブ対象アカウントの統治ホワイトリスト）
CREATE TABLE IF NOT EXISTS whitelist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,                        -- 種別: "account" | "keyword"
    value TEXT NOT NULL UNIQUE,                -- 値: "msluo14" 等のスクリーンネーム
    is_active BOOLEAN NOT NULL DEFAULT 1       -- 有効フラグ
);
```

---

## 1.3 高速化インデックス定義 (10 Optimizations)
ローカルでの無限スクロール描画、会話スレッドの逆引き走査、およびStashappのUUIDとのバインド突合処理において、 **ミリ秒単位のゼロレイテンシレスポンス** を100%保証するため、以下の10大インデックスを厳格に適用します [16, 58]。

```sql
-- 1. 【お気に入りタイムライン高速化】
-- 目的: いいねした投稿（Bookmarks）のみを瞬時に逆順フィルタリングする [58]
CREATE INDEX IF NOT EXISTS idx_articles_is_liked_created ON articles(is_liked, created_at DESC) WHERE is_liked = 1;

-- 2. 【アカウント別タイムライン高速化】
-- 目的: 特定アカウントのタイムラインを無限スクロールで超高速に表示・ページネーションする [58]
CREATE INDEX IF NOT EXISTS idx_articles_account_created ON articles(account_id, created_at DESC);

-- 3. 【会話ツリー爆速スキャン】
-- 目的: 特定の conversation_id に属する関連投稿を時系列順に一括解決する
CREATE INDEX IF NOT EXISTS idx_articles_conversation ON articles(conversation_id, created_at ASC);

-- 4. 【リプライ親参照インデックス】
-- 目的: 特定の親記事に対する子リプライ（スレッドツリーの下流）を瞬時に検索する
CREATE INDEX IF NOT EXISTS idx_articles_reply_to ON articles(reply_to_status_id);

-- 5. 【統合タイムライン高速化】
-- 目的: 登録されている全アカウントの投稿を時系列に結合した統合タイムラインの描画ソートを O(1) 化する [58]
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);

-- 6. 【アバター履歴世代・逆引きミリ秒解決】
-- 目的: 記事の投稿日時に最も近いアバターの「世代キー」をGORM BeforeFindフックで瞬時に逆引き解決する [68]
CREATE INDEX IF NOT EXISTS idx_history_lookup ON account_profile_history(account_id, avatar_seq DESC);

-- 7. 【ユーザー名インクリメンタル検索】
-- 目的: ハンドルネーム（username）から numeric_id やプロフィール情報をミリ秒検索する
CREATE INDEX IF NOT EXISTS idx_accounts_username ON accounts(username);

-- 8. 【メディア紐付け（Reconciliation）自動高速化】
-- 目的: 添付メディアのロード、およびStash IDの逆引きインポートをバースト検索可能にする
CREATE INDEX IF NOT EXISTS idx_media_article ON media(article_id);

-- 9. 【StashビデオUUIDユニーク参照】
-- 目的: StashappのUUIDとSQLite3 mediaレコードの一対一整合性を最速で照合・監視する [52]
CREATE UNIQUE INDEX IF NOT EXISTS idx_media_stash_scene ON media(stash_scene_id) WHERE stash_scene_id IS NOT NULL;

-- 10. 【StashイメージUUIDユニーク参照】
-- 目的: Stashappの静止画UUIDとのバインド突合をミリ秒で判定し、二重インポートを100%防止する [52]
CREATE UNIQUE INDEX IF NOT EXISTS idx_media_stash_image ON media(stash_image_id) WHERE stash_image_id IS NOT NULL;
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
