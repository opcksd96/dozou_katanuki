// driver/repo_articles.go (100行以下)
package driver

import (
	"dozou_katanuki/models"
	"gorm.io/gorm"
)

func (r *Repository) applyAccountFilter(query *gorm.DB, accountID string) *gorm.DB {
	if accountID == "" || accountID == "all" { return query }
	if len(accountID) > 6 && accountID[:6] == "group:" {
		grp := accountID[6:]
		return query.Where("account_id IN (SELECT numeric_id FROM accounts WHERE group_name = ?)", grp)
	}
	return query.Where("account_id = ? OR account_id IN (SELECT numeric_id FROM accounts WHERE alias_of = (SELECT username FROM accounts WHERE numeric_id = ?) OR (alias_of != '' AND alias_of = (SELECT alias_of FROM accounts WHERE numeric_id = ?)))", accountID, accountID, accountID)
}

func (r *Repository) FetchArticles(accountID, filter string, limit, offset int) ([]models.Article, error) {
	query := r.db.Model(&models.Article{}).Preload("Account").Preload("Account.ProfileHistory").Preload("Media").Preload("Media.Variants").Preload("UrlRedirects").Order("created_at DESC")
	query = query.Where("is_trash = ? AND (is_repost = ? OR reply_to_handle IN (SELECT value FROM whitelists WHERE is_active = ?))", false, false, true)
	query = r.applyAccountFilter(query, accountID)
	switch filter {
	case "reposts": query = query.Where("is_repost = ?", true)
	case "media": query = query.Joins("JOIN media ON media.article_id = articles.id").Group("articles.id")
	case "bookmarks": query = query.Where("is_liked = ?", true)
	}
	var articles []models.Article
	return articles, query.Limit(limit).Offset(offset).Find(&articles).Error
}

func (r *Repository) UpsertArticleTx(art *models.Article) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if art.Account.NumericID != "" {
			if err := tx.Save(&art.Account).Error; err != nil { return err }
		}
		if err := tx.Save(art).Error; err != nil { return err }
		for i := range art.Media {
			art.Media[i].ArticleID = art.ID
			if err := tx.Save(&art.Media[i]).Error; err != nil { return err }
		}
		for i := range art.UrlRedirects {
			art.UrlRedirects[i].ArticleID = art.ID
			if err := tx.Save(&art.UrlRedirects[i]).Error; err != nil { return err }
		}
		return nil
	})
}

func (r *Repository) SearchArticles(searchQuery, accountID, filter string, limit, offset int) ([]models.Article, int64, error) {
	query := r.db.Model(&models.Article{}).Preload("Account").Preload("Account.ProfileHistory").Preload("Media").Preload("Media.Variants").Preload("UrlRedirects")
	if searchQuery != "" {
		pat := "%" + searchQuery + "%"
		query = query.Where("full_text LIKE ? OR full_text_ja LIKE ? OR full_text_en LIKE ? OR full_text_zh LIKE ? OR id LIKE ?", pat, pat, pat, pat, pat)
	}
	query = r.applyAccountFilter(query, accountID)
	if filter == "trash" {
		query = query.Where("is_trash = ?", true)
	} else if filter == "all_with_trash" {
		// ゴミ箱データも含めて全件検索
	} else {
		query = query.Where("is_trash = ?", false)
		switch filter {
		case "reposts": query = query.Where("is_repost = ?", true)
		case "media": query = query.Joins("JOIN media ON media.article_id = articles.id").Group("articles.id")
		case "bookmarks": query = query.Where("is_liked = ?", true)
		}
	}
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	var articles []models.Article
	if limit <= 0 { limit = 20 }
	return articles, total, query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&articles).Error
}


func (r *Repository) GetArticleByID(id string) (*models.Article, error) {
	var a models.Article
	err := r.db.Preload("Account").Preload("Account.ProfileHistory").Preload("Media").Preload("Media.Variants").Preload("UrlRedirects").Where("id = ?", id).First(&a).Error
	return &a, err
}

func (r *Repository) GetArticlesByConversationID(convID string) ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Account").Preload("Account.ProfileHistory").Preload("Media").Preload("Media.Variants").Preload("UrlRedirects").Where("conversation_id = ?", convID).Order("created_at ASC").Find(&articles).Error
	return articles, err
}

func (r *Repository) UpdateArticleTranslations(id string, ja, en, zh string) error {
	up := map[string]interface{}{}
	if ja != "" { up["full_text_ja"] = ja } else { up["full_text_ja"] = nil }
	if en != "" { up["full_text_en"] = en } else { up["full_text_en"] = nil }
	if zh != "" { up["full_text_zh"] = zh } else { up["full_text_zh"] = nil }
	return r.db.Model(&models.Article{}).Where("id = ?", id).Updates(up).Error
}
