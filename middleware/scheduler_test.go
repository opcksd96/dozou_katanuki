package middleware_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestSchedulerDB(t *testing.T) (*driver.Repository, string) {
	tempFile := filepath.Join(os.TempDir(), "test_sched_"+time.Now().Format("20060102150405.000000")+".db")
	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	_ = db.AutoMigrate(&models.Account{}, &models.Article{}, &models.Media{})
	return driver.NewRepository(db), tempFile
}

func TestSchedulerServiceLifecycleAndTriggers(t *testing.T) {
	repo, tempDB := setupTestSchedulerDB(t)
	defer os.Remove(tempDB)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var emittedEvents []string
	emitter := func(event string, data ...interface{}) {
		emittedEvents = append(emittedEvents, event)
	}

	orch := middleware.NewJobOrchestrator(ctx, emitter)
	defer orch.Close()

	cfg := models.SchedulerConfig{
		PollIntervalSec:      1, // 1秒周期でテスト
		BackupIntervalHours:  1,
		MaxBackupGenerations: 3,
	}

	sched := middleware.NewSchedulerService(cfg, repo, orch, emitter)
	sched.Start(ctx)

	// 手動バックアップトリガーのテスト
	backupPath, err := sched.TriggerBackup()
	if err != nil {
		t.Fatalf("TriggerBackup failed: %v", err)
	}
	if backupPath == "" {
		t.Fatalf("Expected backup path, got empty")
	}
	defer os.Remove(backupPath)

	// 手動ポーリングトリガーのテスト
	job, err := sched.TriggerPoll()
	if err != nil {
		t.Fatalf("TriggerPoll failed: %v", err)
	}
	if job == nil || job.Type != models.JobTypeMediaPoll {
		t.Fatalf("Expected JobTypeMediaPoll, got: %+v", job)
	}

	// 少し待って自動ポーリングのループが稼働することを確認
	time.Sleep(1200 * time.Millisecond)

	sched.Stop()
}
