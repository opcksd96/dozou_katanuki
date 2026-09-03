package sqlite

import "time"

// AccountSchema is the GORM-annotated struct used only for DB interactions
type AccountSchema struct {
	NumericID    string     `gorm:"primaryKey;column:numeric_id;type:text"`
	Username     string     `gorm:"index;column:username;type:text;not null"`
	DisplayName  string     `gorm:"column:display_name;type:text;not null"`
	AvatarURL    string     `gorm:"column:avatar_url;type:text;not null"`
	AvatarBase64 string     `gorm:"column:avatar_base64;type:text"`
	Description  string     `gorm:"column:description;type:text"`
	GroupName    string     `gorm:"column:group_name;type:text;default:''"`
	AliasOf      string     `gorm:"column:alias_of;type:text;default:''"`
	IsWhitelist  bool       `gorm:"column:is_whitelist;default:true"`
	IsTrash      bool       `gorm:"index;column:is_trash;type:boolean;not null;default:false"`
	TrashedBy    string     `gorm:"column:trashed_by;type:text"`
	TrashReason  string     `gorm:"column:trash_reason;type:text"`
	TrashedAt    *time.Time `gorm:"column:trashed_at;type:datetime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;not null"`
}

// TableName overrides the table name for GORM
func (AccountSchema) TableName() string {
	return "accounts"
}
