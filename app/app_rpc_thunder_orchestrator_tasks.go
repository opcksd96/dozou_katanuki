// app/app_rpc_thunder_orchestrator_tasks.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

func (a *App) buildThunderOrchestratorTasks() []*models.ThunderOrchestratorTask {
	var tasks []*models.ThunderOrchestratorTask
	if a.Repo == nil { return tasks }

	// 迅雷オーケストレーターは ESCALATED のメディアのみを対象とする
	items, _, _, err := a.Repo.FetchRawMediaItems("", "ESCALATED", "", 1000, 0)
	if err != nil || len(items) == 0 { return tasks }

	var dTasks []models.DownloadTask
	seenMedia := make(map[string]bool)

	for _, item := range items {
		m := item.Media
		if seenMedia[m.MediaID] || m.IsTrash { continue }
		seenMedia[m.MediaID] = true

		cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
		ext := ".jpg"
		if m.Type == "video" || strings.Contains(m.DownloadURL, ".mp4") || strings.HasSuffix(strings.ToLower(m.MediaID), ".mp4") {
			ext = ".mp4"
		}
		candidates := BuildCandidateURLsFromMediaWithArticle(m.MediaID, m.DownloadURL, m.Type, m.ArticleID)
		for _, c := range candidates {
			tID := fmt.Sprintf("%s-%s", cleanID, c.Type)
			fileName := fmt.Sprintf("%s_%s%s", cleanID, c.Type, ext)

			tasks = append(tasks, &models.ThunderOrchestratorTask{
				ID:             tID,
				MediaID:        m.MediaID,
				ArticleID:      m.ArticleID,
				AccountID:      item.AccountID,
				Username:       item.Username,
				ResolutionType: c.Type,
				URL:            c.URL,
				FileName:       fileName,
				Status:         "pending",
				SlotIndex:      -1,
			})

		}
		dTasks = append(dTasks, models.DownloadTask{
			MediaID:   m.MediaID,
			ArticleID: m.ArticleID,
			Stage:     models.StageThunder,
			URL:       m.DownloadURL,
			FileName:  m.MediaID,
			Status:    models.TaskPending,
		})
	}

	if len(dTasks) > 0 {
		_ = a.Repo.BatchUpsertDownloadTasks(dTasks)
	}

	return tasks
}
