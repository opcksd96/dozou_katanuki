// app/app_rpc_thunder_orchestrator_worker.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

func (a *App) runThunderOrchestrationWorker() {
	interval := time.Duration(orchState.config.IntervalSeconds) * time.Second
	if interval <= 0 { interval = 3 * time.Second }
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-orchState.stopCh:
			return
		case <-ticker.C:
			orchState.mu.Lock()
			if !orchState.isRunning || orchState.isPaused {
				orchState.mu.Unlock()
				continue
			}

			// 迅雷の既存タスク一覧（ファイル名）を取得して重複投入を防止
			existingFileMap := a.fetchExistingThunderFilesMap()

			var nextTask *models.ThunderOrchestratorTask
			for _, t := range orchState.queue {
				if t.Status == "pending" {
					if existingFileMap[t.FileName] || orchState.processedMap[t.ID] {
						t.Status = "running"
						continue
					}
					nextTask = t
					break
				}
			}

			if nextTask == nil {
				orchState.isRunning = false
				orchState.mu.Unlock()
				continue
			}

			a.dispatchTaskDirectly(nextTask)
			orchState.mu.Unlock()
		}
	}
}

func (a *App) fetchExistingThunderFilesMap() map[string]bool {
	m := make(map[string]bool)
	st := a.GetThunderCDPStatus()
	for _, t := range st.CapturedTasks {
		if t.FileName != "" { m[t.FileName] = true }
	}
	return m
}

func (a *App) dispatchTaskDirectly(task *models.ThunderOrchestratorTask) {
	now := time.Now()
	task.Status = "running"
	task.DispatchedAt = &now
	orchState.processedMap[task.ID] = true

	if a.Repo != nil && task.MediaID != "" {
		_ = a.Repo.UpdateMediaMetadata(task.MediaID, "ESCALATED", "", "", fmt.Sprintf("迅雷投入中 (%s)", task.ResolutionType))
		_ = a.Repo.MarkThunderTaskRunning(task.ID)
	}

	destDir := a.getMediaDownloadDir()
	go func(t *models.ThunderOrchestratorTask, dest string) {
		_ = AddTaskViaThunderCOM(t.URL, t.FileName, dest)
	}(task, destDir)

	orchState.recentTasks = append([]models.ThunderOrchestratorTask{*task}, orchState.recentTasks...)
	if len(orchState.recentTasks) > 30 {
		orchState.recentTasks = orchState.recentTasks[:30]
	}
}
