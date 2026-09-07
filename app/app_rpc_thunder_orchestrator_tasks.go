// app/app_rpc_thunder_orchestrator_tasks.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

func (a *App) buildThunderOrchestratorTasks() []*models.ThunderOrchestratorTask {
	var tasks []*models.ThunderOrchestratorTask
	if a.Repo == nil { return tasks }

	items, _, _, err := a.Repo.FetchRawMediaItems("", "ESCALATED", "", 1000, 0)
	if err != nil || len(items) == 0 { return tasks }

	destDir := a.getMediaDownloadDir()
	existingStatusMap := make(map[string]models.ThunderTaskStatus)
	activeMediaMap := make(map[string]bool)
	if a.Repo != nil && a.Repo.DB() != nil {
		var existing []models.ThunderTask
		_ = a.Repo.DB().Select("id, media_id, status").Find(&existing).Error
		for _, ex := range existing {
			existingStatusMap[ex.ID] = ex.Status
			if ex.Status == models.ThunderTaskOnboarded || ex.Status == models.ThunderTaskRunning || ex.Status == models.ThunderTaskHolding {
				activeMediaMap[ex.MediaID] = true
			}
		}
	}

	var dTasks []models.DownloadTask
	var tTasks []models.ThunderTask
	seenMedia := make(map[string]bool)

	for _, item := range items {
		m := item.Media
		if seenMedia[m.MediaID] || m.IsTrash || activeMediaMap[m.MediaID] { continue }
		seenMedia[m.MediaID] = true

		cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
		ext := ".jpg"
		if m.Type == "video" || strings.Contains(m.DownloadURL, ".mp4") || strings.HasSuffix(strings.ToLower(m.MediaID), ".mp4") { ext = ".mp4" }

		// ローカルディスクに実ファイルが既に存在する場合は COMPLETED に昇格して投入回避
		if fi, err := os.Stat(filepath.Join(destDir, cleanID+ext)); err == nil && fi.Size() > 0 {
			_ = a.Repo.UpdateMediaMetadata(m.MediaID, "COMPLETED", "", "", "実ファイル確認済み(自動昇格)")
			continue
		}

		candidates := BuildCandidateURLsFromMediaWithArticle(m.MediaID, m.DownloadURL, m.Type, m.ArticleID)
		// 同一メディアからは未着手(PENDING)の最優先1候補のみを単一選択(多重投入・重複ダイアログの根絶)
		for _, c := range candidates {
			tID := fmt.Sprintf("%s-%s", cleanID, c.Type)
			fileName := fmt.Sprintf("%s_%s%s", cleanID, c.Type, ext)

			// 過去に一度でも投入・処理されたタスクは除外
			if st, exists := existingStatusMap[tID]; exists && st != models.ThunderTaskPending { continue }

			tasks = append(tasks, &models.ThunderOrchestratorTask{
				ID: tID, MediaID: m.MediaID, ArticleID: m.ArticleID, AccountID: item.AccountID,
				Username: item.Username, ResolutionType: c.Type, URL: c.URL, FileName: fileName, Status: "pending", SlotIndex: -1,
			})
			tTasks = append(tTasks, models.ThunderTask{
				ID: tID, MediaID: m.MediaID, ArticleID: m.ArticleID, ResolutionType: string(c.Type),
				URL: c.URL, FileName: fileName, Status: models.ThunderTaskPending,
			})
			break // 1メディアにつき最優先の1候補のみ投入！
		}
		dTasks = append(dTasks, models.DownloadTask{
			MediaID: m.MediaID, ArticleID: m.ArticleID, Stage: models.StageThunder,
			URL: m.DownloadURL, FileName: m.MediaID, Status: models.TaskPending,
		})
	}

	if len(dTasks) > 0 { _ = a.Repo.BatchUpsertDownloadTasks(dTasks) }
	if len(tTasks) > 0 { _ = a.Repo.BatchUpsertThunderTasks(tTasks) }
	return tasks
}
