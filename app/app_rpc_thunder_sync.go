// app/app_rpc_thunder_sync.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
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
	if tempDir == "" {
		tempDir = `D:\迅雷下载`
	}

	destRoot, syncedCount := a.getMediaDownloadDir(), 0
	if c, err := a.processDirectoryFiles(tempDir, destRoot); err == nil {
		syncedCount += c
	}
	if c, err := a.processDirectoryFiles(destRoot, destRoot); err == nil {
		syncedCount += c
	}
	escalateDir := filepath.Join(destRoot, "_escalate")
	if c, err := a.processDirectoryFiles(escalateDir, destRoot); err == nil {
		syncedCount += c
	}

	if syncedCount > 0 {
		a.TriggerStashAllPipelines()
	}
	return syncedCount, nil
}

func (a *App) processDirectoryFiles(srcDir, destRoot string) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".mp4" && ext != ".webm" {
			continue
		}

		srcPath := filepath.Join(srcDir, name)
		mediaID := resolveMediaIDFromFileName(name)
		destFileName := mediaID
		if !strings.HasSuffix(strings.ToLower(destFileName), ext) {
			destFileName = destFileName + ext
		}

		targetSubDir := filepath.Join(destRoot, "_escalate")
		if a.Repo != nil {
			owner, err := a.Repo.GetMediaOwnerUsername(mediaID)
			if err != nil || owner == "" {
				owner, _ = a.Repo.GetMediaOwnerUsername(name)
			}
			if owner != "" {
				targetSubDir = filepath.Join(destRoot, owner, "X(Twitter)", "_assets")
			}
		}
		_ = os.MkdirAll(targetSubDir, 0755)
		destPath := filepath.Join(targetSubDir, destFileName)

		if srcPath == destPath {
			continue
		}

		if err := moveFileSafe(srcPath, destPath); err == nil {
			if a.Repo != nil {
				count++
				a.AppendPipelineLog("THUNDER", "SUCCESS", fmt.Sprintf("⚡ 迅雷完了回収: %s -> %s", name, destPath))
				go a.CoordinateTaskCompletion(mediaID, name, models.StageThunder)
				go a.ReapCompletedDuplicates(mediaID, name)
			}
		}
	}
	return count, nil
}

func resolveMediaIDFromFileName(name string) string {
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	noExt := strings.TrimSuffix(base, ext)

	// 重複連番 (1), (5) 等を除去
	if idx := strings.LastIndex(noExt, "("); idx != -1 && strings.HasSuffix(noExt, ")") {
		noExt = strings.TrimSpace(noExt[:idx])
	}
	// StreamSaver / 連番 _0, _0_0 等を除去
	for strings.HasSuffix(noExt, "_0") || strings.HasSuffix(noExt, "_1") {
		noExt = strings.TrimSuffix(strings.TrimSuffix(noExt, "_0"), "_1")
	}

	for _, sfx := range []string{
		"_wayback_plain", "_wayback_colon", "_wayback_orig",
		"_colon_orig", "_orig", "_large", "_plain",
		"_wayback", ":orig", ":large",
		"_motrix", "_requests", "_thunder",
	} {
		if strings.HasSuffix(noExt, sfx) {
			noExt = strings.TrimSuffix(noExt, sfx)
			break
		}
	}
	return noExt
}
