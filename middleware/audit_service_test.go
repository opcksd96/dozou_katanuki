package middleware

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dozou_katanuki/driver"
)

func TestAuditService(t *testing.T) {
	testDBPath := filepath.Join(os.TempDir(), "test_audit_svc_"+time.Now().Format("20060102150405")+".db")
	defer os.Remove(testDBPath)

	db, err := driver.InitDB(testDBPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	repo := driver.NewRepository(db)

	tempDir := t.TempDir()
	stashDir := filepath.Join(tempDir, "stash")
	blobsDir := filepath.Join(tempDir, "blobs")
	_ = os.MkdirAll(filepath.Join(stashDir, "scenes"), 0755)
	_ = os.MkdirAll(blobsDir, 0755)

	orphanFile := filepath.Join(stashDir, "scenes", "orphan_test.mp4")
	_ = os.WriteFile(orphanFile, []byte("test orphan content"), 0644)

	var emittedEvents []string
	emitter := func(event string, data ...interface{}) {
		emittedEvents = append(emittedEvents, event)
	}

	svc := NewAuditService(repo, emitter)

	// 1. パージなし監査
	report, err := svc.RunAudit(context.Background(), stashDir, blobsDir, false, false)
	if err != nil {
		t.Fatalf("RunAudit failed: %v", err)
	}
	if !report.IntegrityOK {
		t.Errorf("expected IntegrityOK true")
	}
	if len(report.OrphanFiles) != 1 {
		t.Errorf("expected 1 orphan file, got %d", len(report.OrphanFiles))
	}
	if report.PurgedFileCount != 0 {
		t.Errorf("expected 0 purged files, got %d", report.PurgedFileCount)
	}

	// 2. パージあり監査
	report2, err := svc.RunAudit(context.Background(), stashDir, blobsDir, true, false)
	if err != nil {
		t.Fatalf("RunAudit with purge failed: %v", err)
	}
	if report2.PurgedFileCount != 1 {
		t.Errorf("expected 1 purged file, got %d", report2.PurgedFileCount)
	}

	// イベントが発火しているか
	if len(emittedEvents) < 2 {
		t.Errorf("expected at least 2 events emitted, got %d", len(emittedEvents))
	}
}
