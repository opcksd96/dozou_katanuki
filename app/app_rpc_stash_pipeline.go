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

// TriggerStashFullPipeline は スキャン ➔ 生成 ➔ メタデータ充填/バインドの一連のパイプラインをバックグラウンド実行します
func (a *App) TriggerStashFullPipeline() {
	go func() {
		// 1. Stash へファイルスキャンを指示
		_, _ = a.TriggerStashScan(nil)
		time.Sleep(3 * time.Second)

		// 2. VTTスプライト・サムネイル・プレビュー動画を自動生成
		_, _ = a.TriggerStashGenerate(nil, nil)
		time.Sleep(2 * time.Second)

		// 3. Stash の ID と dozou の Media/Article をバインドしてツイートメタデータを充填
		_, _ = a.ReconcileStashMedia()
	}()
}
