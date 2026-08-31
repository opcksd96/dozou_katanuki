// models/download_task.go (100行以下 - SPEC-PRINCIPLE-001)
package models

import "time"

type PipelineStage string

const (
	StageRequests PipelineStage = "REQUESTS"
	StageMotrix   PipelineStage = "MOTRIX"
	StageThunder  PipelineStage = "THUNDER"
	StageStash    PipelineStage = "STASH"
)

type DownloadTaskStatus string

const (
	TaskPending   DownloadTaskStatus = "PENDING"
	TaskRunning   DownloadTaskStatus = "RUNNING"
	TaskCompleted DownloadTaskStatus = "COMPLETED"
	TaskFailed    DownloadTaskStatus = "FAILED"
	TaskDepleted  DownloadTaskStatus = "DEPLETED"
	TaskReaped    DownloadTaskStatus = "REAPED"
)

// DownloadTask は 全パイプラインステージ（Requests/Motrix/Thunder/Stash）で共通の子タスクモデルです
type DownloadTask struct {
	ID             string             `gorm:"primaryKey" json:"id"` // {media_id}-{stage}-{resolution_type}
	MediaID        string             `gorm:"index;not null" json:"media_id"`
	ArticleID      string             `gorm:"index" json:"article_id"`
	Stage          PipelineStage      `gorm:"index;not null" json:"stage"`
	ResolutionType string             `gorm:"not null" json:"resolution_type"` // orig, large, medium, small, thumb, wayback_orig
	URL            string             `gorm:"not null" json:"url"`
	FileName       string             `gorm:"index;not null" json:"file_name"`
	Status         DownloadTaskStatus `gorm:"index;default:'PENDING'" json:"status"`
	FailedReason   string             `json:"failed_reason,omitempty"`
	DispatchedAt   *time.Time         `json:"dispatched_at,omitempty"`
	CompletedAt    *time.Time         `json:"completed_at,omitempty"`
	ReapedAt       *time.Time         `json:"reaped_at,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func (DownloadTask) TableName() string {
	return "download_tasks"
}
