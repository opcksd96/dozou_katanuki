// driver/repo_thunder_tasks_status.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"time"

	"dozou_katanuki/models"
)

// MarkThunderTaskOnboarded は タスクが迅雷へ投入されたことをマークします
func (r *Repository) MarkThunderTaskOnboarded(taskID, summarySize string) error {
	if r.db == nil || taskID == "" { return nil }
	now := time.Now()
	updates := map[string]interface{}{"status": models.ThunderTaskOnboarded, "dispatched_at": &now, "last_attempt_at": &now}
	if summarySize != "" { updates["summary_size"] = summarySize }
	return r.db.Model(&models.ThunderTask{}).Where("id = ?", taskID).Updates(updates).Error
}

// MarkThunderTaskHolding は サマリ>1Bの探索継続タスクをHOLDINGに更新します
func (r *Repository) MarkThunderTaskHolding(fileName, summarySize, reason string) error {
	if r.db == nil || fileName == "" { return nil }
	now := time.Now()
	return r.db.Model(&models.ThunderTask{}).Where("file_name = ?", fileName).
		Updates(map[string]interface{}{
			"status": models.ThunderTaskHolding, "summary_size": summarySize, "error_reason": reason, "last_attempt_at": &now,
		}).Error
}

// MarkThunderTaskRetiredAndCheckAll は 候補タスクをRETIREDにし、全候補が全滅(ALL_TRUE)したかを判定します
func (r *Repository) MarkThunderTaskRetiredAndCheckAll(fileName, reason string) (allRetired bool, mediaID string, err error) {
	if r.db == nil || fileName == "" { return false, "", nil }
	var task models.ThunderTask
	if err := r.db.Where("file_name = ?", fileName).First(&task).Error; err != nil {
		return false, "", err
	}

	now := time.Now()
	_ = r.db.Model(&models.ThunderTask{}).Where("id = ?", task.ID).
		Updates(map[string]interface{}{
			"status": models.ThunderTaskRetired, "error_reason": reason, "reaped_at": &now,
		}).Error

	// 同一 media_id の全タスクを検査: PENDING / ONBOARDED / RUNNING / HOLDING が1つでもあれば false
	var activeCount int64
	_ = r.db.Model(&models.ThunderTask{}).
		Where("media_id = ? AND status IN ?", task.MediaID, []string{
			string(models.ThunderTaskPending), string(models.ThunderTaskOnboarded),
			string(models.ThunderTaskRunning), string(models.ThunderTaskHolding),
		}).Count(&activeCount).Error

	// すべてのダウンロード候補が RETIRED (または COMPLETED/REAPED 以外で全滅) の時のみ ALL_TRUE
	return activeCount == 0, task.MediaID, nil
}

// UpdateThunderTaskCooldown は 429等の異常時にクールダウン開始時刻を記録します
func (r *Repository) UpdateThunderTaskCooldown(fileName, errReason string) error {
	if r.db == nil || fileName == "" { return nil }
	now := time.Now()
	return r.db.Model(&models.ThunderTask{}).Where("file_name = ?", fileName).
		Updates(map[string]interface{}{"error_reason": errReason, "last_attempt_at": &now}).Error
}
