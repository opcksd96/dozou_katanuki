// app/app_rpc_thunder_orchestrator_api.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"dozou_katanuki/models"
)

// StartThunderOrchestrator は RETAINED メディアから303個の厳選タスクを抽出し、間欠投入を開始します
func (a *App) StartThunderOrchestrator(maxSlots, intervalSec int) (models.ThunderOrchestratorStatus, error) {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()

	if orchState.isRunning {
		return a.getOrchestratorStatusLocked(), nil
	}

	if maxSlots <= 0 { maxSlots = 12 }
	if intervalSec <= 0 { intervalSec = 5 }
	orchState.config.MaxConcurrentSlots = maxSlots
	orchState.config.IntervalSeconds = intervalSec

	orchState.slots = make([]models.ThunderOrchestratorSlot, maxSlots)
	for i := 0; i < maxSlots; i++ {
		orchState.slots[i] = models.ThunderOrchestratorSlot{Index: i, IsOccupied: false}
	}

	orchState.queue = a.buildThunderOrchestratorTasks()
	orchState.isRunning = true
	orchState.isPaused = false
	orchState.stopCh = make(chan struct{})
	orchState.pauseCh = make(chan struct{})
	orchState.resumeCh = make(chan struct{})
	orchState.processedMap = make(map[string]bool)

	go a.runThunderOrchestrationWorker()
	return a.getOrchestratorStatusLocked(), nil
}

// ResetAndRebuildThunderQueue は動画等のステータスを RETAINED へ差し戻してキューを再構築します
func (a *App) ResetAndRebuildThunderQueue(resetVideos bool) (models.ThunderOrchestratorStatus, error) {
	orchState.mu.Lock()
	defer orchState.mu.Unlock()

	if orchState.isRunning {
		orchState.isRunning = false
		orchState.isPaused = false
		close(orchState.stopCh)
	}

	if a.Repo != nil {
		if resetVideos {
			_, _ = a.Repo.ResetVideosToRetained()
		} else {
			_, _ = a.Repo.ResetAllFailedToRetained()
		}
	}

	maxSlots := orchState.config.MaxConcurrentSlots
	if maxSlots <= 0 { maxSlots = 12 }
	orchState.slots = make([]models.ThunderOrchestratorSlot, maxSlots)
	for i := 0; i < maxSlots; i++ {
		orchState.slots[i] = models.ThunderOrchestratorSlot{Index: i, IsOccupied: false}
	}

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
	orchState.isRunning = false
	orchState.isPaused = false
	close(orchState.stopCh)
	return true
}
