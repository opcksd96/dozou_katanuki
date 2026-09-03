// adapters/driving/dto/account_dto.go (100行以下 - SPEC-PRINCIPLE-001)
package dto

import "time"

// AccountDTO はフロントエンドへアカウント情報を返却するためのデータ転送オブジェクトです
type AccountDTO struct {
	NumericID    string     `json:"numeric_id"`
	Username     string     `json:"username"`
	DisplayName  string     `json:"display_name"`
	AvatarURL    string     `json:"avatar_url"`
	AvatarBase64 string     `json:"avatar_base64"`
	Description  string     `json:"description"`
	GroupName    string     `json:"group_name"`
	AliasOf      string     `json:"alias_of"`
	IsWhitelist  bool       `json:"is_whitelist"`
	IsTrash      bool       `json:"is_trash"`
	TrashedBy    string     `json:"trashed_by,omitempty"`
	TrashReason  string     `json:"trash_reason,omitempty"`
	TrashedAt    *time.Time `json:"trashed_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	PostCount    int64      `json:"post_count"`
}

// AccountProfileHistoryDTO はアカウントの過去のプロフィール変更履歴です
type AccountProfileHistoryDTO struct {
	ID                uint      `json:"id"`
	AccountID         string    `json:"account_id"`
	DisplayName       string    `json:"display_name"`
	Description       string    `json:"description"`
	AvatarOriginalURL string    `json:"avatar_original_url"`
	AvatarSeq         int       `json:"avatar_seq"`
	AvatarVirtualKey  string    `json:"avatar_virtual_key"`
	AvatarBase64      string    `json:"avatar_base64"`
	ObservedAt        time.Time `json:"observed_at"`
}

// AccountDetailResult はアカウント詳細画面（Wails/WebUI両用）への返却モデルです
type AccountDetailResult struct {
	Account   AccountDTO                 `json:"account"`
	Histories []AccountProfileHistoryDTO `json:"histories"`
	PostCount int64                      `json:"post_count"`
}
