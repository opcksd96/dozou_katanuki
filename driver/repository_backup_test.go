package driver_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, string) {
	tempFile := filepath.Join(t.TempDir(), "test_repo_backup.db")
	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil && sqlDB != nil {
			_ = sqlDB.Close()
		}
	})

	err = db.AutoMigrate(
		&models.Account{},
		&models.Article{},
		&models.Media{},
	)
	if err != nil {
		t.Fatalf("Failed to automigrate: %v", err)
	}

	return db, tempFile
}

func TestBackupDatabaseAndPurge(t *testing.T) {
	db, _ := setupTestDB(t)

	backupDir := filepath.Join(t.TempDir(), "test_backups")

	repo := driver.NewRepository(db)

	// 1回目バックアップ
	path1, err := repo.BackupDatabase(backupDir, 2)
	if err != nil {
		t.Fatalf("BackupDatabase #1 failed: %v", err)
	}
	if _, err := os.Stat(path1); err != nil {
		t.Fatalf("Backup file 1 not found: %v", err)
	}

	time.Sleep(1100 * time.Millisecond)

	// 2回目バックアップ
	path2, err := repo.BackupDatabase(backupDir, 2)
	if err != nil {
		t.Fatalf("BackupDatabase #2 failed: %v", err)
	}
	if _, err := os.Stat(path2); err != nil {
		t.Fatalf("Backup file 2 not found: %v", err)
	}

	time.Sleep(1100 * time.Millisecond)

	// 3回目バックアップ（世代数2を指定しているため、1回目のバックアップがパージされるはず）
	path3, err := repo.BackupDatabase(backupDir, 2)
	if err != nil {
		t.Fatalf("BackupDatabase #3 failed: %v", err)
	}
	if _, err := os.Stat(path3); err != nil {
		t.Fatalf("Backup file 3 not found: %v", err)
	}

	// path1 は削除され、path2 と path3 が残っているか検証
	if _, err := os.Stat(path1); !os.IsNotExist(err) {
		t.Fatalf("Expected path1 (%s) to be purged, but it still exists", path1)
	}
	if _, err := os.Stat(path2); err != nil {
		t.Fatalf("Expected path2 (%s) to exist: %v", path2, err)
	}
	if _, err := os.Stat(path3); err != nil {
		t.Fatalf("Expected path3 (%s) to exist: %v", path3, err)
	}
}
