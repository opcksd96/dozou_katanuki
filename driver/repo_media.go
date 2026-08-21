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

func (r *Repository) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	return DeleteDBMediaByIDs(r.db, mediaIDs)
}
