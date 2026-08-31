// models/thunder_task.go (100行以下 - SPEC-PRINCIPLE-001)
package models

import "time"

// ThunderTaskStatus はタスクのライフサイクル状態を表します
type ThunderTaskStatus string

const (
	ThunderTaskPending   ThunderTaskStatus = "PENDING"
	ThunderTaskRunning   ThunderTaskStatus = "RUNNING"
	ThunderTaskCompleted ThunderTaskStatus = "COMPLETED"
	ThunderTaskDepleted  ThunderTaskStatus = "DEPLETED"
	ThunderTaskReaped    ThunderTaskStatus = "REAPED"
)

// ThunderTask は 1つのメディアに対して投入される候補URL（ワーカータスク）を表すモデルです
type ThunderTask struct {
	ID             string            `gorm:"primaryKey" json:"id"`
	MediaID        string            `gorm:"index;not null" json:"media_id"`
	ArticleID      string            `gorm:"index" json:"article_id"`
	ResolutionType string            `gorm:"not null" json:"resolution_type"`
	URL            string            `gorm:"not null" json:"url"`
	FileName       string            `gorm:"index;not null" json:"file_name"`
	Status         ThunderTaskStatus `gorm:"index;default:'PENDING'" json:"status"`
	DispatchedAt   *time.Time        `json:"dispatched_at,omitempty"`
	CompletedAt    *time.Time        `json:"completed_at,omitempty"`
	ReapedAt       *time.Time        `json:"reaped_at,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func (ThunderTask) TableName() string {
	return "thunder_tasks"
}
