// driver/repo_media.go (100行以下)
package driver

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

func (r *Repository) ResetMediaStatus(mediaID string) error {
	return r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Updates(map[string]interface{}{
		"download_status": "QUEUED", "failed_reason": nil,
	}).Error
}

func (r *Repository) GetMediaByID(mediaID string) (*models.Media, error) {
	var m models.Media
	err := r.db.Where("media_id = ?", mediaID).First(&m).Error
	return &m, err
}

func (r *Repository) AuditDatabase(stashDir, blobsDir string) (*models.AuditReport, error) {
	rep := &models.AuditReport{ExecutedAt: time.Now()}
	iOK, iMsgs, err := RunIntegrityCheck(r.db)
	if err != nil { rep.IntegrityOK = false; rep.IntegrityErrors = []string{err.Error()} } else { rep.IntegrityOK = iOK; rep.IntegrityErrors = iMsgs }

	fkV, err := RunForeignKeyCheck(r.db)
	if err != nil { rep.ForeignKeyOK = false } else { rep.ForeignKeyOK = len(fkV) == 0; rep.ForeignKeyErrors = fkV }

	if oMed, err := FindOrphanDBMedia(r.db, stashDir, blobsDir); err == nil { rep.OrphanDBMedia = oMed }
	if keys, err := GetKnownMediaIdentifiers(r.db); err == nil {
		if oFiles, err := ScanOrphanFiles(stashDir, blobsDir, keys); err == nil { rep.OrphanFiles = oFiles }
	}

	summary := "健全"
	if !rep.IntegrityOK { summary = "要修復 (SQLite ページ/インデックス破損検知)" } else if !rep.ForeignKeyOK { summary = fmt.Sprintf("注意 (外部キー違反: %d件)", len(rep.ForeignKeyErrors)) } else if len(rep.OrphanDBMedia) > 0 || len(rep.OrphanFiles) > 0 {
		summary = fmt.Sprintf("要クレンジング (孤立DB:%d件, 孤立ファイル:%d件)", len(rep.OrphanDBMedia), len(rep.OrphanFiles))
	}
	rep.Summary = summary
	return rep, nil
}

func (r *Repository) PurgeOrphanFiles(paths []string) (int, error) {
	cnt := 0
	for _, p := range paths {
		if err := MoveToRecycleBin(p); err == nil { cnt++ }
	}
	return cnt, nil
}

func (r *Repository) PurgeOrphanDBMedia(trashDir string, mediaIDs []string) (int, error) {
	return BackupAndPurgeDBMedia(r.db, trashDir, mediaIDs)
}

func (r *Repository) RollbackLastDBPurge(trashDir string) (int, error) {
	return RollbackLastDBPurge(r.db, trashDir)
}

func (r *Repository) CanRollbackDBPurge(trashDir string) bool {
	return CanRollbackDBPurge(trashDir)
}

// PurgeMedia deletes a single media record from DB
func (r *Repository) PurgeMedia(mediaID string) error {
	return r.db.Where("media_id = ?", mediaID).Delete(&models.Media{}).Error
}

// PurgeMediaByStatus batch purges media records by status (EXCLUDED, UNLINKED, DEAD_404)
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

// MigrateExcludedMedia migrates legacy DEAD_404 records that failed due to Whitelist to EXCLUDED
func (r *Repository) MigrateExcludedMedia() (int64, error) {
	res := r.db.Model(&models.Media{}).
		Where("download_status = 'DEAD_404' AND (failed_reason LIKE '%Whitelist外%' OR failed_reason LIKE '%ダウンロード対象外%')").
		Update("download_status", "EXCLUDED")
	return res.RowsAffected, res.Error
}
