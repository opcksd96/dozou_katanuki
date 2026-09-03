// app/app_rpc_motrix_reconcile.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// ReconcileMotrixTasks はDB上で OUTSOURCED (Motrix担当) になっているが Motrix キューに存在しないタスクを検知し、再投入またはエスカレーションします
func (a *App) ReconcileMotrixTasks() (int, error) {
	if a.Repo == nil || a.Repo.DB() == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	fullQueue, err := a.FetchMotrixFullQueue()
	if err != nil {
		return 0, err
	}

	activeMap := make(map[string]bool)
	for _, t := range fullQueue.ActiveTasks {
		if t.FileName != "" { activeMap[t.FileName] = true }
	}
	for _, t := range fullQueue.WaitingTasks {
		if t.FileName != "" { activeMap[t.FileName] = true }
	}

	var outsourcedMedias []models.Media
	if err := a.Repo.DB().Where("download_status = 'OUTSOURCED' AND (is_trash = 0 OR is_trash IS NULL)").Find(&outsourcedMedias).Error; err != nil {
		return 0, err
	}

	reconciled := 0
	for _, m := range outsourcedMedias {
		if activeMap[m.MediaID] {
			continue
		}
		a.AppendPipelineLog("MOTRIX", "WARN", fmt.Sprintf("Motrixキューから脱落したタスクを検出: %s", m.MediaID))
		
		owner := "unknown"
		if o, _ := a.Repo.GetMediaOwnerUsername(m.MediaID); o != "" { owner = o }
		targetDir := filepath.Join(a.getMediaDownloadDir(), owner, "X(Twitter)", "_assets")
		cands := BuildCandidateURLsFromMediaWithArticle(m.MediaID, m.DownloadURL, m.Type, m.ArticleID)
		
		var urls []string
		for _, c := range cands {
			if !strings.Contains(c.URL, "web.archive.org") {
				urls = append(urls, c.URL)
			}
		}

		if len(urls) > 0 {
			if gid, err := a.AddMotrixDownload(urls, targetDir, m.MediaID); err == nil && gid != "" {
				a.AppendPipelineLog("MOTRIX", "INFO", fmt.Sprintf("Motrixへ再投入成功: %s (GID: %s)", m.MediaID, gid))
			} else {
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "ESCALATED", "", "", "Motrix再投入失敗→迅雷エスカレーション待機")
				_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageThunder)
			}
		} else {
			_ = a.Repo.UpdateMediaMetadata(m.MediaID, "ESCALATED", "", "", "有効なURLなし→迅雷エスカレーション待機")
			_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageThunder)
		}
		reconciled++
	}
	return reconciled, nil
}
