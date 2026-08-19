package models

import "database/sql"

// Media represents media table (Stashapp Mapping & 3-Stage Recovery)
type Media struct {
	MediaID        string         `gorm:"primaryKey;column:media_id;type:text" json:"media_id"`
	ArticleID      string         `gorm:"index;column:article_id;type:text;not null" json:"article_id"`
	Type           string         `gorm:"column:type;type:text;not null" json:"type"`
	DownloadURL    string         `gorm:"column:download_url;type:text;not null" json:"download_url"`
	Width          int            `gorm:"column:width;type:integer;not null" json:"width"`
	Height         int            `gorm:"column:height;type:integer;not null" json:"height"`
	DownloadStatus string         `gorm:"column:download_status;type:text;not null;default:'QUEUED'" json:"download_status"`
	FailedReason   sql.NullString `gorm:"column:failed_reason;type:text" json:"failed_reason"`
	StashSceneID   sql.NullString `gorm:"uniqueIndex;column:stash_scene_id;type:text" json:"stash_scene_id"`
	StashImageID   sql.NullString `gorm:"uniqueIndex;column:stash_image_id;type:text" json:"stash_image_id"`
}
