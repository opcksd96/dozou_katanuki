// models/media_excluded.go (100行以下 - SPEC-PRINCIPLE-001)
package models

import "time"

// MediaExcluded represents media_excluded table for quarantined/non-whitelist media
type MediaExcluded struct {
	MediaID          string    `gorm:"primaryKey;column:media_id;type:text" json:"media_id"`
	ArticleID        string    `gorm:"column:article_id;type:text;not null" json:"article_id"`
	AccountID        string    `gorm:"column:account_id;type:text" json:"account_id,omitempty"`
	Type             string    `gorm:"column:type;type:text;not null" json:"type"`
	DownloadURL      string    `gorm:"column:download_url;type:text;not null" json:"download_url"`
	Width            int       `gorm:"column:width;type:integer" json:"width"`
	Height           int       `gorm:"column:height;type:integer" json:"height"`
	DownloadStatus   string    `gorm:"column:download_status;type:text" json:"download_status"`
	FailedReason     string    `gorm:"column:failed_reason;type:text" json:"failed_reason,omitempty"`
	TweetURLs        string    `gorm:"column:tweet_urls;type:text" json:"tweet_urls,omitempty"`
	ThumbnailURL     string    `gorm:"column:thumbnail_url;type:text" json:"thumbnail_url,omitempty"`
	QuarantinedAt    time.Time `gorm:"column:quarantined_at;type:datetime" json:"quarantined_at"`
	QuarantineReason string    `gorm:"column:quarantine_reason;type:text" json:"quarantine_reason"`
}

// TableName overrides default table name for MediaExcluded
func (MediaExcluded) TableName() string {
	return "media_excluded"
}
