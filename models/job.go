package models

import (
	"time"
)

// JobType はジョブの種類を表します
type JobType string

const (
	JobTypeSalvage       JobType = "salvage"
	JobTypeImportManual  JobType = "import-manual"
	JobTypeMediaDownload JobType = "download"
	JobTypeMediaPoll     JobType = "poll"
	JobTypeRestore       JobType = "restore"
	JobTypeScrape        JobType = "scrape"
	JobTypeTranslate     JobType = "translate"
	JobTypeCustom        JobType = "custom"
)

// JobStatus はジョブの実行ステータスを表します
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// JobRequest はジョブの起動要求を表します
type JobRequest struct {
	ID         string            `json:"id"`
	Type       JobType           `json:"type"`
	Platform   string            `json:"platform,omitempty"`
	Account    string            `json:"account,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	WARCPath   string            `json:"warc_path,omitempty"`
	Offline    bool              `json:"offline,omitempty"`
	ScriptPath string            `json:"script_path,omitempty"`
	Args       []string          `json:"args,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

// JobProgress はリアルタイムなジョブ進捗とステータスを表します（オンメモリ追跡用）
type JobProgress struct {
	ID         string     `json:"id"`
	Type       JobType    `json:"type"`
	Status     JobStatus  `json:"status"`
	Current    int        `json:"current"`
	Total      int        `json:"total"`
	Percentage float64    `json:"percentage"`
	Message    string     `json:"message"`
	Logs       []string   `json:"logs"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Error      string     `json:"error,omitempty"`
}
