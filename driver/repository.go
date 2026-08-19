package driver

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetArticles は条件に応じたタイムライン投稿を時系列降順で取得します
func (r *Repository) GetArticles(accountID, filter string, limit, offset int) ([]models.Article, error) {
	query := r.db.Model(&models.Article{}).
		Preload("Account").
		Preload("Media").
		Order("created_at DESC")

	if accountID != "all" && accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}

	switch filter {
	case "reposts":
		query = query.Where("is_repost = ?", true)
	case "media":
		query = query.Joins("JOIN media ON media.article_id = articles.id").Group("articles.id")
	case "bookmarks":
		query = query.Where("is_liked = ?", true)
	}

	var articles []models.Article
	if err := query.Limit(limit).Offset(offset).Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}
	return articles, nil
}

// UpsertArticleTx は共通中間JSONデータ（記事・アカウント・メディア）を一括保存します
func (r *Repository) UpsertArticleTx(article *models.Article) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. アカウント & アバター世代監査
		if _, err := AuditAndResolveAvatar(tx, &article.Account); err != nil {
			return err
		}
		article.Account.UpdatedAt = time.Now()
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&article.Account).Error; err != nil {
			return err
		}

		// 2. 記事本体のUpsert
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(article).Error; err != nil {
			return err
		}

		// 3. 添付メディアのUpsert
		for i := range article.Media {
			article.Media[i].ArticleID = article.ID
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&article.Media[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
