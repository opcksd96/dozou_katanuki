// app/app_rpc_downloaders.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"dozou_katanuki/models"
)

// GetDownloaderDashboardStatus は Motrix (Aria2) と Thunder の最新ステータスを取得します
func (a *App) GetDownloaderDashboardStatus() (models.DownloaderDashboardStatus, error) {
	var result models.DownloaderDashboardStatus
	result.Motrix = a.fetchMotrixStatus()

	thunderPath := `C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe`
	if _, err := os.Stat(thunderPath); err == nil {
		result.Thunder.IsInstalled = true
		result.Thunder.Executable = thunderPath
	}
	if stats, err := a.Repo.FetchDownloadStatusStats(""); err == nil && stats != nil {
		result.Thunder.RetainedCount = stats.Retained
	}
	return result, nil
}

// ControlMotrix は Motrix に対する制御コマンドを実行します
func (a *App) ControlMotrix(action string) (bool, error) {
	var method string
	var params []interface{}

	switch action {
	case "pause_all": method = "aria2.forcePauseAll"
	case "unpause_all": method = "aria2.unpauseAll"
	case "purge_all": method = "aria2.purgeDownloadResult"
	case "safe_limits":
		method = "aria2.changeGlobalOption"
		params = []interface{}{map[string]string{
			"max-concurrent-downloads": "2", "max-connection-per-server": "1",
			"split": "1", "retry-wait": "5", "max-tries": "3",
		}}
	default: return false, nil
	}

	_, err := callMotrixRPC(method, params)
	return err == nil, err
}

// LaunchThunder は Thunder.exe を起動します
func (a *App) LaunchThunder() (bool, error) {
	thunderPath := `C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe`
	if _, err := os.Stat(thunderPath); err != nil { return false, err }
	cmd := exec.Command(thunderPath)
	return cmd.Start() == nil, nil
}

func (a *App) fetchMotrixStatus() models.MotrixGlobalStat {
	raw, err := callMotrixRPC("aria2.getGlobalStat", nil)
	if err != nil { return models.MotrixGlobalStat{IsOnline: false} }

	var res struct {
		Result struct {
			DownloadSpeed string `json:"downloadSpeed"`
			UploadSpeed   string `json:"uploadSpeed"`
			NumActive     string `json:"numActive"`
			NumWaiting    string `json:"numWaiting"`
			NumStopped    string `json:"numStopped"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil { return models.MotrixGlobalStat{IsOnline: false} }

	ds, _ := strconv.ParseInt(res.Result.DownloadSpeed, 10, 64)
	us, _ := strconv.ParseInt(res.Result.UploadSpeed, 10, 64)
	na, _ := strconv.Atoi(res.Result.NumActive)
	nw, _ := strconv.Atoi(res.Result.NumWaiting)
	ns, _ := strconv.Atoi(res.Result.NumStopped)
	tasks := a.fetchMotrixActiveTasks()
	return models.MotrixGlobalStat{IsOnline: true, DownloadSpeed: ds, UploadSpeed: us, NumActive: na, NumWaiting: nw, NumStopped: ns, ActiveTasks: tasks}
}

func (a *App) fetchMotrixActiveTasks() []models.DownloaderTaskInfo {
	raw, err := callMotrixRPC("aria2.tellActive", []interface{}{[]string{"gid", "status", "totalLength", "completedLength", "downloadSpeed", "files", "errorMessage"}})
	if err != nil { return nil }

	var res struct {
		Result []struct {
			GID, Status, TotalLength, CompletedLength, DownloadSpeed, ErrorMessage string
			Files []struct{ Path string `json:"path"` } `json:"files"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &res); err != nil { return nil }

	var tasks []models.DownloaderTaskInfo
	for _, t := range res.Result {
		tl, _ := strconv.ParseInt(t.TotalLength, 10, 64)
		cl, _ := strconv.ParseInt(t.CompletedLength, 10, 64)
		ds, _ := strconv.ParseInt(t.DownloadSpeed, 10, 64)
		fn := ""
		if len(t.Files) > 0 { fn = filepath.Base(t.Files[0].Path) }
		prog := 0.0
		if tl > 0 { prog = float64(cl) / float64(tl) * 100 }
		tasks = append(tasks, models.DownloaderTaskInfo{GID: t.GID, Status: t.Status, FileName: fn, TotalLength: tl, CompletedLength: cl, DownloadSpeed: ds, Progress: prog, ErrorMessage: t.ErrorMessage})
	}
	return tasks
}
