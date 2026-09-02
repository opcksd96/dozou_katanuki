// app/app_rpc_pipeline_ignite.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"

	"dozou_katanuki/models"
)

type PipelineIgniteResult struct {
	QueuedCount    int64 `json:"queued_count"`
	EscalatedCount int64 `json:"escalated_count"`
	IsIgnited      bool  `json:"is_ignited"`
}

// IgnitePipeline は DBの進捗状態に基づき、全ステージのキューを途中から一斉に再発火します
func (a *App) IgnitePipeline() (*PipelineIgniteResult, error) {
	if a.Repo == nil || a.Repo.DB() == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	db := a.Repo.DB()

	var qCount, eCount int64
	_ = db.Model(&models.Media{}).Where("download_status = 'QUEUED' AND (is_trash = 0 OR is_trash IS NULL)").Count(&qCount).Error
	_ = db.Model(&models.Media{}).Where("download_status = 'ESCALATED' AND (is_trash = 0 OR is_trash IS NULL)").Count(&eCount).Error

	res := &PipelineIgniteResult{
		QueuedCount:    qCount,
		EscalatedCount: eCount,
		IsIgnited:      true,
	}

	a.AppendPipelineLog("SYSTEM", "INFO", fmt.Sprintf("🔥 パイプライン点火(Ignite): QUEUED %d 件, ESCALATED %d 件 を各ステージから再発火！", qCount, eCount))

	// 1. QUEUED メディアがあれば第1ステージ (Requests ➔ Motrix) を非同期発火
	if qCount > 0 {
		_, _ = a.ProcessQueuedViaRequests()
	}

	// 2. Motrix (OUTSOURCED) の完了同期をキック
	_, _ = a.SyncCompletedDownloads()

	// 3. ESCALATED メディアがあれば第3ステージ (迅雷オーケストレーター) を発火
	if eCount > 0 && !a.isThunderOrchestratorRunning() {
		_, _ = a.StartThunderOrchestrator(3, 4)
	}

	// 4. 完了済み未同期があれば第4ステージ (Stash) を同期
	go func() {
		_, _ = a.SyncThunderDownloads("")
	}()

	return res, nil
}
