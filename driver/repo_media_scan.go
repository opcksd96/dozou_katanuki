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
	IsTrash bool `gorm:"column:is_trash"`; TrashedBy string `gorm:"column:trashed_by"`; TrashReason string `gorm:"column:trash_reason"`
	AccountID string `gorm:"column:account_id"`; Username string `gorm:"column:username"`
	DisplayName string `gorm:"column:display_name"`; CreatedAt *time.Time `gorm:"column:created_at"`
	FullText string `gorm:"column:full_text"`; FullTextJA string `gorm:"column:full_text_ja"`
	WaybackURL string `gorm:"column:wayback_url"`; AvatarBase64 string `gorm:"column:avatar_base64"`
}

func (r *Repository) FetchRawMediaItems(accountID, status, mediaType string, limit, offset int) ([]models.MediaScanItem, int64, models.MediaSearchStats, error) {
	newQ := func() *gorm.DB {
		q := r.db.Table("media").
			Joins("LEFT JOIN articles ON articles.id = media.article_id").
			Joins("LEFT JOIN accounts ON accounts.numeric_id = articles.account_id")
		if accountID != "" && accountID != "all" { q = q.Where("accounts.numeric_id = ? OR accounts.username = ? OR articles.account_id = ?", accountID, accountID, accountID) }
		if status == "TRASH" || status == "trash" { q = q.Where("media.is_trash = ?", true) } else {
			q = q.Where("media.is_trash = 0 AND articles.is_trash = 0")
			if status != "" && status != "all" { q = q.Where("media.download_status = ?", status) }
		}
		if mediaType == "image" { q = q.Where("media.type = 'image'") } else if mediaType == "video" { q = q.Where("media.type != 'image'") }
		return q
	}

	var stats models.MediaSearchStats
	_ = newQ().Count(&stats.TotalCount).Error
	_ = newQ().Where("media.type = 'image'").Count(&stats.ImageCount).Error
	stats.VideoCount = stats.TotalCount - stats.ImageCount

	total := stats.TotalCount
	if limit <= 0 { limit = 24 }

	var rows []mediaScanRow
	err := newQ().Select("media.media_id, media.article_id, media.type, media.download_url, media.width, media.height, media.download_status, media.failed_reason, media.stash_scene_id, media.stash_image_id, media.is_bookmarked, media.media_quality, media.is_trash, media.trashed_by, media.trash_reason, articles.created_at, articles.full_text, articles.full_text_ja, accounts.username, accounts.display_name, accounts.numeric_id as account_id, articles.wayback_url, accounts.avatar_base64").
		Order("COALESCE(articles.created_at, '1970-01-01') DESC, media.media_id DESC").Limit(limit).Offset(offset).Scan(&rows).Error
	if err != nil { return nil, 0, stats, err }

	rawItems := make([]models.MediaScanItem, 0, len(rows))
	for _, row := range rows {
		t := time.Time{}; if row.CreatedAt != nil { t = *row.CreatedAt }
		m := models.Media{
			MediaID: row.MediaID, ArticleID: row.ArticleID, Type: row.Type, DownloadURL: row.DownloadURL,
			Width: row.Width, Height: row.Height, DownloadStatus: row.DownloadStatus,
			IsBookmarked: row.IsBookmarked, MediaQuality: row.MediaQuality, IsTrash: row.IsTrash, TrashedBy: row.TrashedBy, TrashReason: row.TrashReason,
		}
		if row.FailedReason != "" { m.FailedReason = sql.NullString{String: row.FailedReason, Valid: true} }
		if row.StashSceneID != "" { m.StashSceneID = sql.NullString{String: row.StashSceneID, Valid: true} }
		if row.StashImageID != "" { m.StashImageID = sql.NullString{String: row.StashImageID, Valid: true} }
		rawItems = append(rawItems, models.MediaScanItem{
			Media: m, ArticleID: row.ArticleID, AccountID: row.AccountID,
			Username: row.Username, DisplayName: row.DisplayName, CreatedAt: t,
			FullText: row.FullText, FullTextJA: row.FullTextJA, WaybackURL: row.WaybackURL,
			AvatarBase64: row.AvatarBase64,
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
