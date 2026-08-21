// app_rpc_article.go (100行以下)
package main

import (
	"log"

	"dozou_katanuki/middleware"
	"dozou_katanuki/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SearchArticles は保存済み記事の検索および一覧取得を行う Wails バインドメソッドです
func (a *App) SearchArticles(query, accountID, filter string, limit, offset int) (*models.ArticleSearchResult, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	articles, total, err := a.repo.SearchArticles(query, accountID, filter, limit, offset)
	if err != nil {
		log.Printf("[Wails RPC] SearchArticles error: %v", err)
		return nil, err
	}

	items := make([]models.RenderTree, len(articles))
	for i, art := range articles {
		items[i] = middleware.ToRenderTree(art, "twitter")
	}
	return &models.ArticleSearchResult{Items: items, Total: total}, nil
}

// GetArticle は指定されたIDの記事詳細を取得する Wails バインドメソッドです
func (a *App) GetArticle(id string) (*models.RenderTree, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	art, err := a.repo.GetArticleByID(id)
	if err != nil {
		log.Printf("[Wails RPC] GetArticle error (id=%s): %v", id, err)
		return nil, err
	}
	renderTree := middleware.ToRenderTree(*art, "twitter")
	return &renderTree, nil
}

// UpdateArticleTranslations は記事の日本語・英語・中国語翻訳テキストを更新する Wails バインドメソッドです
func (a *App) UpdateArticleTranslations(id, ja, en, zh string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	err := a.repo.UpdateArticleTranslations(id, ja, en, zh)
	if err != nil {
		log.Printf("[Wails RPC] UpdateArticleTranslations error (id=%s): %v", id, err)
		return err
	}
	return nil
}

// RetryMediaDownload は指定されたメディアのダウンロードステータスをリセットし再試行ジョブをキックする Wails バインドメソッドです
func (a *App) RetryMediaDownload(mediaID string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	if a.repo.ResetMediaStatus(mediaID) != nil {
		return a.repo.ResetMediaStatus(mediaID)
	}
	_, err := a.jobOrchestrator.EnqueueMediaDownload("twitter", mediaID)
	if err != nil {
		log.Printf("[Wails RPC] RetryMediaDownload EnqueueMediaDownload error: %v", err)
		return err
	}
	runtime.EventsEmit(a.ctx, "media:retried", map[string]string{"media_id": mediaID})
	return nil
}
