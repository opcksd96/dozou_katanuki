// driver/repo_account_trash.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"fmt"
	"strings"
	"time"

	"dozou_katanuki/models"
)

// TrashAccount は指定されたアカウントを論理削除（ゴミ箱へ移動）します
func (r *Repository) TrashAccount(numericID, reason, trashedBy string) error {
	id := strings.TrimSpace(numericID)
	if id == "" {
		return fmt.Errorf("numericID must not be empty")
	}
	now := time.Now()
	res := r.db.Model(&models.Account{}).Where("numeric_id = ? OR username = ? OR LOWER(username) = LOWER(?)", id, id, id).Updates(map[string]interface{}{
		"is_trash":     true,
		"trash_reason": reason,
		"trashed_by":   trashedBy,
		"trashed_at":   &now,
		"is_whitelist": false,
		"updated_at":   now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("account not found: %s", id)
	}
	return nil
}

// RestoreAccount はゴミ箱に入っているアカウントを通常状態へ復元します
func (r *Repository) RestoreAccount(numericID string) error {
	id := strings.TrimSpace(numericID)
	if id == "" {
		return fmt.Errorf("numericID must not be empty")
	}
	now := time.Now()
	res := r.db.Model(&models.Account{}).Where("numeric_id = ? OR username = ? OR LOWER(username) = LOWER(?)", id, id, id).Updates(map[string]interface{}{
		"is_trash":     false,
		"trash_reason": nil,
		"trashed_by":   nil,
		"trashed_at":   nil,
		"updated_at":   now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("account not found: %s", id)
	}
	return nil
}

// GetTrashedAccounts はゴミ箱に入っているアカウント一覧を取得します
func (r *Repository) GetTrashedAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Select("accounts.*, (SELECT count(*) FROM articles WHERE articles.account_id = accounts.numeric_id) as post_count").
		Where("accounts.is_trash = ?", true).Preload("ProfileHistory").Order("trashed_at DESC, username ASC").Find(&accounts).Error
	return accounts, err
}
