// app/app_rpc_account_trash.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"dozou_katanuki/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TrashAccount moves the specified account to the trash bin (soft-delete)
func (a *App) TrashAccount(numericID, reason string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if reason == "" {
		reason = "ユーザーによる手動退避"
	}
	if err := a.Repo.TrashAccount(numericID, reason, "admin_ui"); err != nil {
		return err
	}
	if a.Ctx != nil {
		runtime.EventsEmit(a.Ctx, "account:trashed", map[string]string{
			"numeric_id": numericID,
			"reason":     reason,
		})
	}
	return nil
}

// RestoreAccount restores an account from the trash bin back to active status
func (a *App) RestoreAccount(numericID string) error {
	if err := a.WaitForReady(); err != nil {
		return err
	}
	if err := a.Repo.RestoreAccount(numericID); err != nil {
		return err
	}
	if a.Ctx != nil {
		runtime.EventsEmit(a.Ctx, "account:restored", map[string]string{
			"numeric_id": numericID,
		})
	}
	return nil
}

// ListTrashedAccounts returns all accounts currently in the trash bin
func (a *App) ListTrashedAccounts() ([]models.Account, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.Repo.GetTrashedAccounts()
}

// ListAccountsWithTrash returns accounts, optionally including trashed ones
func (a *App) ListAccountsWithTrash(includeTrash bool) ([]models.Account, error) {
	if err := a.WaitForReady(); err != nil {
		return nil, err
	}
	return a.Repo.GetAllAccounts(includeTrash)
}
