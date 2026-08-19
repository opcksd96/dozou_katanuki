package driver

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dozou_katanuki/models"
)

func TestAuditDatabase(t *testing.T) {
	testDBPath := filepath.Join(os.TempDir(), "test_audit_"+time.Now().Format("20060102150405")+".db")
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	repo := NewRepository(db)

	// テスト用の一時ディレクトリ（Stash / Blobs）
	tempDir := t.TempDir()
	stashDir := filepath.Join(tempDir, "stash")
	blobsDir := filepath.Join(tempDir, "blobs")
	_ = os.MkdirAll(filepath.Join(stashDir, "scenes"), 0755)
	_ = os.MkdirAll(filepath.Join(stashDir, "images"), 0755)
	_ = os.MkdirAll(blobsDir, 0755)

	// 1. 正常なデータの登録
	acc := models.Account{
		NumericID:   "1001",
		Username:    "mashu",
		DisplayName: "Mash Kyrielight",
	}
	art := models.Article{
		ID:        "art_100",
		AccountID: "1001",
		FullText:  "Test post for audit",
		Media: []models.Media{
			{
				MediaID:        "media_valid_1",
				ArticleID:      "art_100",
				Type:           "image",
				DownloadURL:    "https://example.com/img1.jpg",
				DownloadStatus: "COMPLETED",
				StashImageID:   sql.NullString{String: "img_valid_1", Valid: true},
			},
		},
	}
	_ = db.Save(&acc).Error
	_ = repo.UpsertArticleTx(&art)

	// 正当なファイルの配置
	_ = os.WriteFile(filepath.Join(stashDir, "images", "img_valid_1.jpg"), []byte("valid image data"), 0644)
	_ = os.WriteFile(filepath.Join(blobsDir, "media_valid_1.jpg"), []byte("valid blob data"), 0644)

	// 2. 孤立DBレコード（親記事のない media - 外部キー制約を一時解除して過去の破損をシミュレート）
	_ = db.Exec("PRAGMA foreign_keys = OFF;").Error
	orphanMedia := models.Media{
		MediaID:        "media_orphan_db",
		ArticleID:      "art_non_existent",
		Type:           "video",
		DownloadURL:    "https://example.com/orphan.mp4",
		DownloadStatus: "QUEUED",
	}
	if err := db.Save(&orphanMedia).Error; err != nil {
		t.Fatalf("failed to insert orphanMedia: %v", err)
	}
	_ = db.Exec("PRAGMA foreign_keys = ON;").Error

	// 3. 孤立ファイル（DBに存在しないファイル）
	orphanFilePath := filepath.Join(stashDir, "scenes", "scene_orphan_999.mp4")
	_ = os.WriteFile(orphanFilePath, []byte("dummy video content"), 0644)
	orphanBlobPath := filepath.Join(blobsDir, "unknown_blob_888.jpg")
	_ = os.WriteFile(orphanBlobPath, []byte("dummy blob content"), 0644)

	// 監査の実行
	report, err := repo.AuditDatabase(stashDir, blobsDir)
	if err != nil {
		t.Fatalf("AuditDatabase failed: %v", err)
	}

	// 検証
	if !report.IntegrityOK {
		t.Errorf("expected IntegrityOK to be true, got false")
	}

	// 孤立DBメディアの検出確認
	foundDBOrphan := false
	for _, m := range report.OrphanDBMedia {
		if m.MediaID == "media_orphan_db" {
			foundDBOrphan = true
			break
		}
	}
	if !foundDBOrphan {
		t.Errorf("expected to find media_orphan_db in OrphanDBMedia")
	}

	// 孤立ファイルの検出確認
	foundFileOrphan := 0
	for _, f := range report.OrphanFiles {
		if f.FileName == "scene_orphan_999.mp4" || f.FileName == "unknown_blob_888.jpg" {
			foundFileOrphan++
		}
	}
	if foundFileOrphan != 2 {
		t.Errorf("expected 2 orphan files detected, got %d", foundFileOrphan)
	}

	// パージの実行テスト (Orphan Files)
	purgedFiles, err := repo.PurgeOrphanFiles([]string{orphanFilePath, orphanBlobPath})
	if err != nil {
		t.Fatalf("PurgeOrphanFiles failed: %v", err)
	}
	if purgedFiles != 2 {
		t.Errorf("expected 2 purged files, got %d", purgedFiles)
	}

	// パージの実行テスト (Orphan DB)
	purgedDB, err := repo.PurgeOrphanDBMedia([]string{"media_orphan_db"})
	if err != nil {
		t.Fatalf("PurgeOrphanDBMedia failed: %v", err)
	}
	if purgedDB != 1 {
		t.Errorf("expected 1 purged DB record, got %d", purgedDB)
	}
}
