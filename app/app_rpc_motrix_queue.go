// app/app_rpc_motrix_queue.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"

	"dozou_katanuki/models"
)

// MotrixFullQueueStatus は Motrix の全ステータス別タスクリストです
type MotrixFullQueueStatus struct {
	ActiveTasks  []models.DownloaderTaskInfo `json:"active_tasks"`
	WaitingTasks []models.DownloaderTaskInfo `json:"waiting_tasks"`
	StoppedTasks []models.DownloaderTaskInfo `json:"stopped_tasks"`
}

// FetchMotrixFullQueue は Active / Waiting / Stopped の全タスクを取得します
func (a *App) FetchMotrixFullQueue() (MotrixFullQueueStatus, error) {
	var full MotrixFullQueueStatus
	keys := []string{"gid", "status", "totalLength", "completedLength", "downloadSpeed", "files", "errorMessage"}

	if raw, err := callMotrixRPC("aria2.tellActive", []interface{}{keys}); err == nil {
		full.ActiveTasks = parseAria2Tasks(raw)
	}
	if raw, err := callMotrixRPC("aria2.tellWaiting", []interface{}{0, 100, keys}); err == nil {
		full.WaitingTasks = parseAria2Tasks(raw)
	}
	if raw, err := callMotrixRPC("aria2.tellStopped", []interface{}{0, 100, keys}); err == nil {
		full.StoppedTasks = parseAria2Tasks(raw)
	}
	return full, nil
}

// AddMotrixDownload は新規ダウンロードタスクを投入します
func (a *App) AddMotrixDownload(urls []string, saveDir, fileName string) (string, error) {
	if len(urls) == 0 {
		return "", fmt.Errorf("urls required")
	}
	opts := map[string]interface{}{}
	if saveDir != "" {
		opts["dir"] = saveDir
	}
	if fileName != "" {
		opts["out"] = fileName
	}
	raw, err := callMotrixRPC("aria2.addUri", []interface{}{urls, opts})
	if err != nil {
		return "", err
	}
	var res struct {
		Result string `json:"result"`
	}
	_ = json.Unmarshal(raw, &res)
	return res.Result, nil
}
