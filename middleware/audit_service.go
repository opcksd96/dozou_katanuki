package middleware

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

// AuditService は SPEC-AUDIT-001 に基づくデータベースおよびストレージの整合性監査を統括するサービスです
type AuditService struct {
	repo    *driver.Repository
	emitter func(event string, data ...interface{})
	mu      sync.Mutex
}

func NewAuditService(repo *driver.Repository, emitter func(event string, data ...interface{})) *AuditService {
	return &AuditService{
		repo:    repo,
		emitter: emitter,
	}
}

// RunAudit は SQLite3 整合性監査および孤立メディア・ファイルを走査し、オプションで自動パージを実行します
func (s *AuditService) RunAudit(ctx context.Context, stashDir, blobsDir string, autoPurgeFiles, autoPurgeDB bool) (*models.AuditReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("[AuditService] Starting comprehensive audit (stashDir=%s, blobsDir=%s, purgeFiles=%v, purgeDB=%v)",
		stashDir, blobsDir, autoPurgeFiles, autoPurgeDB)

	if s.emitter != nil {
		s.emitter("audit:started", map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	report, err := s.repo.AuditDatabase(stashDir, blobsDir)
	if err != nil {
		log.Printf("[AuditService] Audit failed: %v", err)
		if s.emitter != nil {
			s.emitter("audit:error", err.Error())
		}
		return nil, fmt.Errorf("audit execution failed: %w", err)
	}

	// 自動パージ: 孤立ファイルの退避
	if autoPurgeFiles && len(report.OrphanFiles) > 0 {
		var paths []string
		for _, f := range report.OrphanFiles {
			paths = append(paths, f.Path)
		}
		purgedCount, purgeErr := s.repo.PurgeOrphanFiles(paths)
		if purgeErr != nil {
			log.Printf("[AuditService] Warning: failed to purge some orphan files: %v", purgeErr)
		}
		report.PurgedFileCount = purgedCount
		log.Printf("[AuditService] Successfully moved %d orphan files to Recycle Bin", purgedCount)
	}

	// 自動パージ: 孤立 DB メディアの削除
	if autoPurgeDB && len(report.OrphanDBMedia) > 0 {
		var ids []string
		for _, m := range report.OrphanDBMedia {
			ids = append(ids, m.MediaID)
		}
		purgedDBCount, dbErr := s.repo.PurgeOrphanDBMedia(ids)
		if dbErr != nil {
			log.Printf("[AuditService] Warning: failed to delete some orphan DB records: %v", dbErr)
		}
		report.PurgedDBMediaCount = purgedDBCount
		log.Printf("[AuditService] Successfully deleted %d orphan DB media records", purgedDBCount)
	}

	if s.emitter != nil {
		s.emitter("audit:finished", report)
	}

	log.Printf("[AuditService] Audit finished successfully: %s (DB orphans: %d, File orphans: %d, Purged files: %d)",
		report.Summary, len(report.OrphanDBMedia), len(report.OrphanFiles), report.PurgedFileCount)

	return report, nil
}

// PurgeOrphanFiles は指定されたファイルパスのリストをOSのごみ箱へ退避します
func (s *AuditService) PurgeOrphanFiles(paths []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.PurgeOrphanFiles(paths)
}

// PurgeOrphanDBMedia は指定された media_id の DB レコードを物理削除します
func (s *AuditService) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.PurgeOrphanDBMedia(mediaIDs)
}
