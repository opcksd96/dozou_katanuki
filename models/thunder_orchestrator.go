package models

import "time"

// ThunderTaskResolutionType は厳選された解像度種別です
type ThunderResolutionType string

const (
	ResolutionOrig        ThunderResolutionType = "orig"
	ResolutionLarge       ThunderResolutionType = "large"
	ResolutionWaybackOrig ThunderResolutionType = "wayback_orig"
)

// ThunderOrchestratorTask はオーケストレーターが管理する単一ジョブです
type ThunderOrchestratorTask struct {
	ID             string                `json:"id"`
	MediaID        string                `json:"media_id"`
	ArticleID      string                `json:"article_id"`
	ResolutionType ThunderResolutionType `json:"resolution_type"`
	URL            string                `json:"url"`
	FileName       string                `json:"file_name"`
	Status         string                `json:"status"` // pending, running, success, failed, skipped
	SlotIndex      int                   `json:"slot_index"`
	DispatchedAt   *time.Time            `json:"dispatched_at,omitempty"`
	CompletedAt    *time.Time            `json:"completed_at,omitempty"`
	ErrorMessage   string                `json:"error_message,omitempty"`
}

// ThunderOrchestratorSlot は最大12スロットの実行枠です
type ThunderOrchestratorSlot struct {
	Index        int                      `json:"index"`
	IsOccupied   bool                     `json:"is_occupied"`
	CurrentTask  *ThunderOrchestratorTask `json:"current_task,omitempty"`
	DispatchedAt *time.Time               `json:"dispatched_at,omitempty"`
}

// ThunderOrchestratorConfig はオーケストレーションの制御設定です
type ThunderOrchestratorConfig struct {
	MaxConcurrentSlots int `json:"max_concurrent_slots"` // デフォルト: 12
	IntervalSeconds    int `json:"interval_seconds"`      // 間欠ディスパッチ間隔 (秒, デフォルト: 5)
	TopResolutionsOnly bool `json:"top_resolutions_only"` // 厳選3種類 (orig, large, wayback_orig) のみ対象
}

// ThunderOrchestratorStatus はフロントエンドへ配信するリアルタイムステータスです
type ThunderOrchestratorStatus struct {
	IsRunning       bool                      `json:"is_running"`
	IsPaused        bool                      `json:"is_paused"`
	Config          ThunderOrchestratorConfig `json:"config"`
	TotalJobs       int                       `json:"total_jobs"`
	PendingJobs     int                       `json:"pending_jobs"`
	RunningJobs     int                       `json:"running_jobs"`
	SuccessJobs     int                       `json:"success_jobs"`
	FailedJobs      int                       `json:"failed_jobs"`
	TotalMediaCount int                       `json:"total_media_count"`
	Slots           []ThunderOrchestratorSlot `json:"slots"`
	RecentTasks     []ThunderOrchestratorTask `json:"recent_tasks"`
}
