// driver/audit_rollback_test.go (100行以下)
package driver

import (
	"database/sql"
	"path/filepath"
	"testing"

	"dozou_katanuki/models"
)

func TestPurgeAndRollback(t *testing.T) {
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_rollback.db")
	trashDir := filepath.Join(tempDir, "dumps", "_trash")

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil && sqlDB != nil { _ = sqlDB.Close() }
	})
	repo := NewRepository(db)

	acc := models.Account{NumericID: "1001", Username: "mashu", DisplayName: "Mash Kyrielight"}
	art := models.Article{ID: "art_dummy", AccountID: "1001", FullText: "Dummy post for rollback test"}
	if err := db.Save(&acc).Error; err != nil { t.Fatalf("Save account failed: %v", err) }
	if err := db.Save(&art).Error; err != nil { t.Fatalf("Save article failed: %v", err) }

	// 1. テスト用メディアレコードを作成
	m1 := models.Media{
		MediaID:        "media_rb_1",
		ArticleID:      "art_dummy",
		Type:           "image",
		DownloadURL:    "https://example.com/rb1.jpg",
		DownloadStatus: "COMPLETED",
		StashImageID:   sql.NullString{String: "img_rb_1", Valid: true},
	}
	m2 := models.Media{
		MediaID:        "media_rb_2",
		ArticleID:      "art_dummy",
		Type:           "video",
		DownloadURL:    "https://example.com/rb2.mp4",
		DownloadStatus: "COMPLETED",
	}
	if err := db.Save(&m1).Error; err != nil { t.Fatalf("Save m1 failed: %v", err) }
	if err := db.Save(&m2).Error; err != nil { t.Fatalf("Save m2 failed: %v", err) }

	if repo.CanRollbackDBPurge(trashDir) {
		t.Errorf("Expected CanRollback to be false initially")
	}

	// 2. パージ（自動退避付き）
	purged, err := repo.PurgeOrphanDBMedia(trashDir, []string{"media_rb_1", "media_rb_2"})
	if err != nil || purged != 2 {
		t.Fatalf("PurgeOrphanDBMedia failed: count=%d, err=%v", purged, err)
	}

	// レコードが削除されたことを確認
	var count int64
	db.Table("media").Where("media_id IN ?", []string{"media_rb_1", "media_rb_2"}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 records after purge, got %d", count)
	}

	if !repo.CanRollbackDBPurge(trashDir) {
		t.Errorf("Expected CanRollback to be true after purge")
	}

	// 3. ロールバック（復元）
	restored, err := repo.RollbackLastDBPurge(trashDir)
	if err != nil || restored != 2 {
		t.Fatalf("RollbackLastDBPurge failed: count=%d, err=%v", restored, err)
	}

	// レコードが完全に復元されていることを確認
	db.Table("media").Where("media_id IN ?", []string{"media_rb_1", "media_rb_2"}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 records after rollback, got %d", count)
	}

	var fetched models.Media
	db.Table("media").Where("media_id = ?", "media_rb_1").First(&fetched)
	if fetched.StashImageID.String != "img_rb_1" {
		t.Errorf("Expected restored StashImageID 'img_rb_1', got '%s'", fetched.StashImageID.String)
	}
}
