// app_rpc_audit_test.go (100行以下)
package main

import (
	"os"
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

func TestSchedulerRPC(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)

	repo := driver.NewRepository(db)
	timeline := middleware.NewTimelineService(repo)
	readyChan := make(chan struct{})
	close(readyChan)

	app := &App{repo: repo, timelineService: timeline, ready: readyChan}
	orch := middleware.NewJobOrchestrator(t.Context(), func(string, ...interface{}) {})
	defer orch.Close()
	app.jobOrchestrator = orch

	sched := middleware.NewSchedulerService(models.SchedulerConfig{PollIntervalSec: 10, BackupIntervalHours: 24, MaxBackupGenerations: 3}, repo, orch, func(string, ...interface{}) {})
	app.scheduler = sched

	backupPath, err := app.TriggerBackup()
	if err != nil || backupPath == "" { t.Fatalf("TriggerBackup failed: %v", err) }
	defer os.Remove(backupPath)

	job, err := app.TriggerPoll()
	if err != nil || job == nil || job.Type != models.JobTypeMediaPoll { t.Fatalf("TriggerPoll failed: %+v", job) }
}

func TestAuditRPC(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)

	repo := driver.NewRepository(db)
	timeline := middleware.NewTimelineService(repo)
	readyChan := make(chan struct{})
	close(readyChan)

	app := &App{repo: repo, timelineService: timeline, ready: readyChan}
	app.auditService = middleware.NewAuditService(repo, func(string, ...interface{}) {})

	report, err := app.RunAudit(false, false)
	if err != nil || report == nil || !report.IntegrityOK { t.Fatalf("RunAudit failed: %+v", report) }

	purgedFiles, err := app.PurgeOrphanFiles([]string{})
	if err != nil || purgedFiles != 0 { t.Fatalf("PurgeOrphanFiles failed: %v", err) }

	purgedDB, err := app.PurgeOrphanDBMedia([]string{})
	if err != nil || purgedDB != 0 { t.Fatalf("PurgeOrphanDBMedia failed: %v", err) }
}
