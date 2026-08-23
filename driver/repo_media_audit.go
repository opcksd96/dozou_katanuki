// driver/repo_media_audit.go (100行以下)
package driver

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

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
	if !rep.IntegrityOK {
		summary = "要修復 (SQLite ページ/インデックス破損検知)"
	} else if !rep.ForeignKeyOK {
		summary = fmt.Sprintf("注意 (外部キー違反: %d件)", len(rep.ForeignKeyErrors))
	} else if len(rep.OrphanDBMedia) > 0 || len(rep.OrphanFiles) > 0 {
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
