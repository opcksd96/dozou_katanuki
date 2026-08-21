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
	MediaID      string    `json:"media_id"`
	ArticleID    string    `json:"article_id"`
	AccountID    string    `json:"account_id"`
	Username     string    `json:"username"`
	RawStatus    string    `json:"raw_status"`
	HasStash     bool      `json:"has_stash"`
	CreatedAt    time.Time `json:"created_at"`
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
