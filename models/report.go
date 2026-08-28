// models/report.go (100行以下 - SPEC-PRINCIPLE-001)
package models

import "time"

// MissionEscalation は失敗時のエスカレーション・リカバリー先を表します
type MissionEscalation struct {
	ItemID   string `json:"item_id"`
	Reason   string `json:"reason"`
	Action   string `json:"action"` // "OUTSOURCED_ARIA2", "SMART_RECOVERY", "SKIP"
	TargetID string `json:"target_id,omitempty"`
}

// MissionReport はサルベージ/インポート完遂時の 5W1H 構造化レポートです
type MissionReport struct {
	ID           string              `json:"id"`
	FileName     string              `json:"file_name"`
	Type         string              `json:"type"`          // "salvage" | "import-manual"
	Platform     string              `json:"platform"`      // "twitter" | "bsky"
	Account      string              `json:"account"`       // "@sayapom4"
	Source       string              `json:"source"`        // "wayback+sotwe"
	StartedAt    time.Time           `json:"started_at"`
	FinishedAt   time.Time           `json:"finished_at"`
	DurationSec  float64             `json:"duration_sec"`
	TotalScanned int                 `json:"total_scanned"`
	SuccessCount int                 `json:"success_count"`
	FailedCount  int                 `json:"failed_count"`
	MarkdownText string              `json:"markdown_text"`
	Escalations  []MissionEscalation `json:"escalations,omitempty"`
}

// SystemJournalEntry はインフラ常駐処理用のオンメモリ構造化ジャーナルです
type SystemJournalEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Component string                 `json:"component"` // "stash" | "downloader" | "crawler" | "audit"
	Level     string                 `json:"level"`     // "INFO" | "WARN" | "ERROR"
	Event     string                 `json:"event"`     // "heartbeat_ok", "outsourced_polled" 等
	Message   string                 `json:"message"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}
