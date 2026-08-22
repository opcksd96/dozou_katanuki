// app_rpc_stash.go (100行以下)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"dozou_katanuki/models"
)

func (a *App) getStashGraphQLEndpoint() string {
	port := 9999
	if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Network.StashPort > 0 {
		port = cfg.Network.StashPort
	}
	return fmt.Sprintf("http://127.0.0.1:%d/graphql", port)
}

func (a *App) queryStashGraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	endpoint := a.getStashGraphQLEndpoint()
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("stash graphql connection error: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data   map[string]interface{} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", result.Errors[0].Message)
	}
	return result.Data, nil
}

// GetStashMetadata は Stash の GraphQL API から Scene または Image の詳細メタデータを取得します
func (a *App) GetStashMetadata(sceneID, imageID string) (*models.StashMetadataResult, error) {
	if sceneID != "" {
		return a.fetchStashScene(sceneID)
	}
	if imageID != "" {
		return a.fetchStashImage(imageID)
	}
	return nil, fmt.Errorf("neither sceneID nor imageID provided")
}

func (a *App) fetchStashScene(id string) (*models.StashMetadataResult, error) {
	q := `query FindScene($id: ID!) {
		findScene(id: $id) {
			id title details url date rating100
			files { path size duration video_codec audio_codec width height frame_rate bit_rate }
			tags { id name } studio { name }
		}
	}`
	data, err := a.queryStashGraphQL(q, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	sceneMap, ok := data["findScene"].(map[string]interface{})
	if !ok || sceneMap == nil {
		return nil, fmt.Errorf("scene not found: %s", id)
	}
	return parseStashMetadata(sceneMap, true), nil
}

func (a *App) fetchStashImage(id string) (*models.StashMetadataResult, error) {
	q := `query FindImage($id: ID!) {
		findImage(id: $id) {
			id title details url date rating100
			files { path size width height }
			tags { id name } studio { name }
		}
	}`
	data, err := a.queryStashGraphQL(q, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	imgMap, ok := data["findImage"].(map[string]interface{})
	if !ok || imgMap == nil {
		return nil, fmt.Errorf("image not found: %s", id)
	}
	return parseStashMetadata(imgMap, false), nil
}

func parseStashMetadata(m map[string]interface{}, isScene bool) *models.StashMetadataResult {
	res := &models.StashMetadataResult{
		ID:        getString(m, "id"),
		IsScene:   isScene,
		Title:     getString(m, "title"),
		Details:   getString(m, "details"),
		URL:       getString(m, "url"),
		Date:      getString(m, "date"),
		Rating100: int(getFloat(m, "rating100")),
	}
	if studioMap, ok := m["studio"].(map[string]interface{}); ok {
		res.Studio = getString(studioMap, "name")
	}
	if filesList, ok := m["files"].([]interface{}); ok {
		for _, fItem := range filesList {
			if fMap, ok := fItem.(map[string]interface{}); ok {
				res.Files = append(res.Files, models.StashFileDetails{
					Path:       getString(fMap, "path"),
					Size:       int64(getFloat(fMap, "size")),
					Duration:   getFloat(fMap, "duration"),
					VideoCodec: getString(fMap, "video_codec"),
					AudioCodec: getString(fMap, "audio_codec"),
					Width:      int(getFloat(fMap, "width")),
					Height:     int(getFloat(fMap, "height")),
					FrameRate:  getFloat(fMap, "frame_rate"),
					BitRate:    int64(getFloat(fMap, "bit_rate")),
				})
			}
		}
	}
	if tagsList, ok := m["tags"].([]interface{}); ok {
		for _, tItem := range tagsList {
			if tMap, ok := tItem.(map[string]interface{}); ok {
				res.Tags = append(res.Tags, models.StashTag{
					ID:   getString(tMap, "id"),
					Name: getString(tMap, "name"),
				})
			}
		}
	}
	return res
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}

// UpdateStashMetadata は Stash の GraphQL API を通じてメタデータを安全に更新します
func (a *App) UpdateStashMetadata(isScene bool, id, title, details string, rating100 int) (*models.StashMetadataResult, error) {
	if isScene {
		m := `mutation SceneUpdate($input: SceneUpdateInput!) {
			sceneUpdate(input: $input) {
				id title details url date rating100
				files { path size duration video_codec audio_codec width height frame_rate bit_rate }
				tags { id name } studio { name }
			}
		}`
		input := map[string]interface{}{
			"id":        id,
			"title":     title,
			"details":   details,
			"rating100": rating100,
		}
		data, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
		if err != nil {
			return nil, err
		}
		if sMap, ok := data["sceneUpdate"].(map[string]interface{}); ok {
			return parseStashMetadata(sMap, true), nil
		}
	} else {
		m := `mutation ImageUpdate($input: ImageUpdateInput!) {
			imageUpdate(input: $input) {
				id title details url date rating100
				files { path size width height }
				tags { id name } studio { name }
			}
		}`
		input := map[string]interface{}{
			"id":        id,
			"title":     title,
			"details":   details,
			"rating100": rating100,
		}
		data, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
		if err != nil {
			return nil, err
		}
		if imgMap, ok := data["imageUpdate"].(map[string]interface{}); ok {
			return parseStashMetadata(imgMap, false), nil
		}
	}
	return a.GetStashMetadata(id, id)
}

// TriggerStashScan は Stash の GraphQL API を通じてメタデータスキャンをトリガーします
func (a *App) TriggerStashScan(paths []string) (bool, error) {
	if err := a.waitForReady(); err != nil {
		return false, err
	}
	input := map[string]interface{}{"rescan": false}
	if len(paths) > 0 {
		input["paths"] = paths
	}
	m := `mutation ScanMetadata($input: ScanMetadataInput!) {
		metadataScan(input: $input)
	}`
	_, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReconcileStashMedia は Stashapp 内の全メディアから標準タイトルおよびファイル名を走査し、DB の media レコードへ逆引き自動バインドします (SPEC-STORAGE-001)
func (a *App) ReconcileStashMedia() (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	q := `query {
		allScenes { id title files { path } }
		allImages { id title files { path } }
	}`
	data, err := a.queryStashGraphQL(q, nil)
	if err != nil {
		return 0, err
	}

	titleRegex := regexp.MustCompile(`^([A-Za-z0-9_]+)\s\(@([A-Za-z0-9_]+)\):\s([A-Za-z]+)\s([0-9]+)$`)
	boundCount := 0

	db := a.repo.DB()
	if db == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	// 1. Process Scenes
	if scnList, ok := data["allScenes"].([]interface{}); ok {
		for _, item := range scnList {
			if m, ok := item.(map[string]interface{}); ok {
				title := getString(m, "title")
				sID := getString(m, "id")
				if matches := titleRegex.FindStringSubmatch(title); len(matches) == 5 {
					postID := matches[4]
					res := db.Exec("UPDATE media SET stash_scene_id = ?, download_status = 'COMPLETED' WHERE article_id = ? AND stash_scene_id IS NULL", sID, postID)
					if res.Error == nil {
						boundCount += int(res.RowsAffected)
					}
				}
				if filesList, ok := m["files"].([]interface{}); ok {
					for _, fItem := range filesList {
						if fMap, ok := fItem.(map[string]interface{}); ok {
							p := getString(fMap, "path")
							if p != "" {
								base := filepath.Base(p)
								res := db.Exec("UPDATE media SET stash_scene_id = ?, download_status = 'COMPLETED' WHERE media_id = ? AND stash_scene_id IS NULL", sID, base)
								if res.Error == nil {
									boundCount += int(res.RowsAffected)
								}
							}
						}
					}
				}
			}
		}
	}

	// 2. Process Images
	if imgList, ok := data["allImages"].([]interface{}); ok {
		for _, item := range imgList {
			if m, ok := item.(map[string]interface{}); ok {
				title := getString(m, "title")
				imgID := getString(m, "id")
				if matches := titleRegex.FindStringSubmatch(title); len(matches) == 5 {
					postID := matches[4]
					res := db.Exec("UPDATE media SET stash_image_id = ?, download_status = 'COMPLETED' WHERE article_id = ? AND stash_image_id IS NULL", imgID, postID)
					if res.Error == nil {
						boundCount += int(res.RowsAffected)
					}
				}
				if filesList, ok := m["files"].([]interface{}); ok {
					for _, fItem := range filesList {
						if fMap, ok := fItem.(map[string]interface{}); ok {
							p := getString(fMap, "path")
							if p != "" {
								base := filepath.Base(p)
								res := db.Exec("UPDATE media SET stash_image_id = ?, download_status = 'COMPLETED' WHERE media_id = ? AND stash_image_id IS NULL", imgID, base)
								if res.Error == nil {
									boundCount += int(res.RowsAffected)
								}
							}
						}
					}
				}
			}
		}
	}

	return boundCount, nil
}
