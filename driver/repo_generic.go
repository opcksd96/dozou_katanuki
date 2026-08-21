// driver/repo_generic.go (100行以下)
package driver

import (
	"fmt"
	"time"

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

type mediaScanRow struct {
	models.Media
	ArticleID string    `gorm:"column:article_id"`
	AccountID string    `gorm:"column:account_id"`
	Username  string    `gorm:"column:username"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (r *Repository) SearchMediaDetails(accountID, status, mediaType string, limit, offset int) (*models.MediaSearchResult, error) {
	baseQ := r.db.Table("media").
		Joins("JOIN articles ON articles.id = media.article_id").
		Joins("JOIN accounts ON accounts.numeric_id = articles.account_id")

	if accountID != "" && accountID != "all" {
		baseQ = baseQ.Where("accounts.numeric_id = ? OR accounts.username = ?", accountID, accountID)
	}

	// 統計カウント (同一アカウントフィルタ下での全体・画像・動画件数)
	var stats models.MediaSearchStats
	_ = baseQ.Count(&stats.TotalCount).Error
	_ = r.db.Table("media").
		Joins("JOIN articles ON articles.id = media.article_id").
		Joins("JOIN accounts ON accounts.numeric_id = articles.account_id").
		Where(func() string {
			if accountID != "" && accountID != "all" { return "(accounts.numeric_id = '" + accountID + "' OR accounts.username = '" + accountID + "') AND " }
			return ""
		}() + "media.type = 'image'").Count(&stats.ImageCount).Error
	stats.VideoCount = stats.TotalCount - stats.ImageCount

	// フィルタ条件適用
	q := baseQ
	if status != "" && status != "all" {
		q = q.Where("media.download_status = ?", status)
	}
	if mediaType == "image" {
		q = q.Where("media.type = 'image'")
	} else if mediaType == "video" {
		q = q.Where("media.type != 'image'")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, err }
	if limit <= 0 { limit = 20 }

	var rows []mediaScanRow
	err := q.Select("media.*, articles.created_at as created_at, accounts.username as username, accounts.numeric_id as account_id").
		Order("articles.created_at DESC").Limit(limit).Offset(offset).Scan(&rows).Error
	if err != nil { return nil, err }

	items := make([]models.MediaItemDetail, 0, len(rows))
	for _, row := range rows {
		rm := models.BuildRenderMedia(row.Media)
		hasStash := row.StashSceneID.Valid && row.StashSceneID.String != "" || row.StashImageID.Valid && row.StashImageID.String != ""
		items = append(items, models.MediaItemDetail{
			RenderMedia: rm,
			MediaID:     row.Media.MediaID,
			ArticleID:   row.ArticleID,
			AccountID:   row.AccountID,
			Username:    row.Username,
			RawStatus:   row.Media.DownloadStatus,
			HasStash:    hasStash,
			CreatedAt:   row.CreatedAt,
		})
	}

	return &models.MediaSearchResult{Items: items, Total: total, Stats: stats}, nil
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

