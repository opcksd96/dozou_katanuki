// app/app_rpc_pipeline_requests_worker.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"dozou_katanuki/models"
)

var (
	reqWorkerMu     sync.Mutex
	isReqWorkerBusy bool
)

type PipelineRunResult struct {
	QueuedTotal int  `json:"queued_total"`
	IsStarted   bool `json:"is_started"`
}

// ProcessQueuedViaRequests は QUEUED メディアに対して直接HTTPフェッチを試行し、失敗時は Motrix Next へ自動移管します
func (a *App) ProcessQueuedViaRequests() (*PipelineRunResult, error) {
	reqWorkerMu.Lock()
	if isReqWorkerBusy {
		reqWorkerMu.Unlock()
		return &PipelineRunResult{QueuedTotal: 0, IsStarted: false}, nil
	}
	isReqWorkerBusy = true
	reqWorkerMu.Unlock()
	unlock := func() { reqWorkerMu.Lock(); isReqWorkerBusy = false; reqWorkerMu.Unlock() }

	if a.Repo == nil || a.Repo.DB() == nil {
		unlock()
		return nil, fmt.Errorf("database not initialized")
	}
	var medias []models.Media
	if err := a.Repo.DB().Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Find(&medias).Error; err != nil {
		unlock()
		return nil, err
	}
	total := len(medias)
	if total == 0 {
		unlock()
		return &PipelineRunResult{QueuedTotal: 0, IsStarted: false}, nil
	}
	a.AppendPipelineLog("REQUESTS", "INFO", fmt.Sprintf("Requests開始: %d 件を順次処理します", total))
	go func() {
		defer unlock()
		a.processQueuedWorker(medias)
	}()
	return &PipelineRunResult{QueuedTotal: total, IsStarted: true}, nil
}

func (a *App) processQueuedWorker(medias []models.Media) {
	client, destRoot := &http.Client{Timeout: 8 * time.Second}, a.getMediaDownloadDir()
	ok, outsourced, escalated, total := 0, 0, 0, len(medias)
	for i, m := range medias {
		owner := "unknown"
		if o, _ := a.Repo.GetMediaOwnerUsername(m.MediaID); o != "" {
			owner = o
		}
		pos := fmt.Sprintf("[%d/%d]", i+1, total)
		a.AppendPipelineLog("REQUESTS", "INFO", fmt.Sprintf("%s 探索開始: %s (所有者: %s)", pos, m.MediaID, owner))
		_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageRequests)
		targetDir := filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
		_ = os.MkdirAll(targetDir, 0755)
		cands := BuildCandidateURLsFromMediaWithArticle(m.MediaID, m.DownloadURL, m.Type, m.ArticleID)
		fetched, urls, destPath := false, []string{}, filepath.Join(targetDir, m.MediaID)

		for cIdx, c := range cands {
			urls = append(urls, c.URL)
			// Wayback URL は直接HTTPで叩かない (IP BAN 防止) → Motrix/迅雷フォールバック専用
			if strings.Contains(c.URL, "web.archive.org") {
				continue
			}
			res := a.tryDirectFetchDetailed(client, c.URL, destPath)
			if res.Success {
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "COMPLETED", "", "", "Requests取得成功("+string(c.Type)+")")
				_ = a.Repo.MarkTaskCompleted(m.MediaID)
				a.AppendPipelineLog("REQUESTS", "SUCCESS", fmt.Sprintf("%s ✅ 取得成功 [%s]: %d bytes -> %s", pos, c.Type, res.Bytes, c.URL))
				ok, fetched = ok+1, true
				a.TriggerStashAllPipelines()
				break
			}
			a.AppendPipelineLog("REQUESTS", "DEBUG", fmt.Sprintf("%s   ↳ 試行 [%d/%d %s]: %s -> %s", pos, cIdx+1, len(cands), c.Type, c.URL, res.ErrorMsg))
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(200 * time.Millisecond)
		if !fetched {
			if gid, err := a.AddMotrixDownload(urls, targetDir, m.MediaID); err == nil && gid != "" {
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "OUTSOURCED", "", "", fmt.Sprintf("Motrix Next移管 (GID: %s)", gid))
				_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageMotrix)
				a.AppendPipelineLog("MOTRIX", "INFO", fmt.Sprintf("%s 🚀 Requests全滅 ➔ Motrix Next投入 (%d候補, GID: %s)", pos, len(urls), gid))
				outsourced++
			} else {
				_ = a.Repo.UpdateMediaMetadata(m.MediaID, "ESCALATED", "", "", "Motrixオフライン→迅雷エスカレーション待機")
				_ = a.Repo.UpdateMediaCheckpointTime(m.MediaID, models.StageThunder)
				a.AppendPipelineLog("REQUESTS", "WARN", fmt.Sprintf("%s ⚡ Motrix利用不可 ➔ 迅雷エスカレーション待機: %s", pos, m.MediaID))
				escalated++
			}
		}
	}
	a.AppendPipelineLog("REQUESTS", "SUCCESS", fmt.Sprintf("Requests完了: 成功 %d / Motrix移管 %d / 迅雷直接 %d (総計 %d)", ok, outsourced, escalated, total))
	if outsourced > 0 || escalated > 0 {
		go func() {
			time.Sleep(2 * time.Second)
			_, _ = a.SyncCompletedDownloads()
			if escalated > 0 && !a.isThunderOrchestratorRunning() {
				_, _ = a.StartThunderOrchestrator(3, 4)
			}
		}()
	}
}
