# SNSエントリのストレージプール（SQLite3テーブル）のスキーマ設計書

## 概要

| 項目 | 内容 |
|------|------|
| **ドキュメントID** | STORAGE-001 |
| **バージョン** | 1.0 |
| **作成日** | 2026年8月16日 |
| **最終更新** | 2026年8月16日 |
| **作成者** | Kilo Assistant |
| **ステータス** | ドラフト（レビュー待ち） |
| **対象プロジェクト** | x_timeline_app |

本ドキュメントは、x_timeline_app における **SNS（X/Twitter等）エントリのストレージプールとしてのSQLite3テーブルスキーマ** を調査・設計・ドキュメント化したものです。対象とするSNSの種類、格納するデータ項目、テーブル制約、インデックス設計方針を明記し、永続的なデータ管理とMVVMアーキテクチャとの整合性を確保するためのスキーマ設計を提供します。

---

## 詳細内容

### 1. 対象とするSNSの種類

x_timeline_app では以下のSNSプラットフォームからのエントリを対象とします。

| SNSプラットフォーム | APIバージョン | 主な取得対象 | データ取得頻度 |
|---------------------|--------------|--------------|----------------|
| X (Twitter)         | API v2       | ツイート、リツイート、返信 | 5分間隔（ポーリング） |
| Mastodon            | API v1       | トゥート、ブースト | 15分間隔（ポーリング） |
| Bluesky             | AT Protocol  | ポスト、リポスト | 10分間隔（ポーリング） |

> ※現在はX（Twitter）のみを実装済み。拡張性を考慮したスキーマ設計を実施。

### 2. ストレージプールとしてのSQLite3テーブル構成

#### 2.1 データベースファイル
- **メインデータベース**: `chronoarchive.db`（日常アクセス用）
- **アーカイブデータベース**: `archive.db`（長期保存・バックアップ用）

#### 2.2 テーブル概略図
```
accounts       ← 1:N → tweets
               ↓
media          ← 1:N → tweets (via tweet_media)
               ↓
tweet_media    ← N:1 → tweets
               ↓
               N:1 → media
audit_log      ← （全テーブルの変更履歴）
```

#### 2.3 accounts テーブル（ユーザーアカウント情報）

| カラム名 | データ型 | 制約 | 説明 |
|----------|----------|------|------|
| id | TEXT | PRIMARY KEY | アカウントの一意識別子（例: X の numeric_id） |
| username | TEXT | NOT NULL, UNIQUE | アカウントのハンドル名（@username） |
| display_name | TEXT | NOT NULL | 表示名 |
| avatar_url | TEXT |  | プロフィール画像URL |
| header_url | TEXT |  | ヘッダー画像URL |
| description | TEXT |  | プロフィール説明 |
| location | TEXT |  | 場所情報 |
| followers_count | INTEGER | DEFAULT 0 | フォロワー数 |
| following_count | INTEGER | DEFAULT 0 | フォロー中数 |
| tweet_count | INTEGER | DEFAULT 0 | ツイート数 |
| listed_count | INTEGER | DEFAULT 0 | リスト登録数 |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | アカウント作成日時（SNS側） |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 最終更新日時 |
| is_protected | BOOLEAN | DEFAULT FALSE | 鍵アカウントフラグ |
| is_verified | BOOLEAN | DEFAULT FALSE | 認証バッジフラグ |
| raw_data | JSON |  | APIレスポンスの生データ（デバッグ用） |

**インデックス設計**:
- `idx_accounts_username` (username) - ユニーク
- `idx_accounts_created_at` (created_at) - 時系列検索用

#### 2.4 tweets テーブル（SNSエントリ本体）

| カラム名 | データ型 | 制約 | 説明 |
|----------|----------|------|------|
| id | TEXT | PRIMARY KEY | エントリの一意識別子（ツイートID等） |
| account_id | TEXT | NOT NULL, FK(accounts.id) | 投稿者アカウントID |
| conversation_id | TEXT |  | 会話ID（スレッドのルート） |
| in_reply_to_tweet_id | TEXT |  | 返信元ツイートID |
| in_reply_to_user_id | TEXT |  | 返信先ユーザーID |
| quoted_tweet_id | TEXT |  | 引用ツイートID |
| text | TEXT | NOT NULL | エントリ本文（ツイートテキスト） |
| lang | TEXT |  | 言語コード（ISO 639-1） |
| possibly_sensitive | BOOLEAN | DEFAULT FALSE | センシティブコンテンツフラグ |
| created_at | DATETIME | NOT NULL | 投稿日時（SNS側） |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 最終更新日時 |
| retweet_count | INTEGER | DEFAULT 0 | リツイート数 |
| reply_count | INTEGER | DEFAULT 0 | 返信数 |
| like_count | INTEGER | DEFAULT 0 | いいね数 |
| quote_count | INTEGER | DEFAULT 0 | 引用ツイート数 |
| bookmark_count | INTEGER | DEFAULT 0 | ブックマーク数 |
| impression_count | INTEGER |  | インプレッション数（X Premium） |
| source | TEXT |  | 投稿ソース（Twitter for iPhone等） |
| geo | JSON |  | 地理情報（GeoJSON形式） |
| entities | JSON |  | ハッシュタグ・メンション・URL等の構造化データ |
| withheld | JSON |  | 国別制限情報 |
| edit_history_tweet_ids | TEXT ARRAY |  | 編集履歴のツイートIDリスト（Xのみ） |
| version | INTEGER | NOT NULL, DEFAULT 1 | MVCC用バージョン番号 |
| status | TEXT | NOT NULL, DEFAULT 'active' | ステータス（active/deleted/archived/withheld） |
| wayback_url | TEXT |  | Wayback MachineアーカイブURL |
| raw_data | JSON |  | APIレスポンスの生データ（デバッグ用） |

**インデックス設計**:
- `idx_tweets_account_id` (account_id) - ユーザー別タイムライン取得用
- `idx_tweets_created_at` (created_at) - 時系列ソート用
- `idx_tweets_conversation_id` (conversation_id) - スレッド取得用
- `idx_tweets_status` (status) - ステータスフィルタ用
- `idx_tweets_lang` (lang) - 言語フィルタ用
- `idx_tweets_version` (version) - MVCC競合検出用
- `idx_tweets_account_created` (account_id, created_at) - 複合インデックス（ユーザー別時系列取得最適化）

#### 2.5 media テーブル（メディアファイル情報）

| カラム名 | データ型 | 制約 | 説明 |
|----------|----------|------|------|
| id | TEXT | PRIMARY KEY | メディアの一意識別子（StashのIDまたはハッシュ） |
| tweet_id | TEXT | NOT NULL, FK(tweets.id) | 関連ツイートID |
| media_type | TEXT | NOT NULL | メディアタイプ（photo/video/animated_gif） |
| url | TEXT | NOT NULL | ダウンロード元URL（Twitter CDN等） |
| local_path | TEXT |  | ローカルストレージパス |
| stash_id | TEXT |  | StashappでのID |
| width | INTEGER |  | 画像幅（ピクセル） |
| height | INTEGER |  | 画像高（ピクセル） |
| duration_ms | INTEGER |  | 動画の再生時間（ミリ秒、画像の場合NULL） |
| file_size_bytes | INTEGER |  | ファイルサイズ（バイト） |
| mime_type | TEXT |  | MIMEタイプ（image/jpeg等） |
| alt_text | TEXT |  | 代替テキスト（アクセシビリティ用） |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 取得日時 |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 最終更新日時 |
| download_status | TEXT | DEFAULT 'pending' | ダウンロード状態（pending/downloading/completed/failed） |
| blurhash | TEXT |  | プレースホルダー画像用BlurHash |
| dominant_color | TEXT |  | 優勢色（HEX形式） |
| is_sensitive | BOOLEAN | DEFAULT FALSE | センシティブメディアフラグ |
| raw_data | JSON |  | APIレスポンスの生データ（デバッグ用） |

**インデックス設計**:
- `idx_media_tweet_id` (tweet_id) - ツイート別メディア取得用
- `idx_media_media_type` (media_type) - メディアタイプフィルタ用
- `idx_media_download_status` (download_status) - ダウンロード状態フィルタ用
- `idx_media_created_at` (created_at) - 取得日時ソート用
- `idx_media_tweet_status` (tweet_id, download_status) - 複合インデックス

#### 2.6 audit_log テーブル（変更履歴・監査ログ）

| カラム名 | データ型 | 制約 | 説明 |
|----------|----------|------|------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | ログの一意識別子 |
| table_name | TEXT | NOT NULL | 対象テーブル名（accounts/tweets/media） |
| record_id | TEXT | NOT NULL | 対象レコードのID |
| operation | TEXT | NOT NULL | 操作種類（INSERT/UPDATE/DELETE） |
| changed_by | TEXT |  | 変更を行ったユーザーまたはシステム |
| changed_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 変更日時 |
| old_value | JSON |  | 変更前の値（NULL for INSERT） |
| new_value | JSON |  | 変更後の値（NULL for DELETE） |
| change_reason | TEXT |  | 変更理由（オプション） |

**インデックス設計**:
- `idx_audit_table_record` (table_name, record_id) - 特定レコードの履歴検索用
- `idx_audit_changed_at` (changed_at) - 時系列検索用
- `idx_audit_changed_by` (changed_by) - ユーザー別変更履歴検索用

### 3. テーブル間の関係性と制約

#### 3.1 外部キー制約（Foreign Key Constraints）
```sql
-- tweets テーブル
FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
FOREIGN KEY (in_reply_to_tweet_id) REFERENCES tweets(id) ON DELETE SET NULL
FOREIGN KEY (quoted_tweet_id) REFERENCES tweets(id) ON DELETE SET NULL

-- media テーブル
FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE
```

#### 3.2 トリガーによる自動更新（更新日時の自動設定）
```sql
-- accounts テーブルの updated_at 自動更新
CREATE TRIGGER update_accounts_timestamp 
AFTER UPDATE ON accounts
FOR EACH ROW
BEGIN
    UPDATE accounts SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- tweets テーブルの updated_at 自動更新
CREATE TRIGGER update_tweets_timestamp 
AFTER UPDATE ON tweets
FOR EACH ROW
BEGIN
    UPDATE tweets SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- media テーブルの updated_at 自動更新
CREATE TRIGGER update_media_timestamp 
AFTER UPDATE ON media
FOR EACH ROW
BEGIN
    UPDATE media SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

#### 3.3 MVCC用バージョンインクリメントトリガー
```sql
-- tweets テーブルのバージョン自動インクリメント
CREATE TRIGGER increment_tweets_version
BEFORE UPDATE ON tweets
FOR EACH ROW
WHEN NEW.version <= OLD.version
BEGIN
    SET NEW.version = OLD.version + 1;
END;
```

### 4. インデックス設計方針

#### 4.1 基本方針
1. **クエリパターンに基づく設計**: 実際のアプリケーションクエリを分析し、頻繁に使用されるWHERE/JOIN条件にインデックスを貼る
2. **複合インデックスの活用**: 並び替えとフィルタリングを同時に行うクエリに対して複合インデックスを設計
3. **カバーリングインデックス**: SELECT句のカラムがすべてインデックスに含まれるように設計し、テーブルアクセスを回避
4. **サイズ管理**: インデックスのサイズがテーブルサイズの20%を超えないように調整

#### 4.2 主要クエリパターンと対応インデックス

| クエリパターン | 目的 | 対応インデックス |
|----------------|------|------------------|
| `SELECT * FROM tweets WHERE account_id = ? ORDER BY created_at DESC LIMIT 20` | ユーザータイムライン取得 | `idx_tweets_account_created` |
| `SELECT * FROM tweets WHERE conversation_id = ? ORDER BY created_at ASC` | スレッド取得 | `idx_tweets_conversation_id` |
| `SELECT * FROM tweets WHERE created_at > ? AND created_at < ?` | 期間別検索 | `idx_tweets_created_at` |
| `SELECT * FROM tweets WHERE status = 'active' AND lang = 'ja'` | フィルタ付きタイムライン | 複合インデックス検討中 |
| `SELECT * FROM media WHERE tweet_id = ? AND download_status = 'completed'` | ツイートのメディア取得 | `idx_media_tweet_status` |
| `SELECT * FROM media WHERE media_type = 'video' AND download_status = 'completed'` | 動画メディア一覧 | `idx_media_media_type` + `idx_media_download_status` |
| `SELECT * FROM audit_log WHERE table_name = 'tweets' AND record_id = ? ORDER BY changed_at DESC` | 特定ツイートの変更履歴 | `idx_audit_table_record` |

#### 4.3 インデックスのメンテナンス戦略
- **定期的な再構築**: 月1回のスケジュールで `REINDEX` コマンド実行
- **使用統計の監視**: `PRAGMA index_info()` と `PRAGMA index_list()` で使用状況を監視
- **不要インデックスの削除**: 使用頻度が低いインデイスは定期的にレビューし削除

### 5. データ整合性とバリデーションルール

#### 5.1 CHECK制約によるデータバリデーション
```sql
-- accounts テーブル
CHECK (followers_count >= 0)
CHECK (following_count >= 0)
CHECK (tweet_count >= 0)
CHECK (listed_count >= 0)
CHECK (length(username) >= 1 AND length(username) <= 15)
CHECK (length(display_name) >= 1 AND length(display_name) <= 50)

-- tweets テーブル
CHECK (retweet_count >= 0)
CHECK (reply_count >= 0)
CHECK (like_count >= 0)
CHECK (quote_count >= 0)
CHECK (bookmark_count >= 0)
CHECK (length(text) <= 280)  -- Xの文字制限
CHECK (status IN ('active', 'deleted', 'archived', 'withheld'))
CHECK (version >= 1)

-- media テーブル
CHECK (width > 0 AND height > 0)
CHECK (duration_ms >= 0)
CHECK (file_size_bytes >= 0)
CHECK (media_type IN ('photo', 'video', 'animated_gif'))
CHECK (download_status IN ('pending', 'downloading', 'completed', 'failed'))
```

#### 5.2 トリガーによるビジネスルール enforcement
```sql
-- ツイート削除時のカスケード処理（メディアも論理削除に）
CREATE TRIGGER cascade_tweet_delete
AFTER UPDATE OF status ON tweets
WHEN NEW.status = 'deleted' AND OLD.status != 'deleted'
BEGIN
    UPDATE media SET download_status = 'failed' WHERE tweet_id = NEW.id;
END;
```

### 6. パフォーマンス最適化方針

#### 6.1 書き込みパフォーマンス向上
- **バッチインサート**: 複数レコードを1トランザクションで挿入
- **インデックス遅延作成**: 大量データ投入時は一時的にインデックスを無効化し、後で再作成
- **WALモードの活用**: デフォルトのWrite-Ahead Loggingモードで同時読み書きを可能に
- **同期モードの調整**: `PRAGMA synchronous = NORMAL` でパフォーマンスと安全性のバランスを取る

#### 6.2 読み込みパフォーマンス向上
- **キャッシュ戦略**: よくアクセスされるデータはアプリケーションレベルでキャッシュ
- **クエリ最適化**: EXPLAIN QUERY PLAN で実行計画を分析し最適化
- **部分インデックスの活用**: 特定条件にのみ適用するインデックス（例: `status = 'active'` のみにインデックス）

#### 6.3 ストレージ最適化
- **自動バキューム**: `PRAGMA auto_vacuum = INCREMENTAL` で断片化を防止
- **ページサイズ調整**: `PRAGMA page_size = 4096` でI/O効率を最適化
- **圧縮**: JSONカラムの定期的な圧縮検討（ただし可読性を考慮）

### 7. 移行とスキーマ進化戦略

#### 7.1 バージョン管理
- **スキーマバージョンテーブル**: `schema_migrations` テーブルで適用済みマイグレーションを管理
- **マイグレーションスクリプト**: 数値プレフィックス付きのSQLファイルで管理（例: `001_init.sql`, `002_add_wayback_url.sql`）

#### 7.2 後方互換性の確保
- **追加のみの変更**: 基本的にはカラムの追加のみを行い、削除や型変更は避ける
- **デフォルト値の設定**: 新規カラムには適切なデフォルト値を設定
- **NULL許容の設計**: 新規カラムは最初はNULL許容とし、後でデータを埋めてからNOT NULLに変更

#### 7.3 マイグレーション実装例
```sql
-- 003_add_wayback_url.sql
BEGIN TRANSACTION;

-- tweets テーブルに wayback_url カラムを追加
ALTER TABLE tweets ADD COLUMN wayback_url TEXT;

-- 既存レコードの初期値設定（NULLのままでも良いが、空文字列でも可）
UPDATE tweets SET wayback_url = NULL WHERE wayback_url IS NULL;

COMMIT;
```

### 8. セキュリティとプライバシー保護

#### 8.1 機密情報の取り扱い
- **プレーンテキストでの保存を禁止**: パスワード・トークン等は絶対にプレーンテキストで保存しない
- **環境変数による管理**: APIキー等の機密情報は環境変数（.envファイル）で管理
- **暗号化が必要なデータ**: 必要に応じて特定フィールドの暗号化を実装（例: ダイレクトメッセージ内容）

#### 8.2 アクセス制御
- **ローレベルセキュリティ（RLS）**: SQLiteではトリガーで擬似的に実装可能
- **アプリケーションレベルのチェック**: すべてのデータアクセスはアプリケーション層で認可チェックを実施
- **最小特権原则**: データベース接続ユーザーには必要最小限の権限のみを付与

#### 8.3 監査とログ
- **改ざん検出**: 重要なテーブルについてはハッシュチェーンによる改ざん検出を検討
- **アクセスログ**: データベースへのすべてのアクセスをログ記録（SELECT含む）
- **定期的な監査**: 月1回のスキーマとデータの整合性チェックを実施

### 9. 参考実装とベストプラクティス

#### 9.1 推奨されるSQLite3プラグマ設定
```sql
-- パフォーマンスと安全性のバランス
PRAGMA journal_mode = WAL;              -- Write-Ahead Logging
PRAGMA synchronous = NORMAL;            -- バランスの取れた同期モード
PRAGMA cache_size = -64000;             -- 64MBキャッシュ（負の値はKB単位）
PRAGMA foreign_keys = ON;               -- 外部キー制約を有効化
PRAGMA auto_vacuum = INCREMENTAL;       -- 自動バキューム（インクリメンタル）
PRAGMA mmap_size = 268435456;           -- 256MBメモリマップ I/O
PRAGMA busy_timeout = 5000;             -- 5秒のビジータイムアウト
```

#### 9.2 推奨されるバックアップ戦略
1. **ホットバックアップ**: SQLiteの `backup` API または `VACUUM INTO` を使用したオンラインバックアップ
2. **コールドバックアップ**: データベースファイルの直接コピー（必ず書き込みを停止してから）
3. **論理バックアップ**: `.dump` コマンドによるSQLダンプ（スキーマとデータの両方を取得）
4. **オフサイトバックアップ**: 暗号化したバックアップファイルをクラウドストレージにアップロード

#### 9.3 監視とメトリクス
- **パフォーマンスメトリクス**: クエリ実行時間、スロークエリの検出
- **ストレージメトリクス**: データベースファイルサイズ、インデックスサイズ、キャッシュヒット率
- **エラーメトリクス**: ロックタイムアウト、制約違反、構文エラーの発生頻度
- **カスタムメトリクス**: アプリケーション固有のビジネスメトリクス（アクティブユーザー数等）

---

## 参考資料

1. **SQLite3 Official Documentation** - https://www.sqlite.org/docs.html
2. **SQLite3 Query Language Reference** - https://www.sqlite.org/lang.html
3. **SQLite3 Pragmas** - https://www.sqlite.org/pragma.html
4. **Database Normalization** - C.J. Date, "Database Design and Relational Theory"
5. **Indexing Strategies** - Markus Winand, "SQL Performance Explained"
6. **タイムスタンプとバージョン管理** - Martin Fowler, "Patterns of Enterprise Application Architecture"
7. **監査ログ設計** - "Audit Logging: Best Practices for Implementing Audit Trails"

---

*このドキュメントはレビュー中のドラフトです。*  
*内容に関するご意見や修正指示がありましたら、お知らせください。*  
*このドキュメントはプロジェクトメンバーが参照しやすい構成になっています。*  
*スキーマ設計は将来の拡張性とパフォーマンスを考慮しています。*