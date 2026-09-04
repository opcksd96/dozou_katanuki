// models/account.go (100行以下 - SPEC-PRINCIPLE-001)
package models

import "time"

// Account represents accounts table (SSOT Profile Raw Data)
type Account struct {
	NumericID    string     `gorm:"primaryKey;column:numeric_id;type:text" json:"numeric_id"`
	Username     string     `gorm:"index;column:username;type:text;not null" json:"username"`
	DisplayName  string     `gorm:"column:display_name;type:text;not null" json:"display_name"`
	AvatarURL    string     `gorm:"column:avatar_url;type:text;not null" json:"avatar_url"`
	AvatarBase64 string     `gorm:"column:avatar_base64;type:text" json:"avatar_base64"`
	Description  string     `gorm:"column:description;type:text" json:"description"`
	GroupName    string     `gorm:"column:group_name;type:text;default:''" json:"group_name"`
	AliasOf      string     `gorm:"column:alias_of;type:text;default:''" json:"alias_of"`
	IsWhitelist  bool       `gorm:"column:is_whitelist;default:true" json:"is_whitelist"`
	IsTrash      bool       `gorm:"index;column:is_trash;type:boolean;not null;default:false" json:"is_trash"`
	TrashedBy    string     `gorm:"column:trashed_by;type:text" json:"trashed_by,omitempty"`
	TrashReason  string     `gorm:"column:trash_reason;type:text" json:"trash_reason,omitempty"`
	TrashedAt    *time.Time `gorm:"column:trashed_at;type:datetime" json:"trashed_at,omitempty"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
	PostCount    int64      `gorm:"->" json:"post_count"`

	ProfileHistory []AccountProfileHistory `gorm:"foreignKey:AccountID;references:NumericID;constraint:OnDelete:CASCADE;" json:"profile_history,omitempty"`
	Articles       []Article               `gorm:"foreignKey:AccountID;references:NumericID;constraint:OnDelete:CASCADE;" json:"articles,omitempty"`
}
