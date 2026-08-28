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
	ticker := time.NewTicker(time.Duration(orchState.config.IntervalSeconds) * time.Second)
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
			if freeIdx := a.findOrFreeSlot(); freeIdx != -1 {
				a.dispatchNextTaskToSlot(freeIdx)
			}
			orchState.mu.Unlock()
		}
	}
}

func (a *App) findOrFreeSlot() int {
	for i := 0; i < len(orchState.slots); i++ {
		if !orchState.slots[i].IsOccupied { return i }
	}
	tempDir := `D:\迅雷下载`
	if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.ThunderDownloadDir != "" {
		tempDir = cfg.Storage.ThunderDownloadDir
	}
	now := time.Now()
	for i := 0; i < len(orchState.slots); i++ {
		s := &orchState.slots[i]
		if !s.IsOccupied || s.DispatchedAt == nil { continue }
		if now.Sub(*s.DispatchedAt) >= 30*time.Second {
			if s.CurrentTask != nil && a.Repo != nil {
				target := filepath.Join(tempDir, s.CurrentTask.FileName)
				hasXltd := isFile(target+".xltd") || isFile(target+".td")
				if isFile(target) {
					s.CurrentTask.Status = "success"
					go func() { _, _ = a.SyncThunderDownloads(tempDir) }()
				} else if hasXltd {
					s.CurrentTask.Status = "running"
					_ = a.Repo.UpdateMediaMetadata(s.CurrentTask.MediaID, "ESCALATED", "", "", "迅雷ダウンロード中 (*.xltd)")
				} else {
					s.CurrentTask.Status = "retained"
					_ = a.Repo.UpdateMediaMetadata(s.CurrentTask.MediaID, "RETAINED", "", "", "即時キャッシュなし・長期待機")
				}
			}
			s.IsOccupied, s.CurrentTask = false, nil
			return i
		}
	}
	return -1
}

func isFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}

func (a *App) dispatchNextTaskToSlot(slotIdx int) {
	var nextTask *models.ThunderOrchestratorTask
	for _, t := range orchState.queue {
		if t.Status == "pending" {
			nextTask = t
			break
		}
	}
	if nextTask == nil { return }
	now := time.Now()
	nextTask.Status, nextTask.SlotIndex, nextTask.DispatchedAt = "running", slotIdx, &now
	orchState.slots[slotIdx].IsOccupied, orchState.slots[slotIdx].CurrentTask, orchState.slots[slotIdx].DispatchedAt = true, nextTask, &now

	if a.Repo != nil && nextTask.MediaID != "" {
		_ = a.Repo.UpdateMediaMetadata(nextTask.MediaID, "ESCALATED", "", "", fmt.Sprintf("迅雷投入中 (%s)", nextTask.ResolutionType))
	}
	destDir := a.getMediaDownloadDir()
	go func(t *models.ThunderOrchestratorTask, dest string) {
		_ = AddTaskViaThunderCOM(t.URL, t.FileName, dest)
	}(nextTask, destDir)
	orchState.recentTasks = append([]models.ThunderOrchestratorTask{*nextTask}, orchState.recentTasks...)
	if len(orchState.recentTasks) > 20 { orchState.recentTasks = orchState.recentTasks[:20] }
}
