// app/app_rpc_thunder_reconciler.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"dozou_katanuki/models"
)

var parenRegex = regexp.MustCompile(`\(\d+\)(\.[a-zA-Z0-9]+)$`)

func cleanThunderFileName(fn string) string { return parenRegex.ReplaceAllString(fn, "$1") }

// ReconcileThunderTasksWithDB は 迅雷タスク一覧とDBを突合し、アクティブ実行中タスクのマップを返します
func (a *App) ReconcileThunderTasksWithDB() (int, map[string]bool) {
	activeMap := make(map[string]bool)
	if a.Repo != nil && a.Repo.DB() != nil {
		var activeTasks []models.ThunderTask
		_ = a.Repo.DB().Select("file_name").Where("status IN ?", []string{
			string(models.ThunderTaskOnboarded), string(models.ThunderTaskRunning), string(models.ThunderTaskHolding),
		}).Find(&activeTasks).Error
		for _, at := range activeTasks {
			activeMap[at.FileName] = true; activeMap[cleanThunderFileName(at.FileName)] = true
		}
	}
	status := a.GetThunderCDPStatus()
	if !status.IsConnected || len(status.CapturedTasks) == 0 { return len(activeMap), activeMap }

	updatedCount := 0
	for _, item := range status.CapturedTasks {
		dbFileName := cleanThunderFileName(item.FileName)
		if dbFileName == "" { continue }
		activeMap[item.FileName] = true; activeMap[dbFileName] = true

		currTask := a.getDozouThunderTask(dbFileName)
		if currTask == nil { continue }

		eval := EvaluateThunderTaskError(item.RawText)
		if (currTask.SummarySize != "" && currTask.SummarySize != "0B") || currTask.Status == models.ThunderTaskHolding {
			eval.Decision = DecisionHold
			if eval.SummarySize == "0B" || eval.SummarySize == "" { eval.SummarySize = currTask.SummarySize }
		}
		if currTask.DispatchedAt != nil && time.Since(*currTask.DispatchedAt) < 30*time.Second && eval.Decision == DecisionRetire {
			continue
		}

		switch eval.Decision {
		case DecisionRetire:
			if strings.Contains(item.RawText, "正在下载") || (strings.Contains(item.RawText, "/s") && !strings.Contains(item.RawText, "0B/s")) { continue }
			if currTask.Status != models.ThunderTaskRetired && a.Repo != nil {
				allRetired, mediaID, err := a.Repo.MarkThunderTaskRetiredAndCheckAll(dbFileName, eval.Reason)
				if err == nil && allRetired && mediaID != "" {
					_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "全候補タスクRETIREDによりRETAINED退避")
					a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("📦 全候補枯渇のため退避: %s", mediaID))
				}
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("🛑 ジョブ終了(RETIRED): %s (%s)", item.FileName, eval.Reason))
				updatedCount++
			}
			if wsURL, err := FetchThunderMainRendererWSUrl(9222); err == nil && wsURL != "" { a.deleteTaskByFileNameSilent(wsURL, item.FileName) }

		case DecisionHold:
			if currTask.Status != models.ThunderTaskHolding && a.Repo != nil {
				_ = a.Repo.MarkThunderTaskHolding(dbFileName, eval.SummarySize, eval.Reason)
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("⚡ サイズ確定により継続枠へ昇格(スロット解放): %s (%s)", dbFileName, eval.SummarySize))
				updatedCount++
			}

		case DecisionCooldown:
			if a.Repo != nil { _ = a.Repo.UpdateThunderTaskCooldown(dbFileName, eval.Reason) }

		default:
			isActive := strings.Contains(item.RawText, "正在") || strings.Contains(item.RawText, "连接") || (eval.HasSummary && eval.SummarySize != "0B")
			if currTask.Status == models.ThunderTaskRetired && isActive && a.Repo != nil {
				_ = a.Repo.MarkThunderTaskOnboarded(currTask.ID, currTask.SummarySize)
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("⚡ ユーザー再開を検知しタスク復帰: %s", dbFileName))
				updatedCount++
			}
			activeMap[item.FileName] = true; activeMap[dbFileName] = true
		}
	}
	return updatedCount, activeMap
}

func (a *App) getDozouThunderTask(fileName string) *models.ThunderTask {
	if a.Repo == nil || a.Repo.DB() == nil || fileName == "" { return nil }
	var t models.ThunderTask
	if err := a.Repo.DB().Where("file_name = ?", fileName).First(&t).Error; err != nil { return nil }
	return &t
}

func (a *App) CheckAndReonboardMissingTasks(activeMap map[string]bool, maxSlots int) int { return 0 }
