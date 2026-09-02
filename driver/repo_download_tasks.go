// driver/repo_download_tasks.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm/clause"
)

// BatchUpsertDownloadTasks は メディアごとの単一ダウンロードタスクを一括登録・更新します
func (r *Repository) BatchUpsertDownloadTasks(tasks []models.DownloadTask) error {
	if len(tasks) == 0 || r.db == nil { return nil }
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "media_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"article_id", "url", "file_name", "stage", "status", "updated_at"}),
	}).Create(&tasks).Error
}

// MarkTaskCompleted は メディアのダウンロードタスクを完了状態に更新します
func (r *Repository) MarkTaskCompleted(mediaID string) error {
	if r.db == nil || mediaID == "" { return nil }
	now := time.Now()
	return r.db.Model(&models.DownloadTask{}).Where("media_id = ?", mediaID).
		Updates(map[string]interface{}{"status": models.TaskCompleted, "completed_at": &now, "updated_at": now}).Error
}

// UpdateMediaCheckpointTime は 1メディア1行のタスクレコードにステージ通過時刻を記録します
func (r *Repository) UpdateMediaCheckpointTime(mediaID string, stage models.PipelineStage) error {
	if r.db == nil || mediaID == "" { return nil }
	now := time.Now()
	updates := map[string]interface{}{"stage": stage, "status": models.TaskRunning, "updated_at": now}
	switch stage {
	case models.StageRequests:
		updates["requests_at"] = &now
	case models.StageMotrix:
		updates["motrix_at"] = &now
	case models.StageThunder:
		updates["thunder_at"] = &now
	case models.StageStash:
		updates["stash_at"] = &now
	}
	return r.db.Model(&models.DownloadTask{}).Where("media_id = ?", mediaID).Updates(updates).Error
}
