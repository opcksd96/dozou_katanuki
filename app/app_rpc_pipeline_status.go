// app/app_rpc_pipeline_status.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"dozou_katanuki/models"
)

type CheckpointStatus struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	IsOnline    bool   `json:"is_online"`
	ActiveCount int    `json:"active_count"`
	TotalCount  int    `json:"total_count"`
	SpeedText   string `json:"speed_text,omitempty"`
	StatusText  string `json:"status_text"`
	SummaryText string `json:"summary_text,omitempty"`
}

type PipelineOverviewResult struct {
	Checkpoints         []CheckpointStatus     `json:"checkpoints"`
	TotalMedia          int64                  `json:"total_media"`
	Completed           int64                  `json:"completed"`
	Escalated           int64                  `json:"escalated"`
	Outsourced          int64                  `json:"outsourced"`
	Retained            int64                  `json:"retained"`
	OverallProgress     float64                `json:"overall_progress"`
	IsAutoEngineRunning bool                   `json:"is_auto_engine_running"`
	Metrics             models.PipelineMetrics `json:"metrics"`
}

func (a *App) GetPipelineOverview() (*PipelineOverviewResult, error) {
	res := &PipelineOverviewResult{Checkpoints: make([]CheckpointStatus, 4), IsAutoEngineRunning: a.IsPipelineAutoEngineRunning(), Metrics: a.CalculatePipelineMetrics()}
	if a.Repo != nil && a.Repo.DB() != nil {
		db := a.Repo.DB()
		_ = db.Model(&models.Media{}).Count(&res.TotalMedia).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'COMPLETED'").Count(&res.Completed).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'ESCALATED'").Count(&res.Escalated).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'OUTSOURCED'").Count(&res.Outsourced).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'RETAINED'").Count(&res.Retained).Error
		if res.TotalMedia > 0 { res.OverallProgress = float64(res.Completed) / float64(res.TotalMedia) * 100.0 }
	}

	// 1. Requests
	res.Checkpoints[0] = CheckpointStatus{
		Name: "Requests / 内蔵HTTP", Key: "requests", IsOnline: true,
		ActiveCount: int(res.Retained), TotalCount: int(res.TotalMedia), StatusText: "STANDBY",
		SummaryText: fmt.Sprintf("リキュー待ち: %d 件 / 全体: %d 件", res.Retained, res.TotalMedia),
	}

	// 2. Motrix Next
	mStat := a.fetchMotrixStatus()
	mText := "OFFLINE"
	if mStat.IsOnline { mText = "ONLINE" }
	res.Checkpoints[1] = CheckpointStatus{
		Name: "Motrix Next / Aria2", Key: "motrix", IsOnline: mStat.IsOnline,
		ActiveCount: mStat.NumActive, TotalCount: mStat.NumActive + mStat.NumWaiting + mStat.NumStopped,
		StatusText: mText, SummaryText: fmt.Sprintf("実行: %d / 待機: %d / 停止: %d", mStat.NumActive, mStat.NumWaiting, mStat.NumStopped),
	}

	// 3. Thunder
	thunderProc := isThunderProcessRunning()
	thunderCDP := false
	if _, err := FetchThunderMainRendererWSUrl(9222); err == nil { thunderCDP = true }
	tText := "OFFLINE"
	if orchState.isRunning {
		if orchState.isPaused { tText = "PAUSED" } else { tText = "RUNNING" }
	} else if thunderCDP { tText = "CDP CONNECTED" } else if thunderProc { tText = "ONLINE" }

	tCounts, _ := a.Repo.GetThunderTaskCounts()
	tActive := tCounts.Running
	if len(orchState.recentTasks) > 0 { tActive = len(orchState.recentTasks) }
	tSummary := fmt.Sprintf("進行: %d / 待機: %d / 枯渇: %d", tCounts.Running, tCounts.Pending, tCounts.Failed)
	if len(orchState.recentTasks) > 0 {
		tSummary = fmt.Sprintf("アプリ: %d件 (進行: %d / 待機: %d)", len(orchState.recentTasks), tCounts.Running, tCounts.Pending)
	}
	res.Checkpoints[2] = CheckpointStatus{
		Name: "迅雷 (Thunder) P2SP", Key: "thunder", IsOnline: thunderProc || thunderCDP,
		ActiveCount: tActive, TotalCount: tCounts.Total, StatusText: tText, SummaryText: tSummary,
	}

	// 4. Stash
	sText := "OFFLINE"
	if a.isStashServerOnline() { sText = "ONLINE" }
	res.Checkpoints[3] = CheckpointStatus{
		Name: "Stashapp DB / Assets", Key: "stash", IsOnline: a.isStashServerOnline(),
		ActiveCount: int(res.Completed), TotalCount: int(res.TotalMedia), StatusText: sText,
		SummaryText: fmt.Sprintf("登録済: %d 件 (進捗: %.1f%%)", res.Completed, res.OverallProgress),
	}

	return res, nil
}
