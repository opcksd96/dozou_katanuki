// app/app_rpc_thunder_orchestrator_api.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"dozou_katanuki/models"
)

// StartThunderOrchestrator は RETAINED メディアから303個の厳選タスクを抽出し、間欠投入を開始します
func (a *App) StartThunderOrchestrator(maxSlots, intervalSec int) (models.ThunderOrchestratorStatus, error) {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()

	if orchState.isRunning { return a.getOrchestratorStatusLocked(), nil }
	if maxSlots <= 0 || maxSlots > 3 { maxSlots = 3 }
	if intervalSec <= 0 { intervalSec = 4 }
	orchState.config.MaxConcurrentSlots = maxSlots
	orchState.config.IntervalSeconds = intervalSec

	orchState.slots = make([]models.ThunderOrchestratorSlot, maxSlots)
	for i := 0; i < maxSlots; i++ {
		orchState.slots[i] = models.ThunderOrchestratorSlot{Index: i, IsOccupied: false}
	}

	orchState.queue = a.buildThunderOrchestratorTasks()
	orchState.isRunning, orchState.isPaused = true, false
	orchState.stopCh, orchState.pauseCh, orchState.resumeCh = make(chan struct{}), make(chan struct{}), make(chan struct{})
	orchState.processedMap = make(map[string]bool)

	go a.StartThunderCDPAdaptivePoller()
	go a.runThunderOrchestrationWorker()
	return a.getOrchestratorStatusLocked(), nil
}

// isThunderOrchestratorRunning はオーケストレータが稼働中かを返します
func (a *App) isThunderOrchestratorRunning() bool {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()
	return orchState.isRunning && !orchState.isPaused
}

// ResetAndRebuildThunderQueue はステータスを RETAINED へ差し戻してキューを再構築します
func (a *App) ResetAndRebuildThunderQueue(resetVideos bool) (models.ThunderOrchestratorStatus, error) {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()
	if orchState.isRunning { orchState.isRunning, orchState.isPaused = false, false; close(orchState.stopCh) }
	if a.Repo != nil {
		if resetVideos { _, _ = a.Repo.ResetVideosToRetained() } else { _, _ = a.Repo.ResetAllFailedToRetained() }
	}
	maxSlots := orchState.config.MaxConcurrentSlots
	if maxSlots <= 0 { maxSlots = 12 }
	orchState.slots = make([]models.ThunderOrchestratorSlot, maxSlots)
	for i := 0; i < maxSlots; i++ { orchState.slots[i] = models.ThunderOrchestratorSlot{Index: i, IsOccupied: false} }
	orchState.queue = a.buildThunderOrchestratorTasks()
	orchState.recentTasks = nil
	return a.getOrchestratorStatusLocked(), nil
}

// PauseThunderOrchestrator は一時停止します
func (a *App) PauseThunderOrchestrator() bool {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()
	if !orchState.isRunning || orchState.isPaused { return false }
	orchState.isPaused = true
	return true
}

// ResumeThunderOrchestrator は再開します
func (a *App) ResumeThunderOrchestrator() bool {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()
	if !orchState.isRunning || !orchState.isPaused { return false }
	orchState.isPaused = false
	return true
}

// StopThunderOrchestrator はオーケストレーターを完全停止します
func (a *App) StopThunderOrchestrator() bool {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()
	if !orchState.isRunning { return false }
	orchState.isRunning, orchState.isPaused = false, false
	close(orchState.stopCh)
	return true
}
