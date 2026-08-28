// models/downloaders.go (100行以下 - SPEC-PRINCIPLE-001)
package models

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

// DownloaderDashboardStatus は管理コンソール向けの総合ステータスです
type DownloaderDashboardStatus struct {
	Motrix  MotrixGlobalStat  `json:"motrix"`
	Thunder ThunderGlobalStat `json:"thunder"`
}
