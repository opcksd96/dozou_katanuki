// domain/ports/repository.go (100行以下 - SPEC-PRINCIPLE-001)
package ports

import (
	"context"
	"dozou_katanuki/domain/entities"
)

// AccountRepository defines operations for Account entity
type AccountRepository interface {
	GetByNumericID(ctx context.Context, id string) (*entities.Account, error)
	GetByHandle(ctx context.Context, handle string) (*entities.Account, error)
	ListAll(ctx context.Context, includeTrash bool) ([]*entities.Account, error)
	Save(ctx context.Context, account *entities.Account) error
}

// ArticleRepository defines operations for Article entity
type ArticleRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Article, error)
	GetThread(ctx context.Context, id string) ([]*entities.Article, error)
	Save(ctx context.Context, article *entities.Article) error
}

// MediaRepository defines operations for Media entity
type MediaRepository interface {
	GetByID(ctx context.Context, mediaID string) (*entities.Media, error)
	GetByArticleID(ctx context.Context, articleID string) ([]*entities.Media, error)
	ListByStatus(ctx context.Context, status string) ([]*entities.Media, error)
	Save(ctx context.Context, media *entities.Media) error
}
