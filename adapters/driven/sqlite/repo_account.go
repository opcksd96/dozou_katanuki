package sqlite

import (
	"context"
	"dozou_katanuki/domain/account"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepositoryImpl(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db: db}
}

func toDomain(schema *AccountSchema) *account.Account {
	if schema == nil {
		return nil
	}
	return &account.Account{
		NumericID:    schema.NumericID,
		Username:     schema.Username,
		DisplayName:  schema.DisplayName,
		AvatarURL:    schema.AvatarURL,
		AvatarBase64: schema.AvatarBase64,
		Description:  schema.Description,
		GroupName:    schema.GroupName,
		AliasOf:      schema.AliasOf,
		IsWhitelist:  schema.IsWhitelist,
		IsTrash:      schema.IsTrash,
		TrashedBy:    schema.TrashedBy,
		TrashReason:  schema.TrashReason,
		TrashedAt:    schema.TrashedAt,
		UpdatedAt:    schema.UpdatedAt,
	}
}

func toSchema(entity *account.Account) *AccountSchema {
	if entity == nil {
		return nil
	}
	return &AccountSchema{
		NumericID:    entity.NumericID,
		Username:     entity.Username,
		DisplayName:  entity.DisplayName,
		AvatarURL:    entity.AvatarURL,
		AvatarBase64: entity.AvatarBase64,
		Description:  entity.Description,
		GroupName:    entity.GroupName,
		AliasOf:      entity.AliasOf,
		IsWhitelist:  entity.IsWhitelist,
		IsTrash:      entity.IsTrash,
		TrashedBy:    entity.TrashedBy,
		TrashReason:  entity.TrashReason,
		TrashedAt:    entity.TrashedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func (r *AccountRepositoryImpl) Save(ctx context.Context, acc *account.Account) error {
	schema := toSchema(acc)
	return r.db.WithContext(ctx).Save(schema).Error
}

func (r *AccountRepositoryImpl) GetByID(ctx context.Context, numericID string) (*account.Account, error) {
	var schema AccountSchema
	if err := r.db.WithContext(ctx).Where("numeric_id = ?", numericID).First(&schema).Error; err != nil {
		return nil, err
	}
	return toDomain(&schema), nil
}

func (r *AccountRepositoryImpl) GetByUsername(ctx context.Context, username string) (*account.Account, error) {
	var schema AccountSchema
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&schema).Error; err != nil {
		return nil, err
	}
	return toDomain(&schema), nil
}

func (r *AccountRepositoryImpl) ListAll(ctx context.Context, includeTrash bool) ([]*account.Account, error) {
	var schemas []AccountSchema
	q := r.db.WithContext(ctx)
	if !includeTrash {
		q = q.Where("is_trash = ?", false)
	}
	if err := q.Find(&schemas).Error; err != nil {
		return nil, err
	}
	
	result := make([]*account.Account, len(schemas))
	for i, s := range schemas {
		schemaCopy := s
		result[i] = toDomain(&schemaCopy)
	}
	return result, nil
}

func (r *AccountRepositoryImpl) Delete(ctx context.Context, numericID string) error {
	return r.db.WithContext(ctx).Where("numeric_id = ?", numericID).Delete(&AccountSchema{}).Error
}
