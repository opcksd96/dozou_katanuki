// app/app_test.go (100行以下)
package app

import (
	"path/filepath"
	"testing"

	"dozou_katanuki/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, string) {
	tempFile := filepath.Join(t.TempDir(), "test_archive.db")
	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("Failed to open test db: %v", err) }
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil && sqlDB != nil { _ = sqlDB.Close() }
	})
	_ = db.AutoMigrate(
		&models.Account{}, &models.AccountProfileHistory{},
		&models.Article{}, &models.Media{},
		&models.UrlRedirect{}, &models.Whitelist{},
	)
	return db, tempFile
}
