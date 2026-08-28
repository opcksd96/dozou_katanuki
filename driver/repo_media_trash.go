// driver/repo_media_trash.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"fmt"
	"strings"
	"time"

	"dozou_katanuki/models"
)

// TrashMedia は単一メディアを論理削除（ゴミ箱へ移動）します
func (r *Repository) TrashMedia(mediaID, reason, trashedBy string) error {
	id := strings.TrimSpace(mediaID)
	if id == "" { return fmt.Errorf("mediaID must not be empty") }
	now := time.Now()
	res := r.db.Model(&models.Media{}).Where("media_id = ?", id).Updates(map[string]interface{}{
		"is_trash": true, "trash_reason": reason, "trashed_by": trashedBy, "trashed_at": &now,
	})
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return fmt.Errorf("media not found: %s", id) }
	return nil
}

// RestoreMedia はゴミ箱に入っている単一メディアを通常状態へ復元します
func (r *Repository) RestoreMedia(mediaID string) error {
	id := strings.TrimSpace(mediaID)
	if id == "" { return fmt.Errorf("mediaID must not be empty") }
	res := r.db.Model(&models.Media{}).Where("media_id = ?", id).Updates(map[string]interface{}{
		"is_trash": false, "trash_reason": "", "trashed_by": "", "trashed_at": nil,
	})
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return fmt.Errorf("media not found: %s", id) }
	return nil
}

// BatchTrashMedia は複数メディアを一括で論理削除します
func (r *Repository) BatchTrashMedia(mediaIDs []string, reason, trashedBy string) error {
	if len(mediaIDs) == 0 { return nil }
	now := time.Now()
	return r.db.Model(&models.Media{}).Where("media_id IN ?", mediaIDs).Updates(map[string]interface{}{
		"is_trash": true, "trash_reason": reason, "trashed_by": trashedBy, "trashed_at": &now,
	}).Error
}

// BatchRestoreMedia は複数メディアを一括復元します
func (r *Repository) BatchRestoreMedia(mediaIDs []string) error {
	if len(mediaIDs) == 0 { return nil }
	return r.db.Model(&models.Media{}).Where("media_id IN ?", mediaIDs).Updates(map[string]interface{}{
		"is_trash": false, "trash_reason": "", "trashed_by": "", "trashed_at": nil,
	}).Error
}
