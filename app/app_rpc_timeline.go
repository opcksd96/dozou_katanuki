// app/app_rpc_timeline.go (100行以下)
package app

import (
	"log"
	
	"dozou_katanuki/adapters/driving/dto"
	"dozou_katanuki/models"
)

// GetAccounts は登録されている全アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]*dto.AccountDTO, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	// Use Application Layer (UseCase)
	return a.AccountUseCase.ListAllAccounts(a.Ctx)
}

// GetTimeline はフロントエンドへ RenderTree 配列を供給する Wails バインドメソッドです
// FIXME: TimelineUseCase に本格的なクエリ機能を実装するまで、一時的に古いServiceを経由
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	res, err := a.TimelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	log.Printf("[Wails RPC] GetTimeline(platform=%s, accountID=%s, filter=%s) -> 取得件数: %d (err: %v)",
		platform, accountID, filter, len(res), err)
	return res, err
}

// GetArticleDetail は指定された個別記事およびスレッド会話ツリーを取得する Wails バインドメソッドです
func (a *App) GetArticleDetail(platform, id string) ([]dto.RenderTree, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	// Use Application Layer (UseCase)
	return a.TimelineUseCase.GetThread(a.Ctx, id)
}
