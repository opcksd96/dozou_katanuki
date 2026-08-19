// driver/repository.go (100行以下)
package driver

import (
	"dozou_katanuki/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Preload("ProfileHistory").Order("username ASC").Find(&accounts).Error
	return accounts, err
}

func (r *Repository) GetAccountHistories(accountID string) ([]models.AccountProfileHistory, error) {
	var histories []models.AccountProfileHistory
	err := r.db.Where("account_id = ?", accountID).Order("avatar_seq desc").Find(&histories).Error
	return histories, err
}

// FetchArticles はインデックスを活用して最大50件の生記事データを取得します
func (r *Repository) FetchArticles(accountID, filter string, limit, offset int) ([]models.Article, error) {
	query := r.db.Model(&models.Article{}).
		Preload("Account").
		Preload("Account.ProfileHistory").
		Preload("Media").
		Order("created_at DESC")

	// 外向きリツイートの排除（リツイート対象がアクティブなホワイトリストに存在しない投稿を除外）
	query = query.Where("is_repost = ? OR reply_to_handle IN (SELECT value FROM whitelists WHERE is_active = ?)", false, true)

	if accountID != "all" {
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
	err := query.Limit(limit).Offset(offset).Find(&articles).Error
	return articles, err
}

// UpsertArticleTx は記事・アカウント・メディアをトランザクション内で保存します
func (r *Repository) UpsertArticleTx(art *models.Article) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if art.Account.NumericID != "" {
			if err := tx.Save(&art.Account).Error; err != nil {
				return err
			}
		}
		if err := tx.Save(art).Error; err != nil {
			return err
		}
		for i := range art.Media {
			art.Media[i].ArticleID = art.ID
			if err := tx.Save(&art.Media[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
