// driver/repo_article_batch_test.go (Under 100 lines)
package driver_test

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupBatchTestDB(t *testing.T) *driver.Repository {
	tempFile := filepath.Join(t.TempDir(), "test_batch.db")
	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil { t.Fatalf("Failed to open test db: %v", err) }
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil && sqlDB != nil { _ = sqlDB.Close() }
	})
	if err := db.AutoMigrate(&models.Account{}, &models.AccountProfileHistory{}, &models.Article{}, &models.Media{}, &models.UrlRedirect{}, &models.Whitelist{}); err != nil {
		t.Fatalf("Failed to automigrate: %v", err)
	}
	return driver.NewRepository(db)
}

func TestBatchArticleOps(t *testing.T) {
	repo := setupBatchTestDB(t)
	art1 := models.Article{
		ID: "batch_001", AccountID: "acc_001", ConversationID: "conv_001", CreatedAt: time.Now(),
		FullText: "Song 1", FullTextJA: sql.NullString{String: "曲 1", Valid: true}, WaybackURL: "https://web.archive.org/1",
	}
	art2 := models.Article{
		ID: "batch_002", AccountID: "acc_001", ConversationID: "conv_001", CreatedAt: time.Now(),
		FullText: "Song 2", FullTextJA: sql.NullString{String: "曲 2", Valid: true}, WaybackURL: "https://web.archive.org/2",
	}
	_ = repo.UpsertArticleTx(&art1)
	_ = repo.UpsertArticleTx(&art2)

	// 1. 一括翻訳リセットの検証
	if err := repo.BatchResetTranslations([]string{"batch_001", "batch_002"}); err != nil {
		t.Fatalf("BatchResetTranslations failed: %v", err)
	}
	fetched, _ := repo.GetArticlesByIDs([]string{"batch_001", "batch_002"})
	for _, a := range fetched {
		if a.FullTextJA.Valid && a.FullTextJA.String != "" {
			t.Fatalf("Expected translation to be reset, got: %s", a.FullTextJA.String)
		}
	}

	// 2. 一括ゴミ箱移動の検証
	if err := repo.BatchTrashArticles([]string{"batch_001", "batch_002"}, "admin", "一括整理"); err != nil {
		t.Fatalf("BatchTrashArticles failed: %v", err)
	}
	normalItems, normalTotal, _ := repo.SearchArticles("", "all", "all", 10, 0)
	if normalTotal != 0 || len(normalItems) != 0 {
		t.Fatalf("Expected 0 articles in normal search, got %d", normalTotal)
	}
	allWithTrash, allTotal, _ := repo.SearchArticles("", "all", "all_with_trash", 10, 0)
	if allTotal != 2 || len(allWithTrash) != 2 {
		t.Fatalf("Expected 2 articles in all_with_trash, got %d", allTotal)
	}

	// 3. 一括復元の検証
	if err := repo.BatchRestoreArticles([]string{"batch_001", "batch_002"}); err != nil {
		t.Fatalf("BatchRestoreArticles failed: %v", err)
	}
	restored, restoredTotal, _ := repo.SearchArticles("", "all", "all", 10, 0)
	if restoredTotal != 2 || len(restored) != 2 {
		t.Fatalf("Expected 2 restored articles, got %d", restoredTotal)
	}
}
