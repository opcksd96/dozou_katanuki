// app_rpc_audit.go (100行以下)
package main

import (
	"log"
	"path/filepath"

	"dozou_katanuki/models"
)

func (a *App) getTrashDir() string {
	dumpsDir := "./backups/dumps"
	cfg, err := a.GetConfig()
	if err == nil && cfg != nil && cfg.Storage.DumpsDir != "" { dumpsDir = cfg.Storage.DumpsDir }
	return filepath.Join(dumpsDir, "_trash")
}

func (a *App) TriggerBackup() (string, error) {
	if err := a.waitForReady(); err != nil { return "", err }
	if a.scheduler == nil { return "", nil }
	return a.scheduler.TriggerBackup()
}

func (a *App) TriggerPoll() (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	if a.scheduler == nil { return nil, nil }
	return a.scheduler.TriggerPoll()
}

func (a *App) RunAudit(purgeFiles, purgeDB bool) (*models.AuditReport, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	if a.auditService == nil { return nil, nil }

	stashDir, blobsDir := "./stash", "./blobs"
	if cfg, err := a.GetConfig(); err == nil && cfg != nil {
		if cfg.Storage.StashDir != "" { stashDir = cfg.Storage.StashDir }
		if cfg.Storage.LocalMediaDir != "" { blobsDir = cfg.Storage.LocalMediaDir }
	}
	return a.auditService.RunAudit(a.ctx, stashDir, blobsDir, purgeFiles, purgeDB)
}

func (a *App) PurgeOrphanFiles(paths []string) (int, error) {
	if err := a.waitForReady(); err != nil { return 0, err }
	if a.auditService == nil { return 0, nil }
	log.Printf("[Wails RPC] PurgeOrphanFiles: requested %d files to purge", len(paths))
	count, err := a.auditService.PurgeOrphanFiles(paths)
	log.Printf("[Wails RPC] PurgeOrphanFiles: purged %d files (err: %v)", count, err)
	return count, err
}

func (a *App) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	if err := a.waitForReady(); err != nil { return 0, err }
	if a.auditService == nil { return 0, nil }
	log.Printf("[Wails RPC] PurgeOrphanDBMedia: requested %d media IDs to delete", len(mediaIDs))
	count, err := a.auditService.PurgeOrphanDBMedia(a.getTrashDir(), mediaIDs)
	log.Printf("[Wails RPC] PurgeOrphanDBMedia: deleted %d records (err: %v)", count, err)
	return count, err
}

// RollbackLastPurge は直前に削除された孤立DBレコードを復元します
func (a *App) RollbackLastPurge() (int, error) {
	if err := a.waitForReady(); err != nil { return 0, err }
	if a.auditService == nil { return 0, nil }
	log.Printf("[Wails RPC] RollbackLastPurge: starting rollback from trash...")
	count, err := a.auditService.RollbackLastDBPurge(a.getTrashDir())
	log.Printf("[Wails RPC] RollbackLastPurge: restored %d records (err: %v)", count, err)
	return count, err
}

// CanRollback はロールバック可能なスナップショットが存在するか判定します
func (a *App) CanRollback() (bool, error) {
	if err := a.waitForReady(); err != nil { return false, err }
	if a.auditService == nil { return false, nil }
	return a.auditService.CanRollbackDBPurge(a.getTrashDir()), nil
}

func (a *App) TriggerRestore(dumpsDir string, resetDB bool) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	if dumpsDir == "" {
		dumpsDir = "./backups/dumps"
		if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.DumpsDir != "" { dumpsDir = cfg.Storage.DumpsDir }
	}
	if resetDB && a.repo != nil {
		log.Println("[Wails RPC] TriggerRestore: Resetting database before restore...")
		_ = a.repo.ResetDatabase()
	}
	return a.jobOrchestrator.EnqueueRestore(dumpsDir)
}
