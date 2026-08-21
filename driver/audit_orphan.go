// driver/audit_orphan.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
	"gorm.io/gorm"
)

func FindOrphanDBMedia(db *gorm.DB, stashDir, blobsDir string) ([]models.OrphanDBMedia, error) {
	var orphans []models.OrphanDBMedia
	var missingParent []models.Media
	if err := db.Raw("SELECT m.* FROM media m LEFT JOIN articles a ON m.article_id = a.id WHERE a.id IS NULL").Scan(&missingParent).Error; err != nil {
		return nil, err
	}
	for _, m := range missingParent {
		orphans = append(orphans, models.OrphanDBMedia{
			MediaID: m.MediaID, ArticleID: m.ArticleID, Type: m.Type, DownloadURL: m.DownloadURL,
			Status: m.DownloadStatus, StashSceneID: m.StashSceneID.String, StashImageID: m.StashImageID.String,
			Reason: fmt.Sprintf("親記事 (%s) が存在しません（参照切れ孤立レコード）", m.ArticleID),
		})
	}
	return orphans, nil
}

func GetKnownMediaIdentifiers(db *gorm.DB) (map[string]bool, error) {
	known := make(map[string]bool)
	type idRecord struct {
		MediaID string `gorm:"column:media_id"`; StashSceneID sql.NullString `gorm:"column:stash_scene_id"`; StashImageID sql.NullString `gorm:"column:stash_image_id"`
	}
	var recs []idRecord
	if err := db.Model(&models.Media{}).Select("media_id, stash_scene_id, stash_image_id").Scan(&recs).Error; err != nil { return nil, err }
	for _, r := range recs {
		if r.MediaID != "" { known[strings.ToLower(r.MediaID)] = true }
		if r.StashSceneID.Valid && r.StashSceneID.String != "" { known[strings.ToLower(r.StashSceneID.String)] = true }
		if r.StashImageID.Valid && r.StashImageID.String != "" { known[strings.ToLower(r.StashImageID.String)] = true }
	}
	return known, nil
}

func ScanOrphanFiles(stashDir, blobsDir string, knownKeys map[string]bool) ([]models.OrphanFile, error) {
	var list []models.OrphanFile
	scanDir := func(dirPath, category string) {
		if dirPath == "" { return }
		_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() { return nil }
			base := info.Name()
			if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".json") || strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".yaml") { return nil }
			name := strings.ToLower(strings.TrimSuffix(base, filepath.Ext(base)))
			isKnown := knownKeys[name]
			if !isKnown {
				for k := range knownKeys { if len(k) >= 4 && strings.Contains(name, k) { isKnown = true; break } }
			}
			if !isKnown { list = append(list, models.OrphanFile{Path: path, FileName: base, FileSize: info.Size(), Category: category}) }
			return nil
		})
	}
	if stashDir != "" {
		scanDir(filepath.Join(stashDir, "scenes"), "stash_scene"); scanDir(filepath.Join(stashDir, "images"), "stash_image"); scanDir(filepath.Join(stashDir, "generated"), "stash_generated")
	}
	if blobsDir != "" { scanDir(blobsDir, "blob") }
	return list, nil
}

func DeleteDBMediaByIDs(db *gorm.DB, mediaIDs []string) (int, error) {
	if len(mediaIDs) == 0 { return 0, nil }
	res := db.Table("media").Where("media_id IN ?", mediaIDs).Delete(&models.Media{})
	if res.Error != nil {
		return 0, fmt.Errorf("failed to delete media records: %w", res.Error)
	}
	return int(res.RowsAffected), nil
}
