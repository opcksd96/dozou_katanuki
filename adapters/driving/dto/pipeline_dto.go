// adapters/driving/dto/pipeline_dto.go (100行以下 - SPEC-PRINCIPLE-001)
package dto

import "time"

// CheckpointStatus はダッシュボードの各サービスのオンライン状態を表します
type CheckpointStatus struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	IsOnline    bool   `json:"is_online"`
	ActiveCount int    `json:"active_count"`
	TotalCount  int    `json:"total_count"`
	SpeedText   string `json:"speed_text"`
	StatusText  string `json:"status_text"`
}

// PipelineOverviewDTO はパイプライン全体の稼働状態とメディア集計です
type PipelineOverviewDTO struct {
	Checkpoints     []CheckpointStatus `json:"checkpoints"`
	TotalMedia      int64              `json:"total_media"`
	Completed       int64              `json:"completed"`
	Escalated       int64              `json:"escalated"`
	Outsourced      int64              `json:"outsourced"`
	Retained        int64              `json:"retained"`
	OverallProgress float64            `json:"overall_progress"`
}

// DownloaderTaskInfo は Motrix/Aria2 内の個別タスク情報です
type DownloaderTaskInfo struct {
	GID             string  `json:"gid"`
	Status          string  `json:"status"`
	FileName        string  `json:"file_name"`
	URL             string  `json:"url,omitempty"`
	TotalLength     int64   `json:"total_length"`
	CompletedLength int64   `json:"completed_length"`
	DownloadSpeed   int64   `json:"download_speed"`
	Progress        float64 `json:"progress"`
	ErrorMessage    string  `json:"error_message,omitempty"`
}

// MotrixGlobalStat は Aria2 全体の統計情報です
type MotrixGlobalStat struct {
	IsOnline      bool                 `json:"is_online"`
	DownloadSpeed int64                `json:"download_speed"`
	UploadSpeed   int64                `json:"upload_speed"`
	NumActive     int                  `json:"num_active"`
	NumWaiting    int                  `json:"num_waiting"`
	NumStopped    int                  `json:"num_stopped"`
	ActiveTasks   []DownloaderTaskInfo `json:"active_tasks"`
}

// ThunderGlobalStat は Thunder (迅雷) の状態情報です
type ThunderGlobalStat struct {
	IsInstalled    bool   `json:"is_installed"`
	Executable     string `json:"executable"`
	EscalatedCount int64  `json:"escalated_count"`
	RetainedCount  int64  `json:"retained_count"`
}

// DownloaderDashboardDTO は管理コンソール向けの総合ステータスです
type DownloaderDashboardDTO struct {
	Motrix  MotrixGlobalStat  `json:"motrix"`
	Thunder ThunderGlobalStat `json:"thunder"`
}

// JobProgressDTO はリアルタイムなジョブ進捗とステータスを表します
type JobProgressDTO struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	Current    int        `json:"current"`
	Total      int        `json:"total"`
	Percentage float64    `json:"percentage"`
	Message    string     `json:"message"`
	Logs       []string   `json:"logs"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Error      string     `json:"error,omitempty"`
}
