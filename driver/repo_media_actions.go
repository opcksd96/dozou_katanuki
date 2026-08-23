// driver/repo_media_actions.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// ResolveMediaFilePath はメディアIDまたはStash IDからローカル上の実体ファイルパスを探索・解決します
func (r *Repository) ResolveMediaFilePath(mediaID string) (string, error) {
	var m models.Media
	if err := r.db.Where("media_id = ?", mediaID).First(&m).Error; err != nil {
		return "", err
	}

	searchDirs := []string{"blobs", "assets", "G:/Media_Storage/Influencers", "backups"}
	cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))

	for _, baseDir := range searchDirs {
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			continue
		}
		var foundPath string
		_ = filepath.Walk(baseDir, func(p string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}
			name := info.Name()
			if strings.EqualFold(name, m.MediaID) || strings.EqualFold(strings.TrimSuffix(name, filepath.Ext(name)), cleanID) {
				foundPath = p
				return filepath.SkipAll
			}
			return nil
		})
		if foundPath != "" {
			return filepath.Abs(foundPath)
		}
	}
	return "", fmt.Errorf("media file not found on disk for ID: %s", mediaID)
}

// ToggleMediaBookmark はメディアのブックマーク（お気に入り）状態を反転させます
func (r *Repository) ToggleMediaBookmark(mediaID string) (bool, error) {
	var m models.Media
	if err := r.db.Where("media_id = ?", mediaID).First(&m).Error; err != nil {
		return false, err
	}
	newVal := !m.IsBookmarked
	err := r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Update("is_bookmarked", newVal).Error
	return newVal, err
}

// RenameMedia はメディアIDおよび紐づく実体ファイル名をリネームします
func (r *Repository) RenameMedia(mediaID, newMediaID string) error {
	if newMediaID == "" || mediaID == newMediaID {
		return fmt.Errorf("invalid new media ID")
	}
	filePath, err := r.ResolveMediaFilePath(mediaID)
	if err == nil && filePath != "" {
		newPath := filepath.Join(filepath.Dir(filePath), newMediaID+filepath.Ext(filePath))
		_ = os.Rename(filePath, newPath)
	}
	return r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Update("media_id", newMediaID).Error
}

// RequeueMediaByStatus は指定ステータスのメディアを一括で QUEUED 状態に戻して再ダウンロード対象にします
func (r *Repository) RequeueMediaByStatus(status, accountID string) (int64, error) {
	q := r.db.Table("media")
	if accountID != "" && accountID != "all" {
		q = q.Where("article_id IN (SELECT id FROM articles WHERE account_id = ? OR account_id IN (SELECT numeric_id FROM accounts WHERE username = ?))", accountID, accountID)
	}
	if status == "DEAD_404" {
		q = q.Where("download_status = 'DEAD_404'")
	} else if status == "EXCLUDED" {
		q = q.Where("download_status = 'EXCLUDED' OR failed_reason LIKE '%Whitelist外%'")
	} else if status != "" && status != "all" {
		q = q.Where("download_status = ?", status)
	}
	res := q.Updates(map[string]interface{}{"download_status": "QUEUED", "failed_reason": nil})
	return res.RowsAffected, res.Error
}
