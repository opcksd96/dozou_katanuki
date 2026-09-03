package account

import (
	"context"
	"dozou_katanuki/domain"
	"dozou_katanuki/domain/account"
)

type MergeAccountUseCase struct {
	repo account.AccountRepository
}

func NewMergeAccountUseCase(repo account.AccountRepository) *MergeAccountUseCase {
	return &MergeAccountUseCase{repo: repo}
}

func (uc *MergeAccountUseCase) Execute(ctx context.Context, newAccountData *account.Account) (*account.Account, error) {
	// Guard: This is a management operation, ensure caller has Admin scope.
	if err := domain.EnsureAdmin(ctx); err != nil {
		return nil, err
	}

	existing, err := uc.repo.GetByID(ctx, newAccountData.NumericID)
	if err != nil {
		// Assume not found, save as new
		if errSave := uc.repo.Save(ctx, newAccountData); errSave != nil {
			return nil, errSave
		}
		return newAccountData, nil
	}

	// Merge updates into existing entity
	existing.MergeUpdates(newAccountData)

	if errSave := uc.repo.Save(ctx, existing); errSave != nil {
		return nil, errSave
	}

	return existing, nil
}
