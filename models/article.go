package models

import (
	"database/sql"
	"time"
)

// Article represents articles table (Generic Timeline Item)
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

	// Relationships
	Account Account `gorm:"foreignKey:AccountID;references:NumericID" json:"account"`
	Media   []Media `gorm:"foreignKey:ArticleID;references:ID" json:"media"`
}
