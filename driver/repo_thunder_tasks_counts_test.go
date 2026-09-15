// driver/repo_thunder_tasks_counts_test.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"testing"

	"dozou_katanuki/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestGetThunderTaskCounts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil { t.Fatalf("failed to open in-memory db: %v", err) }
	_ = db.AutoMigrate(&models.ThunderTask{})
	repo := NewRepository(db)

	tasks := []models.ThunderTask{
		{ID: "t1", Status: models.ThunderTaskPending},
		{ID: "t2", Status: models.ThunderTaskPending},
		{ID: "t3", Status: models.ThunderTaskRunning},
		{ID: "t4", Status: models.ThunderTaskOnboarded},
		{ID: "t5", Status: models.ThunderTaskHolding},
		{ID: "t6", Status: models.ThunderTaskCompleted},
		{ID: "t7", Status: models.ThunderTaskRetired},
		{ID: "t8", Status: models.ThunderTaskReaped},
	}
	_ = repo.BatchUpsertThunderTasks(tasks)

	counts, err := repo.GetThunderTaskCounts()
	if err != nil { t.Fatalf("GetThunderTaskCounts error: %v", err) }

	if counts.Total != 8 { t.Errorf("expected total=8, got %d", counts.Total) }
	if counts.Pending != 2 { t.Errorf("expected pending=2, got %d", counts.Pending) }
	if counts.Running != 2 { t.Errorf("expected running=2 (RUNNING+ONBOARDED), got %d", counts.Running) }
	if counts.Holding != 1 { t.Errorf("expected holding=1, got %d", counts.Holding) }
	if counts.Completed != 1 { t.Errorf("expected completed=1, got %d", counts.Completed) }
	if counts.Failed != 2 { t.Errorf("expected failed=2 (RETIRED+REAPED), got %d", counts.Failed) }
}
