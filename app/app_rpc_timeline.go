// app/app_rpc_timeline.go (100行以下)
package app

import (
	"log"
	
	"dozou_katanuki/adapters/driving/dto"
	"dozou_katanuki/models"
)

// GetAccounts は登録されているホワイトリスト対象アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]*dto.AccountDTO, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	dtos, err := a.AccountUseCase.ListAllAccounts(a.Ctx)
	if err != nil {
		return nil, err
	}
	var res []*dto.AccountDTO
	for _, d := range dtos {
		if d.IsWhitelist {
			res = append(res, d)
		}
	}
	return res, nil
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
func (a *App) GetArticleDetail(platform, id string) (*models.ArticleDetailResult, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	res, err := a.TimelineService.GetArticleDetail(platform, id)
	log.Printf("[Wails RPC] GetArticleDetail(platform=%s, id=%s) -> res: %v (err: %v)", platform, id, res != nil, err)
	return res, err
}
