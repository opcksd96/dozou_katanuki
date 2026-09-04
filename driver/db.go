// driver/db.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return nil, err
	}

	_ = MigrateAvatarsToBase64(db)
	log.Printf("[Driver] Database initialized (WAL mode & Migrations applied): %s", dbPath)
	return db, nil
}
