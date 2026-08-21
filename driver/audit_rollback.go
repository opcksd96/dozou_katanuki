// driver/audit_rollback.go (100行以下 - SPEC-AUDIT-ROLLBACK)
package driver

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurgeBackupPackage struct {
	ExecutedAt time.Time      `json:"executed_at"`
	Count      int            `json:"count"`
	Records    []models.Media `json:"records"`
}

func BackupAndPurgeDBMedia(db *gorm.DB, trashDir string, mediaIDs []string) (int, error) {
	if len(mediaIDs) == 0 { return 0, nil }
	if trashDir == "" { trashDir = "./backups/dumps/_trash" }
	_ = os.MkdirAll(trashDir, 0755)

	var targets []models.Media
	if err := db.Table("media").Where("media_id IN ?", mediaIDs).Find(&targets).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch records for backup: %w", err)
	}

	if len(targets) > 0 {
		pkg := PurgeBackupPackage{ExecutedAt: time.Now(), Count: len(targets), Records: targets}
		if data, err := json.MarshalIndent(pkg, "", "  "); err == nil {
			fn := fmt.Sprintf("purge_media_%s.json", time.Now().Format("20060102_150405"))
			_ = os.WriteFile(filepath.Join(trashDir, fn), data, 0644)
		}
	}

	res := db.Table("media").Where("media_id IN ?", mediaIDs).Delete(&models.Media{})
	if res.Error != nil { return 0, fmt.Errorf("failed to delete media records: %w", res.Error) }
	return int(res.RowsAffected), nil
}

func RollbackLastDBPurge(db *gorm.DB, trashDir string) (int, error) {
	if trashDir == "" { trashDir = "./backups/dumps/_trash" }
	entries, err := os.ReadDir(trashDir)
	if err != nil { return 0, fmt.Errorf("trash directory not found: %w", err) }

	var jsonFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "purge_media_") && strings.HasSuffix(e.Name(), ".json") {
			jsonFiles = append(jsonFiles, filepath.Join(trashDir, e.Name()))
		}
	}
	if len(jsonFiles) == 0 { return 0, fmt.Errorf("no rollback snapshots available") }

	sort.Slice(jsonFiles, func(i, j int) bool { return jsonFiles[i] > jsonFiles[j] })
	latestFile := jsonFiles[0]

	data, err := os.ReadFile(latestFile)
	if err != nil { return 0, fmt.Errorf("failed to read rollback snapshot: %w", err) }

	var pkg PurgeBackupPackage
	if err := json.Unmarshal(data, &pkg); err != nil { return 0, fmt.Errorf("failed to parse rollback snapshot: %w", err) }

	restoredCount := 0
	for _, m := range pkg.Records {
		if err := db.Table("media").Clauses(clause.OnConflict{UpdateAll: true}).Create(&m).Error; err == nil {
			restoredCount++
		} else {
			log.Printf("[Rollback] Warning failed to restore media %s: %v", m.MediaID, err)
		}
	}

	_ = os.Rename(latestFile, latestFile+".restored")
	return restoredCount, nil
}

func CanRollbackDBPurge(trashDir string) bool {
	if trashDir == "" { trashDir = "./backups/dumps/_trash" }
	entries, err := os.ReadDir(trashDir)
	if err != nil { return false }
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "purge_media_") && strings.HasSuffix(e.Name(), ".json") { return true }
	}
	return false
}
