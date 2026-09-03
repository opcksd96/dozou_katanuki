// app/app_rpc_pipeline_status.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"dozou_katanuki/models"
)

type CheckpointStatus struct {
	Name        string `json:"name"` // Requests, Motrix Next, Thunder, Stash
	Key         string `json:"key"`  // requests, motrix, thunder, stash
	IsOnline    bool   `json:"is_online"`
	ActiveCount int    `json:"active_count"`
	TotalCount  int    `json:"total_count"`
	SpeedText   string `json:"speed_text"`
	StatusText  string `json:"status_text"`
}

type PipelineOverviewResult struct {
	Checkpoints     []CheckpointStatus `json:"checkpoints"`
	TotalMedia      int64              `json:"total_media"`
	Completed       int64              `json:"completed"`
	Escalated       int64              `json:"escalated"`
	Outsourced      int64              `json:"outsourced"`
	Retained        int64              `json:"retained"`
	OverallProgress float64            `json:"overall_progress"`
}

// GetPipelineOverview は パイプライン全体の4大チェックポイント稼働状態とメディア集計を一括返却します
func (a *App) GetPipelineOverview() (*PipelineOverviewResult, error) {
	res := &PipelineOverviewResult{Checkpoints: make([]CheckpointStatus, 4)}

	// 1. Requests (内蔵)
	res.Checkpoints[0] = CheckpointStatus{Name: "Requests / 内蔵HTTP", Key: "requests", IsOnline: true, StatusText: "🟢 STANDBY"}

	// 2. Motrix Next
	mStat := a.fetchMotrixStatus()
	mText := "🔴 OFFLINE"
	if mStat.IsOnline {
		mText = "🟢 ONLINE"
	}
	res.Checkpoints[1] = CheckpointStatus{Name: "Motrix Next / Aria2", Key: "motrix", IsOnline: mStat.IsOnline, ActiveCount: mStat.NumActive, StatusText: mText}

	// 3. Thunder
	thunderProc := isThunderProcessRunning()
	thunderCDP := false
	if _, err := FetchThunderMainRendererWSUrl(9222); err == nil {
		thunderCDP = true
	}
	thunderOnline := thunderProc || thunderCDP

	tText := "🔴 OFFLINE"
	if orchState.isRunning {
		tText = "⚡ RUNNING"
	} else if thunderCDP {
		tText = "⚡ CDP CONNECTED"
	} else if thunderProc {
		tText = "🟢 ONLINE"
	}
	res.Checkpoints[2] = CheckpointStatus{Name: "迅雷 (Thunder) P2SP", Key: "thunder", IsOnline: thunderOnline, ActiveCount: len(orchState.recentTasks), StatusText: tText}

	// 4. Stash
	stashOnline := a.isStashServerOnline()
	sText := "🔴 OFFLINE"
	if stashOnline {
		sText = "🟢 ONLINE"
	}
	res.Checkpoints[3] = CheckpointStatus{Name: "Stashapp DB / Assets", Key: "stash", IsOnline: stashOnline, StatusText: sText}

	// DB 集計
	if a.Repo != nil && a.Repo.DB() != nil {
		db := a.Repo.DB()
		_ = db.Model(&models.Media{}).Count(&res.TotalMedia).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'COMPLETED'").Count(&res.Completed).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'ESCALATED'").Count(&res.Escalated).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'OUTSOURCED'").Count(&res.Outsourced).Error
		_ = db.Model(&models.Media{}).Where("download_status = 'RETAINED'").Count(&res.Retained).Error

		if res.TotalMedia > 0 {
			res.OverallProgress = float64(res.Completed) / float64(res.TotalMedia) * 100.0
		}
	}

	return res, nil
}
