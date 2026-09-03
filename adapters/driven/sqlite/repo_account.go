package sqlite

import (
	"context"
	"dozou_katanuki/domain/entities"
	"dozou_katanuki/domain/ports"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

// Compile-time assertion that AccountRepositoryImpl implements ports.AccountRepository
var _ ports.AccountRepository = (*AccountRepositoryImpl)(nil)

func NewAccountRepositoryImpl(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db: db}
}

func (r *AccountRepositoryImpl) Save(ctx context.Context, acc *entities.Account) error {
	schema := &AccountSchema{}
	schema.FromEntity(acc)
	return r.db.WithContext(ctx).Save(schema).Error
}

func (r *AccountRepositoryImpl) GetByNumericID(ctx context.Context, numericID string) (*entities.Account, error) {
	var schema AccountSchema
	if err := r.db.WithContext(ctx).Where("numeric_id = ?", numericID).First(&schema).Error; err != nil {
		return nil, err
	}
	return schema.ToEntity(), nil
}

func (r *AccountRepositoryImpl) GetByHandle(ctx context.Context, handle string) (*entities.Account, error) {
	var schema AccountSchema
	if err := r.db.WithContext(ctx).Where("username = ?", handle).First(&schema).Error; err != nil {
		return nil, err
	}
	return schema.ToEntity(), nil
}

func (r *AccountRepositoryImpl) ListAll(ctx context.Context, includeTrash bool) ([]*entities.Account, error) {
	var schemas []AccountSchema
	q := r.db.WithContext(ctx)
	if !includeTrash {
		q = q.Where("is_trash = ?", false)
	}
	if err := q.Find(&schemas).Error; err != nil {
		return nil, err
	}
	
	result := make([]*entities.Account, len(schemas))
	for i, s := range schemas {
		schemaCopy := s
		result[i] = schemaCopy.ToEntity()
	}
	return result, nil
}
