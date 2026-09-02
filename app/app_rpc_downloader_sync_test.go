// app/app_rpc_downloader_sync_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func TestSyncCompletedDownloadsEscalate(t *testing.T) {
	db, _ := setupTestDB(t)
	_ = db.AutoMigrate(&models.DownloadTask{})
	repo := driver.NewRepository(db)
	app := &App{Repo: repo}

	// 拡張子なしの media_id を持つメディアを作成
	m := models.Media{
		MediaID:        "Go4PdTDb0AAmdRV",
		ArticleID:      "art_motrix_test",
		DownloadURL:    "https://pbs.twimg.com/media/Go4PdTDb0AAmdRV.jpg",
		DownloadStatus: "OUTSOURCED",
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create media: %v", err)
	}

	// Motrix側で返るファイル名（拡張子付き）
	fileName := "Go4PdTDb0AAmdRV.jpg"
	cleanID := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	var med models.Media
	err := db.Where("media_id = ? OR media_id = ? OR download_url LIKE ?", fileName, cleanID, "%/"+fileName).First(&med).Error
	if err != nil {
		t.Fatalf("failed to find media by fileName/cleanID: %v", err)
	}
	if med.MediaID != "Go4PdTDb0AAmdRV" {
		t.Fatalf("expected media_id Go4PdTDb0AAmdRV, got %s", med.MediaID)
	}

	// エスカレーション更新とログ
	_ = app.Repo.UpdateMediaMetadata(med.MediaID, "ESCALATED", "", "", "Motrixエラーテスト")
	_ = app.Repo.UpdateMediaCheckpointTime(med.MediaID, models.StageThunder)

	var updated models.Media
	_ = db.Where("media_id = ?", med.MediaID).First(&updated)
	if updated.DownloadStatus != "ESCALATED" {
		t.Fatalf("expected status ESCALATED, got %s", updated.DownloadStatus)
	}
}

func TestSyncCompletedDownloadsFakeFileRejected(t *testing.T) {
	tempDir := t.TempDir()
	db, _ := setupTestDB(t)
	repo := driver.NewRepository(db)

	fakeFile := filepath.Join(tempDir, "fake_image.jpg")
	_ = os.WriteFile(fakeFile, []byte("<!DOCTYPE html><html><body>Blocked</body></html>"), 0644)

	// RegisterCompletedMediaFile は偽ファイルを拒否してエラーを返すこと
	err := repo.RegisterCompletedMediaFile(fakeFile)
	if err == nil {
		t.Fatalf("expected error for fake html file, got nil")
	}
	t.Logf("Fake file rejected as expected: %v", err)
}
