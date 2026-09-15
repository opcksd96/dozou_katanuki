// models/candidate_task_enriched.go (100行以下 - SPEC-PRINCIPLE-001)
package models

// CandidateTaskEnriched は DBのThunderTaskに迅雷CDPの完全な生タスク情報(THUNDER_TASK_SCHEMA)を付与したモデル
type CandidateTaskEnriched struct {
	ThunderTask
	CDP *ActiveThunderTask `json:"cdp,omitempty"`
}

// MediaWithTasksEnriched は作品と候補URLタスク(CDP生情報付き)の階層モデル
type MediaWithTasksEnriched struct {
	MediaID        string                  `json:"media_id"`
	ArticleID      string                  `json:"article_id"`
	OriginalURL    string                  `json:"original_url"`
	Type           string                  `json:"type"`
	Width          int                     `json:"width"`
	Height         int                     `json:"height"`
	DownloadStatus string                  `json:"download_status"`
	FailedReason   string                  `json:"failed_reason"`
	AccountID      string                  `json:"account_id"`
	MediaQuality   string                  `json:"media_quality"`
	CandidateTasks []CandidateTaskEnriched `json:"candidate_tasks"`
}
