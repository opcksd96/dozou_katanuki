// driver/repo_download_tasks.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm/clause"
)

// BatchUpsertDownloadTasks は 全ステージの候補タスクを一括登録します
func (r *Repository) BatchUpsertDownloadTasks(tasks []models.DownloadTask) error {
	if len(tasks) == 0 || r.db == nil { return nil }
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"stage", "url", "file_name", "status", "updated_at"}),
	}).Create(&tasks).Error
}

// MarkTaskCompletedAndReapAllOthers は 勝利タスクを完了にし、他ステージ・他候補タスクを一括でREAPEDにします
func (r *Repository) MarkTaskCompletedAndReapAllOthers(mediaID, completedFileName string) ([]models.DownloadTask, error) {
	var reaps []models.DownloadTask
	if r.db == nil || mediaID == "" { return reaps, nil }
	now := time.Now()

	_ = r.db.Model(&models.DownloadTask{}).Where("media_id = ? AND file_name = ?", mediaID, completedFileName).
		Updates(map[string]interface{}{"status": models.TaskCompleted, "completed_at": &now}).Error

	_ = r.db.Where("media_id = ? AND file_name != ? AND status IN ?", mediaID, completedFileName, []string{string(models.TaskPending), string(models.TaskRunning)}).
		Find(&reaps).Error

	if len(reaps) > 0 {
		_ = r.db.Model(&models.DownloadTask{}).Where("media_id = ? AND file_name != ?", mediaID, completedFileName).
			Updates(map[string]interface{}{"status": models.TaskReaped, "reaped_at": &now}).Error
	}
	return reaps, nil
}

// MarkStageTaskDepleted は 特定ステージのタスクを枯渇済みにし、同一ステージの全滅判定を返します
func (r *Repository) MarkStageTaskDepleted(fileName string, stage models.PipelineStage) (stageAllDepleted bool, mediaID string, err error) {
	if r.db == nil || fileName == "" { return false, "", nil }
	var task models.DownloadTask
	if err := r.db.Where("file_name = ? AND stage = ?", fileName, stage).First(&task).Error; err != nil {
		return false, "", err
	}
	now := time.Now()
	_ = r.db.Model(&models.DownloadTask{}).Where("id = ?", task.ID).
		Updates(map[string]interface{}{"status": models.TaskDepleted, "reaped_at": &now}).Error

	var activeCount int64
	_ = r.db.Model(&models.DownloadTask{}).
		Where("media_id = ? AND stage = ? AND status IN ?", task.MediaID, stage, []string{string(models.TaskPending), string(models.TaskRunning)}).
		Count(&activeCount).Error

	return activeCount == 0, task.MediaID, nil
}

// UpdateMediaCheckpointTime は メディアの各チェックポイント通過時刻を記録します
func (r *Repository) UpdateMediaCheckpointTime(mediaID string, stage models.PipelineStage) error {
	if r.db == nil || mediaID == "" { return nil }
	now := time.Now()
	field := "requests_at"
	switch stage {
	case models.StageMotrix:
		field = "motrix_at"
	case models.StageThunder:
		field = "thunder_at"
	case models.StageStash:
		field = "stash_at"
	}
	return r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Update(field, &now).Error
}
