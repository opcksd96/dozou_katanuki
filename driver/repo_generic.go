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
	ArticleID   string    `gorm:"column:article_id"`
	AccountID   string    `gorm:"column:account_id"`
	Username    string    `gorm:"column:username"`
	DisplayName string    `gorm:"column:display_name"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

// FetchRawMediaItems はリポジトリ層でメディアとアカウント・世代履歴を取得します
func (r *Repository) FetchRawMediaItems(accountID, status, mediaType string, limit, offset int) ([]models.MediaScanItem, int64, models.MediaSearchStats, error) {
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
	if err := q.Count(&total).Error; err != nil { return nil, 0, stats, err }
	if limit <= 0 { limit = 20 }

	var rows []mediaScanRow
	err := q.Select("media.*, articles.created_at as created_at, accounts.username as username, accounts.display_name as display_name, accounts.numeric_id as account_id").
		Order("articles.created_at DESC").Limit(limit).Offset(offset).Scan(&rows).Error
	if err != nil { return nil, 0, stats, err }

	// 取得した各アカウントの世代履歴を一括取得
	accIDs := make([]string, 0, len(rows))
	accIDMap := make(map[string]bool)
	for _, row := range rows {
		if row.AccountID != "" && !accIDMap[row.AccountID] {
			accIDMap[row.AccountID] = true
			accIDs = append(accIDs, row.AccountID)
		}
	}

	var allHistories []models.AccountProfileHistory
	if len(accIDs) > 0 {
		_ = r.db.Where("account_id IN ?", accIDs).Order("observed_at ASC, avatar_seq ASC").Find(&allHistories).Error
	}
	historyByAccount := make(map[string][]models.AccountProfileHistory)
	for _, h := range allHistories {
		historyByAccount[h.AccountID] = append(historyByAccount[h.AccountID], h)
	}

	rawItems := make([]models.MediaScanItem, 0, len(rows))
	for _, row := range rows {
		rawItems = append(rawItems, models.MediaScanItem{
			Media:          row.Media,
			ArticleID:      row.ArticleID,
			AccountID:      row.AccountID,
			Username:       row.Username,
			DisplayName:    row.DisplayName,
			CreatedAt:      row.CreatedAt,
			ProfileHistory: historyByAccount[row.AccountID],
		})
	}

	return rawItems, total, stats, nil
}

func (r *Repository) SearchMediaDetails(accountID, status, mediaType string, limit, offset int) (*models.MediaSearchResult, error) {
	rawItems, total, stats, err := r.FetchRawMediaItems(accountID, status, mediaType, limit, offset)
	if err != nil { return nil, err }

	items := make([]models.MediaItemDetail, 0, len(rawItems))
	for _, item := range rawItems {
		rm := models.BuildRenderMedia(item.Media)
		hasStash := item.Media.StashSceneID.Valid && item.Media.StashSceneID.String != "" || item.Media.StashImageID.Valid && item.Media.StashImageID.String != ""
		items = append(items, models.MediaItemDetail{
			RenderMedia: rm,
			MediaID:     item.Media.MediaID,
			ArticleID:   item.ArticleID,
			AccountID:   item.AccountID,
			Username:    item.Username,
			DisplayName: item.DisplayName,
			RawStatus:   item.Media.DownloadStatus,
			HasStash:    hasStash,
			CreatedAt:   item.CreatedAt,
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

