// app/app_rpc_pipeline_continuous_engine.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"sync"
	"time"
)

type PipelineEngineState struct {
	mu        sync.Mutex
	isRunning bool
	stopCh    chan struct{}
}

var pipeEngineState = PipelineEngineState{}

// TogglePipelineAutoEngine はパイプラインの完全自動運転ループを開始/停止します
func (a *App) TogglePipelineAutoEngine(enable bool) (bool, error) {
	pipeEngineState.mu.Lock()
	defer pipeEngineState.mu.Unlock()

	if enable && !pipeEngineState.isRunning {
		pipeEngineState.isRunning = true
		pipeEngineState.stopCh = make(chan struct{})
		a.AppendPipelineLog("SYSTEM", "INFO", "🚀 パイプライン完全自動運転エンジンを起動しました")
		go a.runContinuousPipelineLoop(pipeEngineState.stopCh)
	} else if !enable && pipeEngineState.isRunning {
		pipeEngineState.isRunning = false
		close(pipeEngineState.stopCh)
		a.AppendPipelineLog("SYSTEM", "WARN", "⏸️ パイプライン完全自動運転エンジンを停止しました")
	}
	return pipeEngineState.isRunning, nil
}

// IsPipelineAutoEngineRunning は自動運転中かどうかを返します
func (a *App) IsPipelineAutoEngineRunning() bool {
	pipeEngineState.mu.Lock()
	defer pipeEngineState.mu.Unlock()
	return pipeEngineState.isRunning
}

func (a *App) runContinuousPipelineLoop(stopCh chan struct{}) {
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			a.executeAutonomousPipelineCycle()
		}
	}
}

func (a *App) executeAutonomousPipelineCycle() {
	if a.Repo == nil || a.Repo.DB() == nil { return }
	db := a.Repo.DB()

	// 1. QUEUED があれば Requests ➔ Motrix 投入を自動キック
	var qCount int64
	_ = db.Model(&models_Media{}).Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Count(&qCount).Error
	if qCount > 0 {
		_, _ = a.ProcessQueuedViaRequests()
	}

	// 2. Motrix (OUTSOURCED) の完了同期を自動実行
	_, _ = a.SyncCompletedDownloads()

	// 3. 迅雷ダウンロード完了の回収 ＆ Stash自動連携
	_, _ = a.SyncThunderDownloads("")
}

type models_Media struct {
	DownloadStatus string `gorm:"column:download_status"`
	IsTrash        int    `gorm:"column:is_trash"`
}

func (models_Media) TableName() string { return "media" }
