// app/app_rpc_media_trash.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TrashMedia は指定メディアを論理削除（ゴミ箱へ移動）します
func (a *App) TrashMedia(mediaID, reason string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if reason == "" {
		reason = "手動整理"
	}
	if err := a.Repo.TrashMedia(mediaID, reason, "admin_ui"); err != nil {
		log.Printf("[Wails RPC] TrashMedia error (%s): %v", mediaID, err)
		return err
	}
	if a.Ctx != nil {
		runtime.EventsEmit(a.Ctx, "media:trashed", map[string]string{
			"media_id": mediaID,
			"reason":   reason,
		})
	}
	return nil
}

// RestoreMedia はゴミ箱のメディアを通常状態へ復元します
func (a *App) RestoreMedia(mediaID string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.RestoreMedia(mediaID); err != nil {
		log.Printf("[Wails RPC] RestoreMedia error (%s): %v", mediaID, err)
		return err
	}
	if a.Ctx != nil {
		runtime.EventsEmit(a.Ctx, "media:restored", map[string]string{
			"media_id": mediaID,
		})
	}
	return nil
}

// BatchTrashMedia は複数メディアを一括で論理削除します
func (a *App) BatchTrashMedia(mediaIDs []string, reason string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if reason == "" {
		reason = "一括整理"
	}
	if err := a.Repo.BatchTrashMedia(mediaIDs, reason, "admin_ui"); err != nil {
		log.Printf("[Wails RPC] BatchTrashMedia error: %v", err)
		return err
	}
	if a.Ctx != nil {
		runtime.EventsEmit(a.Ctx, "media:batch_trashed", map[string]interface{}{
			"count":  len(mediaIDs),
			"reason": reason,
		})
	}
	return nil
}

// BatchRestoreMedia は複数メディアを一括復元します
func (a *App) BatchRestoreMedia(mediaIDs []string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.BatchRestoreMedia(mediaIDs); err != nil {
		log.Printf("[Wails RPC] BatchRestoreMedia error: %v", err)
		return err
	}
	if a.Ctx != nil {
		runtime.EventsEmit(a.Ctx, "media:batch_restored", map[string]interface{}{
			"count": len(mediaIDs),
		})
	}
	return nil
}
