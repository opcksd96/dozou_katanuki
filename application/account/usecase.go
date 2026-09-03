// application/account/usecase.go (100行以下 - SPEC-PRINCIPLE-001)
package account

import (
	"context"
	"dozou_katanuki/adapters/driving/dto"
	"dozou_katanuki/domain/ports"
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
	"time"
)

type AccountUseCase interface {
	GetAccountOverview(ctx context.Context, handle string) (*dto.AccountDTO, error)
	ListAllAccounts(ctx context.Context) ([]*dto.AccountDTO, error)
}

type accountUseCaseImpl struct {
	repo ports.AccountRepository
}

func NewAccountUseCase(repo ports.AccountRepository) AccountUseCase {
	return &accountUseCaseImpl{repo: repo}
}

func (u *accountUseCaseImpl) GetAccountOverview(ctx context.Context, handle string) (*dto.AccountDTO, error) {
	acc, err := u.repo.GetByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}
	
	// Map Entity to DTO
	return &dto.AccountDTO{
		NumericID:   acc.NumericID,
		Username:    acc.Username,
		DisplayName: acc.DisplayName,
		AvatarURL:   middleware.ResolveAccountAvatar("twitter", time.Now(), models.Account{AvatarBase64: acc.AvatarBase64}),
		Description: acc.Description,
		GroupName:   acc.GroupName,
		IsWhitelist: acc.IsWhitelist,
	}, nil
}

func (u *accountUseCaseImpl) ListAllAccounts(ctx context.Context) ([]*dto.AccountDTO, error) {
	accs, err := u.repo.ListAll(ctx, false)
	if err != nil {
		return nil, err
	}
	
	dtos := make([]*dto.AccountDTO, len(accs))
	for i, acc := range accs {
		dtos[i] = &dto.AccountDTO{
			NumericID:   acc.NumericID,
			Username:    acc.Username,
			DisplayName: acc.DisplayName,
			AvatarURL:   middleware.ResolveAccountAvatar("twitter", time.Now(), models.Account{AvatarBase64: acc.AvatarBase64}),
			Description: acc.Description,
			GroupName:   acc.GroupName,
			IsWhitelist: acc.IsWhitelist,
		}
	}
	return dtos, nil
}
