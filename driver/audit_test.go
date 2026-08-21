// driver/audit_test.go (100行以下)
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
	testDBPath := filepath.Join(os.TempDir(), "test_audit_"+time.Now().Format("20060102150405.000")+".db")
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil { t.Fatalf("InitDB failed: %v", err) }
	repo := NewRepository(db)

	tempDir := t.TempDir()
	stashDir, blobsDir := filepath.Join(tempDir, "stash"), filepath.Join(tempDir, "blobs")
	_ = os.MkdirAll(filepath.Join(stashDir, "scenes"), 0755)
	_ = os.MkdirAll(filepath.Join(stashDir, "images"), 0755)
	_ = os.MkdirAll(blobsDir, 0755)

	acc := models.Account{NumericID: "1001", Username: "mashu", DisplayName: "Mash Kyrielight"}
	art := models.Article{
		ID: "art_100", AccountID: "1001", FullText: "Test post for audit",
		Media: []models.Media{{
			MediaID: "media_valid_1", ArticleID: "art_100", Type: "image",
			DownloadURL: "https://example.com/img1.jpg", DownloadStatus: "COMPLETED",
			StashImageID: sql.NullString{String: "img_valid_1", Valid: true},
		}},
	}
	_ = db.Save(&acc).Error
	_ = repo.UpsertArticleTx(&art)

	_ = os.WriteFile(filepath.Join(stashDir, "images", "img_valid_1.jpg"), []byte("valid"), 0644)
	_ = os.WriteFile(filepath.Join(blobsDir, "media_valid_1.jpg"), []byte("valid"), 0644)

	_ = db.Exec("PRAGMA foreign_keys = OFF;").Error
	orphanMedia := models.Media{MediaID: "media_orphan_db", ArticleID: "art_non_existent", Type: "video", DownloadStatus: "QUEUED"}
	_ = db.Save(&orphanMedia).Error
	_ = db.Exec("PRAGMA foreign_keys = ON;").Error

	orphanFile := filepath.Join(stashDir, "scenes", "scene_orphan_999.mp4")
	_ = os.WriteFile(orphanFile, []byte("dummy video"), 0644)

	report, err := repo.AuditDatabase(stashDir, blobsDir)
	if err != nil || !report.IntegrityOK { t.Fatalf("AuditDatabase failed: %v", err) }

	foundDBOrphan, foundFileOrphan := false, false
	for _, m := range report.OrphanDBMedia { if m.MediaID == "media_orphan_db" { foundDBOrphan = true } }
	for _, f := range report.OrphanFiles { if f.FileName == "scene_orphan_999.mp4" { foundFileOrphan = true } }

	if !foundDBOrphan || !foundFileOrphan { t.Errorf("Orphan detection failed: DB=%v, File=%v", foundDBOrphan, foundFileOrphan) }

	purgedFiles, err := repo.PurgeOrphanFiles([]string{orphanFile})
	if err != nil || purgedFiles != 1 { t.Errorf("PurgeOrphanFiles failed: %v, count=%d", err, purgedFiles) }

	purgedDB, err := repo.PurgeOrphanDBMedia([]string{"media_orphan_db"})
	if err != nil || purgedDB != 1 { t.Errorf("PurgeOrphanDBMedia failed: %v, count=%d", err, purgedDB) }
}
