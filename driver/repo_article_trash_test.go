// driver/repo_article_trash_test.go (Under 100 lines)
package driver_test

import (
	"path/filepath"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTrashTestDB(t *testing.T) *driver.Repository {
	tempFile := filepath.Join(t.TempDir(), "test_trash.db")
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

func TestTrashArticleAndRestore(t *testing.T) {
	repo := setupTrashTestDB(t)
	art := models.Article{
		ID: "art_trash_001", AccountID: "acc_001", ConversationID: "conv_001",
		CreatedAt: time.Now(), FullText: "Famicom retro music reproduction", Lang: "en",
		WaybackURL: "https://web.archive.org/web/123",
	}
	if err := repo.UpsertArticleTx(&art); err != nil { t.Fatalf("UpsertArticleTx failed: %v", err) }

	// 1. 通常検索でヒットすることを確認
	items, total, err := repo.SearchArticles("Famicom", "all", "all", 10, 0)
	if err != nil || total != 1 || len(items) != 1 { t.Fatalf("Expected 1 article, got total=%d, err=%v", total, err) }

	// 2. ゴミ箱へ移動 (trashedBy="admin", reason="デバッグデータの整理")
	reason := "デバッグデータの整理"
	if err := repo.TrashArticle("art_trash_001", "admin", reason); err != nil {
		t.Fatalf("TrashArticle failed: %v", err)
	}

	// 3. 通常検索から除外されていることを確認
	itemsAfter, totalAfter, err := repo.SearchArticles("Famicom", "all", "all", 10, 0)
	if err != nil || totalAfter != 0 || len(itemsAfter) != 0 {
		t.Fatalf("Expected 0 articles in normal search, got total=%d", totalAfter)
	}

	// 4. trash フィルターで取得できること、および trashed_by / trash_reason の検証
	trashedItems, trashedTotal, err := repo.SearchArticles("Famicom", "all", "trash", 10, 0)
	if err != nil || trashedTotal != 1 || len(trashedItems) != 1 {
		t.Fatalf("Expected 1 article in trash search, got total=%d, err=%v", trashedTotal, err)
	}
	if !trashedItems[0].IsTrash || trashedItems[0].TrashedBy.String != "admin" || trashedItems[0].TrashReason.String != reason {
		t.Fatalf("Trash metadata mismatch: is_trash=%v, trashed_by=%s, reason=%s",
			trashedItems[0].IsTrash, trashedItems[0].TrashedBy.String, trashedItems[0].TrashReason.String)
	}

	// 5. 復元
	if err := repo.RestoreArticle("art_trash_001"); err != nil {
		t.Fatalf("RestoreArticle failed: %v", err)
	}
	restoredItems, restoredTotal, err := repo.SearchArticles("Famicom", "all", "all", 10, 0)
	if err != nil || restoredTotal != 1 || len(restoredItems) != 1 {
		t.Fatalf("Expected 1 restored article in normal search, got total=%d", restoredTotal)
	}
	if restoredItems[0].IsTrash { t.Fatalf("Expected restored article is_trash=false") }
}
