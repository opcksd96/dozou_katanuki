// app/app_rpc_thunder_orchestrator_worker.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"time"

	"dozou_katanuki/models"
)

func (a *App) runThunderOrchestrationWorker() {
	interval := time.Duration(orchState.config.IntervalSeconds) * time.Second
	if interval <= 0 { interval = 4 * time.Second }
	ticker := time.NewTicker(interval); defer ticker.Stop()
	cycle := 0

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

			cycle++
			if cycle%3 == 0 { _, _ = a.ReapFailedEmptyThunderTasks() }
			if cycle%15 == 0 { _, _ = a.ReactivateThunderTasks(true) }

			_, activeMap := a.ReconcileThunderTasksWithDB()
			runningCount := 0
			activeMedia := make(map[string]bool)
			for _, t := range orchState.queue {
				if t.Status == "running" || t.Status == "holding" {
					if a.CheckMediaFileExists(t.MediaID, t.FileName) {
						t.Status = "completed"; continue
					}
				}
				if t.Status == "running" {
					isRecent := t.DispatchedAt != nil && time.Since(*t.DispatchedAt) < 15*time.Second
					if !activeMap[t.FileName] && !isRecent { t.Status = "holding" } else { runningCount++; activeMedia[t.MediaID] = true }
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
				if a.CheckMediaFileExists(nextTask.MediaID, nextTask.FileName) {
					nextTask.Status = "completed"
					if a.Repo != nil { _ = a.Repo.UpdateMediaMetadata(nextTask.MediaID, "COMPLETED", "", "", "実ファイル確認済み") }
					continue
				}
				a.dispatchTaskDirectly(nextTask)
				activeMedia[nextTask.MediaID], activeMap[nextTask.FileName] = true, true
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
	destDir := a.ResolveMediaTargetDir(task.MediaID)
	items := []models.ThunderTaskInputItem{{URL: task.URL, FileName: task.FileName}}
	go func(its []models.ThunderTaskInputItem, dest string) { _ = AddTasksViaDirectCDP(9222, its, dest) }(items, destDir)
	orchState.recentTasks = append([]models.ThunderOrchestratorTask{*task}, orchState.recentTasks...)
	if len(orchState.recentTasks) > 30 { orchState.recentTasks = orchState.recentTasks[:30] }
}
