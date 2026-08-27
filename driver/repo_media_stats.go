// driver/repo_media_stats.go (100行以下)
package driver

import "dozou_katanuki/models"

type statusCountRow struct {
	Status string `gorm:"column:status"`
	Count  int64  `gorm:"column:count"`
}

// FetchDownloadStatusStats はダウンロードキュー全体をステータス別に集計します
func (r *Repository) FetchDownloadStatusStats(accountID string) (*models.DownloadStatusStats, error) {
	q := r.db.Table("media").
		Joins("JOIN articles ON articles.id = media.article_id").
		Joins("JOIN accounts ON accounts.numeric_id = articles.account_id")
	if accountID != "" && accountID != "all" {
		q = q.Where("accounts.numeric_id = ? OR accounts.username = ?", accountID, accountID)
	}

	var rows []statusCountRow
	if err := q.Select("media.download_status as status, COUNT(*) as count").Group("media.download_status").Scan(&rows).Error; err != nil {
		return nil, err
	}

	stats := models.DownloadStatusStats{}
	for _, row := range rows {
		stats.Total += row.Count
		switch row.Status {
		case "QUEUED", "DOWNLOADING":
			stats.Queued += row.Count
		case "COMPLETED":
			stats.Completed += row.Count
		case "DEAD_404":
			stats.Dead404 += row.Count
		case "OUTSOURCED":
			stats.Outsourced += row.Count
		case "RETAINED":
			stats.Retained += row.Count
		case "FAILED":
			stats.Failed += row.Count
		}
	}
	return &stats, nil
}
