// driver/repo_thunder_tasks.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm/clause"
)

// BatchUpsertThunderTasks は 候補タスクを一括登録します
func (r *Repository) BatchUpsertThunderTasks(tasks []models.ThunderTask) error {
	if len(tasks) == 0 || r.db == nil { return nil }
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"url", "file_name", "status", "updated_at"}),
	}).Create(&tasks).Error
}

// MarkThunderTaskRunning は タスクの投入時刻とRUNNINGステータスを更新します
func (r *Repository) MarkThunderTaskRunning(taskID string) error {
	if r.db == nil { return nil }
	now := time.Now()
	return r.db.Model(&models.ThunderTask{}).Where("id = ?", taskID).
		Updates(map[string]interface{}{"status": models.ThunderTaskRunning, "dispatched_at": &now}).Error
}

// MarkThunderTaskCompletedAndReapOthers は 勝利タスクを完了にし、他候補を取り下げ対象として取得します
func (r *Repository) MarkThunderTaskCompletedAndReapOthers(mediaID, completedFileName string) ([]models.ThunderTask, error) {
	var reaps []models.ThunderTask
	if r.db == nil || mediaID == "" { return reaps, nil }
	now := time.Now()

	// 完了したタスクを特定・更新
	_ = r.db.Model(&models.ThunderTask{}).Where("media_id = ? AND file_name = ?", mediaID, completedFileName).
		Updates(map[string]interface{}{"status": models.ThunderTaskCompleted, "completed_at": &now}).Error

	// 他の未完了タスクを取得して REAPED へ更新
	_ = r.db.Where("media_id = ? AND file_name != ? AND status IN ?", mediaID, completedFileName, []string{string(models.ThunderTaskPending), string(models.ThunderTaskRunning)}).
		Find(&reaps).Error

	if len(reaps) > 0 {
		_ = r.db.Model(&models.ThunderTask{}).Where("media_id = ? AND file_name != ?", mediaID, completedFileName).
			Updates(map[string]interface{}{"status": models.ThunderTaskReaped, "reaped_at": &now}).Error
	}
	return reaps, nil
}

// MarkThunderTaskDepleted は タスクを枯渇済みにし、全候補が全滅したかを返します
func (r *Repository) MarkThunderTaskDepleted(fileName string) (allDepleted bool, mediaID string, err error) {
	if r.db == nil || fileName == "" { return false, "", nil }
	var task models.ThunderTask
	if err := r.db.Where("file_name = ?", fileName).First(&task).Error; err != nil {
		return false, "", err
	}
	now := time.Now()
	_ = r.db.Model(&models.ThunderTask{}).Where("id = ?", task.ID).
		Updates(map[string]interface{}{"status": models.ThunderTaskDepleted, "reaped_at": &now}).Error

	// 同一 media_id でまだ生きている (PENDING または RUNNING) タスクがあるか確認
	var activeCount int64
	_ = r.db.Model(&models.ThunderTask{}).
		Where("media_id = ? AND status IN ?", task.MediaID, []string{string(models.ThunderTaskPending), string(models.ThunderTaskRunning)}).
		Count(&activeCount).Error

	return activeCount == 0, task.MediaID, nil
}
