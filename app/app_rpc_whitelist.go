// app/app_rpc_whitelist.go (100行以下)
package app

import (
	"log"

	"dozou_katanuki/models"
)

// GetWhitelists は全ホワイトリスト項目を取得する Wails バインドメソッドです
func (a *App) GetWhitelists() ([]models.Whitelist, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.Repo.GetWhitelists()
}

// AddWhitelist はホワイトリスト項目を追加する Wails バインドメソッドです
func (a *App) AddWhitelist(itemType, value string) (*models.Whitelist, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	item, err := a.Repo.AddWhitelist(itemType, value)
	if err != nil {
		log.Printf("[Wails RPC] AddWhitelist error: %v", err)
		return nil, err
	}
	log.Printf("[Wails RPC] AddWhitelist added: [%s] %s (id: %d)", item.Type, item.Value, item.ID)
	return item, nil
}

// UpdateWhitelist はホワイトリスト項目を更新する Wails バインドメソッドです
func (a *App) UpdateWhitelist(id uint, itemType, value string, isActive bool) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	err := a.Repo.UpdateWhitelist(id, itemType, value, isActive)
	if err != nil {
		log.Printf("[Wails RPC] UpdateWhitelist error (id: %d): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] UpdateWhitelist updated: (id: %d) [%s] %s (active: %v)", id, itemType, value, isActive)
	return nil
}

// DeleteWhitelist はホワイトリスト項目を削除する Wails バインドメソッドです
func (a *App) DeleteWhitelist(id uint) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	err := a.Repo.DeleteWhitelist(id)
	if err != nil {
		log.Printf("[Wails RPC] DeleteWhitelist error (id: %d): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] DeleteWhitelist deleted: (id: %d)", id)
	return nil
}

// ToggleWhitelist はホワイトリスト項目の有効/無効を切り替える Wails バインドメソッドです
func (a *App) ToggleWhitelist(id uint) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	err := a.Repo.ToggleWhitelist(id)
	if err != nil {
		log.Printf("[Wails RPC] ToggleWhitelist error (id: %d): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] ToggleWhitelist toggled: (id: %d)", id)
	return nil
}
