// app/app_rpc_downloaders.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

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

// ControlMotrix は Motrix (Aria2) に対する制御コマンドを実行します (pause_all / unpause_all / purge_all / safe_limits)
func (a *App) ControlMotrix(action string) (bool, error) {
	endpoints := []string{"http://localhost:16800/jsonrpc", "http://localhost:6800/jsonrpc"}
	client := &http.Client{Timeout: 3 * time.Second}
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

	for _, ep := range endpoints {
		payload, _ := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": "ctrl", "method": method, "params": params})
		if resp, err := client.Post(ep, "application/json", bytes.NewReader(payload)); err == nil {
			resp.Body.Close()
			return true, nil
		}
	}
	return false, nil
}

// LaunchThunder は Thunder.exe を起動します
func (a *App) LaunchThunder() (bool, error) {
	thunderPath := `C:\Program Files (x86)\Thunder Network\Thunder\Program\Thunder.exe`
	if _, err := os.Stat(thunderPath); err != nil { return false, err }
	cmd := exec.Command(thunderPath)
	return cmd.Start() == nil, nil
}

func (a *App) fetchMotrixStatus() models.MotrixGlobalStat {
	endpoints := []string{"http://localhost:16800/jsonrpc", "http://localhost:6800/jsonrpc"}
	client := &http.Client{Timeout: 2 * time.Second}
	for _, ep := range endpoints {
		payload, _ := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": "stat", "method": "aria2.getGlobalStat"})
		resp, err := client.Post(ep, "application/json", bytes.NewReader(payload))
		if err != nil { continue }
		defer resp.Body.Close()

		var res struct {
			Result struct {
				DownloadSpeed string `json:"downloadSpeed"`
				UploadSpeed   string `json:"uploadSpeed"`
				NumActive     string `json:"numActive"`
				NumWaiting    string `json:"numWaiting"`
				NumStopped    string `json:"numStopped"`
			} `json:"result"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
			ds, _ := strconv.ParseInt(res.Result.DownloadSpeed, 10, 64)
			us, _ := strconv.ParseInt(res.Result.UploadSpeed, 10, 64)
			na, _ := strconv.Atoi(res.Result.NumActive)
			nw, _ := strconv.Atoi(res.Result.NumWaiting)
			ns, _ := strconv.Atoi(res.Result.NumStopped)
			tasks := a.fetchMotrixActiveTasks(client, ep)
			return models.MotrixGlobalStat{IsOnline: true, DownloadSpeed: ds, UploadSpeed: us, NumActive: na, NumWaiting: nw, NumStopped: ns, ActiveTasks: tasks}
		}
	}
	return models.MotrixGlobalStat{IsOnline: false}
}

func (a *App) fetchMotrixActiveTasks(client *http.Client, ep string) []models.DownloaderTaskInfo {
	payload, _ := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": "active", "method": "aria2.tellActive", "params": []interface{}{[]string{"gid", "status", "totalLength", "completedLength", "downloadSpeed", "files", "errorMessage"}}})
	resp, err := client.Post(ep, "application/json", bytes.NewReader(payload))
	if err != nil { return nil }
	defer resp.Body.Close()

	var res struct {
		Result []struct {
			GID             string `json:"gid"`
			Status          string `json:"status"`
			TotalLength     string `json:"totalLength"`
			CompletedLength string `json:"completedLength"`
			DownloadSpeed   string `json:"downloadSpeed"`
			ErrorMessage    string `json:"errorMessage"`
			Files           []struct {
				Path string `json:"path"`
			} `json:"files"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil { return nil }

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
