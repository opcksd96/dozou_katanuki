// app_rpc_timeline.go (100行以下)
package main

import (
	"log"

	"dozou_katanuki/models"
)

// GetAccounts は登録されている全アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]models.RenderAuthor, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.timelineService.GetAccounts(platform)
}

// GetTimeline はフロントエンドへ RenderTree 配列を供給する Wails バインドメソッドです
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	res, err := a.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	log.Printf("[Wails RPC] GetTimeline(platform=%s, accountID=%s, filter=%s) -> 取得件数: %d (err: %v)",
		platform, accountID, filter, len(res), err)
	return res, err
}

// GetArticleDetail は指定された個別記事およびスレッド会話ツリーを取得する Wails バインドメソッドです
func (a *App) GetArticleDetail(platform, id string) (*models.ArticleDetailResult, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.timelineService.GetArticleDetail(platform, id)
}
