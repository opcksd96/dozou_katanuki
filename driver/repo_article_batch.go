// driver/repo_article_batch.go (Under 100 lines - SPEC-PRINCIPLE-001)
package driver

import (
	"database/sql"
	"time"

	"dozou_katanuki/models"
)

// BatchTrashArticles は指定された複数記事を一括で論理削除（ゴミ箱移動）します
func (r *Repository) BatchTrashArticles(ids []string, trashedBy, reason string) error {
	if len(ids) == 0 { return nil }
	now := time.Now()
	updates := map[string]interface{}{
		"is_trash":     true,
		"trashed_by":   sql.NullString{String: trashedBy, Valid: trashedBy != ""},
		"trash_reason": sql.NullString{String: reason, Valid: reason != ""},
		"trashed_at":   sql.NullTime{Time: now, Valid: true},
	}
	return r.db.Model(&models.Article{}).Where("id IN ?", ids).Updates(updates).Error
}

// BatchRestoreArticles はゴミ箱にある指定複数記事を一括復元します
func (r *Repository) BatchRestoreArticles(ids []string) error {
	if len(ids) == 0 { return nil }
	updates := map[string]interface{}{
		"is_trash":     false,
		"trashed_by":   nil,
		"trash_reason": nil,
		"trashed_at":   nil,
	}
	return r.db.Model(&models.Article{}).Where("id IN ?", ids).Updates(updates).Error
}

// BatchResetTranslations は指定された複数記事の翻訳データ（JA/EN/ZH）を一括初期化します
func (r *Repository) BatchResetTranslations(ids []string) error {
	if len(ids) == 0 { return nil }
	updates := map[string]interface{}{
		"full_text_ja": nil,
		"full_text_en": nil,
		"full_text_zh": nil,
	}
	return r.db.Model(&models.Article{}).Where("id IN ?", ids).Updates(updates).Error
}

// GetArticlesByIDs は指定された複数IDの記事一覧を取得します（Undo/Redoスナップショット用）
func (r *Repository) GetArticlesByIDs(ids []string) ([]models.Article, error) {
	if len(ids) == 0 { return []models.Article{}, nil }
	var articles []models.Article
	err := r.db.Where("id IN ?", ids).Find(&articles).Error
	return articles, err
}
