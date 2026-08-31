// app/app_rpc_pipeline_dispatch.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// DispatchQueuedToMotrix は QUEUED のメディアを Motrix (Aria2) へ一括投入して第1ステージを開始します
func (a *App) DispatchQueuedToMotrix() (int, error) {
	if a.Repo == nil || a.Repo.DB() == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	db := a.Repo.DB()

	var queuedMedias []models.Media
	if err := db.Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Find(&queuedMedias).Error; err != nil {
		return 0, err
	}

	destRoot := a.getMediaDownloadDir()
	dispatched := 0

	for _, m := range queuedMedias {
		if m.DownloadURL == "" { continue }
		candidates := BuildMediaCandidateURLs(m.DownloadURL)
		if len(candidates) == 0 { continue }

		var urls []string
		for _, c := range candidates { urls = append(urls, c.URL) }

		cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
		ext := ".jpg"
		if m.Type == "video" || strings.Contains(m.DownloadURL, ".mp4") || strings.HasSuffix(strings.ToLower(m.MediaID), ".mp4") {
			ext = ".mp4"
		}
		fileName := fmt.Sprintf("%s_motrix%s", cleanID, ext)

		saveDir := filepath.Join(destRoot, "_motrix_temp")
		if owner, err := a.Repo.GetMediaOwnerUsername(m.MediaID); err == nil && owner != "" {
			saveDir = filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
		}

		gid, err := a.AddMotrixDownload(urls, saveDir, fileName)
		if err == nil && gid != "" {
			_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageMotrix)
			_ = a.Repo.UpdateMediaMetadata(m.MediaID, "OUTSOURCED", "", "", fmt.Sprintf("Motrix投入済 (GID: %s)", gid))
			dispatched++
		}
	}

	a.AppendPipelineLog("MOTRIX", "INFO", fmt.Sprintf("Motrixへ %d 件のQUEUEDメディアを一括ディスパッチ完了", dispatched))
	return dispatched, nil
}
