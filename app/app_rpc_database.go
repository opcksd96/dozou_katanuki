// app/app_rpc_database.go (100行以下)
package app

import (
	"log"

	"dozou_katanuki/models"
)

// GetAccountDetail は指定されたアカウントの詳細とアバター世代履歴を取得する Wails バインドメソッドです
func (a *App) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	res, err := a.Repo.GetAccountDetail(numericID)
	if err != nil {
		log.Printf("[Wails RPC] GetAccountDetail error: %v", err)
		return nil, err
	}
	return res, nil
}

// GetMediaList はメディア一覧（Stashステータス・アカウント情報・種別フィルタ付き）を取得する Wails バインドメソッドです
func (a *App) GetMediaList(accountID, status, mediaType string, limit, offset int) (*models.MediaSearchResult, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	_, _ = a.Repo.MigrateExcludedMedia() // レガシーWhitelist外DEAD_404をEXCLUDEDへ安全移行
	res, err := a.TimelineService.SearchMediaDetails(accountID, status, mediaType, limit, offset)
	if err != nil {
		log.Printf("[Wails RPC] GetMediaList error: %v", err)
		return nil, err
	}
	return res, nil
}

// PurgeMedia は指定された単一メディアレコードをデータベースから物理削除する Wails バインドメソッドです
func (a *App) PurgeMedia(mediaID string) error {
	if err := a.WaitForReady(); err != nil { return err }
	return a.Repo.PurgeMedia(mediaID)
}

// PurgeMediaByStatus は指定ステータス（EXCLUDED, UNLINKED, DEAD_404等）のメディアを一括削除する Wails バインドメソッドです
func (a *App) PurgeMediaByStatus(status, accountID string) (int64, error) {
	if err := a.WaitForReady(); err != nil { return 0, err }
	return a.Repo.PurgeMediaByStatus(status, accountID)
}

// UpdateMediaMetadata は指定されたメディアのメタデータ（ステータス・Stash ID・失敗理由）を更新する Wails バインドメソッドです
func (a *App) UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason string) error {
	if err := a.WaitForReady(); err != nil { return err }
	return a.Repo.UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason)
}

// GetTableRecords は汎用テーブルの生カラム・ロウ形式スプレッドシートデータを取得する Wails バインドメソッドです
func (a *App) GetTableRecords(tableName string, limit, offset int, search string) (*models.TableRecordResult, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	res, err := a.Repo.GetTableRecords(tableName, limit, offset, search)
	if err != nil {
		log.Printf("[Wails RPC] GetTableRecords error: %v", err)
		return nil, err
	}
	return res, nil
}

// ListAllAccounts は登録済みアカウント一覧を取得する Wails バインドメソッドです
func (a *App) ListAllAccounts() ([]models.Account, error) {
	if err := a.WaitForReady(); err != nil { return nil, err }
	return a.Repo.ListAccounts()
}

// MergeAccounts は source アカウントの全記事を target へ移動し、source の alias_of に target を設定する Wails バインドメソッドです
func (a *App) MergeAccounts(sourceNumericID, targetNumericID string) error {
	if err := a.WaitForReady(); err != nil { return err }
	if err := a.Repo.MergeAccounts(sourceNumericID, targetNumericID); err != nil {
		log.Printf("[Wails RPC] MergeAccounts error: %v", err)
		return err
	}
	return nil
}
