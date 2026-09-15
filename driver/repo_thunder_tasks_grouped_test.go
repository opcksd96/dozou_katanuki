// driver/repo_thunder_tasks_grouped_test.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"testing"

	"dozou_katanuki/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestFetchMediaWithCandidateTasks(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil { t.Fatalf("failed to open in-memory db: %v", err) }
	_ = db.AutoMigrate(&models.Media{}, &models.ThunderTask{})
	repo := NewRepository(db)

	m1 := models.Media{MediaID: "m1", ArticleID: "a1", DownloadURL: "https://example.com/pic1.jpg", DownloadStatus: "ESCALATED"}
	m2 := models.Media{MediaID: "m2", ArticleID: "a2", DownloadURL: "https://example.com/pic2.jpg", DownloadStatus: "ESCALATED"}
	_ = db.Create(&m1).Error
	_ = db.Create(&m2).Error

	t1 := models.ThunderTask{ID: "task1", MediaID: "m1", URL: "https://example.com/pic1_orig.jpg", ResolutionType: "orig", Status: models.ThunderTaskRunning}
	t2 := models.ThunderTask{ID: "task2", MediaID: "m1", URL: "https://example.com/pic1_large.jpg", ResolutionType: "large", Status: models.ThunderTaskPending}
	t3 := models.ThunderTask{ID: "task3", MediaID: "m2", URL: "https://example.com/pic2_orig.jpg", ResolutionType: "orig", Status: models.ThunderTaskPending}
	_ = repo.BatchUpsertThunderTasks([]models.ThunderTask{t1, t2, t3})

	items, err := repo.FetchMediaWithCandidateTasks("ESCALATED", 10)
	if err != nil { t.Fatalf("FetchMediaWithCandidateTasks failed: %v", err) }
	if len(items) != 2 { t.Fatalf("expected 2 media items, got %d", len(items)) }

	if items[0].MediaID != "m1" { t.Errorf("expected first media m1, got %s", items[0].MediaID) }
	if len(items[0].CandidateTasks) != 2 { t.Errorf("expected m1 to have 2 tasks, got %d", len(items[0].CandidateTasks)) }
	if items[1].MediaID != "m2" { t.Errorf("expected second media m2, got %s", items[1].MediaID) }
	if len(items[1].CandidateTasks) != 1 { t.Errorf("expected m2 to have 1 task, got %d", len(items[1].CandidateTasks)) }
}
