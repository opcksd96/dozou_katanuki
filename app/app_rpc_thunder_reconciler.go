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

// ReconcileThunderTasksWithDB は 迅雷タスク一覧とDBを突合し、実態に基づきDBを同期します
func (a *App) ReconcileThunderTasksWithDB() (int, map[string]bool) {
	status := a.GetThunderCDPStatus()
	if !status.IsConnected { return 0, a.getDBActiveMap() }

	liveMap := make(map[string]bool)
	updatedCount := 0

	// 1. 迅雷実機タスクの評価と同期
	for _, item := range status.CapturedTasks {
		dbFileName := cleanThunderFileName(item.FileName)
		if dbFileName == "" { continue }
		liveMap[item.FileName] = true; liveMap[dbFileName] = true

		currTask := a.getDozouThunderTask(dbFileName)
		if currTask == nil { continue }

		eval := EvaluateThunderTaskError(item.RawText)
		if eval.Decision != DecisionRetire && (currTask.SummarySize != "" || currTask.Status == models.ThunderTaskHolding) {
			eval.Decision = DecisionHold
			if eval.SummarySize == "" { eval.SummarySize = currTask.SummarySize }
		}
		if currTask.DispatchedAt != nil && time.Since(*currTask.DispatchedAt) < 30*time.Second && eval.Decision == DecisionRetire {
			continue
		}

		if eval.Decision == DecisionRetire {
			if strings.Contains(item.RawText, "正在下载") || (strings.Contains(item.RawText, "/s") && !strings.Contains(item.RawText, "0B/s")) { continue }
			if currTask.Status != models.ThunderTaskRetired && a.Repo != nil {
				allRet, mID, err := a.Repo.MarkThunderTaskRetiredAndCheckAll(dbFileName, eval.Reason)
				if err == nil && allRet && mID != "" { _ = a.Repo.UpdateMediaMetadata(mID, "RETAINED", "", "", "全候補RETIREDにより退避") }
				a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("🛑 枯渇タスクをRETIRED: %s (%s)", item.FileName, eval.Reason))
				updatedCount++
			}
			if wsURL, err := FetchThunderMainRendererWSUrl(9222); err == nil && wsURL != "" { a.deleteTaskByFileNameSilent(wsURL, item.FileName) }
		} else if eval.Decision == DecisionHold && currTask.Status != models.ThunderTaskHolding && a.Repo != nil {
			_ = a.Repo.MarkThunderTaskHolding(dbFileName, eval.SummarySize, eval.Reason)
			updatedCount++
		}
	}

	// 2. ⚡ 幽霊タスクの回収: DBはONBOARDEDだが迅雷実機から消えたタスクを同期
	if a.Repo != nil && a.Repo.DB() != nil {
		var ghostTasks []models.ThunderTask
		_ = a.Repo.DB().Where("status IN ?", []string{string(models.ThunderTaskOnboarded), string(models.ThunderTaskRunning)}).Find(&ghostTasks).Error

		for _, gt := range ghostTasks {
			if liveMap[gt.FileName] || liveMap[cleanThunderFileName(gt.FileName)] { continue }
			if gt.DispatchedAt != nil && time.Since(*gt.DispatchedAt) < 20*time.Second { continue }

			if a.CheckMediaFileExists(gt.MediaID, gt.FileName) {
				_ = a.Repo.DB().Model(&models.ThunderTask{}).Where("id = ?", gt.ID).Update("status", models.ThunderTaskCompleted).Error
				a.AppendPipelineLog("THUNDER", "SUCCESS", fmt.Sprintf("🎉 実体ファイル確認完了: %s", gt.FileName))
			} else {
				allRet, mID, err := a.Repo.MarkThunderTaskRetiredAndCheckAll(gt.FileName, "迅雷不在（消滅検知）")
				if err == nil && allRet && mID != "" { _ = a.Repo.UpdateMediaMetadata(mID, "RETAINED", "", "", "全候補不在により退避") }
				a.AppendPipelineLog("THUNDER", "WARN", fmt.Sprintf("🗑️ 迅雷不在タスクをRETIRED同期: %s", gt.FileName))
			}
			updatedCount++
		}
	}
	return updatedCount, liveMap
}

func (a *App) getDBActiveMap() map[string]bool {
	m := make(map[string]bool)
	if a.Repo == nil || a.Repo.DB() == nil { return m }
	var tasks []models.ThunderTask
	_ = a.Repo.DB().Select("file_name").Where("status IN ?", []string{string(models.ThunderTaskOnboarded), string(models.ThunderTaskRunning)}).Find(&tasks).Error
	for _, t := range tasks { m[t.FileName] = true; m[cleanThunderFileName(t.FileName)] = true }
	return m
}

func (a *App) getDozouThunderTask(fileName string) *models.ThunderTask {
	if a.Repo == nil || a.Repo.DB() == nil || fileName == "" { return nil }
	var t models.ThunderTask
	if err := a.Repo.DB().Where("file_name = ?", fileName).First(&t).Error; err != nil { return nil }
	return &t
}
