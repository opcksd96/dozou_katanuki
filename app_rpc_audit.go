// app_rpc_audit.go (100行以下)
package main

import (
	"log"

	"dozou_katanuki/models"
)

// TriggerBackup は即時オンラインバックアップ (VACUUM INTO) を実行する Wails バインドメソッドです
func (a *App) TriggerBackup() (string, error) {
	if err := a.waitForReady(); err != nil {
		return "", err
	}
	if a.scheduler == nil {
		return "", nil
	}
	return a.scheduler.TriggerBackup()
}

// TriggerPoll は即時第3段階ポーリング (Motrix ➔ Stash 自動回収) を実行する Wails バインドメソッドです
func (a *App) TriggerPoll() (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	if a.scheduler == nil {
		return nil, nil
	}
	return a.scheduler.TriggerPoll()
}

// RunAudit は SQLite3 整合性監査 (PRAGMA) および孤立ファイルの検出・パージを実行する Wails バインドメソッドです
func (a *App) RunAudit(purgeFiles, purgeDB bool) (*models.AuditReport, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	if a.auditService == nil {
		return nil, nil
	}

	stashDir := "./stash"
	blobsDir := "./blobs"
	cfg, err := a.GetConfig()
	if err == nil && cfg != nil {
		if cfg.Storage.StashDir != "" {
			stashDir = cfg.Storage.StashDir
		}
		if cfg.Storage.LocalMediaDir != "" {
			blobsDir = cfg.Storage.LocalMediaDir
		}
	}
	return a.auditService.RunAudit(a.ctx, stashDir, blobsDir, purgeFiles, purgeDB)
}

// PurgeOrphanFiles は指定された孤立ファイルをOSのごみ箱へ退避する Wails バインドメソッドです
func (a *App) PurgeOrphanFiles(paths []string) (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	if a.auditService == nil {
		return 0, nil
	}
	return a.auditService.PurgeOrphanFiles(paths)
}

// PurgeOrphanDBMedia は指定された media_id の DB レコードを削除する Wails バインドメソッドです
func (a *App) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	if a.auditService == nil {
		return 0, nil
	}
	return a.auditService.PurgeOrphanDBMedia(mediaIDs)
}

// TriggerRestore は Layer 2 (dumps/) からの完全オフライン自動リストア (SPEC-RECOVERY-001) をキックする Wails バインドメソッドです
func (a *App) TriggerRestore(dumpsDir string, resetDB bool) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	if dumpsDir == "" {
		dumpsDir = "./backups/dumps"
		cfg, err := a.GetConfig()
		if err == nil && cfg != nil && cfg.Storage.DumpsDir != "" {
			dumpsDir = cfg.Storage.DumpsDir
		}
	}

	if resetDB && a.repo != nil {
		log.Println("[Wails RPC] TriggerRestore: Resetting database before restore...")
		if err := a.repo.ResetDatabase(); err != nil {
			log.Printf("[Wails RPC] TriggerRestore: Database reset warning: %v", err)
		}
	}
	return a.jobOrchestrator.EnqueueRestore(dumpsDir)
}
