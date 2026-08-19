// models/account.go (100行以下)
package models

import "time"

// Account represents accounts table (SSOT Profile Raw Data)
type Account struct {
	NumericID   string    `gorm:"primaryKey;column:numeric_id;type:text" json:"numeric_id"`
	Username    string    `gorm:"index;column:username;type:text;not null" json:"username"`
	DisplayName string    `gorm:"column:display_name;type:text;not null" json:"display_name"`
	AvatarURL   string    `gorm:"column:avatar_url;type:text;not null" json:"avatar_url"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`

	ProfileHistory []AccountProfileHistory `gorm:"foreignKey:AccountID;references:NumericID" json:"profile_history,omitempty"`
	Articles       []Article               `gorm:"foreignKey:AccountID;references:NumericID" json:"articles,omitempty"`
}
