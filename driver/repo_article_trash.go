// driver/repo_article_trash.go (Under 100 lines - SPEC-PRINCIPLE-001)
package driver

import (
	"database/sql"
	"time"

	"dozou_katanuki/models"
)

// TrashArticle は指定された記事を論理削除（ゴミ箱へ移動）し、削除実行主体および理由を記録します
func (r *Repository) TrashArticle(id, trashedBy, reason string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_trash":     true,
		"trashed_by":   sql.NullString{String: trashedBy, Valid: trashedBy != ""},
		"trash_reason": sql.NullString{String: reason, Valid: reason != ""},
		"trashed_at":   sql.NullTime{Time: now, Valid: true},
	}
	return r.db.Model(&models.Article{}).Where("id = ?", id).Updates(updates).Error
}

// RestoreArticle はゴミ箱にある指定記事を復元します
func (r *Repository) RestoreArticle(id string) error {
	updates := map[string]interface{}{
		"is_trash":     false,
		"trashed_by":   nil,
		"trash_reason": nil,
		"trashed_at":   nil,
	}
	return r.db.Model(&models.Article{}).Where("id = ?", id).Updates(updates).Error
}
