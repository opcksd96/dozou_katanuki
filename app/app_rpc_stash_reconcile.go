// app/app_rpc_stash_reconcile.go (100行以下)
package app

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"dozou_katanuki/models"
)

// ReconcileStashMedia は Stashapp 内の全メディアから標準タイトルおよびファイル名を走査し、DB の media レコードへ逆引き自動バインドします (SPEC-STORAGE-001)
func (a *App) ReconcileStashMedia() (int, error) {
	if err := a.WaitForReady(); err != nil { return 0, err }
	q := `query {
		allScenes { id title details files { path } }
		allImages { id title details files { path } }
	}`
	data, err := a.queryStashGraphQL(q, nil)
	if err != nil { return 0, err }

	titleRegex := regexp.MustCompile(`^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([0-9]+)$`)
	boundCount := 0
	db := a.Repo.DB()
	if db == nil { return 0, fmt.Errorf("database not initialized") }

	if scnList, ok := data["allScenes"].([]interface{}); ok {
		for _, item := range scnList {
			if m, ok := item.(map[string]interface{}); ok {
				if a.bindStashItem(m, true, titleRegex) { boundCount++ }
			}
		}
	}

	if imgList, ok := data["allImages"].([]interface{}); ok {
		for _, item := range imgList {
			if m, ok := item.(map[string]interface{}); ok {
				if a.bindStashItem(m, false, titleRegex) { boundCount++ }
			}
		}
	}
	return boundCount, nil
}

func (a *App) bindStashItem(m map[string]interface{}, isScene bool, titleRegex *regexp.Regexp) bool {
	sID := getString(m, "id")
	db := a.Repo.DB()
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
						field := "stash_image_id"
						if isScene { field = "stash_scene_id" }
						if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{field: sID, "download_status": "COMPLETED"}).Error == nil {
							bound = true
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
			typeCond := "type = 'image'"
			field := "stash_image_id"
			if isScene { typeCond = "type != 'image'"; field = "stash_scene_id" }
			if err := db.Where("article_id = ? AND "+typeCond, postID).First(&med).Error; err == nil {
				if db.Model(&models.Media{}).Where("media_id = ?", med.MediaID).Updates(map[string]interface{}{field: sID, "download_status": "COMPLETED"}).Error == nil {
					bound = true
				}
			}
		}
	}

	if articleID != "" {
		a.syncArticleDetailsToStash(sID, articleID, isScene)
	}
	return bound
}
