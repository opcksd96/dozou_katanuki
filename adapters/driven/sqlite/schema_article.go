// adapters/driven/sqlite/schema_article.go (100行以下 - SPEC-PRINCIPLE-001)
package sqlite

import (
	"database/sql"
	"dozou_katanuki/domain/entities"
	"time"
)

// ArticleSchema is the GORM-annotated struct used only for DB interactions
type ArticleSchema struct {
	ID             string         `gorm:"primaryKey;column:id;type:text"`
	AccountID      string         `gorm:"index;column:account_id;type:text;not null"`
	ConversationID string         `gorm:"index;column:conversation_id;type:text;not null"`
	ReplyToID      sql.NullString `gorm:"column:reply_to_id;type:text"`
	ReplyToHandle  sql.NullString `gorm:"column:reply_to_handle;type:text"`
	CreatedAt      time.Time      `gorm:"index;column:created_at;type:datetime;not null"`
	FullText       string         `gorm:"column:full_text;type:text;not null"`
	Lang           string         `gorm:"column:lang;type:text;not null;default:'ja'"`
	FullTextJA     sql.NullString `gorm:"column:full_text_ja;type:text"`
	FullTextEN     sql.NullString `gorm:"column:full_text_en;type:text"`
	FullTextZH     sql.NullString `gorm:"column:full_text_zh;type:text"`
	Via            string         `gorm:"column:via;type:text;not null"`
	IsRepost       bool           `gorm:"column:is_repost;type:boolean;not null;default:false"`
	IsLiked        bool           `gorm:"index;column:is_liked;type:boolean;not null;default:false"`
	WaybackURL     string         `gorm:"column:wayback_url;type:text;not null"`
	SourceDomain   sql.NullString `gorm:"column:source_domain;type:text"`
	OriginalURL    sql.NullString `gorm:"column:original_url;type:text"`
	SotweURL       sql.NullString `gorm:"column:sotwe_url;type:text"`
	NitterURL      sql.NullString `gorm:"column:nitter_url;type:text"`
	TwistalkerURL  sql.NullString `gorm:"column:twistalker_url;type:text"`
	SourceName     sql.NullString `gorm:"column:source_name;type:text"`
	IsTrash        bool           `gorm:"index;column:is_trash;type:boolean;not null;default:false"`
	TrashedBy      sql.NullString `gorm:"index;column:trashed_by;type:text"`
	TrashReason    sql.NullString `gorm:"column:trash_reason;type:text"`
	TrashedAt      sql.NullTime   `gorm:"column:trashed_at;type:datetime"`
}

// TableName overrides the table name for GORM
func (ArticleSchema) TableName() string {
	return "articles"
}

// nullString returns *string or nil
func nullString(ns sql.NullString) *string {
	if ns.Valid {
		v := ns.String
		return &v
	}
	return nil
}

// toNullString converts *string to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// ToEntity converts GORM schema to pure Domain Entity
func (s *ArticleSchema) ToEntity() *entities.Article {
	return &entities.Article{
		ID:             s.ID,
		AccountID:      s.AccountID,
		ConversationID: s.ConversationID,
		ReplyToID:      nullString(s.ReplyToID),
		ReplyToHandle:  nullString(s.ReplyToHandle),
		CreatedAt:      s.CreatedAt,
		FullText:       s.FullText,
		Lang:           s.Lang,
		FullTextJA:     nullString(s.FullTextJA),
		FullTextEN:     nullString(s.FullTextEN),
		FullTextZH:     nullString(s.FullTextZH),
		Via:            s.Via,
		IsRepost:       s.IsRepost,
		IsLiked:        s.IsLiked,
		WaybackURL:     s.WaybackURL,
		SourceDomain:   nullString(s.SourceDomain),
		OriginalURL:    nullString(s.OriginalURL),
		SotweURL:       nullString(s.SotweURL),
		NitterURL:      nullString(s.NitterURL),
		TwistalkerURL:  nullString(s.TwistalkerURL),
		SourceName:     nullString(s.SourceName),
		IsTrash:        s.IsTrash,
		TrashedBy:      nullString(s.TrashedBy),
		TrashReason:    nullString(s.TrashReason),
		TrashedAt:      func() *time.Time { if s.TrashedAt.Valid { t := s.TrashedAt.Time; return &t }; return nil }(),
	}
}

// FromEntity populates GORM schema from pure Domain Entity
func (s *ArticleSchema) FromEntity(e *entities.Article) {
	s.ID = e.ID
	s.AccountID = e.AccountID
	s.ConversationID = e.ConversationID
	s.ReplyToID = toNullString(e.ReplyToID)
	s.ReplyToHandle = toNullString(e.ReplyToHandle)
	s.CreatedAt = e.CreatedAt
	s.FullText = e.FullText
	s.Lang = e.Lang
	s.FullTextJA = toNullString(e.FullTextJA)
	s.FullTextEN = toNullString(e.FullTextEN)
	s.FullTextZH = toNullString(e.FullTextZH)
	s.Via = e.Via
	s.IsRepost = e.IsRepost
	s.IsLiked = e.IsLiked
	s.WaybackURL = e.WaybackURL
	s.SourceDomain = toNullString(e.SourceDomain)
	s.OriginalURL = toNullString(e.OriginalURL)
	s.SotweURL = toNullString(e.SotweURL)
	s.NitterURL = toNullString(e.NitterURL)
	s.TwistalkerURL = toNullString(e.TwistalkerURL)
	s.SourceName = toNullString(e.SourceName)
	s.IsTrash = e.IsTrash
	s.TrashedBy = toNullString(e.TrashedBy)
	s.TrashReason = toNullString(e.TrashReason)
	s.TrashedAt = func() sql.NullTime { if e.TrashedAt != nil { return sql.NullTime{Time: *e.TrashedAt, Valid: true} }; return sql.NullTime{} }()
}
