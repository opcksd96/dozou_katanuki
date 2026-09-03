// app/app_rpc_thunder_watcher.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CheckThunderDirectoryStatus は D:\迅雷下载 の *.xltd および実体ファイルを走査してステータスを同期します
func (a *App) CheckThunderDirectoryStatus(tempDir string) {
	if tempDir == "" {
		tempDir = `D:\迅雷下载`
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil || a.Repo == nil {
		return
	}

	xltdMediaIDs := make(map[string]bool)
	hasCompletedFiles := false

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// *.xltd や *.td (一時ダウンロードファイル) を検知
		if ext == ".xltd" || ext == ".td" {
			baseName := strings.TrimSuffix(name, ext)
			mediaID := resolveMediaIDFromFileName(baseName)
			xltdMediaIDs[mediaID] = true
		} else if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".mp4" {
			hasCompletedFiles = true
		}
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
		_, _ = a.SyncThunderDownloads(tempDir)
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
