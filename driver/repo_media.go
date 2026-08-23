// driver/repo_media.go (100行以下)
package driver

import (
	"fmt"

	"dozou_katanuki/models"
)

func (r *Repository) ResetMediaStatus(mediaID string) error {
	return r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Updates(map[string]interface{}{
		"download_status": "QUEUED", "failed_reason": nil,
	}).Error
}

func (r *Repository) UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason string) error {
	updates := map[string]interface{}{
		"download_status": downloadStatus, "failed_reason": nil,
		"stash_scene_id": nil, "stash_image_id": nil,
	}
	if failedReason != "" { updates["failed_reason"] = failedReason }
	if stashSceneID != "" { updates["stash_scene_id"] = stashSceneID }
	if stashImageID != "" { updates["stash_image_id"] = stashImageID }
	return r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Updates(updates).Error
}

func (r *Repository) GetMediaByID(mediaID string) (*models.Media, error) {
	var m models.Media
	err := r.db.Where("media_id = ?", mediaID).First(&m).Error
	return &m, err
}

func (r *Repository) PurgeMedia(mediaID string) error {
	return r.db.Where("media_id = ?", mediaID).Delete(&models.Media{}).Error
}

func (r *Repository) PurgeMediaByStatus(status, accountID string) (int64, error) {
	q := r.db.Table("media")
	if accountID != "" && accountID != "all" {
		q = q.Where("article_id IN (SELECT id FROM articles WHERE account_id = ? OR account_id IN (SELECT numeric_id FROM accounts WHERE username = ?))", accountID, accountID)
	}

	if status == "EXCLUDED" {
		q = q.Where("download_status = 'EXCLUDED' OR failed_reason LIKE '%Whitelist外%' OR failed_reason LIKE '%ダウンロード対象外%'")
	} else if status == "UNLINKED" {
		q = q.Where("download_status = 'COMPLETED' AND (stash_scene_id IS NULL OR stash_scene_id = '') AND (stash_image_id IS NULL OR stash_image_id = '')")
	} else if status == "DEAD_404" {
		q = q.Where("download_status = 'DEAD_404'")
	} else if status != "" && status != "all" {
		q = q.Where("download_status = ?", status)
	} else {
		return 0, fmt.Errorf("status is required for batch purge")
	}

	res := q.Delete(&models.Media{})
	return res.RowsAffected, res.Error
}

func (r *Repository) MigrateExcludedMedia() (int64, error) {
	res := r.db.Model(&models.Media{}).
		Where("download_status = 'DEAD_404' AND (failed_reason LIKE '%Whitelist外%' OR failed_reason LIKE '%ダウンロード対象外%')").
		Update("download_status", "EXCLUDED")
	return res.RowsAffected, res.Error
}
