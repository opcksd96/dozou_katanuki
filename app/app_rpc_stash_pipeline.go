// app/app_rpc_stash_pipeline.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"time"
)

// TriggerStashGenerate は サムネイル、VTTスプライト、プレビュー動画の生成タスクをトリガーします
func (a *App) TriggerStashGenerate(sceneIDs, imageIDs []string) (bool, error) {
	if err := a.WaitForReady(); err != nil { return false, err }
	input := map[string]interface{}{
		"sprites":         true,
		"previews":        true,
		"imageThumbnails": true,
		"covers":          true,
	}
	if len(sceneIDs) > 0 { input["sceneIDs"] = sceneIDs }
	if len(imageIDs) > 0 { input["imageIDs"] = imageIDs }

	m := `mutation GenerateMetadata($input: GenerateMetadataInput!) {
		metadataGenerate(input: $input)
	}`
	_, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
	return err == nil, err
}

// TriggerStashScenePipeline は動画専用のパイプライン（スキャン ➔ VTT/Preview生成 ➔ バインド）を実行します
func (a *App) TriggerStashScenePipeline() {
	go func() {
		// 1. スキャン
		_, _ = a.TriggerStashScan(nil)
		time.Sleep(3 * time.Second)

		// 2. 動画特有の生成物 (VTT/Preview)
		_, _ = a.TriggerStashGenerate(nil, []string{}) // trigger for scenes if nil
		time.Sleep(2 * time.Second)

		// 3. 動画専用の逆引きバインド (variants対応)
		_, _ = a.ReconcileStashScenes()
	}()
}

// TriggerStashImagePipeline は画像専用のパイプライン（スキャン ➔ Image生成 ➔ バインド）を実行します
func (a *App) TriggerStashImagePipeline() {
	go func() {
		// 1. スキャン
		_, _ = a.TriggerStashScan(nil)
		time.Sleep(3 * time.Second)

		// 2. 画像特有の生成物 (Thumbnails)
		_, _ = a.TriggerStashGenerate([]string{}, nil) // trigger for images if nil
		time.Sleep(2 * time.Second)

		// 3. 画像専用の逆引きバインド
		_, _ = a.ReconcileStashImages()
	}()
}

// TriggerStashAllPipelines は互換性のために両方のパイプラインを順に呼び出します
func (a *App) TriggerStashAllPipelines() {
	a.TriggerStashScenePipeline()
	// Note: since they are goroutines, they run concurrently.
	// But they both call scan, which might be redundant. For simplicity, we just trigger them.
	a.TriggerStashImagePipeline()
}
