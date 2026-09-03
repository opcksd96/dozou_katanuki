// domain/entities/account.go (100行以下 - SPEC-PRINCIPLE-001)
package entities

import "time"

// Account は純粋なドメインエンティティです（インフラ・DB・UIへの依存なし）
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
}

// IsValid はアカウントエンティティが有効な状態か判定します
func (a *Account) IsValid() bool {
	if a.NumericID == "" {
		return false
	}
	return true
}

// MarkAsTrash はアカウントをゴミ箱状態へ遷移させます
func (a *Account) MarkAsTrash(reason, by string) {
	a.IsTrash = true
	a.TrashReason = reason
	a.TrashedBy = by
	now := time.Now()
	a.TrashedAt = &now
	a.UpdatedAt = now
}
