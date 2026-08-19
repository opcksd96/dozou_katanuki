package models
package models

import "time"

// Account represents accounts table (SSOT Profile Raw Data)
type Account struct {
	NumericID   string    `gorm:"primaryKey;column:numeric_id;type:text" json:"numeric_id"`
	Username    string    `gorm:"index;column:username;type:text;not null" json:"username"`
	DisplayName string    `gorm:"column:display_name;type:text;not null" json:"display_name"`
	AvatarURL   string    `gorm:"column:avatar_url;type:text;not null" json:"avatar_url"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`

	// Relationships
	ProfileHistory []AccountProfileHistory `gorm:"foreignKey:AccountID;references:NumericID" json:"profile_history,omitempty"`
	Articles       []Article               `gorm:"foreignKey:AccountID;references:NumericID" json:"articles,omitempty"`
}

// AccountProfileHistory represents account_profile_history table
type AccountProfileHistory struct {
	ID                uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	AccountID         string    `gorm:"index;column:account_id;type:text;not null" json:"account_id"`
	DisplayName       string    `gorm:"column:display_name;type:text;not null" json:"display_name"`
	AvatarOriginalURL string    `gorm:"column:avatar_original_url;type:text;not null" json:"avatar_original_url"`
	AvatarSeq         int       `gorm:"column:avatar_seq;type:integer;not null" json:"avatar_seq"`
	AvatarVirtualKey  string    `gorm:"column:avatar_virtual_key;type:text;not null" json:"avatar_virtual_key"`
	ObservedAt        time.Time `gorm:"column:observed_at;type:datetime;not null" json:"observed_at"`
}
