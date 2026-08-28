// app/app_rpc_media_trash_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

func TestMediaTrashAndRestore(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := driver.NewRepository(db)
	timeline := middleware.NewTimelineService(repo)
	a := &App{Repo: repo, TimelineService: timeline, Ready: make(chan struct{})}
	close(a.Ready)

	// 1. テストデータの挿入
	m := models.Media{
		MediaID:        "media_test_001",
		ArticleID:      "article_001",
		Type:           "image",
		DownloadURL:    "https://example.com/test.jpg",
		DownloadStatus: "COMPLETED",
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create test media: %v", err)
	}

	// 2. ゴミ箱へ移動
	if err := a.TrashMedia("media_test_001", "不要なテストメディア"); err != nil {
		t.Fatalf("TrashMedia failed: %v", err)
	}

	// 3. ゴミ箱状態の検証
	var saved models.Media
	if err := db.Where("media_id = ?", "media_test_001").First(&saved).Error; err != nil {
		t.Fatalf("failed to query media: %v", err)
	}
	if !saved.IsTrash {
		t.Errorf("expected IsTrash=true, got %v", saved.IsTrash)
	}
	if saved.TrashReason != "不要なテストメディア" {
		t.Errorf("expected reason='不要なテストメディア', got %s", saved.TrashReason)
	}

	// 4. 復元
	if err := a.RestoreMedia("media_test_001"); err != nil {
		t.Fatalf("RestoreMedia failed: %v", err)
	}

	// 5. 復元状態の検証
	if err := db.Where("media_id = ?", "media_test_001").First(&saved).Error; err != nil {
		t.Fatalf("failed to query restored media: %v", err)
	}
	if saved.IsTrash {
		t.Errorf("expected IsTrash=false, got %v", saved.IsTrash)
	}
	if saved.TrashReason != "" {
		t.Errorf("expected empty reason, got %s", saved.TrashReason)
	}
}
