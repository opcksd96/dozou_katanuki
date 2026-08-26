// driver/repo_media_scan.go (100行以下)
package driver

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

type mediaScanRow struct {
	models.Media
	ArticleID   string    `gorm:"column:article_id"`
	AccountID   string    `gorm:"column:account_id"`
	Username    string    `gorm:"column:username"`
	DisplayName string    `gorm:"column:display_name"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	FullText    string    `gorm:"column:full_text"`
	FullTextJA  string    `gorm:"column:full_text_ja"`
	WaybackURL  string    `gorm:"column:wayback_url"`
}

func (r *Repository) FetchRawMediaItems(accountID, status, mediaType string, limit, offset int) ([]models.MediaScanItem, int64, models.MediaSearchStats, error) {
	baseQ := r.db.Table("media").
		Joins("JOIN articles ON articles.id = media.article_id").
		Joins("JOIN accounts ON accounts.numeric_id = articles.account_id")

	if accountID != "" && accountID != "all" {
		baseQ = baseQ.Where("accounts.numeric_id = ? OR accounts.username = ?", accountID, accountID)
	}

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

	q := baseQ
	if status != "" && status != "all" { q = q.Where("media.download_status = ?", status) }
	if mediaType == "image" { q = q.Where("media.type = 'image'") } else if mediaType == "video" { q = q.Where("media.type != 'image'") }

	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, stats, err }
	if limit <= 0 { limit = 24 }

	var rows []mediaScanRow
	err := q.Select("media.*, articles.created_at as created_at, articles.full_text as full_text, articles.full_text_ja as full_text_ja, accounts.username as username, accounts.display_name as display_name, accounts.numeric_id as account_id, articles.wayback_url as wayback_url").
		Order("articles.created_at DESC").Limit(limit).Offset(offset).Scan(&rows).Error
	if err != nil { return nil, 0, stats, err }

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
			Media: row.Media, ArticleID: row.ArticleID, AccountID: row.AccountID,
			Username: row.Username, DisplayName: row.DisplayName, CreatedAt: row.CreatedAt,
			FullText: row.FullText, FullTextJA: row.FullTextJA, ProfileHistory: historyByAccount[row.AccountID],
			WaybackURL: row.WaybackURL,
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
			RenderMedia: rm, MediaID: item.Media.MediaID, ArticleID: item.ArticleID, AccountID: item.AccountID,
			Username: item.Username, DisplayName: item.DisplayName, RawStatus: item.Media.DownloadStatus,
			HasStash: hasStash, CreatedAt: item.CreatedAt, Title: fmt.Sprintf("X (@%s): Tweet %s", item.Username, item.ArticleID),
			FullText: item.FullText, FullTextJA: item.FullTextJA, TweetDate: item.CreatedAt.Format("2006-01-02"),
			WaybackURL: item.WaybackURL,
		})
	}
	return &models.MediaSearchResult{Items: items, Total: total, Stats: stats}, nil
}
