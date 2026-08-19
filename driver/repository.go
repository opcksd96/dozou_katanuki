// driver/repository.go
package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
		Preload("UrlRedirects").
		Order("created_at DESC")

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

// UpsertArticleTx は記事・アカウント・メディア・短縮URLをトランザクション内で保存します
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
		for i := range art.UrlRedirects {
			art.UrlRedirects[i].ArticleID = art.ID
			if err := tx.Save(&art.UrlRedirects[i]).Error; err != nil {
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
		Preload("Media").
		Preload("UrlRedirects")

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
		Preload("UrlRedirects").
		Where("id = ?", id).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticlesByConversationID は指定されたスレッド(conversation_id)に属する全記事を時系列順に取得します
func (r *Repository) GetArticlesByConversationID(conversationID string) ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Account").
		Preload("Account.ProfileHistory").
		Preload("Media").
		Preload("UrlRedirects").
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&articles).Error
	return articles, err
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

// ResetMediaStatus は指定されたメディアのダウンロードステータスを QUEUED にリセットします
func (r *Repository) ResetMediaStatus(mediaID string) error {
	return r.db.Model(&models.Media{}).Where("media_id = ?", mediaID).Updates(map[string]interface{}{
		"download_status": "QUEUED",
		"failed_reason":   nil,
	}).Error
}

// GetMediaByID は指定されたメディアIDのレコードを取得します
func (r *Repository) GetMediaByID(mediaID string) (*models.Media, error) {
	var m models.Media
	err := r.db.Where("media_id = ?", mediaID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// BackupDatabase は SQLite の VACUUM INTO を用いてオンラインバックアップを作成し、世代管理上限を超えた古いファイルを削除します
func (r *Repository) BackupDatabase(destDir string, maxGenerations int) (string, error) {
	if destDir == "" {
		destDir = filepath.Join("backups", "database")
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup dir: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("archive_%s.db", timestamp)
	destPath := filepath.Join(destDir, backupFileName)

	// SQLite の VACUUM INTO コマンドはスラッシュ区切りのパスを安全に受け付けます
	normalizedPath := filepath.ToSlash(destPath)
	if err := r.db.Exec(fmt.Sprintf("VACUUM INTO '%s'", normalizedPath)).Error; err != nil {
		return "", fmt.Errorf("failed to execute VACUUM INTO: %w", err)
	}

	if maxGenerations > 0 {
		_ = r.purgeOldBackups(destDir, maxGenerations)
	}

	return destPath, nil
}

func (r *Repository) purgeOldBackups(destDir string, maxGenerations int) error {
	entries, err := os.ReadDir(destDir)
	if err != nil {
		return err
	}

	var backupFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "archive_") && strings.HasSuffix(entry.Name(), ".db") {
			backupFiles = append(backupFiles, filepath.Join(destDir, entry.Name()))
		}
	}

	if len(backupFiles) <= maxGenerations {
		return nil
	}

	// 昇順（古い順）にソート
	sort.Strings(backupFiles)
	toDelete := len(backupFiles) - maxGenerations
	for i := 0; i < toDelete; i++ {
		_ = os.Remove(backupFiles[i])
	}
	return nil
}

// AuditDatabase は SQLite3 整合性監査 (PRAGMA) および孤立メディア・ファイルを総合検査します
func (r *Repository) AuditDatabase(stashDir, blobsDir string) (*models.AuditReport, error) {
	report := &models.AuditReport{
		ExecutedAt: time.Now(),
	}

	// 1. PRAGMA integrity_check
	integrityOK, integrityMsgs, err := RunIntegrityCheck(r.db)
	if err != nil {
		report.IntegrityOK = false
		report.IntegrityErrors = []string{err.Error()}
	} else {
		report.IntegrityOK = integrityOK
		report.IntegrityErrors = integrityMsgs
	}

	// 2. PRAGMA foreign_key_check
	fkViolations, err := RunForeignKeyCheck(r.db)
	if err != nil {
		report.ForeignKeyOK = false
	} else {
		report.ForeignKeyOK = len(fkViolations) == 0
		report.ForeignKeyErrors = fkViolations
	}

	// 3. 孤立 DB メディア検出
	orphanMedia, err := FindOrphanDBMedia(r.db, stashDir, blobsDir)
	if err == nil {
		report.OrphanDBMedia = orphanMedia
	}

	// 4. 孤立ファイル検出
	knownKeys, err := GetKnownMediaIdentifiers(r.db)
	if err == nil {
		orphanFiles, err := ScanOrphanFiles(stashDir, blobsDir, knownKeys)
		if err == nil {
			report.OrphanFiles = orphanFiles
		}
	}

	// サマリー生成
	summary := "健全"
	if !report.IntegrityOK {
		summary = "要修復 (SQLite ページ/インデックス破損検知)"
	} else if !report.ForeignKeyOK {
		summary = fmt.Sprintf("注意 (外部キー違反: %d件)", len(report.ForeignKeyErrors))
	} else if len(report.OrphanDBMedia) > 0 || len(report.OrphanFiles) > 0 {
		summary = fmt.Sprintf("要クレンジング (孤立DB:%d件, 孤立ファイル:%d件)", len(report.OrphanDBMedia), len(report.OrphanFiles))
	}
	report.Summary = summary

	return report, nil
}

// PurgeOrphanFiles は指定された孤立ファイルをOSのごみ箱へ退避します
func (r *Repository) PurgeOrphanFiles(paths []string) (int, error) {
	count := 0
	for _, p := range paths {
		if err := MoveToRecycleBin(p); err == nil {
			count++
		}
	}
	return count, nil
}

// PurgeOrphanDBMedia は指定された media_id の DB レコードを物理削除します
func (r *Repository) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	return DeleteDBMediaByIDs(r.db, mediaIDs)
}

