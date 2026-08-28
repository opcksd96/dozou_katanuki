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
	if a.Repo == nil {
		return tasks
	}
	var allItems []models.MediaScanItem
	for _, st := range []string{"OUTSOURCED", "RETAINED", "ESCALATED"} {
		items, _, _, err := a.Repo.FetchRawMediaItems("", st, "", 1000, 0)
		if err == nil && len(items) > 0 {
			allItems = append(allItems, items...)
		}
	}

	seenMedia := make(map[string]bool)
	for _, item := range allItems {
		m := item.Media
		if m.DownloadURL == "" || seenMedia[m.MediaID] { continue }
		seenMedia[m.MediaID] = true

		cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
		ext := ".jpg"
		if m.Type == "video" || strings.Contains(m.DownloadURL, ".mp4") || strings.HasSuffix(strings.ToLower(m.MediaID), ".mp4") {
			ext = ".mp4"
		}
		candidates := BuildThunderTop3CandidateURLs(m.DownloadURL)
		for _, c := range candidates {
			tasks = append(tasks, &models.ThunderOrchestratorTask{
				ID:             fmt.Sprintf("%s-%s", cleanID, c.Type),
				MediaID:        m.MediaID,
				ArticleID:      m.ArticleID,
				ResolutionType: c.Type,
				URL:            c.URL,
				FileName:       fmt.Sprintf("%s_%s%s", cleanID, c.Type, ext),
				Status:         "pending",
				SlotIndex:      -1,
			})
		}
	}
	return tasks
}
