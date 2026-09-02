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

// DownloadTask は 1メディア1行でライフサイクル・ステージ通過時刻を管理するタスクモデルです
type DownloadTask struct {
	MediaID      string             `gorm:"primaryKey;column:media_id" json:"media_id"`
	ArticleID    string             `gorm:"index;column:article_id" json:"article_id"`
	URL          string             `gorm:"column:url;not null" json:"url"` // 原本URL (Single Source of Truth)
	FileName     string             `gorm:"column:file_name;not null" json:"file_name"`
	Stage        PipelineStage      `gorm:"column:stage;index;default:'REQUESTS'" json:"stage"`
	Status       DownloadTaskStatus `gorm:"column:status;index;default:'PENDING'" json:"status"`
	FailedReason string             `gorm:"column:failed_reason" json:"failed_reason,omitempty"`
	RequestsAt   *time.Time         `gorm:"column:requests_at" json:"requests_at,omitempty"`
	MotrixAt     *time.Time         `gorm:"column:motrix_at" json:"motrix_at,omitempty"`
	ThunderAt    *time.Time         `gorm:"column:thunder_at" json:"thunder_at,omitempty"`
	StashAt      *time.Time         `gorm:"column:stash_at" json:"stash_at,omitempty"`
	CompletedAt  *time.Time         `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt    time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"column:updated_at" json:"updated_at"`
}

func (DownloadTask) TableName() string {
	return "download_tasks"
}
