// app/app_rpc_stash_sync_worker.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"dozou_katanuki/models"
)

var (
	stashWorkerMu   sync.Mutex
	isStashSyncBusy bool
)

// ScanUnsyncedMediaAndTriggerStash はダウンローダー側から分離された Stash 同期専用の非同期ワーカです
func (a *App) ScanUnsyncedMediaAndTriggerStash() {
	stashWorkerMu.Lock()
	if isStashSyncBusy {
		stashWorkerMu.Unlock()
		return
	}
	isStashSyncBusy = true
	stashWorkerMu.Unlock()

	defer func() {
		stashWorkerMu.Lock()
		isStashSyncBusy = false
		stashWorkerMu.Unlock()
	}()

	if a.Repo == nil || a.Repo.DB() == nil {
		return
	}

	var unsynced []models.Media
	// COMPLETED かつ、Stashに未登録のものを検索
	err := a.Repo.DB().Where("download_status = 'COMPLETED' AND (is_trash = 0 OR is_trash IS NULL) AND (stash_scene_id IS NULL OR stash_scene_id = '') AND (stash_image_id IS NULL OR stash_image_id = '')").
		Limit(50).Find(&unsynced).Error
	if err != nil || len(unsynced) == 0 {
		return
	}

	destRoot := a.getMediaDownloadDir()
	pathSet := make(map[string]bool)

	for _, m := range unsynced {
		owner := "unknown"
		if o, _ := a.Repo.GetMediaOwnerUsername(m.MediaID); o != "" {
			owner = o
		}
		targetDir := filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
		pathSet[targetDir] = true
	}

	var pathsToSync []string
	for p := range pathSet {
		pathsToSync = append(pathsToSync, p)
	}

	a.AppendPipelineLog("STASH_WORKER", "INFO", fmt.Sprintf("未同期メディア %d 件を発見。%d 個のディレクトリのStash同期を開始します", len(unsynced), len(pathsToSync)))
	
	// 非同期パイプラインを開始
	a.TriggerStashPipelineForPaths(pathsToSync)
	
	// パイプラインが完了するまで少し待機 (TriggerStashPipelineForPaths 内のSleepを考慮)
	time.Sleep(12 * time.Second)
}
