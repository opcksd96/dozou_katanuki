// app/app_rpc_thunder_watcher.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var thunderWatchCounter int

// CheckThunderDirectoryStatus は D:\迅雷下载 の *.xltd および実体ファイルを走査してステータスを同期します
func (a *App) CheckThunderDirectoryStatus(tempDir string) {
	if tempDir == "" { tempDir = `D:\迅雷下载` }
	entries, err := os.ReadDir(tempDir)
	if err != nil || a.Repo == nil { return }

	xltdMediaIDs := make(map[string]bool)
	hasCompletedFiles := false

	for _, e := range entries {
		if e.IsDir() { continue }
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		switch ext {
		case ".xltd", ".td":
			baseName := strings.TrimSuffix(name, ext)
			xltdMediaIDs[resolveMediaIDFromFileName(baseName)] = true
		case ".jpg", ".jpeg", ".png", ".webp", ".mp4":
			hasCompletedFiles = true
		}
	}

	thunderWatchCounter++
	if len(xltdMediaIDs) > 0 || hasCompletedFiles || thunderWatchCounter%6 == 1 {
		a.AppendPipelineLog("THUNDER", "INFO",
			fmt.Sprintf("迅雷監視走査: %s (進行中xltd:%d, 完了候補:%v)", tempDir, len(xltdMediaIDs), hasCompletedFiles))
	}

	// 1. *.xltd が生えているメディアを ESCALATED に同期 (RETAINED からの再浮上復帰)
	for mediaID := range xltdMediaIDs {
		if m, err := a.Repo.GetMediaByID(mediaID); err == nil && m != nil {
			if m.DownloadStatus == "RETAINED" || m.DownloadStatus == "OUTSOURCED" {
				_ = a.Repo.UpdateMediaMetadata(mediaID, "ESCALATED", "", "", "迅雷 P2SP キャッシュ捕捉 (*.xltd 生起)")
			}
		}
	}

	// 2. 実体ファイルが完成していればアカウントフォルダへ自動移動 & COMPLETED 同期
	if hasCompletedFiles {
		synced, err := a.SyncThunderDownloads(tempDir)
		if err == nil && synced > 0 {
			a.AppendPipelineLog("THUNDER", "INFO",
				fmt.Sprintf("迅雷ダウンロード成果物 %d 件を取り込みました", synced))
		}
	}
}

// StartThunderBackgroundWatcher は定期的に CheckThunderDirectoryStatus を実行する監視ループです
func (a *App) StartThunderBackgroundWatcher() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			tempDir := `D:\迅雷下载`
			if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.ThunderDownloadDir != "" {
				tempDir = cfg.Storage.ThunderDownloadDir
			}
			a.CheckThunderDirectoryStatus(tempDir)
		}
	}()
}
