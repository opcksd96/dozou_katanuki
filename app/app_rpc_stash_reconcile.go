// app/app_rpc_stash_reconcile.go (100行以下)
package app

import (
	"path/filepath"
	"regexp"
	"strings"

	"dozou_katanuki/models"
)

// ReconcileStashScenes は Stashapp 内の全Scene(動画)から標準タイトルおよびファイル名を走査し、DB の media レコードへ逆引き自動バインドします
func (a *App) ReconcileStashScenes() (int, error) {
	if err := a.WaitForReady(); err != nil { return 0, err }
	q := `query { allScenes { id title details files { path } } }`
	data, err := a.queryStashGraphQL(q, nil)
	if err != nil { return 0, err }

	titleRegex := regexp.MustCompile(`^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([0-9]+)$`)
	boundCount := 0
	if scnList, ok := data["allScenes"].([]interface{}); ok {
		for _, item := range scnList {
			if m, ok := item.(map[string]interface{}); ok {
				if a.bindStashScene(m, titleRegex) { boundCount++ }
			}
		}
	}
	return boundCount, nil
}

// ReconcileStashImages は Stashapp 内の全Image(画像)から標準タイトルおよびファイル名を走査し、DB の media レコードへ逆引き自動バインドします
func (a *App) ReconcileStashImages() (int, error) {
	if err := a.WaitForReady(); err != nil { return 0, err }
	q := `query { allImages { id title details files { path } } }`
	data, err := a.queryStashGraphQL(q, nil)
	if err != nil { return 0, err }

	titleRegex := regexp.MustCompile(`^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([0-9]+)$`)
	boundCount := 0
	if imgList, ok := data["allImages"].([]interface{}); ok {
		for _, item := range imgList {
			if m, ok := item.(map[string]interface{}); ok {
				if a.bindStashImage(m, titleRegex) { boundCount++ }
			}
		}
	}
	return boundCount, nil
}

func (a *App) bindStashScene(m map[string]interface{}, titleRegex *regexp.Regexp) bool {
	sID := getString(m, "id")
	db := a.Repo.DB()
	if db == nil { return false }
	bound := false
	var articleID string

	if filesList, ok := m["files"].([]interface{}); ok {
		for _, fItem := range filesList {
			if fMap, ok := fItem.(map[string]interface{}); ok {
				p := getString(fMap, "path")
				if p != "" {
					base := filepath.Base(p)
					cleanBase := strings.TrimSuffix(base, filepath.Ext(base))
					var med models.Media
					if err := db.Where("(media_id = ? OR media_id = ? OR download_url LIKE ?)", base, cleanBase, "%/"+base).First(&med).Error; err == nil {
						articleID = med.ArticleID
					} else {
						var variant models.MediaVariant
						if err := db.Where("variant_hash = ?", cleanBase).First(&variant).Error; err == nil {
							if err := db.Where("media_id = ?", variant.MediaID).First(&med).Error; err == nil {
								articleID = med.ArticleID
							}
						}
					}
					if articleID != "" {
						if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{"stash_scene_id": sID, "download_status": "COMPLETED"}).Error == nil {
							bound = true
							break
						}
					}
				}
			}
		}
	}

	title := getString(m, "title")
	if !bound && title != "" {
		if matches := titleRegex.FindStringSubmatch(title); len(matches) == 5 {
			postID := matches[4]
			articleID = postID
			var med models.Media
			if err := db.Where("article_id = ? AND type != 'image'", postID).First(&med).Error; err == nil {
				if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{"stash_scene_id": sID, "download_status": "COMPLETED"}).Error == nil {
					bound = true
				}
			}
		}
	}

	if articleID != "" { a.syncArticleDetailsToStash(sID, articleID, true) }
	return bound
}

func (a *App) bindStashImage(m map[string]interface{}, titleRegex *regexp.Regexp) bool {
	sID := getString(m, "id")
	db := a.Repo.DB()
	if db == nil { return false }
	bound := false
	var articleID string

	if filesList, ok := m["files"].([]interface{}); ok {
		for _, fItem := range filesList {
			if fMap, ok := fItem.(map[string]interface{}); ok {
				p := getString(fMap, "path")
				if p != "" {
					base := filepath.Base(p)
					var med models.Media
					if err := db.Where("(media_id = ? OR media_id = ? OR download_url LIKE ?)", base, strings.TrimSuffix(base, filepath.Ext(base)), "%/"+base).First(&med).Error; err == nil {
						articleID = med.ArticleID
						if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{"stash_image_id": sID, "download_status": "COMPLETED"}).Error == nil {
							bound = true
							break
						}
					}
				}
			}
		}
	}

	title := getString(m, "title")
	if !bound && title != "" {
		if matches := titleRegex.FindStringSubmatch(title); len(matches) == 5 {
			postID := matches[4]
			articleID = postID
			var med models.Media
			if err := db.Where("article_id = ? AND type = 'image'", postID).First(&med).Error; err == nil {
				if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{"stash_image_id": sID, "download_status": "COMPLETED"}).Error == nil {
					bound = true
				}
			}
		}
	}

	if articleID != "" { a.syncArticleDetailsToStash(sID, articleID, false) }
	return bound
}
