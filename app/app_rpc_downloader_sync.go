// app/app_rpc_downloader_sync.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

// SyncCompletedDownloads queries Aria2/Motrix completed/error tasks and updates DB records
func (a *App) SyncCompletedDownloads() (int, error) {
	raw, err := callMotrixRPC("aria2.tellStopped", []interface{}{0, 50, []string{"gid", "status", "files", "errorMessage"}})
	if err != nil {
		return 0, err
	}

	var res struct {
		Result []struct {
			GID, Status, ErrorMessage string
			Files                     []struct {
				Path string `json:"path"`
			} `json:"files"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, err
	}

	syncedCount := 0

	for _, t := range res.Result {
		filePath, fileName := "", ""
		if len(t.Files) > 0 && t.Files[0].Path != "" {
			filePath = t.Files[0].Path
			fn := filepath.Base(filePath)
			ext := filepath.Ext(fn)
			baseName := strings.TrimSuffix(fn, ext)
			for _, sfx := range []string{"_motrix", "_requests", "_thunder"} {
				baseName = strings.TrimSuffix(baseName, sfx)
			}
			fileName = baseName + ext
		}
		if t.Status == "complete" && fileName != "" {
			if a.Repo != nil {
				if err := a.Repo.RegisterCompletedMediaFile(filePath); err == nil {
					_ = a.Repo.MarkTaskCompleted(fileName)
					a.AppendPipelineLog("MOTRIX", "SUCCESS", fmt.Sprintf("✅ Motrix完了回収: %s (%s)", fileName, filePath))
					syncedCount++
					a.TriggerStashPipelineForPaths([]string{filepath.Dir(filePath)})
				} else {
					a.AppendPipelineLog("MOTRIX", "WARN", fmt.Sprintf("⚠️ Motrix偽ファイル検知 (%v): %s ➔ 迅雷エスカレーション", err, fileName))
					_ = driver.MoveToRecycleBin(filePath)
					cleanID := strings.TrimSuffix(fileName, filepath.Ext(fileName))
					var med models.Media
					if errDB := a.Repo.DB().Where("media_id = ? OR media_id = ? OR download_url LIKE ?", fileName, cleanID, "%/"+fileName).First(&med).Error; errDB == nil {
						_ = a.Repo.UpdateMediaMetadata(med.MediaID, "ESCALATED", "", "", "偽ファイル検知("+err.Error()+")→迅雷エスカレーション")
						_ = a.Repo.UpdateMediaCheckpointTime(med.MediaID, models.StageThunder)
						_, _ = a.EscalateToThunder(med.MediaID, med.DownloadURL)
					}
					syncedCount++
				}
			}
			_, _ = callMotrixRPC("aria2.removeDownloadResult", []interface{}{t.GID})
		} else if t.Status == "error" && fileName != "" {
			_, _ = callMotrixRPC("aria2.removeDownloadResult", []interface{}{t.GID})
			if a.Repo != nil {
				cleanID := strings.TrimSuffix(fileName, filepath.Ext(fileName))
				var med models.Media
				if err := a.Repo.DB().Where("media_id = ? OR media_id = ? OR download_url LIKE ?", fileName, cleanID, "%/"+fileName).First(&med).Error; err == nil {
					_ = a.Repo.UpdateMediaMetadata(med.MediaID, "ESCALATED", "", "", "Motrixエラー("+t.ErrorMessage+")→迅雷自動エスカレーション")
					_ = a.Repo.UpdateMediaCheckpointTime(med.MediaID, models.StageThunder)
					a.AppendPipelineLog("MOTRIX", "WARN", fmt.Sprintf("⚠️ Motrix失敗 (%s) ➔ 迅雷エスカレーション投入: %s", t.ErrorMessage, med.MediaID))
					_, _ = a.EscalateToThunder(med.MediaID, med.DownloadURL)
				} else {
					_ = a.Repo.UpdateMediaMetadata(fileName, "ESCALATED", "", "", "Motrixエラー("+t.ErrorMessage+")→迅雷エスカレーション待機")
					_ = a.Repo.UpdateMediaCheckpointTime(fileName, models.StageThunder)
					a.AppendPipelineLog("MOTRIX", "WARN", fmt.Sprintf("⚠️ Motrix失敗 (%s) ➔ 迅雷エスカレーション待機: %s", t.ErrorMessage, fileName))
					_, _ = a.EscalateToThunder(fileName, "")
				}
				syncedCount++
			}
		}
	}
	return syncedCount, nil
}
