package models

import "time"

// DownloadReserve はキューから削除された、または安全退避されたタスクのメタデータです
type DownloadReserve struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GID         string    `gorm:"index;size:64" json:"gid"`
	URL         string    `gorm:"type:text;not null" json:"url"`
	FileName    string    `gorm:"size:255" json:"file_name"`
	ArticleID   string    `gorm:"index;size:64" json:"article_id"`
	MediaID     string    `gorm:"index;size:64" json:"media_id"`
	MirrorURLs  string    `gorm:"type:text" json:"mirror_urls"` // JSON string array of fallback URLs
	Status      string    `gorm:"index;size:32;default:'reserved'" json:"status"` // reserved, retrying, escalated, completed
	Reason      string    `gorm:"size:255" json:"reason"`
	TotalLength int64     `json:"total_length"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MirrorDownloadTask は1つの論理タスクに束ねられたミラー候補群です
type MirrorDownloadTask struct {
	PrimaryURL      string   `json:"primary_url"`
	MirrorURLs      []string `json:"mirror_urls"`
	FileName        string   `json:"file_name"`
	MediaID         string   `json:"media_id"`
	ArticleID       string   `json:"article_id"`
	DestinationPath string   `json:"destination_path"`
	Status          string   `json:"status"`
}
