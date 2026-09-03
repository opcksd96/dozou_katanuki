// app/app_rpc_thunder_reconciler.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

// ReconcileThunderTasksWithDB は 迅雷のCDPタスク一覧とDBを突合・更新し、エラー判定と取り下げを行います
func (a *App) ReconcileThunderTasksWithDB() (int, map[string]bool) {
	existingMap := make(map[string]bool)
	status := a.GetThunderCDPStatus()
	if !status.IsConnected || len(status.CapturedTasks) == 0 {
		return 0, existingMap
	}

	wsURL := status.ActiveWSUrl
	for _, item := range status.CapturedTasks {
		if item.FileName == "" {
			continue
		}
		if !a.isDozouManagedTask(item.FileName) {
			continue
		} // ユーザー独自タスクは除外
		existingMap[item.FileName] = true

		// ⑦ エラー文言 × サマリサイズ(>1B vs 0B) の判定
		eval := EvaluateThunderTaskError(item.RawText)
		switch eval.Decision {
		case DecisionRetire:
			// 原始资源不存在 かつ 0B: RETIRED付与、迅雷から削除、全候補全滅(ALL_TRUE)時のみRETAINED
			if a.Repo != nil {
				allRetired, mediaID, err := a.Repo.MarkThunderTaskRetiredAndCheckAll(item.FileName, eval.Reason)
				if err == nil && allRetired && mediaID != "" {
					_ = a.Repo.UpdateMediaMetadata(mediaID, "RETAINED", "", "", "全候補タスクRETIREDによりRETAINED退避")
					a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("📦 全候補枯渇(ALL_TRUE)のため退避: %s", mediaID))
				}
			}
			a.deleteTaskByFileNameSilent(wsURL, item.FileName)

		case DecisionHold:
			// 暂无任何有效资源 かつ >1B: タスク維持 & media ESCALATED維持
			if a.Repo != nil {
				_ = a.Repo.MarkThunderTaskHolding(item.FileName, eval.SummarySize, eval.Reason)
			}

		case DecisionCooldown:
			// 429 / ネットワーク異常: 10分クールダウン記録
			if a.Repo != nil {
				_ = a.Repo.UpdateThunderTaskCooldown(item.FileName, eval.Reason)
			}
		}
	}

	return len(existingMap), existingMap
}

// isDozouManagedTask は ユーザーが個別に追加した独自タスクを排除し、dozou管轄タスクのみ判定します
func (a *App) isDozouManagedTask(fileName string) bool {
	if a.Repo == nil || a.Repo.DB() == nil || fileName == "" {
		return false
	}
	var count int64
	_ = a.Repo.DB().Model(&models.ThunderTask{}).Where("file_name = ?", fileName).Count(&count).Error
	return count > 0
}

// CheckAndReonboardMissingTasks は DB上ONBOARDEDだが迅雷から消失したタスクを3個上限内で再投入します
func (a *App) CheckAndReonboardMissingTasks(existingMap map[string]bool, maxSlots int) int {
	if a.Repo == nil || a.Repo.DB() == nil {
		return 0
	}
	if maxSlots <= 0 || maxSlots > 3 {
		maxSlots = 3
	}

	currentCount := len(existingMap)
	if currentCount >= maxSlots {
		return 0
	}

	var onboardedTasks []models.ThunderTask
	_ = a.Repo.DB().Where("status = ?", models.ThunderTaskOnboarded).Limit(10).Find(&onboardedTasks).Error

	readded := 0
	destDir := a.getMediaDownloadDir()

	for _, t := range onboardedTasks {
		if currentCount >= maxSlots {
			break
		}
		if existingMap[t.FileName] {
			continue
		}
		if IsThunderCooldownActive(t.LastAttemptAt) {
			continue
		} // 10分冷却中ならスキップ

		if AddTaskViaThunderCOM(t.URL, t.FileName, destDir) {
			existingMap[t.FileName] = true
			currentCount++
			readded++
			_ = a.Repo.MarkThunderTaskOnboarded(t.ID, t.SummarySize)
			a.AppendPipelineLog("THUNDER", "INFO", fmt.Sprintf("⚡ 欠損タスク再登録: %s", t.FileName))
		}
	}
	return readded
}
