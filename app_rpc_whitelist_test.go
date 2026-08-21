// app_rpc_whitelist_test.go (100行以下)
package main

import (
	"os"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, string) {
	tempFile := "test_archive_" + time.Now().Format("20060102150405.000") + ".db"
	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("Failed to open test db: %v", err) }
	_ = db.AutoMigrate(&models.Account{}, &models.AccountProfileHistory{}, &models.Article{}, &models.Media{}, &models.UrlRedirect{}, &models.Whitelist{})
	return db, tempFile
}

func TestWhitelistCRUD(t *testing.T) {
	db, tempFile := setupTestDB(t)
	defer os.Remove(tempFile)
	repo := driver.NewRepository(db)

	item, err := repo.AddWhitelist("account", "mashu_dev")
	if err != nil || item.ID == 0 || item.Value != "mashu_dev" || !item.IsActive {
		t.Fatalf("AddWhitelist failed: %v", err)
	}

	list, err := repo.GetWhitelists()
	if err != nil || len(list) != 1 { t.Fatalf("GetWhitelists failed: %v", err) }

	if err := repo.ToggleWhitelist(item.ID); err != nil { t.Fatalf("ToggleWhitelist failed: %v", err) }
	list2, _ := repo.GetWhitelists()
	if list2[0].IsActive != false { t.Fatalf("Expected IsActive to be false") }

	if err := repo.UpdateWhitelist(item.ID, "keyword", "retro_famicom", true); err != nil { t.Fatalf("UpdateWhitelist failed: %v", err) }
	list3, _ := repo.GetWhitelists()
	if list3[0].Type != "keyword" || list3[0].Value != "retro_famicom" { t.Fatalf("UpdateWhitelist failed: %+v", list3[0]) }

	if err := repo.DeleteWhitelist(item.ID); err != nil { t.Fatalf("DeleteWhitelist failed: %v", err) }
	list4, _ := repo.GetWhitelists()
	if len(list4) != 0 { t.Fatalf("Expected len 0, got %d", len(list4)) }
}
