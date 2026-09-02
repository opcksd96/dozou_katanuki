// driver/repo_media_register.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// RegisterCompletedMediaFile はダウンロード完了した実体ファイルパスから該当メディアレコードを特定し COMPLETED に更新します
func (r *Repository) RegisterCompletedMediaFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("empty file path")
	}
	if err := ValidateMediaFile(filePath); err != nil {
		return err
	}
	base := filepath.Base(filePath)
	cleanID := strings.TrimSuffix(base, filepath.Ext(base))

	res := r.db.Model(&models.Media{}).
		Where("media_id = ? OR media_id = ? OR download_url LIKE ?", base, cleanID, "%/"+base).
		Updates(map[string]interface{}{
			"download_status": "COMPLETED",
			"failed_reason":   nil,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no matching media found for %s", base)
	}
	return nil
}
