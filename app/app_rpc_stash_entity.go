// app/app_rpc_stash_entity.go (100行以下)
package app

import (
	"fmt"
	"regexp"
	"strings"

	"dozou_katanuki/models"
)

func (a *App) syncArticleDetailsToStash(stashID, articleID string, isScene bool) {
	db := a.Repo.DB()
	var art models.Article
	if err := db.Preload("Account").Where("id = ?", articleID).First(&art).Error; err != nil { return }

	textToSync := art.FullText
	if art.FullTextJA.Valid && art.FullTextJA.String != "" { textToSync = art.FullTextJA.String + "\n\n" + art.FullText }
	uname, dname := art.Account.Username, art.Account.DisplayName

	urlsList := []string{fmt.Sprintf("https://twitter.com/%s/status/%s", uname, articleID)}
	if art.WaybackURL != "" { urlsList = append(urlsList, art.WaybackURL) }
	urlsList = append(urlsList, fmt.Sprintf("http://localhost:9999/plugin/x-timeline-middleware/index.html?view=x-timeline&performer=%s&jump_to_tweet=%s", uname, articleID))

	cFieldsMap := map[string]interface{}{
		"tweet_id": articleID, "original_name": dname, "source_system": "X_Wayback",
		"wayback_url": "[]", "dead_media": "[]",
	}
	if art.WaybackURL != "" {
		cFieldsMap["wayback_url"] = fmt.Sprintf("[\"%s\"]", art.WaybackURL)
		if match := regexp.MustCompile(`/web/(\d{14})`).FindStringSubmatch(art.WaybackURL); len(match) == 2 {
			cFieldsMap["wayback_timestamp"] = match[1]
		}
	}

	input := map[string]interface{}{
		"id": stashID, "title": fmt.Sprintf("X (@%s): Tweet %s", uname, articleID),
		"details": textToSync, "urls": urlsList, "url": urlsList[0], "date": art.CreatedAt.Format("2006-01-02"),
		"custom_fields": map[string]interface{}{"partial": cFieldsMap},
	}
	if sIDObj := a.findOrCreateStudio(uname); sIDObj != "" { input["studio_id"] = sIDObj }
	if pIDObj := a.findOrCreatePerformer(uname, dname); pIDObj != "" { input["performer_ids"] = []string{pIDObj} }

	mut := `mutation ImageUpdate($input: ImageUpdateInput!) { imageUpdate(input: $input) { id } }`
	if isScene { mut = `mutation SceneUpdate($input: SceneUpdateInput!) { sceneUpdate(input: $input) { id } }` }
	_, _ = a.queryStashGraphQL(mut, map[string]interface{}{"input": input})
}

func (a *App) findOrCreateStudio(name string) string {
	if name == "" { return "" }
	q := `query FindStudios($q: String!) { findStudios(filter: { q: $q, per_page: 1 }) { studios { id name } } }`
	if data, err := a.queryStashGraphQL(q, map[string]interface{}{"q": name}); err == nil && data != nil {
		if fMap, ok := data["findStudios"].(map[string]interface{}); ok {
			if sList, ok := fMap["studios"].([]interface{}); ok {
				for _, sItem := range sList {
					if sm, ok := sItem.(map[string]interface{}); ok && strings.EqualFold(getString(sm, "name"), name) {
						return getString(sm, "id")
					}
				}
			}
		}
	}
	cMut := `mutation StudioCreate($input: StudioCreateInput!) { studioCreate(input: $input) { id } }`
	if cData, err := a.queryStashGraphQL(cMut, map[string]interface{}{"input": map[string]interface{}{"name": name}}); err == nil && cData != nil {
		if cMap, ok := cData["studioCreate"].(map[string]interface{}); ok { return getString(cMap, "id") }
	}
	return ""
}

func (a *App) findOrCreatePerformer(name, disambiguation string) string {
	if name == "" { return "" }
	q := `query FindPerformers($q: String!) { findPerformers(filter: { q: $q, per_page: 1 }) { performers { id name } } }`
	if data, err := a.queryStashGraphQL(q, map[string]interface{}{"q": name}); err == nil && data != nil {
		if fMap, ok := data["findPerformers"].(map[string]interface{}); ok {
			if pList, ok := fMap["performers"].([]interface{}); ok {
				for _, pItem := range pList {
					if pm, ok := pItem.(map[string]interface{}); ok && strings.EqualFold(getString(pm, "name"), name) {
						return getString(pm, "id")
					}
				}
			}
		}
	}
	cMut := `mutation PerformerCreate($input: PerformerCreateInput!) { performerCreate(input: $input) { id } }`
	inp := map[string]interface{}{"name": name}
	if disambiguation != "" { inp["disambiguation"] = disambiguation }
	if cData, err := a.queryStashGraphQL(cMut, map[string]interface{}{"input": inp}); err == nil && cData != nil {
		if cMap, ok := cData["performerCreate"].(map[string]interface{}); ok { return getString(cMap, "id") }
	}
	return ""
}
