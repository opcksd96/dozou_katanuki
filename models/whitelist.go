// models/whitelist.go (100行以下)
package models

type Whitelist struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Type      string `gorm:"column:type;type:text;not null" json:"type"`
	Value     string `gorm:"uniqueIndex;column:value;type:text;not null" json:"value"`
	GroupName string `gorm:"column:group_name;type:text;default:''" json:"group_name"`
	AliasOf   string `gorm:"column:alias_of;type:text;default:''" json:"alias_of"`
	IsActive  bool   `gorm:"column:is_active;type:boolean;not null;default:true" json:"is_active"`
}

