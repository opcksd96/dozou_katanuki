// app/app_rpc_article_batch.go (Under 100 lines - SPEC-PRINCIPLE-001)
package app

import (
	"log"

	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

// BatchTrashArticles は指定された複数記事を一括論理削除します
func (a *App) BatchTrashArticles(ids []string, trashedBy, reason string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.BatchTrashArticles(ids, trashedBy, reason); err != nil {
		log.Printf("[Wails RPC] BatchTrashArticles error: %v", err)
		return err
	}
	if a.Ctx != nil {
		a.EmitEvent("article:batch_trashed", map[string]interface{}{
			"ids": ids, "trashed_by": trashedBy, "reason": reason,
		})
	}
	return nil
}

// BatchRestoreArticles は指定された複数記事を一括復元します
func (a *App) BatchRestoreArticles(ids []string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.BatchRestoreArticles(ids); err != nil {
		log.Printf("[Wails RPC] BatchRestoreArticles error: %v", err)
		return err
	}
	if a.Ctx != nil {
		a.EmitEvent("article:batch_restored", map[string]interface{}{"ids": ids})
	}
	return nil
}

// BatchResetTranslations は指定された複数記事の翻訳を一括初期化します
func (a *App) BatchResetTranslations(ids []string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.BatchResetTranslations(ids); err != nil {
		log.Printf("[Wails RPC] BatchResetTranslations error: %v", err)
		return err
	}
	if a.Ctx != nil {
		a.EmitEvent("article:batch_reset_translations", map[string]interface{}{"ids": ids})
	}
	return nil
}

// GetArticlesByIDs は指定複数IDの記事詳細（RenderTree配列）を返します
func (a *App) GetArticlesByIDs(ids []string) ([]models.RenderTree, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	articles, err := a.Repo.GetArticlesByIDs(ids)
	if err != nil {
		return nil, err
	}
	items := make([]models.RenderTree, len(articles))
	for i, art := range articles {
		items[i] = middleware.ToRenderTree(art, "twitter")
	}
	return items, nil
}
