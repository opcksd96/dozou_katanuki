// driver/repo_media_destination.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GetMediaOwnerUsername はメディアIDから所属するアカウントのusernameを取得します
func (r *Repository) GetMediaOwnerUsername(mediaID string) (string, error) {
	if r.db == nil || mediaID == "" {
		return "", fmt.Errorf("invalid args")
	}
	cleanID := strings.TrimSuffix(mediaID, filepath.Ext(mediaID))
	for _, sfx := range []string{"_orig", "_large", "_wayback_orig", "_wayback"} {
		cleanID = strings.TrimSuffix(cleanID, sfx)
	}

	var result struct {
		Username string
	}
	err := r.db.Table("media").
		Select("accounts.username").
		Joins("JOIN articles ON articles.id = media.article_id").
		Joins("JOIN accounts ON (accounts.numeric_id = articles.account_id OR accounts.username = articles.account_id)").
		Where("media.media_id = ? OR media.media_id = ? OR media.media_id LIKE ?", mediaID, cleanID, cleanID+"%").
		First(&result).Error
	if err != nil {
		return "", err
	}
	return result.Username, nil
}
