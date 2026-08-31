// app/app_rpc_thunder_sync.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// SyncThunderDownloads は迅雷フォルダおよび _escalate の完了ファイルをアカウント所定フォルダへ移動・DB同期・Stash連携します
func (a *App) SyncThunderDownloads(customDir string) (int, error) {
	tempDir := customDir
	if tempDir == "" {
		if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.ThunderDownloadDir != "" {
			tempDir = cfg.Storage.ThunderDownloadDir
		}
	}
	if tempDir == "" { tempDir = `D:\迅雷下载` }

	destRoot := a.getMediaDownloadDir()
	syncedCount := 0

	if count, err := a.processDirectoryFiles(tempDir, destRoot); err == nil { syncedCount += count }
	if count, err := a.processDirectoryFiles(destRoot, destRoot); err == nil { syncedCount += count }
	escalateDir := filepath.Join(destRoot, "_escalate")
	if count, err := a.processDirectoryFiles(escalateDir, destRoot); err == nil { syncedCount += count }

	if syncedCount > 0 {
		a.TriggerStashFullPipeline()
	}

	return syncedCount, nil
}

func (a *App) processDirectoryFiles(srcDir, destRoot string) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil { return 0, err }

	count := 0
	for _, entry := range entries {
		if entry.IsDir() { continue }
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".mp4" && ext != ".webm" { continue }

		srcPath := filepath.Join(srcDir, name)
		mediaID := resolveMediaIDFromFileName(name)
		destFileName := mediaID
		if !strings.HasSuffix(strings.ToLower(destFileName), ext) { destFileName = destFileName + ext }

		targetSubDir := filepath.Join(destRoot, "_escalate")
		if a.Repo != nil {
			if owner, err := a.Repo.GetMediaOwnerUsername(mediaID); err == nil && owner != "" {
				targetSubDir = filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
			}
		}
		_ = os.MkdirAll(targetSubDir, 0755)
		destPath := filepath.Join(targetSubDir, destFileName)

		if srcPath == destPath { continue }

		if err := moveFileSafe(srcPath, destPath); err == nil {
			if a.Repo != nil {
				count++
				// 【新SSOT】全ステージ協調刈り取り ＆ Stash自動パイプライン連携
				go a.CoordinateTaskCompletion(mediaID, name, models.StageThunder)
			}
		}
	}
	return count, nil
}

func resolveMediaIDFromFileName(name string) string {
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	noExt := strings.TrimSuffix(base, ext)

	for _, sfx := range []string{"_orig", "_large", "_wayback_orig", "_wayback", ":orig", ":large"} {
		if strings.HasSuffix(noExt, sfx) {
			noExt = strings.TrimSuffix(noExt, sfx)
			break
		}
	}
	return noExt
}
