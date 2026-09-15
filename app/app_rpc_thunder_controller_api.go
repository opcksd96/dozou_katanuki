// app/app_rpc_thunder_controller_api.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

// RegisterThunderTasksFromDB は thunder_tasks テーブルの PENDING タスクを迅雷へ投入します (要件①)
func (a *App) RegisterThunderTasksFromDB(limit int, saveDir string) (int, error) {
	if a.Repo == nil { return 0, fmt.Errorf("リポジトリ未初期化") }
	if limit <= 0 { limit = 3 }
	if saveDir == "" { saveDir = a.getMediaDownloadDir() }

	tasks, err := a.Repo.FetchPendingThunderTasks(limit)
	if err != nil || len(tasks) == 0 { return 0, err }

	var items []models.ThunderTaskInputItem
	for _, t := range tasks {
		fn := t.FileName
		if fn == "" { fn = resolveMediaIDFromFileName(t.URL) }
		items = append(items, models.ThunderTaskInputItem{URL: t.URL, FileName: fn})
	}

	if err := AddTasksViaDirectCDP(9222, items, saveDir); err != nil {
		return 0, fmt.Errorf("Direct CDP 投入エラー: %w", err)
	}

	for _, t := range tasks {
		_ = a.Repo.MarkThunderTaskRunning(t.ID)
		a.AppendPipelineLog("THUNDER", "SUCCESS", fmt.Sprintf("🚀 迅雷Direct投入: %s", t.FileName))
	}
	return len(tasks), nil
}

// GetThunderActiveTasks は 迅雷内のアクティブタスク一覧・リアルタイム詳細を取得します (要件②)
func (a *App) GetThunderActiveTasks() ([]models.ActiveThunderTask, error) {
	return GetTasksViaDirectCDP(9222)
}

// ReapFailedEmptyThunderTasks は 0B/0Bかつタイムアウト以外の失敗タスクをゴミ箱へ移動します (要件③)
func (a *App) ReapFailedEmptyThunderTasks() (int, error) {
	if a.Repo == nil { return 0, nil }
	tasks, err := GetTasksViaDirectCDP(9222)
	if err != nil { return 0, err }

	reapedCount := 0
	for _, t := range tasks {
		// ⚡ 巻き添え防止: thunder_tasks DB に登録されている管理対象タスク以外は絶対に触らない！
		if !a.Repo.IsManagedThunderTask(t.FileName) { continue }
		if !isFailedEmptyTask(t) { continue }

		success, err := DeleteTaskViaDirectCDP(9222, t.FileName)
		if err == nil && success {
			reapedCount++
			allReaped, mediaID, _ := a.Repo.MarkThunderTaskReaped(t.FileName, t.DetailText)
			if allReaped && mediaID != "" {
				_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "全候補タスク枯渇によりRETAINED退避")
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("📦 全候補枯渇のため退避: %s", mediaID))
			}
			a.AppendPipelineLog("THUNDER", "WARN", fmt.Sprintf("🗑️ 停止中枯渇タスクをゴミ箱へ移動: %s (%s)", t.FileName, t.DetailText))
		}
	}
	return reapedCount, nil
}

// ReactivateThunderTasks は 一時停止中・タイムアウトタスクを再活性化します (要件④)
func (a *App) ReactivateThunderTasks(includeDB bool) (int, error) {
	resumedCount := 0
	// 迅雷内部の再活性化: DB登録済みの停止タスクのみを対象に安全に再開
	if a.Repo != nil {
		tasks, err := GetTasksViaDirectCDP(9222)
		if err == nil {
			for _, t := range tasks {
				if a.Repo.IsManagedThunderTask(t.FileName) && (t.TaskStatusCode == 7 || t.TaskStatusCode == 8) {
					if ok, _ := ResumeTaskViaDirectCDP(9222, t.FileName); ok { resumedCount++ }
				}
			}
		}
		if includeDB {
			if n, err := a.Repo.ReactivatePausedThunderTasks(); err == nil && n > 0 {
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("♻️ DB内の保留タスクを再活性化: %d 件", n))
			}
		}
	}
	return resumedCount, nil
}
