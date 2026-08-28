// app/app_rpc_downloader_retry.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
)

// RetryFailedDownloads unpauses and restarts error/paused tasks in Motrix / Aria2
func (a *App) RetryFailedDownloads() (int, error) {
	raw, err := callMotrixRPC("aria2.tellStopped", []interface{}{0, 50, []string{"gid", "status", "errorMessage"}})
	if err != nil {
		return 0, err
	}

	var res struct {
		Result []struct {
			GID    string `json:"gid"`
			Status string `json:"status"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, err
	}

	retriedCount := 0
	for _, t := range res.Result {
		if t.Status == "paused" || t.Status == "error" {
			if _, err := callMotrixRPC("aria2.unpause", []interface{}{t.GID}); err == nil {
				retriedCount++
			}
		}
	}
	return retriedCount, nil
}
