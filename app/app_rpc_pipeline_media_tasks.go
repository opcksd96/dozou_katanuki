// app/app_rpc_pipeline_media_tasks.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"regexp"
	"strings"

	"dozou_katanuki/models"
)

var reTrailingCopy = regexp.MustCompile(`\(\d+\)(\.[a-zA-Z0-9]+)$`)

// GetMediaWithCandidateTasks は 作品ごとの候補URLタスク群に迅雷CDP生タスク情報(THUNDER_TASK_SCHEMA)を付与して取得します
func (a *App) GetMediaWithCandidateTasks(status string, limit int) ([]models.MediaWithTasksEnriched, error) {
	if a.Repo == nil { return []models.MediaWithTasksEnriched{}, nil }
	if limit <= 0 { limit = 50 }
	if status == "" { status = "ESCALATED" }

	if isThunderProcessRunning() {
		_, _ = a.ReconcileThunderTasksWithDB()
	}
	_, _ = a.Repo.ReconcileDepletedMediaToRetained()

	rawItems, err := a.Repo.FetchMediaWithCandidateTasks(status, limit)
	if err != nil { return nil, err }

	cdpMap := make(map[string]*models.ActiveThunderTask)
	if isThunderProcessRunning() {
		if cdpTasks, err := GetTasksViaDirectCDP(9222); err == nil && len(cdpTasks) > 0 {
			for i := range cdpTasks {
				t := &cdpTasks[i]
				if u := strings.TrimSpace(t.URL); u != "" { cdpMap[u] = t }
				if name := strings.TrimSpace(t.FileName); name != "" {
					cdpMap[name] = t
					cdpMap[strings.ToLower(name)] = t
					clean := reTrailingCopy.ReplaceAllString(name, "$1")
					cdpMap[clean] = t
					cdpMap[strings.ToLower(clean)] = t
				}
			}
		}
	}

	result := make([]models.MediaWithTasksEnriched, len(rawItems))
	for i, m := range rawItems {
		enrichedTasks := make([]models.CandidateTaskEnriched, len(m.CandidateTasks))
		for j, ct := range m.CandidateTasks {
			enriched := models.CandidateTaskEnriched{ThunderTask: ct}
			if c, ok := cdpMap[ct.URL]; ok {
				enriched.CDP = c
			} else if c, ok := cdpMap[ct.FileName]; ok {
				enriched.CDP = c
			} else if c, ok := cdpMap[strings.ToLower(ct.FileName)]; ok {
				enriched.CDP = c
			}

			// DB永続化 / DBからの逆復元
			if enriched.CDP != nil && enriched.CDP.RawJSON != "" && ct.RawJSON == "" {
				_ = a.Repo.DB().Model(&models.ThunderTask{}).Where("id = ?", ct.ID).Update("raw_json", enriched.CDP.RawJSON).Error
			} else if enriched.CDP == nil && ct.RawJSON != "" {
				var restored models.ActiveThunderTask
				if err := json.Unmarshal([]byte(ct.RawJSON), &restored); err == nil { enriched.CDP = &restored }
			}
			enrichedTasks[j] = enriched
		}
		result[i] = models.MediaWithTasksEnriched{
			MediaID: m.MediaID, ArticleID: m.ArticleID, OriginalURL: m.OriginalURL,
			Type: m.Type, Width: m.Width, Height: m.Height, DownloadStatus: m.DownloadStatus,
			FailedReason: m.FailedReason, AccountID: m.AccountID, MediaQuality: m.MediaQuality,
			CandidateTasks: enrichedTasks,
		}
	}
	return result, nil
}
