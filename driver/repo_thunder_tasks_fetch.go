// driver/repo_thunder_tasks_fetch.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"dozou_katanuki/models"
)

var dupParenReg = regexp.MustCompile(`\s*\(\d+\)(\.[a-zA-Z0-9]+)$`)

// IsManagedThunderTask は 指定ファイル名が thunder_tasks DB に登録されている管理対象かを返します
func (r *Repository) IsManagedThunderTask(fileName string) bool {
	if r.db == nil || fileName == "" { return false }
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".mp4" {
		return false // Twitter メディア以外の拡張子 (zip等) は絶対に dozou 管理対象ではない
	}
	clean := strings.TrimSpace(dupParenReg.ReplaceAllString(fileName, "$1"))
	var count int64
	_ = r.db.Model(&models.ThunderTask{}).
		Where("file_name = ? OR file_name = ?", fileName, clean).
		Count(&count).Error
	return count > 0
}

// FetchPendingThunderTasks は PENDING 状態のタスクを指定件数取得します
func (r *Repository) FetchPendingThunderTasks(limit int) ([]models.ThunderTask, error) {
	var tasks []models.ThunderTask
	if r.db == nil { return tasks, nil }
	if limit <= 0 { limit = 10 }

	err := r.db.Where("status = ?", models.ThunderTaskPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// MarkThunderTaskReaped は タスクを REAPED（0B/0B, status=9, err=402）に更新し、全滅(ALL_FAILED)かを返します
func (r *Repository) MarkThunderTaskReaped(fileName, reason string) (allReaped bool, mediaID string, err error) {
	if r.db == nil || fileName == "" { return false, "", nil }
	var task models.ThunderTask
	if err := r.db.Where("file_name = ?", fileName).First(&task).Error; err != nil {
		return false, "", err
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":           models.ThunderTaskReaped,
		"reaped_at":        &now,
		"file_size":        0,
		"download_size":    0,
		"task_status_code": 9,
		"error_code":       402,
		"detail_text":      reason,
		"error_reason":     reason,
	}
	_ = r.db.Model(&models.ThunderTask{}).Where("id = ?", task.ID).Updates(updates).Error

	var compCount int64
	_ = r.db.Model(&models.ThunderTask{}).Where("media_id = ? AND status = ?", task.MediaID, models.ThunderTaskCompleted).Count(&compCount).Error
	if compCount > 0 { return false, task.MediaID, nil }

	var activeCount int64
	_ = r.db.Model(&models.ThunderTask{}).Where("media_id = ? AND status IN ?", task.MediaID, []string{
		string(models.ThunderTaskPending), string(models.ThunderTaskOnboarded),
		string(models.ThunderTaskRunning), string(models.ThunderTaskHolding),
	}).Count(&activeCount).Error

	return activeCount == 0, task.MediaID, nil
}

// ReactivatePausedThunderTasks は タイムアウトや保留中のタスクを PENDING に差し戻して再活性化します
func (r *Repository) ReactivatePausedThunderTasks() (int64, error) {
	if r.db == nil { return 0, nil }
	now := time.Now()
	res := r.db.Model(&models.ThunderTask{}).
		Where("status IN ? AND (error_reason LIKE ? OR error_reason LIKE ?)", 
			[]string{string(models.ThunderTaskHolding), string(models.ThunderTaskRetired)},
			"%超时%", "%timeout%").
		Updates(map[string]interface{}{
			"status":     models.ThunderTaskPending,
			"updated_at": now,
		})
	return res.RowsAffected, res.Error
}
