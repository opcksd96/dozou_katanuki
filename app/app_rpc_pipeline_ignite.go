// app/app_rpc_pipeline_ignite.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
)

type PipelineIgniteResult struct {
	QueuedCount    int64 `json:"queued_count"`
	EscalatedCount int64 `json:"escalated_count"`
	IsIgnited      bool  `json:"is_ignited"`
}

// IgnitePipeline は統括契約に基づき、全ステージを即時1サイクル完走します
func (a *App) IgnitePipeline() (*PipelineIgniteResult, error) {
	cycleRes, err := a.ExecutePipelineCycleNow()
	if err != nil {
		return nil, err
	}
	if cycleRes == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	a.AppendPipelineLog("SYSTEM", "INFO", fmt.Sprintf("⚡ パイプライン即時実行: QUEUED %d 件, ESCALATED %d 件 を順次処理しました", cycleRes.QueuedCount, cycleRes.EscalatedCount))

	return &PipelineIgniteResult{
		QueuedCount:    cycleRes.QueuedCount,
		EscalatedCount: cycleRes.EscalatedCount,
		IsIgnited:      cycleRes.Success,
	}, nil
}
