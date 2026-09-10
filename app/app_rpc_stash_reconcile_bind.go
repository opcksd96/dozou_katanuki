// app/app_rpc_stash_reconcile_bind.go (100行以下)
package app

import (
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

func (a *App) bindStashScene(m map[string]interface{}) bool {
	sID := getString(m, "id")
	db := a.Repo.DB()
	if db == nil || sID == "" {
		return false
	}
	bound, articleID := false, ""
	filesList, ok := m["files"].([]interface{})
	if !ok {
		return false
	}
	for _, fItem := range filesList {
		fMap, ok := fItem.(map[string]interface{})
		if !ok {
			continue
		}
		p := getString(fMap, "path")
		if p == "" {
			continue
		}
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
			if !med.StashSceneID.Valid || med.StashSceneID.String != sID {
				if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{"stash_scene_id": sID, "download_status": "COMPLETED"}).Error == nil {
					bound = true
					break
				}
			}
		}
	}
	if articleID != "" {
		a.syncArticleDetailsToStash(sID, articleID, true)
	}
	return bound
}

func (a *App) bindStashImage(m map[string]interface{}) bool {
	sID := getString(m, "id")
	db := a.Repo.DB()
	if db == nil || sID == "" {
		return false
	}
	bound, articleID := false, ""
	filesList, ok := m["files"].([]interface{})
	if !ok {
		return false
	}
	for _, fItem := range filesList {
		fMap, ok := fItem.(map[string]interface{})
		if !ok {
			continue
		}
		p := getString(fMap, "path")
		if p == "" {
			continue
		}
		base := filepath.Base(p)
		var med models.Media
		if err := db.Where("(media_id = ? OR media_id = ? OR download_url LIKE ?)", base, strings.TrimSuffix(base, filepath.Ext(base)), "%/"+base).First(&med).Error; err == nil {
			articleID = med.ArticleID
			if !med.StashImageID.Valid || med.StashImageID.String != sID {
				if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{"stash_image_id": sID, "download_status": "COMPLETED"}).Error == nil {
					bound = true
					break
				}
			}
		}
	}
	if articleID != "" {
		a.syncArticleDetailsToStash(sID, articleID, false)
	}
	return bound
}
