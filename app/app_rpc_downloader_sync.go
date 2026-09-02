// app/app_rpc_downloader_sync.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// SyncCompletedDownloads queries Aria2/Motrix completed/error tasks and updates DB records
func (a *App) SyncCompletedDownloads() (int, error) {
	raw, err := callMotrixRPC("aria2.tellStopped", []interface{}{0, 50, []string{"gid", "status", "files", "errorMessage"}})
	if err != nil { return 0, err }

	var res struct {
		Result []struct {
			GID, Status, ErrorMessage string
			Files []struct { Path string `json:"path"` } `json:"files"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil { return 0, err }

	syncedCount := 0

	for _, t := range res.Result {
		fileName := ""
		if len(t.Files) > 0 && t.Files[0].Path != "" {
			fn := filepath.Base(t.Files[0].Path)
			ext := filepath.Ext(fn)
			baseName := strings.TrimSuffix(fn, ext)
			for _, sfx := range []string{"_motrix", "_requests", "_thunder"} { baseName = strings.TrimSuffix(baseName, sfx) }
			fileName = baseName + ext
		}
		if t.Status == "complete" && fileName != "" {
			if a.Repo != nil && a.Repo.RegisterCompletedMediaFile(t.Files[0].Path) == nil {
				_ = a.Repo.MarkTaskCompleted(fileName)
				a.AppendPipelineLog("MOTRIX", "SUCCESS", fmt.Sprintf("✅ Motrix完了回収: %s (%s)", fileName, t.Files[0].Path))
				syncedCount++
				a.TriggerStashAllPipelines()
			}
			_, _ = callMotrixRPC("aria2.removeDownloadResult", []interface{}{t.GID})
		} else if t.Status == "error" && fileName != "" {
			_, _ = callMotrixRPC("aria2.removeDownloadResult", []interface{}{t.GID})
			if a.Repo != nil {
				_ = a.Repo.UpdateMediaMetadata(fileName, "ESCALATED", "", "", "Motrixエラー("+t.ErrorMessage+")→迅雷エスカレーション待機")
				_ = a.Repo.UpdateMediaCheckpointTime(fileName, models.StageThunder)
				a.AppendPipelineLog("MOTRIX", "WARN", fmt.Sprintf("⚠️ Motrix失敗 (%s) ➔ 迅雷エスカレーション待機: %s", t.ErrorMessage, fileName))
				syncedCount++
			}
		}
	}
	return syncedCount, nil
}
