package account

import (
	"context"
	"dozou_katanuki/domain/account"
)

type GetAccountUseCase struct {
	repo account.AccountRepository
}

func NewGetAccountUseCase(repo account.AccountRepository) *GetAccountUseCase {
	return &GetAccountUseCase{repo: repo}
}

func (uc *GetAccountUseCase) Execute(ctx context.Context, numericID string) (*account.Account, error) {
	return uc.repo.GetByID(ctx, numericID)
}

func (uc *GetAccountUseCase) ListAll(ctx context.Context, includeTrash bool) ([]*account.Account, error) {
	return uc.repo.ListAll(ctx, includeTrash)
}
