// app/app_rpc_stash_mutate.go (100行以下)
package app

import (
	"dozou_katanuki/models"
)

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
	if err := a.WaitForReady(); err != nil {
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

// DestroyStashScene は Stash の GraphQL API を通じて Scene を削除します (ファイル削除なし)
func (a *App) DestroyStashScene(id string) error {
	if err := a.WaitForReady(); err != nil { return err }
	m := `mutation SceneDestroy($input: SceneDestroyInput!) { sceneDestroy(input: $input) }`
	input := map[string]interface{}{ "id": id, "delete_file": false, "delete_generated": true }
	_, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
	return err
}

// DestroyStashImage は Stash の GraphQL API を通じて Image を削除します (ファイル削除なし)
func (a *App) DestroyStashImage(id string) error {
	if err := a.WaitForReady(); err != nil { return err }
	m := `mutation ImageDestroy($input: ImageDestroyInput!) { imageDestroy(input: $input) }`
	input := map[string]interface{}{ "id": id, "delete_file": false, "delete_generated": true }
	_, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
	return err
}
