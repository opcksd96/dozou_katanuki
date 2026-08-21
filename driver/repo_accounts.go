// driver/repo_accounts.go (100行以下)
package driver

import (
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
