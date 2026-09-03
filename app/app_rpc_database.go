// app/app_rpc_database.go (100行以下)
package app

import (
	"log"

	"dozou_katanuki/models"
)

// GetAccountDetail は指定されたアカウントの詳細とアバター世代履歴を取得する Wails バインドメソッドです
func (a *App) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	res, err := a.Repo.GetAccountDetail(numericID)
	if err != nil {
		log.Printf("[Wails RPC] GetAccountDetail error: %v", err)
		return nil, err
	}
	return res, nil
}

// ToggleAccountWhitelist は指定されたアカウントの Whitelist（巡回対象）フラグを切り替える Wails バインドメソッドです
func (a *App) ToggleAccountWhitelist(numericID string, isWhitelist bool) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.ToggleAccountWhitelist(numericID, isWhitelist); err != nil {
		log.Printf("[Wails RPC] ToggleAccountWhitelist error: %v", err)
		return err
	}
	return nil
}

// GetMediaList はメディア一覧（Stashステータス・アカウント情報・種別フィルタ付き）を取得する Wails バインドメソッドです
func (a *App) GetMediaList(accountID, status, mediaType string, limit, offset int) (*models.MediaSearchResult, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	_, _ = a.Repo.MigrateExcludedMedia()
	return a.TimelineService.SearchMediaDetails(accountID, status, mediaType, limit, offset)
}

// GetMediaDownloadStatusStats はダウンロードキュー全体のステータス別集計を取得する Wails バインドメソッドです
func (a *App) GetMediaDownloadStatusStats(accountID string) (*models.DownloadStatusStats, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.Repo.FetchDownloadStatusStats(accountID)
}

// PurgeMedia は指定された単一メディアレコードをデータベースから物理削除する Wails バインドメソッドです
func (a *App) PurgeMedia(mediaID string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	return a.Repo.PurgeMedia(mediaID)
}

// PurgeMediaByStatus は指定ステータスのメディアを一括削除する Wails バインドメソッドです
func (a *App) PurgeMediaByStatus(status, accountID string) (int64, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	return a.Repo.PurgeMediaByStatus(status, accountID)
}

// UpdateMediaMetadata は指定されたメディアのメタデータを更新する Wails バインドメソッドです
func (a *App) UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if downloadStatus == "QUEUED" {
		failedReason = ""
	}
	if err := a.Repo.UpdateMediaMetadata(mediaID, downloadStatus, stashSceneID, stashImageID, failedReason); err != nil {
		return err
	}
	if downloadStatus == "QUEUED" {
		_, _ = a.ResetSpecificMediasToQueued([]string{mediaID})
	} else if a.Repo != nil && a.Repo.DB() != nil {
		_ = a.Repo.DB().Where("media_id = ?", mediaID).Delete(&models.DownloadTask{}).Error
	}
	return nil
}

// ListAllAccounts は登録済みアカウント一覧を取得する Wails バインドメソッドです
func (a *App) ListAllAccounts() ([]models.Account, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.Repo.ListAccounts()
}

// MergeAccounts は source アカウントの全記事を target へ移動する Wails バインドメソッドです
func (a *App) MergeAccounts(sourceNumericID, targetNumericID string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	return a.Repo.MergeAccounts(sourceNumericID, targetNumericID)
}
