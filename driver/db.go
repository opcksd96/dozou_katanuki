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

	// GAP-1: Wiki仕様 10大インデックス完全適用 (SPEC-DATABASE-001 §1.3)
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_articles_is_liked_created ON articles(is_liked, created_at DESC) WHERE is_liked = 1;",
		"CREATE INDEX IF NOT EXISTS idx_articles_account_created ON articles(account_id, created_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_articles_conversation ON articles(conversation_id, created_at ASC);",
		"CREATE INDEX IF NOT EXISTS idx_articles_reply_to ON articles(reply_to_id);",
		"CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_history_lookup ON account_profile_histories(account_id, avatar_seq DESC);",
		"CREATE INDEX IF NOT EXISTS idx_accounts_username ON accounts(username);",
		"CREATE INDEX IF NOT EXISTS idx_media_article ON media(article_id);",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_media_stash_scene ON media(stash_scene_id) WHERE stash_scene_id IS NOT NULL;",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_media_stash_image ON media(stash_image_id) WHERE stash_image_id IS NOT NULL;",
		"CREATE INDEX IF NOT EXISTS idx_media_status_type ON media(download_status, type);",
		"CREATE INDEX IF NOT EXISTS idx_whitelist_type_active ON whitelists(type, is_active);",
		"CREATE INDEX IF NOT EXISTS idx_articles_is_trash ON articles(is_trash, created_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_articles_trashed_by ON articles(trashed_by) WHERE is_trash = 1;",
	}
	for _, idxSql := range indexes { _ = db.Exec(idxSql).Error }

	_ = MigrateAvatarsToBase64(db)
	log.Printf("[Driver] Database initialized (WAL mode & 10 Major Indexes applied): %s", dbPath)
	return db, nil
}
