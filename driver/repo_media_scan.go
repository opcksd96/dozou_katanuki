// driver/repo_media_scan.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"database/sql"
	"fmt"
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm"
)

type mediaScanRow struct {
	MediaID string `gorm:"column:media_id"`; ArticleID string `gorm:"column:article_id"`
	Type string `gorm:"column:type"`; DownloadURL string `gorm:"column:download_url"`
	Width int `gorm:"column:width"`; Height int `gorm:"column:height"`
	DownloadStatus string `gorm:"column:download_status"`; FailedReason string `gorm:"column:failed_reason"`
	StashSceneID string `gorm:"column:stash_scene_id"`; StashImageID string `gorm:"column:stash_image_id"`
	IsBookmarked bool `gorm:"column:is_bookmarked"`; MediaQuality string `gorm:"column:media_quality"`
	AccountID string `gorm:"column:account_id"`; Username string `gorm:"column:username"`
	DisplayName string `gorm:"column:display_name"`; CreatedAt *time.Time `gorm:"column:created_at"`
	FullText string `gorm:"column:full_text"`; FullTextJA string `gorm:"column:full_text_ja"`
	WaybackURL string `gorm:"column:wayback_url"`
}

func (r *Repository) FetchRawMediaItems(accountID, status, mediaType string, limit, offset int) ([]models.MediaScanItem, int64, models.MediaSearchStats, error) {
	newQ := func() *gorm.DB {
		q := r.db.Table("media").
			Joins("LEFT JOIN articles ON articles.id = media.article_id").
			Joins("LEFT JOIN accounts ON accounts.numeric_id = articles.account_id")
		if accountID != "" && accountID != "all" {
			q = q.Where("accounts.numeric_id = ? OR accounts.username = ? OR articles.account_id = ?", accountID, accountID, accountID)
		}
		return q
	}

	var stats models.MediaSearchStats
	_ = newQ().Count(&stats.TotalCount).Error
	_ = newQ().Where("media.type = 'image'").Count(&stats.ImageCount).Error
	stats.VideoCount = stats.TotalCount - stats.ImageCount

	filterQ := newQ()
	if status != "" && status != "all" { filterQ = filterQ.Where("media.download_status = ?", status) }
	if mediaType == "image" { filterQ = filterQ.Where("media.type = 'image'") } else if mediaType == "video" { filterQ = filterQ.Where("media.type != 'image'") }

	var total int64
	if err := filterQ.Count(&total).Error; err != nil { return nil, 0, stats, err }
	if limit <= 0 { limit = 24 }

	var rows []mediaScanRow
	err := filterQ.Select("media.media_id, media.article_id, media.type, media.download_url, media.width, media.height, media.download_status, media.failed_reason, media.stash_scene_id, media.stash_image_id, media.is_bookmarked, media.media_quality, articles.created_at as created_at, articles.full_text as full_text, articles.full_text_ja as full_text_ja, accounts.username as username, accounts.display_name as display_name, accounts.numeric_id as account_id, articles.wayback_url as wayback_url").
		Order("COALESCE(articles.created_at, '1970-01-01') DESC, media.media_id DESC").Limit(limit).Offset(offset).Scan(&rows).Error
	if err != nil { return nil, 0, stats, err }

	accIDs := make([]string, 0, len(rows)); accIDMap := make(map[string]bool)
	for _, row := range rows {
		if row.AccountID != "" && !accIDMap[row.AccountID] { accIDMap[row.AccountID] = true; accIDs = append(accIDs, row.AccountID) }
	}

	var allHistories []models.AccountProfileHistory
	if len(accIDs) > 0 { _ = r.db.Where("account_id IN ?", accIDs).Order("observed_at ASC, avatar_seq ASC").Find(&allHistories).Error }
	historyByAccount := make(map[string][]models.AccountProfileHistory)
	for _, h := range allHistories { historyByAccount[h.AccountID] = append(historyByAccount[h.AccountID], h) }

	rawItems := make([]models.MediaScanItem, 0, len(rows))
	for _, row := range rows {
		t := time.Time{}; if row.CreatedAt != nil { t = *row.CreatedAt }
		m := models.Media{
			MediaID: row.MediaID, ArticleID: row.ArticleID, Type: row.Type, DownloadURL: row.DownloadURL,
			Width: row.Width, Height: row.Height, DownloadStatus: row.DownloadStatus,
			IsBookmarked: row.IsBookmarked, MediaQuality: row.MediaQuality,
		}
		if row.FailedReason != "" { m.FailedReason = sql.NullString{String: row.FailedReason, Valid: true} }
		if row.StashSceneID != "" { m.StashSceneID = sql.NullString{String: row.StashSceneID, Valid: true} }
		if row.StashImageID != "" { m.StashImageID = sql.NullString{String: row.StashImageID, Valid: true} }

		rawItems = append(rawItems, models.MediaScanItem{
			Media: m, ArticleID: row.ArticleID, AccountID: row.AccountID,
			Username: row.Username, DisplayName: row.DisplayName, CreatedAt: t,
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
		tweetDate := "-"; if !item.CreatedAt.IsZero() { tweetDate = item.CreatedAt.Format("2006-01-02") }
		items = append(items, models.MediaItemDetail{
			RenderMedia: rm, MediaID: item.Media.MediaID, ArticleID: item.ArticleID, AccountID: item.AccountID,
			Username: item.Username, DisplayName: item.DisplayName, RawStatus: item.Media.DownloadStatus,
			HasStash: hasStash, CreatedAt: item.CreatedAt, Title: fmt.Sprintf("X (@%s): Tweet %s", item.Username, item.ArticleID),
			FullText: item.FullText, FullTextJA: item.FullTextJA, TweetDate: tweetDate, WaybackURL: item.WaybackURL,
		})
	}
	return &models.MediaSearchResult{Items: items, Total: total, Stats: stats}, nil
}
