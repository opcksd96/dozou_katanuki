// app/app_rpc_pipeline_reset.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

// ResetAllToQueuedAndBootstrap は ESCALATED/OUTSOURCED/RETAINED/FAILED の全タスクを一括で QUEUED へ差し戻し、download_tasks を再構築します
func (a *App) ResetAllToQueuedAndBootstrap() (int64, error) {
	if a.Repo == nil || a.Repo.DB() == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	db := a.Repo.DB()

	// 1. 対象メディアを一括で QUEUED にリセット
	targetStatuses := []string{"ESCALATED", "OUTSOURCED", "RETAINED", "FAILED", "ERROR"}
	res := db.Model(&models.Media{}).
		Where("download_status IN ? AND (is_trash = 0 OR is_trash IS NULL)", targetStatuses).
		Updates(map[string]interface{}{
			"download_status": "QUEUED",
			"failed_reason":   nil,
		})
	if res.Error != nil { return 0, res.Error }
	count := res.RowsAffected

	// 2. QUEUED メディア全件を取得して download_tasks に Top3 候補を展開登録
	var queuedMedias []models.Media
	_ = db.Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Find(&queuedMedias).Error

	var allTasks []models.DownloadTask
	for _, m := range queuedMedias {
		tasks := ExpandMediaCandidateTasks(m, models.StageRequests)
		allTasks = append(allTasks, tasks...)
	}

	if len(allTasks) > 0 {
		_ = a.Repo.BatchUpsertDownloadTasks(allTasks)
	}

	a.AppendPipelineLog("SYSTEM", "INFO", fmt.Sprintf("パイプライン全件初期化: %d 件のメディアを QUEUED に差し戻し、%d 件のタスクを登録", count, len(allTasks)))

	return count, nil
}
