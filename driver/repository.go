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

// GetWhitelists は登録されているすべてのホワイトリストを取得します
func (r *Repository) GetWhitelists() ([]models.Whitelist, error) {
	var list []models.Whitelist
	err := r.db.Order("id ASC").Find(&list).Error
	return list, err
}

// AddWhitelist はホワイトリスト項目を追加します
func (r *Repository) AddWhitelist(itemType, value string) (*models.Whitelist, error) {
	item := models.Whitelist{
		Type:     itemType,
		Value:    value,
		IsActive: true,
	}
	if err := r.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateWhitelist はホワイトリスト項目の内容を更新します
func (r *Repository) UpdateWhitelist(id uint, itemType, value string, isActive bool) error {
	return r.db.Model(&models.Whitelist{}).Where("id = ?", id).Updates(map[string]interface{}{
		"type":      itemType,
		"value":     value,
		"is_active": isActive,
	}).Error
}

// DeleteWhitelist はホワイトリスト項目を削除します
func (r *Repository) DeleteWhitelist(id uint) error {
	return r.db.Delete(&models.Whitelist{}, id).Error
}

// ToggleWhitelist はホワイトリスト項目の有効/無効を切り替えます
func (r *Repository) ToggleWhitelist(id uint) error {
	var item models.Whitelist
	if err := r.db.First(&item, id).Error; err != nil {
		return err
	}
	item.IsActive = !item.IsActive
	return r.db.Save(&item).Error
}

// SearchArticles はキーワード・アカウント・各種フィルターによる柔軟な記事検索と総件数取得を行います
func (r *Repository) SearchArticles(searchQuery, accountID, filter string, limit, offset int) ([]models.Article, int64, error) {
	query := r.db.Model(&models.Article{}).
		Preload("Account").
		Preload("Account.ProfileHistory").
		Preload("Media")

	if searchQuery != "" {
		likePattern := "%" + searchQuery + "%"
		query = query.Where("full_text LIKE ? OR full_text_ja LIKE ? OR full_text_en LIKE ? OR full_text_zh LIKE ? OR id LIKE ?",
			likePattern, likePattern, likePattern, likePattern, likePattern)
	}

	if accountID != "" && accountID != "all" {
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

	var total int64
	// カウント計算（Group によるカウント不整合を避けるため、Articles テーブルベースでカウント）
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []models.Article
	if limit <= 0 {
		limit = 20
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, total, err
}

// GetArticleByID は指定されたIDの単一記事を取得します
func (r *Repository) GetArticleByID(id string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Account").
		Preload("Account.ProfileHistory").
		Preload("Media").
		Where("id = ?", id).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// UpdateArticleTranslations は指定された記事の日本語・英語・中国語の翻訳テキストを更新します
func (r *Repository) UpdateArticleTranslations(id string, ja, en, zh string) error {
	updates := map[string]interface{}{}
	if ja != "" {
		updates["full_text_ja"] = ja
	} else {
		updates["full_text_ja"] = nil
	}
	if en != "" {
		updates["full_text_en"] = en
	} else {
		updates["full_text_en"] = nil
	}
	if zh != "" {
		updates["full_text_zh"] = zh
	} else {
		updates["full_text_zh"] = nil
	}

	return r.db.Model(&models.Article{}).Where("id = ?", id).Updates(updates).Error
}

