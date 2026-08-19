package driver

import (
	"fmt"
	"log"
	"time"

	"dozou_katanuki/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB はSQLite3データベースへの接続を確立し、WALモードと自動マイグレーションを適用します
func InitDB(dbPath string) (*gorm.DB, error) {
	// WALモードおよびビジータイムアウトを有効化するDSNクエリ
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL", dbPath)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// コネクションプールの設定（SQLiteのロック競合を防止）
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自動マイグレーション（テーブル作成・更新）
	err = db.AutoMigrate(
		&models.Account{},
		&models.AccountProfileHistory{},
		&models.Article{},
		&models.Media{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Printf("[Driver] Database initialized successfully (WAL mode enabled): %s", dbPath)
	return db, nil
}
