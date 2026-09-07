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
	ticker := time.NewTicker(interval); defer ticker.Stop()

	for {
		select {
		case <-orchState.stopCh: return
		case <-ticker.C:
			orchState.mu.Lock()
			if !orchState.isRunning || orchState.isPaused { orchState.mu.Unlock(); continue }
			maxSlots := orchState.config.MaxConcurrentSlots
			if maxSlots <= 0 || maxSlots > 3 { maxSlots = 3 }
			if !a.GetThunderCDPStatus().IsConnected {
				a.AppendPipelineLog("THUNDER", "ERROR", "❌ 迅雷CDP (port 9222) 未接続のため待機します。")
				orchState.mu.Unlock(); continue
			}

			_, activeMap := a.ReconcileThunderTasksWithDB()
			a.CheckAndReonboardMissingTasks(activeMap, maxSlots)

			destDir, runningCount := a.getMediaDownloadDir(), 0
			activeMedia := make(map[string]bool)
			for _, t := range orchState.queue {
				if t.Status == "running" || t.Status == "holding" {
					if fi, err := os.Stat(filepath.Join(destDir, t.FileName)); err == nil && fi.Size() > 0 {
						t.Status = "completed"; continue
					}
				}
				if t.Status == "running" {
					isRecent := t.DispatchedAt != nil && time.Since(*t.DispatchedAt) < 15*time.Second
					if !activeMap[t.FileName] && !isRecent {
						t.Status = "holding"
					} else {
						runningCount++
						activeMedia[t.MediaID] = true
					}
				}
			}

			// スロットの空き枠分だけ、未着手の別メディアタスクを順次ディスパッチ
			for runningCount < maxSlots {
				var nextTask *models.ThunderOrchestratorTask
				for _, t := range orchState.queue {
					if t.Status == "pending" && !orchState.processedMap[t.ID] && !activeMap[t.FileName] && !activeMedia[t.MediaID] {
						nextTask = t; break
					}
				}
				if nextTask == nil {
					for _, mt := range a.buildThunderOrchestratorTasks() {
						if !orchState.processedMap[mt.ID] { orchState.queue = append(orchState.queue, mt) }
					}
					for _, t := range orchState.queue {
						if t.Status == "pending" && !orchState.processedMap[t.ID] && !activeMap[t.FileName] && !activeMedia[t.MediaID] {
							nextTask = t; break
						}
					}
				}
				if nextTask == nil { break }
				if fi, err := os.Stat(filepath.Join(destDir, nextTask.FileName)); err == nil && fi.Size() > 0 {
					nextTask.Status = "completed"
					if a.Repo != nil { _ = a.Repo.UpdateMediaMetadata(nextTask.MediaID, "COMPLETED", "", "", "実ファイル確認済み") }
					continue
				}
				a.dispatchTaskDirectly(nextTask)
				activeMedia[nextTask.MediaID] = true
				activeMap[nextTask.FileName] = true
				runningCount++
			}
			orchState.mu.Unlock()
		}
	}
}

func (a *App) dispatchTaskDirectly(task *models.ThunderOrchestratorTask) {
	now := time.Now()
	task.Status, task.DispatchedAt = "running", &now
	orchState.processedMap[task.ID] = true
	if a.Repo != nil && task.MediaID != "" {
		_ = a.Repo.UpdateMediaMetadata(task.MediaID, "ESCALATED", "", "", fmt.Sprintf("迅雷投入中 (%s)", task.ResolutionType))
		_ = a.Repo.MarkThunderTaskOnboarded(task.ID, "")
	}
	destDir := a.getMediaDownloadDir()
	go func(t *models.ThunderOrchestratorTask, dest string) { _ = AddTaskViaThunderCOM(t.URL, t.FileName, dest) }(task, destDir)
	orchState.recentTasks = append([]models.ThunderOrchestratorTask{*task}, orchState.recentTasks...)
	if len(orchState.recentTasks) > 30 { orchState.recentTasks = orchState.recentTasks[:30] }
}
