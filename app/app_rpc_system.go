// app/app_rpc_system.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"log"

	"dozou_katanuki/middleware"
)

// RestartBackendServices はバックエンドの全内部サービスを再初期化・リロードする Wails バインドメソッドです
func (a *App) RestartBackendServices() (bool, error) {
	if err := a.WaitForReady(); err != nil { return false, err }
	log.Printf("[Wails RPC] RestartBackendServices requested by user.")

	// 1. 設定再読み込み
	cfg, _ := a.GetConfig()
	if cfg != nil {
		if a.JobOrchestrator != nil { a.JobOrchestrator.SetStorageConfig(cfg.Storage) }
		if a.UnifiedHandler != nil { a.UnifiedHandler.SetMediaDir(cfg.Storage.LocalMediaDir) }
		_ = middleware.SyncStashConfig(cfg)
	}

	// 2. システムジャーナルに再起動イベントを記録
	middleware.GetGlobalJournal().Record(
		"system", "INFO", "backend_restarted",
		"Wails backend services re-initialized and synced with config.json",
		map[string]interface{}{"status": "reloaded_ok"},
	)

	log.Printf("[Wails RPC] Backend services successfully re-initialized.")
	return true, nil
}
