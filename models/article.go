// models/article.go (100行以下)
package models

import (
	"database/sql"
	"time"
)

// Article represents articles table (Generic Timeline Item Specification)
type Article struct {
	ID             string         `gorm:"primaryKey;column:id;type:text" json:"id"`
	AccountID      string         `gorm:"index;column:account_id;type:text;not null" json:"account_id"`
	ConversationID string         `gorm:"index;column:conversation_id;type:text;not null" json:"conversation_id"`
	ReplyToID      sql.NullString `gorm:"column:reply_to_id;type:text" json:"reply_to_id"`
	ReplyToHandle  sql.NullString `gorm:"column:reply_to_handle;type:text" json:"reply_to_handle"`
	CreatedAt      time.Time      `gorm:"index;column:created_at;type:datetime;not null" json:"created_at"`
	FullText       string         `gorm:"column:full_text;type:text;not null" json:"full_text"`
	Lang           string         `gorm:"column:lang;type:text;not null;default:'ja'" json:"lang"`
	FullTextJA     sql.NullString `gorm:"column:full_text_ja;type:text" json:"full_text_ja"`
	FullTextEN     sql.NullString `gorm:"column:full_text_en;type:text" json:"full_text_en"`
	FullTextZH     sql.NullString `gorm:"column:full_text_zh;type:text" json:"full_text_zh"`
	Via            string         `gorm:"column:via;type:text;not null" json:"via"`
	IsRepost       bool           `gorm:"column:is_repost;type:boolean;not null;default:false" json:"is_repost"`
	IsLiked        bool           `gorm:"index;column:is_liked;type:boolean;not null;default:false" json:"is_liked"`
	WaybackURL     string         `gorm:"column:wayback_url;type:text;not null" json:"wayback_url"`
	SourceDomain   sql.NullString `gorm:"column:source_domain;type:text" json:"source_domain,omitempty"`
	OriginalURL    sql.NullString `gorm:"column:original_url;type:text" json:"original_url,omitempty"`
	SotweURL       sql.NullString `gorm:"column:sotwe_url;type:text" json:"sotwe_url,omitempty"`
	NitterURL      sql.NullString `gorm:"column:nitter_url;type:text" json:"nitter_url,omitempty"`
	TwistalkerURL  sql.NullString `gorm:"column:twistalker_url;type:text" json:"twistalker_url,omitempty"`
	SourceName     sql.NullString `gorm:"column:source_name;type:text" json:"source_name,omitempty"`
	IsTrash        bool           `gorm:"index;column:is_trash;type:boolean;not null;default:false" json:"is_trash"`
	TrashedBy      sql.NullString `gorm:"index;column:trashed_by;type:text" json:"trashed_by,omitempty"`
	TrashReason    sql.NullString `gorm:"column:trash_reason;type:text" json:"trash_reason,omitempty"`
	TrashedAt      sql.NullTime   `gorm:"column:trashed_at;type:datetime" json:"trashed_at,omitempty"`

	Account      Account       `gorm:"foreignKey:AccountID;references:NumericID" json:"account,omitempty"`
	Media        []Media       `gorm:"foreignKey:ArticleID;references:ID" json:"media,omitempty"`
	UrlRedirects []UrlRedirect `gorm:"foreignKey:ArticleID;references:ID" json:"url_redirects,omitempty"`
}
