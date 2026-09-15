// app/app_rpc_thunder_orchestrator_state.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
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
		MaxConcurrentSlots: 3,
		IntervalSeconds:    5,
		TopResolutionsOnly: true,
	},
	slots: make([]models.ThunderOrchestratorSlot, 3),
}

func init() {
	for i := 0; i < 3; i++ {
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
	total, pending, running, holding, success, failed := 0, 0, 0, 0, 0, 0
	if a.Repo != nil {
		if c, err := a.Repo.GetThunderTaskCounts(); err == nil {
			total, pending, running, holding, success, failed = c.Total, c.Pending, c.Running, c.Holding, c.Completed, c.Failed
		}
	}

	cdpTasks, _ := GetTasksViaDirectCDP(9222)
	recent := orchState.recentTasks
	// ⚡ 投入中リストのフォールバック: オーケストレーター起動前でも迅雷内のリアルタイムタスクを表示
	if len(recent) == 0 && len(cdpTasks) > 0 {
		for i, ct := range cdpTasks {
			st := "running"
			if ct.TaskStatusCode == 11 { st = "completed" } else if ct.TaskStatusCode == 9 { st = "depleted" }
			recent = append(recent, models.ThunderOrchestratorTask{
				ID: fmt.Sprintf("%d", ct.TaskID), MediaID: ct.FileName, FileName: ct.FileName,
				URL: ct.URL, Status: st, SlotIndex: i % 3, ErrorMessage: ct.DetailText,
			})
		}
	}

	occupied := 0
	for _, s := range orchState.slots { if s.IsOccupied { occupied++ } }
	maxSlots := orchState.config.MaxConcurrentSlots
	if maxSlots <= 0 { maxSlots = 3 }

	return models.ThunderOrchestratorStatus{
		IsRunning:       orchState.isRunning,
		IsPaused:        orchState.isPaused,
		Config:          orchState.config,
		TotalJobs:       total,
		PendingJobs:     pending,
		RunningJobs:     running,
		HoldingJobs:     holding,
		SuccessJobs:     success,
		FailedJobs:      failed,
		OccupiedSlots:   occupied,
		TotalSlots:      maxSlots,
		TotalMediaCount: total,
		Slots:           orchState.slots,
		RecentTasks:     recent,
		CDPTasks:        cdpTasks,
	}
}
