// driver/repo_accounts.go (100行以下)
package driver

import (
	"time"

	"dozou_katanuki/models"
)

func (r *Repository) GetAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Preload("ProfileHistory").Order("username ASC").Find(&accounts).Error
	return accounts, err
}

func (r *Repository) GetAccountHistories(accountID string) ([]models.AccountProfileHistory, error) {
	var histories []models.AccountProfileHistory
	err := r.db.Where("account_id = ?", accountID).Order("avatar_seq desc").Find(&histories).Error
	return histories, err
}

func (r *Repository) UpdateAccount(numericID, displayName, username, avatarURL, description, aliasOf, groupName string) error {
	return r.db.Model(&models.Account{}).Where("numeric_id = ?", numericID).Updates(map[string]interface{}{
		"display_name": displayName,
		"username":     username,
		"avatar_url":   avatarURL,
		"description":  description,
		"alias_of":     aliasOf,
		"group_name":   groupName,
		"updated_at":   time.Now(),
	}).Error
}

func (r *Repository) UpdateAvatarBase64ByVirtualKey(virtualKey, b64Data string) error {
	// 1. account_profile_histories を更新
	var hist models.AccountProfileHistory
	if err := r.db.Where("avatar_virtual_key = ? OR avatar_virtual_key = ?", virtualKey, virtualKey+".jpg").First(&hist).Error; err == nil {
		_ = r.db.Model(&models.AccountProfileHistory{}).Where("id = ?", hist.ID).Update("avatar_base64", b64Data).Error
		// 2. 所属アカウントの最新アバターも更新
		_ = r.db.Model(&models.Account{}).Where("numeric_id = ?", hist.AccountID).Update("avatar_base64", b64Data).Error
		return nil
	}
	// virtualKey がアカウントのアバターURLまたはusernameの場合
	_ = r.db.Model(&models.Account{}).Where("username = ? OR avatar_url LIKE ?", virtualKey, "%"+virtualKey+"%").Update("avatar_base64", b64Data).Error
	return nil
}
