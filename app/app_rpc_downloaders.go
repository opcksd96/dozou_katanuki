// app/app_rpc_downloaders.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"os"
	"strconv"

	"dozou_katanuki/models"
)

// GetDownloaderDashboardStatus は Motrix と Thunder の実態ステータスを取得します
func (a *App) GetDownloaderDashboardStatus() (models.DownloaderDashboardStatus, error) {
	var result models.DownloaderDashboardStatus
	result.Motrix = a.fetchMotrixStatus()
	thunderPath := `C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe`
	if _, err := os.Stat(thunderPath); err == nil {
		result.Thunder.IsInstalled, result.Thunder.Executable = true, thunderPath
	}
	result.Thunder.IsConnected = a.GetThunderCDPStatus().IsConnected

	if a.Repo != nil {
		if c, err := a.Repo.GetThunderTaskCounts(); err == nil {
			result.Thunder.TotalCount = c.Total
			result.Thunder.PendingCount = c.Pending
			result.Thunder.RunningCount = c.Running
			result.Thunder.FailedCount = c.Failed
		}
		if stats, err := a.Repo.FetchDownloadStatusStats(""); err == nil && stats != nil {
			result.Thunder.EscalatedCount = stats.Escalated
			result.Thunder.RetainedCount = stats.Retained
		}
	}
	return result, nil
}

// ControlMotrix は Motrix に対する制御コマンドを実行します
func (a *App) ControlMotrix(action string) (bool, error) {
	var method string
	var params []interface{}
	switch action {
	case "pause_all":
		method = "aria2.forcePauseAll"
	case "unpause_all":
		method = "aria2.unpauseAll"
	case "purge_all":
		method = "aria2.purgeDownloadResult"
	case "safe_limits":
		method = "aria2.changeGlobalOption"
		params = []interface{}{map[string]string{
			"max-concurrent-downloads": "2", "max-connection-per-server": "1",
			"split": "1", "retry-wait": "5", "max-tries": "3",
		}}
	default:
		return false, nil
	}
	_, err := callMotrixRPC(method, params)
	return err == nil, err
}

func (a *App) fetchMotrixStatus() models.MotrixGlobalStat {
	raw, err := callMotrixRPC("aria2.getGlobalStat", nil)
	if err != nil { return models.MotrixGlobalStat{IsOnline: false} }
	var res struct {
		Result struct { DownloadSpeed, UploadSpeed, NumActive, NumWaiting, NumStopped string } `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil { return models.MotrixGlobalStat{IsOnline: false} }
	ds, _ := strconv.ParseInt(res.Result.DownloadSpeed, 10, 64)
	us, _ := strconv.ParseInt(res.Result.UploadSpeed, 10, 64)
	na, _ := strconv.Atoi(res.Result.NumActive)
	nw, _ := strconv.Atoi(res.Result.NumWaiting)
	ns, _ := strconv.Atoi(res.Result.NumStopped)
	return models.MotrixGlobalStat{
		IsOnline: true, DownloadSpeed: ds, UploadSpeed: us,
		NumActive: na, NumWaiting: nw, NumStopped: ns,
		ActiveTasks: a.fetchMotrixActiveTasks(),
	}
}
