// driver/repo_thunder_tasks_status_test.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"testing"

	"dozou_katanuki/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestMarkThunderTaskRetiredAndCheckAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	_ = db.AutoMigrate(&models.ThunderTask{}, &models.Media{})
	repo := NewRepository(db)

	mediaID := "test-media-123"
	tasks := []models.ThunderTask{
		{ID: "t-orig", MediaID: mediaID, FileName: "test_orig.jpg", ResolutionType: "orig", Status: models.ThunderTaskOnboarded},
		{ID: "t-large", MediaID: mediaID, FileName: "test_large.jpg", ResolutionType: "large", Status: models.ThunderTaskOnboarded},
		{ID: "t-thumb", MediaID: mediaID, FileName: "test_thumb.jpg", ResolutionType: "thumb", Status: models.ThunderTaskPending},
	}
	if err := repo.BatchUpsertThunderTasks(tasks); err != nil {
		t.Fatalf("failed to insert tasks: %v", err)
	}

	// 1. 1つ目をRETIRED ➔ まだ残りがあるため allRetired は false
	all1, mID1, err := repo.MarkThunderTaskRetiredAndCheckAll("test_orig.jpg", "404")
	if err != nil || all1 || mID1 != mediaID {
		t.Errorf("expected allRetired=false for 1st task, got %v, err=%v", all1, err)
	}

	// 2. 2つ目をRETIRED ➔ まだ thumb(PENDING) があるため false
	all2, _, _ := repo.MarkThunderTaskRetiredAndCheckAll("test_large.jpg", "404")
	if all2 {
		t.Errorf("expected allRetired=false for 2nd task, got true")
	}

	// 3. 最後の thumb も RETIRED ➔ 全候補全滅 (ALL_TRUE) で true になる！
	all3, mID3, _ := repo.MarkThunderTaskRetiredAndCheckAll("test_thumb.jpg", "404")
	if !all3 || mID3 != mediaID {
		t.Errorf("expected allRetired=true (ALL_TRUE) when all candidates retired, got %v", all3)
	}
}
