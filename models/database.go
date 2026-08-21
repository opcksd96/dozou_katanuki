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

// MediaItemDetail はメディア管理ビュー用のレコードモデルです
type MediaItemDetail struct {
	MediaID        string     `json:"media_id"`
	ArticleID      string     `json:"article_id"`
	AccountID      string     `json:"account_id"`
	Username       string     `json:"username"`
	Type           string     `json:"type"`
	DownloadURL    string     `json:"download_url"`
	Width          int        `json:"width"`
	Height         int        `json:"height"`
	DownloadStatus string     `json:"download_status"`
	FailedReason   *string    `json:"failed_reason"`
	StashSceneID   *string    `json:"stash_scene_id"`
	StashImageID   *string    `json:"stash_image_id"`
	CreatedAt      time.Time  `json:"created_at"`
}

// MediaSearchResult はメディア一覧のページネーション結果です
type MediaSearchResult struct {
	Items []MediaItemDetail `json:"items"`
	Total int64             `json:"total"`
}
