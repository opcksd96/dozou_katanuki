// app/app_rpc_pipeline_reset_batch.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

// ResetSpecificMediasToQueued は指定されたメディアのリストを QUEUED へ差し戻し、download_tasks を再構築します
func (a *App) ResetSpecificMediasToQueued(mediaIDs []string) (int64, error) {
	if len(mediaIDs) == 0 {
		return 0, nil
	}
	if a.Repo == nil || a.Repo.DB() == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	db := a.Repo.DB()

	// 1. 対象メディアを QUEUED にリセット (個別明示指定のためCOMPLETED含む)
	res := db.Model(&models.Media{}).
		Where("media_id IN ?", mediaIDs).
		Updates(map[string]interface{}{
			"download_status": "QUEUED",
			"failed_reason":   nil,
		})
	if res.Error != nil {
		return 0, res.Error
	}
	count := res.RowsAffected
	if count == 0 {
		return 0, nil
	}

	// 2. 指定されたメディアIDに関連する旧汚染タスクを一掃
	// まずメディアの実体を取得
	var targetMedias []models.Media
	_ = db.Where("media_id IN ?", mediaIDs).
		Where("download_status = 'QUEUED'").Find(&targetMedias).Error

	var targetMediaIDs []string
	for _, m := range targetMedias {
		if m.MediaID != "" {
			targetMediaIDs = append(targetMediaIDs, m.MediaID)
		}
	}

	if len(targetMediaIDs) > 0 {
		_ = db.Where("media_id IN ?", targetMediaIDs).Delete(&models.DownloadTask{}).Error
	}

	// 3. QUEUED に差し戻されたメディアの download_tasks を再登録
	var allTasks []models.DownloadTask
	for _, m := range targetMedias {
		tasks := ExpandMediaCandidateTasks(m, models.StageRequests)
		allTasks = append(allTasks, tasks...)
	}

	if len(allTasks) > 0 {
		_ = a.Repo.BatchUpsertDownloadTasks(allTasks)
	}

	a.AppendPipelineLog("SYSTEM", "INFO", fmt.Sprintf("バッチ差し戻し: %d 件のメディアを QUEUED に差し戻し、%d 件の正規タスクを再構築", count, len(allTasks)))

	// 4. パイプラインを自動点火して処理を開始
	go func() {
		_, _ = a.IgnitePipeline()
	}()

	return count, nil
}
