// app/app_rpc_motrix_batch.go (100行以下 - SPEC-PRINCIPLE-001)
package app

// BatchControlMotrix は複数のタスクに対して一括操作を実行します
func (a *App) BatchControlMotrix(action string, gids []string) (int, error) {
	successCount := 0
	for _, gid := range gids {
		if gid == "" {
			continue
		}
		var err error
		switch action {
		case "pause":
			_, err = callMotrixRPC("aria2.pause", []interface{}{gid})
			if err != nil {
				_, err = callMotrixRPC("aria2.forcePause", []interface{}{gid})
			}
		case "unpause":
			_, err = callMotrixRPC("aria2.unpause", []interface{}{gid})
		case "remove":
			_, err = callMotrixRPC("aria2.remove", []interface{}{gid})
			if err != nil {
				_, err = callMotrixRPC("aria2.removeDownloadResult", []interface{}{gid})
			}
		}
		if err == nil {
			successCount++
		}
	}
	return successCount, nil
}
