// app/app_rpc_stash_reconcile_media_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func TestReconcileStashMedia(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := driver.NewRepository(db)
	app := &App{Repo: repo}

	// Stashがオフラインの場合でもパニックせずエラーまたは0を返すことを検証
	count, err := app.ReconcileStashMedia()
	t.Logf("ReconcileStashMedia count: %d, err: %v", count, err)
	if count != 0 {
		t.Fatalf("expected 0 for offline stash, got %d", count)
	}
}

func TestCoordinateTaskDepletionEscalate(t *testing.T) {
	db, _ := setupTestDB(t)
	_ = db.AutoMigrate(&models.DownloadTask{})
	repo := driver.NewRepository(db)
	app := &App{Repo: repo}

	m := models.Media{
		MediaID:        "media_deplete_test",
		ArticleID:      "art_deplete",
		DownloadURL:    "https://example.com/test.mp4",
		DownloadStatus: "OUTSOURCED",
	}
	_ = db.Create(&m)

	app.CoordinateTaskDepletion("media_deplete_test", models.StageMotrix)

	var updated models.Media
	_ = db.Where("media_id = ?", "media_deplete_test").First(&updated)
	if updated.DownloadStatus != "ESCALATED" {
		t.Fatalf("expected status ESCALATED, got %s", updated.DownloadStatus)
	}
}
