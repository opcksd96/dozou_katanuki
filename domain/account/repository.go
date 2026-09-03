package account

import "context"

// AccountRepository defines the outbound port for Account persistence
type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, numericID string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
	ListAll(ctx context.Context, includeTrash bool) ([]*Account, error)
	Delete(ctx context.Context, numericID string) error
}
