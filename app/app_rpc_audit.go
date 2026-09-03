// app/app_rpc_audit.go (100行以下)
package app

import (
	"log"
	"path/filepath"

	"dozou_katanuki/models"
)

func (a *App) getTrashDir() string {
	dumpsDir := "./backups/dumps"
	cfg, err := a.GetConfig()
	if err == nil && cfg != nil && cfg.Storage.DumpsDir != "" {
		dumpsDir = cfg.Storage.DumpsDir
	}
	return filepath.Join(dumpsDir, "_trash")
}

func (a *App) TriggerBackup() (string, error) {
	if err := a.WaitForReady(); err != nil {
		return "", err
	}
	if a.Scheduler == nil {
		return "", nil
	}
	return a.Scheduler.TriggerBackup()
}

func (a *App) TriggerPoll() (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	if a.Scheduler == nil {
		return nil, nil
	}
	return a.Scheduler.TriggerPoll()
}

func (a *App) RunAudit(purgeFiles, purgeDB bool) (*models.AuditReport, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	if a.AuditService == nil {
		return nil, nil
	}

	stashDir, blobsDir := "./stash", "./blobs"
	if cfg, err := a.GetConfig(); err == nil && cfg != nil {
		if cfg.Storage.StashDir != "" {
			stashDir = cfg.Storage.StashDir
		}
		if cfg.Storage.LocalMediaDir != "" {
			blobsDir = cfg.Storage.LocalMediaDir
		}
	}
	return a.AuditService.RunAudit(a.Ctx, stashDir, blobsDir, purgeFiles, purgeDB)
}

func (a *App) PurgeOrphanFiles(paths []string) (int, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	if a.AuditService == nil {
		return 0, nil
	}
	log.Printf("[Wails RPC] PurgeOrphanFiles: requested %d files to purge", len(paths))
	count, err := a.AuditService.PurgeOrphanFiles(paths)
	log.Printf("[Wails RPC] PurgeOrphanFiles: purged %d files (err: %v)", count, err)
	return count, err
}

func (a *App) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	if a.AuditService == nil {
		return 0, nil
	}
	log.Printf("[Wails RPC] PurgeOrphanDBMedia: requested %d media IDs to delete", len(mediaIDs))
	count, err := a.AuditService.PurgeOrphanDBMedia(a.getTrashDir(), mediaIDs)
	log.Printf("[Wails RPC] PurgeOrphanDBMedia: deleted %d records (err: %v)", count, err)
	return count, err
}

// RollbackLastPurge は直前に削除された孤立DBレコードを復元します
func (a *App) RollbackLastPurge() (int, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	if a.AuditService == nil {
		return 0, nil
	}
	log.Printf("[Wails RPC] RollbackLastPurge: starting rollback from trash...")
	count, err := a.AuditService.RollbackLastDBPurge(a.getTrashDir())
	log.Printf("[Wails RPC] RollbackLastPurge: restored %d records (err: %v)", count, err)
	return count, err
}

// CanRollback はロールバック可能なスナップショットが存在するか判定します
func (a *App) CanRollback() (bool, error) {
	if err := a.WaitForReady(); err != nil {
		return false, err
	}
	if a.AuditService == nil {
		return false, nil
	}
	return a.AuditService.CanRollbackDBPurge(a.getTrashDir()), nil
}

func (a *App) TriggerRestore(dumpsDir string, resetDB bool) (*models.JobProgress, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	if dumpsDir == "" {
		dumpsDir = "./backups/dumps"
		if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Storage.DumpsDir != "" {
			dumpsDir = cfg.Storage.DumpsDir
		}
	}
	if resetDB && a.Repo != nil {
		log.Println("[Wails RPC] TriggerRestore: Resetting database before restore...")
		_ = a.Repo.ResetDatabase()
	}
	return a.JobOrchestrator.EnqueueRestore(dumpsDir)
}
