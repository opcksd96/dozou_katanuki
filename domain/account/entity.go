package account

import "time"

// Account represents the pure domain entity for a user account
type Account struct {
	NumericID    string
	Username     string
	DisplayName  string
	AvatarURL    string
	AvatarBase64 string
	Description  string
	GroupName    string
	AliasOf      string
	IsWhitelist  bool
	IsTrash      bool
	TrashedBy    string
	TrashReason  string
	TrashedAt    *time.Time
	UpdatedAt    time.Time
	PostCount    int64
}

// MergeUpdates merges another account's non-empty fields into this account
func (a *Account) MergeUpdates(other *Account) {
	if other.Username != "" {
		a.Username = other.Username
	}
	if other.DisplayName != "" {
		a.DisplayName = other.DisplayName
	}
	if other.AvatarURL != "" {
		a.AvatarURL = other.AvatarURL
	}
	if other.Description != "" {
		a.Description = other.Description
	}
	a.UpdatedAt = time.Now()
}
