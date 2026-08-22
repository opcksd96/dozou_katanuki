// driver/repository.go (100行以下)
package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// BackupDatabase は SQLite の VACUUM INTO を用いてオンラインバックアップを作成します
func (r *Repository) BackupDatabase(destDir string, maxGenerations int) (string, error) {
	if destDir == "" { destDir = filepath.Join("backups", "database") }
	if err := os.MkdirAll(destDir, 0755); err != nil { return "", fmt.Errorf("failed to create backup dir: %w", err) }

	timestamp := time.Now().Format("20060102_150405")
	destPath := filepath.Join(destDir, fmt.Sprintf("archive_%s.db", timestamp))
	normalizedPath := filepath.ToSlash(destPath)
	if err := r.db.Exec(fmt.Sprintf("VACUUM INTO '%s'", normalizedPath)).Error; err != nil {
		return "", fmt.Errorf("failed to execute VACUUM INTO: %w", err)
	}
	if maxGenerations > 0 { _ = r.purgeOldBackups(destDir, maxGenerations) }
	return destPath, nil
}

func (r *Repository) purgeOldBackups(destDir string, maxGenerations int) error {
	entries, err := os.ReadDir(destDir)
	if err != nil { return err }

	var backupFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "archive_") && strings.HasSuffix(entry.Name(), ".db") {
			backupFiles = append(backupFiles, filepath.Join(destDir, entry.Name()))
		}
	}
	if len(backupFiles) <= maxGenerations { return nil }

	sort.Strings(backupFiles)
	toDelete := len(backupFiles) - maxGenerations
	for i := 0; i < toDelete; i++ { _ = os.Remove(backupFiles[i]) }
	return nil
}

// ResetDatabase は全テーブルのデータを安全に初期化（物理削除）します
func (r *Repository) ResetDatabase() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, table := range []string{"media", "url_redirects", "articles", "account_profile_histories", "accounts"} {
			if err := tx.Exec("DELETE FROM " + table).Error; err != nil { return err }
		}
		return nil
	})
}
