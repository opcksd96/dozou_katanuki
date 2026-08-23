// driver/repo_generic.go (100行以下)
package driver

import (
	"fmt"

	"dozou_katanuki/models"
)

var allowedTables = map[string]bool{
	"accounts": true, "account_profile_histories": true, "articles": true,
	"media": true, "url_redirects": true, "whitelists": true,
}

func (r *Repository) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	var acc models.Account
	if err := r.db.Where("numeric_id = ? OR username = ?", numericID, numericID).First(&acc).Error; err != nil {
		return nil, err
	}
	var hist []models.AccountProfileHistory
	_ = r.db.Where("account_id = ?", acc.NumericID).Order("avatar_seq ASC").Find(&hist).Error

	var postCount int64
	_ = r.db.Model(&models.Article{}).Where("account_id = ?", acc.NumericID).Count(&postCount).Error
	return &models.AccountDetailResult{Account: acc, Histories: hist, PostCount: postCount}, nil
}

func (r *Repository) GetTableRecords(tableName string, limit, offset int, search string) (*models.TableRecordResult, error) {
	if !allowedTables[tableName] { return nil, fmt.Errorf("table not allowed: %s", tableName) }
	q := r.db.Table(tableName)
	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, err }

	var rows []map[string]interface{}
	if limit <= 0 { limit = 50 }
	if err := q.Limit(limit).Offset(offset).Find(&rows).Error; err != nil { return nil, err }

	var cols []string
	if len(rows) > 0 {
		for k := range rows[0] { cols = append(cols, k) }
	}
	return &models.TableRecordResult{TableName: tableName, Columns: cols, Rows: rows, Total: total}, nil
}

func (r *Repository) ListAccounts() ([]models.Account, error) {
	var accs []models.Account
	err := r.db.Order("updated_at DESC").Find(&accs).Error
	return accs, err
}
