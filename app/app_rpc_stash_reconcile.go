// app/app_rpc_stash_reconcile.go (100行以下)
package app

// ReconcileStashScenes は Stashapp 内の全Scene(動画)からファイル名を走査し、DB の media レコードへ逆引き自動バインドします
func (a *App) ReconcileStashScenes() (int, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	q := `query { findScenes(filter: {sort: "created_at", direction: DESC, per_page: 200}) { scenes { id title details files { path } } } }`
	data, err := a.queryStashGraphQL(q, nil)
	if err != nil {
		return 0, err
	}

	boundCount := 0
	if findScenes, ok := data["findScenes"].(map[string]interface{}); ok {
		if scnList, ok := findScenes["scenes"].([]interface{}); ok {
			for _, item := range scnList {
				if m, ok := item.(map[string]interface{}); ok {
					if a.bindStashScene(m) {
						boundCount++
					}
				}
			}
		}
	}
	return boundCount, nil
}

// ReconcileStashImages は Stashapp 内の全Image(画像)からファイル名を走査し、DB の media レコードへ逆引き自動バインドします
func (a *App) ReconcileStashImages() (int, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	q := `query { findImages(filter: {sort: "created_at", direction: DESC, per_page: 200}) { images { id title details files { path } } } }`
	data, err := a.queryStashGraphQL(q, nil)
	if err != nil {
		return 0, err
	}

	boundCount := 0
	if findImages, ok := data["findImages"].(map[string]interface{}); ok {
		if imgList, ok := findImages["images"].([]interface{}); ok {
			for _, item := range imgList {
				if m, ok := item.(map[string]interface{}); ok {
					if a.bindStashImage(m) {
						boundCount++
					}
				}
			}
		}
	}
	return boundCount, nil
}
