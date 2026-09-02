// app/app_rpc_thunder_orchestrator_worker.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"dozou_katanuki/models"
)

func (a *App) runThunderOrchestrationWorker() {
	interval := time.Duration(orchState.config.IntervalSeconds) * time.Second
	if interval <= 0 { interval = 4 * time.Second }
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

			existingFileMap := a.fetchExistingThunderFilesMap()
			destDir := a.getMediaDownloadDir()

			runningCount := 0
			for _, t := range orchState.queue {
				if t.Status == "running" {
					if fi, err := os.Stat(filepath.Join(destDir, t.FileName)); err == nil && fi.Size() > 0 {
						t.Status = "completed"
					} else if len(existingFileMap) > 0 && !existingFileMap[t.FileName] {
						t.Status = "depleted"
					} else {
						runningCount++
					}
				}
			}

			maxSlots := orchState.config.MaxConcurrentSlots
			if maxSlots <= 0 || maxSlots > 3 { maxSlots = 3 }
			if runningCount >= maxSlots || len(existingFileMap) >= maxSlots {
				orchState.mu.Unlock()
				continue
			}

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
				if runningCount == 0 { orchState.isRunning = false }
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
	for _, t := range a.GetThunderCDPStatus().CapturedTasks {
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
	if len(orchState.recentTasks) > 30 { orchState.recentTasks = orchState.recentTasks[:30] }
}
