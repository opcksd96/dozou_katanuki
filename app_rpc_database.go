// app_rpc_database.go (100行以下)
package main

import (
	"log"

	"dozou_katanuki/models"
)

// GetAccountDetail は指定されたアカウントの詳細とアバター世代履歴を取得する Wails バインドメソッドです
func (a *App) GetAccountDetail(numericID string) (*models.AccountDetailResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	res, err := a.repo.GetAccountDetail(numericID)
	if err != nil {
		log.Printf("[Wails RPC] GetAccountDetail error: %v", err)
		return nil, err
	}
	return res, nil
}

// GetMediaList はメディア一覧（Stashステータス・アカウント情報付き）を取得する Wails バインドメソッドです
func (a *App) GetMediaList(accountID, status string, limit, offset int) (*models.MediaSearchResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	res, err := a.repo.SearchMediaDetails(accountID, status, limit, offset)
	if err != nil {
		log.Printf("[Wails RPC] GetMediaList error: %v", err)
		return nil, err
	}
	return res, nil
}

// GetTableRecords は汎用テーブルの生カラム・ロウ形式スプレッドシートデータを取得する Wails バインドメソッドです
func (a *App) GetTableRecords(tableName string, limit, offset int, search string) (*models.TableRecordResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	res, err := a.repo.GetTableRecords(tableName, limit, offset, search)
	if err != nil {
		log.Printf("[Wails RPC] GetTableRecords error: %v", err)
		return nil, err
	}
	return res, nil
}

// ListAllAccounts は登録済みアカウント一覧を取得する Wails バインドメソッドです
func (a *App) ListAllAccounts() ([]models.Account, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	return a.repo.ListAccounts()
}

