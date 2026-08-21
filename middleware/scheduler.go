// middleware/scheduler.go (100行以下)
package middleware

import (
	"context"
	"log"
	"sync"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

type SchedulerService struct {
	cfg          models.SchedulerConfig
	repo         *driver.Repository
	orchestrator *JobOrchestrator
	emitter      EventEmitter
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex
	running      bool
}

func NewSchedulerService(cfg models.SchedulerConfig, repo *driver.Repository, orch *JobOrchestrator, emitter EventEmitter) *SchedulerService {
	if cfg.PollIntervalSec <= 0 { cfg.PollIntervalSec = 300 }
	if cfg.BackupIntervalHours <= 0 { cfg.BackupIntervalHours = 24 }
	if cfg.MaxBackupGenerations <= 0 { cfg.MaxBackupGenerations = 7 }
	return &SchedulerService{cfg: cfg, repo: repo, orchestrator: orch, emitter: emitter}
}

func (s *SchedulerService) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running { s.mu.Unlock(); return }
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.running = true
	s.mu.Unlock()

	log.Printf("[Scheduler] Started: Poll=%ds, Backup=%dh, Gen=%d", s.cfg.PollIntervalSec, s.cfg.BackupIntervalHours, s.cfg.MaxBackupGenerations)
	go s.runPollingLoop()
	go s.runBackupLoop()
}

func (s *SchedulerService) Stop() {
	s.mu.Lock(); defer s.mu.Unlock()
	if !s.running { return }
	s.cancel(); s.running = false
	log.Println("[Scheduler] Stopped")
}

func (s *SchedulerService) TriggerPoll() (*models.JobProgress, error) {
	if s.orchestrator == nil { return nil, nil }
	return s.orchestrator.EnqueueMediaPoll("twitter")
}

func (s *SchedulerService) TriggerBackup() (string, error) {
	if s.repo == nil { return "", nil }
	path, err := s.repo.BackupDatabase("backups/database", s.cfg.MaxBackupGenerations)
	if err != nil {
		if s.emitter != nil { s.emitter("scheduler:backup_failed", map[string]string{"error": err.Error()}) }
		return "", err
	}
	if s.emitter != nil { s.emitter("scheduler:backup_completed", map[string]string{"path": path}) }
	return path, nil
}

func (s *SchedulerService) runPollingLoop() {
	ticker := time.NewTicker(time.Duration(s.cfg.PollIntervalSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done(): return
		case <-ticker.C: _, _ = s.TriggerPoll()
		}
	}
}

func (s *SchedulerService) runBackupLoop() {
	ticker := time.NewTicker(time.Duration(s.cfg.BackupIntervalHours) * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done(): return
		case <-ticker.C: _, _ = s.TriggerBackup()
		}
	}
}
