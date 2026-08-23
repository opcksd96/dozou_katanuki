// app/app_rpc_stash_query.go (100行以下)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

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
