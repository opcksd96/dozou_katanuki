// driver/db.go (100行以下)
package driver

import (
	"fmt"
	"log"
	"strings"

	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	dsn := dbPath
	sep := "?"
	if strings.Contains(dsn, "?") { sep = "&" }
	dsn += fmt.Sprintf("%s_pragma=busy_timeout(30000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", sep)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { return nil, err }

	sqlDB, err := db.DB()
	if err != nil { return nil, err }
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	err = db.AutoMigrate(
		&models.Account{}, &models.AccountProfileHistory{},
		&models.Article{}, &models.Media{},
		&models.UrlRedirect{}, &models.Whitelist{},
	)
	if err != nil { return nil, err }

	// GAP-1: 5つの複合インデックスを適用 (SPEC-STORAGE-001)
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_articles_account_created ON articles (account_id, created_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_media_article_type ON media (article_id, type);",
		"CREATE INDEX IF NOT EXISTS idx_media_status_type ON media (download_status, type);",
		"CREATE INDEX IF NOT EXISTS idx_history_account_recorded ON account_profile_history (account_id, recorded_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_whitelist_type_active ON whitelists (type, is_active);",
	}
	for _, idxSql := range indexes { _ = db.Exec(idxSql).Error }

	_ = MigrateAvatarsToBase64(db)
	log.Printf("[Driver] Database initialized successfully (WAL mode & Composite Indexes applied): %s", dbPath)
	return db, nil
}
