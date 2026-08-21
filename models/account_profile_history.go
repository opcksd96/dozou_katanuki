// models/account_profile_history.go (100行以下)
package models

import "time"

// AccountProfileHistory represents account_profile_history table
type AccountProfileHistory struct {
	ID                uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	AccountID         string    `gorm:"index;column:account_id;type:text;not null" json:"account_id"`
	DisplayName       string    `gorm:"column:display_name;type:text;not null" json:"display_name"`
	AvatarOriginalURL string    `gorm:"column:avatar_original_url;type:text;not null" json:"avatar_original_url"`
	AvatarSeq         int       `gorm:"column:avatar_seq;type:integer;not null" json:"avatar_seq"`
	AvatarVirtualKey  string    `gorm:"column:avatar_virtual_key;type:text;not null" json:"avatar_virtual_key"`
	AvatarBase64      string    `gorm:"column:avatar_base64;type:text" json:"avatar_base64"`
	ObservedAt        time.Time `gorm:"column:observed_at;type:datetime;not null" json:"observed_at"`
}
