// models/account.go (100行以下)
package models

import "time"

// Account represents accounts table (SSOT Profile Raw Data)
type Account struct {
	NumericID    string    `gorm:"primaryKey;column:numeric_id;type:text" json:"numeric_id"`
	Username     string    `gorm:"index;column:username;type:text;not null" json:"username"`
	DisplayName  string    `gorm:"column:display_name;type:text;not null" json:"display_name"`
	AvatarURL    string    `gorm:"column:avatar_url;type:text;not null" json:"avatar_url"`
	AvatarBase64 string    `gorm:"column:avatar_base64;type:text" json:"avatar_base64"`
	Description  string    `gorm:"column:description;type:text" json:"description"`
	GroupName    string    `gorm:"column:group_name;type:text;default:''" json:"group_name"`
	AliasOf      string    `gorm:"column:alias_of;type:text;default:''" json:"alias_of"`
	IsWhitelist  bool      `gorm:"column:is_whitelist;default:true" json:"is_whitelist"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`

	ProfileHistory []AccountProfileHistory `gorm:"foreignKey:AccountID;references:NumericID" json:"profile_history,omitempty"`
	Articles       []Article               `gorm:"foreignKey:AccountID;references:NumericID" json:"articles,omitempty"`
}
