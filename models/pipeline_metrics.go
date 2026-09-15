// models/pipeline_metrics.go (100行以下 - SPEC-PRINCIPLE-001)
package models

// PipelineMetrics はパイプラインの厳格な関係式モデルを保持するDTO
type PipelineMetrics struct {
	MediaCount      int64 `json:"media_count"`      // 作品数 (現在探索対象作品数)
	TotalMediaCount int64 `json:"total_media_count"`// 全作品数 (全件 2622等)
	TotalTasks      int64 `json:"total_tasks"`      // タスク数 (作品数 * 7)
	RegisteredTasks int64 `json:"registered_tasks"` // 下載登録数 (稼働数 + 継続探査数)
	ActiveTasks     int64 `json:"active_tasks"`     // 稼働数 (ONBOARDED)
	PendingTasks    int64 `json:"pending_tasks"`    // 継続探査数 (PENDING)
	SuccessTasks    int64 `json:"success_tasks"`    // 成功数 (COMPLETED / REAPED)
	FailedTasks     int64 `json:"failed_tasks"`     // 失敗数 (RETIRED / RETAINED)
}
