// app/app_rpc_pipeline_metrics.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"dozou_katanuki/models"
)

// CalculatePipelineMetrics はユーザー定義の関係式に基づいてメトリクスを集計・算出します
func (a *App) CalculatePipelineMetrics() models.PipelineMetrics {
	var metrics models.PipelineMetrics
	if a.Repo == nil || a.Repo.DB() == nil {
		return metrics
	}
	db := a.Repo.DB()

	// 1. 作品数: mediaテーブルにおける代表IDの数 (探索中・全体)
	_ = db.Model(&models.Media{}).Count(&metrics.TotalMediaCount).Error
	var escalatedCount int64
	_ = db.Model(&models.Media{}).Where("download_status = 'ESCALATED'").Count(&escalatedCount).Error
	metrics.MediaCount = escalatedCount
	if metrics.MediaCount == 0 && metrics.TotalMediaCount > 0 {
		metrics.MediaCount = metrics.TotalMediaCount
	}

	// 2. タスク数: 作品数 * 7
	metrics.TotalTasks = metrics.MediaCount * 7

	// 3. 稼働数 (ONBOARDED) & 継続探査数 (PENDING)
	_ = db.Model(&models.ThunderTask{}).Where("status = 'ONBOARDED'").Count(&metrics.ActiveTasks).Error
	_ = db.Model(&models.ThunderTask{}).Where("status = 'PENDING'").Count(&metrics.PendingTasks).Error

	// 下載登録数 = 稼働数 + 継続探査数
	metrics.RegisteredTasks = metrics.ActiveTasks + metrics.PendingTasks

	// 4. 成功数 (COMPLETED / REAPED) & 失敗数 (RETIRED)
	_ = db.Model(&models.ThunderTask{}).Where("status IN ('REAPED', 'COMPLETED')").Count(&metrics.SuccessTasks).Error
	_ = db.Model(&models.ThunderTask{}).Where("status IN ('RETIRED', 'DEPLETED')").Count(&metrics.FailedTasks).Error

	// 関係式整合: タスク数 ＝ 下載登録数 + 成功数 + 失敗数
	sumTasks := metrics.RegisteredTasks + metrics.SuccessTasks + metrics.FailedTasks
	if sumTasks > metrics.TotalTasks {
		metrics.TotalTasks = sumTasks
	}

	return metrics
}
