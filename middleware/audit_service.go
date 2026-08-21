// middleware/audit_service.go (100行以下)
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

type AuditService struct {
	repo    *driver.Repository
	emitter func(event string, data ...interface{})
	mu      sync.Mutex
}

func NewAuditService(repo *driver.Repository, emitter func(event string, data ...interface{})) *AuditService {
	return &AuditService{repo: repo, emitter: emitter}
}

func (s *AuditService) RunAudit(ctx context.Context, stashDir, blobsDir string, autoPurgeFiles, autoPurgeDB bool) (*models.AuditReport, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	log.Printf("[AuditService] Starting audit: stash=%s, blobs=%s, purgeFiles=%v, purgeDB=%v", stashDir, blobsDir, autoPurgeFiles, autoPurgeDB)
	if s.emitter != nil { s.emitter("audit:started", map[string]interface{}{"timestamp": time.Now().Format(time.RFC3339)}) }

	report, err := s.repo.AuditDatabase(stashDir, blobsDir)
	if err != nil {
		if s.emitter != nil { s.emitter("audit:error", err.Error()) }
		return nil, fmt.Errorf("audit failed: %w", err)
	}

	if autoPurgeFiles && len(report.OrphanFiles) > 0 {
		var paths []string
		for _, f := range report.OrphanFiles { paths = append(paths, f.Path) }
		purgedCount, purgeErr := s.repo.PurgeOrphanFiles(paths)
		if purgeErr != nil { log.Printf("[AuditService] Purge files warning: %v", purgeErr) }
		report.PurgedFileCount = purgedCount
	}

	if autoPurgeDB && len(report.OrphanDBMedia) > 0 {
		var ids []string
		for _, m := range report.OrphanDBMedia { ids = append(ids, m.MediaID) }
		purgedDBCount, dbErr := s.repo.PurgeOrphanDBMedia(ids)
		if dbErr != nil { log.Printf("[AuditService] Purge DB records warning: %v", dbErr) }
		report.PurgedDBMediaCount = purgedDBCount
	}

	if s.emitter != nil { s.emitter("audit:finished", report) }
	log.Printf("[AuditService] Audit finished: %s (DB orphans: %d, File orphans: %d)", report.Summary, len(report.OrphanDBMedia), len(report.OrphanFiles))
	return report, nil
}

func (s *AuditService) PurgeOrphanFiles(paths []string) (int, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	return s.repo.PurgeOrphanFiles(paths)
}

func (s *AuditService) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	return s.repo.PurgeOrphanDBMedia(mediaIDs)
}
