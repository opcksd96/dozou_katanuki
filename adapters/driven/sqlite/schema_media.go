// adapters/driven/sqlite/schema_media.go (100行以下 - SPEC-PRINCIPLE-001)
package sqlite

import (
	"database/sql"
	"dozou_katanuki/domain/entities"
	"time"
)

// MediaVariantSchema represents different resolutions/hashes for a single media item
type MediaVariantSchema struct {
	VariantHash string `gorm:"primaryKey;column:variant_hash;type:text"`
	MediaID     string `gorm:"index;column:media_id;type:text;not null"`
	ArticleID   string `gorm:"index;column:article_id;type:text;not null"`
	DownloadURL string `gorm:"column:download_url;type:text;not null"`
	BitRate     int    `gorm:"column:bit_rate;type:integer"`
}

// TableName overrides the table name for GORM
func (MediaVariantSchema) TableName() string {
	return "media_variants"
}

// MediaSchema represents media table (Stashapp Mapping and 3-Stage Recovery status)
type MediaSchema struct {
	MediaID        string               `gorm:"primaryKey;column:media_id;type:text"`
	ArticleID      string               `gorm:"index;column:article_id;type:text;not null"`
	AccountID      string               `gorm:"index;column:account_id;type:text"`
	Type           string               `gorm:"column:type;type:text;not null"`
	DownloadURL    string               `gorm:"column:download_url;type:text;not null"`
	Width          int                  `gorm:"column:width;type:integer;not null"`
	Height         int                  `gorm:"column:height;type:integer;not null"`
	DownloadStatus string               `gorm:"column:download_status;type:text;not null;default:'QUEUED'"`
	FailedReason   sql.NullString       `gorm:"column:failed_reason;type:text"`
	StashSceneID   sql.NullString       `gorm:"index;column:stash_scene_id;type:text"`
	StashImageID   sql.NullString       `gorm:"index;column:stash_image_id;type:text"`
	IsBookmarked   bool                 `gorm:"column:is_bookmarked;type:boolean;default:false"`
	MediaQuality   string               `gorm:"column:media_quality;type:text;default:''"`
	IsTrash        bool                 `gorm:"index;column:is_trash;type:boolean;not null;default:false"`
	TrashedBy      string               `gorm:"column:trashed_by;type:text"`
	TrashReason    string               `gorm:"column:trash_reason;type:text"`
	TrashedAt      *time.Time           `gorm:"column:trashed_at;type:datetime"`
	Variants       []MediaVariantSchema `gorm:"foreignKey:MediaID"`
}

// TableName overrides the table name for GORM
func (MediaSchema) TableName() string {
	return "media"
}

func (s *MediaSchema) ToEntity() *entities.Media {
	variants := make([]entities.MediaVariant, len(s.Variants))
	for i, v := range s.Variants {
		variants[i] = entities.MediaVariant{
			VariantHash: v.VariantHash, MediaID: v.MediaID, ArticleID: v.ArticleID,
			DownloadURL: v.DownloadURL, BitRate: v.BitRate,
		}
	}
	return &entities.Media{
		MediaID: s.MediaID, ArticleID: s.ArticleID, AccountID: s.AccountID,
		Type: s.Type, DownloadURL: s.DownloadURL, Width: s.Width, Height: s.Height,
		DownloadStatus: s.DownloadStatus, FailedReason: nullString(s.FailedReason),
		StashSceneID: nullString(s.StashSceneID), StashImageID: nullString(s.StashImageID),
		IsBookmarked: s.IsBookmarked, MediaQuality: s.MediaQuality,
		IsTrash: s.IsTrash, TrashedBy: s.TrashedBy, TrashReason: s.TrashReason, TrashedAt: s.TrashedAt,
		Variants: variants,
	}
}

func (s *MediaSchema) FromEntity(e *entities.Media) {
	s.MediaID = e.MediaID
	s.ArticleID = e.ArticleID
	s.AccountID = e.AccountID
	s.Type = e.Type
	s.DownloadURL = e.DownloadURL
	s.Width = e.Width
	s.Height = e.Height
	s.DownloadStatus = e.DownloadStatus
	s.FailedReason = toNullString(e.FailedReason)
	s.StashSceneID = toNullString(e.StashSceneID)
	s.StashImageID = toNullString(e.StashImageID)
	s.IsBookmarked = e.IsBookmarked
	s.MediaQuality = e.MediaQuality
	s.IsTrash = e.IsTrash
	s.TrashedBy = e.TrashedBy
	s.TrashReason = e.TrashReason
	s.TrashedAt = e.TrashedAt
	
	variants := make([]MediaVariantSchema, len(e.Variants))
	for i, v := range e.Variants {
		variants[i] = MediaVariantSchema{
			VariantHash: v.VariantHash, MediaID: v.MediaID, ArticleID: v.ArticleID,
			DownloadURL: v.DownloadURL, BitRate: v.BitRate,
		}
	}
	s.Variants = variants
}
