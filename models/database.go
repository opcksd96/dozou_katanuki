// models/database.go (100行以下)
package models

import "time"

// TableRecordResult はスプレッドシート表示用の汎用テーブルデータモデルです
type TableRecordResult struct {
	TableName string                   `json:"table_name"`
	Columns   []string                 `json:"columns"`
	Rows      []map[string]interface{} `json:"rows"`
	Total     int64                    `json:"total"`
}

// AccountDetailResult はアカウント管理ビュー用の詳細モデルです
type AccountDetailResult struct {
	Account   Account                 `json:"account"`
	Histories []AccountProfileHistory `json:"histories"`
	PostCount int64                   `json:"post_count"`
}

// MediaItemDetail はメディア管理ビュー用のレコードモデルです（RenderMediaを埋め込み統一）
type MediaItemDetail struct {
	RenderMedia
	MediaID     string    `json:"media_id"`
	ArticleID   string    `json:"article_id"`
	AccountID   string    `json:"account_id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	RawStatus   string    `json:"raw_status"`
	HasStash    bool      `json:"has_stash"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	FullText    string    `json:"full_text"`
	FullTextJA  string    `json:"full_text_ja,omitempty"`
	TweetDate   string    `json:"tweet_date"`
	WaybackURL  string    `json:"wayback_url,omitempty"`
}

// MediaScanItem はリポジトリ層からサービス層へ引き渡す中間スキャンモデルです
type MediaScanItem struct {
	Media          Media
	ArticleID      string
	AccountID      string
	Username       string
	DisplayName    string
	CreatedAt      time.Time
	FullText       string
	FullTextJA     string
	ProfileHistory []AccountProfileHistory
	WaybackURL     string
	AvatarBase64   string
}

// MediaSearchStats はメディア種別の統計カウントです
type MediaSearchStats struct {
	TotalCount int64 `json:"total_count"`
	ImageCount int64 `json:"image_count"`
	VideoCount int64 `json:"video_count"`
}

// MediaSearchResult はメディア一覧のページネーション結果です
type MediaSearchResult struct {
	Items []MediaItemDetail `json:"items"`
	Total int64             `json:"total"`
	Stats MediaSearchStats  `json:"stats"`
}

// DownloadStatusStats はダウンロードキュー全体のステータス別集計です
type DownloadStatusStats struct {
	Queued     int64 `json:"queued"`
	Completed  int64 `json:"completed"`
	Dead404    int64 `json:"dead_404"`
	Outsourced int64 `json:"outsourced"`
	Escalated  int64 `json:"escalated"`
	Retained   int64 `json:"retained"`
	Failed     int64 `json:"failed"`
	Total      int64 `json:"total"`
}
