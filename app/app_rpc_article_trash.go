// app/app_rpc_article_trash.go (Under 100 lines - SPEC-PRINCIPLE-001)
package app

import (
	"log"
)

// TrashArticle は指定された記事を論理削除（ゴミ箱へ移動）する Wails バインドメソッドです
func (a *App) TrashArticle(id, trashedBy, reason string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.TrashArticle(id, trashedBy, reason); err != nil {
		log.Printf("[Wails RPC] TrashArticle error (id=%s): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] Article trashed (id=%s, by=%s, reason=%s)", id, trashedBy, reason)
	if a.Ctx != nil {
		a.EmitEvent("article:trashed", map[string]string{
			"id": id, "trashed_by": trashedBy, "reason": reason,
		})
	}
	return nil
}

// RestoreArticle はゴミ箱内の指定記事を復元する Wails バインドメソッドです
func (a *App) RestoreArticle(id string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.RestoreArticle(id); err != nil {
		log.Printf("[Wails RPC] RestoreArticle error (id=%s): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] Article restored (id=%s)", id)
	if a.Ctx != nil {
		a.EmitEvent("article:restored", map[string]string{"id": id})
	}
	return nil
}
