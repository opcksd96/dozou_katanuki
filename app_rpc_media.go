// app_rpc_media.go (100行以下 - SPEC-PRINCIPLE-001)
package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// OpenInExplorer はメディアの実体ファイルを Windows エクスプローラーで選択状態で開きます
func (a *App) OpenInExplorer(mediaID string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	filePath, err := a.repo.ResolveMediaFilePath(mediaID)
	if err != nil {
		return err
	}
	absPath, _ := filepath.Abs(filePath)
	return exec.Command("explorer", fmt.Sprintf("/select,%s", absPath)).Start()
}

// OpenWithDefaultApp はメディアの実体ファイルを OS 既定のビューアまたはプレイヤーで開きます
func (a *App) OpenWithDefaultApp(mediaID string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	filePath, err := a.repo.ResolveMediaFilePath(mediaID)
	if err != nil {
		return err
	}
	absPath, _ := filepath.Abs(filePath)
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", absPath).Start()
}

// ToggleMediaBookmark はメディアのブックマーク（お気に入り）状態を切り替えます
func (a *App) ToggleMediaBookmark(mediaID string) (bool, error) {
	if err := a.waitForReady(); err != nil {
		return false, err
	}
	return a.repo.ToggleMediaBookmark(mediaID)
}

// RenameMedia はメディアIDおよびファイル名を変更します
func (a *App) RenameMedia(mediaID, newName string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	return a.repo.RenameMedia(mediaID, newName)
}

// RequeueMediaByStatus は指定ステータスのメディアを一括で QUEUED 状態に戻します
func (a *App) RequeueMediaByStatus(status, accountID string) (int64, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	return a.repo.RequeueMediaByStatus(status, accountID)
}

// MergeDuplicateMedia は同一ファイル名や重複URLのメディアを自動統合します (Intelligence)
func (a *App) MergeDuplicateMedia() (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	return a.repo.MergeDuplicateMedia()
}

// PurgeLowerResolutionDuplicates は同一コンテンツの低解像度バリエーションをゴミ箱へ退避します (Intelligence)
func (a *App) PurgeLowerResolutionDuplicates() (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	return a.repo.PurgeLowerResolutionDuplicates("backups/trash")
}
