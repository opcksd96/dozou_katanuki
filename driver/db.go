// driver/db.go (100行以下)
package driver

import (
	"log"

	"dozou_katanuki/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, pragma := range pragmas {
		if _, err := sqlDB.Exec(pragma); err != nil {
			log.Fatalf("[FATAL] Failed to apply pragma '%s': %v", pragma, err)
		}
	}

	err = db.AutoMigrate(
		&models.Account{},
		&models.AccountProfileHistory{},
		&models.Article{},
		&models.Media{},
		&models.UrlRedirect{},
		&models.Whitelist{},
	)
	if err != nil {
		return nil, err
	}

	// 既存 assets/ 内のアバター画像を DB の avatar_base64 に移行
	_ = MigrateAvatarsToBase64(db)

	log.Printf("[Driver] Database initialized successfully (WAL mode enabled): %s", dbPath)
	return db, nil
}
