// driver/repo_thunder_tasks_grouped.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"dozou_katanuki/models"
)

// MediaWithCandidateTasks は作品(メディア)とそれに紐づく全候補URLタスクの階層モデル
type MediaWithCandidateTasks struct {
	MediaID        string               `json:"media_id"`
	ArticleID      string               `json:"article_id"`
	OriginalURL    string               `json:"original_url"`
	Type           string               `json:"type"`
	Width          int                  `json:"width"`
	Height         int                  `json:"height"`
	DownloadStatus string               `json:"download_status"`
	FailedReason   string               `json:"failed_reason"`
	AccountID      string               `json:"account_id"`
	MediaQuality   string               `json:"media_quality"`
	CandidateTasks []models.ThunderTask `json:"candidate_tasks"`
}

// FetchMediaWithCandidateTasks は 指定ステータス(ESCALATED等)の作品と各候補URLタスク群を取得します
func (r *Repository) FetchMediaWithCandidateTasks(status string, limit int) ([]MediaWithCandidateTasks, error) {
	var result []MediaWithCandidateTasks
	if r.db == nil { return result, nil }
	if limit <= 0 { limit = 50 }
	if status == "" { status = "ESCALATED" }

	var medias []models.Media
	err := r.db.Where("download_status = ?", status).
		Order("(SELECT COUNT(*) FROM thunder_tasks WHERE thunder_tasks.media_id = media.media_id) DESC, media_id ASC").
		Limit(limit).
		Find(&medias).Error
	if err != nil || len(medias) == 0 { return result, err }

	mediaIDs := make([]string, len(medias))
	for i, m := range medias { mediaIDs[i] = m.MediaID }

	var tasks []models.ThunderTask
	_ = r.db.Where("media_id IN ?", mediaIDs).Order("created_at ASC").Find(&tasks).Error

	taskMap := make(map[string][]models.ThunderTask)
	for _, t := range tasks {
		taskMap[t.MediaID] = append(taskMap[t.MediaID], t)
	}

	for _, m := range medias {
		reason := ""
		if m.FailedReason.Valid { reason = m.FailedReason.String }
		cTasks := taskMap[m.MediaID]
		if cTasks == nil { cTasks = []models.ThunderTask{} }
		result = append(result, MediaWithCandidateTasks{
			MediaID:        m.MediaID,
			ArticleID:      m.ArticleID,
			OriginalURL:    m.DownloadURL,
			Type:           m.Type,
			Width:          m.Width,
			Height:         m.Height,
			DownloadStatus: m.DownloadStatus,
			FailedReason:   reason,
			AccountID:      m.AccountID,
			MediaQuality:   m.MediaQuality,
			CandidateTasks: cTasks,
		})
	}
	return result, nil
}
