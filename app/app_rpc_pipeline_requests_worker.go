// app/app_rpc_pipeline_requests_worker.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dozou_katanuki/models"
)

type PipelineRunResult struct {
	QueuedTotal int  `json:"queued_total"`
	IsStarted   bool `json:"is_started"`
}

// ProcessQueuedViaRequests は QUEUED メディアに対して直接HTTPフェッチを試行し、失敗時は Motrix Next へ自動移管します
func (a *App) ProcessQueuedViaRequests() (*PipelineRunResult, error) {
	if a.Repo == nil || a.Repo.DB() == nil { return nil, fmt.Errorf("database not initialized") }
	var medias []models.Media
	if err := a.Repo.DB().Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Find(&medias).Error; err != nil {
		return nil, err
	}
	total := len(medias)
	if total == 0 { return &PipelineRunResult{QueuedTotal: 0, IsStarted: false}, nil }

	a.AppendPipelineLog("REQUESTS", "INFO", fmt.Sprintf("Requests開始: %d 件を処理します", total))
	go a.processQueuedWorker(medias)
	return &PipelineRunResult{QueuedTotal: total, IsStarted: true}, nil
}

func (a *App) processQueuedWorker(medias []models.Media) {
	client := &http.Client{Timeout: 8 * time.Second}
	destRoot := a.getMediaDownloadDir()
	ok, outsourced, escalated, total := 0, 0, 0, len(medias)

	for i, m := range medias {
		u := m.DownloadURL
		if u == "" { u = fmt.Sprintf("https://pbs.twimg.com/media/%s", m.MediaID) }
		pos := fmt.Sprintf("[%d/%d]", i+1, total)
		a.AppendPipelineLog("REQUESTS", "INFO", fmt.Sprintf("%s 処理中: %s", pos, m.MediaID))
		_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageRequests)
		owner, _ := a.Repo.GetMediaOwnerUsername(m.MediaID)
		if owner == "" { owner = "unknown" }
		cleanID := strings.TrimSuffix(m.MediaID, filepath.Ext(m.MediaID))
		ext := ".jpg"
		if m.Type == "video" || strings.Contains(u, ".mp4") { ext = ".mp4" }
		targetDir := filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
		_ = os.MkdirAll(targetDir, 0755)

		candidates, fetched, urls := BuildMediaCandidateURLs(u), false, []string{}
		for _, c := range candidates {
			urls = append(urls, c.URL)
			destPath := filepath.Join(targetDir, cleanID+"_"+string(c.Type)+ext)
			if a.tryDirectFetch(client, c.URL, destPath) {
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "COMPLETED", "", "", "Requests取得成功("+string(c.Type)+")")
				a.AppendPipelineLog("REQUESTS", "SUCCESS", fmt.Sprintf("%s ✅ 取得成功: %s (%s)", pos, m.MediaID, string(c.Type)))
				ok, fetched = ok+1, true
				break
			}
		}
		if !fetched {
			gid, err := a.AddMotrixDownload(urls, targetDir, cleanID+"_motrix"+ext)
			if err == nil && gid != "" {
				_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageMotrix)
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "OUTSOURCED", "", "", fmt.Sprintf("Motrix Next移管 (GID: %s)", gid))
				a.AppendPipelineLog("MOTRIX", "INFO", fmt.Sprintf("%s 🚀 Motrix Nextへ移管: %s (GID: %s)", pos, m.MediaID, gid))
				outsourced++
			} else {
				_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageThunder)
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "ESCALATED", "", "", "Requests枯渇/Motrix未稼働→迅雷エスカレーション")
				a.AppendPipelineLog("REQUESTS", "WARN", fmt.Sprintf("%s ⚡ 迅雷へ (Motrix:%v): %s", pos, err, m.MediaID))
				escalated++
			}
		}
	}
	a.AppendPipelineLog("REQUESTS", "SUCCESS", fmt.Sprintf("Requests完了: 成功 %d 件 / Motrix移管 %d 件 / 迅雷 %d 件 (計 %d)", ok, outsourced, escalated, total))
}

func (a *App) tryDirectFetch(client *http.Client, url, destPath string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return false }
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil { return false }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { return false }
	out, err := os.Create(destPath)
	if err != nil { return false }
	defer out.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil || n < 1024 { _ = os.Remove(destPath); return false }
	return true
}
