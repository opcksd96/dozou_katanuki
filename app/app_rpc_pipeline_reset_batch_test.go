// app/app_rpc_pipeline_reset_batch_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func TestResetSpecificMediasToQueued(t *testing.T) {
	db, _ := setupTestDB(t)
	_ = db.AutoMigrate(&models.DownloadTask{})
	repo := driver.NewRepository(db)
	app := &App{Repo: repo}

	// COMPLETED なメディアを1件作成
	m := models.Media{
		MediaID:        "media_test_001",
		ArticleID:      "art_001",
		Type:           "video",
		DownloadURL:    "https://example.com/video.mp4",
		DownloadStatus: "COMPLETED",
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create media: %v", err)
	}

	// ResetSpecificMediasToQueued を実行
	count, err := app.ResetSpecificMediasToQueued([]string{"media_test_001"})
	if err != nil {
		t.Fatalf("ResetSpecificMediasToQueued error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}

	var updated models.Media
	if err := db.Where("media_id = ?", "media_test_001").First(&updated).Error; err != nil {
		t.Fatalf("failed to query media: %v", err)
	}
	if updated.DownloadStatus != "QUEUED" {
		t.Fatalf("expected status QUEUED, got %s", updated.DownloadStatus)
	}
	if updated.FailedReason.Valid {
		t.Fatalf("expected failed_reason to be nil/invalid, got %v", updated.FailedReason)
	}

	// 候補タスクが download_tasks に登録されていること
	var taskCount int64
	_ = db.Model(&models.DownloadTask{}).Where("media_id = ?", "media_test_001").Count(&taskCount).Error
	if taskCount == 0 {
		t.Errorf("expected download_tasks to be generated, got %d", taskCount)
	}
}

func TestExpandMediaCandidateTasks_Comprehensive(t *testing.T) {
	m := models.Media{
		MediaID:     "123456789_0",
		ArticleID:   "art_999",
		Type:        "video",
		DownloadURL: "https://video.twimg.com/ext_tw_video/123456789/pu/vid/1280x720/test.mp4",
	}
	tasks := ExpandMediaCandidateTasks(m, models.StageRequests)
	if len(tasks) <= 1 {
		t.Errorf("expected multiple candidate tasks, got %d", len(tasks))
	}
	foundWayback := false
	for _, task := range tasks {
		if task.URL != "" && len(task.URL) > 0 {
			if task.MediaID != "123456789_0" {
				t.Errorf("unexpected media_id in task: %s", task.MediaID)
			}
		}
	}
	_ = foundWayback
}
