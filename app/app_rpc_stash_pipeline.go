// app/app_rpc_stash_pipeline.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"time"
)

// TriggerStashGenerate は サムネイル、VTTスプライト、プレビュー動画の生成タスクをトリガーします
func (a *App) TriggerStashGenerate(sceneIDs, imageIDs []string) (bool, error) {
	if err := a.WaitForReady(); err != nil {
		return false, err
	}
	input := map[string]interface{}{
		"sprites":         true,
		"previews":        true,
		"imageThumbnails": true,
		"covers":          true,
	}
	if len(sceneIDs) > 0 {
		input["sceneIDs"] = sceneIDs
	}
	if len(imageIDs) > 0 {
		input["imageIDs"] = imageIDs
	}

	m := `mutation GenerateMetadata($input: GenerateMetadataInput!) {
		metadataGenerate(input: $input)
	}`
	_, err := a.queryStashGraphQL(m, map[string]interface{}{"input": input})
	return err == nil, err
}

// TriggerStashPipelineForPaths は指定パス（または全体）に対してスキャンからバインドまでを順序正しく実行します
func (a *App) TriggerStashPipelineForPaths(paths []string) {
	go func() {
		// 1. スキャン（多重呼び出しを避けて1回だけ発行）
		_, _ = a.TriggerStashScan(paths)
		time.Sleep(3 * time.Second)

		// 2. 生成タスク (Thumbnails, Previews, Sprites)
		// FIXME: 12秒ごとの自動パイプラインでStash全件に対して生成タスクが走るとハングアップの原因になるため、自動トリガーは無効化
		// _, _ = a.TriggerStashGenerate(nil, nil)

		// 3. バインド（スキャンの非同期完了ラグを吸収するためインターバル付きで試行）
		for i := 0; i < 4; i++ {
			time.Sleep(2 * time.Second)
			total, _ := a.ReconcileStashMedia()
			if total > 0 {
				break
			}
		}
	}()
}

// TriggerStashScenePipeline は動画パイプラインを順序正しく実行します
func (a *App) TriggerStashScenePipeline() {
	a.TriggerStashPipelineForPaths(nil)
}

// TriggerStashImagePipeline は画像パイプラインを順序正しく実行します
func (a *App) TriggerStashImagePipeline() {
	a.TriggerStashPipelineForPaths(nil)
}

// TriggerStashAllPipelines は互換性のために全体スキャン＆照合を実行します
func (a *App) TriggerStashAllPipelines() {
	a.TriggerStashPipelineForPaths(nil)
}
