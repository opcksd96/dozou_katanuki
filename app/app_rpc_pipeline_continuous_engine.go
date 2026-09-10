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
		a.ResumeThunderOrchestrator()
		go a.runContinuousPipelineLoop(pipeEngineState.stopCh)
	} else if !enable && pipeEngineState.isRunning {
		pipeEngineState.isRunning = false
		close(pipeEngineState.stopCh)
		a.AppendPipelineLog("SYSTEM", "WARN", "⏸️ パイプライン完全自動運転エンジンを停止しました")
		a.PauseThunderOrchestrator()
	}
	return pipeEngineState.isRunning, nil
}

// IsPipelineAutoEngineRunning は自動運転中かどうかを返します
func (a *App) IsPipelineAutoEngineRunning() bool {
	pipeEngineState.mu.Lock()
	defer pipeEngineState.mu.Unlock()
	return pipeEngineState.isRunning
}

type PipelineCycleResult struct {
	QueuedCount    int64 `json:"queued_count"`
	EscalatedCount int64 `json:"escalated_count"`
	Success        bool  `json:"success"`
}

func (a *App) runContinuousPipelineLoop(stopCh chan struct{}) {
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			_, _ = a.ExecutePipelineCycleNow()
		}
	}
}

// ExecutePipelineCycleNow は統括契約に基づき、4つの副次契約（Requests ➔ Motrix ➔ 迅雷 ➔ Stash）を1サイクル完走します
func (a *App) ExecutePipelineCycleNow() (*PipelineCycleResult, error) {
	if a.Repo == nil || a.Repo.DB() == nil {
		return &PipelineCycleResult{Success: false}, nil
	}
	db := a.Repo.DB()

	var qCount, eCount int64
	_ = db.Model(&models_Media{}).Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Count(&qCount).Error
	_ = db.Model(&models_Media{}).Where("download_status = 'ESCALATED' AND (is_trash = 0 OR is_trash IS NULL)").Count(&eCount).Error

	// 1. 突貫突撃契約 (Requests)
	if qCount > 0 {
		_, _ = a.ProcessQueuedViaRequests()
	}

	// 2. バッチ契約 (Motrix Next)
	_, _ = a.SyncCompletedDownloads()
	_, _ = a.ReconcileMotrixTasks()

	// 3. P2SP捜索契約 (迅雷 COM+CDP)
	if eCount > 0 && !a.isThunderOrchestratorRunning() {
		_, _ = a.StartThunderOrchestrator(3, 4)
	}
	_, _ = a.SyncThunderDownloads("")

	// 4. 資産取り込み契約 (Stash)
	go a.ScanUnsyncedMediaAndTriggerStash()

	return &PipelineCycleResult{QueuedCount: qCount, EscalatedCount: eCount, Success: true}, nil
}

type models_Media struct {
	DownloadStatus string `gorm:"column:download_status"`
	IsTrash        int    `gorm:"column:is_trash"`
}
func (models_Media) TableName() string { return "media" }
