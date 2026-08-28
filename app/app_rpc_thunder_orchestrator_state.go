// app/app_rpc_thunder_orchestrator_state.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"sync"

	"dozou_katanuki/models"
)

type thunderOrchestratorState struct {
	mu           sync.RWMutex
	isRunning    bool
	isPaused     bool
	config       models.ThunderOrchestratorConfig
	queue        []*models.ThunderOrchestratorTask
	slots        []models.ThunderOrchestratorSlot
	recentTasks  []models.ThunderOrchestratorTask
	stopCh       chan struct{}
	pauseCh      chan struct{}
	resumeCh     chan struct{}
	processedMap map[string]bool
}

var orchState = &thunderOrchestratorState{
	config: models.ThunderOrchestratorConfig{
		MaxConcurrentSlots: 12,
		IntervalSeconds:    5,
		TopResolutionsOnly: true,
	},
	slots: make([]models.ThunderOrchestratorSlot, 12),
}

func init() {
	for i := 0; i < 12; i++ {
		orchState.slots[i] = models.ThunderOrchestratorSlot{Index: i, IsOccupied: false}
	}
}

// GetThunderOrchestratorStatus は現在のリアルタイム稼働状況を返します
func (a *App) GetThunderOrchestratorStatus() models.ThunderOrchestratorStatus {
	orchState.mu.RLock()
	defer orchState.mu.RUnlock()
	return a.getOrchestratorStatusLocked()
}

func (a *App) getOrchestratorStatusLocked() models.ThunderOrchestratorStatus {
	total := len(orchState.queue)
	pending, running, success, failed := 0, 0, 0, 0
	for _, t := range orchState.queue {
		switch t.Status {
		case "pending":
			pending++
		case "running":
			running++
		case "success":
			success++
		case "failed":
			failed++
		}
	}

	return models.ThunderOrchestratorStatus{
		IsRunning:       orchState.isRunning,
		IsPaused:        orchState.isPaused,
		Config:          orchState.config,
		TotalJobs:       total,
		PendingJobs:     pending,
		RunningJobs:     running,
		SuccessJobs:     success,
		FailedJobs:      failed,
		TotalMediaCount: total / 3,
		Slots:           orchState.slots,
		RecentTasks:     orchState.recentTasks,
	}
}
