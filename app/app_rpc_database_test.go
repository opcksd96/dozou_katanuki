// app/app_rpc_database_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"database/sql"
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

func TestUpdateMediaMetadata_QueuedRebuild(t *testing.T) {
	db, _ := setupTestDB(t)
	_ = db.AutoMigrate(&models.DownloadTask{})
	repo := driver.NewRepository(db)
	timeline := middleware.NewTimelineService(repo)
	app := &App{Repo: repo, TimelineService: timeline, Ready: make(chan struct{})}
	close(app.Ready)

	m := models.Media{
		MediaID:        "media_test_requeue",
		ArticleID:      "art_001",
		Type:           "image",
		DownloadURL:    "https://example.com/image.jpg",
		DownloadStatus: "FAILED",
		FailedReason:   sql.NullString{String: "Motrixエラー", Valid: true},
		StashImageID:   sql.NullString{String: "2735", Valid: true},
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create media: %v", err)
	}

	// UpdateMediaMetadata で QUEUED へ更新
	err := app.UpdateMediaMetadata("media_test_requeue", "QUEUED", "", "2735", "Motrixエラー")
	if err != nil {
		t.Fatalf("UpdateMediaMetadata error: %v", err)
	}

	// メディアの検証
	var updated models.Media
	if err := db.Where("media_id = ?", "media_test_requeue").First(&updated).Error; err != nil {
		t.Fatalf("failed to fetch updated media: %v", err)
	}
	if updated.DownloadStatus != "QUEUED" {
		t.Errorf("expected QUEUED, got %s", updated.DownloadStatus)
	}
	if updated.FailedReason.Valid && updated.FailedReason.String != "" {
		t.Errorf("expected failed_reason to be cleared, got %s", updated.FailedReason.String)
	}
	if !updated.StashImageID.Valid || updated.StashImageID.String != "2735" {
		t.Errorf("expected stash_image_id 2735, got %v", updated.StashImageID)
	}

	// download_tasks が生成されていること
	var taskCount int64
	_ = db.Model(&models.DownloadTask{}).Where("media_id = ?", "media_test_requeue").Count(&taskCount).Error
	if taskCount != 1 {
		t.Errorf("expected 1 download_task, got %d", taskCount)
	}
}
