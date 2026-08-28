// driver/repo_accounts.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"strings"
	"time"

	"dozou_katanuki/models"
)

func (r *Repository) GetAccounts() ([]models.Account, error) {
	return r.GetAllAccounts(false)
}

func (r *Repository) GetAllAccounts(includeTrash bool) ([]models.Account, error) {
	var accounts []models.Account
	q := r.db.Select("accounts.*, (SELECT count(*) FROM articles WHERE articles.account_id = accounts.numeric_id AND (articles.is_trash = 0 OR articles.is_trash IS NULL)) as post_count").
		Preload("ProfileHistory").Order("is_whitelist DESC, username ASC")
	if !includeTrash {
		q = q.Where("accounts.is_trash = ?", false)
	}
	err := q.Find(&accounts).Error
	return accounts, err
}

func (r *Repository) ToggleAccountWhitelist(numericID string, isWhitelist bool) error {
	var acc models.Account
	if err := r.db.Where("numeric_id = ?", numericID).First(&acc).Error; err != nil {
		return err
	}
	if err := r.db.Model(&acc).Update("is_whitelist", isWhitelist).Error; err != nil {
		return err
	}
	lowerUser := strings.ToLower(acc.Username)
	var wl models.Whitelist
	if err := r.db.Where("LOWER(value) = ?", lowerUser).First(&wl).Error; err == nil {
		_ = r.db.Model(&wl).Update("is_active", isWhitelist).Error
	} else if isWhitelist {
		_ = r.db.Create(&models.Whitelist{
			Type: "account", Value: lowerUser, GroupName: acc.GroupName, AliasOf: acc.AliasOf, IsActive: true,
		}).Error
	}
	return nil
}

func (r *Repository) UpdateAccount(numericID, displayName, username, avatarURL, description, aliasOf, groupName string) error {
	lowerUser, lowerAlias := strings.ToLower(username), strings.ToLower(aliasOf)
	return r.db.Model(&models.Account{}).Where("numeric_id = ?", numericID).Updates(map[string]interface{}{
		"display_name": displayName, "username": lowerUser, "avatar_url": avatarURL,
		"description": description, "alias_of": lowerAlias, "group_name": groupName, "updated_at": time.Now(),
	}).Error
}

func (r *Repository) UpdateAvatarBase64ByVirtualKey(virtualKey, b64Data string) error {
	var hist models.AccountProfileHistory
	if err := r.db.Where("avatar_virtual_key = ? OR avatar_virtual_key = ?", virtualKey, virtualKey+".jpg").First(&hist).Error; err == nil {
		_ = r.db.Model(&models.AccountProfileHistory{}).Where("id = ?", hist.ID).Update("avatar_base64", b64Data).Error
		_ = r.db.Model(&models.Account{}).Where("numeric_id = ?", hist.AccountID).Update("avatar_base64", b64Data).Error
		return nil
	}
	_ = r.db.Model(&models.Account{}).Where("username = ? OR avatar_url LIKE ?", strings.ToLower(virtualKey), "%"+virtualKey+"%").Update("avatar_base64", b64Data).Error
	return nil
}
