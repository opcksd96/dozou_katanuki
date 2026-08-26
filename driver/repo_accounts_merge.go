// driver/repo_accounts_merge.go (100行以下)
package driver

import (
	"fmt"

	"dozou_katanuki/models"
	"gorm.io/gorm"
)

// MergeAccounts transfers all articles from source to target, then marks the
// source account as an alias of the target (sets alias_of). The whole operation
// runs inside a single transaction so that articles and the alias flag stay
// consistent.
func (r *Repository) MergeAccounts(sourceNumericID, targetNumericID string) error {
	if sourceNumericID == "" || targetNumericID == "" {
		return fmt.Errorf("both sourceNumericID and targetNumericID must be non-empty")
	}
	if sourceNumericID == targetNumericID {
		return fmt.Errorf("source and target accounts must differ")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Article{}).Where("account_id = ?", sourceNumericID).Update("account_id", targetNumericID).Error; err != nil {
			return fmt.Errorf("failed to reassign articles: %w", err)
		}
		if err := tx.Model(&models.Account{}).Where("numeric_id = ?", sourceNumericID).Update("alias_of", targetNumericID).Error; err != nil {
			return fmt.Errorf("failed to set alias_of: %w", err)
		}
		return nil
	})
}
