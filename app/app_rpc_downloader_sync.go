// app/app_rpc_downloader_sync.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
)

// SyncCompletedDownloads queries Aria2/Motrix completed tasks and updates DB records
func (a *App) SyncCompletedDownloads() (int, error) {
	raw, err := callMotrixRPC("aria2.tellStopped", []interface{}{0, 50, []string{"gid", "status", "files"}})
	if err != nil {
		return 0, err
	}

	var res struct {
		Result []struct {
			GID    string `json:"gid"`
			Status string `json:"status"`
			Files  []struct {
				Path string `json:"path"`
			} `json:"files"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, err
	}

	syncedCount := 0
	for _, t := range res.Result {
		if t.Status == "complete" && len(t.Files) > 0 {
			filePath := t.Files[0].Path
			if a.Repo != nil {
				if err := a.Repo.RegisterCompletedMediaFile(filePath); err == nil {
					syncedCount++
				}
			}
		}
	}
	return syncedCount, nil
}
